---
intent_id: DOC-LEGACY
owner: docs-core
status: active
last_reviewed_at: 2025-10-28
next_review_due: 2025-11-28
---

# Birdseye リファレンス

Birdseye は、Workflow Cookbook の知識マップを統合的に参照する仕組みです。
`index.json`・`caps/*.json`・`hot.json` を併せて読み解くことで、主要ノードの関係性と鮮度を一元確認できます。
本書は Guardrails からのフォールバック参照起点として、Edges・ホットリスト・更新手順を 1 か所に整理します。

## Edges（主要ノードの隣接関係）

`docs/birdseye/index.json` の `edges` から、Guardrails がフォールバック時に ±1 hop を推定しやすいよう主要ノードを抜粋しています。
`docs/BIRDSEYE.md` へのリンクは `docs/birdseye/index.json` にも `README.md`／`GUARDRAILS.md`／`tools/codemap/README.md` との双方向エッジと登録され、
ホットリスト専用のフォールバック導線に加えて JSON からも辿れる。

- `README.md`
  - 主要 Edges: `GUARDRAILS.md`, `HUB.codex.md`, `BLUEPRINT.md`, `RUNBOOK.md`, `EVALUATION.md`, `docs/birdseye/index.json`
  - 用途: 初動ガイドと Birdseye 読込順序の提示
- `GUARDRAILS.md`
  - 主要 Edges: `README.md`, `HUB.codex.md`, `RUNBOOK.md`, `EVALUATION.md`, `docs/birdseye/index.json`
  - 用途: 行動指針・鮮度管理・フォールバック手順
- `HUB.codex.md`
  - 主要 Edges: `README.md`, `GUARDRAILS.md`, `RUNBOOK.md`, `EVALUATION.md`, `CHECKLISTS.md`, `docs/birdseye/index.json`
  - 用途: 仕様・タスクの依存ハブ
- `RUNBOOK.md`
  - 主要 Edges: `README.md`, `GUARDRAILS.md`, `EVALUATION.md`, `docs/IN-20250115-001.md`
  - 用途: 運用および Birdseye 再生成 SOP
- `EVALUATION.md`
  - 主要 Edges: `README.md`, `GUARDRAILS.md`, `RUNBOOK.md`, `CHECKLISTS.md`
  - 用途: 受入基準と品質メトリクス
- `docs/birdseye/index.json`
  - 主要 Edges: `README.md`, `GUARDRAILS.md`, `HUB.codex.md`, `docs/birdseye/caps/`, `tools/codemap/`
  - 用途: Birdseye hop 計算の基盤
- `tools/codemap/README.md`
  - 主要 Edges: `tools/codemap/update.py`, `docs/birdseye/index.json`, `docs/birdseye/caps/`
  - 用途: 再生成コマンドと契約

> 詳細なエッジリストは `docs/birdseye/index.json` を参照してください。フォールバック中でも JSON を第一読者とし、ここは要約に留めます。

## Hot List（主要ノードの即時参照）

`docs/birdseye/hot.json` に定義されたホットリストの要旨です。鮮度確認や調査時の着手順に活用してください。

- `README.md`: リポジトリの導入手順と Birdseye 読込フローの起点。
- `GUARDRAILS.md`: 行動指針と鮮度管理ポリシー。必ず Birdseye の再生成条件を確認。
- `HUB.codex.md`: 依存関係とタスク分割のハブ。±hop の整合確認に必須。
- `RUNBOOK.md`: 標準オペレーションと再生成ステップ。
- `EVALUATION.md`: 受入基準と品質観点の参照元。
- `BLUEPRINT.md`: 機能要件・境界の定義。
- `CHECKLISTS.md`: デリバリー確認項目と鮮度チェックの連携。
- `docs/BIRDSEYE.md`: 人間向けフォールバック。JSON が取得できない場合の最終ライン。
- `docs/birdseye/index.json`: hop 抽出の基盤データ。`hot.json` と同じターゲットで更新。
- `tools/codemap/README.md`: `codemap.update` の契約とパラメータ。

ホットリストの `generated_at` は、`index.json` 再生成のたびに同じターゲットで更新し、鮮度を揃えてください。必要に応じてホットリスト項目の `last_verified_at` を更新対象に含めます。

`docs/birdseye/*.json` の `generated_at` は、初回生成が古いものから 5 桁ゼロ埋めの連番（例: `00001`, `00002`）で管理し、再生成時も番号を引き継いでください。

## Birdseye 更新手順

1. 変更対象や鮮度が落ちたノードを整理し、`codemap.update` のターゲットに指定します。
2. 以下を実行し、カンマ区切りで指定した `index.json` と `hot.json` を同時にターゲットへ指定しつつ、`index` を出力します
   （`caps` も自動生成されますが追加の手作業は不要です）。

   ```bash
   python tools/codemap/update.py \
     --targets docs/birdseye/index.json,docs/birdseye/hot.json \
     --emit index+caps
   ```

3. `docs/birdseye/index.json.generated_at` と `docs/birdseye/hot.json.generated_at` が最新コミットに追随しているか確認します。
4. `docs/birdseye/hot.json` の `refresh_command` と `index_snapshot` が現行手順を反映しているか点検します。
   `index.json` と `hot.json` を同一ターゲットで再生成し、両者の鮮度を揃えてください。
5. ホットリスト項目に `last_verified_at` を保持している場合は、対象ノードの最新確認日を反映しているか確認します。
6. `GUARDRAILS.md` の [鮮度管理](../GUARDRAILS.md#鮮度管理staleness-handling) に従って、必要に応じて人間へ再生成依頼またはインシデント共有を行います。

## フォールバック運用

- 自動ツールが利用できない場合は、上記 Edges と Hot List を参考に読込対象を最小化しつつ、
  Birdseye JSON を第一読者として `docs/birdseye/index.json` → `caps/*.json` → `docs/birdseye/hot.json` の順に確認してください。
- JSON が取得できない場合でも、`docs/BIRDSEYE.md` に記載された Edges/Hot/更新手順を用いて暫定判断を行い、
  可能な限り早く `tools/codemap/update.py` を実行できる環境へエスカレーションします。
- インシデントレベルの齟齬や破損が見つかった場合は `docs/IN-20250115-001.md` の手順で共有し、`RUNBOOK.md` の標準オペレーションに沿って復旧します。

> ここに記載した情報は JSON の要約であり、最新状態は常に `docs/birdseye/index.json`・`docs/birdseye/hot.json`・`docs/birdseye/caps/` を参照してください。
