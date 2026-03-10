# ADR-0002: v1必須3エンドポイント固定方針

- Status: Accepted
- Date: 2026-03-04

## Context
- v1 系では後方互換を最優先とする一方、拡張エンドポイントの追加余地があり、必須契約の境界が曖昧になりやすい。
- CLI `--json` と API 契約の同型維持を守るため、最小公開面を固定する必要がある。

## Decision
- v1 MUST の API 契約を以下3エンドポイントに固定する。
  - `POST /v1/notes:ingest`
  - `POST /v1/notes:search`
  - `GET /v1/notes/{id}`
- 上記3件の入力/出力/エラー契約は v1 系で破壊的変更を禁止する。
- `POST /v1/gc:run` は SHOULD（feature flag 前提）として扱い、v1 MUST に昇格しない。

## Consequences
- リリース判定時の必須確認対象が明確化され、契約レビューが短縮される。
- 新規APIは追加してもよいが、既存3契約への影響有無を常に検証する運用が必要になる。
- CLI `--json` 互換テストは、最低でも3エンドポイント相当の整合確認を継続する必要がある。
