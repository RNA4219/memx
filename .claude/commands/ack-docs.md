# Acknowledge reading a document

Record that a document has been read for a specific task.

## Usage

```
/ack-docs --task-id TASK_ID --doc-id DOC_ID [--version VERSION]
```

## Parameters

- `--task-id`: Task ID (required)
- `--doc-id`: Document ID (required)
- `--version`: Document version (optional, uses current if not specified)

## Example

```
/ack-docs --task-id task:feature:local:123 --doc-id doc:spec:memory-import
```

## Implementation

```bash
mem docs ack --task-id "$ARGUMENTS" --json
```

## Output

Returns a read receipt with:
- `task_id`: Task ID
- `doc_id`: Document ID
- `version`: Version that was read
- `read_at`: Timestamp of acknowledgment
- `reader`: Who read the document