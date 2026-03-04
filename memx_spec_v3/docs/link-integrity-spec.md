# Link Integrity Spec

## 1. 目的
本仕様は、設計文書群の章内リンク・章間リンク・運用リンクの到達性を検証し、  
Phase 3 のリンク健全性判定を再現可能な基準で実施するための手順を定義する。

また、受け入れレビューで使用するリンク不達件数および詳細レポートの算出方法を定義する。

判定語彙（重大度・Status 遷移）は `docs/TASKS.md` と同一語彙を使用する。

---

# 2. 検証対象

## 2.1 入力対象

以下の設計文書群を対象とする。

固定対象:

- `memx_spec_v3/docs/design.md`
- `memx_spec_v3/docs/interfaces.md`
- `memx_spec_v3/docs/requirements.md`
- `memx_spec_v3/docs/traceability.md`
- `memx_spec_v3/docs/operations-spec.md`
- `RUNBOOK.md`
- `EVALUATION.md`

ワイルドカード対象:

- `memx_spec_v3/docs/*.md`
- `docs/*.md`
- `docs/birdseye/index.json`

---

## 2.2 リンク分類

抽出された参照リンクは次の3種類に分類して検証する。

### 1. 章内アンカー（Markdown Anchor）

例:

```

#phase-3
./requirements.md#2-2-変更タイプ別チェックリスト

```

### 2. 相対パス（Relative Path）

例:

```

../memx_spec_v3/docs/design.md
./IN-20260303-001.md

```

### 3. 契約 JSON Pointer（Contract JSON Pointer）

例:

```

memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1memories/get
memx_spec_v3/docs/contracts/cli-json.schema.json#/properties/output

```

---

# 3. リンク不達の定義

リンク不達（unreachable link）とは、次のいずれかを満たすリンクである。

- 相対パス先ファイルが存在しない
- Markdownアンカーが存在しない
- 参照先セクション名が一致しない
- JSON Pointer が解決できない
- Pointer が RFC6901 構文を満たさない
- 相対パス解決後にリポジトリルート外へ逸脱する

---

# 4. 判定ルール

## 4.1 章内アンカー

条件:

- 同一ファイルまたはリンク先 Markdown 内に対象見出しが存在すること
- 見出しは次の正規化後に一致判定する

正規化:

- 小文字化
- 空白 → `-`
- 記号削除

---

## 4.2 相対パス

条件:

- 正規化後のパスが存在する
- リポジトリルート配下に収まる
- Markdown / JSON 参照の場合

許可拡張子:

```

.md
.json

```

---

## 4.3 契約 JSON Pointer

条件:

- 参照先ファイルが存在する
- JSON/YAML としてパース可能
- Pointer が **1件のノードに解決可能**

Pointer仕様:

- RFC 6901 準拠
- エスケープ

```

~0
~1

```

---

# 5. 失敗条件

以下をリンク健全性チェックの失敗とする。

| failure_type | 内容 |
|---|---|
| missing_target | 相対パス先ファイル不存在 |
| unresolved_anchor | Markdownアンカー未解決 |
| bad_extension | 不正拡張子 |
| path_escape | リポジトリルート外参照 |
| pointer_not_found | JSON Pointer 未解決 |
| pointer_syntax_error | JSON Pointer 構文誤り |

---

# 6. 重大度（severity）

### high

- 相対パス階層逸脱
- JSON Pointer 未解決
- JSON Pointer 構文誤り
- 主要仕様リンクの参照先不存在

### medium

- アンカー未解決
- 補助説明リンクの参照先不存在

### low

- 拡張子誤り
- 軽微な表記揺れ
- 同一実体への重複リンク

---

# 7. 判定

## 7.1 受け入れ判定

| status | 条件 |
|---|---|
| pass | link_unreachable_count = 0 |
| fail | link_unreachable_count ≥ 1 |

---

# 8. Status 遷移規定

`docs/TASKS.md` 語彙を使用。

許可Status:

```

planned
active
in_progress
reviewing
blocked
done

```

遷移規定:

1. `reviewing -> blocked`
   - `severity: high` が 1件以上

2. `reviewing` 継続
   - `severity: medium` のみ残存

3. `reviewing -> done`
   - high=0
   - medium=0
   - low=0

4. `in_progress -> reviewing`
   - 本仕様によるリンクチェック結果を Task Seed に反映済

---

# 9. 出力

## 9.1 集計出力

```

link_total
link_unreachable_count

````

---

## 9.2 詳細レポート

```yaml
check_id: LI-YYYYMMDD-001
severity: low | medium | high
status_transition: reviewing->blocked
source_ref: memx_spec_v3/docs/design.md#L10-L40
target_ref: ../docs/unknown.md#phase-1
failure_type: missing_target | unresolved_anchor | bad_extension | path_escape | pointer_not_found | pointer_syntax_error
action: "リンク先追加または参照先修正"
````

---

# 10. Phase 3 運用

`orchestration/memx-design-docs-authoring.md` の Phase 3 では
本仕様をリンク健全性チェックの正規基準として参照する。

チェック結果は以下に反映する。

* Task Seed `Requirements`
* Task Seed `Commands`
* Task Seed `Status`