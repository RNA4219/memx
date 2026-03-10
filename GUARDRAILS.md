---
intent_id: INT-001
owner: memx-resolver
status: active
last_reviewed_at: 2026-03-10
next_review_due: 2026-04-10
---

# Guardrails & 行動指針

## 目的

- docs/requirements.md, docs/interfaces.md, docs/design.md の仕様を遵守する
- 変更は最小差分で行い、Public APIを破壊しない
- テスト駆動開発を基本とし、テストを先に記述する

## スコープとドキュメント

1. 目的を一文で定義し、誰のどの課題をなぜ今扱うかを明示する
2. Scopeを固定し、In/Outの境界を先に決めて記録する
3. I/O契約をBLUEPRINT.mdに整理する
4. Acceptance CriteriaをEVALUATION.mdに列挙する
5. 最小フローをRUNBOOK.mdに記す
6. HUB.codex.mdの自動タスク分割フローに従う

## 実装原則

- 型安全：新規・変更シグネチャには必ず型を付与
- 例外設計：既存errors階層に合わせ、再試行可否を区別
- 後方互換：CLI/JSON出力は互換性を維持
- スコープ上限：1回の変更は合計100行または2ファイルまで

## Birdseye / Minimal Context Intake Guardrails

### MUST（必須）

1. まずREADME.mdのLLM-BOOTSTRAPブロックのみ読む
2. `docs/birdseye/index.json`を読み、対象変更ファイル±2 hopのノードID集合を得る
3. 対応する`docs/birdseye/caps/*.json`だけを読み込む
4. 生成物ではノードID（パス）を明示し出典を示す

### SHOULD（推奨）

- 2 hopの合計が1,200 tokensを超えそうなら1 hopに縮小
- 読み順はentrypoints → application → domain → infra → ui

### MUST NOT（禁止）

- `node_modules`, `.venv`, `dist`, `build`, `coverage`等の重量ディレクトリを直読みしない
- 偽の読込結果を作らない

## 鮮度管理（Staleness Handling）

- `index.json.generated_at`が最新コミットより古い場合、再生成を要求
- 依頼フロー：人間がcodemapスクリプトをローカルで実行し、成果物をコミット

## 生成物に関する要求（出力契約）

- **plan**：読み込んだCapsノードID一覧とhop、抜粋理由
- **patch**：変更対象ファイルの相対パスを先頭コメントで明記
- **tests**：対象ノードのtests/*を参照して増補
- **commands**：読込に使ったツールと再現手順
- **notes**：鮮度判断、スコープ外ファイル、既知リスク