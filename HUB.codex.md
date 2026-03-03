# HUB.codex.md

## 目的
エージェント応答の契約を統一し、実行環境差異（tools 可用性差）に依存しない最小互換出力を保証する。

## 必須出力（固定順）
以下 5 セクションを必須とする（通常時）。

1. `plan`
   - 実施方針を箇条書きで記載。
   - 変更対象・非対象を明記。
2. `patch`
   - 変更ファイルと要点差分を記載。
   - 実変更がない場合は `no-op` と明記。
3. `tests`
   - 実行した検証コマンドと結果（pass/fail/warn）を記載。
   - 未実行の場合は理由を 1 行で記載。
4. `commands`
   - 実行・提案コマンドを列挙。
   - **正本（canonical source）** は `memx_spec_v3/docs/quickstart.md` の API 起動/投入/検索/表示の各コマンド例とする。
   - 参照リンク: [`memx_spec_v3/docs/quickstart.md`](memx_spec_v3/docs/quickstart.md)
5. `notes`
   - 判断理由、制約、未解決事項を最小限で記載。
   - 競合解消がある場合は「双方の意図をどう最小統合したか」を 1 行で記載。

## 失敗時出力契約
- ツール未使用・利用不能・実行基盤制約で処理不能な場合は、
  - `plan` と要求情報（tool 呼び出しの `request envelope` JSON）のみを返す。
  - `patch/tests/commands/notes` は省略可。
- 通常時（ツール実行成功時）は `plan/patch/tests/commands/notes` の 5 セクションを必須とする。
- 禁止事項: ツール結果の推測・捏造。

## 実行環境差異の統一方針
- Native function-calling 可能環境: ツール呼び出しを実行し、同内容の request envelope を併記可。
- 非対応環境: request envelope を一次成果物として返却。
- どの環境でも、最終的な契約解釈は本ファイルを優先する。

- オーケストレーション入力ソースとして `orchestration/*.md` を参照する。

## 言語ポリシー
- デフォルト言語は日本語。
- コード識別子（変数名・関数名・型名・CLI フラグ・JSON キー）は英語を維持する。
- 外部仕様や既存 API 名は原文尊重で改変しない。
