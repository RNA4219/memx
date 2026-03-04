---
intent_id: memx-operations-runbook-v1
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# RUNBOOK

## v1 必須運用フロー
1. **ingest**: `go run ./memx_spec_v3/go/cmd/mem in short`（または `POST /v1/notes:ingest`）で short へ投入。
2. **search**: `go run ./memx_spec_v3/go/cmd/mem out search`（または `POST /v1/notes:search`）で FTS 検索。
3. **show**: `go run ./memx_spec_v3/go/cmd/mem out show`（または `GET /v1/notes/{id}`）で単一ノート参照。

## `go run ./memx_spec_v3/go/cmd/mem in short` 手順
1. CLI で `file` または stdin から本文を読み込み、request を生成。
2. API へ request を送信（in-proc または HTTP）。
3. Service 側で Gatekeeper(`memory_store`)・ノート保存・タグ/埋め込み更新・`short_meta` 更新を実行。
4. CLI が note id 等を整形表示。

### CLI `--json` と API レスポンスの変換ルール（v1必須3エンドポイント）

契約の正本は `BLUEPRINT.md` の「v1 API Contract」とし、本節は `--json` 運用上の写像ルールを規定する。

現行実装では差分なし（CLI `--json` は API レスポンス JSON をそのまま出力）。

- `mem in short --json` ⇔ `POST /v1/notes:ingest`
  - 変換ルール: なし（`NotesIngestResponse` をそのまま表示）。
- `mem out search ... --json` ⇔ `POST /v1/notes:search`
  - 変換ルール: なし（`NotesSearchResponse` をそのまま表示）。
- `mem out show <id> --json` ⇔ `GET /v1/notes/{id}`
  - 変換ルール: なし（`Note` をそのまま表示）。

運用ルール:
- 互換性維持のため、上記3コマンドの `--json` は API と同一スキーマを維持する。
- 差分が必要な場合は、CLI 側で明示的なバージョンフラグを導入し、既定出力は維持する。

## `mem out recall` 手順
1. クエリを EmbeddingClient で埋め込み化。
2. 対象ストアの `note_embeddings` で類似度計算し、閾値以上を抽出。
3. 上位 `top-k` を anchor として `created_at` 前後 `range` を連結取得。
4. `--stores` 入力を正規化し、不正値は 400 系で失敗。
5. 埋め込みクライアント未設定時はデフォルトエラー、明示フラグ時のみ FTS フォールバック。

## `mem gc short` 手順
1. `short_meta` と `memory_policy.yaml.gc.short` からトリガ判定。
2. soft/hard limit 条件と `min_interval_minutes` で実行可否を決定。
3. 実行時は正確値を再計算して閾値を再確認。
4. `--dry-run` は DB 変更せず予定操作 JSON のみ返却。


## 設計更新完了時の最終確認手順
1. `memx_spec_v3/docs/design-deliverables-package-spec.md` の必須成果物表を上から順に確認し、対象タスクで `design.md` / `interfaces.md` / `traceability.md` / `memx_spec_v3/docs/reviews/*.md` / `CHANGELOG.md` / `memx_spec_v3/CHANGES.md` が更新済みかを点検する。
2. `docs/TASKS.md` の完了前チェック（`Release Note Draft` / `Status` / `Moved-to-CHANGES`）を順に実施し、Task Seed の記載と差分有無を照合する。
3. `orchestration/memx-design-docs-authoring.md` の Phase Done Criteria と照合し、最終判定に必要な成果物・証跡が不足していないことを確認する。
4. 判定ログに確認日時と確認者を記録し、`Status: done` 更新後に CHANGES 移送を完了する。

## 障害時手順（要件ID紐付け）

- 適用要件: `REQ-NFR-002` / `REQ-NFR-003` / `REQ-NFR-004` / `REQ-NFR-005` / `REQ-NFR-006`

### 障害時手順のNFR-ID対応表

