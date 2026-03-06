# Contract Alignment Spec

## 1. 比較対象
契約同期の比較対象は以下の5ファイルとする。

- `memx_spec_v3/docs/design.md`
- `memx_spec_v3/docs/interfaces.md`
- `memx_spec_v3/docs/contracts/openapi.yaml`
- `memx_spec_v3/docs/contracts/cli-json.schema.json`
- `memx_spec_v3/docs/error-contract.md`

## 2. 差分分類
契約差分は以下の3分類で判定する。

### 2.1 互換維持差分（severity: low）
- 文言補足・注記追加・例示追加のみ
- 必須項目、型、意味論、既存入出力に変更なし
- 既存クライアント・既存運用手順へ影響なし

### 2.2 要注意差分（severity: medium）
- 説明不一致（`design.md` と `interfaces.md`、または契約ファイル間で記述が食い違う）
- 実装上は互換の可能性があるが、解釈が分かれる状態
- 仕様レビューで整合コメントまたは補正が必要

### 2.3 破壊差分（severity: high）
- 必須項目の削除
- 型変更（例: `string` → `integer`）
- 既存 `ErrorCode` の意味変更
- 後方互換を損なうレスポンス/CLI出力変更

## 3. REQ-ID 連動ルール
差分判定は要件IDと必ず紐づける。

- API 契約差分は `REQ-API-001` に連動させる
  - 主対象: `openapi.yaml` と `design.md` / `interfaces.md`
- CLI JSON 契約差分は `REQ-CLI-001` に連動させる
  - 主対象: `cli-json.schema.json` と `design.md` / `interfaces.md`
- エラー契約差分は `REQ-ERR-001` に連動させる
  - 主対象: `error-contract.md` と API/CLI のエラー記述

判定ルール:
1. 差分ごとに `affected_req` を最低1件付与する。
2. `severity: high`（破壊差分）は必ず該当 REQ-ID を明示し、互換性対応方針（移行・非推奨期間・代替）を `action` に記載する。
3. REQ-ID が特定できない差分は `severity: medium` とし、先に要件トレーサビリティの補完タスクを起票する。

## 4. 判定結果の出力フォーマット
契約同期チェック結果は、差分1件につき以下フォーマットで出力する。

```yaml
diff_id: CA-YYYYMMDD-001
severity: low | medium | high
affected_req:
  - REQ-API-001
source_refs:
  - memx_spec_v3/docs/design.md#L120-L150
  - memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1memories/get
action: "文言統一のみ実施。互換性影響なし。"
```

フィールド定義:
- `diff_id`: 契約差分を一意に識別するID
- `severity`: 差分分類（low/medium/high）
- `affected_req`: 影響する REQ-ID（複数可）
- `source_refs`: 判定根拠の参照元（ファイル行またはスキーマパス）
- `action`: 解消方針・保留理由・移行手順


## 4.1 共通メタキー追記ルール
`memx_spec_v3/docs/design-evidence-schema-spec.md` 準拠で、4章の出力レコードには以下キーを必須追記する。

- `run_id`
- `generated_at`
- `source_commit`
- `chapter_id`
- `tool`
- `status`
- `severity_summary`
- `evidence_paths`

`severity_summary` は差分全体の件数サマリを保持し、各 `diff_id` レコードは `evidence_paths` に根拠ファイルを列挙する。

## 5. 運用ルール
- Phase 3 の契約同期判定は本仕様を基準に実施する。
- 判定結果は章別 Task Seed の `Requirements` / `Commands` に反映する。
- `severity: high` が 1 件でも残る場合、Phase 3 は未完了扱いとする。
