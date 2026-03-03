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

## LLM クライアント運用
- `EmbeddingClient`: 埋め込み生成。
- `MiniLLMClient`: タグ・スコア・機密度推定。
- `ReflectLLMClient`: Observer/Reflector 要約更新。
- タイムアウト 15 秒、最大 2 回リトライ（指数バックオフ）、再試行可/不可を区別して実装する。

## 関連ドキュメント
- エラー契約: `memx_spec_v3/docs/error-contract.md`

## 性能再計測手順（EVALUATION.md 同条件）
1. テストデータ投入（10,000 件 / 1件 約500文字）を実施。
2. 計測環境がローカル単体（4 vCPU / 16GB RAM / NVMe SSD / Linux x86_64）であることを確認。
3. ウォームアップとして各エンドポイントを 20 回実行。
4. 本計測として各エンドポイントを 200 回実行し、P50/P95 を算出。

### 実行前提
- 必須バイナリ: `go`(1.22+), `python3`(3.10+)
- 実行ディレクトリ: リポジトリルート（`/workspace/memx`）
- コマンド正本: リポジトリルート起点の `go run ./memx_spec_v3/go/cmd/mem ...`
- 代替表記: `cd memx_spec_v3/go` 後に `go run ./cmd/mem ...`（正本ではない）
- 入力データ形式: UTF-8 プレーンテキスト（1ノート約500文字、1行1ノートで生成）
- 計測対象コマンド実体:
  - ingest: `go run ./memx_spec_v3/go/cmd/mem in short`
  - search: `go run ./memx_spec_v3/go/cmd/mem out search`
  - show: `go run ./memx_spec_v3/go/cmd/mem out show`

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


## Birdseye 鮮度不足時の復旧手順
`docs/birdseye/index.json.generated_at` が判定時刻から7日を超える場合は、以下を順に実行する。

1. index を更新する。
```bash
python workflow-cookbook/tools/codemap/update.py --targets docs/birdseye/index.json --emit index
```
2. capsule を更新する。
```bash
python workflow-cookbook/tools/codemap/update.py --targets docs/birdseye/caps --emit caps
```
3. index/caps を再生成して再実行状態をそろえる。
```bash
python workflow-cookbook/tools/codemap/update.py --targets docs/birdseye/index.json,docs/birdseye/caps --emit index+caps
```

## Observability / 確認手順
1. 性能閾値の正本は `EVALUATION.md` とし、`governance/metrics.yaml` はその同期先として一致を維持する。
2. 日次確認では `response_time` / `compatibility` / `error_classification` / `recall_threshold` の breach 有無を確認する。
3. breach 発生時は `governance/metrics.yaml` の `action_on_breach` に従ってインシデントを起票する。

## リリース前確認（Release Drafter）
1. マージ済み PR に `feature` / `fix` / `chore` / `breaking` ラベルが正しく付与されていることを確認する。
2. GitHub の Releases 画面で Draft Release を開き、カテゴリ分類（Features/Fixes/Chores/Breaking Changes）とタイトルを確認する。
3. 誤分類や欠落がある場合は PR ラベルを修正し、Release Drafter の再実行でドラフトを更新する。
