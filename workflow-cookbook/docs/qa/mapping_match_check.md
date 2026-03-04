---
intent_id: DOC-LEGACY
owner: docs-core
status: active
last_reviewed_at: 2026-03-04
next_review_due: 2026-04-04
---

# Mapping Match Check

## 判定基準

全 REQ について `design_ref` / `contract_ref` / `test_ref` が `docs/qa/mapping_matrix.yaml` に存在すること。

## 運用ルール

- `gap` セクションの件数が 0 件なら `pass`。
- `gap` セクションに 1 件以上あれば `fail`。
- gap は「REQ に対して design/contract/test のいずれかが未接続、または参照先が無効」の項目を列挙する。

## 判定結果

- result: **pass**
- assessed_at: 2026-03-04
- assessed_file: `docs/qa/mapping_matrix.yaml`

## gap

- 0 件
