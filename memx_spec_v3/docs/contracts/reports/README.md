# Contract Alignment Reports 運用 README

このディレクトリは Phase 3（契約整合）の正本レポート保存先とする。

## 1. 必須成果物（固定）
- `CONTRACT-ALIGN-YYYYMMDD-###.md`（詳細レポート）
- `LATEST.md`（直近レポートへのポインタ）

上記以外は補助ファイル扱いとし、運用上の必須成果物に数えない。

## 2. 命名規則
### 2.1 詳細レポート
- 形式: `CONTRACT-ALIGN-YYYYMMDD-###.md`
- `YYYYMMDD`: 判定実施日（ローカル日付）
- `###`: 3桁ゼロ埋め連番（`001` 開始）

例:
- `CONTRACT-ALIGN-20260304-001.md`
- `CONTRACT-ALIGN-20260304-002.md`

### 2.2 最新ポインタ
- ファイル名は固定で `LATEST.md` とする。
- 直近 1 件の `report_id`・判定・参照先を保持する。

## 3. 連番採番規則
1. 当日（`YYYYMMDD`）の既存 `CONTRACT-ALIGN-YYYYMMDD-###.md` を列挙する。
2. 最大連番を取得する（存在しない場合は `000` とみなす）。
3. 次回作成時は `+1` した 3 桁番号を採番する。
4. 既存番号の再利用・上書きは禁止。

## 4. 最新レポート探索ルール
最新レポートの参照は次の優先順位で解決する。

1. `LATEST.md` の `report_id` と `report_path` を正本として利用する。
2. `LATEST.md` が欠落または不整合の場合、`CONTRACT-ALIGN-*.md` を日付降順・連番降順で探索し、先頭を暫定最新とする。
3. 暫定最新を採用した場合は `LATEST.md` を即時再生成し、整合状態へ戻す。

## 5. `LATEST.md` 必須項目
`LATEST.md` は契約整合レポート作成のたびに毎回更新し、以下 5 キーを必須で記載する。

```md
report_id: CONTRACT-ALIGN-YYYYMMDD-###
report_path: ./CONTRACT-ALIGN-YYYYMMDD-###.md
decision_date: YYYY-MM-DD
high_count: <number>
phase3_status: done | blocked | reviewing
```

`phase3_status` は `docs/TASKS.md` と同一語彙（`done/blocked/reviewing`）に合わせる。
1 キーでも未更新・欠落がある場合は運用違反として `fail` とする。
