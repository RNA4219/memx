---
owner: memx-core
status: active
last_reviewed_at: 2026-03-06
next_review_due: 2026-06-06
---

# memx 要求事項 - Release & Traceability

> 本書は `requirements.md` から分割された一部です。正本は `requirements.md` を参照してください。

## 0-1. Release Scope Matrix

### CLI

| MUST (v1) | SHOULD (v1.x) | FUTURE (v2+) |
| --- | --- | --- |
| `mem in short`, `mem out search`, `mem out show` ✅ | `mem gc short` ✅（`--enable-gc` または `--dry-run` で有効化） | `mem out recall`, `mem working`, `mem tag`, `mem meta`, `mem lineage`, `mem distill`, `mem out context` |

### API

| MUST (v1) | SHOULD (v1.x) | FUTURE (v2+) |
| --- | --- | --- |
| `POST /v1/notes:ingest` ✅, `POST /v1/notes:search` ✅, `GET /v1/notes/{id}` ✅ | `POST /v1/gc:run` ✅（dry_run オプション対応） | Recall/Working/Tag/Meta/Lineage 系 API |

### 実装状況（2026-03-06 更新）

| 領域 | 要件ID | 状態 | 備考 |
| --- | --- | --- | --- |
| CLI v1必須 | `REQ-CLI-001` | ✅ 完了 | in/search/show |
| API v1必須 | `REQ-API-001` | ✅ 完了 | ingest/search/get |
| GC dry-run | `REQ-GC-001` | ✅ 完了 | Phase0/Phase3実装 |
| fail-closed | `REQ-SEC-001` | ✅ 完了 | Gatekeeper実装 |
| エラーモデル | `REQ-ERR-001` | ✅ 完了 | 400/404/409/403/500 |

受け入れ判定で適用する品質ゲートの言語境界は Go を対象とし、Python/Node は現行運用では対象外とする（判定基準・運用コマンドは `docs/QUALITY_GATES.md` を正とする）。

| 区分 | 条件 |
| --- | --- |
| v1必須 | 入出力互換（CLI→API の入出力マッピングが保持されること） |
| v1必須 | エラーコード（入力不備: 400系 / 内部障害: 500系 を返すこと） |
| v1必須 | 最小性能目標（`ingest`/`search`/`show` がローカル単体で実用応答時間を維持すること） |

## 0-2. バージョニングと段階移行

本節は、機能追加・仕様変更・廃止を `MUST(v1)` / `SHOULD(v1.x)` / `FUTURE(v2+)` の3段階で運用するための正本要件とする。
`v1` は後方互換維持を最優先とし、破壊変更は `FUTURE(v2+)` へ隔離して段階移行する。

### 0-2-1. 段階別ルール（許可変更 / 禁止変更 / 廃止条件）

| 区分 | 許可変更 | 禁止変更 | 廃止条件 |
| --- | --- | --- | --- |
| MUST (v1) | 後方互換を維持した拡張のみ追加可能（任意フィールド追加、任意オプション追加、任意パラメータ追加） | 既存 CLI/API 入出力の型・意味・必須性の変更、既存エラーコード削除、既存コマンド/エンドポイント削除、`--json` 既定出力の非同型化 | v1 系では廃止不可。廃止は `FUTURE(v2+)` へ昇格して予告し、`CHANGELOG.md` と `memx_spec_v3/CHANGES.md` に破壊変更チェックリストを記載したうえで次メジャーで実施 |
| SHOULD (v1.x) | 実験機能として追加可能（feature flag 既定 OFF、既定挙動に影響しないこと） | 既定 ON 化、flag なし常時有効化、MUST と同名 I/F の上書き、flag 未指定での出力仕様変更 | まず feature flag を deprecated 扱いにし 1 つ以上のマイナー期間で警告、次メジャーで削除 |
| FUTURE (v2+) | 次メジャー向けに仕様追加・再設計可能（互換フラグ/移行導線付き） | v1 系へ逆流させる破壊的導入（互換フラグなし）、移行手順未定義のままの強制切替 | `v1 -> v2` 移行手順を明示し、互換期間の並行提供方針を定義してから廃止 |

### 0-2-2. エラーコード拡張の昇格条件（service sentinel 連動）

- `CONFLICT` / `GATEKEEP_DENY` / `FEATURE_DISABLED` は、**service 層に対応する sentinel error が実装済みである場合のみ** `INTERNAL` から個別コードへ昇格してよい。
- 適用条件（実装有無との対応）は次の通り。
  - `CONFLICT`: `service.ErrConflict`（同等 sentinel）実装済み時のみ適用。未実装時は `INTERNAL`。
  - `GATEKEEP_DENY`: `service.ErrGatekeepDeny`（同等 sentinel）実装済み時のみ適用。未実装時は `INTERNAL`。
  - `FEATURE_DISABLED`: `service.ErrFeatureDisabled`（同等 sentinel）実装済み時のみ適用。未実装時は `INTERNAL`。
