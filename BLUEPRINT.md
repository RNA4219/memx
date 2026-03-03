---
intent_id: memx-governance-blueprint-v1
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# BLUEPRINT

## 目的とスコープ
- memx は個人用・ローカル運用の「知識＋記憶」管理システムとして、short / chronicle / memopedia / archive の4ストアを分離運用する。
- CLI + API + SQLite を前提に、実行環境依存を減らした知識基盤を提供する。
- v1 非ゴール: Web UI、マルチユーザー公開運用、常駐必須設計、完全自動エージェント。

## レイヤリング方針
- 基本フロー: Human CLI → Tool/AI API → Service(Usecase) → DB/LLM/Gatekeeper。
- CLI は入力整形と表示のみを担当し、DB へ直接アクセスしない。
- API は安定 JSON I/F を提供し、Service を唯一のビジネスロジック入口とする。

## ストア・スキーマ方針
- 物理 DB: `short.db`, `chronicle.db`, `memopedia.db`, `archive.db`。
- 共通テーブル: `notes`, `tags`, `note_tags`, `note_embeddings`, `notes_fts`（archive は一部省略可）。
- short 固有: `short_meta`（GCトリガ用メタ）と `lineage`（蒸留/昇格/退避の系譜）。
- スキーマバージョンは `PRAGMA user_version` を利用し、破壊的/非互換 DDL 時のみインクリメントする。

## 検索・評価方針
- 検索は FTS5 + 埋め込みベクター検索を併用し、将来の SQLite 拡張差し替えを前提に API を固定する。
- 評価軸は `relevance`, `quality`, `novelty`, `importance_static`。
- `memory_policy.yaml` で閾値・禁止パターン・decay を管理する。

## Gatekeeper 方針
- `memory_store` / `memory_output` の2タイミングで判定フックを持つ。
- 判定は `allow` / `deny` / `needs_human`、v1.3 では `needs_human` を deny 相当で fail-closed 扱い。
- エラーマッピングは service 層で sentinel error 化し、API 側で集約マッピングする。

## v1 API Contract

### 対象エンドポイント
- `POST /v1/notes:ingest`
- `POST /v1/notes:search`
- `GET /v1/notes/{id}`

### 最小 JSON 例

#### `POST /v1/notes:ingest`
request:
```json
{
  "store": "short",
  "title": "daily-log",
  "body": "今日は検索性能を確認した"
}
```

response (`200 OK`):
```json
{
  "id": "n_123",
  "store": "short",
  "created_at": "2026-03-03T10:00:00Z"
}
```

#### `POST /v1/notes:search`
request:
```json
{
  "store": "short",
  "query": "検索性能",
  "limit": 3
}
```

response (`200 OK`):
```json
{
  "items": [
    {
      "id": "n_123",
      "title": "daily-log",
      "snippet": "今日は検索性能を確認した"
    }
  ],
  "total": 1
}
```

#### `GET /v1/notes/{id}`
response (`200 OK`):
```json
{
  "id": "n_123",
  "store": "short",
  "title": "daily-log",
  "body": "今日は検索性能を確認した",
  "created_at": "2026-03-03T10:00:00Z"
}
```

### エラー分類と HTTP/コード方針
- 入力不正: `400 Bad Request` を返し、`error.code` は `INVALID_ARGUMENT` 系（例: `INVALID_ARGUMENT`, `INVALID_STORE`, `INVALID_LIMIT`）を使用する。
- ポリシー拒否: `403 Forbidden` を返し、`error.code` は `POLICY_DENIED`（`needs_human` を含む fail-closed）を使用する。
- 内部障害: `500 Internal Server Error` を返し、`error.code` は `INTERNAL_ERROR` を使用する。

### 互換性禁止事項（v1 固定）
- 必須フィールドの削除を行わない。
- 既存フィールドの意味を変更しない。
- 成功レスポンスのトップレベル構造を変更しない。

`--json` 運用との対応は `RUNBOOK.md` の「CLI `--json` と API レスポンスの変換ルール（v1必須3エンドポイント）」を正本とする。

## API/互換性方針
- v1 必須 API: `POST /v1/notes:ingest`, `POST /v1/notes:search`, `GET /v1/notes/{id}`。
- v1 互換性は後方互換維持を原則とし、必須フィールド削除・既存意味変更・成功レスポンス構造破壊を禁止する。
- 破壊変更は並行エンドポイント追加または `/v2` 新設で段階移行する。
