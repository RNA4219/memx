---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
priority: high
---

# memx 運用仕様（operations-spec）

本書は `memx_spec_v3/docs/requirements.md` の運用要件を、運用実装時に参照しやすい形で要件ID単位に再編した補助仕様である。正本は `requirements.md`。

## 1. インシデント起票条件（要件ID単位）

| Requirement ID | 起票トリガ（必須条件） | 起票時の必須アクション |
| --- | --- | --- |
| `REQ-NFR-005` | `short→archive` 補償フローが障害検知から 30 分以内に収束しない（`pending_compensation_count > 0` を解消できない、または `S3` 到達）。 | `docs/IN-<実日付>-<連番>.md` を起票し、再計画チケットを発行する。 |
| `REQ-NFR-005` | 再処理が `最大2回` を超えても `archive+lineage` の整合が回復しない。 | 未収束として `docs/IN-*.md` 起票を必須化し、打ち切り理由を記録する。 |
| `REQ-NFR-006` | インシデント記録を監査証跡として扱う必要がある事象（セキュリティ/品質障害）。 | 受入対象は `docs/IN-<実日付>-<連番>.md` 形式のみとし、テンプレートID文書は実績証跡に使わない。 |

関連実行手順: [RUNBOOK 障害時手順（Detect/Retry/Rollback/Re-plan）](../../RUNBOOK.md#障害時手順要件id紐付け)

## 2. waiver 記録必須項目（`REQ-NFR-006`）

waiver を一時許容する場合でも記録媒体は `docs/IN-<実日付>-<連番>.md` に限定し、以下 7 項目を必須記録する。

1. waiver対象要件ID
2. waiver理由（技術的制約/外部依存/緊急運用の別）
3. 期限（UTC、失効日時）
4. 暫定リスク受容者（承認者）
5. 代替統制（監視強化・手動運用手順・追加検証）
6. 解除条件（どの証跡が揃えば waiver を解消するか）
7. 関連証跡パス（`artifacts/ops/incident-summary.json`、`artifacts/ops/recovery-log.ndjson` など）

関連実行手順: [RUNBOOK 障害時手順 3) 再計画（Re-plan）](../../RUNBOOK.md#3-再計画re-plan)

## 3. RTO/RPO 判定（`REQ-NFR-002`）

| Requirement ID | 判定ルール | 合否条件 |
| --- | --- | --- |
| `REQ-NFR-002` | `RTO` は `detected_at` から `mitigated_at`（暫定復旧）または `resolved_at`（恒久復旧）の先着時刻まで。 | `rto_minutes <= 30` |
| `REQ-NFR-002` | `RPO` は障害復旧後に再投入が必要だった最古データ時刻と `detected_at` の差分。 | `rpo_minutes <= 5` |
| `REQ-NFR-002` | 同一インシデントで RTO/RPO を同時判定。 | `rto_minutes <= 30` かつ `rpo_minutes <= 5` の同時成立を必須とする。 |

関連実行手順: [RUNBOOK 障害時手順 2) ロールバック（Rollback）](../../RUNBOOK.md#2-ロールバックrollback)

## 4. 補償フロー収束条件（`REQ-NFR-005`）

収束判定は次の 1〜5 を同時充足した場合のみ成立とする。

1. `pending_compensation_count == 0`
2. 各 `src_note_id` で `archive 実在 + archived_from lineage 実在` が `1組以上`
3. `dup_archive_count <= 1`（削減不能時は理由を `docs/IN-*.md` に記録）
4. `short_delete_ready_ratio == 1.0`
5. 障害検知から 30 分以内に上記を満たす（未達時は未収束としてインシデント起票）

状態遷移の終端条件:

- `S4`（収束完了）: 上記 1〜4 を満たす。
- `S5`（要起票終端）: `docs/IN-*.md` 起票 + 再計画チケット発行済み。

関連実行手順: [RUNBOOK 障害時手順 2) ロールバック（Rollback） / 3) 再計画（Re-plan）](../../RUNBOOK.md#障害時手順要件id紐付け)

## 5. 時系列フロー（検知→一次切り分け→緩和→復旧→事後レビュー）

運用は以下の時系列で固定し、各フェーズで証跡更新を必須化する。

1. 検知（Detect）
   - `detected_at` を `incident-summary.json` と `recovery-log.ndjson` の `detect` イベントへ同時記録。
2. 一次切り分け（Triage）
   - 影響範囲、暫定対応可否、再試行対象/非対象を判定。
   - `retry_count` 初期値を確定し、再試行上限（2回）超過見込みなら即時に再計画へ分岐。
3. 緩和（Mitigation）
   - 暫定復旧時刻を `mitigated_at` として記録。
   - `detected_at` 起点で 15 分以内を満たさない場合、未達理由を `docs/IN-*.md` に記録。
4. 復旧（Recovery）
   - 恒久復旧時刻 `resolved_at`、`rto_minutes`、`rpo_minutes` を記録。
   - `pending_compensation_count == 0` と `short_delete_ready_ratio == 1.0` を確認。
5. 事後レビュー（Postmortem）
   - `docs/IN-*.md` に 5 Whys・再発防止策・waiver有無/期限を確定記録。
   - 再計画チケットと復旧証跡パスを相互参照。

関連実行手順: [RUNBOOK 障害時手順（Detect/Retry/Rollback/Re-plan）](../../RUNBOOK.md#障害時手順要件id紐付け)

---

## 6. 必須証跡ファイル一覧とキー定義

運用判定に使う証跡は以下 3 種に固定する。

### 6.1 `artifacts/ops/incident-summary.json`

- 必須キー: `incident_id`, `detected_at`, `mitigated_at`, `resolved_at`, `rto_minutes`, `rpo_minutes`, `retry_count`
- キー定義:
  - `incident_id`: `IN-YYYYMMDD-XXX` 形式の識別子。
  - `detected_at` / `mitigated_at` / `resolved_at`: UTC ISO-8601。
  - `rto_minutes`: `detected_at` から `mitigated_at` または `resolved_at` までの分。
  - `rpo_minutes`: 復旧後に必要だった最古再投入時刻と `detected_at` の差分（分）。
  - `retry_count`: 実行した再試行回数（0〜2）。

### 6.2 `artifacts/ops/recovery-log.ndjson`

- 必須イベント: `detect`, `retry`, `rollback`（実施時）, `replan`（実施時）, `mitigate`, `resolve`
- 必須キー: `event`, `timestamp`, `pending_compensation_count`, `short_delete_ready_ratio`
- キー定義:
  - `event`: 上記必須イベントのみ使用。
  - `timestamp`: UTC ISO-8601。
  - `pending_compensation_count`: 未補償件数（収束時は 0）。
  - `short_delete_ready_ratio`: short delete 実行可能比率（収束時は 1.0）。

### 6.3 `docs/IN-*.md`

- 受入可能形式: `docs/IN-<実日付>-<連番>.md`
- 必須記録: 検知/緩和/復旧時刻、再試行回数、ロールバック有無、再計画チケット、waiver情報（適用時）

関連実行手順: [RUNBOOK インシデント対応要件](../../RUNBOOK.md#インシデント対応要件)

---

## 7. waiver 運用（発動条件・期限・解除条件・未解除時エスカレーション）

### 7.1 発動条件

- `REQ-NFR-005`/`REQ-NFR-006` の必須条件を期限内に満たせないことが確定した場合に限る。
- `docs/IN-<実日付>-<連番>.md` の起票と同時にのみ発動できる。

### 7.2 期限

- waiver には UTC 失効時刻を必須設定する。
- 期限未設定の waiver は無効扱い（fail）。

### 7.3 解除条件

- 解除条件として定義した証跡（`incident-summary.json` / `recovery-log.ndjson` / `docs/IN-*.md`）が揃った時点で即時解除する。
- 解除時は `docs/IN-*.md` に解除時刻・確認者・根拠証跡パスを追記する。

### 7.4 未解除時のエスカレーション

- 失効時刻到達時点で未解除の場合、`replan` を再発行し、次営業日内に責任者レビューを必須化する。
- 期限切れ waiver を残したままの運用継続は禁止する。

関連実行手順: [RUNBOOK 障害時手順 3) 再計画（Re-plan）](../../RUNBOOK.md#3-再計画re-plan)

---

## 8. 合否判定の正本参照（`EVALUATION.md` 固定）

- 合否判定の正本は常に `../../EVALUATION.md` とする。
- 本書（`operations-spec.md`）の閾値・判定ルールは運用参照用であり、最終判定は `EVALUATION.md` を優先する。
- 判定不一致時は `EVALUATION.md` を正として本書を同期修正する。

関連判定基準: [EVALUATION 正本](../../EVALUATION.md)

---

## 9. エラー契約変更時レビュー項目

エラー契約（`ErrorCode` / HTTP status / retryable）を変更する PR は、受け入れレビューで次のチェックポイントを必須確認する。

| チェックID | レビュー項目 | 合格条件 |
| --- | --- | --- |
| `REV-ERR-001` | 必須更新対象の同期 | `go/api/errors.go` / `go/service/errors.go` / `docs/error-contract.md` / `docs/interfaces.md` 4章が同一PRで更新され、差分説明がある。 |
| `REV-ERR-002` | retryable責務境界の遵守 | retryable 判定が service 層で定義され、transport 層が `ErrorCode -> HTTP` マッピングのみに限定されている。 |
| `REV-ERR-003` | 破壊的変更禁止の遵守 | `REQ-ERR-001`（既存 code 意味変更禁止）と `REQ-ERR-002`（既存 status 変更禁止）に抵触しない。 |
| `REV-ERR-004` | 契約テーブル整合 | `error-contract.md` の運用マトリクスとインタフェース仕様（`docs/interfaces.md` 4章）の code/status/retryable が一致する。 |

---

## 10. 契約差分チェック手順（レビュー必須）

契約変更を含む PR では、以下を**必須確認**とする。

1. 変更対象の特定  
   - `git diff -- memx_spec_v3/docs/contracts/openapi.yaml memx_spec_v3/docs/contracts/cli-json.schema.json` で正本スキーマ差分を確認する。

2. 変更順序の整合  
   - `requirements.md` → `contracts/*` → `interfaces.md` / `CONTRACTS.md` → `EVALUATION*` / `operations-spec.md` の順で更新されていることを確認する。

3. 重複定義の排除  
   - `CONTRACTS.md` が索引・抜粋のみで、フィールド定義の再定義を持たないことを確認する。

4. CLI/API 互換性の確認  
   - `interfaces.md` の互換方針（必須削除禁止・意味変更禁止・トップレベル構造変更禁止）に反していないことを確認する。

5. レビュー記録  
   - PR コメントまたはレビュー本文に「契約差分チェック実施済み」を明記する。