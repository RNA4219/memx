# memx Help

Display help information for memx CLI commands.

## Usage

```
/memx-help
```

## What it does

Shows all available memx commands and their usage.

## Available Skills

| Skill | Description |
|-------|-------------|
| `/remember` | Store information in short-term memory |
| `/recall` | Search and retrieve information |
| `/journal` | Add a journal entry |
| `/knowledge` | Add a knowledge base entry |
| `/show` | Display a specific note by ID |
| `/memx-help` | Show this help |

## Store Overview

| Store | Purpose | typed_ref type |
|-------|---------|----------------|
| short | Temporary notes, logs | evidence |
| journal | Time-series logs (requires scope) | evidence |
| knowledge | Knowledge base (requires scope) | knowledge |
| archive | Archived notes (no search) | evidence |

## CLI Direct Access

```bash
cd C:/Users/ryo-n/Codex_dev/memx-core/memx_spec_v3/go
go run ./cmd/mem in short --title "Title" --body "Body"
go run ./cmd/mem out search --json "query"
go run ./cmd/mem out show <id>
```

## API Server

```bash
go run ./cmd/mem api serve --addr 127.0.0.1:7766
```

Then use `--api-url http://127.0.0.1:7766` for CLI commands.
