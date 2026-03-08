from __future__ import annotations

import json
import os
import platform
import statistics
import subprocess
import sys
import time
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[2]
GO_ROOT = REPO_ROOT / "memx_spec_v3" / "go"
PERF_DIR = REPO_ROOT / "artifacts" / "perf"
API_URL = os.environ.get("MEMX_PERF_API_URL", "http://127.0.0.1:7791")
MEMX_PERF_BIN = os.environ.get("MEMX_PERF_BIN", "")
SEED_COUNT = 10_000
WARMUP = 20
RUNS = 200
QUERY = "alpha"
BODY = ("alpha " * 84)[:500]


def percentile(values: list[float], p: float) -> float:
    ordered = sorted(values)
    index = int((len(ordered) - 1) * p)
    return round(ordered[index], 2)


def mem_command(args: list[str]) -> list[str]:
    if MEMX_PERF_BIN:
        return [MEMX_PERF_BIN, *args]
    return ["go", "run", "./cmd/mem", *args]


def run_mem(args: list[str], *, stdin_body: str | None = None, capture_stdout: bool = False) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        mem_command(args),
        cwd=GO_ROOT,
        input=stdin_body,
        text=True,
        capture_output=capture_stdout,
        check=True,
    )


def timed_mem(args: list[str], *, stdin_body: str | None = None) -> float:
    started = time.perf_counter()
    subprocess.run(
        mem_command(args),
        cwd=GO_ROOT,
        input=stdin_body,
        text=True,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
        check=True,
    )
    return (time.perf_counter() - started) * 1000


def ingest_args(title: str) -> list[str]:
    return [
        "in",
        "short",
        "--stdin",
        "--title",
        title,
        "--no-llm",
        "--api-url",
        API_URL,
    ]


def search_args() -> list[str]:
    return [
        "out",
        "search",
        "--api-url",
        API_URL,
        QUERY,
    ]


def show_args(note_id: str) -> list[str]:
    return [
        "out",
        "show",
        "--api-url",
        API_URL,
        note_id,
    ]


def main() -> int:
    PERF_DIR.mkdir(parents=True, exist_ok=True)

    seed_ids: list[str] = []
    for index in range(SEED_COUNT):
        completed = run_mem(ingest_args(f"perf-seed-{index:05d}"), stdin_body=BODY, capture_stdout=True)
        line = completed.stdout.strip()
        if not line.startswith("ok id="):
            raise RuntimeError(f"unexpected ingest output at seed {index}: {line!r}")
        seed_ids.append(line.removeprefix("ok id="))

    known_id = seed_ids[0]
    (PERF_DIR / "seed-result.json").write_text(
        json.dumps(
            {
                "count": len(seed_ids),
                "body_chars": len(BODY),
                "query": QUERY,
                "api_url": API_URL,
                "known_id": known_id,
                "note_ids": seed_ids[:5],
            },
            ensure_ascii=False,
            indent=2,
        ),
        encoding="utf-8",
    )

    for index in range(WARMUP):
        timed_mem(ingest_args(f"perf-warmup-{index:02d}"), stdin_body=BODY)
        timed_mem(search_args())
        timed_mem(show_args(known_id))

    (PERF_DIR / "warmup-result.json").write_text(
        json.dumps(
            {
                "warmup": WARMUP,
                "query": QUERY,
                "known_id": known_id,
                "body_chars": len(BODY),
                "api_url": API_URL,
            },
            ensure_ascii=False,
            indent=2,
        ),
        encoding="utf-8",
    )

    ingest = [timed_mem(ingest_args(f"perf-bench-{index:03d}"), stdin_body=BODY) for index in range(RUNS)]
    search = [timed_mem(search_args()) for _ in range(RUNS)]
    show = [timed_mem(show_args(known_id)) for _ in range(RUNS)]

    result = {
        "environment": {
            "cpu": f"{os.cpu_count()} logical CPUs",
            "ram": "unknown",
            "storage": "unknown",
            "os": platform.platform(),
        },
        "dataset": {
            "store": "short",
            "note_count": SEED_COUNT,
            "body_chars": len(BODY),
        },
        "results": {
            "ingest": {
                "p50_ms": percentile(ingest, 0.50),
                "p95_ms": percentile(ingest, 0.95),
                "runs": RUNS,
                "mean_ms": round(statistics.fmean(ingest), 2),
            },
            "search": {
                "p50_ms": percentile(search, 0.50),
                "p95_ms": percentile(search, 0.95),
                "runs": RUNS,
                "mean_ms": round(statistics.fmean(search), 2),
            },
            "show": {
                "p50_ms": percentile(show, 0.50),
                "p95_ms": percentile(show, 0.95),
                "runs": RUNS,
                "mean_ms": round(statistics.fmean(show), 2),
            },
        },
    }
    (PERF_DIR / "perf-result.json.tmp").write_text(
        json.dumps(result, ensure_ascii=False, indent=2),
        encoding="utf-8",
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