- 昇格時の必須条件は次の通り。
  1. `go/service` に sentinel error を追加し、再試行可否の意味が固定されていること。
  2. `go/api/errors.go`（または同等の `mapError`）に明示マッピングを追加すること。
  3. CLI `--json` 出力に同型の `code` が反映されること。
  4. 本要件書と変更履歴（`CHANGELOG.md` / `memx_spec_v3/CHANGES.md`）へ昇格理由を記録すること。
- 上記を満たさない段階では `INTERNAL` フォールバックを維持する。

### 0-2-3. CLI `--json` と API レスポンス同型維持の例外条件

- 既定動作では、CLI `--json` は API レスポンスと**同型（同一キー体系・同一意味）**を維持しなければならない。
- 同型を崩せるのは、利用者が互換性逸脱を許可する**明示フラグ**を指定した場合に限定する（例: 互換オフ/人間可読優先モード）。
- 明示フラグ未指定時に、CLI 側都合のみでフィールド名変更・構造変更・意味変更を行ってはならない。
- 例外モードを導入した場合も、API の canonical 形は維持し、CLI ヘルプに「非互換モード」であることを明示する。

#### 例外適用の最小条件

- 例外は次の 3 条件を**すべて**満たす場合に限り許可する。
  1. 既定動作では同型性を保持し、明示フラグ指定時のみ非同型を許可する。
  2. API 側の canonical schema（キー名・型・意味）は変更しない。
  3. 非同型モードの利用目的（可読性優先・デバッグ用途など）を CLI ヘルプ/変更履歴に記載する。
- 上記 3 条件のいずれかを満たさない場合は、例外を認めず同型を維持する。

### 0-2-4. 破壊変更時の必須チェックリスト

- 次のいずれかに該当する変更は「破壊変更」とみなし、マージ前にチェックリスト完了を必須とする。
  - CLI/API の既存必須フィールド削除、型変更、意味変更
  - 既存コマンド/エンドポイント/エラーコードの削除または互換なし改名
  - `--json` 既定出力の同型性を崩す変更
- 破壊変更時は、`CHANGELOG.md` と `memx_spec_v3/CHANGES.md` の双方に、同一日付で以下を必ず記載する。
  - 対象 I/F、変更種別、影響範囲、移行先、移行期限、移行手順、互換期間中の挙動
  - 「明示フラグでのみ新挙動を有効化する」かどうか

### 0-2-5. 破壊変更が必要な場合の Task Seed 追記テンプレート

破壊変更を含むタスクは `docs/TASKS.md` の必須項目に加えて、以下テンプレートを Task Seed 本文へ追記する。

