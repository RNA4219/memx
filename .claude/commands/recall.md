# Recall Information

Search and retrieve information from memx memory stores.

## Usage

```
/recall <search query>
```

## What it does

1. Searches across short, journal, and knowledge stores
2. Returns matching notes with their IDs and titles
3. Use `/show <id>` to view full content

## Example

```
/recall authentication API
```

## Implementation

Use the Bash tool to run:

```bash
cd C:/Users/ryo-n/Codex_dev/memx-core/memx_spec_v3/go && go run ./cmd/mem out search --json "<query>"
```

The top-level search already aggregates short, journal, and knowledge stores. Use store-specific search only when you need to narrow the scope:

```bash
go run ./cmd/mem out journal search --json "<query>"
go run ./cmd/mem out knowledge search --json "<query>"
```
