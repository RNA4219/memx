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

## 5. 不整合時の復旧手順（優先順位ルール正規化）
`LATEST.md` と詳細レポートの不整合は、以下の優先順位で復旧する。

1. `LATEST.md` の `report_id` と `report_path` が実在し、内容整合する場合は `LATEST.md` を正本として採用する。
2. `LATEST.md` が欠落または不整合の場合は、`CONTRACT-ALIGN-*.md` を日付降順・連番降順で探索し、先頭を暫定最新として採用する。
3. 暫定最新を採用した場合、`LATEST.md` を即時再生成し、`report_id/report_path/decision_date/high_count/phase3_status` を再記録する。
4. 復旧後に `docs/TASKS.md` と `EVALUATION.md` のレポートID整合を再確認し、未整合なら Phase 3 を `reviewing` または `blocked` のまま維持する。

