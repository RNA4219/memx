# Resolve documents for a feature, task, or topic

Resolve required and recommended documents based on feature name, task ID, or topic.

## Usage

```
/resolve-docs [--feature FEATURE] [--task-id TASK_ID] [--topic TOPIC] [--limit N]
```

## Parameters

- `--feature`: Feature key to resolve documents for
- `--task-id`: Task ID to resolve documents for
- `--topic`: Topic query to search documents
- `--limit`: Maximum number of results (default: 10)

## Example

```
/resolve-docs --feature memory-import
/resolve-docs --task-id task:feature:local:123
/resolve-docs --topic "acceptance criteria"
```

## Implementation

```bash
mem docs resolve --feature "$ARGUMENTS" --json 2>/dev/null || \
mem docs resolve --task-id "$ARGUMENTS" --json 2>/dev/null || \
mem docs resolve --topic "$ARGUMENTS" --json 2>/dev/null
```

## Output

Returns a JSON object with:
- `required`: Array of required documents
- `recommended`: Array of recommended documents

Each entry contains:
- `doc_id`: Document ID
- `title`: Document title
- `version`: Document version
- `importance`: "required" or "recommended"
- `reason`: Why this document was selected
- `top_chunks`: Array of top chunk IDs to read