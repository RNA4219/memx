# SPDX-License-Identifier: Apache-2.0
# Copyright 2025 RNA4219

"""Birdseye再生成ツール。"""

from __future__ import annotations

import argparse
import json
import re
import subprocess
import sys
from collections import deque
from dataclasses import dataclass, replace
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Iterable, Mapping, Protocol, Sequence


_REPO_ROOT = Path(__file__).resolve().parents[2]


_BIRDSEYE_REGENERATE_COMMAND = (
    "python tools/codemap/update.py --targets "
    "docs/birdseye/index.json,docs/birdseye/hot.json --emit index+caps"
)


CapsuleEntry = tuple[Path, dict[str, Any], str]
CapsuleState = dict[str, CapsuleEntry]
Graph = dict[str, list[str]]


def _resolve_focus_nodes_for_targets(
    root_targets: Sequence[Path],
    root: Path,
    graph_out: Mapping[str, Sequence[str]],
    graph_in: Mapping[str, Sequence[str]],
    caps_state: CapsuleState,
    cap_path_lookup: Mapping[Path, str],
    available_caps: Mapping[str, Path],
) -> set[str]:
    combined_caps: dict[str, Path] = dict(available_caps)
    for cap_id, (cap_path, _cap_data, _cap_original) in caps_state.items():
        combined_caps.setdefault(cap_id, cap_path)
    if not combined_caps:
        return set()

    focus_nodes: set[str] = set()
    root_resolved = root.resolve()
    index_resolved = (root / "index.json").resolve()
    caps_dir_resolved = (root / "caps").resolve()
    hot_resolved = (root / "hot.json").resolve()
    special_roots = {root_resolved, index_resolved, caps_dir_resolved, hot_resolved}

    for candidate in root_targets:
        resolved = candidate.resolve()
        if resolved in special_roots:
            return set(combined_caps)
        cap_id = cap_path_lookup.get(resolved)
        if cap_id:
            focus_nodes.add(cap_id)

    if not focus_nodes:
        focus_nodes = set(caps_state) or set(combined_caps)

    seen: set[str] = set()
    queue: deque[tuple[str, int]] = deque((node, 0) for node in focus_nodes)
    while queue:
        node, distance = queue.popleft()
        if node in seen or distance > 2:
            continue
        seen.add(node)
        if distance == 2:
            continue
        for neighbour in graph_out.get(node, ()):
            if neighbour not in seen:
                queue.append((neighbour, distance + 1))
        for neighbour in graph_in.get(node, ()):
            if neighbour not in seen:
                queue.append((neighbour, distance + 1))

    return {node for node in seen if node in combined_caps}


class TargetResolutionError(RuntimeError):
    """Raised when Birdseye targets cannot be resolved."""


class DiffResolver(Protocol):
    def resolve(self, reference: str) -> tuple[Path, ...]:
        ...


@dataclass(frozen=True)
class UpdateOptions:
    targets: tuple[Path, ...]
    emit: str
    dry_run: bool = False
    since: str | None = None
    diff_resolver: DiffResolver | None = None

    def resolve_targets(self) -> tuple[Path, ...]:
        resolved = [_normalise_target(path) for path in self.targets]
        if self.since:
            if self.diff_resolver is None:
                raise TargetResolutionError("Diff resolver is required when --since is used")
            try:
                derived = self.diff_resolver.resolve(self.since)
            except subprocess.CalledProcessError as exc:
                raise TargetResolutionError(
                    f"Failed to resolve git diff for --since: {exc}"
                ) from exc
            resolved.extend(_normalise_target(path) for path in derived)
        if not resolved:
            resolved.extend(_default_birdseye_targets())
        unique_targets = tuple(dict.fromkeys(resolved))
        if not unique_targets:
            raise TargetResolutionError("Specify --targets, --since, or both")
        return unique_targets


@dataclass(frozen=True)
class UpdateReport:
    generated_at: str
    planned_writes: tuple[Path, ...]
    performed_writes: tuple[Path, ...]


