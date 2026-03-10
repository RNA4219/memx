# Read chunks from a document

Retrieve specific chunks from a document by doc ID, heading, or query.

## Usage

```
/read-chunks --doc-id DOC_ID [--heading HEADING] [--query QUERY] [--limit N]
```

## Parameters

- `--doc-id`: Document ID (required)
- `--heading`: Filter by heading name
- `--query`: Filter by content query
- `--limit`: Maximum number of chunks to return

## Example

```
/read-chunks --doc-id doc:spec:memory-import
/read-chunks --doc-id doc:spec:memory-import --heading "Acceptance Criteria"
/read-chunks --doc-id doc:spec:memory-import --query "validation"
```

## Implementation

```bash
mem docs chunks --doc-id "$ARGUMENTS" --json
```

## Output

Returns chunks with:
- `chunk_id`: Unique chunk identifier
- `heading`: Section heading
- `heading_path`: Full path of headings
- `body`: Chunk content
- `importance`: "required", "recommended", or "reference"