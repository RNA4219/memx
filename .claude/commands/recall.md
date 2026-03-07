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
cd C:/Users/ryo-n/Codex_dev/memx-core/memx_spec_v3/go && go run ./cmd/mem out search "<query>" --json
```

Also search journal and knowledge stores if relevant:

```bash
go run ./cmd/mem out journal search "<query>" --json
go run ./cmd/mem out knowledge search "<query>" --json
```