def _format_diff_reference(reference: str) -> str:
    stripped = reference.strip()
    if ".." in stripped:
        return stripped
    return f"{stripped}...HEAD"


class GitDiffParser:
    def __init__(self, *, repo_root: Path | None = None) -> None:
        self._repo_root = repo_root or _REPO_ROOT

    def parse(self, payload: str) -> tuple[Path, ...]:
        diff_entries: list[str] = []
        for raw_line in payload.splitlines():
            stripped = raw_line.strip()
            if not stripped:
                continue
            columns = [segment.strip() for segment in stripped.split("\t")]
            if not columns:
                continue
            status = columns[0]
            if not status:
                continue
            kind = status[0].upper()
            payloads: list[str]
            if kind in {"R", "C"} and len(columns) >= 3:
                payloads = columns[1:3]
            elif kind in {"A", "M", "D", "T", "U"} and len(columns) >= 2:
                payloads = [columns[1]]
            else:
                continue
            for value in payloads:
                candidate = value.replace("\\", "/").strip()
                if not candidate:
                    continue
                diff_entries.append(candidate)
        derived = _derive_targets_from_since(diff_entries, repo_root=self._repo_root)
        return tuple(
            path
            for path in derived
            if len(path.parts) >= 2 and path.parts[0] == "docs" and path.parts[1] == "birdseye"
        )


class _GitDiffRunner(Protocol):
    def __call__(self, args: Sequence[str]) -> subprocess.CompletedProcess[str]:
        ...


def _run_git_diff(args: Sequence[str]) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        args,
        capture_output=True,
        text=True,
        check=True,
        cwd=_REPO_ROOT,
    )


class GitDiffResolver:
    def __init__(
        self,
        *,
        parser: GitDiffParser | None = None,
        runner: _GitDiffRunner | None = None,
    ) -> None:
        self._parser = parser or GitDiffParser()
        self._runner = runner or _run_git_diff

    def resolve(self, reference: str) -> tuple[Path, ...]:
        diff_reference = _format_diff_reference(reference)
        command = [
            "git",
            "diff",
            "--name-status",
            "--find-renames",
            "--find-copies",
            diff_reference,
        ]
        result = self._runner(command)
        return self._parser.parse(result.stdout)


def _derive_targets_from_since(
    diff_paths: Iterable[str | Path], *, repo_root: Path | None = None
) -> tuple[Path, ...]:
    base_root = repo_root or _REPO_ROOT
    derived: list[Path] = []
    seen: set[Path] = set()

    def _add(path: Path) -> None:
        if path not in seen:
            seen.add(path)
            derived.append(path)

    for path in diff_paths:
        if isinstance(path, Path):
            raw_entry = path.as_posix()
        else:
            raw_entry = path
        raw_entry = raw_entry.replace("\\", "/").strip()
        if not raw_entry:
            continue
        segments: list[str] = []
        for marker in (" -> ", " => "):
            if marker in raw_entry:
                left, right = raw_entry.split(marker, 1)
                segments.extend([left.strip(), right.strip()])
                break
        if not segments:
            segments.append(raw_entry)
        for segment in segments:
            candidate = segment
            if not candidate:
                continue
            while candidate.startswith("./"):
                candidate = candidate[2:]
            candidate = candidate.rstrip("/")
            if not candidate:
                continue
            normalised = Path(candidate)
            if normalised.is_absolute():
                try:
                    normalised = normalised.relative_to(base_root)
                except ValueError:
                    continue
            if normalised.parts[:2] == ("docs", "birdseye"):
                _add(normalised)
                continue
            candidate_slug = normalised.as_posix() if normalised.parts else candidate
            capsule_slug = candidate_slug.replace("/", ".")
            capsule_path = Path("docs/birdseye/caps") / f"{capsule_slug}.json"
            if (base_root / capsule_path).is_file():
                _add(capsule_path)
    return tuple(derived)


