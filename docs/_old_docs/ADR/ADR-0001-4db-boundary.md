# ADR-0001: 4DB分割と責務境界

- Status: Accepted
- Date: 2026-03-04

## Context
- memx v1.3 では `short/journal/knowledge/archive` の4ストアを前提に運用しているが、設計意図を1箇所で固定していなかった。
- ストア間で責務が混在すると、GC・検索対象・保持期間の運用判断が不安定になる。

## Decision
- 4DB分割を正式採用し、責務境界を以下で固定する。
  - `short.db`: 一次投入と短期保持の起点。
  - `journal.db`: 時系列ログ。
  - `knowledge.db`: 時系列非依存の知識。
  - `archive.db`: 退避保管（通常検索対象外）。
- 共通テーブルは `notes/tags/note_tags` を必須、`note_embeddings/notes_fts` は store 特性に応じて採用可とする。
- `short.db` のみ `short_meta/lineage` を必須とし、GC起点と系譜追跡の責務を持たせる。

## Consequences
- スキーマ変更時に「どの store の責務か」を判断しやすくなり、非互換変更の混入を抑制できる。
- archive の通常検索除外を前提に、運用監査と復元導線を分離して設計できる。
- store 跨ぎ機能は境界越えルール（昇格・退避）を明示する必要がある。