```md
## Breaking Change Addendum

### Impacted Interface
- CLI/API 名称:
- 互換性影響（削除/型変更/意味変更）:

### Migration Plan
- 互換フラグ名（必須）:
- 既定値:
- 有効化条件:
- 利用者移行手順:
- 互換期間の終了条件:

### Checklist
- [ ] Source は `path#Section` で記載済み（`docs/TASKS.md` 準拠）
- [ ] Node IDs を記載済み（依存照合対象なら必須）
- [ ] Requirements に後方互換/非機能制約を明記済み
- [ ] エラーコード変更時は `memx_spec_v3/docs/requirements.md` と `memx_spec_v3/docs/error-contract.md` を更新対象に含めた
- [ ] Commands に検証コマンドを順序付きで記載済み
- [ ] Release Note Draft を記載済み
- [ ] `memx_spec_v3/CHANGES.md` と `CHANGELOG.md` への反映項目を記載済み
- [ ] `Status: done` 前に `Moved-to-CHANGES: YYYY-MM-DD` を追記する
```

- 本テンプレートは `docs/TASKS.md` の「Task Seed 必須項目」「CHANGES 連携ルール」と矛盾しないことを必須条件とする。
- `Source` にはテンプレートID（例: `IN-YYYYMMDD-001`, `IN-202603xx-001`）および `TBD` を記載してはならない。
- インシデントを要件根拠に使う場合、`Source` で許可するのは `docs/IN-<実日付>-<連番>.md` のみとする。`docs/IN-BASELINE.md` / `docs/IN-YYYYMMDD-001.md` / `docs/IN-202603xx-001.md` はテンプレート資料扱いで、根拠参照は禁止する。
- セキュリティ/保持に関わる Task Seed では、Requirements 欄に `REQ-SEC-001` / `REQ-RET-001` / `REQ-SEC-AUD-001` / `REQ-SEC-AUD-002` / `REQ-SEC-GRD-001` の該当IDを必ず列挙する。
- 特に `docs/TASKS.md` の `Requirements` / `Release Note Draft` / `Status: done` 条件（`Moved-to-CHANGES`）と同一基準で運用する。

---

<a id="requirements-traceability"></a>

## 0-3. 要件トレーサビリティ

<a id="主要要件id固定"></a>

### 主要要件ID（CLI/API/GC/Security/Error 固定）

| 要件領域 | Requirement ID | 受入基準（期待結果） | 検証コマンド（RUNBOOK 1:1） | EVALUATION 相互参照 |
| --- | --- | --- | --- | --- |
| CLI | `REQ-CLI-001` | `mem out search` の `--json` 出力が API 契約と同型で返る。 | [`trace-req-cli-001`](../../RUNBOOK.md#trace-req-cli-001) | [REQ-CLI-001](../../EVALUATION.md#req-cli-001-passfail) |
| API | `REQ-API-001` | `POST /v1/notes:ingest` が v1 契約（入力/出力/HTTP）を維持する。 | [`trace-req-api-001`](../../RUNBOOK.md#trace-req-api-001) | [REQ-API-001](../../EVALUATION.md#req-api-001-passfail) |
| GC | `REQ-GC-001` | `mem gc short --dry-run` が DB 非更新で判定結果のみ返す。 | [`trace-req-gc-001`](../../RUNBOOK.md#trace-req-gc-001) | [REQ-GC-001](../../EVALUATION.md#req-gc-001-passfail) |
| Security | `REQ-SEC-001` | `sensitivity=secret` 相当入力を fail-closed（保存禁止）で拒否する。 | [`trace-req-sec-001`](../../RUNBOOK.md#trace-req-sec-001) | [REQ-SEC-001](../../EVALUATION.md#req-sec-001-passfail) |
| Error | `REQ-ERR-001` | `NOT_FOUND`/`INVALID_ARGUMENT`/`INTERNAL` の契約と再試行可否が整合する。 | [`trace-req-err-001`](../../RUNBOOK.md#trace-req-err-001) | [REQ-ERR-001](../../EVALUATION.md#req-err-001-passfail) |

<a id="task-seed-source-fixed"></a>

### Task Seed 転記用固定表（Source / Requirements 直接引用用・表形式固定）

| Requirement ID | Source（Task Seed で直接参照する固定値） | Requirements（引用用） |
| --- | --- | --- |
| `REQ-CLI-001` | `memx_spec_v3/docs/requirements-cli.md#3-cli-要件` | `CLI v1必須3コマンドのJSON互換を維持し、mem out search の --json 出力は API 契約と同型を維持する。` |
| `REQ-API-001` | `memx_spec_v3/docs/requirements-api.md#6-api-要件` | `API v1必須3エンドポイント契約を維持し、POST /v1/notes:ingest は v1 契約を維持する。` |
| `REQ-GC-001` | `memx_spec_v3/docs/requirements-cli.md#3-5-mem-gc-short` | `GC dry-run/閾値判定/DB非更新契約を満たし、mem gc short --dry-run は DB を更新せず判定結果のみ返す。` |
| `REQ-SEC-001` | `memx_spec_v3/docs/requirements-data-model.md#2-7-security--retention-requirements` | `sensitivity判定をfail-closedで適用し、secret は fail-closed で拒否する。` |
| `REQ-RET-001` | `memx_spec_v3/docs/requirements-data-model.md#2-7-security--retention-requirements` | `archive退避/削除と監査ログ要件を満たす。` |
| `REQ-SEC-AUD-001` | `memx_spec_v3/docs/requirements-data-model.md#2-7-2-actor--approval--audit-責任分界表` | `archive_moveの監査ログ固定項目（証跡ファイルパス/必須キー/保持期間）を満たす。` |
| `REQ-SEC-AUD-002` | `memx_spec_v3/docs/requirements-data-model.md#2-7-2-actor--approval--audit-責任分界表` | `archive_purgeの監査ログ固定項目（証跡ファイルパス/必須キー/保持期間）を満たす。` |
| `REQ-SEC-GRD-001` | `memx_spec_v3/docs/requirements-data-model.md#2-7-5-guardrails-fail-closed-との整合チェック要件` | `GUARDRAILS fail-closed整合チェックを満たす。` |
| `REQ-ERR-001` | `memx_spec_v3/docs/requirements-api.md#6-4-エラーモデル` | `ErrorCode契約とretryableルール（再試行可否）を維持する。` |
| `REQ-NFR-001` | `memx_spec_v3/docs/requirements-nfr.md#5-1-性能目標` | `性能閾値（ingest/search/show）を満たす。` |