def _default_birdseye_targets() -> tuple[Path, ...]:
    birdseye_root = _REPO_ROOT / "docs" / "birdseye"
    candidates = (
        birdseye_root / "index.json",
        birdseye_root / "hot.json",
        birdseye_root / "caps",
    )
    fallback: list[Path] = []
    for candidate in candidates:
        if candidate.is_file() or candidate.is_dir():
            fallback.append(candidate.resolve())
    return tuple(fallback)


def _normalise_target(target: Path) -> Path:
    if target.is_absolute():
        return target.resolve()
    return (_REPO_ROOT / target).resolve()


def parse_args(argv: Iterable[str] | None = None) -> UpdateOptions:
    parser = argparse.ArgumentParser(
        description="Regenerate Birdseye index and capsules.",
    )
    parser.add_argument(
        "--targets",
        type=str,
        help="Comma-separated list of Birdseye resources to analyse.",
    )
    parser.add_argument(
        "--since",
        type=str,
        nargs="?",
        const="main",
        help="Derive targets from git diff since the given reference (default: main).",
    )
    parser.add_argument(
        "--emit",
        type=str,
        choices=("index", "caps", "index+caps"),
        default="index+caps",
        help="Select which artefacts to write.",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Compute updates without writing to disk.",
    )
    args = parser.parse_args(list(argv) if argv is not None else None)
    target_paths: list[Path] = []
    if args.targets:
        target_paths.extend(
            Path(value.strip()) for value in args.targets.split(",") if value.strip()
        )
    if not target_paths and not args.since:
        parser.error("Specify --targets, --since, or both")
    normalised_targets = tuple(
        dict.fromkeys(_normalise_target(path) for path in target_paths)
    )
    return UpdateOptions(
        targets=normalised_targets,
        emit=args.emit,
        dry_run=args.dry_run,
        since=args.since,
    )


def ensure_python_version() -> None:
    if sys.version_info < (3, 11):
        print("[ERROR] Python 3.11 or newer is required.")
        raise SystemExit(1)


def utc_now() -> datetime:
    return datetime.now(timezone.utc)


def _format_timestamp(moment: datetime) -> str:
    return moment.replace(microsecond=0).isoformat().replace("+00:00", "Z")


def _resolve_root(target: Path) -> Path:
    resolved = target.resolve()
    if resolved.is_dir():
        if resolved.name == "caps":
            return resolved.parent.resolve()
        return resolved
    parent = resolved.parent
    if parent.name == "caps":
        return parent.parent.resolve()
    return parent.resolve()


def _load_json(path: Path) -> tuple[Any, str]:
    raw = path.read_text(encoding="utf-8")
    return json.loads(raw), raw


def _dump_json(data: Any) -> str:
    return json.dumps(data, ensure_ascii=False, indent=2) + "\n"


def _sorted_unique(items: Sequence[str]) -> list[str]:
    return sorted(dict.fromkeys(items))


_SERIAL_PATTERN = re.compile(r"\d{5}")


def _coerce_serial(candidate: Any) -> int | None:
    if isinstance(candidate, str) and _SERIAL_PATTERN.fullmatch(candidate):
        return int(candidate)
    return None


class _SerialAllocator:
    __slots__ = ("max_serial", "_next_serial")

    def __init__(self) -> None:
        self.max_serial = 0
        self._next_serial: int | None = None

    def observe(self, candidate: Any) -> None:
        value = _coerce_serial(candidate)
        if value is not None and value > self.max_serial:
            self.max_serial = value

    def allocate(self, existing: Any) -> str:
        candidate = _coerce_serial(existing)
        if candidate is not None and candidate > self.max_serial:
            self.max_serial = candidate
        if self._next_serial is None:
            self._next_serial = self.max_serial + 1 if self.max_serial else 1
        target = self._next_serial
        if target > self.max_serial:
            self.max_serial = target
        return f"{target:05d}"


def _next_generated_at(existing: Any, fallback: str, *, allocator: _SerialAllocator) -> str:
    del existing, allocator
    return fallback


