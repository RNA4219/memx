# Requirements Coverage Spec

## 1. 目的

本仕様は、`memx_spec_v3/docs/requirements.md` に定義された要件と  
設計関連成果物の対応関係を検証し、要件ID網羅率（coverage）を算出する手順を定義する。

本仕様は以下の目的で使用する。

- 受け入れレビューで使用する要件網羅率の算出
- Phase 3 完了判定の一意化
- 設計トレーサビリティの完全性検証

---

# 2. 入力

以下の設計成果物を入力とする。

- `memx_spec_v3/docs/requirements.md`
- `memx_spec_v3/docs/traceability.md`
- `memx_spec_v3/docs/design.md`
- `memx_spec_v3/docs/interfaces.md`

---

# 3. 用語定義

### total_req（母数）

`requirements.md` に存在する有効 `REQ-*` の総数

### mapped_req（分子）

`traceability.md` において対象 `REQ-*` が存在し、  
以下 5 列がすべて非空である件数

```

Source
Design Mapping
Interface Mapping
Evaluation Mapping
Contract Mapping

```

### coverage（網羅率）

```

coverage = mapped_req / total_req

```

表示時は以下で算出する。

```

coverage_rate = coverage * 100

```

小数第2位四捨五入。

---

# 4. 算出ルール（必須）

1. `requirements.md` から `REQ-*` を抽出する
2. 有効REQ一覧を作成する
3. 各REQについて `traceability.md` の行を確認する
4. 以下条件を満たす場合 `mapped_req` に加算する

条件:

- 対象REQの行が存在
- 以下5列がすべて非空

```

Source
Design Mapping
Interface Mapping
Evaluation Mapping
Contract Mapping

```

5. 次式で coverage を算出する

```

coverage = mapped_req / total_req

```

---

# 5. 除外ルール

## 5.1 廃止REQ

以下を満たす場合に限り coverage 集計から除外可能。

- `req-id-lifecycle-spec.md` 4章の必須記録
- 置換先REQ-ID
- 移行期限
- `CHANGES.md` 連携

これらが欠ける場合

```

廃止REQは無効

```

として **total_req に含める**

---

## 5.2 waiver中REQ

waiver は **判定猶予でありマッピング免除ではない**

したがって

```

waiver中REQも total_req に含める

```

ただし

```

5列マッピング未完

```

の場合は

```

mapped_req に含めない

```

---

# 6. 判定

## 6.1 受け入れレビュー判定

| status | 条件 |
|---|---|
| pass | coverage = 100% |
| fail | coverage < 100% |

---

## 6.2 Phase3 完了判定

Phase3 完了条件

```

mapped_req == total_req

```

---

# 7. Status 遷移規定

以下の場合

```

Status: blocked

```

へ遷移する。

1. `coverage < 100%`
2. `traceability.md` に REQ 行欠落
3. 5列の空セル存在
4. `total_req = 0`
5. 入力ソース不整合

---

# 8. 出力

## 8.1 集計出力

```

coverage_total
coverage_matched
coverage_rate
missing_req_ids

````

---

## 8.2 YAML 証跡

```yaml
coverage_report:
  snapshot_at: "2026-03-04T00:00:00Z"
  source_requirements: "memx_spec_v3/docs/requirements.md"
  source_traceability: "memx_spec_v3/docs/traceability.md"
  total_req: 24
  mapped_req: 24
  coverage: 1.0
  coverage_rate: 100
  threshold: 1.0
  phase3_gate: pass
  status: done
  excluded:
    deprecated: []
    waiver: []
  blocked_reasons: []
````

---

## 8.3 Markdown テーブル証跡

| snapshot_at          | total_req | mapped_req | coverage | threshold | phase3_gate | status | blocked_reasons |
| -------------------- | --------: | ---------: | -------: | --------: | ----------- | ------ | --------------- |
| 2026-03-04T00:00:00Z |        24 |         24 |     100% |      100% | pass        | done   | -               |

---

# 9. `design-review-spec.md` 参照規約

`memx_spec_v3/docs/design-review-spec.md` の判定根拠セクションでは
本仕様による coverage 証跡への参照を必須とする。

最低限記録する項目

```
total_req
mapped_req
coverage
coverage_rate
status
blocked_reasons
```