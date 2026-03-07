# memx エージェントガイド

**AIエージェント向けの利用案内** - このドキュメントは、AIエージェントが memx を理解し活用するための最初のガイドです。

---

## memx とは？

**memx** は、ローカルLLM/エージェント向けの **メモリ基盤** です。
エージェントが情報を保存・検索・参照するための軽量なストアを提供します。

### なぜ必要か？

LLMエージェントには「記憶」がありません：
- 長い会話の文脈を忘れる
- 過去の知識を再利用できない
- ユーザー固有の情報を保持できない

**memx はこれらを解決する「外部メモリ」を提供します。**

---

## 基本概念：4つのストア

memx は用途別に4つのストアを提供します：

| ストア | 用途 | 保持期間 |
|--------|------|----------|
| **short** | 短期記憶（作業中のメモ、一時的な情報） | 日〜週単位 |
| **journal** | 長期記憶（重要な出来事、プロジェクト履歴） | 月〜年単位 |
| **knowledge** | 知識ベース（FAQ、手順書、概念説明） | 永続 |
| **archive** | アーカイブ（不要だが保持する情報） | 無期限 |

**v1.3 では全ストアが実装されています。**

---

## Claude Code スキル

memx は Claude Code 用のスキルを提供しています。これらを使用することで、自然言語で memx を操作できます。

| スキル | 説明 | 例 |
|--------|------|-----|
| `/remember` | 情報を short ストアに保存 | `/remember API認証は JWT を使用` |
| `/recall` | 情報を検索 | `/recall 認証` |
| `/journal` | journal ストアにログ追加 | `/journal --scope project:memx 実装完了` |
| `/knowledge` | knowledge ストアに知識追加 | `/knowledge --scope glossary JWT = JSON Web Token` |
| `/show` | ノート詳細表示 | `/show abc123...` |
| `/memx-help` | ヘルプ表示 | `/memx-help` |

### スキル定義ファイル

スキルは `.claude/commands/` ディレクトリに定義されています：

```
.claude/commands/
├── remember.md    # /remember スキル
├── recall.md      # /recall スキル
├── journal.md     # /journal スキル
├── knowledge.md   # /knowledge スキル
├── show.md        # /show スキル
└── memx-help.md   # /memx-help スキル
```

---

## クイックスタート

### インストール

```bash
cd memx_spec_v3/go
go build ./cmd/mem
```

### 基本的な使い方

#### 1. メモを保存する（ingest）

**short ストア**（短期記憶）：

```bash
# タイトルと本文を指定
mem in short --title "会議メモ" --body "明日の10時に打ち合わせ"

# 標準入力から読み込み
echo "重要な情報" | mem in short --title "メモ" --stdin

# タグを付ける
mem in short --title "バグ報告" --body "詳細..." --tag bug --tag priority-high
```

**journal ストア**（時系列ログ）：

```bash
# working_scope が必須
mem in journal --title "進捗ログ" --body "API実装完了" --scope project:memx
```

**knowledge ストア**（知識ベース）：

```bash
# working_scope が必須
mem in knowledge --title "用語定義" --body "JWT = JSON Web Token" --scope glossary

# ピン留め付き
mem in knowledge --title "重要な設定" --body "環境変数..." --scope config --pinned
```

#### 2. メモを検索する（search）

**short ストア**：

```bash
# キーワード検索
mem out search "会議"

# 結果数を指定
mem out search "バグ" -k 5
```

**journal ストア**：

```bash
mem out journal search "実装"
mem out journal list --scope project:memx
```

**knowledge ストア**：

```bash
mem out knowledge search "JWT"
mem out knowledge list --scope glossary
mem out knowledge pinned  # ピン留め一覧
```

#### 3. メモを詳細表示する（show）

```bash
# IDを指定して詳細表示
mem out show <NOTE_ID>
```

#### 4. 要約を生成する（summarize）

```bash
# 単一メモの要約
mem summarize <NOTE_ID>

# 複数メモの一括要約
mem summarize --ids id1,id2,id3
```

#### 5. 古いメモを整理する（GC）

```bash
# dry-run: 整理対象を確認（DBは更新しない）
mem gc short --dry-run

# 実行: 古いメモをarchiveへ退避
mem gc short --enable-gc
```

**GC機能**:
- デフォルトで無効（feature flag）
- `--dry-run` で判定結果のみ表示
- `--enable-gc` で実際に実行
- トリガ条件: 1200ノート超過（soft limit）または2000ノート超過（hard limit）