def _resolve_node_mtime(node_id: str) -> str | None:
    target = (_REPO_ROOT / node_id).resolve()
    if not target.exists():
        return None
    return _format_timestamp(datetime.fromtimestamp(target.stat().st_mtime, tz=timezone.utc))


@dataclass(frozen=True)
class PlannedWrite:
    path: Path
    content: str
    original: str


@dataclass(frozen=True)
class BirdseyePlan:
    generated_at: str
    writes: tuple[PlannedWrite, ...]


@dataclass(frozen=True)
class BirdseyeRootPlan:
    loads: tuple[str, ...]
    writes: tuple[PlannedWrite, ...]
    remembered: str | None = None

    def apply(self, session: "BirdseyeUpdateSession") -> None:
        session._writes.extend(self.writes)
        if self.remembered is not None:
            session._remember_generated(self.remembered)

class BirdseyeFocusResolver:
    def __init__(self, *, radius: int = 2) -> None:
        self.radius = radius

    def resolve(
        self,
        root_targets: Sequence[Path],
        root: Path,
        graph_out: Mapping[str, Sequence[str]],
        graph_in: Mapping[str, Sequence[str]],
        caps_state: CapsuleState,
        cap_path_lookup: Mapping[Path, str],
        available_caps: Mapping[str, Path],
    ) -> set[str]:
        combined_caps = self._merge_capsules(caps_state, available_caps)
        if not combined_caps:
            return set()
        if self._targets_include_root(root_targets, root):
            return set(combined_caps)
        focus_nodes = self._initial_focus_nodes(root_targets, cap_path_lookup)
        if not focus_nodes:
            focus_nodes = set(caps_state) or set(combined_caps)
        seen = self._expand_focus(focus_nodes, graph_out, graph_in)
        return {node for node in seen if node in combined_caps}

    def _merge_capsules(
        self,
        caps_state: CapsuleState,
        available_caps: Mapping[str, Path],
    ) -> dict[str, Path]:
        combined: dict[str, Path] = dict(available_caps)
        for cap_id, (cap_path, _cap_data, _cap_original) in caps_state.items():
            combined.setdefault(cap_id, cap_path)
        return combined

    def _targets_include_root(
        self, root_targets: Sequence[Path], root: Path
    ) -> bool:
        root_resolved = root.resolve()
        index_resolved = (root / "index.json").resolve()
        caps_dir_resolved = (root / "caps").resolve()
        hot_resolved = (root / "hot.json").resolve()
        special_roots = {
            root_resolved,
            index_resolved,
            caps_dir_resolved,
            hot_resolved,
        }
        for candidate in root_targets:
            resolved = candidate.resolve()
            if resolved in special_roots:
                return True
        return False

    def _initial_focus_nodes(
        self,
        root_targets: Sequence[Path],
        cap_path_lookup: Mapping[Path, str],
    ) -> set[str]:
        seeds: set[str] = set()
        for candidate in root_targets:
            resolved = candidate.resolve()
            cap_id = cap_path_lookup.get(resolved)
            if cap_id:
                seeds.add(cap_id)
        return seeds

    def _expand_focus(
        self,
        focus_nodes: set[str],
        graph_out: Mapping[str, Sequence[str]],
        graph_in: Mapping[str, Sequence[str]],
    ) -> set[str]:
        if self.radius < 0:
            return set()
        seen: set[str] = set()
        queue: deque[tuple[str, int]] = deque((node, 0) for node in focus_nodes)
        empty: Sequence[str] = ()
        while queue:
            node, distance = queue.popleft()
            if node in seen or distance > self.radius:
                continue
            seen.add(node)
            if distance == self.radius:
                continue
            for neighbour in graph_out.get(node, empty):
                if neighbour not in seen:
                    queue.append((neighbour, distance + 1))
            for neighbour in graph_in.get(node, empty):
                if neighbour not in seen:
                    queue.append((neighbour, distance + 1))
        return seen