| 手順 | 主対象NFR-ID | 判定に使う主要証跡 |
| --- | --- | --- |
| 検知・初期化（Detect） | `REQ-NFR-002`, `REQ-NFR-003` | `incident-summary.json.detected_at`, `recovery-log.ndjson.detect` |
| 再試行（Retry） | `REQ-NFR-003`, `REQ-NFR-004` | `recovery-log.ndjson.retry_count`, `incident-summary.json.mitigated_at` |
| ロールバック（Rollback） | `REQ-NFR-002`, `REQ-NFR-005` | `incident-summary.json.rto_minutes/rpo_minutes`, `recovery-log.ndjson.pending_compensation_count` |
| 再計画（Re-plan） | `REQ-NFR-003`, `REQ-NFR-005`, `REQ-NFR-006` | `docs/IN-*.md`, `recovery-log.ndjson.replan_ticket_id` |

### 0) 検知・初期化（Detect）
- 要件紐付け: `REQ-NFR-002`, `REQ-NFR-003`
- 障害検知時点を `incident-summary.json.detected_at` と `recovery-log.ndjson` の `detect` イベントで固定記録する。
- 後続の RTO/RPO 判定、15分暫定復旧判定はこの `detected_at` を起点に算出する。

### 1) 再試行（Retry）
- 要件紐付け: `REQ-NFR-003`, `REQ-NFR-004`
- 初動で一時障害を判定した場合のみ再試行を実施する。
- 1 リクエスト（または 1 ノート）あたり再試行は最大 2 回までとし、3 回目は実施しない。
- 障害検知から暫定復旧まで 15 分を超過する見込みの場合、再試行を打ち切ってロールバック/再計画へ移行する。

### 2) ロールバック（Rollback）
- 要件紐付け: `REQ-NFR-002`, `REQ-NFR-005`
- `archive 実在 + 対応 lineage 実在` を満たさない状態では short 側 Delete を禁止し、ロールバックしてデータ喪失を回避する。
- 復旧手順の実行後、`RTO <= 30分` / `RPO <= 5分` を `incident-summary.json` で確認する。
- 補償収束の確認として `pending_compensation_count == 0` と `short_delete_ready_ratio == 1.0` を `recovery-log.ndjson` で確認する。

### 3) 再計画（Re-plan）
- 要件紐付け: `REQ-NFR-003`, `REQ-NFR-005`, `REQ-NFR-006`
- 再試行上限到達または 30 分以内に収束しない場合、`docs/IN-*.md` を起票して再計画チケットを発行する。
- `docs/IN-*.md` には検知時刻/暫定復旧完了時刻/恒久復旧完了時刻・再試行回数・ロールバック実施有無を必須記録する。
- waiver を適用する場合は `REQ-NFR-006` として、`docs/IN-<実日付>-<連番>.md` に waiver対象要件ID/理由/期限/承認者/代替統制/解除条件/証跡パスを必須記録する。

## LLM クライアント運用
- `EmbeddingClient`: 埋め込み生成。
- `MiniLLMClient`: タグ・スコア・機密度推定。
- `ReflectLLMClient`: Observer/Reflector 要約更新。
- タイムアウト 15 秒、最大 2 回リトライ（指数バックオフ）、再試行可/不可を区別して実装する。

## エラーコード別 再試行戦略（運用）
- 即時再試行（1 回のみ）:
  - 直前の接続瞬断や軽微なネットワーク揺らぎで、`INTERNAL` かつ一時障害と判定できる場合。
- 指数バックオフ再試行（最大 2 回）:
  - `INTERNAL` で原因が DB ロック、LLM タイムアウト、HTTP 429/502/503/504 の場合。
  - 推奨待機: 1s → 2s（ジッタ許容）。
- 再試行禁止:
  - `INVALID_ARGUMENT`（入力不備）、`NOT_FOUND`（対象不存在）、`GATEKEEP_DENY`（ポリシー deny）。
  - 恒久的な `INTERNAL`（設定不備・スキーマ不整合等）は再試行せず原因修正を優先。

## 関連ドキュメント
- エラー契約: `memx_spec_v3/docs/error-contract.md`

## 要件トレーサビリティ用 検証コマンド（正本）

