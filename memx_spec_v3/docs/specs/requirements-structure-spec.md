# Requirements Structure Spec

## 1. 目的
本仕様は `memx_spec_v3/docs/requirements.md` の章構造・REQ-ID運用・関連文書同期の統一ルールを定義し、判定対象と説明情報の混在を防止する。

## 2. 章分類ルール（Normative / Informative）

### 2.1 定義
- **Normative（判定対象）**: 実装適合性を判定する要求本文を記載する章。
- **Informative（背景・実装注記・運用補足）**: 背景説明、設計意図、運用メモ、例示を記載する章。

### 2.2 混在禁止ルール
1. 1つの節は `Normative` または `Informative` のどちらか一方に固定し、同一節内で混在させてはならない。
2. `Normative` 節には判定対象の要求のみを記載し、背景説明・実装注記・運用補足は記載してはならない。
3. `Informative` 節には判定対象の要求文（MUST/SHALL相当）を記載してはならない。
4. 背景説明が必要な場合は、`Normative` 節から `Informative` 節へ参照リンクのみを置き、本文の混在を禁止する。

## 3. REQ-ID 配置規約

### 3.1 REQ-ID を許可する章
- `Normative` と明示された節のみ REQ-ID を付与できる。
- 許可形式は `REQ-<カテゴリ>-<連番>`（例: `REQ-CORE-001`）とする。

### 3.2 REQ-ID を禁止する章
- `Informative` 節では REQ-ID を付与してはならない。
- 以下の説明専用章（例）では REQ-ID を禁止する。
  - 用語集
  - 背景・前提
  - 実装ノート
  - 運用補足
  - 変更履歴

### 3.3 判定対象節マーカー標準
`requirements.md` の各節冒頭に以下の2行を必須記載とする。

```md
分類: Normative | Informative
判定対象: yes | no
```

- `分類: Normative` の場合、`判定対象: yes` を必須とする。
- `分類: Informative` の場合、`判定対象: no` を必須とする。
- 組み合わせ不一致（例: `Normative` + `no`）は構造違反として扱う。

## 4. 変更時チェック手順（同期チェック）
`requirements.md` の変更時は、以下を同一変更セット内で確認する。

1. **`traceability.md` 同期**
   - 追加/更新/削除した REQ-ID が `traceability.md` に反映済みであること。
   - REQ-ID のリンク先（設計・実装・検証）が欠落していないこと。
2. **`EVALUATION.md` 同期**
   - 判定対象 REQ-ID の評価観点・評価結果参照が更新されていること。
   - 新規 REQ-ID に未評価状態が残っていないこと。
3. **`interfaces.md` 同期**
   - インターフェース契約へ影響する REQ-ID の差分が `interfaces.md` に反映されていること。
   - I/O制約・エラー契約・互換性注記の齟齬がないこと。

## 5. 受け入れ条件（DR-001 再発防止）

### 5.1 レビュー観点
- すべての節に「分類」「判定対象」マーカーがある。
- `Normative` 節にのみ REQ-ID が存在する。
- `Informative` 節に REQ-ID が存在しない。
- `requirements.md` 差分に対し、`traceability.md` / `EVALUATION.md` / `interfaces.md` の同期確認結果がレビュー記録に残っている。

### 5.2 fail 条件
以下のいずれか1件でも該当した場合、DR-001 再発として受け入れを `fail` とする。

1. 節マーカー未記載または `分類` と `判定対象` の不一致。
2. `Informative` 節への REQ-ID 記載。
3. `Normative` 節で REQ-ID 未付与の判定要求が存在。
4. `requirements.md` の REQ-ID 差分が `traceability.md` / `EVALUATION.md` / `interfaces.md` のいずれかに未反映。
5. 上記 1〜4 の確認結果がレビュー証跡に残っていない。