class BirdseyeRootBuilder:
    def __init__(
        self,
        *,
        root: Path,
        root_targets: Sequence[Path],
        emit_index: bool,
        emit_caps: bool,
        timestamp: str,
        serial_allocator: _SerialAllocator,
    ) -> None:
        self.root = root
        self.root_targets = root_targets
        self.emit_index = emit_index
        self.emit_caps = emit_caps
        self.timestamp = timestamp
        self.serial_allocator = serial_allocator
        self.index_path = root / "index.json"
        self.hot_path = root / "hot.json"
        self.caps_dir = root / "caps"
        self._loads: list[str] = []
        self._writes: list[PlannedWrite] = []
        self._focus_resolver = BirdseyeFocusResolver(radius=2)
        self._first_generated: str | None = None

    def build(self) -> BirdseyeRootPlan:
        self._validate()
        index_data, index_original = self._load_index()
        graph_out, graph_in = _build_graph(index_data)

        hot_data: dict[str, Any] | None = None
        hot_original: str | None = None
        if self.emit_index or self.emit_caps:
            hot_data, hot_original = self._load_hot()

        caps_state: CapsuleState = {}
        cap_path_lookup: dict[Path, str] = {}
        available_caps: dict[str, Path] = {}
        if self.emit_caps:
            caps_state, cap_path_lookup, available_caps = self._load_capsules(index_data)

        if self.emit_index:
            self._plan_index(index_data, index_original)
            if hot_data is not None and hot_original is not None:
                self._plan_hot(hot_data, hot_original)

        if self.emit_caps and available_caps:
            focus_nodes = self._resolve_focus_nodes(
                graph_out,
                graph_in,
                caps_state,
                cap_path_lookup,
                available_caps,
            )
            self._ensure_placeholders(focus_nodes, caps_state, available_caps)
            for cap_id in sorted(focus_nodes):
                self._plan_capsule(cap_id, caps_state[cap_id], graph_out, graph_in)

        return BirdseyeRootPlan(
            loads=tuple(self._loads),
            writes=tuple(self._writes),
            remembered=self._first_generated,
        )

    def _validate(self) -> None:
        if not self.index_path.is_file():
            raise FileNotFoundError(self.index_path)
        if self.emit_index and not self.hot_path.exists():
            raise FileNotFoundError(
                f"{self.hot_path} is missing. Regenerate via: {_BIRDSEYE_REGENERATE_COMMAND}"
            )
        if self.emit_caps and not self.caps_dir.is_dir():
            raise FileNotFoundError(self.caps_dir)

    def _load_index(self) -> tuple[dict[str, Any], str]:
        self._loads.append("index")
        loaded_index, index_original = _load_json(self.index_path)
        index_data = loaded_index if isinstance(loaded_index, dict) else {}
        self.serial_allocator.observe(index_data.get("generated_at"))
        return index_data, index_original

    def _load_hot(self) -> tuple[dict[str, Any] | None, str | None]:
        self._loads.append("hot")
        if not self.hot_path.exists():
            return None, None
        loaded_hot, hot_original = _load_json(self.hot_path)
        if not isinstance(loaded_hot, dict):
            return None, None
        self.serial_allocator.observe(loaded_hot.get("generated_at"))
        return loaded_hot, hot_original

    def _load_capsules(
        self, index_data: Mapping[str, Any]
    ) -> tuple[CapsuleState, dict[Path, str], dict[str, Path]]:
        self._loads.append("caps")
        caps_state: CapsuleState = {}
        cap_path_lookup: dict[Path, str] = {}
        available_caps: dict[str, Path] = {}
        raw_nodes = index_data.get("nodes", {})
        if isinstance(raw_nodes, Mapping):
            iterable_nodes = raw_nodes.items()
        elif isinstance(raw_nodes, Sequence):
            iterable_nodes = (
                (node_payload.get("node_id"), node_payload)
                for node_payload in raw_nodes
                if isinstance(node_payload, Mapping)
            )
        else:
            iterable_nodes = ()

        for node_id, node_payload in iterable_nodes:
            if not isinstance(node_id, str) or not isinstance(node_payload, Mapping):
                continue
            caps_ref = node_payload.get("caps")
            if not isinstance(caps_ref, str) or not caps_ref:
                caps_ref = node_payload.get("capsule")
            if not isinstance(caps_ref, str) or not caps_ref:
                continue
            candidate_path = Path(caps_ref)
            if candidate_path.is_absolute():
                resolved_candidate = candidate_path.resolve()
            else:
                resolved_candidate = (_REPO_ROOT / candidate_path).resolve()
            available_caps.setdefault(node_id, resolved_candidate)
            cap_path_lookup.setdefault(resolved_candidate, node_id)
        for cap_path in sorted(self.caps_dir.glob("*.json")):
            if not cap_path.is_file():
                continue
            cap_data, cap_original = _load_json(cap_path)
            if not isinstance(cap_data, dict):
                continue
            cap_id = cap_data.get("id")
            if not isinstance(cap_id, str):
                continue
            cap_path_resolved = cap_path.resolve()
            caps_state[cap_id] = (cap_path_resolved, cap_data, cap_original)
            cap_path_lookup[cap_path_resolved] = cap_id
            available_caps.setdefault(cap_id, cap_path_resolved)
            self.serial_allocator.observe(cap_data.get("generated_at"))
        return caps_state, cap_path_lookup, available_caps

    def _plan_index(self, index_data: dict[str, Any], index_original: str) -> None:
        new_generated = _next_generated_at(
            index_data.get("generated_at"),
            self.timestamp,
            allocator=self.serial_allocator,
        )
        if index_data.get("generated_at") != new_generated:
            index_data["generated_at"] = new_generated
            self._remember_generated(new_generated)

        nodes = index_data.get("nodes")
        if isinstance(nodes, dict):
            iterable_nodes = nodes.items()
        elif isinstance(nodes, Sequence):
            iterable_nodes = (
                (payload.get("node_id"), payload)
                for payload in nodes
                if isinstance(payload, dict)
            )
        else:
            iterable_nodes = ()

        for node_id, payload in iterable_nodes:
            if not isinstance(node_id, str) or not isinstance(payload, dict):
                continue
            mtime = _resolve_node_mtime(node_id)
            if mtime is None:
                continue
            if payload.get("mtime") != mtime:
                payload["mtime"] = mtime

        serialized = _dump_json(index_data)
        if serialized != index_original:
            self._writes.append(
                PlannedWrite(path=self.index_path, content=serialized, original=index_original)
            )

    def _plan_hot(self, hot_data: dict[str, Any], hot_original: str) -> None:
        new_generated = _next_generated_at(
            hot_data.get("generated_at"),
            self.timestamp,
            allocator=self.serial_allocator,
        )
        if hot_data.get("generated_at") != new_generated:
            hot_data["generated_at"] = new_generated
            self._remember_generated(new_generated)
        serialized = _dump_json(hot_data)
        if serialized != hot_original:
            self._writes.append(
                PlannedWrite(path=self.hot_path, content=serialized, original=hot_original)
            )

    def _resolve_focus_nodes(
        self,
        graph_out: Mapping[str, Sequence[str]],
        graph_in: Mapping[str, Sequence[str]],
        caps_state: CapsuleState,
        cap_path_lookup: Mapping[Path, str],
        available_caps: Mapping[str, Path],
    ) -> set[str]:
        return _resolve_focus_nodes_for_targets(
            self.root_targets,
            self.root,
            graph_out,
            graph_in,
            caps_state,
            cap_path_lookup,
            available_caps,
        )

    def _ensure_placeholders(
        self,
        focus_nodes: Iterable[str],
        caps_state: CapsuleState,
        available_caps: Mapping[str, Path],
    ) -> None:
        for cap_id in focus_nodes:
            if cap_id in caps_state:
                continue
            cap_path = available_caps.get(cap_id)
            if cap_path is None:
                continue
            placeholder_data: dict[str, Any] = {
                "id": cap_id,
                "role": "doc",
                "public_api": [],
                "summary": cap_id,
                "deps_out": [],
                "deps_in": [],
                "risks": [],
                "tests": [],
            }
            caps_state[cap_id] = (cap_path, placeholder_data, "")

    def _plan_capsule(
        self,
        cap_id: str,
        capsule: CapsuleEntry,
        graph_out: Mapping[str, Sequence[str]],
        graph_in: Mapping[str, Sequence[str]],
    ) -> None:
        cap_path, cap_data, cap_original = capsule
        expected_out = _sorted_unique(graph_out.get(cap_id, []))
        expected_in = _sorted_unique(graph_in.get(cap_id, []))
        updated = False
        if cap_data.get("deps_out") != expected_out:
            cap_data["deps_out"] = expected_out
            updated = True
        if cap_data.get("deps_in") != expected_in:
            cap_data["deps_in"] = expected_in
            updated = True
        new_generated = _next_generated_at(
            cap_data.get("generated_at"),
            self.timestamp,
            allocator=self.serial_allocator,
        )
        if cap_data.get("generated_at") != new_generated:
            cap_data["generated_at"] = new_generated
            updated = True
            self._remember_generated(new_generated)
        if updated:
            serialized = _dump_json(cap_data)
            if serialized != cap_original:
                self._writes.append(
                    PlannedWrite(path=cap_path, content=serialized, original=cap_original)
                )

    def _remember_generated(self, value: str) -> None:
        if self._first_generated is None:
            self._first_generated = value

