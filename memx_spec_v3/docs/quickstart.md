# memx v1.3 Quickstart（最小）

> この ZIP は「要件＋参照実装の骨格」です。
> `short.db` の ingest/search/get が最小で動くようにしてあります。

## 1) API サーバー起動

```bash
go run ./memx_spec_v3/go/cmd/mem api serve --addr 127.0.0.1:7766 --short ./short.db
```

代替表記（`go` ディレクトリへ移動して実行）:

```bash
cd memx_spec_v3/go
go run ./cmd/mem api serve --addr 127.0.0.1:7766 --short ./short.db
```

## 2) ノート投入

```bash
echo "hello memx" | go run ./memx_spec_v3/go/cmd/mem in short --title "test" --stdin --api-url http://127.0.0.1:7766
```

代替表記（`go` ディレクトリへ移動して実行）:

```bash
cd memx_spec_v3/go
echo "hello memx" | go run ./cmd/mem in short --title "test" --stdin --api-url http://127.0.0.1:7766
```

## 3) 検索

```bash
go run ./memx_spec_v3/go/cmd/mem out search "hello" --api-url http://127.0.0.1:7766
```

代替表記（`go` ディレクトリへ移動して実行）:

```bash
cd memx_spec_v3/go
go run ./cmd/mem out search "hello" --api-url http://127.0.0.1:7766
```

## 4) 単体表示

```bash
go run ./memx_spec_v3/go/cmd/mem out show <NOTE_ID> --api-url http://127.0.0.1:7766
```

代替表記（`go` ディレクトリへ移動して実行）:

```bash
cd memx_spec_v3/go
go run ./cmd/mem out show <NOTE_ID> --api-url http://127.0.0.1:7766
```

## 補足

- `--api-url` を省略すると CLI は in-proc で Service を呼びます（HTTP を経由しない）。
- SQLite ドライバは `modernc.org/sqlite` を想定しています。環境に合わせて差し替えてください。

## 関連ドキュメント
- 仕様参照の起点: `./spec.md`
- 要件正本: `./requirements.md`（本書 `quickstart.md` は運用手順の補助文書）
- API 契約正本: `./contracts/openapi.yaml`
- CLI `--json` 契約正本: `./contracts/cli-json.schema.json`
- エラー契約要約: `./error-contract.md`
- v1必須3エンドポイント契約表: `./requirements.md` の「6-3-1. v1必須3エンドポイント契約（`requirements.md` × `go/api/types.go` 照合）」
