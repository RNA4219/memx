from __future__ import annotations

import json
import re
from dataclasses import replace
from datetime import datetime, timezone
from pathlib import Path
import sys
from types import SimpleNamespace
from typing import Iterable, Mapping, Sequence

import pytest

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from tools.codemap import update


def test_ensure_python_version_exits(monkeypatch, capsys):
    monkeypatch.setattr(update, "sys", SimpleNamespace(version_info=(3, 10, 0)))

    with pytest.raises(SystemExit) as excinfo:
        update.ensure_python_version()

    assert excinfo.value.code == 1
    captured = capsys.readouterr()
    assert "Python 3.11 or newer is required" in captured.out


def test_since_command_resolves_capsules_for_non_birdseye_diff(monkeypatch):
    options = update.parse_args(["--since", "--emit", "index+caps"])

    def fake_run(args, capture_output, text, check, cwd):
        assert args == [
            "git",
            "diff",
            "--name-status",
            "--find-renames",
            "--find-copies",
            "main...HEAD",
        ]
        assert cwd == update._REPO_ROOT
        return SimpleNamespace(stdout="M\tREADME.md\n")

    monkeypatch.setattr(update.subprocess, "run", fake_run)
    options = replace(options, diff_resolver=update.GitDiffResolver())

    resolved = options.resolve_targets()

    expected = update._REPO_ROOT / "docs" / "birdseye" / "caps" / "README.md.json"
    assert expected in resolved


def test_since_command_falls_back_to_default_targets(monkeypatch):
    options = update.parse_args(["--since"])

    options = replace(options, diff_resolver=SimpleNamespace(resolve=lambda _: ()))

    resolved = options.resolve_targets()

    expected_index = (update._REPO_ROOT / "docs" / "birdseye" / "index.json").resolve()
    expected_hot = (update._REPO_ROOT / "docs" / "birdseye" / "hot.json").resolve()
    expected_caps = (update._REPO_ROOT / "docs" / "birdseye" / "caps").resolve()

    assert resolved == (expected_index, expected_hot, expected_caps)


def test_git_diff_parser_handles_complex_statuses(tmp_path: Path):
    repo_root = tmp_path / "repo"
    caps_dir = repo_root / "docs" / "birdseye" / "caps"
    caps_dir.mkdir(parents=True, exist_ok=True)
    for slug in ("README.md", "RUNBOOK.md", "GUIDE.md", "GUIDE_COPY.md"):
        (caps_dir / f"{slug}.json").write_text("{}", encoding="utf-8")

    parser = update.GitDiffParser(repo_root=repo_root)

    payload = "\n".join(
        (
            "R100\tREADME.md\tRUNBOOK.md",
            "C100\tGUIDE.md\tGUIDE_COPY.md",
            "D\tdocs/birdseye/index.json",
        )
    )

    resolved = parser.parse(payload)

    assert resolved == (
        Path("docs/birdseye/caps/README.md.json"),
        Path("docs/birdseye/caps/RUNBOOK.md.json"),
        Path("docs/birdseye/caps/GUIDE.md.json"),
        Path("docs/birdseye/caps/GUIDE_COPY.md.json"),
        Path("docs/birdseye/index.json"),
    )


def test_git_diff_resolver_parses_type_change(monkeypatch):
    parser = update.GitDiffParser()

    def fake_runner(args):
        assert args[-1] == "main...HEAD"
        return SimpleNamespace(stdout="T\tREADME.md\n")

    resolver = update.GitDiffResolver(parser=parser, runner=fake_runner)
    resolved = resolver.resolve("main")

    assert resolved == (Path("docs/birdseye/caps/README.md.json"),)


def test_git_diff_resolver_parses_unmerged_status(monkeypatch):
    parser = update.GitDiffParser()

    def fake_runner(args):
        return SimpleNamespace(stdout="U\tREADME.md\n")

    resolver = update.GitDiffResolver(parser=parser, runner=fake_runner)
    resolved = resolver.resolve("main")

    assert resolved == (Path("docs/birdseye/caps/README.md.json"),)


def test_git_diff_resolver_uses_repo_root():
    captured: dict[str, object] = {}

    def fake_runner(args):
        captured["args"] = args
        return SimpleNamespace(stdout="")

    resolver = update.GitDiffResolver(parser=update.GitDiffParser(), runner=fake_runner)
    resolver.resolve("feature")

    assert captured["args"] == [
        "git",
        "diff",
        "--name-status",
        "--find-renames",
        "--find-copies",
        "feature...HEAD",
    ]


def test_git_diff_resolver_uses_user_provided_range():
    captured: list[list[str]] = []

    def fake_runner(args):
        captured.append(args)
        return SimpleNamespace(stdout="")

    resolver = update.GitDiffResolver(parser=update.GitDiffParser(), runner=fake_runner)
    resolver.resolve("origin/main..HEAD")
    resolver.resolve("feature...staging")

    assert captured == [
        [
            "git",
            "diff",
            "--name-status",
            "--find-renames",
            "--find-copies",
            "origin/main..HEAD",
        ],
        [
            "git",
            "diff",
            "--name-status",
            "--find-renames",
            "--find-copies",
            "feature...staging",
        ],
    ]


def test_derive_targets_from_since_accepts_absolute_paths():
    repo_root = Path(__file__).resolve().parents[1]
    absolute_readme = (repo_root / "README.md").resolve()

    derived = update._derive_targets_from_since((absolute_readme,), repo_root=repo_root)

    assert derived == (Path("docs/birdseye/caps/README.md.json"),)


def test_derive_targets_from_since_handles_rename_notation():
    derived = update._derive_targets_from_since(("README.md -> RUNBOOK.md",))

    assert derived == (
        Path("docs/birdseye/caps/README.md.json"),
        Path("docs/birdseye/caps/RUNBOOK.md.json"),
    )


def test_derive_targets_from_since_deduplicates_old_and_new_paths():
    derived = update._derive_targets_from_since(("README.md -> README.md",))

    assert derived == (Path("docs/birdseye/caps/README.md.json"),)


def _write_json(path, payload):
    path.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def _caps_payload(cap_id: str, *, deps_out=None, deps_in=None) -> dict[str, object]:
    return {
        "id": cap_id,
        "role": "doc",
        "public_api": [],
        "summary": cap_id,
        "deps_out": list(deps_out or []),
        "deps_in": list(deps_in or []),
        "risks": [],
        "tests": [],
    }


def _prepare_birdseye_root(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch, *, include_hot: bool = True
) -> dict[str, Path]:
    repo_root = tmp_path / "repo"
    repo_root.mkdir(parents=True, exist_ok=True)
    monkeypatch.setattr(update, "_REPO_ROOT", repo_root)

    root = repo_root / "docs" / "birdseye"
    caps_dir = root / "caps"
    caps_dir.mkdir(parents=True, exist_ok=True)

    index_path = root / "index.json"
    hot_path = root / "hot.json"
    capsule_path = caps_dir / "README.md.json"

    index_payload = {
        "generated_at": "00001",
        "nodes": {
            "README.md": {
                "caps": "docs/birdseye/caps/README.md.json",
            }
        },
        "edges": [],
    }
    _write_json(index_path, index_payload)

    if include_hot:
        _write_json(hot_path, {"generated_at": "00001"})
    elif hot_path.exists():
        hot_path.unlink()

    capsule_payload = _caps_payload("README.md")
    capsule_payload["generated_at"] = "00001"
    _write_json(capsule_path, capsule_payload)

    return {
        "root": root,
        "index": index_path,
        "hot": hot_path,
        "capsule": capsule_path,
    }


