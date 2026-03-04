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

## 2. waiver 記録必須項目（`REQ-NFR-006`）

waiver を一時許容する場合でも記録媒体は `docs/IN-<実日付>-<連番>.md` に限定し、以下 7 項目を必須記録する。

1. waiver対象要件ID
2. waiver理由（技術的制約/外部依存/緊急運用の別）
3. 期限（UTC、失効日時）
4. 暫定リスク受容者（承認者）
5. 代替統制（監視強化・手動運用手順・追加検証）
6. 解除条件（どの証跡が揃えば waiver を解消するか）
7. 関連証跡パス（`artifacts/ops/incident-summary.json`、`artifacts/ops/recovery-log.ndjson` など）

## 3. RTO/RPO 判定（`REQ-NFR-002`）

| Requirement ID | 判定ルール | 合否条件 |
| --- | --- | --- |
| `REQ-NFR-002` | `RTO` は `detected_at` から `mitigated_at`（暫定復旧）または `resolved_at`（恒久復旧）の先着時刻まで。 | `rto_minutes <= 30` |
| `REQ-NFR-002` | `RPO` は障害復旧後に再投入が必要だった最古データ時刻と `detected_at` の差分。 | `rpo_minutes <= 5` |
| `REQ-NFR-002` | 同一インシデントで RTO/RPO を同時判定。 | `rto_minutes <= 30` かつ `rpo_minutes <= 5` の同時成立を必須とする。 |

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