---

## API サーバーとして使う

### サーバー起動

```bash
mem api serve --addr 127.0.0.1:7766
```

### API エンドポイント

**short ストア**：

| エンドポイント | 用途 |
|----------------|------|
| `POST /v1/notes:ingest` | メモ保存 |
| `POST /v1/notes:search` | メモ検索 |
| `GET /v1/notes/{id}` | メモ取得 |
| `POST /v1/notes:summarize` | 要約生成 |
| `POST /v1/gc:run` | GC実行 |

**journal ストア**：

| エンドポイント | 用途 |
|----------------|------|
| `POST /v1/journal:ingest` | ログ保存 |
| `POST /v1/journal:search` | ログ検索 |
| `GET /v1/journal/{id}` | ログ取得 |
| `POST /v1/journal:list-by-scope` | スコープ別一覧 |

**knowledge ストア**：

| エンドポイント | 用途 |
|----------------|------|
| `POST /v1/knowledge:ingest` | 知識保存 |
| `POST /v1/knowledge:search` | 知識検索 |
| `GET /v1/knowledge/{id}` | 知識取得 |
| `POST /v1/knowledge/{id}:pin` | ピン留め |
| `POST /v1/knowledge/{id}:unpin` | ピン解除 |
| `POST /v1/knowledge:list-pinned` | ピン留め一覧 |

**archive ストア**：

| エンドポイント | 用途 |
|----------------|------|
| `GET /v1/archive` | アーカイブ一覧 |
| `GET /v1/archive/{id}` | アーカイブ取得 |
| `POST /v1/archive/{id}:restore` | 復元 |

### CLI から API を使う

```bash
# APIサーバー経由で実行
mem in short --title "test" --body "body" --api-url http://127.0.0.1:7766
mem out search "test" --api-url http://127.0.0.1:7766
```

---

## JSON 出力

CLI は `--json` フラグで JSON 出力をサポートします。API レスポンスと同型です。

```bash
mem out search "test" --json
```

---

## セキュリティ機能

### sensitivity（機密度）

メモには機密度を設定できます：

| 値 | 説明 |
|----|------|
| `public` | 公開可能 |
| `internal` | 内部利用（既定値） |
| `confidential` | 機密 |
| `secret` | 極秘（**保存拒否**） |

**fail-closed 方針**: `secret` 指定のメモは保存を拒否します。

```bash
# これはエラーになる（secret は保存拒否）
mem in short --title "secret" --body "..." --sensitivity secret
```

---

## アーキテクチャ概要

```
┌─────────────────────────────────────────────┐
│                  CLI / API                   │
├─────────────────────────────────────────────┤
│                Service Layer                │
│  ┌─────────────┐  ┌──────────────────────┐  │
│  │ Gatekeeper  │  │   Validation Layer   │  │
│  └─────────────┘  └──────────────────────┘  │
├─────────────────────────────────────────────┤
│                   DB Layer                  │
│  ┌───────┐ ┌───────────┐ ┌──────────┐ ┌────┐│
│  │short  │ │journal  │ │knowledge │ │arch││
│  └───────┘ └───────────┘ └──────────┘ └────┘│
└─────────────────────────────────────────────┘
```

---

## ドキュメント構成

エージェントが参照すべき正本ドキュメント：

| 目的 | ドキュメント |
|------|--------------|
| 要件を確認する | `memx_spec_v3/docs/requirements.md` |
| 設計を理解する | `memx_spec_v3/docs/design.md` |
| API仕様を見る | `memx_spec_v3/docs/contracts/openapi.yaml` |
| CLI仕様を見る | `memx_spec_v3/docs/contracts/cli-json.schema.json` |

**重要**: 正本（Normative）と補助（Secondary）の区別があります。
正本が優先されます。詳細は `memx_spec_v3/docs/spec.md` を参照してください。

---

## よくあるタスク

### 新しいメモを保存する

```bash
mem in short --title "タイトル" --body "本文"
```

### 過去のメモを探す

```bash
mem out search "キーワード"
```

### 特定のメモの詳細を見る

```bash
mem out show <ID>
```

### メモを要約する

```bash
mem summarize <ID>
```

---

## 次のステップ

- 詳細な要件: [`requirements.md`](./memx_spec_v3/docs/requirements.md)
- 設計の詳細: [`design.md`](./memx_spec_v3/docs/design.md)
- API 契約: [`contracts/openapi.yaml`](./memx_spec_v3/docs/contracts/openapi.yaml)