def _is_iso8601_utc(value: str) -> bool:
    return bool(re.fullmatch(r"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z", value))


@pytest.mark.parametrize(
    ("emit", "expected_loads", "expected_writes"),
    (
        (
            "index+caps",
            ("index", "hot", "caps"),
            {"index.json", "hot.json", "README.md.json"},
        ),
        ("index", ("index", "hot"), {"index.json", "hot.json"}),
        ("caps", ("index", "hot", "caps"), {"README.md.json"}),
    ),
)
def test_birdseye_root_plan_records_loads_and_writes(
    tmp_path: Path,
    monkeypatch: pytest.MonkeyPatch,
    emit: str,
    expected_loads: tuple[str, ...],
    expected_writes: set[str],
) -> None:
    fixture = _prepare_birdseye_root(tmp_path, monkeypatch)
    emit_index = emit in {"index", "index+caps"}
    emit_caps = emit in {"caps", "index+caps"}
    builder = update.BirdseyeRootBuilder(
        root=fixture["root"],
        root_targets=(fixture["index"],),
        emit_index=emit_index,
        emit_caps=emit_caps,
        timestamp="2025-01-01T00:00:00Z",
        serial_allocator=update._SerialAllocator(),
    )
    plan = builder.build()

    assert plan.loads == expected_loads
    assert {write.path.name for write in plan.writes} == expected_writes


