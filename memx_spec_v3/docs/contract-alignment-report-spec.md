# Contract Alignment Report Spec

## 1. 目的
`contract-alignment-spec.md` の「4. 判定結果の出力フォーマット」を、Phase 3 実運用で直接使える記録テンプレートと運用手順に具体化する。

## 2. 正本保存先（正本化）
- 正本保存先は **`memx_spec_v3/docs/contracts/reports/`** とする。
- `memx_spec_v3/docs/reviews/` は設計レビュー用のため、契約同期判定レポートの正本には使わない。
- 既存で `docs/reviews/` に契約同期判定を保存している場合は、次回更新時に `docs/contracts/reports/` へ移し、旧ファイル側に移転先を追記する。

## 3. ファイル命名規則
- ファイル名は `CONTRACT-ALIGN-YYYYMMDD-###.md` とする。
  - `YYYYMMDD`: 判定実施日（ローカル日付）
  - `###`: 同日内連番（`001` 始まり）
- 例: `CONTRACT-ALIGN-20260304-001.md`

## 4. レポート最小記載項目
差分 1 件につき、以下 5 項目を必須とする。

- `diff_id`
- `severity`
- `affected_req`
- `source_refs`
- `action`

### 4.1 項目定義
- `diff_id`: 差分識別子。`CA-YYYYMMDD-###` を推奨。
- `severity`: `low | medium | high`。
- `affected_req`: 影響する REQ-ID（1 件以上）。
- `source_refs`: 判定根拠（ファイル行番号 or スキーマパス）。
- `action`: 解消方針・保留理由・移行手順。

## 5. 実運用テンプレート
以下を 1 レポートファイル内に列挙して記録する。

```yaml
report_id: CONTRACT-ALIGN-YYYYMMDD-###
phase: 3
summary:
  total: <number>
  low: <number>
  medium: <number>
  high: <number>
results:
  - diff_id: CA-YYYYMMDD-001
    severity: low | medium | high
    affected_req:
      - REQ-API-001
    source_refs:
      - memx_spec_v3/docs/design.md#L120-L150
      - memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1memories/get
    action: "文言統一のみ実施。互換性影響なし。"
```

## 6. 受け入れ条件（Phase 3 完了条件）
- Phase 3 を完了扱いにできる条件は **`severity: high = 0`** のみ。
- `high` が 1 件以上ある場合、Phase 3 は未完了。

## 7. `EVALUATION.md` / `docs/TASKS.md` 反映手順
1. レポート作成後、`EVALUATION.md` に以下を追記する。
   - 判定日
   - レポートID
   - `high` 件数
   - Phase 3 判定（完了 / 未完了）
2. `docs/TASKS.md` の Phase 3 関連タスクを更新する。
   - `high = 0`: 状態を `done` に更新
   - `high > 0`: 状態を `blocked` もしくは `reviewing` に更新し、未解消 `diff_id` を列挙
3. `docs/TASKS.md` には、対応タスクのチェックリストを差分単位で記載する。
   - 例: `- [ ] CA-20260304-003 (REQ-API-001): openapi 必須項目復元`
4. Phase 3 を `done` に更新する前に、`EVALUATION.md` と `docs/TASKS.md` のレポートIDが一致していることを確認する。

## 8. 保存・更新ルール
- 新規判定ごとに新規ファイルを作成し、既存レポートは上書きしない。
- 追記修正時は同一ファイル内に `Updated:` 行を残し、差分履歴を保持する。
- 参照元は必ず実在パスを記載し、曖昧な記述（例: "関連箇所"）は禁止。