### trace-req-* と I/F 項目ID 対応
| trace anchor | requirement | 対応コマンド | I/F 項目ID |
| --- | --- | --- | --- |
| `trace-req-cli-001` | `REQ-CLI-001` | `mem out search` | `IF-CLI-SEARCH-REQ`, `IF-CLI-SEARCH-RES` |
| `trace-req-api-001` | `REQ-API-001` | `mem in short` | `IF-CLI-INGEST-REQ`, `IF-API-INGEST-REQ`, `IF-API-INGEST-RES` |
| `trace-req-gc-001` | `REQ-GC-001` | `mem gc short` | `IF-GC-SHORT-REQ`, `IF-GC-SHORT-RES` |
| `trace-req-sec-001` | `REQ-SEC-001` | `mem in short` | `IF-CLI-INGEST-REQ`, `IF-API-INGEST-REQ` |
| `trace-req-err-001` | `REQ-ERR-001` | `mem out show 999999` | `IF-CLI-SHOW-REQ`, `IF-ERR-MATRIX` |

<a id="trace-lint"></a>

### lint（ruff）
```bash
python3 -m ruff check .
```

<a id="trace-type"></a>

### type（mypy/strict）
```bash
python3 -m mypy --strict memx_spec_v3
```

<a id="trace-test-pytest"></a>

### test（pytest）
```bash
python3 -m pytest -q
```

<a id="trace-test-node"></a>

### test（node:test）
```bash
node --test
```

<a id="trace-manual"></a>

### manual（CLI/API/GC/Security/Error）
```bash
# CLI/API ingest
printf '%s' 'traceability-sample' | go run ./memx_spec_v3/go/cmd/mem in short --stdin --title traceability --api-url http://127.0.0.1:7766

# CLI/API search
go run ./memx_spec_v3/go/cmd/mem out search 'traceability' --api-url http://127.0.0.1:7766

# CLI/API show
go run ./memx_spec_v3/go/cmd/mem out show 1 --api-url http://127.0.0.1:7766

# GC dry-run (flag ON 想定)
go run ./memx_spec_v3/go/cmd/mem gc short --dry-run --api-url http://127.0.0.1:7766
# GC API (flag OFF 想定: route公開時は INTERNAL, 非公開時は NOT_FOUND)
curl -sS -X POST http://127.0.0.1:7766/v1/gc:run -H "content-type: application/json" -d '{"target":"short","options":{"dry_run":true}}'
```

<a id="trace-req-cli-001"></a>

### REQ-CLI-001 検証コマンド（IF: IF-CLI-SEARCH-REQ / IF-CLI-SEARCH-RES）
```bash
go run ./memx_spec_v3/go/cmd/mem out search 'traceability' --api-url http://127.0.0.1:7766
```

<a id="trace-req-api-001"></a>

### REQ-API-001 検証コマンド（IF: IF-CLI-INGEST-REQ / IF-API-INGEST-REQ / IF-API-INGEST-RES）
```bash
go run ./memx_spec_v3/go/cmd/mem in short --stdin --title traceability --api-url http://127.0.0.1:7766
```

<a id="trace-req-gc-001"></a>

### REQ-GC-001 検証コマンド（IF: IF-GC-SHORT-REQ / IF-GC-SHORT-RES）
```bash
# Case 1: flag ON（dry-run 成功）
go run ./memx_spec_v3/go/cmd/mem gc short --dry-run --api-url http://127.0.0.1:7766

# Case 2: flag OFF（期待: route非公開なら NOT_FOUND / route公開なら INTERNAL）
curl -i -sS -X POST http://127.0.0.1:7766/v1/gc:run -H "content-type: application/json" -d '{"target":"short","options":{"dry_run":true}}'
```

<a id="trace-req-sec-001"></a>

### REQ-SEC-001 検証コマンド（IF: IF-CLI-INGEST-REQ / IF-API-INGEST-REQ）
```bash
printf '%s' 'secret-token-for-trace' | go run ./memx_spec_v3/go/cmd/mem in short --stdin --title trace-sec --api-url http://127.0.0.1:7766
```

<a id="trace-req-err-001"></a>

### REQ-ERR-001 検証コマンド（IF: IF-CLI-SHOW-REQ / IF-ERR-MATRIX）
```bash
go run ./memx_spec_v3/go/cmd/mem out show 999999 --api-url http://127.0.0.1:7766
```

<a id="trace-perf"></a>

## 性能再計測手順（EVALUATION.md 同条件）
1. テストデータ投入（10,000 件 / 1件 約500文字）を実施。
2. 計測環境がローカル単体（4 vCPU / 16GB RAM / NVMe SSD / Linux x86_64）であることを確認。
3. ウォームアップとして各エンドポイントを 20 回実行。
4. 本計測として各エンドポイントを 200 回実行し、P50/P95 を算出。