def _finalise(paths: set[Path]) -> tuple[Path, ...]:
    return tuple(sorted(paths, key=lambda candidate: candidate.as_posix()))


def _group_targets(targets: Iterable[Path]) -> dict[Path, list[Path]]:
    grouped: dict[Path, list[Path]] = {}
    for target in targets:
        normalised = _normalise_target(target)
        root = _resolve_root(normalised)
        grouped.setdefault(root, []).append(normalised)
    return grouped


def _build_graph(index_data: Mapping[str, Any]) -> tuple[Graph, Graph]:
    raw_nodes = index_data.get("nodes", {})
    if isinstance(raw_nodes, Mapping):
        node_ids = [node for node in raw_nodes if isinstance(node, str)]
    elif isinstance(raw_nodes, Sequence):
        node_ids = [
            node_payload.get("node_id")
            for node_payload in raw_nodes
            if isinstance(node_payload, Mapping)
            and isinstance(node_payload.get("node_id"), str)
        ]
    else:
        node_ids = []
    graph_out: Graph = {node: [] for node in node_ids}
    graph_in: Graph = {node: [] for node in node_ids}
    for raw_edge in index_data.get("edges", []):
        if not isinstance(raw_edge, Sequence) or len(raw_edge) != 2:
            continue
        source, destination = raw_edge
        if not isinstance(source, str) or not isinstance(destination, str):
            continue
        graph_out.setdefault(source, []).append(destination)
        graph_in.setdefault(destination, []).append(source)
        graph_out.setdefault(destination, graph_out.get(destination, []))
        graph_in.setdefault(source, graph_in.get(source, []))
    for values in graph_out.values():
        values.sort()
    for values in graph_in.values():
        values.sort()
    return graph_out, graph_in



