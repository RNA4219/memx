# ADR 運用ガイド

## 目的
- アーキテクチャ判断（ADR）の記録場所と運用ルールを明確化し、設計変更の意思決定経緯を追跡可能にする。

## 更新責任者
- Architecture Owner（不在時は Repo Maintainer 代行）。

## 更新頻度
- ADR 追加・更新が発生した都度（PR 単位）。

## 関連リンク（BLUEPRINT/RUNBOOK/GUARDRAILS/EVALUATION）
- [BLUEPRINT.md](../../BLUEPRINT.md)
- [RUNBOOK.md](../../RUNBOOK.md)
- [GUARDRAILS.md](../../GUARDRAILS.md)
- [EVALUATION.md](../../EVALUATION.md)

## ADR 一覧
- [ADR-0001: 4DB分割と責務境界](./ADR-0001-4db-boundary.md)
- [ADR-0002: v1必須3エンドポイント固定方針](./ADR-0002-v1-required-endpoints.md)
- [ADR-0003: ErrorCodeとretryable設計（再試行可/不可の境界）](./ADR-0003-errorcode-retryable-boundary.md)

## 索引更新手順（ADR追加時の必須更新項目）
1. `docs/ADR/` に `ADR-xxxx-*.md` を追加し、`Status` と `Date` を記載する。
2. 本READMEの「ADR 一覧」にリンクを追記する（番号昇順）。
3. `memx_spec_v3/docs/design.md` の該当節に ADR リンクを追記する。
4. `memx_spec_v3/docs/requirements.md` の該当節に ADR リンクを追記する。
5. 仕様レビュー時は `spec.md -> requirements.md -> design.md -> docs/ADR/*.md` の順で差分確認する。
