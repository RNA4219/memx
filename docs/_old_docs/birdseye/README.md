# Birdseye インデックス運用

本書は HUB（`HUB.codex.md`）から参照する Birdseye 入口ドキュメント（canonical path: `docs/birdseye/README.md`）である。

`docs/birdseye/index.json` は memx_spec_v3 の主要ドキュメント/実装ノードの依存関係を管理する HUB です。

## 構成
- `index.json`: HUB とノード一覧（`node_id`, `role`, `depends_on`, `generated_at` を保持）
- `caps/*.json`: 主要文書ごとの capsule（要点・依存ノードを保持）
- `hot.json`: `optional`（現行は未運用。導入時のみ生成/参照対象）

## `hot.json` の運用前提（HUB/RUNBOOK 整合）
- 現行運用（既定）: `docs/birdseye/hot.json` は `optional`（未運用）で、欠損していても処理継続可能。
- 欠損時の記録: `notes.readiness_status=ready` を維持しつつ、`notes.missing_files` に `docs/birdseye/hot.json` を記録する。
- 移行運用: `hot.json` を運用開始する場合のみ、再生成手順に `--emit hot` を追加する。

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

## index参照とcaps実体の整合チェック手順
1. `index.json` に記載された `nodes[].capsule` を棚卸しする。
   ```bash
   jq -r '.nodes[].capsule' docs/birdseye/index.json | sort -u
   ```
2. 棚卸し結果を順に実体確認し、欠落ファイルを抽出する。
   ```bash
   while read -r cap; do
     if [ -f "$cap" ]; then
       echo "OK  $cap"
     else
       echo "MISSING  $cap"
     fi
   done < <(jq -r '.nodes[].capsule' docs/birdseye/index.json)
   ```
3. 欠落が 1 件でもあれば不整合として扱い、`docs/TASKS.md` の Status 運用に従って `blocked` へ遷移し、Task Seed に欠落 capsule 一覧を記録する。

### 現在の棚卸し結果（`index.json` 기준）
- 欠損 0 件（`MISSING` なし）

## 変更内容と capsule 更新対応表

| 変更ファイル | 更新する capsule |
| --- | --- |
| `memx_spec_v3/docs/requirements.md` | `docs/birdseye/caps/requirements.json` |
| `memx_spec_v3/docs/quickstart.md` | `docs/birdseye/caps/quickstart.json` |
| `memx_spec_v3/go/service/service.go` | `docs/birdseye/caps/service.json` |
| `memx_spec_v3/go/api/http_server.go` | `docs/birdseye/caps/api.json` |
| `GUARDRAILS.md` | `docs/birdseye/caps/GUARDRAILS.md.json` |
| `RUNBOOK.md` | `docs/birdseye/caps/RUNBOOK.md.json` |
| `CHECKLISTS.md` | `docs/birdseye/caps/CHECKLISTS.md.json` |
| `EVALUATION.md` | `docs/birdseye/caps/EVALUATION.md.json` |
| `governance/policy.yaml` | `docs/birdseye/caps/governance.policy.yaml.json` |
| `governance/prioritization.yaml` | `docs/birdseye/caps/governance.prioritization.yaml.json` |



## 更新責任者・更新トリガ
- 更新責任者: memx-core（変更PRの作成者が更新を実施）
- 更新トリガ: PR 作成前（必須）
- 更新トリガ: 週次定期メンテナンス（少なくとも週1回）


## 復旧手順の参照先
- 正本参照: `./workflow-cookbook/tools/codemap/README.md`
- 別名パス: `tools/codemap/README.md`（注記扱い）
