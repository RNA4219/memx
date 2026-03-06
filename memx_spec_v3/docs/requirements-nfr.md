---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# memx 要求事項 - 非機能要件

> 本書は `requirements.md` から分割された一部です。正本は `requirements.md` を参照してください。

## 5. 非機能要件

### 5-1. 性能目標（v1必須3エンドポイント）

計測条件は以下で固定する。

- データセット条件: short ストア 10,000 件、本文 1 ノートあたり約 500 文字（UTF-8 プレーンテキスト）
- 実行環境: ローカル単体（4 vCPU / 16GB RAM / NVMe SSD / Linux x86_64）
- ウォームアップ: 各エンドポイント 20 回
- 本計測: 各エンドポイント 200 回（ウォームアップ除外）

| 操作 | API | p50 目標 | p95 目標 |
| --- | --- | --- | --- |
| ingest | `POST /v1/notes:ingest` | `<= 120ms` | `<= 250ms` |
| search | `POST /v1/notes:search` | `<= 80ms` | `<= 180ms` |
| show | `GET /v1/notes/{id}` | `<= 40ms` | `<= 90ms` |
- Requirement ID: `REQ-NFR-001`

#### 計測プロトコル

- データ生成条件（固定）:
  - 対象ストアは `short` 固定。
  - 事前投入データは 10,000 件固定。
  - 各ノート本文は約 500 文字（UTF-8 プレーンテキスト）固定。
- 計測回数（固定）:
  - 各エンドポイントごとにウォームアップ 20 回を先に実行。
  - 本計測は各エンドポイント 200 回で固定し、ウォームアップを集計から除外する。
- 除外条件（固定）:
  - ウォームアップ 20 回は母集団に含めない。
  - `REQ-NFR-001` 判定では `ingest` / `search` / `show` の 3 エンドポイント以外を除外する。
  - ローカル単体（4 vCPU / 16GB RAM / NVMe SSD / Linux x86_64）以外の環境で得た計測値は正式判定から除外する。
- 集計方法（固定）:
  - 本計測 200 回のレイテンシ分布から `p50` と `p95` を算出する。
  - 合否判定値は `p50_ms` / `p95_ms` のみとし、平均値や最大値は参考情報として扱う。
  - 出力フォーマットは `RUNBOOK.md` の `artifacts/perf/perf-result.json` に定義されたキーを正本とする。

### Dependencies

- `BLUEPRINT.md`
- `EVALUATION.md`
- `GUARDRAILS.md`
- `RUNBOOK.md`

- OS：ローカル（Linux / macOS / Windows）で動作する CLI を想定。
- DB：SQLite3（WAL モード、foreign_keys ON）。
- 言語：Go（単一バイナリビルドを前提）。
- 依存：
  - SQLite ドライバ（例：`modernc.org/sqlite` / `github.com/mattn/go-sqlite3`）。
  - CLI は標準 `flag` でもよい（薄いラッパが前提）。
  - HTTP API は標準 `net/http` でよい。
- セキュリティ：
  - APIキーや秘密情報は `memory_policy.yaml` + Gatekeeper により保存前にブロック。
- 拡張性：
  - `chronicle.db` / `memopedia.db` / `archive.db` のスキーマは、`short.db` と同様の `notes` / `tags` / `note_tags` / `note_embeddings` / `notes_fts` 構造をベースとする。
  - 将来、Working Memory／プロジェクトスコープ／セッションスコープを追加しても、既存の CLI と DB 構造を壊さない。
- トランザクションと一貫性：
  - ATTACH を跨いだ完全な原子的トランザクションは SQLite の仕様上保証できない。
- 設計上「データ喪失より重複を許容する」ポリシーとし、lineage による追跡・再蒸留で整合性を取り戻せるようにする。

### 5-2. 可用性・復旧・整合性回復（運用NFR）

| Requirement ID | 区分 | 要件 |
| --- | --- | --- |
| `REQ-NFR-002` | 可用性/復旧目標 | 障害時の目標復旧時間は `RTO <= 30分`、目標復旧時点は `RPO <= 5分` を満たすこと。 |
| `REQ-NFR-003` | 検知〜暫定復旧 | 障害検知から暫定復旧（サービス縮退を含む）完了までの上限時間を `15分` とする。 |
| `REQ-NFR-004` | 再処理上限 | 再試行/再処理は 1 リクエスト（または 1 ノート）あたり `最大 2 回` までとし、3 回目以降は自動再試行を禁止して運用エスカレーションする。 |
| `REQ-NFR-005` | 整合性回復時間 | `short→archive` 補償フローの整合性回復は、障害検知から `30分以内` に収束判定（または `docs/IN-*.md` 起票）へ到達すること。 |

補足:
- `RPO` は障害復旧後に再投入が必要なデータ欠損許容時間を指す。
- `RTO` は障害検知時点から正常系または暫定系サービス復帰までの許容時間を指す。
- 再処理方針は `at-least-once` を採用し、重複は許容するが欠損は許容しない（重複は `REQ-NFR-005` の収束条件で解消する）。

#### 5-2-1. RTO/RPO 判定の固定ルール

- `RTO` の起点は `detected_at`、終点は `mitigated_at`（暫定復旧）または `resolved_at`（恒久復旧）の先着時刻とする。
- `RPO` は障害復旧後に再投入が必要だった最古データ時刻と `detected_at` の差分で算出する。
- `REQ-NFR-002` の合否は、同一インシデントに対して `rto_minutes <= 30` かつ `rpo_minutes <= 5` の同時成立を必須とする。
- 判定証跡は `artifacts/ops/incident-summary.json` を正本とし、手計算値は参考情報扱いとする。

