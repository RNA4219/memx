---
intent_id: INT-001
owner: your-handle
status: active   # draft|active|deprecated
last_reviewed_at: 2025-10-14
next_review_due: 2025-11-14
---

- HUB優先順: ルート [`../HUB.codex.md`](../HUB.codex.md) を基準とし、`workflow-cookbook/` 配下では [`HUB.codex.md`](HUB.codex.md)（`HUB_SCOPE_DECLARATION`）を優先。
# Workflow Cookbook / Codex Task Kit

[![License: Apache-2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

This repo defines QA/Governance-first workflows (not application code).
AI agents implement changes under these policies with acceptance tests and
canary rules.

<!-- LLM-BOOTSTRAP v1 -->
読む順番:

1. docs/birdseye/index.json  …… ノード一覧・隣接関係（軽量）
2. docs/birdseye/caps/`<path>`.json …… 必要ノードだけ point read（個別カプセル）

フォーカス手順:

- 直近変更ファイル±2hopのノードIDを index.json から取得
- 対応する caps/*.json のみ読み込み

<!-- /LLM-BOOTSTRAP -->

任意のリポジトリに貼って入力時に本リポジトリを参照させるだけで、**仕様→実装→検収**まで一貫して回せるコンテキストエンジニアリングが可能なMD群。

- 人間にもエージェント（Codex等）にも読ませやすい最小フォーマット
- 言語・技術スタック非依存（存在するコマンドだけ使用）

- 品質ゲート参照先: `workflow-cookbook/` の変更は [`docs/QUALITY_GATES.md`](docs/QUALITY_GATES.md) を参照（memx 本体スコープとは別運用）。

## 使い方（最短）

参加前に [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) を読み、連絡窓口と遵守事項を確認してください。

1. リポジトリ直下に本ファイルをコピー
   - フォークで運用する場合は [`docs/FORK_NOTES.md`](docs/FORK_NOTES.md) をテンプレ通りに初期化し、差分管理を開始
2. `BLUEPRINT.md` で要件と制約を1ページに集約
3. 実行手順は `RUNBOOK.md`、評価基準は `EVALUATION.md` に記述し、
   以下で Birdseye の最小読込とタスク分割の前提を共有

    - [`docs/addenda/A_Glossary.md`](docs/addenda/A_Glossary.md) …… Intent や Task Seed など頻出用語の定義を確認

    - [`GUARDRAILS.md`](GUARDRAILS.md) …… 行動指針と Birdseye の `deps_out`
      と整合する最小読込ガードレールを確認
    - [`docs/BIRDSEYE.md`](docs/BIRDSEYE.md) …… Birdseye の Edges / Hot / 更新手順をフォールバック用に把握
    - [`tools/codemap/README.md`](tools/codemap/README.md) …… Birdseye カプセル
      再生成前提と `codemap.update` の流れを把握
    - [`tools/codemap/update.py`](tools/codemap/update.py) ……
      `python tools/codemap/update.py` で `codemap.update` を実行し
      Birdseye カプセルを再生成する。標準では直近変更ファイルから±2 hop の
      カプセルのみ更新し、全カプセルを再生成したい場合は
      `--targets docs/birdseye/index.json,docs/birdseye/hot.json`
      を指定する（`GUARDRAILS.md` の
      [鮮度管理](GUARDRAILS.md#鮮度管理staleness-handling)
      参照）。

      ```sh
      # 例: main との差分から対象カプセルを自動抽出
      python tools/codemap/update.py --since --emit index+caps

      # 例: Birdseye リソースを明示的に指定（従来挙動）
      python tools/codemap/update.py --targets docs/birdseye/index.json,docs/birdseye/hot.json --emit index+caps
      ```

      `docs/birdseye/hot.json` が欠落している場合は上記コマンドの実行時に
      `FileNotFoundError` が発生し、再生成コマンドがメッセージで提示されます。

    - [`HUB.codex.md`](HUB.codex.md) …… 仕様集約とタスク分割ハブを整備し、Birdseye カプセルの依存関係を維持
    - [`docs/IN-20250115-001.md`](docs/IN-20250115-001.md) …… インシデントログを参照し
      Birdseye カプセル要約で指示される `deps_out` を照合
4. タスクごとに `TASK.codex.md` を複製して内容を埋め、エージェントに渡す
   - 雛形との差分を確認したい場合は `examples/TASK.sample.md` を参照し、実在の値が埋め込まれたダミーサンプルと比較する
   - 完了した `TASK.*` の成果は `[Unreleased](CHANGELOG.md#unreleased)` へ通番付きで転記し、該当 Task Seed から成果差分へのリンクを貼る
5. リリースは `CHECKLISTS.md` をなぞり、差分は `CHANGELOG.md` に追記しつつ、`[Unreleased](CHANGELOG.md#unreleased)` に集約した Task Seed 成果をリリースノートへ昇格させる

利用時は [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) を確認し、行動規範の遵守を徹底してください。

## 変更履歴の更新ルール {#changelog-update-rules}

- **更新タイミング**:
  リリース判定が `CHECKLISTS.md` の [Release](CHECKLISTS.md#release) を通過し、レビュー承認で確定した直後に [`CHANGELOG.md`](CHANGELOG.md)
  を更新する。承認前に書き始めないことで、記録の正確性と監査性を確保する。
- **記載形式**:
  [`CHANGELOG.md`](CHANGELOG.md) ではセマンティックバージョニングに従い、`## x.y.z - YYYY-MM-DD` の見出し配下へ `### Added`・
  `### Changed` などのカテゴリ小見出しを用いて差分を整理する。最新リリースを先頭に追記し、既存節の体裁を崩さない。
  各箇条書きの先頭へ 4 桁ゼロ埋めの通番（例: `0001`）を付与し、既存項目の最大値に 1 を加えて採番する。
- **突合手順**:
  1. `CHECKLISTS.md` の [Release](CHECKLISTS.md#release) を順に確認し、完了済みチェック項目と未了項目を照合する。
  2. チェックリストに記録した内容を [`CHANGELOG.md`](CHANGELOG.md) の該当リリース節へ反映し、必要に応じて `RUNBOOK.md` や関連資料の更新有無をメモする。
  3. 反映後に再度チェックリストへ戻り、記録済みであることをコメントまたは添付リンクで明示してリリース完了とする。

詳細な手順は [`docs/addenda/M_Versioning_Release.md`](docs/addenda/M_Versioning_Release.md) を参照する。

### 最小導入セット

- [`BLUEPRINT.md`](BLUEPRINT.md) …… Intent と仕様全体の骨子を提示
- [`RUNBOOK.md`](RUNBOOK.md) …… 実装および運用手順を逐次記載
- [`EVALUATION.md`](EVALUATION.md) …… 受入基準と検証観点を定義（代表シナリオは [付録I](docs/addenda/I_Test_Cases.md) を参照）
- [`GUARDRAILS.md`](GUARDRAILS.md) …… 行動指針と Birdseye 連携の制約を明示
- [`HUB.codex.md`](HUB.codex.md) …… タスク分割と依存グラフの中核ハブを維持
- [`CHECKLISTS.md`](CHECKLISTS.md) …… リリースとレビューフローのチェックリストを提供
- [`governance-gate.yml`](.github/workflows/governance-gate.yml)
  …… Intent 検証 CI を常時有効化
- [`docs/security/SAC.md`](docs/security/SAC.md) …… セキュリティ非機能要件を契約として明文化
- [`docs/security/Security_Review_Checklist.md`](docs/security/Security_Review_Checklist.md) …… 準備→実装→レビューで実施するセキュリティ審査手順を提供
- [`reusable/python-ci.yml`](.github/workflows/reusable/python-ci.yml) /
  [`reusable/security-ci.yml`](.github/workflows/reusable/security-ci.yml)
  …… 他リポから `workflow_call` で利用できる最小CIセット
- Intent ゲートは
  [`tools/ci/check_governance_gate.py`](tools/ci/check_governance_gate.py) により
  自動適用されるため、CI の設定だけで運用に組み込めます
- SRC の主要言語に応じて、以下の CI テストセットを組み合わせると、導入直後から
  最低限の品質・安全・可搬性が確保できます

### 再利用CIの呼び出し例（下流リポ側）

```yaml
name: example CI
on: [push, pull_request]
jobs:
  python:
    uses: RNA4219/workflow-cookbook/.github/workflows/reusable/python-ci.yml@v0.1
    with:
      python-version: '3.11'
  security:
    uses: RNA4219/workflow-cookbook/.github/workflows/reusable/security-ci.yml@v0.1
    with:
      python-version: '3.11'
    secrets: inherit
  governance:
    uses: RNA4219/workflow-cookbook/.github/workflows/governance-gate.yml@v0.1
```

#### 言語別 CI テストセット（鉄板構成）

<!-- markdownlint-disable MD013 -->

| 言語 | カテゴリ | コマンド | 目的 |
| :--- | :--- | :--- | :--- |
| Rust | Format | `cargo fmt --all -- --check` | コード整形確認 |
| Rust | Lint | `cargo clippy --all-targets --all-features -D warnings` | 警告・アンチパターン検出 |
| Rust | Test | `cargo test --all-features` | 単体・統合テスト |
| Rust | Build | `cargo build --release` | リリースビルド確認 |
| Rust | Security | `cargo audit` / `cargo deny check` | 依存脆弱性・ライセンス確認 |
| Rust | Docs | `cargo doc -D warnings` | ドキュメント構文確認 |
| Rust (任意) | Coverage | `cargo llvm-cov` | カバレッジ測定 |
| Python | Format | `black --check .` | コード整形確認 |
| Python | Lint | `ruff check .` / `flake8 .` / `pylint src/` | コード品質 |
| Python | Typing | `mypy src/` | 型整合性 |
| Python | Test | `pytest --maxfail=1 --disable-warnings -q` | 単体・統合テスト |
| Python | Security | `bandit -r src/` / `pip-audit` | 静的セキュリティ解析 |
| Python | Coverage | `pytest --cov=src` | カバレッジ |
| Python | Docs | `pydocstyle src/` | docstring 構文確認 |
| Node.js/TS | Lint | `eslint .` / `tsc --noEmit` | コード・型検査 |
| Node.js/TS | Format | `prettier --check .` | 整形検証 |
| Node.js/TS | Test | `npm test` / `vitest run` / `jest --ci` | ユニットテスト |
| Node.js/TS | Build | `npm run build` | 本番ビルド |
| Node.js/TS | Security | `npm audit --audit-level=moderate` | 依存脆弱性 |

> CI 設定全体と最新実行のみを保持する自動キャンセル構成は、[`docs/ci-config.md`](docs/ci-config.md) を参照してください。
| Node.js/TS (任意) | Coverage | `npm run coverage` | カバレッジ |
| Node.js/TS | Docs | `typedoc` / `markdownlint` | API/MD構文確認 |
| Go | Format | `gofmt -l .` | 整形確認 |
| Go | Lint | `golangci-lint run` | 静的解析 |
| Go | Vet | `go vet ./...` | 構文・型安全検査 |
| Go | Test | `go test -v ./...` | 単体テスト |
| Go | Build | `go build ./...` | ビルド保証 |
| Go | Security | `gosec ./...` | 脆弱性スキャン |
| Go (任意) | Coverage | `go test -cover ./...` | カバレッジ測定 |
| Java | Format | `mvn formatter:validate` / `spotless:check` | 整形確認 |
| Java | Lint | `mvn checkstyle:check` / `spotbugs:check` | 静的解析 |
| Java | Test | `mvn test` / `gradle test` | 単体テスト |
| Java | Coverage | `mvn jacoco:report` | カバレッジ |
| Java | Security | `mvn dependency-check:check` | 依存脆弱性 |
| Java | Build | `mvn package` / `gradle build` | 本番ビルド |
| C/C++ | Format | `clang-format --dry-run -Werror` | 整形確認 |
| C/C++ | Lint | `cppcheck --enable=all` | 静的解析 |
| C/C++ | Build | `cmake . && make` / `meson compile` | ビルド確認 |
| C/C++ | Test | `ctest --output-on-failure` / `gtest` | 単体テスト |
| C/C++ | Security | `clang --analyze` / `cppcheck --addon=cert` | セキュリティ |
| C/C++ | Coverage | `lcov` / `gcov` | カバレッジ測定 |

#### 共通モジュール4種（全リポ共通）

| 領域/モジュール | 目的・役割 | 主要仕様書 | 備考 |
| :--- | :--- | :--- | :--- |
| CodeQL | 静的解析・脆弱性検出を CI に組み込む | [docs/spec.md](docs/spec.md) / [docs/ci-config.md](docs/ci-config.md) | `github/codeql-action` ワークフローで SAST ゲートを維持 |
| Dependabot | 依存更新の自動 PR を定期化する | [docs/requirements.md](docs/requirements.md) / [docs/spec.md](docs/spec.md) | 週次スケジュールで依存差分を検知し CI と連動 |
| Pre-commit Hooks | Lint / Format をローカルで再現する | [docs/design.md](docs/design.md) / [docs/requirements.md](docs/requirements.md) | `.pre-commit-config.yaml` でチーム基準を固定 |
| Artifact Upload | テスト結果やログを共有する | [docs/ci-config.md](docs/ci-config.md) / [EVALUATION.md](EVALUATION.md#acceptance-criteria) | CI 実行痕跡をアーカイブしてレビューへ提示 |
| `examples/` | レシピ参照実装と設計・仕様の整合を確認する | [docs/design.md](docs/design.md) / [docs/spec.md](docs/spec.md) | サンプル更新時は Birdseye 同期[^birdseye] |
| `styles/` | QA ルールによる表記統一・禁止用語を管理する | [docs/design.md](docs/design.md) / [docs/requirements.md](docs/requirements.md) | `styles/qa/QA.yml` の用語ルールを適用[^styles] |
| `tools/` | ドキュメント同期と検証スクリプトを運用する | [docs/design.md](docs/design.md) / [RUNBOOK.md](RUNBOOK.md#execute) | `tools/codemap/update.py` で Birdseye を再生成[^birdseye] |
| `docs/security/` | セキュリティレビューと SAC 手順を集約する | [docs/security/Security_Review_Checklist.md](docs/security/Security_Review_Checklist.md) / [docs/security/SAC.md](docs/security/SAC.md) | リリース審査の証跡を更新 |

[^birdseye]: `python tools/codemap/update.py --since --emit index+caps` で Birdseye インデックスとカプセルを更新し、必要に応じて `--targets docs/birdseye/index.json,docs/birdseye/hot.json` などを併用して生成対象を明示する。`--emit` の指定により出力形式を切り替える。詳細は [tools/codemap/README.md](tools/codemap/README.md#実行手順) と [GUARDRAILS.md](GUARDRAILS.md#鮮度管理staleness-handling) を参照。
[^styles]: `styles/qa/QA.yml` の禁止用語・表記揺れルールをレビューで適用し、検知結果を `CHECKLISTS.md` のリリース手順へ反映する。

![lint](https://github.com/RNA4219/workflow-cookbook/actions/workflows/markdown.yml/badge.svg)
![links](https://github.com/RNA4219/workflow-cookbook/actions/workflows/links.yml/badge.svg)
![lead_time_p95_hours](https://img.shields.io/badge/lead__time__p95__hours-24h-blue)
![mttr_p95_minutes](https://img.shields.io/badge/mttr__p95__minutes-30m-blue)
![change_failure_rate_max](https://img.shields.io/badge/change__failure__rate__max-0.20-blue)
<!-- markdownlint-enable MD013 -->

> バッジ値は `governance/policy.yaml` の `slo` と同期。更新時は同ファイルの値を修正し、上記3つのバッジ表示を揃える。

## License

Apache-2.0. Unless noted otherwise,
files copied from this repo into other projects remain Apache-2.0 and require retaining NOTICE text in redistributions.

### Commit message guide

- feat: 〜 を追加
- fix: 〜 を修正
- chore/docs: 〜 を整備
- semver:major/minor/patch ラベルでリリース自動分類

### Pull Request checklist (CI 必須項目)

- PR 本文に `Intent: INT-xxx`（例: `Intent: INT-123`）を含めること。
- `EVALUATION` 見出し（例:
  `[Acceptance Criteria](EVALUATION.md#acceptance-criteria)`）へのリンクを本文に
  明示すること。
- 可能であれば [`Priority Score`](docs/addenda/A_Glossary.md#priority-score):
  `<number>` を追記し、`governance/prioritization.yaml` の重み付けを参照する。
- ローカルでゲートを確認する場合は `PR_BODY` に PR 本文を渡してから
  `python tools/ci/check_governance_gate.py` を実行する。
  （任意のディレクトリからは `python /path/to/workflow-cookbook/tools/ci/check_governance_gate.py`
  のように実行しても同様に検証できる。）

  ```sh
  PR_BODY=$(cat <<'EOF'
  Intent: INT-123
  ## EVALUATION
  - [Acceptance Criteria](EVALUATION.md#acceptance-criteria)
  Priority Score: 1
  EOF
  ) python tools/ci/check_governance_gate.py
  ```
