# memx-core Interfaces

## 1. Purpose
本書は、memx-core の最小公開インターフェースを定義する。  
memx-core は **Memory Substrate** として、Evidence / Knowledge / Artifact / Lineage に対する保存・取得・検索・要約・退避の操作を提供する。

---

## 2. Design Rules

### 2.1 Interface Principles
- インターフェースは task/work 概念を持たない
- 外部 tracker 概念を持たない
- すべての entity は stable ID を持つ
- 参照は typed ref を使う
- raw data と distilled knowledge を分離する
- 破壊的変更より append / new version を優先する

### 2.2 Typed Ref Format
他 repo からの参照は以下の形式を使う。

```text
memx:<entity_type>:<id>
````

例:

```text
memx:evidence:01HXXXXXXX
memx:knowledge:01HYYYYYYY
memx:artifact:01HZZZZZZZ
memx:lineage:01HAAAAAAA
```

---

## 3. Entity Types

### 3.1 Evidence

生データの親 entity。

Fields:

* id
* store
* kind
* title
* source_uri
* source_hash
* recorded_at
* metadata_json
* created_at
* updated_at

---

### 3.2 EvidenceChunk

Evidence の分割単位。

Fields:

* id
* evidence_id
* seq
* text
* token_count
* embedding_ref
* created_at

---

### 3.3 KnowledgeCard

圧縮済み知識。

Fields:

* id
* scope
* kind
* title
* body
* confidence
* valid_from
* valid_to
* metadata_json
* created_at
* updated_at

---

### 3.4 Artifact

成果物メタデータ。

Fields:

* id
* kind
* title
* uri
* version
* content_hash
* metadata_json
* created_at
* updated_at

---

### 3.5 LineageEdge

由来関係。

Fields:

* id
* from_ref
* edge_type
* to_ref
* weight
* metadata_json
* created_at

---

## 4. Commands / Operations

## 4.1 Evidence Operations

### `evidence.create`

新しい Evidence を登録する。

#### Input

```json
{
  "store": "journal",
  "kind": "transcript",
  "title": "meeting-2026-03-07",
  "source_uri": "file://notes/meeting-2026-03-07.md",
  "source_hash": "sha256:...",
  "recorded_at": "2026-03-07T10:00:00Z",
  "metadata": {}
}
```

#### Output

```json
{
  "id": "01H...",
  "ref": "memx:evidence:01H..."
}
```

---

### `evidence.chunk.append`

Evidence に chunk を追加する。

#### Input

```json
{
  "evidence_id": "01H...",
  "seq": 0,
  "text": "raw text chunk",
  "token_count": 128
}
```

#### Output

```json
{
  "id": "01HCHUNK...",
  "ref": "memx:evidence_chunk:01HCHUNK..."
}
```

---

### `evidence.get`

Evidence 本体を取得する。

#### Input

```json
{
  "id": "01H..."
}
```

#### Output

```json
{
  "id": "01H...",
  "store": "journal",
  "kind": "transcript",
  "title": "meeting-2026-03-07",
  "source_uri": "file://notes/meeting-2026-03-07.md",
  "metadata": {}
}
```

---

### `evidence.chunks.list`

Evidence に属する chunk を取得する。

#### Input

```json
{
  "evidence_id": "01H..."
}
```

#### Output

```json
{
  "items": [
    {
      "id": "01HCHUNK...",
      "seq": 0,
      "text": "raw text chunk"
    }
  ]
}
```

---

## 4.2 Knowledge Operations

### `knowledge.create`

新しい KnowledgeCard を登録する。

#### Input

```json
{
  "scope": "project",
  "kind": "summary",
  "title": "image-api-summary",
  "body": "要約本文",
  "confidence": "medium",
  "valid_from": null,
  "valid_to": null,
  "metadata": {}
}
```

#### Output

```json
{
  "id": "01H...",
  "ref": "memx:knowledge:01H..."
}
```

---

### `knowledge.get`

KnowledgeCard を取得する。

#### Input

```json
{
  "id": "01H..."
}
```

#### Output

```json
{
  "id": "01H...",
  "scope": "project",
  "kind": "summary",
  "title": "image-api-summary",
  "body": "要約本文",
  "confidence": "medium"
}
```

---

### `knowledge.search`

KnowledgeCard を検索する。

#### Input

```json
{
  "query": "normalization failure pattern",
  "scope": "project",
  "kind": "failure_pattern",
  "limit": 10
}
```

#### Output

```json
{
  "items": [
    {
      "id": "01H...",
      "ref": "memx:knowledge:01H...",
      "title": "input normalization regression",
      "confidence": "high"
    }
  ]
}
```

---

## 4.3 Artifact Operations

### `artifact.register`

Artifact metadata を登録する。

#### Input

```json
{
  "kind": "spec",
  "title": "image-summary-api-spec",
  "uri": "repo://docs/spec.md",
  "version": "git:abc123",
  "content_hash": "sha256:...",
  "metadata": {}
}
```

#### Output

```json
{
  "id": "01H...",
  "ref": "memx:artifact:01H..."
}
```

---

### `artifact.get`

Artifact metadata を取得する。

#### Input

```json
{
  "id": "01H..."
}
```

#### Output

```json
{
  "id": "01H...",
  "kind": "spec",
  "title": "image-summary-api-spec",
  "uri": "repo://docs/spec.md",
  "version": "git:abc123"
}
```

---

## 4.4 Lineage Operations

### `lineage.link`

由来関係を登録する。

#### Input

```json
{
  "from_ref": "memx:evidence:01H...",
  "edge_type": "summarizes",
  "to_ref": "memx:knowledge:01H..."
}
```

#### Output

```json
{
  "id": "01H...",
  "ref": "memx:lineage:01H..."
}
```

---

### `lineage.trace`

ある entity から lineage を辿る。

#### Input

```json
{
  "ref": "memx:knowledge:01H...",
  "direction": "upstream",
  "depth": 3
}
```

#### Output

```json
{
  "items": [
    {
      "from_ref": "memx:evidence:01H...",
      "edge_type": "summarizes",
      "to_ref": "memx:knowledge:01H..."
    }
  ]
}
```

---

## 4.5 Summarize / Distill Operations

### `summarize.run`

Evidence から summary を生成する。

#### Input

```json
{
  "evidence_refs": [
    "memx:evidence:01H..."
  ],
  "target_scope": "project",
  "title": "meeting-summary-2026-03-07"
}
```

#### Output

```json
{
  "knowledge_ref": "memx:knowledge:01H...",
  "lineage_refs": [
    "memx:lineage:01H..."
  ]
}
```

---

### `distill.run`

Evidence / Knowledge から distilled knowledge を生成する。

#### Input

```json
{
  "source_refs": [
    "memx:evidence:01H...",
    "memx:knowledge:01H..."
  ],
  "kind": "failure_pattern",
  "scope": "project",
  "title": "normalization-failure-pattern"
}
```

#### Output

```json
{
  "knowledge_ref": "memx:knowledge:01H..."
}
```

---

## 4.6 Archive / GC Operations

### `archive.move`

entity を archive に移す。

#### Input

```json
{
  "ref": "memx:evidence:01H..."
}
```

#### Output

```json
{
  "ref": "memx:evidence:01H...",
  "store": "archive"
}
```

---

### `gc.run`

GC を実行する。

#### Input

```json
{
  "policy": "default"
}
```

#### Output

```json
{
  "deleted_count": 0,
  "archived_count": 3
}
```

---

## 5. Retrieval Contract

### 5.1 Retrieve by Ref

typed ref から entity を取得できなければならない。

### 5.2 Retrieve by Query

query / kind / scope / store により relevant な entity を検索できなければならない。

### 5.3 Summary-first Retrieval

上位層は通常、knowledge を優先参照し、必要時に evidence へ降りる。
memx-core はこの運用を支える検索インターフェースを提供する。

---

## 6. Error Model

### Error Codes

* `NOT_FOUND`
* `INVALID_INPUT`
* `INVALID_REF`
* `CONFLICT`
* `STORE_MISMATCH`
* `UNSUPPORTED_OPERATION`
* `INTERNAL_ERROR`

### Error Response Example

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "knowledge not found",
    "details": {
      "id": "01H..."
    }
  }
}
```

---

## 7. Non-Goals in Interface Layer

以下の操作は memx-core interface に含めない。

* task.create
* task.update_status
* decision.accept
* question.resolve
* issue.sync
* issue.comment.create
* execution.plan.build

---

## 8. Compatibility Rules

* append / additive change を優先する
* typed ref format は安定させる
* entity type 名は短期で変更しない
* field 削除ではなく deprecate を優先する

---

## 9. Minimal Initial Surface

初期実装で必須とする操作:

* evidence.create
* evidence.chunk.append
* evidence.get
* evidence.chunks.list
* knowledge.create
* knowledge.get
* knowledge.search
* artifact.register
* artifact.get
* lineage.link
* lineage.trace
* summarize.run
* distill.run
* archive.move
* gc.run