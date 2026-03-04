# Contract Alignment Report Spec

## 1. 目的
`contract-alignment-spec.md` の「4. 判定結果の出力フォーマット」を、Phase 3 実運用で直接使える記録テンプレートと運用手順に具体化する。

## 2. 正本保存先（正本化）
- 正本保存先は **`memx_spec_v3/docs/contracts/reports/`** とする。
- `memx_spec_v3/docs/reviews/` は設計レビュー用のため、契約同期判定レポートの正本には使わない。
- 既存で `docs/reviews/` に契約同期判定を保存している場合は、次回更新時に `docs/contracts/reports/` へ移し、旧ファイル側に移転先を追記する。

## 3. 必須成果物（固定）
Phase 3 の契約整合運用で作成する成果物は、次の 2 点に固定する。

1. `CONTRACT-ALIGN-YYYYMMDD-###.md`（詳細レポート）
2. `LATEST.md`（直近レポートIDと判定のポインタ）

補助運用ルール（命名規則・連番採番・最新レポート探索）は `memx_spec_v3/docs/contracts/reports/README.md` を正本とする。

## 4. 詳細レポート命名規則
- ファイル名は `CONTRACT-ALIGN-YYYYMMDD-###.md` とする。
  - `YYYYMMDD`: 判定実施日（ローカル日付）
  - `###`: 同日内連番（`001` 始まり）
- 例: `CONTRACT-ALIGN-20260304-001.md`

## 5. レポート最小記載項目
差分 1 件につき、以下 5 項目を必須とする。

- `diff_id`
- `severity`
- `affected_req`
- `source_refs`
- `action`

### 5.1 項目定義
- `diff_id`: 差分識別子。`CA-YYYYMMDD-###` を推奨。
- `severity`: `low | medium | high`。
- `affected_req`: 影響する REQ-ID（1 件以上）。
- `source_refs`: 判定根拠（ファイル行番号 or スキーマパス）。
- `action`: 解消方針・保留理由・移行手順。

## 6. 実運用テンプレート
以下を `CONTRACT-ALIGN-YYYYMMDD-###.md` 内に記録する。

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

`LATEST.md` には最低限、次を記載する。
- 最新 `report_id`
- 判定日
- `high` 件数
- Phase 3 判定（完了 / 未完了）
- 詳細レポートへの相対パス

上記は必須キー（`report_id/report_path/decision_date/high_count/phase3_status`）として固定し、契約整合レポート更新のたびに毎回更新する。

## 7. 受け入れ条件（Phase 3 完了条件）
- Phase 3 を完了扱いにできる条件は **`severity: high = 0`** のみ。
- `high` が 1 件以上ある場合、Phase 3 は未完了。

## 8. 統合運用手順（`CONTRACT-ALIGN` / `LATEST` / `EVALUATION`）
生成・更新は以下の順序で固定し、1つの手順として実施する。

1. **`CONTRACT-ALIGN-YYYYMMDD-###.md` を作成する**
   - 当日初回は `CONTRACT-ALIGN-<実日付>-001.md` を作成する。
   - `report_id` と `summary.high` を確定する。
2. **`LATEST.md` を更新する**
   - 必須キー `report_id/report_path/decision_date/high_count/phase3_status` を全て記録する。
   - 1キーでも未記載なら `fail`。
3. **`EVALUATION.md` を照合・更新する**
   - `EVALUATION.md` のレポートID参照と `LATEST.md.report_id` の一致を確認する。
   - `EVALUATION.md` のレポートパス参照と `LATEST.md.report_path` の一致を確認する。
   - ID 不一致時は Phase 3 ステータス遷移を `reviewing` 固定とし、`done` / `blocked` へ遷移しない。
4. **`docs/TASKS.md` の Phase 3 関連タスクを更新する**
   - `docs/TASKS.md` の「2-1-3. Phase 3（契約整合）チェック（必須）」をチェックリストとして実行する。
   - `high = 0` かつ `report_id/report_path` 一致: 状態を `done` に更新。
   - `high > 0`: 状態を `reviewing/blocked` に固定し、未解消 `diff_id` を列挙。
   - `report_id` または `report_path` 不一致: 状態を `reviewing` に固定。
5. **差分チェックリストを記録する**
   - `docs/TASKS.md` に差分単位のチェックリストを記載する。
   - 例: `- [ ] CA-20260304-003 (REQ-API-001): openapi 必須項目復元`

### 8.1 逆順更新の禁止
- 更新順序は **`CONTRACT-ALIGN作成 -> LATEST更新 -> EVALUATION照合`** のみを許可する。
- 逆順（`LATEST` 先行、`EVALUATION` 先行）の更新は運用違反として `fail`。

## 9. 保存・更新ルール
- 新規判定ごとに新規ファイルを作成し、既存レポートは上書きしない。
- `LATEST.md` は毎回更新し、常に最新レポートへのポインタを保つ。
- 追記修正時は同一ファイル内に `Updated:` 行を残し、差分履歴を保持する。
- 参照元は必ず実在パスを記載し、曖昧な記述（例: "関連箇所"）は禁止。