#### 5-2-2. 再試行方針（運用固定）

- 自動再試行の対象は一時障害（DB lock、LLM timeout、HTTP 429/502/503/504）のみとする。
- 再試行回数は 1 リクエスト（または 1 ノート）あたり最大 2 回、待機は指数バックオフ（`1s -> 2s`、ジッタ許容）を推奨値とする。
- `INVALID_ARGUMENT` / `NOT_FOUND` / `GATEKEEP_DENY` / 恒久障害の `INTERNAL` は再試行禁止とし、運用エスカレーションへ遷移する。
- 再試行打ち切り後は `rollback` または `replan` を必須化し、`docs/IN-*.md` に再試行回数と打ち切り理由を記録する。

### 5-3. 整合性回復要件（Archive 補償フロー）

- Requirement ID: `REQ-NFR-005`
- `mem gc short` の再実行で、`lineage` と `archive.notes` の不整合を `30分以内` に収束させる。
- 整合性回復の 1 サイクルで実施する再処理回数は `最大 2 回`（`REQ-NFR-004` 準拠）とし、未収束時は `docs/IN-*.md` 起票を必須とする。
- 収束条件は「`archive 実在 + 対応 lineage 実在` を満たし、short 側 Delete 判定が再開可能」であること。

収束判定（数値/状態）:
1. 対象バッチの `pending_compensation_count == 0`（未補償 0 件）。
2. 対象 `src_note_id` ごとに `archive 実在 + archived_from lineage 実在` が `1 組以上` 成立。
3. 同一 `src_note_id` の重複 archive は `dup_archive_count <= 1`（許容重複 1 件）まで削減済み、または削減不能理由を `docs/IN-*.md` に記録済み。
4. `short_delete_ready_ratio == 1.0`（Delete 判定対象の全件が削除可能状態）。
5. 上記 1〜4 を障害検知から 30 分以内に満たせない場合は「未収束」とし、`docs/IN-*.md` 起票と再計画チケット発行を必須とする。

状態定義（重複許容後の収束条件）:

| 状態ID | 名称 | 判定条件 | 次状態 |
| --- | --- | --- | --- |
| `S0` | 検知直後 | `pending_compensation_count > 0` | `S1` |
| `S1` | 補償実行中 | `retry_count <= 2` かつ `archive+lineage` の片系不足が存在 | `S2` or `S3` |
| `S2` | 重複許容安定 | `dup_archive_count >= 1` を許容しつつ `欠損=0` を維持 | `S4`（`pending_compensation_count==0` かつ `short_delete_ready_ratio==1.0` かつ `dup_archive_count<=1`）or `S3` |
| `S3` | 未収束 | `retry_count > 2` または 30分超過で `pending_compensation_count > 0` | `S5` |
| `S4` | 収束完了 | `pending_compensation_count == 0` かつ `short_delete_ready_ratio == 1.0` かつ `dup_archive_count <= 1` | 終端 |
| `S5` | 要起票終端 | `docs/IN-*.md` 起票 + 再計画チケット発行済み | 終端 |

- `S2` は「重複許容の暫定安定状態」とし、欠損ゼロを維持したまま `S4` へ収束させる中間状態とする。
- `S2` のまま 30 分を超過した場合は `S3`（未収束）へ遷移し、`S5` へ進める。
- `S2 -> S4` 遷移は「重複許容後の最終収束条件（未補償ゼロ・Delete 再開可・重複許容上限内）」を同時充足した場合に限定する。

### 5-4. インシデント記録（`docs/IN-*.md`）最小監査項目

- Requirement ID: `REQ-NFR-006`
- 受入対象の実運用インシデント記録は `docs/IN-<実日付>-<連番>.md` 形式のみとする。
- `docs/IN-*.md` には、以下の監査項目を必須記録する。
  1. 事象識別子: `インシデントID` / `発生日` / `起票日` / `重大度` / `ステータス`
  2. 要件トレーサビリティ: `関連要件ID` / `要件違反有無` / `違反した要件IDまたは節`
  3. 時間監査: `検知日時` / `暫定復旧完了日時` / `恒久復旧完了日時`
  4. 復旧行動監査: `再試行回数` / `ロールバック実施有無` / `再計画チケットID`
  5. 影響監査: `影響対象` / `影響期間` / `影響規模` / `CIA影響`
  6. 証跡: `関連ログ・メトリクス・判定結果ファイル` の保存先パス

#### 5-4-1. waiver 時の必須記録（`docs/IN-*.md` 運用連動）

- waiver を許容する場合でも、記録媒体は必ず `docs/IN-<実日付>-<連番>.md` とする（口頭/チャットのみは不可）。
- 必須項目の記載フォーマットは `docs/IN-BASELINE.md` / `docs/IN-YYYYMMDD-001.md` / `docs/IN-202603xx-001.md` の waiver セクションに準拠する。これら3文書はテンプレート資料であり、要件根拠としては参照禁止とする。
- `EVALUATION.md` の運用NFR合否判定で機械的に参照できること。
- 必須記録項目:
  1. waiver対象要件ID
  2. waiver理由（技術的制約/外部依存/緊急運用の別）
  3. 期限（UTC、失効日時）
  4. 暫定リスク受容者（承認者）
  5. 代替統制（監視強化・手動運用手順・追加検証）
  6. 解除条件（どの証跡が揃えば waiver を解消するか）
  7. 関連証跡パス（`artifacts/ops/incident-summary.json`、`artifacts/ops/recovery-log.ndjson` など）