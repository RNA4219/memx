# Incident Record Operations Spec

## 目的
- インシデント記録のテンプレート利用と実体運用の境界を固定し、`EVALUATION`/Task Seed/監査証跡の誤参照を防止する。

## 1. 利用境界（テンプレートと実体）
- テンプレート: `docs/IN-YYYYMMDD-001.md`
  - 雛形・記入例のみを扱う。
  - 実運用値（実日付/実ID/実時刻/実要件ID/実証跡パス）を記入してはならない。
  - `EVALUATION`、レビュー記録、Task Seed の `Source` に実証跡として参照してはならない。
- 実体: `docs/IN-<実日付>-<連番>.md`
  - 実インシデントの唯一の記録対象。
  - `EVALUATION`、レビュー記録、Task Seed `Source` の証跡として参照可能。
  - 起票時は必ず新規作成し、テンプレートファイル上書き・流用は禁止。

## 2. 実体ファイル作成時の必須項目（`docs/INCIDENT_TEMPLATE.md` 同期定義）
- 以下は `docs/INCIDENT_TEMPLATE.md` と同一定義で、`docs/IN-<実日付>-<連番>.md` に必須記録する。
  1. ID項目
     - `インシデントID`
     - `Task Seed 参照ID（インシデントIDと同値）`
  2. 時刻項目
     - `発生日（UTC）`
     - `起票日（UTC）`
     - `検知日時`
     - `暫定復旧完了日時`
     - `恒久復旧完了日時`
  3. 要件ID項目
     - `関連要件ID`
     - `waiver対象要件ID`（waiver 運用時）
  4. waiver項目（waiver 運用時）
     - `waiver理由`
     - `期限（UTC）`
     - `暫定リスク受容者`
     - `代替統制`
     - `解除条件`
  5. 証跡パス項目
     - `関連証跡パス`（例: `artifacts/ops/incident-summary.json`、`artifacts/ops/recovery-log.ndjson`）

## 3. 誤用防止ルール
- テンプレート（`docs/IN-YYYYMMDD-001.md`）に実運用値を入れない。
- 実体（`docs/IN-<実日付>-<連番>.md`）以外を `EVALUATION` 証跡として参照しない。
- `docs/IN-BASELINE.md` は補助資料であり、実インシデント証跡の代替にしない。

## 4. Incident → Task Seed 変換の固定項目
- Incident から Task Seed へ転記する際、`memx_spec_v3/docs/incident-to-task-traceability-spec.md` を正本とする。
- 固定必須フィールド:
  - `Source`: `docs/IN-<実日付>-<連番>.md#<対象章>`
  - `Requirements`: Incident 由来要件 + 関連要件ID/waiver対象要件ID（該当時）
  - `Commands`: Incident 要件の検証コマンド（1件以上）
  - `Dependencies`: Incident 由来の前提依存（無い場合 `- none`）

## 5. Design Docs Authoring（Phase 1/4）への接続
- Phase 1 の `source evidence` 採用条件:
  - `docs/IN-<実日付>-<連番>.md` のみ採用可。
  - 必須項目（ID/時刻/要件ID/waiver項目/証跡パス）に欠落がある場合は採用不可。
- Phase 4 の `source evidence` 最終判定条件:
  - gate 判定で Incident 証跡を使う場合、実体ファイルのみを根拠として記録する。
  - テンプレート/ベースライン参照が残る場合は `done` 遷移不可。
