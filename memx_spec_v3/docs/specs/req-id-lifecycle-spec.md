# REQ-ID Lifecycle Spec

## 1. 目的
本仕様は、`requirements.md` に定義される `REQ-*` の新規採番・変更・廃止のライフサイクル管理ルールを固定し、`traceability.md` と `CHANGES.md` の整合を維持するための必須要件を定義する。

## 2. 新規採番規則（プレフィックス、連番、スコープ）
- プレフィックスは **`REQ-<カテゴリ>-<3桁連番>`** とする。
  - 例: `REQ-CLI-001`, `REQ-INT-012`
- 連番はカテゴリ単位で 001 から開始し、欠番再利用は禁止する。
- スコープは `requirements.md` を正本（source of truth）とし、同一 REQ-ID を複数要件に再利用してはならない。
- 既存 ID の意味を変更せずに新要件を追加する場合は、新規 ID を採番する。

## 3. 変更種別ごとの扱い
### 3.1 文言修正（意味不変）
- 要件の意味・受け入れ基準が不変な修正（誤字、語順、表現統一）は **REQ-ID を維持**する。
- `traceability.md` は既存マッピングを維持し、必要に応じて参照アンカーのみ更新する。

### 3.2 互換維持変更（意味拡張）
- 既存利用者との互換を維持したまま要件を拡張する変更は **REQ-ID を維持**する。
- 変更内容を `CHANGES.md` に追記し、影響範囲（Design/Interface/Evaluation/Contract）を明記する。
- `traceability.md` の 5 列マッピングは同一 PR で更新する。

### 3.3 破壊変更（非互換）
- 既存利用者に非互換を生む変更は **新規 REQ-ID を採番**し、旧 REQ-ID は廃止扱いにする。
- 旧 REQ-ID には置換先 REQ-ID を明記し、移行期限を設定する。
- `CHANGES.md` には Breaking Change として記録し、段階移行の期限を記載する。

## 4. 廃止時の必須記録
REQ-ID を廃止する場合、以下 3 点をすべて必須とする。

1. **置換先 REQ-ID**
   - 廃止対象ごとに後継 `REQ-*` を 1 つ以上明記する。
2. **移行期限**
   - `YYYY-MM-DD` 形式で期限を定義する。
3. **`CHANGES.md` 連携**
   - 廃止理由、利用者影響、移行手順、期限を `CHANGES.md` に記録する。

上記 3 点のいずれかが欠ける廃止は無効とし、レビュー判定を `fail` とする。

## 5. `traceability.md` 同時更新必須ルール（未マッピング時は fail）
- `requirements.md` で `REQ-*` を追加・変更・廃止した PR は、同一 PR で `traceability.md` を必ず更新する。
- `traceability.md` では対象 REQ-ID の行について、以下 5 列すべてを必須入力とする。
  - `Source`
  - `Design Mapping`
  - `Interface Mapping`
  - `Evaluation Mapping`
  - `Contract Mapping`
- `Design Mapping` は必ず `memx_spec_v3/docs/design.md#...` 形式の参照を含む。
- 未マッピング（行欠落または 5 列の空欄）が 1 件でもある場合、レビュー判定を `fail` とする。

## 6. `design-review-spec.md` 6章との整合
- `memx_spec_v3/docs/design-review-spec.md` 6章「差分レビュー時の未マッピングREQ検出手順」は、本仕様の 5 章を参照して検証を実施する。
- レビュー時は次の 2 点を同時に満たすことを必須とする。
  1. 差分に含まれる対象 REQ-ID が `traceability.md` に存在する。
  2. 対象行の 5 列（`Source / Design Mapping / Interface Mapping / Evaluation Mapping / Contract Mapping`）がすべて充足している。
