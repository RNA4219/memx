# Design HUB Source Coverage Spec

## 1. 目的
- 本仕様は `gate_hub_source_coverage` の判定ロジックと証跡要件の正本を定義する。
- 検索キーと対象パスを固定し、Phase 間での判定ぶれと重複定義を防止する。

## 2. 判定対象（固定）
- 検索キー（判定キー）
  - `Incident`
  - `Orchestration`
  - `TASK`
- 対象パス
  - `docs/IN-*.md`
  - `orchestration/*.md`
  - `TASK.*`

## 3. 判定キー別ルール（pass/warn/fail）

| 判定キー | pass 条件 | warn 条件 | fail 条件 |
| --- | --- | --- | --- |
| `Incident` | `docs/IN-*.md` に検索ヒットが1件以上あり、必要情報を抽出できる | ヒットはあるが有効証跡として採用できない（テンプレート混入、抽出不能、記載不備） | ヒット0件、または必要検索を未実施 |
| `Orchestration` | `orchestration/*.md` に検索ヒットが1件以上あり、依存/実行情報を抽出できる | ヒットはあるが依存関係・手順の特定が不十分 | ヒット0件、または必要検索を未実施 |
| `TASK` | `TASK.*` に検索ヒットが1件以上あり、Task 連携情報を抽出できる | ヒットはあるが Task 連携情報が不十分（要追記） | ヒット0件、または必要検索を未実施 |

## 4. `high/medium/low` への写像規則
- `gate_hub_source_coverage = high`
  - `Incident` または `Orchestration` のいずれかが `fail`。
- `gate_hub_source_coverage = medium`
  - `high` 条件に該当せず、いずれかの判定キーが `warn` または `fail`。
  - 典型例: `TASK` のみ `fail`、または `TASK`/`Incident`/`Orchestration` のいずれかが `warn`。
- `gate_hub_source_coverage = low`
  - 全判定キー（`Incident` / `Orchestration` / `TASK`）が `pass`。

## 5. 証跡の必須項目
判定ログには、各判定キーごとに次を必須記録する。
1. 実行した検索コマンド（対象パスが分かること）。
2. ヒット件数（0件を含む）。
3. 欠落理由（`warn`/`fail` の場合は必須。未実施・抽出不能・記載不備などを具体化）。

## 6. 運用ルール
- 判定時は 3 キーすべてを毎回評価する。
- 任意キーワード追加や任意パス追加での判定上書きは禁止する。
- `memx_spec_v3/docs/design-phase-gate-spec.md` の `gate_hub_source_coverage` は本仕様の写像規則を唯一の正本として適用する。
