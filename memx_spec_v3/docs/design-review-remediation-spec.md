# Design Review Remediation Spec

## 1. 目的
本仕様は `memx_spec_v3/docs/design-review-spec.md` を補完し、指摘の是正運用（状態遷移・エビデンス・Task Seed 連携・完了判定）を固定する。

## 2. 指摘ライフサイクル

### 2.1 状態定義
- `open`: 指摘登録済みで、是正未着手または是正内容未提出。
- `remaining`: 是正内容は提出済みだが、レビュー再確認未完了または不十分。
- `resolved`: 是正内容と根拠エビデンスが受理され、再確認が完了。
- `waived`: 正式な waiver 根拠により是正不要として承認。

### 2.2 状態遷移条件（必須）
- `open -> remaining`
  - 是正差分が存在し、指摘単位の Task Seed 記載（Source/Requirements/Status）が完了していること。
- `remaining -> resolved`
  - 必須エビデンスが重大度ルールを満たしていること。
  - `major` は再レビュー完了を必須とする。
  - `minor` は自己確認完了を必須とし、レビュアー確認はプロジェクト運用ルールで必須化される場合のみ required とする。
- `open|remaining -> waived`
  - `EVALUATION.md` の waiver 条件に整合し、`docs/IN-<実日付>-<連番>.md` 参照を記録していること。
- `resolved|waived -> open`
  - 追補レビューで再発・根拠欠落・関連差分の破壊変更が検知された場合。

## 3. 重大度別の必須エビデンス

### 3.1 major
- 必須:
  1. 是正差分参照（コミット/PR差分）
  2. 自己確認記録（実行コマンドと結果）
  3. **再レビュー記録（必須）**
- 判定:
  - 3点が揃わない限り `resolved` へ遷移不可。

### 3.2 minor
- 必須:
  1. 是正差分参照（コミット/PR差分）
  2. 自己確認記録（実行コマンドと結果）
- レビュアー確認:
  - デフォルトは optional。
  - ただし `design-review-spec.md` の判定で再確認指定が付与された場合は required。

## 4. Task Seed 連携

### 4.1 指摘ごとの必須記載
各指摘（例: `DR-001`）に対して、対応 Task Seed に以下 3 項目を必須記載する。
- `Source`: 指摘ID・レビューID（`DESIGN-REVIEW-YYYYMMDD-###`）
- `Requirements`: 対応 `REQ-*` 一覧
- `Status`: `open|remaining|resolved|waived` の現在値

### 4.2 Moved-to-CHANGES 反映条件
- 指摘単位の `Status` が `resolved` または `waived` となり、レビュー全体判定が `pass` または `waiver` で確定した時点で `Moved-to-CHANGES: YYYY-MM-DD` を追記する。
- `remaining` が 1 件でも残る場合は追記不可。

## 5. 完了判定（DESIGN-REVIEW / DESIGN-ACCEPTANCE 整合）
レビュー是正の完了判定は、`DESIGN-REVIEW-*` と `DESIGN-ACCEPTANCE-*` の整合チェックを全件満たした場合のみ成立する。

### 5.1 整合チェック項目（必須）
1. **判定整合**
   - `DESIGN-REVIEW-*` の最終判定（pass/fail/waiver）と、`DESIGN-ACCEPTANCE-*` の最終判定（pass/fail）に矛盾がないこと。
2. **指摘集約整合**
   - `major` の未解消（`open|remaining`）が 0 件であること。
   - `minor` の未解消が残る場合は `waived` 根拠または再確認計画が記録されていること。
3. **REQ・章参照整合**
   - `DESIGN-REVIEW-*` の対象章・REQ-ID が `DESIGN-ACCEPTANCE-*` の集約対象に包含されていること。
4. **エビデンス参照整合**
   - `evidence_paths` / コマンド結果 / コメントID の参照先が双方で追跡可能であること。
5. **Task Seed整合**
   - 指摘に紐づく Task Seed の `Status` とレビュー記録中の状態が一致し、`Moved-to-CHANGES` の有無が条件と整合すること。

## 6. 参照
- `memx_spec_v3/docs/design-review-spec.md`
- `memx_spec_v3/docs/design-acceptance-report-spec.md`
- `memx_spec_v3/docs/design-review-artifact-sla-spec.md`
- `memx_spec_v3/docs/design-doc-dod-spec.md`