def test_birdseye_root_plan_requires_hot_when_emitting_index(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    fixture = _prepare_birdseye_root(tmp_path, monkeypatch, include_hot=False)
    builder = update.BirdseyeRootBuilder(
        root=fixture["root"],
        root_targets=(fixture["index"],),
        emit_index=True,
        emit_caps=False,
        timestamp="2025-01-01T00:00:00Z",
        serial_allocator=update._SerialAllocator(),
    )

    with pytest.raises(FileNotFoundError):
        builder.build()


def test_next_generated_at_uses_fallback_timestamp():
    allocator = update._SerialAllocator()

    generated = update._next_generated_at(None, "2025-01-01T00:00:00Z", allocator=allocator)
    assert generated == "2025-01-01T00:00:00Z"

    follow_up = update._next_generated_at(
        "2024-12-31T23:59:59Z",
        "2025-01-01T00:00:00Z",
        allocator=allocator,
    )
    assert follow_up == "2025-01-01T00:00:00Z"

    seeded_allocator = update._SerialAllocator()
    seeded_allocator.observe("00041")

    advanced = update._next_generated_at(
        "unexpected-value",
        "2025-01-01T00:00:00Z",
        allocator=seeded_allocator,
    )
    assert advanced == "2025-01-01T00:00:00Z"


_HOT_INDEX_SNAPSHOT = "docs/birdseye/index.json"
_HOT_REFRESH_COMMAND = (
    "python tools/codemap/update.py --targets docs/birdseye/index.json,docs/birdseye/hot.json --emit index+caps"
)
_HOT_CURATION_NOTES = "Birdseye ホットノードのサンプルノート"


_HOT_HOTLIST_METADATA = {
    "index_snapshot": _HOT_INDEX_SNAPSHOT,
    "refresh_command": _HOT_REFRESH_COMMAND,
    "curation_notes": _HOT_CURATION_NOTES,
}


_HOT_NODES_FIXTURE: Sequence[dict[str, object]] = (
    {
        "id": "README.md",
        "role": "bootstrap",
        "reason": "Birdseye エントリーポイントのサンプル",
        "caps": "docs/birdseye/caps/README.md.json",
        "edges": ["GUARDRAILS.md", "docs/birdseye/index.json"],
        "last_verified_at": "2024-01-01T00:00:00Z",
        "refresh_command": _HOT_REFRESH_COMMAND,
        "index_snapshot": _HOT_INDEX_SNAPSHOT,
        "curation_notes": _HOT_CURATION_NOTES,
    },
    {
        "id": "GUARDRAILS.md",
        "role": "policy",
        "reason": "運用ガイドラインのサンプル",
        "caps": "docs/birdseye/caps/GUARDRAILS.md.json",
        "edges": ["README.md", "docs/birdseye/index.json"],
        "last_verified_at": "2024-01-01T00:00:00Z",
        "refresh_command": _HOT_REFRESH_COMMAND,
        "index_snapshot": _HOT_INDEX_SNAPSHOT,
        "curation_notes": _HOT_CURATION_NOTES,
    },
    {
        "id": "docs/birdseye/index.json",
        "role": "index",
        "reason": "Birdseye インデックスのホットビュー",
        "caps": "docs/birdseye/caps/docs.birdseye.index.json.json",
        "edges": ["README.md", "GUARDRAILS.md"],
        "last_verified_at": "2024-01-01T00:00:00Z",
        "refresh_command": _HOT_REFRESH_COMMAND,
        "index_snapshot": _HOT_INDEX_SNAPSHOT,
        "curation_notes": _HOT_CURATION_NOTES,
    },
)


_HOT_NODE_IDS: tuple[str, ...] = tuple(node["id"] for node in _HOT_NODES_FIXTURE)


def test_hot_refresh_command_matches_documentation():
    repo_root = Path(__file__).resolve().parents[1]
    hot_doc = json.loads((repo_root / "docs" / "birdseye" / "hot.json").read_text(encoding="utf-8"))

    assert hot_doc["refresh_command"] == _HOT_REFRESH_COMMAND
    for node in hot_doc.get("nodes", []):
        if "refresh_command" in node:
            assert node["refresh_command"] == _HOT_REFRESH_COMMAND


def test_hot_nodes_caps_refer_to_existing_capsules():
    repo_root = Path(__file__).resolve().parents[1]
    hot_doc = json.loads((repo_root / "docs" / "birdseye" / "hot.json").read_text(encoding="utf-8"))

    birdseye_index_node = next(
        (node for node in hot_doc.get("nodes", []) if node.get("id") == "docs/birdseye/index.json"),
        None,
    )
    assert birdseye_index_node is not None
    assert (
        birdseye_index_node.get("caps") == "docs/birdseye/caps/docs.birdseye.index.json.json"
    )

    for node in hot_doc.get("nodes", []):
        caps_ref = node.get("caps")
        if caps_ref is None:
            continue

        assert isinstance(caps_ref, str), f"caps must be str or null, got {type(caps_ref)!r}"

        caps_path = repo_root / caps_ref
        assert caps_path.is_file(), f"Capsule missing for node {node['id']}: {caps_ref}"

        capsule_payload = json.loads(caps_path.read_text(encoding="utf-8"))
        assert (
            capsule_payload.get("id") == node["id"]
        ), f"Capsule id mismatch for {node['id']}: {capsule_payload.get('id')}"


def test_index_contains_docs_birdseye_node_with_bidirectional_edges():
    repo_root = Path(__file__).resolve().parents[1]
    index_doc = json.loads((repo_root / "docs" / "birdseye" / "index.json").read_text(encoding="utf-8"))

    nodes = index_doc.get("nodes", {})
    assert "docs/BIRDSEYE.md" in nodes

    capsule_path = nodes["docs/BIRDSEYE.md"].get("caps")
    assert capsule_path == "docs/birdseye/caps/docs.BIRDSEYE.md.json"

    edges = {
        tuple(edge)
        for edge in index_doc.get("edges", [])
        if isinstance(edge, list) and len(edge) == 2
    }

    neighbours = ("README.md", "GUARDRAILS.md", "tools/codemap/README.md")
    for neighbour in neighbours:
        assert ("docs/BIRDSEYE.md", neighbour) in edges
        assert (neighbour, "docs/BIRDSEYE.md") in edges


def test_index_mtime_is_iso8601_utc():
    repo_root = Path(__file__).resolve().parents[1]
    index_doc = json.loads((repo_root / "docs" / "birdseye" / "index.json").read_text(encoding="utf-8"))

    mtimes = [node.get("mtime") for node in index_doc.get("nodes", {}).values()]

    assert mtimes, "mtime entries should exist"
    assert all(isinstance(mtime, str) and _is_iso8601_utc(mtime) for mtime in mtimes)


def _prepare_birdseye(
    tmp_path: Path,
    *,
    edges: Iterable[Iterable[str]],
    caps_payloads: Mapping[str, Mapping[str, object]],
    hot_entries: Sequence[str] | None,
    root: Path | None = None,
) -> tuple[Path, Path, Path, dict[str, Path]]:
    root = Path(root) if root is not None else tmp_path / "birdseye"
    caps_dir = root / "caps"
    caps_dir.mkdir(parents=True)
    nodes = {
        cap_id: {"role": payload.get("role", "doc"), "caps": f"docs/birdseye/caps/{cap_id}.json"}
        for cap_id, payload in caps_payloads.items()
    }
    edge_pairs = [list(pair) for pair in edges]
    index_path = root / "index.json"
    _write_json(index_path, {"generated_at": "00001", "nodes": nodes, "edges": edge_pairs})
    cap_paths = {}
    for cap_id, payload in caps_payloads.items():
        cap_path = caps_dir / f"{cap_id}.json"
        cap_paths[cap_id] = cap_path
        _write_json(cap_path, payload)
    hot_path = root / "hot.json"
    defaults_lookup = {entry["id"]: entry for entry in _HOT_NODES_FIXTURE}
    selected_ids = list(hot_entries) if hot_entries is not None else list(_HOT_NODE_IDS)
    neighbours: dict[str, set[str]] = {}
    for pair in edge_pairs:
        if len(pair) != 2:
            continue
        source, destination = pair
        if not isinstance(source, str) or not isinstance(destination, str):
            continue
        neighbours.setdefault(source, set()).add(destination)
        neighbours.setdefault(destination, set()).add(source)
    serialized_nodes = []
    for entry_id in selected_ids:
        defaults = defaults_lookup.get(entry_id, {})
        index_node = nodes.get(entry_id, {})
        resolved_edges = (
            list(defaults["edges"]) if "edges" in defaults else sorted(neighbours.get(entry_id, set()))
        )
        if "caps" in defaults:
            resolved_caps = defaults["caps"]
        else:
            resolved_caps = index_node.get("caps") or f"docs/birdseye/caps/{entry_id}.json"
        payload = {
            "id": entry_id,
            "role": defaults.get("role", index_node.get("role", "doc")),
            "edges": resolved_edges,
            "caps": resolved_caps,
            "last_verified_at": defaults.get("last_verified_at", "2024-01-01T00:00:00Z"),
        }
        if "reason" in defaults:
            payload["reason"] = defaults["reason"]
        if "refresh_command" in defaults:
            payload["refresh_command"] = defaults["refresh_command"]
        if "index_snapshot" in defaults:
            payload["index_snapshot"] = defaults["index_snapshot"]
        else:
            payload["index_snapshot"] = _HOT_INDEX_SNAPSHOT
        if "curation_notes" in defaults:
            payload["curation_notes"] = defaults["curation_notes"]
        else:
            payload["curation_notes"] = _HOT_CURATION_NOTES
        serialized_nodes.append(payload)
    _write_json(
        hot_path,
        {
            "generated_at": "00001",
            **_HOT_HOTLIST_METADATA,
            "nodes": serialized_nodes,
        },
    )
    return root, index_path, hot_path, cap_paths


def test_birdseye_session_plan_targets_root_and_capsule(tmp_path, monkeypatch):
    monkeypatch.chdir(tmp_path)
    monkeypatch.setattr(update, "_REPO_ROOT", tmp_path)

    caps_payloads = {
        "alpha.md": _caps_payload("alpha.md", deps_out=["beta.md"], deps_in=[]),
        "beta.md": _caps_payload("beta.md", deps_out=["gamma.md"], deps_in=["alpha.md"]),
        "gamma.md": _caps_payload("gamma.md", deps_out=["delta.md"], deps_in=["beta.md"]),
        "delta.md": _caps_payload("delta.md", deps_out=[], deps_in=["gamma.md"]),
    }
    _root, index_path, hot_path, cap_paths = _prepare_birdseye(
        tmp_path,
        root=tmp_path / "docs" / "birdseye",
        edges=[
            ["alpha.md", "beta.md"],
            ["beta.md", "alpha.md"],
            ["beta.md", "gamma.md"],
            ["gamma.md", "beta.md"],
            ["gamma.md", "delta.md"],
            ["delta.md", "gamma.md"],
        ],
        caps_payloads=caps_payloads,
        hot_entries=("alpha.md", "beta.md", "gamma.md", "delta.md"),
    )

    options = update.UpdateOptions(
        targets=(index_path, cap_paths["beta.md"]),
        emit="index+caps",
        dry_run=True,
    )

    session = update.BirdseyeUpdateSession(options=options, timestamp="2025-01-01T00:00:00Z")
    plan = session.plan()

    planned_paths = {write.path for write in plan.writes}

    assert plan.generated_at == "2025-01-01T00:00:00Z"
    assert planned_paths == {
        index_path,
        hot_path,
        cap_paths["alpha.md"],
        cap_paths["beta.md"],
        cap_paths["gamma.md"],
        cap_paths["delta.md"],
    }


def test_birdseye_session_plan_computes_expected_writes(tmp_path, monkeypatch):
    monkeypatch.chdir(tmp_path)
    monkeypatch.setattr(update, "_REPO_ROOT", tmp_path)

    caps_payloads = {
        "alpha.md": _caps_payload("alpha.md", deps_out=["beta.md"], deps_in=[]),
        "beta.md": _caps_payload("beta.md", deps_out=[], deps_in=["alpha.md", "gamma.md"]),
        "gamma.md": _caps_payload("gamma.md", deps_out=["delta.md"], deps_in=["beta.md"]),
        "delta.md": _caps_payload("delta.md", deps_out=[], deps_in=["gamma.md"]),
    }
    root, index_path, hot_path, cap_paths = _prepare_birdseye(
        tmp_path,
        root=tmp_path / "docs" / "birdseye",
        edges=[
            ["alpha.md", "beta.md"],
            ["beta.md", "alpha.md"],
            ["beta.md", "gamma.md"],
            ["gamma.md", "beta.md"],
            ["gamma.md", "delta.md"],
            ["delta.md", "gamma.md"],
        ],
        caps_payloads=caps_payloads,
        hot_entries=("alpha.md", "beta.md", "gamma.md", "delta.md"),
    )

    options = update.UpdateOptions(
        targets=(cap_paths["alpha.md"],),
        emit="index+caps",
        dry_run=True,
    )

    session = update.BirdseyeUpdateSession(options=options, timestamp="2025-01-01T00:00:00Z")
    plan = session.plan()

    planned_paths = {write.path for write in plan.writes}

    assert plan.generated_at == "2025-01-01T00:00:00Z"
    assert index_path in planned_paths
    assert hot_path in planned_paths
    assert cap_paths["alpha.md"] in planned_paths
    assert cap_paths["beta.md"] in planned_paths
    assert cap_paths["gamma.md"] in planned_paths
    assert cap_paths["delta.md"] not in planned_paths


def test_run_update_with_since_updates_capsules_within_two_hops(tmp_path, monkeypatch):
    monkeypatch.chdir(tmp_path)
    monkeypatch.setattr(update, "_REPO_ROOT", tmp_path)

    caps_payloads = {
        "alpha.md": _caps_payload("alpha.md", deps_out=["stale"], deps_in=["stale"]),
        "beta.md": _caps_payload("beta.md", deps_out=["stale"], deps_in=["stale"]),
        "gamma.md": _caps_payload("gamma.md", deps_out=["stale"], deps_in=["stale"]),
        "delta.md": _caps_payload("delta.md", deps_out=["stale"], deps_in=["stale"]),
    }
    _, _, _, cap_paths = _prepare_birdseye(
        tmp_path,
        edges=[
            ["alpha.md", "beta.md"],
            ["beta.md", "alpha.md"],
            ["beta.md", "gamma.md"],
            ["gamma.md", "beta.md"],
            ["gamma.md", "delta.md"],
            ["delta.md", "gamma.md"],
        ],
        caps_payloads=caps_payloads,
        hot_entries=_HOT_NODE_IDS,
        root=tmp_path / "docs" / "birdseye",
    )

    options = update.UpdateOptions(
        targets=(),
        emit="caps",
        dry_run=False,
        since="main",
        diff_resolver=SimpleNamespace(
            resolve=lambda _: (Path("docs/birdseye/caps/alpha.md.json"),)
        ),
    )

    report = update.run_update(options)

    expected_caps = {cap_paths[name] for name in ("alpha.md", "beta.md", "gamma.md")}

    assert set(report.planned_writes) == expected_caps
    assert set(report.performed_writes) == expected_caps

    delta_payload = json.loads(cap_paths["delta.md"].read_text(encoding="utf-8"))
    assert delta_payload["deps_out"] == ["stale"]
    assert delta_payload["deps_in"] == ["stale"]


def test_run_update_with_since_handles_git_rename(tmp_path, monkeypatch):
    monkeypatch.chdir(tmp_path)
    monkeypatch.setattr(update, "_REPO_ROOT", tmp_path)

    caps_payloads = {
        "README.md": _caps_payload("README.md"),
        "RUNBOOK.md": _caps_payload("RUNBOOK.md"),
    }
    _, _, _, cap_paths = _prepare_birdseye(
        tmp_path,
        edges=[["README.md", "RUNBOOK.md"], ["RUNBOOK.md", "README.md"]],
        caps_payloads=caps_payloads,
        hot_entries=("README.md", "RUNBOOK.md"),
        root=tmp_path / "docs" / "birdseye",
    )

    def fake_run(args, capture_output, text, check, cwd):
        assert args == [
            "git",
            "diff",
            "--name-status",
            "--find-renames",
            "--find-copies",
            "main...HEAD",
        ]
        assert cwd == tmp_path
        return SimpleNamespace(stdout="R100\tREADME.md\tRUNBOOK.md\n")

    monkeypatch.setattr(update.subprocess, "run", fake_run)

    options = update.UpdateOptions(
        targets=(),
        emit="caps",
        dry_run=False,
        since="main",
        diff_resolver=update.GitDiffResolver(),
    )

    report = update.run_update(options)

    expected_caps = {cap_paths[name] for name in ("README.md", "RUNBOOK.md")}

    assert set(report.planned_writes) == expected_caps
    assert set(report.performed_writes) == expected_caps


def test_run_update_recreates_missing_capsule(tmp_path, monkeypatch):
    monkeypatch.chdir(tmp_path)
    monkeypatch.setattr(update, "_REPO_ROOT", tmp_path)

    caps_payloads = {
        "README.md": _caps_payload("README.md"),
        "RUNBOOK.md": _caps_payload("RUNBOOK.md"),
    }
    _, _, _, cap_paths = _prepare_birdseye(
        tmp_path,
        edges=[["README.md", "RUNBOOK.md"], ["RUNBOOK.md", "README.md"]],
        caps_payloads=caps_payloads,
        hot_entries=("README.md", "RUNBOOK.md"),
        root=tmp_path / "docs" / "birdseye",
    )

    missing_path = cap_paths["README.md"]
    missing_path.unlink()
    assert not missing_path.exists()

    options = update.parse_args(["--since"])
    diff_targets = (Path("docs/birdseye/caps/README.md.json"),)
    options = replace(
        options,
        diff_resolver=SimpleNamespace(resolve=lambda _: diff_targets),
    )

    report = update.run_update(options)

    assert missing_path.exists()
    assert missing_path in report.performed_writes

    payload = json.loads(missing_path.read_text(encoding="utf-8"))
    assert payload["id"] == "README.md"
    assert payload["deps_out"] == ["RUNBOOK.md"]
    assert payload["deps_in"] == ["RUNBOOK.md"]


def test_run_update_resolves_targets_from_repo_root_outside_cwd(tmp_path, monkeypatch):
    repo_root = tmp_path / "repo"
    outside_dir = tmp_path / "outside"
    outside_dir.mkdir()

    caps_payloads = {"README.md": _caps_payload("README.md")}
    _root, index_path, hot_path, _ = _prepare_birdseye(
        repo_root,
        edges=[],
        caps_payloads=caps_payloads,
        hot_entries=("README.md",),
        root=repo_root / "docs" / "birdseye",
    )

    monkeypatch.setattr(update, "_REPO_ROOT", repo_root)
    monkeypatch.chdir(outside_dir)

    options = update.parse_args(["--targets", "docs/birdseye/index.json", "--emit", "index"])

    assert options.targets == (index_path,)
    resolved = options.resolve_targets()
    assert resolved == (index_path,)

    frozen_now = datetime(2025, 1, 2, tzinfo=timezone.utc)
    monkeypatch.setattr(update, "utc_now", lambda: frozen_now)

    report = update.run_update(options)

    expected_writes = {index_path, hot_path}
    assert set(report.planned_writes) == expected_writes
    assert set(report.performed_writes) == expected_writes

    refreshed_index = json.loads(index_path.read_text(encoding="utf-8"))
    refreshed_hot = json.loads(hot_path.read_text(encoding="utf-8"))

    assert refreshed_index["generated_at"] == "2025-01-02T00:00:00Z"
    assert refreshed_hot["generated_at"] == "2025-01-02T00:00:00Z"


@pytest.mark.parametrize("dry_run", [False, True])
def test_run_update_refreshes_metadata_and_dependencies(tmp_path, monkeypatch, dry_run):
    caps_payloads = {
        "alpha.md": _caps_payload("alpha.md", deps_in=["obsolete"]),
        "beta.md": _caps_payload("beta.md", deps_out=["stale"]),
    }
    root, index_path, hot_path, cap_paths = _prepare_birdseye(
        tmp_path,
        edges=[["alpha.md", "beta.md"], ["beta.md", "alpha.md"]],
        caps_payloads=caps_payloads,
        hot_entries=_HOT_NODE_IDS,
    )

    frozen_now = datetime(2025, 1, 1, 9, 30, tzinfo=timezone.utc)
    expected_timestamp = update._format_timestamp(frozen_now)
    monkeypatch.setattr(update, "utc_now", lambda: frozen_now)

    snapshots = None
    if dry_run:
        snapshots = {
            path: path.read_text(encoding="utf-8")
            for path in (index_path, hot_path, *cap_paths.values())
        }

    report = update.run_update(
        update.UpdateOptions(targets=(root,), emit="index+caps", dry_run=dry_run)
    )

    assert _is_iso8601_utc(report.generated_at)
    assert report.generated_at == expected_timestamp
    assert set(report.planned_writes) == {index_path, hot_path, *cap_paths.values()}

    if dry_run:
        assert report.performed_writes == ()
        assert snapshots is not None
        for path, before in snapshots.items():
            assert path.read_text(encoding="utf-8") == before
        return

    assert set(report.performed_writes) == {index_path, hot_path, *cap_paths.values()}

    refreshed_index = json.loads(index_path.read_text(encoding="utf-8"))
    assert refreshed_index["generated_at"] == expected_timestamp

    refreshed_alpha = json.loads(cap_paths["alpha.md"].read_text(encoding="utf-8"))
    assert refreshed_alpha["deps_out"] == ["beta.md"]
    assert refreshed_alpha["deps_in"] == ["beta.md"]

    refreshed_beta = json.loads(cap_paths["beta.md"].read_text(encoding="utf-8"))
    assert refreshed_beta["deps_out"] == ["alpha.md"]
    assert refreshed_beta["deps_in"] == ["alpha.md"]

    refreshed_hot = json.loads(hot_path.read_text(encoding="utf-8"))
    assert refreshed_hot["generated_at"] == expected_timestamp
    expected_hot_nodes = [dict(node) for node in _HOT_NODES_FIXTURE]
    assert refreshed_hot["index_snapshot"] == _HOT_INDEX_SNAPSHOT
    assert refreshed_hot["refresh_command"] == _HOT_REFRESH_COMMAND
    assert refreshed_hot["curation_notes"] == _HOT_CURATION_NOTES
    assert refreshed_hot["nodes"] == expected_hot_nodes
    for node in refreshed_hot["nodes"]:
        assert node["refresh_command"] == _HOT_REFRESH_COMMAND
        assert node["index_snapshot"] == _HOT_INDEX_SNAPSHOT
        assert node["curation_notes"] == _HOT_CURATION_NOTES
    assert len(refreshed_hot["nodes"]) == len(expected_hot_nodes)
    assert refreshed_hot["nodes"][0]["edges"] == expected_hot_nodes[0]["edges"]
    assert refreshed_hot["nodes"][0]["caps"] == expected_hot_nodes[0]["caps"]
    assert refreshed_hot["nodes"][0]["role"] == expected_hot_nodes[0]["role"]
    assert (
        refreshed_hot["nodes"][0]["last_verified_at"]
        == expected_hot_nodes[0]["last_verified_at"]
    )
    hot_caps_lookup = {node["id"]: node["caps"] for node in refreshed_hot["nodes"]}
    assert (
        hot_caps_lookup["docs/birdseye/index.json"]
        == "docs/birdseye/caps/docs.birdseye.index.json.json"
    )


def test_run_update_requires_hot_json_when_emitting_index(tmp_path, monkeypatch):
    caps_payloads = {
        "alpha.md": _caps_payload("alpha.md"),
        "beta.md": _caps_payload("beta.md"),
    }
    root, _, hot_path, _ = _prepare_birdseye(
        tmp_path,
        edges=[["alpha.md", "beta.md"]],
        caps_payloads=caps_payloads,
        hot_entries=_HOT_NODE_IDS,
    )

    hot_path.unlink()

    monkeypatch.setattr(update, "_REPO_ROOT", tmp_path)

    with pytest.raises(FileNotFoundError) as excinfo:
        update.run_update(update.UpdateOptions(targets=(root,), emit="index+caps"))

    message = str(excinfo.value)
    assert str(hot_path) in message
    assert _HOT_REFRESH_COMMAND in message


def test_run_update_preserves_hot_nodes_structure(tmp_path, monkeypatch):
    caps_payloads = {
        "alpha.md": _caps_payload("alpha.md", deps_in=["obsolete"]),
        "beta.md": _caps_payload("beta.md", deps_out=["stale"]),
    }
    root, _, hot_path, _ = _prepare_birdseye(
        tmp_path,
        edges=[["alpha.md", "beta.md"], ["beta.md", "alpha.md"]],
        caps_payloads=caps_payloads,
        hot_entries=_HOT_NODE_IDS,
    )

    baseline_hot = json.loads(hot_path.read_text(encoding="utf-8"))
    assert baseline_hot["index_snapshot"] == _HOT_INDEX_SNAPSHOT
    assert baseline_hot["refresh_command"] == _HOT_REFRESH_COMMAND
    assert baseline_hot["curation_notes"] == _HOT_CURATION_NOTES
    assert (
        {
            node["id"]: node["caps"]
            for node in baseline_hot["nodes"]
            if "caps" in node and "id" in node
        }["docs/birdseye/index.json"]
        == "docs/birdseye/caps/docs.birdseye.index.json.json"
    )
    for node in baseline_hot["nodes"]:
        assert node["refresh_command"] == _HOT_REFRESH_COMMAND
        assert node["index_snapshot"] == _HOT_INDEX_SNAPSHOT
        assert node["curation_notes"] == _HOT_CURATION_NOTES

    frozen_now = datetime(2025, 1, 1, 12, 0, tzinfo=timezone.utc)
    expected_timestamp = update._format_timestamp(frozen_now)
    monkeypatch.setattr(update, "utc_now", lambda: frozen_now)

    report = update.run_update(
        update.UpdateOptions(targets=(root,), emit="index+caps", dry_run=False)
    )

    assert _is_iso8601_utc(report.generated_at)
    assert report.generated_at == expected_timestamp

    refreshed_hot = json.loads(hot_path.read_text(encoding="utf-8"))
    assert refreshed_hot["generated_at"] == expected_timestamp
    assert refreshed_hot["index_snapshot"] == _HOT_INDEX_SNAPSHOT
    assert refreshed_hot["refresh_command"] == _HOT_REFRESH_COMMAND
    assert refreshed_hot["curation_notes"] == _HOT_CURATION_NOTES
    assert refreshed_hot["nodes"] == baseline_hot["nodes"]
    assert (
        {
            node["id"]: node["caps"]
            for node in refreshed_hot["nodes"]
            if "caps" in node and "id" in node
        }["docs/birdseye/index.json"]
        == "docs/birdseye/caps/docs.birdseye.index.json.json"
    )
    for node in refreshed_hot["nodes"]:
        assert node["refresh_command"] == _HOT_REFRESH_COMMAND
        assert node["index_snapshot"] == _HOT_INDEX_SNAPSHOT
        assert node["curation_notes"] == _HOT_CURATION_NOTES
    assert len(refreshed_hot["nodes"]) == len(_HOT_NODES_FIXTURE)
    assert all(
        node["last_verified_at"] == "2024-01-01T00:00:00Z"
        for node in refreshed_hot["nodes"]
    )
    assert refreshed_hot["nodes"][0]["edges"] == baseline_hot["nodes"][0]["edges"]
    assert refreshed_hot["nodes"][0]["caps"] == baseline_hot["nodes"][0]["caps"]
    assert refreshed_hot["nodes"][0]["role"] == baseline_hot["nodes"][0]["role"]


def test_run_update_limits_caps_to_two_hop_scope(tmp_path, monkeypatch):
    caps_payloads = {
        cap_id: _caps_payload(cap_id, deps_out=["stale"], deps_in=["old"])
        for cap_id in ("alpha.md", "beta.md", "gamma.md", "delta.md", "epsilon.md")
    }
    root, index_path, hot_path, cap_paths = _prepare_birdseye(
        tmp_path,
        edges=[
            ["alpha.md", "beta.md"],
            ["beta.md", "gamma.md"],
            ["gamma.md", "delta.md"],
            ["delta.md", "epsilon.md"],
        ],
        caps_payloads=caps_payloads,
        hot_entries=[],
    )

    baseline_hot = json.loads(hot_path.read_text(encoding="utf-8"))
    assert baseline_hot["nodes"] == []

    frozen_now = datetime(2025, 1, 2, tzinfo=timezone.utc)
    monkeypatch.setattr(update, "utc_now", lambda: frozen_now)

    report = update.run_update(
        update.UpdateOptions(targets=(cap_paths["beta.md"],), emit="caps", dry_run=False)
    )

    expected_caps = {
        cap_paths[cap_id]
        for cap_id in ("alpha.md", "beta.md", "gamma.md", "delta.md")
    }
    assert set(report.planned_writes) == expected_caps
    assert set(report.performed_writes) == expected_caps

    expected_deps = {
        "alpha.md": (["beta.md"], []),
        "beta.md": (["gamma.md"], ["alpha.md"]),
        "gamma.md": (["delta.md"], ["beta.md"]),
        "delta.md": (["epsilon.md"], ["gamma.md"]),
    }
    for cap_id, (deps_out, deps_in) in expected_deps.items():
        refreshed = json.loads(cap_paths[cap_id].read_text(encoding="utf-8"))
        assert refreshed["deps_out"] == deps_out
        assert refreshed["deps_in"] == deps_in


def test_run_update_refreshes_caps_generated_at(tmp_path, monkeypatch):
    edges = [
        ["alpha.md", "beta.md"],
        ["beta.md", "gamma.md"],
        ["gamma.md", "delta.md"],
        ["delta.md", "epsilon.md"],
    ]
    expected_deps = {
        "alpha.md": (["beta.md"], []),
        "beta.md": (["gamma.md"], ["alpha.md"]),
        "gamma.md": (["delta.md"], ["beta.md"]),
        "delta.md": (["epsilon.md"], ["gamma.md"]),
        "epsilon.md": ([], ["delta.md"]),
    }
    base_serial = "00009"
    caps_payloads = {}
    for cap_id, (deps_out, deps_in) in expected_deps.items():
        payload = _caps_payload(cap_id, deps_out=deps_out, deps_in=deps_in)
        payload["generated_at"] = base_serial
        caps_payloads[cap_id] = payload

    root, _, hot_path, cap_paths = _prepare_birdseye(
        tmp_path,
        edges=edges,
        caps_payloads=caps_payloads,
        hot_entries=[],
    )

    assert json.loads(hot_path.read_text(encoding="utf-8"))["nodes"] == []

    frozen_now = datetime(2025, 1, 5, tzinfo=timezone.utc)
    monkeypatch.setattr(update, "utc_now", lambda: frozen_now)

    expected_timestamp = update._format_timestamp(frozen_now)
    report = update.run_update(
        update.UpdateOptions(targets=(cap_paths["beta.md"],), emit="caps", dry_run=False)
    )

    expected_caps = {
        cap_paths[cap_id]
        for cap_id in ("alpha.md", "beta.md", "gamma.md", "delta.md")
    }

    assert _is_iso8601_utc(report.generated_at)
    assert report.generated_at == expected_timestamp
    assert set(report.planned_writes) == expected_caps
    assert set(report.performed_writes) == expected_caps

    for cap_id in ("alpha.md", "beta.md", "gamma.md", "delta.md"):
        refreshed = json.loads(cap_paths[cap_id].read_text(encoding="utf-8"))
        assert _is_iso8601_utc(refreshed["generated_at"])
        assert refreshed["generated_at"] == expected_timestamp
        expected_out, expected_in = expected_deps[cap_id]
        assert refreshed["deps_out"] == expected_out
        assert refreshed["deps_in"] == expected_in

    untouched = json.loads(cap_paths["epsilon.md"].read_text(encoding="utf-8"))
    assert untouched["generated_at"] == base_serial
    expected_out, expected_in = expected_deps["epsilon.md"]
    assert untouched["deps_out"] == expected_out
    assert untouched["deps_in"] == expected_in


def test_run_update_recovers_non_serial_generated_at(tmp_path, monkeypatch):
    monkeypatch.chdir(tmp_path)
    monkeypatch.setattr(update, "_REPO_ROOT", tmp_path)

    caps_payloads = {
        "alpha.md": {**_caps_payload("alpha.md"), "generated_at": "00007"},
        "beta.md": {**_caps_payload("beta.md"), "generated_at": "2024-01-01T00:00:00Z"},
    }

    _root, index_path, hot_path, cap_paths = _prepare_birdseye(
        tmp_path,
        edges=[["alpha.md", "beta.md"], ["beta.md", "alpha.md"]],
        caps_payloads=caps_payloads,
        hot_entries=("alpha.md", "beta.md"),
        root=tmp_path / "docs" / "birdseye",
    )

    index_payload = json.loads(index_path.read_text(encoding="utf-8"))
    index_payload["generated_at"] = None
    _write_json(index_path, index_payload)

    hot_payload = json.loads(hot_path.read_text(encoding="utf-8"))
    hot_payload["generated_at"] = "2023-12-31T23:59:59Z"
    _write_json(hot_path, hot_payload)

    frozen_now = datetime(2025, 1, 3, tzinfo=timezone.utc)
    monkeypatch.setattr(update, "utc_now", lambda: frozen_now)

    report = update.run_update(
        update.UpdateOptions(targets=(cap_paths["beta.md"],), emit="index+caps", dry_run=False)
    )

    expected_timestamp = update._format_timestamp(frozen_now)

    assert report.generated_at == expected_timestamp

    refreshed_index = json.loads(index_path.read_text(encoding="utf-8"))
    assert refreshed_index["generated_at"] == expected_timestamp

    refreshed_hot = json.loads(hot_path.read_text(encoding="utf-8"))
    assert refreshed_hot["generated_at"] == expected_timestamp

    refreshed_beta = json.loads(cap_paths["beta.md"].read_text(encoding="utf-8"))
    assert refreshed_beta["generated_at"] == expected_timestamp

    updated_alpha = json.loads(cap_paths["alpha.md"].read_text(encoding="utf-8"))
    assert updated_alpha["generated_at"] == expected_timestamp


def test_run_update_creates_missing_capsule_for_array_nodes(tmp_path, monkeypatch):
    root = tmp_path / "docs" / "birdseye"
    caps_dir = root / "caps"
    caps_dir.mkdir(parents=True)
    index_path = root / "index.json"
    hot_path = root / "hot.json"

    _write_json(
        index_path,
        {
            "generated_at": "2026-03-03T08:48:22Z",
            "nodes": [
                {
                    "node_id": "requirements",
                    "capsule": "docs/birdseye/caps/requirements.json",
                    "depends_on": [],
                },
                {
                    "node_id": "design",
                    "capsule": "docs/birdseye/caps/design.json",
                    "depends_on": ["requirements"],
                },
            ],
            "edges": [["requirements", "design"]],
        },
    )
    _write_json(
        hot_path,
        {
            "generated_at": "2026-03-03T08:48:22Z",
            **_HOT_HOTLIST_METADATA,
            "nodes": [],
        },
    )
    _write_json(caps_dir / "requirements.json", _caps_payload("requirements"))

    monkeypatch.setattr(update, "_REPO_ROOT", tmp_path)

    update.run_update(
        update.UpdateOptions(targets=(index_path,), emit="index+caps", dry_run=False)
    )

    generated_capsule = caps_dir / "design.json"
    assert generated_capsule.exists()
    generated_payload = json.loads(generated_capsule.read_text(encoding="utf-8"))
    assert generated_payload["id"] == "design"


def test_run_update_accepts_caps_directory_target(tmp_path, monkeypatch):
    caps_payloads = {
        cap_id: _caps_payload(cap_id)
        for cap_id in ("alpha.md", "beta.md", "gamma.md", "delta.md", "epsilon.md")
    }
    root_base = tmp_path / "docs" / "birdseye"
    root, index_path, hot_path, cap_paths = _prepare_birdseye(
        tmp_path,
        edges=[
            ["alpha.md", "beta.md"],
            ["beta.md", "gamma.md"],
            ["gamma.md", "delta.md"],
            ["delta.md", "epsilon.md"],
        ],
        caps_payloads=caps_payloads,
        hot_entries=_HOT_NODE_IDS,
        root=root_base,
    )

    monkeypatch.setattr(update, "_REPO_ROOT", tmp_path)
    monkeypatch.chdir(tmp_path)

    frozen_now = datetime(2025, 1, 4, tzinfo=timezone.utc)
    expected_timestamp = update._format_timestamp(frozen_now)
    monkeypatch.setattr(update, "utc_now", lambda: frozen_now)

    report = update.run_update(
        update.UpdateOptions(targets=(Path("docs/birdseye/caps"),), emit="index+caps", dry_run=False)
    )

    expected_caps = {cap_paths[cap_id] for cap_id in cap_paths}
    expected_writes = expected_caps | {index_path, hot_path}

    assert _is_iso8601_utc(report.generated_at)
    assert report.generated_at == expected_timestamp
    assert set(report.planned_writes) == expected_writes
    assert set(report.performed_writes) == expected_writes

    refreshed_index = json.loads(index_path.read_text(encoding="utf-8"))
    assert refreshed_index["generated_at"] == expected_timestamp

    refreshed_hot = json.loads(hot_path.read_text(encoding="utf-8"))
    assert refreshed_hot["generated_at"] == expected_timestamp
    assert len(refreshed_hot["nodes"]) == len(_HOT_NODE_IDS)

    expected_deps = {
        "alpha.md": (["beta.md"], []),
        "beta.md": (["gamma.md"], ["alpha.md"]),
        "gamma.md": (["delta.md"], ["beta.md"]),
        "delta.md": (["epsilon.md"], ["gamma.md"]),
        "epsilon.md": ([], ["delta.md"]),
    }
    for cap_id, (deps_out, deps_in) in expected_deps.items():
        refreshed = json.loads(cap_paths[cap_id].read_text(encoding="utf-8"))
        assert refreshed["deps_out"] == deps_out
        assert refreshed["deps_in"] == deps_in


def test_parse_args_returns_since_without_resolving(monkeypatch):
    monkeypatch.setattr(
        update.subprocess,
        "run",
        lambda *args, **kwargs: (_ for _ in ()).throw(RuntimeError("should not be called")),
    )

    options = update.parse_args(["--since", "feature"])

    assert options.targets == ()
    assert options.since == "feature"


def test_parse_args_supports_since_and_limits_scope(tmp_path, monkeypatch):
    caps_payloads = {
        cap_id: _caps_payload(cap_id, deps_out=["stale"], deps_in=["old"])
        for cap_id in (
            "alpha.md",
            "beta.md",
            "gamma.md",
            "delta.md",
            "epsilon.md",
            "zeta.md",
        )
    }
    root_base = tmp_path / "docs" / "birdseye"
    root, _, _, cap_paths = _prepare_birdseye(
        tmp_path,
        edges=[
            ["alpha.md", "beta.md"],
            ["beta.md", "gamma.md"],
            ["gamma.md", "delta.md"],
            ["delta.md", "epsilon.md"],
            ["epsilon.md", "zeta.md"],
        ],
        caps_payloads=caps_payloads,
        hot_entries=_HOT_NODE_IDS,
        root=root_base,
    )

    frozen_now = datetime(2025, 1, 3, tzinfo=timezone.utc)
    monkeypatch.setattr(update, "utc_now", lambda: frozen_now)

    monkeypatch.chdir(tmp_path)
    monkeypatch.setattr(update, "_REPO_ROOT", tmp_path)

    class FakeResolver:
        def __init__(self) -> None:
            self.calls: list[str] = []

        def resolve(self, reference: str) -> tuple[Path, ...]:
            self.calls.append(reference)
            return (
                Path("docs/birdseye/caps/delta.md.json"),
            )

    options = update.parse_args(["--emit", "caps", "--since"])

    resolver = FakeResolver()
    options = replace(options, diff_resolver=resolver)

    assert resolver.calls == []

    report = update.run_update(options)

    assert resolver.calls == ["main"]

    expected_caps = {
        cap_paths[cap_id]
        for cap_id in ("beta.md", "gamma.md", "delta.md", "epsilon.md", "zeta.md")
    }
    assert set(report.planned_writes) == expected_caps
    assert set(report.performed_writes) == expected_caps

    untouched = json.loads(cap_paths["alpha.md"].read_text(encoding="utf-8"))
    assert untouched["deps_out"] == ["stale"]
    assert untouched["deps_in"] == ["old"]


def test_git_diff_resolver_filters_paths(monkeypatch):
    captured = {}

    def fake_run(args, *, capture_output, text, check, cwd):
        captured.update(
            {
                "args": args,
                "cwd": cwd,
            }
        )
        return SimpleNamespace(
            stdout=(
                "M\tdocs/birdseye/caps/delta.md.json\n"
                "M\tREADME.md\n"
                "M\tdocs/birdseye/index.json\n"
            )
        )

    monkeypatch.setattr(update.subprocess, "run", fake_run)

    resolver = update.GitDiffResolver()

    result = resolver.resolve("develop")

    assert captured["args"] == [
        "git",
        "diff",
        "--name-status",
        "--find-renames",
        "--find-copies",
        "develop...HEAD",
    ]
    assert captured["cwd"] == update._REPO_ROOT
    assert result == (
        Path("docs/birdseye/caps/delta.md.json"),
        Path("docs/birdseye/caps/README.md.json"),
        Path("docs/birdseye/index.json"),
    )


def test_group_targets_resolves_caps_paths(tmp_path):
    root = tmp_path / "docs" / "birdseye"
    caps_dir = root / "caps"
    caps_dir.mkdir(parents=True)
    index_path = root / "index.json"
    hot_path = root / "hot.json"
    cap_file = caps_dir / "alpha.json"
    for path in (index_path, hot_path, cap_file):
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text("{}", encoding="utf-8")

    grouped = update._group_targets((root, index_path, caps_dir, cap_file))

    assert set(grouped) == {root}
    assert grouped[root] == [root, index_path, caps_dir, cap_file]


def test_resolve_focus_nodes_includes_two_hop_neighbours(tmp_path):
    root = tmp_path / "docs" / "birdseye"
    caps_dir = root / "caps"
    caps_dir.mkdir(parents=True)
    cap_alpha = caps_dir / "alpha.json"
    cap_beta = caps_dir / "beta.json"
    cap_gamma = caps_dir / "gamma.json"
    cap_delta = caps_dir / "delta.json"
    cap_paths = {
        "alpha": cap_alpha,
        "beta": cap_beta,
        "gamma": cap_gamma,
        "delta": cap_delta,
    }
    for path in cap_paths.values():
        path.write_text("{}", encoding="utf-8")

    caps_state = {
        cap_id: (cap_path, {"id": cap_id}, "{}")
        for cap_id, cap_path in cap_paths.items()
    }
    cap_lookup = {cap_path.resolve(): cap_id for cap_id, cap_path in cap_paths.items()}
    graph_out = {"alpha": ["beta"], "beta": ["gamma"], "gamma": ["delta"], "delta": []}
    graph_in = {"beta": ["alpha"], "gamma": ["beta"], "delta": ["gamma"], "alpha": []}

    available_caps = {cap_id: cap_path.resolve() for cap_id, cap_path in cap_paths.items()}

    options = update.UpdateOptions(targets=(cap_beta,), emit="caps")
    session = update.BirdseyeUpdateSession(options=options, timestamp="2025-01-01T00:00:00Z")
    focus = session._resolve_focus_nodes(
        (cap_beta,),
        root,
        graph_out,
        graph_in,
        caps_state,
        cap_lookup,
        available_caps,
    )

    assert focus == {"alpha", "beta", "gamma", "delta"}


def test_refresh_hot_updates_serial_and_preserves_metadata(tmp_path):
    hot_path = tmp_path / "hot.json"
    _write_json(
        hot_path,
        {
            "generated_at": "00042",
            **_HOT_HOTLIST_METADATA,
            "nodes": list(_HOT_NODES_FIXTURE),
        },
    )

    timestamp = "2025-01-02T00:00:00Z"
    hot_original = hot_path.read_text(encoding="utf-8")
    hot_payload = json.loads(hot_original)
    options = update.UpdateOptions(targets=(hot_path,), emit="index+caps", dry_run=False)
    session = update.BirdseyeUpdateSession(options=options, timestamp=timestamp)
    session.serial_allocator.observe("00042")
    session._plan_hot(hot_path, hot_payload, hot_original)
    plan = update.BirdseyePlan(
        generated_at=session._generated_at or timestamp,
        writes=tuple(session._writes),
    )
    report = session.execute(plan)

    assert report.planned_writes == (hot_path,)
    assert report.performed_writes == (hot_path,)
    refreshed = json.loads(hot_path.read_text(encoding="utf-8"))
    assert refreshed["generated_at"] == timestamp
    assert refreshed["index_snapshot"] == _HOT_INDEX_SNAPSHOT
    assert refreshed["refresh_command"] == _HOT_REFRESH_COMMAND
    assert refreshed["curation_notes"] == _HOT_CURATION_NOTES
    for node in refreshed["nodes"]:
        assert node["refresh_command"] == _HOT_REFRESH_COMMAND
        assert node["index_snapshot"] == _HOT_INDEX_SNAPSHOT
        assert node["curation_notes"] == _HOT_CURATION_NOTES


def test_refresh_capsule_updates_dependencies_and_serial(tmp_path):
    cap_path = tmp_path / "caps" / "alpha.json"
    cap_path.parent.mkdir(parents=True)
    _write_json(
        cap_path,
        _caps_payload(
            "alpha",
            deps_out=["stale"],
            deps_in=["obsolete"],
        )
        | {"generated_at": "00010"},
    )
    cap_original = cap_path.read_text(encoding="utf-8")
    caps_state = {"alpha": (cap_path, json.loads(cap_original), cap_original)}
    options = update.UpdateOptions(targets=(cap_path,), emit="caps", dry_run=False)
    session = update.BirdseyeUpdateSession(options=options, timestamp="2025-01-02T00:00:00Z")
    session.serial_allocator.observe("00010")
    session._plan_capsule(
        "alpha",
        caps_state["alpha"],
        {"alpha": ["beta"]},
        {"alpha": ["gamma"]},
    )
    plan = update.BirdseyePlan(
        generated_at=session._generated_at or "2025-01-02T00:00:00Z",
        writes=tuple(session._writes),
    )
    report = session.execute(plan)

    assert report.planned_writes == (cap_path,)
    assert report.performed_writes == (cap_path,)
    refreshed = json.loads(cap_path.read_text(encoding="utf-8"))
    assert refreshed["deps_out"] == ["beta"]
    assert refreshed["deps_in"] == ["gamma"]
    assert refreshed["generated_at"] == "2025-01-02T00:00:00Z"


def test_focus_resolver_returns_all_capsules_for_root_target():
    resolver = update.BirdseyeFocusResolver(radius=2)
    root = Path("/repo/docs/birdseye")
    root_targets = [root / "index.json"]
    graph_out: Mapping[str, Sequence[str]] = {}
    graph_in: Mapping[str, Sequence[str]] = {}
    caps_state: update.CapsuleState = {
        "alpha": ((root / "caps" / "alpha.json").resolve(), {}, ""),
        "beta": ((root / "caps" / "beta.json").resolve(), {}, ""),
    }
    cap_path_lookup: Mapping[Path, str] = {
        (root / "caps" / "alpha.json").resolve(): "alpha",
        (root / "caps" / "beta.json").resolve(): "beta",
        (root / "caps" / "gamma.json").resolve(): "gamma",
    }
    available_caps: Mapping[str, Path] = {
        "gamma": (root / "caps" / "gamma.json").resolve(),
    }

    focus = resolver.resolve(
        root_targets=root_targets,
        root=root,
        graph_out=graph_out,
        graph_in=graph_in,
        caps_state=caps_state,
        cap_path_lookup=cap_path_lookup,
        available_caps=available_caps,
    )

    assert focus == {"alpha", "beta", "gamma"}


def test_focus_resolver_includes_two_hop_neighbours():
    resolver = update.BirdseyeFocusResolver(radius=2)
    root = Path("/repo/docs/birdseye")
    root_targets = [(root / "caps" / "b.json").resolve()]
    graph_out: Mapping[str, Sequence[str]] = {
        "a": ["b"],
        "b": ["c"],
        "c": ["d"],
    }
    graph_in: Mapping[str, Sequence[str]] = {
        "b": ["a"],
        "a": ["x"],
    }
    caps_state: update.CapsuleState = {
        "a": ((root / "caps" / "a.json").resolve(), {}, ""),
        "b": ((root / "caps" / "b.json").resolve(), {}, ""),
    }
    cap_path_lookup: Mapping[Path, str] = {
        (root / "caps" / "a.json").resolve(): "a",
        (root / "caps" / "b.json").resolve(): "b",
        (root / "caps" / "c.json").resolve(): "c",
        (root / "caps" / "d.json").resolve(): "d",
        (root / "caps" / "x.json").resolve(): "x",
    }
    available_caps: Mapping[str, Path] = {
        "c": (root / "caps" / "c.json").resolve(),
        "d": (root / "caps" / "d.json").resolve(),
        "x": (root / "caps" / "x.json").resolve(),
    }

    focus = resolver.resolve(
        root_targets=root_targets,
        root=root,
        graph_out=graph_out,
        graph_in=graph_in,
        caps_state=caps_state,
        cap_path_lookup=cap_path_lookup,
        available_caps=available_caps,
    )

    assert focus == {"a", "b", "c", "d", "x"}
