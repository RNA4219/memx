# Contract Alignment Lifecycle Spec

## 1. 目的
Phase 3（契約整合）で扱う成果物と更新タイミング、`LATEST.md` 整合運用を正規化する。

## 2. 対象成果物（固定）
Phase 3 の契約整合ライフサイクルで扱う成果物は次の2点に固定する。

1. `memx_spec_v3/docs/contracts/reports/CONTRACT-ALIGN-YYYYMMDD-###.md`
2. `memx_spec_v3/docs/contracts/reports/LATEST.md`

## 3. 生成・更新トリガー
次のいずれかを満たした時点で、詳細レポートと `LATEST.md` を生成または更新する。

1. Phase 3 を実行したとき。
2. 契約差分で `high > 0` が検出されたとき。
3. 差分解消後の再判定（re-evaluation）を行ったとき。

## 4. `LATEST.md` 必須キー
`LATEST.md` には以下キーを必須で記載する。

```yaml
report_id: CONTRACT-ALIGN-YYYYMMDD-###
report_path: ./CONTRACT-ALIGN-YYYYMMDD-###.md
decision_date: YYYY-MM-DD
high_count: <number>
phase3_status: done | blocked | reviewing
```

- 上記5キーは schema として固定し、1つでも未記載なら Phase 3 チェックを `fail` とする。

## 5. 不整合時の復旧手順（優先順位ルール正規化）
`LATEST.md` と詳細レポートの不整合は、以下の優先順位で復旧する。

1. `LATEST.md` の `report_id` と `report_path` が実在し、内容整合する場合は `LATEST.md` を正本として採用する。
2. `LATEST.md` が欠落または不整合の場合は、`CONTRACT-ALIGN-*.md` を日付降順・連番降順で探索し、先頭を暫定最新として採用する。
3. 暫定最新を採用した場合、`LATEST.md` を即時再生成し、`report_id/report_path/decision_date/high_count/phase3_status` を再記録する。
4. 復旧後に `docs/TASKS.md` と `EVALUATION.md` のレポートID整合を再確認し、未整合なら Phase 3 を `reviewing` または `blocked` のまま維持する。

## 6. 統合運用手順（生成・更新順序固定）
`CONTRACT-ALIGN-<実日付>-001.md` と `LATEST.md` の生成・更新は、以下 3 ステップを 1 つの運用手順として固定する。

1. **`CONTRACT-ALIGN-YYYYMMDD-###.md` を作成/更新する**
   - 当日初回は `CONTRACT-ALIGN-<実日付>-001.md` を作成する。
   - `report_id` を確定し、`summary.high` を含む判定結果を記録する。
2. **`LATEST.md` を更新する**
   - `report_id/report_path/decision_date/high_count/phase3_status` の5キーを全て更新する。
   - 5キーの欠落、または `report_path` が非実在の場合は `fail`。
3. **`EVALUATION.md` を照合・更新する**
   - `EVALUATION.md` が参照するレポートIDと `LATEST.md.report_id` の一致を確認する。
   - ID 不一致時は Phase 3 ステータス遷移を `reviewing` 固定とし、`done` / `blocked` へは遷移させない。

### 6.1 逆順更新の禁止
- 更新順序は **`CONTRACT-ALIGN作成 -> LATEST更新 -> EVALUATION照合`** のみ許可する。
- `LATEST.md` 先行更新、または `EVALUATION.md` 先行更新は運用違反として `fail` とする。