### 実行前提
- 必須バイナリ: `go`(1.22+), `python3`(3.10+)
- 実行ディレクトリ: リポジトリルート（`/workspace/memx`）
- コマンド正本: リポジトリルート起点の `go run ./memx_spec_v3/go/cmd/mem ...` のみを使用する
- 入力データ形式: UTF-8 プレーンテキスト（1ノート約500文字、1行1ノートで生成）
- 計測対象コマンド実体:
  - ingest: `go run ./memx_spec_v3/go/cmd/mem in short`
  - search: `go run ./memx_spec_v3/go/cmd/mem out search`
  - show: `go run ./memx_spec_v3/go/cmd/mem out show`

### 計測コマンドと出力 JSON 項目の対応（正本）

| 計測対象 | 正本コマンド | `artifacts/perf/perf-result.json` の対応項目 |
| --- | --- | --- |
| ingest | `go run ./memx_spec_v3/go/cmd/mem in short --stdin --title <title> --api-url http://127.0.0.1:7766` | `results.ingest.p50_ms` / `results.ingest.p95_ms` / `results.ingest.runs` |
| search | `go run ./memx_spec_v3/go/cmd/mem out search '<query>' --api-url http://127.0.0.1:7766` | `results.search.p50_ms` / `results.search.p95_ms` / `results.search.runs` |
| show | `go run ./memx_spec_v3/go/cmd/mem out show <id> --api-url http://127.0.0.1:7766` | `results.show.p50_ms` / `results.show.p95_ms` / `results.show.runs` |

### 出力 JSON スキーマ（`artifacts/perf/perf-result.json`）
```json
{
  "environment": {
    "cpu": "4 vCPU",
    "ram": "16GB",
    "storage": "NVMe SSD",
    "os": "Linux x86_64"
  },
  "dataset": {
    "store": "short",
    "note_count": 10000,
    "body_chars": 500
  },
  "results": {
    "ingest": { "p50_ms": 0.0, "p95_ms": 0.0, "runs": 200 },
    "search": { "p50_ms": 0.0, "p95_ms": 0.0, "runs": 200 },
    "show": { "p50_ms": 0.0, "p95_ms": 0.0, "runs": 200 }
  }
}
```

### ローカル再現ベンチ（正本コマンド）
```bash
# API 起動（別ターミナル）
go run ./memx_spec_v3/go/cmd/mem api serve --addr 127.0.0.1:7766 --short ./artifacts/perf/short.db

# 同条件データ生成 + warmup(20) + 本計測(200) + 結果出力
# 下段の「実行コマンド例」をそのまま実行する（ローカル再現の正本手順）。
```

### 結果記録フォーマット（最低限）
- 最低限、`artifacts/perf/perf-result.json` に以下を必須記録する。
  - `results.ingest.p50_ms` / `results.ingest.p95_ms`
  - `results.search.p50_ms` / `results.search.p95_ms`
  - `results.show.p50_ms` / `results.show.p95_ms`
- `p50` / `p95` が欠損している結果は `REQ-NFR-001` 判定に使用しない。

