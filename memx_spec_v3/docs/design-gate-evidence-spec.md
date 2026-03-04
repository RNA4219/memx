# Design Gate Evidence Spec

## 1. 目的
本仕様は、Phase 1〜4 の設計ゲート判定で必須とする証跡項目、保存規約、最小記録粒度を固定する。

## 2. 適用範囲
- 対象 Phase: 1 / 2 / 3 / 4
- 対象証跡: `lint` / `type` / `test` / `link` / `contract` / `birdseye` / `coverage`
- 保存先の正本は `memx_spec_v3/docs/reviews/` 配下とし、他ディレクトリへの保管は参照リンクのみ許可する。

## 3. Phase別の必須証跡
| Phase | 必須証跡 |
|---|---|
| Phase 1 | `lint`, `type`, `test`, `coverage` |
| Phase 2 | `lint`, `type`, `test`, `coverage`, `link`, `birdseye` |
| Phase 3 | `lint`, `type`, `test`, `coverage`, `contract`, `link`, `birdseye` |
| Phase 4 | `lint`, `type`, `test`, `coverage`, `contract`, `link`, `birdseye` |

## 4. 証跡の保存規約（固定）
### 4.1 保存先
- すべての証跡は `memx_spec_v3/docs/reviews/` に保存する。
- 章単位で管理する場合は `memx_spec_v3/docs/reviews/<chapter_id>/` の利用を許可する。

### 4.2 命名規則
- 命名形式は `EVIDENCE-<kind>-YYYYMMDD-HHMMSS-<run_id>.md` に固定する。
  - `<kind>`: `lint` / `type` / `test` / `link` / `contract` / `birdseye` / `coverage`
  - `<run_id>`: 同一秒内重複を避ける 3 桁以上の英数字識別子
- 集約ファイルを作成する場合のみ `EVIDENCE-INDEX-YYYYMMDD.md` を許可する。

### 4.3 最小記録粒度（必須）
各証跡ファイルには次の 4 項目を必須記録する。
1. 実行日時（`executed_at`）
2. 実行コマンド（`command`）
3. 実行結果（`result`）
4. 判定（`decision`: `pass` / `fail` / `waiver`）

## 5. 証跡種別ごとの追加要件
- `lint`: 対象言語/対象パスを記録する。
- `type`: 型検査対象モジュールを記録する。
- `test`: 実行スイート（例: unit/integration）を記録する。
- `coverage`: 収集方法と測定値（%）を記録する。
- `link`: 検証対象リンク集合（範囲）を記録する。
- `contract`: 契約差分件数（少なくとも high 件数）を記録する。
- `birdseye`: 参照 `node_id` と issue 件数を記録する。

## 6. 他仕様からの参照ルール
- `memx_spec_v3/docs/design-review-spec.md` と `memx_spec_v3/docs/design-acceptance-report-spec.md` は、証跡の保存先・命名・最小記録粒度を本仕様への参照で定義し、重複定義しない。
