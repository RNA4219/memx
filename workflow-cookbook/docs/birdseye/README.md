# Birdseye データセット運用ガイド

Birdseye は、Workflow Cookbook の知識マップを 3 層（Index / Capsules / Hot）で提供します。
Guardrails の「Bootstrap → Index → Caps」という読込順序を支えるデータセットとして、各層の成果物と鮮度管理ルールをこのディレクトリで管理します。

## ディレクトリ構成

- `index.json`
  - 役割: Birdseye ノード一覧と隣接関係（Edges）の基盤データ
  - Guardrails 連携: `plan` や `notes` で ±hop を抽出する際の一次ソース
- `caps/`
  - 役割: 各ノードのカプセル要約（`<path>.json` を 1 ノード 1 ファイルで保持）
  - Guardrails 連携: `deps_out`・公開 API・リスク情報を提供
- `hot.json`
  - 役割: 主要ノードのホットリスト
  - Guardrails 連携: 「頻出入口ホットリスト」を満たし、即時参照を補助
- `README.md`
  - 役割: データセット運用手順
  - Guardrails 連携: Birdseye の生成・検証フローを共有
- `../BIRDSEYE.md`
  - 役割: フォールバック用の人間向け導線
  - Guardrails 連携: 最終手段として Edges / Hot / 更新手順を提示

## `codemap.update` 実行手順

Birdseye の生成・更新は `tools/codemap/update.py` を介して行います。
Guardrails が要求する「自動再生成＋鮮度確認」を満たすため、以下のコマンドを基準としてください。

```bash
python tools/codemap/update.py \
  --targets docs/birdseye/index.json,docs/birdseye/hot.json \
  --emit index+caps
```

1. 対象ノード（`--targets`）に今回更新したファイルや重要エントリをカンマ区切りで列挙します。
   Birdseye を再生成する場合は `docs/birdseye/` 配下を明示してください。
2. `--emit` で出力対象を指定します。現在は `index+caps` が標準です。
3. `docs/birdseye/index.json` と `docs/birdseye/hot.json` を同一ターゲットで指定すると、両データセットの鮮度が揃います。
   `generated_at` と `index.json.nodes[*].mtime` は ISO 8601 UTC（`YYYY-MM-DDTHH:MM:SSZ`）で出力されるため、
   手動確認時も同フォーマットであることを先に確認してください。
   出力後は `index.json.generated_at` / `hot.json.generated_at` が最新コミットに追随しているか確認し、
   必要に応じてホットリスト項目の `last_verified_at` が対象ノードの最新確認日を反映しているか点検します。
4. 差分をレビューし、`docs/BIRDSEYE.md` のフォールバック情報と矛盾がないことをチェックしてからコミットします。

> 手動編集が必要な場合でも、Birdseye スキーマ（`id`・`role`・`caps`・`edges` など）とパス命名規則（`/` を `.` に置換）を崩さないでください。

## ホットリストと鮮度管理

- `docs/birdseye/hot.json` は `index.json` の再生成時に自動同期され、`refresh_command` と `index_snapshot`
  で更新履歴を追跡します。必要に応じて各ホットリスト項目に記録された `last_verified_at` を確認し、
  対象ノードの確認日が反映されているかを点検します。
- ホットリストの構成は `README.md`・`GUARDRAILS.md`・`HUB.codex.md` など主要導線を中心に選定し、Workflow Cookbook 本流のホットリスト基準を参考に `edges` を明示しています。
- `docs/BIRDSEYE.md` ではホットリストの概要と Edges を人間が参照できるよう再掲しているため、変更時は両ドキュメントの整合性を確認します。

## Guardrails との整合

- `GUARDRAILS.md` が定義する「Birdseye JSON を第一読者とし、人間向けはフォールバック」という方針を満たすため、JSON ファイルは常に最新の情報源となるよう保守します。
- フォールバック時には `docs/BIRDSEYE.md` の Edges / Hot / 更新手順を参照し、必要に応じて本 README の `codemap.update` 手順に合流してください。
- Birdseye を更新した場合は、関連するチェックリストや運用ドキュメント（`CHECKLISTS.md`・`RUNBOOK.md` など）にも鮮度情報を反映し、リポジトリ全体の整合を保ちます。