### 実行コマンド例
```bash
mkdir -p artifacts/perf
rm -f artifacts/perf/short.db

# API サーバー起動（別ターミナル）
go run ./memx_spec_v3/go/cmd/mem api serve \
  --addr 127.0.0.1:7766 \
  --short ./artifacts/perf/short.db

# 1) データ投入（同条件データセット: 10,000件 / 500文字）
python3 - <<'PY'
import json
import subprocess
from pathlib import Path

out = Path("artifacts/perf/seed-result.json")
out.parent.mkdir(parents=True, exist_ok=True)
body = "あ" * 500
ids = []
for i in range(10000):
    cmd = (
        f"printf '%s' '{body}' | "
        f"go run ./memx_spec_v3/go/cmd/mem in short --stdin --title perf-{i} "
        f"--api-url http://127.0.0.1:7766"
    )
    r = subprocess.run(cmd, shell=True, check=True, capture_output=True, text=True)
    ids.append(r.stdout.strip())
out.write_text(json.dumps({"count": len(ids), "note_ids": ids[:5]}, ensure_ascii=False, indent=2), encoding="utf-8")
PY

# 2) ウォームアップ + 3) 本計測（JSON 出力）
python3 - <<'PY'
import json
import statistics
import subprocess
import time
from pathlib import Path

API = "http://127.0.0.1:7766"
WARMUP = 20
RUNS = 200
query = "あ"
known_id = "1"

def ms(cmd: str) -> float:
    t0 = time.perf_counter()
    subprocess.run(cmd, shell=True, check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    return (time.perf_counter() - t0) * 1000

def pct(values, p):
    values = sorted(values)
    i = int((len(values) - 1) * p)
    return round(values[i], 2)

for _ in range(WARMUP):
    ms(f"printf '%s' 'warmup' | go run ./memx_spec_v3/go/cmd/mem in short --stdin --title warmup --api-url {API}")
    ms(f"go run ./memx_spec_v3/go/cmd/mem out search '{query}' --api-url {API}")
    ms(f"go run ./memx_spec_v3/go/cmd/mem out show {known_id} --api-url {API}")

ingest = [ms(f"printf '%s' 'bench' | go run ./memx_spec_v3/go/cmd/mem in short --stdin --title bench --api-url {API}") for _ in range(RUNS)]
search = [ms(f"go run ./memx_spec_v3/go/cmd/mem out search '{query}' --api-url {API}") for _ in range(RUNS)]
show = [ms(f"go run ./memx_spec_v3/go/cmd/mem out show {known_id} --api-url {API}") for _ in range(RUNS)]

Path("artifacts/perf/warmup-result.json").write_text(json.dumps({"warmup": WARMUP}, indent=2), encoding="utf-8")
Path("artifacts/perf/perf-result.json").write_text(json.dumps({
    "environment": {"cpu": "4 vCPU", "ram": "16GB", "storage": "NVMe SSD", "os": "Linux x86_64"},
    "dataset": {"store": "short", "note_count": 10000, "body_chars": 500},
    "results": {
        "ingest": {"p50_ms": pct(ingest, 0.50), "p95_ms": pct(ingest, 0.95), "runs": RUNS},
        "search": {"p50_ms": pct(search, 0.50), "p95_ms": pct(search, 0.95), "runs": RUNS},
        "show": {"p50_ms": pct(show, 0.50), "p95_ms": pct(show, 0.95), "runs": RUNS}
    }
}, ensure_ascii=False, indent=2), encoding="utf-8")
PY
```

### 出力保存先
- シード結果: `artifacts/perf/seed-result.json`
- ウォームアップ結果: `artifacts/perf/warmup-result.json`
- 本計測結果（P50/P95 含む）: `artifacts/perf/perf-result.json`

## インシデント/不具合起票
- 初期運用ベースライン: [`docs/IN-BASELINE.md`](docs/IN-BASELINE.md)
- インシデント記録テンプレート: [`docs/INCIDENT_TEMPLATE.md`](docs/INCIDENT_TEMPLATE.md)
- 不具合起票時は GitHub Issue テンプレートを使用する: [.github/ISSUE_TEMPLATE/bug.yml](.github/ISSUE_TEMPLATE/bug.yml)
- 再現手順・期待値/実際値・影響範囲・関連 Intent ID を必ず記入する。

## 設計書作成前の静的検査手順（最短）
参照解決検証は `memx_spec_v3/docs/design-reference-validation-automation-spec.md` を正本とし、次を順に実行する。

1. `memx_spec_v3/docs/EVALUATION.md` 誤参照の検出。
```bash
rg -n "memx_spec_v3/docs/EVALUATION\.md" TASK.*.md orchestration/*.md memx_spec_v3/docs/design-*.md
```
2. `path#section` 未指定参照の検出（`Source:`/`Dependencies:` 行）。
```bash
rg -n "^(Source|Dependencies):\s+[^#\s]+$" TASK.*.md orchestration/*.md memx_spec_v3/docs/design-*.md
```
3. テンプレート ID (`IN-YYYYMMDD-001`) 混入の検出。
```bash
rg -n "IN-YYYYMMDD-001" TASK.*.md orchestration/*.md memx_spec_v3/docs/design-*.md
```
4. `contracts.md` 表記ゆれ（未解決）の検出。
```bash
rg -n "memx_spec_v3/docs/contracts\.md|\bcontracts\.md\b" TASK.*.md orchestration/*.md memx_spec_v3/docs/design-*.md
```
5. Birdseye 側（node/caps）は別責務として実行。
```bash
python workflow-cookbook/tools/codemap/update.py --targets docs/birdseye/index.json,docs/birdseye/caps --emit index+caps
```

