# Check if task documents are stale

Check if any documents associated with a task have been updated since they were last read.

## Usage

```
/stale-check --task-id TASK_ID
```

## Parameters

- `--task-id`: Task ID to check (required)

## Example

```
/stale-check --task-id task:feature:local:123
```

## Implementation

```bash
mem docs stale --task-id "$ARGUMENTS" --json
```

## Output

Returns stale check result with:
- `task_id`: Task ID
- `stale`: Array of stale reasons (empty if fresh)

Each stale reason contains:
- `doc_id`: Document ID
- `previous_version`: Version that was read
- `current_version`: Current version
- `reason`: "version_mismatch" or "document_missing"