class BirdseyeUpdateSession:
    def __init__(self, *, options: UpdateOptions, timestamp: str | None = None) -> None:
        self.options = options
        self.emit_index = options.emit in {"index", "index+caps"}
        self.emit_caps = options.emit in {"caps", "index+caps"}
        self.timestamp = timestamp or _format_timestamp(utc_now())
        self.serial_allocator = _SerialAllocator()
        self._generated_at: str | None = None
        self._writes: list[PlannedWrite] = []

    def plan(self) -> BirdseyePlan:
        resolved_targets = self.options.resolve_targets()
        grouped = _group_targets(resolved_targets)
        for root, root_targets in grouped.items():
            self._plan_for_root(root, root_targets)
        return BirdseyePlan(
            generated_at=self._generated_at or self.timestamp,
            writes=tuple(self._writes),
        )

    def execute(self, plan: BirdseyePlan) -> UpdateReport:
        planned_paths = {write.path for write in plan.writes}
        performed: set[Path] = set()
        if not self.options.dry_run:
            for write in plan.writes:
                write.path.parent.mkdir(parents=True, exist_ok=True)
                write.path.write_text(write.content, encoding='utf-8')
                performed.add(write.path)
        return UpdateReport(
            generated_at=plan.generated_at,
            planned_writes=_finalise(planned_paths),
            performed_writes=_finalise(performed),
        )

    def _plan_for_root(self, root: Path, root_targets: Sequence[Path]) -> None:
        builder = BirdseyeRootBuilder(
            root=root,
            root_targets=root_targets,
            emit_index=self.emit_index,
            emit_caps=self.emit_caps,
            timestamp=self.timestamp,
            serial_allocator=self.serial_allocator,
        )
        plan = builder.build()
        plan.apply(self)

    def _resolve_focus_nodes(
        self,
        root_targets: Sequence[Path],
        root: Path,
        graph_out: Mapping[str, Sequence[str]],
        graph_in: Mapping[str, Sequence[str]],
        caps_state: CapsuleState,
        cap_path_lookup: Mapping[Path, str],
        available_caps: Mapping[str, Path],
    ) -> set[str]:
        return _resolve_focus_nodes_for_targets(
            root_targets,
            root,
            graph_out,
            graph_in,
            caps_state,
            cap_path_lookup,
            available_caps,
        )

    def _plan_hot(
        self,
        hot_path: Path,
        hot_data: dict[str, Any],
        hot_original: str,
    ) -> None:
        new_generated = _next_generated_at(
            hot_data.get("generated_at"),
            self.timestamp,
            allocator=self.serial_allocator,
        )
        if hot_data.get("generated_at") != new_generated:
            hot_data["generated_at"] = new_generated
            self._remember_generated(new_generated)
        serialized = _dump_json(hot_data)
        if serialized != hot_original:
            self._writes.append(
                PlannedWrite(path=hot_path, content=serialized, original=hot_original)
            )

    def _plan_capsule(
        self,
        cap_id: str,
        capsule: CapsuleEntry,
        graph_out: Mapping[str, Sequence[str]],
        graph_in: Mapping[str, Sequence[str]],
    ) -> None:
        cap_path, cap_data, cap_original = capsule
        expected_out = _sorted_unique(graph_out.get(cap_id, []))
        expected_in = _sorted_unique(graph_in.get(cap_id, []))
        updated = False
        if cap_data.get("deps_out") != expected_out:
            cap_data["deps_out"] = expected_out
            updated = True
        if cap_data.get("deps_in") != expected_in:
            cap_data["deps_in"] = expected_in
            updated = True
        new_generated = _next_generated_at(
            cap_data.get("generated_at"),
            self.timestamp,
            allocator=self.serial_allocator,
        )
        if cap_data.get("generated_at") != new_generated:
            cap_data["generated_at"] = new_generated
            updated = True
            self._remember_generated(new_generated)
        if updated:
            serialized = _dump_json(cap_data)
            if serialized != cap_original:
                self._writes.append(
                    PlannedWrite(path=cap_path, content=serialized, original=cap_original)
                )

    def _remember_generated(self, value: str) -> None:
        if self._generated_at is None:
            self._generated_at = value


def run_update(options: UpdateOptions) -> UpdateReport:
    timestamp = _format_timestamp(utc_now())
    session = BirdseyeUpdateSession(options=options, timestamp=timestamp)
    plan = session.plan()
    return session.execute(plan)

def main(argv: Iterable[str] | None = None) -> int:
    ensure_python_version()
    options = parse_args(argv)
    options = replace(options, diff_resolver=GitDiffResolver())
    try:
        run_update(options)
    except TargetResolutionError as exc:
        print(f"[ERROR] {exc}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