各コマンドの出力が 0 件であることを pass 条件とし、1 件でも検出した場合は設計書作成着手を停止する。


## Birdseye 鮮度不足時の復旧手順
Birdseye 初期理解ドキュメントの canonical path は `docs/birdseye/README.md`（HUB と同一）とする。
`docs/birdseye/index.json.generated_at` が判定時刻から7日を超える場合に実施する。詳細は `./workflow-cookbook/tools/codemap/README.md`（別名パス `tools/codemap/README.md` は注記扱い）の「Birdseyeアクセス異常時の復旧手順」を正本とし、ここでは canonical path `workflow-cookbook/tools/codemap/update.py` の実行順のみ示す。

運用モード（HUB と同一語彙）:
- 現行運用（既定）: `docs/birdseye/hot.json` は `optional`（未運用）。再生成対象に含めない。
- 移行運用（`hot.json` 導入時のみ）: `docs/birdseye/hot.json` を生成対象に含める。
- `hot.json` 欠損時は `notes.readiness_status=ready` のまま継続し、`notes.missing_files` に `docs/birdseye/hot.json` を記録する。

1. index を更新する。
```bash
python workflow-cookbook/tools/codemap/update.py --targets docs/birdseye/index.json --emit index
```
2. capsule を更新する。
```bash
python workflow-cookbook/tools/codemap/update.py --targets docs/birdseye/caps --emit caps
```
3. index/caps を再生成して再実行状態をそろえる（現行運用）。
```bash
python workflow-cookbook/tools/codemap/update.py --targets docs/birdseye/index.json,docs/birdseye/caps --emit index+caps
```
4. `hot.json` 導入時のみ、hot を追加生成する（移行運用）。
```bash
python workflow-cookbook/tools/codemap/update.py --targets docs/birdseye/hot.json --emit hot
```

再発時（`nodes[].capsule` 欠損が再検知された場合）は、Task Seed を `Status: blocked` へ遷移し、運用語彙と必須 notes は [`HUB.codex.md` の「Birdseye Readiness Check」](HUB.codex.md#birdseye-readiness-check) と [`docs/TASKS.md` の「Birdseye Readiness Check」](docs/TASKS.md#birdseye-readiness-check) に従う。

鮮度更新時に `memx_spec_v3/docs/requirements.md` の節構成・Requirement ID（例: `REQ-CLI-001`）を変更した場合は、`docs/birdseye/index.json` の `nodes[].node_id=requirements` と `docs/birdseye/caps/requirements.json`（必要に応じて `docs/birdseye/caps/memx_spec_v3__docs__requirements.md.json`）の `summary`/`depends_on` を同一コミットで更新し、`python workflow-cookbook/tools/codemap/update.py --targets docs/birdseye/index.json,docs/birdseye/caps --emit index+caps` を再実行して要件ノード紐付けを鮮度管理フローへ組み込む。

## Observability / 確認手順
1. 性能閾値の正本は `EVALUATION.md` とし、`governance/metrics.yaml` の `ingest` / `search` / `show` の項目名・閾値を完全一致で同期維持する。
2. 日次確認では `ingest` / `search` / `show` / `compatibility` / `error_classification` / `recall_threshold` の breach 有無を確認する。
3. breach 発生時は `governance/metrics.yaml` の `action_on_breach` に従ってインシデントを起票する。

## リリース前確認（Release Drafter）
1. `git ls-files -- memx_spec_v3/go/go.mod memx_spec_v3/go/go.sum` を実行し、`memx_spec_v3/go/go.mod` と `memx_spec_v3/go/go.sum` の追跡状態を確認する。
2. マージ済み PR に `feature` / `fix` / `chore` / `breaking` ラベルが正しく付与されていることを確認する。
3. GitHub の Releases 画面で Draft Release を開き、カテゴリ分類（Features/Fixes/Chores/Breaking Changes）とタイトルを確認する。
4. 誤分類や欠落がある場合は PR ラベルを修正し、Release Drafter の再実行でドラフトを更新する。
