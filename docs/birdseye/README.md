# Birdseye インデックス運用

`docs/birdseye/index.json` は memx_spec_v3 の主要ドキュメント/実装ノードの依存関係を管理する HUB です。

## 構成
- `index.json`: HUB とノード一覧（`node_id`, `role`, `depends_on`, `generated_at` を保持）
- `caps/*.json`: 主要文書ごとの capsule（要点・依存ノードを保持）

## 更新タイミング
以下の変更が入るたびに `index.json` と関連 `caps/*.json` を同時更新してください。

1. **設計変更時**
   - `memx_spec_v3/docs/requirements.md` の要件・レイヤ構成・ストア責務を更新したとき
   - `memx_spec_v3/docs/quickstart.md` の導線や前提手順を更新したとき
2. **API変更時**
   - `memx_spec_v3/go/api/http_server.go` のエンドポイント/入出力/エラー応答を更新したとき
   - `memx_spec_v3/go/service/service.go` のユースケース（ingest/search/get 等）の依存や責務を更新したとき

## 更新ルール
- `generated_at` は更新時刻（UTC, RFC3339）で揃える。
- `depends_on` は「そのノードが成立するために先に読む/参照するノード」を記述する。
- PR 作成前に `index.json` のノード一覧と `caps/*.json` の内容整合を確認する。
