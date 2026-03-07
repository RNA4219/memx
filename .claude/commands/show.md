# Show Note

Display a specific note by its ID.

## Usage

```
/show <note_id>
```

## What it does

1. Retrieves and displays the full content of a note
2. Updates access count and timestamp
3. Works for notes from any store (short, journal, knowledge, archive)

## Example

```
/show abc123def456...
```

## Implementation

Use the Bash tool to run:

```bash
cd C:/Users/ryo-n/Codex_dev/memx-core/memx_spec_v3/go && go run ./cmd/mem out show "<note_id>"
```

If the note might be in journal or knowledge:

```bash
go run ./cmd/mem out journal show "<note_id>"
go run ./cmd/mem out knowledge show "<note_id>"
```