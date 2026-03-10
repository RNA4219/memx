# ADR-0003: ErrorCodeとretryable設計（再試行可/不可の境界）

- Status: Accepted
- Date: 2026-03-04

## Context
- `INTERNAL` へのフォールバックを許容する現行方針では、呼び出し側が再試行可否を誤判定しやすい。
- API/CLI/運用監査で再試行ポリシーを統一しないと、不要な再実行や障害拡大を招く。

## Decision
- ErrorCodeごとの retryable 境界を固定する。
  - `INVALID_ARGUMENT` / `NOT_FOUND` / `GATEKEEP_DENY`: 再試行不可。
  - `CONFLICT` / `INTERNAL`: 条件付き再試行可（一時障害のみ）。
- `CONFLICT` / `GATEKEEP_DENY` は service sentinel 実装時のみ返却可能とし、未実装時は `INTERNAL` フォールバックを許容する。
- 一時障害の判定対象を `DB lock` / `LLM timeout` / `upstream 502,503,504` に限定する。

## Consequences
- クライアント実装が ErrorCode ベースで再試行制御を実装しやすくなる。
- sentinel 未実装期間は `INTERNAL` に意味が集約されるため、運用ログで原因分類を補完する必要がある。
- 将来の ErrorCode 拡張時は、retryable 境界と sentinel 実装を同時に更新する運用が必須になる。
