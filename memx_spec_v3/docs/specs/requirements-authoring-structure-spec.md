# Requirements Authoring Structure Spec

## 1. 目的
本仕様は `memx_spec_v3/docs/requirements.md` の記述運用を、判定対象（Normative）と補足（Informative）で章レベル分離するための作成規約として定義する。

## 2. 適用範囲
- 対象文書: `memx_spec_v3/docs/requirements.md` / `memx_spec_v3/docs/spec.md` / `EVALUATION.md`
- 対象者: 要件作成者・レビュー担当者・評価担当者

## 3. Normative / Informative の章レベル分離

### 3.1 章分類の原則
1. `requirements.md` の各章（見出し単位）は `Normative` または `Informative` のいずれか一方に固定する。
2. 同一章内で判定対象要求と背景説明を混在させてはならない。
3. 判定対象要求（MUST/SHALL 相当）は `Normative` 章にのみ記載する。
4. 背景・理由・実装注記・運用補足・例示は `Informative` 章にのみ記載する。

### 3.2 既存規定との統合
- 章分類の定義・混在禁止・判定対象マーカーは `requirements-structure-spec.md` を正本として参照する。
- 本書は「作成時の運用ルール（配置順・参照先・リンク固定）」を追加定義する。

## 4. REQ-ID 配置ルール（判定対象の先頭節集約）

### 4.1 基本規約
1. 新規/既存の REQ-ID は、各テーマの **先頭 Normative 節** に集約配置する。
2. 同テーマ内の補足節（Informative 節）へ新規 REQ-ID を追加してはならない。
3. Informative 節で要件補足が必要な場合は、先頭 Normative 節の既存 REQ-ID を参照リンクで再利用する。

### 4.2 禁止事項
- Informative 節での REQ-ID 新設。
- 補足節だけで成立する判定要求の新設。
- 同一要求の REQ-ID 重複採番。

## 5. `spec.md` / `requirements.md` 冒頭の最短参照順ブロック

### 5.1 必須ブロック
`spec.md` および `requirements.md` の先頭（タイトル直下）に、以下 2 区分の「最短参照順」ブロックを必ず配置する。
- CLI利用者向け
- API利用者向け

### 5.2 記載ルール
1. 各区分は 3〜5 行程度の短い参照順で記載する。
2. 先頭リンクは当該利用者が最初に読むべき判定対象章（Normative）を指す。
3. Informative 章へのリンクは末尾に置き、補足であることを明示する。
4. 既存リンクを変更する場合は、`traceability.md` の参照整合を同一変更セットで確認する。

## 6. `EVALUATION.md` の判定参照先ルール

### 6.1 判定参照先の限定
1. `EVALUATION.md` の pass/fail 判定参照は、`requirements.md` の Normative 章のみを参照先とする。
2. Informative 章・付録・注記章を判定参照先にしてはならない。

### 6.2 リンク規約（固定）
- `EVALUATION.md` から `requirements.md` への参照は、以下の形式に固定する。
  - 形式: `./memx_spec_v3/docs/requirements.md#<Normative章アンカー>`
- 判定テーブル（Requirement ID 行）では、REQ-ID ごとに上記形式の固定リンクを必須とする。
- 章名変更でアンカーが変わる場合、`EVALUATION.md` 側リンクを同一コミットで更新する。

## 7. 重複規定の参照統合と優先順位

### 7.1 参照統合
- 以下は本書で再定義せず、`requirements-structure-spec.md` を参照する。
  - Normative / Informative 定義
  - 節マーカー（`分類` / `判定対象`）
  - 基本 fail 条件

### 7.2 衝突時優先順位
規定が衝突した場合の優先順位は以下とする。
1. 本書（`requirements-authoring-structure-spec.md`）の運用規約
2. `requirements-structure-spec.md` の基本構造規約
3. 個別文書内の注記

## 8. 受け入れチェック項目
- `requirements.md` の章が Normative / Informative で分離され、混在がない。
- REQ-ID が各テーマの先頭 Normative 節へ集約され、Informative 節に新規 REQ-ID がない。
- `spec.md` / `requirements.md` 冒頭に「CLI利用者向け」「API利用者向け」最短参照順ブロックがある。
- `EVALUATION.md` の判定参照リンクが Normative 章のみを指し、固定形式に一致する。
- 重複規定は `requirements-structure-spec.md` 参照に統合され、優先順位が明記されている。
