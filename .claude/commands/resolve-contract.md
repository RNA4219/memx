# Resolve contract information for a feature or task

Get acceptance criteria, forbidden patterns, definition of done, and dependencies for a feature or task.

## Usage

```
/resolve-contract [--feature FEATURE] [--task-id TASK_ID]
```

## Parameters

- `--feature`: Feature key to resolve contract for
- `--task-id`: Task ID to resolve contract for

## Example

```
/resolve-contract --feature memory-import
/resolve-contract --task-id task:feature:local:123
```

## Implementation

```bash
mem docs contract --feature "$ARGUMENTS" --json 2>/dev/null || \
mem docs contract --task-id "$ARGUMENTS" --json 2>/dev/null
```

## Output

Returns contract information with:
- `required`: Array of required documents
- `acceptance_criteria`: Array of acceptance criteria items
- `forbidden_patterns`: Array of forbidden patterns
- `definition_of_done`: Array of DoD items
- `dependencies`: Array of dependencies