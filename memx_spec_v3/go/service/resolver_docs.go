package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

type ResolverDocument struct {
	DocID              string   `json:"doc_id"`
	DocType            string   `json:"doc_type"`
	Title              string   `json:"title"`
	SourcePath         string   `json:"source_path"`
	Version            string   `json:"version"`
	UpdatedAt          string   `json:"updated_at"`
	Summary            string   `json:"summary"`
	Body               string   `json:"body,omitempty"`
	Tags               []string `json:"tags,omitempty"`
	FeatureKeys        []string `json:"feature_keys,omitempty"`
	TaskIDs            []string `json:"task_ids,omitempty"`
	AcceptanceCriteria []string `json:"acceptance_criteria,omitempty"`
	ForbiddenPatterns  []string `json:"forbidden_patterns,omitempty"`
	DefinitionOfDone   []string `json:"definition_of_done,omitempty"`
	Dependencies       []string `json:"dependencies,omitempty"`
	Importance         string   `json:"importance"`
}

type ResolverChunk struct {
	ChunkID       string   `json:"chunk_id"`
	DocID         string   `json:"doc_id"`
	Heading       string   `json:"heading,omitempty"`
	HeadingPath   []string `json:"heading_path"`
	Ordinal       int      `json:"ordinal"`
	Body          string   `json:"body"`
	TokenEstimate int      `json:"token_estimate"`
	Importance    string   `json:"importance"`
}

type ResolveEntry struct {
	DocID      string   `json:"doc_id"`
	Title      string   `json:"title"`
	Version    string   `json:"version"`
	Importance string   `json:"importance"`
	Reason     string   `json:"reason"`
	TopChunks  []string `json:"top_chunks"`
}

type ResolverReadReceipt struct {
	TaskID   string   `json:"task_id"`
	DocID    string   `json:"doc_id"`
	Version  string   `json:"version"`
	ChunkIDs []string `json:"chunk_ids,omitempty"`
	Reader   string   `json:"reader"`
	ReadAt   string   `json:"read_at"`
}

type ResolverStaleReason struct {
	TaskID          string `json:"task_id"`
	DocID           string `json:"doc_id"`
	PreviousVersion string `json:"previous_version"`
	CurrentVersion  string `json:"current_version"`
	Reason          string `json:"reason"`
	DetectedAt      string `json:"detected_at"`
}

type ChunkingOptions struct {
	Mode     string `json:"mode,omitempty"`
	MaxChars int    `json:"max_chars,omitempty"`
}

type DocsIngestRequest struct {
	DocID              string          `json:"doc_id,omitempty"`
	DocType            string          `json:"doc_type"`
	Title              string          `json:"title"`
	SourcePath         string          `json:"source_path,omitempty"`
	Version            string          `json:"version"`
	UpdatedAt          string          `json:"updated_at,omitempty"`
	Tags               []string        `json:"tags,omitempty"`
	FeatureKeys        []string        `json:"feature_keys,omitempty"`
	TaskIDs            []string        `json:"task_ids,omitempty"`
	Summary            string          `json:"summary,omitempty"`
	Body               string          `json:"body"`
	Chunking           ChunkingOptions `json:"chunking,omitempty"`
	AcceptanceCriteria []string        `json:"acceptance_criteria,omitempty"`
	ForbiddenPatterns  []string        `json:"forbidden_patterns,omitempty"`
	DefinitionOfDone   []string        `json:"definition_of_done,omitempty"`
	Dependencies       []string        `json:"dependencies,omitempty"`
}

type DocsResolveRequest struct {
	Feature string `json:"feature,omitempty"`
	TaskID  string `json:"task_id,omitempty"`
	Topic   string `json:"topic,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

type ChunksGetRequest struct {
	DocID    string   `json:"doc_id,omitempty"`
	Query    string   `json:"query,omitempty"`
	Heading  string   `json:"heading,omitempty"`
	Limit    int      `json:"limit,omitempty"`
	ChunkIDs []string `json:"chunk_ids,omitempty"`
}

type DocsSearchRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

type ReadsAckRequest struct {
	TaskID   string   `json:"task_id"`
	DocID    string   `json:"doc_id"`
	Version  string   `json:"version,omitempty"`
	ChunkIDs []string `json:"chunk_ids,omitempty"`
	Reader   string   `json:"reader,omitempty"`
}

type DocsStaleCheckRequest struct {
	TaskID string `json:"task_id"`
}

type ContractsResolveRequest struct {
	Feature string `json:"feature,omitempty"`
	TaskID  string `json:"task_id,omitempty"`
}

type resolverSection struct {
	Heading     string
	HeadingPath []string
	Body        string
	Importance  string
}

type scoredResolverDoc struct {
	Doc   ResolverDocument
	Score int
	Why   string
}

var markdownHeadingPattern = regexp.MustCompile(`^(#{1,6})\s+(.*\S)\s*$`)

func (r *DocsIngestRequest) normalize() {
	r.DocType = strings.TrimSpace(strings.ToLower(r.DocType))
	r.Title = strings.TrimSpace(r.Title)
	r.SourcePath = strings.TrimSpace(r.SourcePath)
	r.Version = strings.TrimSpace(r.Version)
	r.UpdatedAt = strings.TrimSpace(r.UpdatedAt)
	r.Summary = strings.TrimSpace(r.Summary)
	r.Body = strings.TrimSpace(r.Body)
	r.DocID = strings.TrimSpace(r.DocID)
	r.Chunking.Mode = strings.TrimSpace(strings.ToLower(r.Chunking.Mode))
	if r.Chunking.Mode == "" {
		r.Chunking.Mode = "heading"
	}
	if r.Chunking.MaxChars <= 0 {
		r.Chunking.MaxChars = 4000
	}
	if r.UpdatedAt == "" {
		r.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	r.Tags = uniqueTrimmedStrings(r.Tags)
	r.FeatureKeys = uniqueTrimmedStrings(r.FeatureKeys)
	r.TaskIDs = uniqueTrimmedStrings(r.TaskIDs)
	r.AcceptanceCriteria = uniqueTrimmedStrings(r.AcceptanceCriteria)
	r.ForbiddenPatterns = uniqueTrimmedStrings(r.ForbiddenPatterns)
	r.DefinitionOfDone = uniqueTrimmedStrings(r.DefinitionOfDone)
	r.Dependencies = uniqueTrimmedStrings(r.Dependencies)
	if r.DocID == "" {
		r.DocID = buildDocID(r.DocType, r.Title)
	}
}

func (r *DocsIngestRequest) validate() error {
	if r.DocType == "" || r.Title == "" || r.Version == "" || r.Body == "" {
		return fmt.Errorf("%w: doc_type/title/version/body is required", ErrInvalidArgument)
	}
	switch r.Chunking.Mode {
	case "heading", "fixed":
	default:
		return fmt.Errorf("%w: unsupported chunking.mode", ErrInvalidArgument)
	}
	return nil
}

func (s *Service) DocsIngest(ctx context.Context, req DocsIngestRequest) (ResolverDocument, int, error) {
	req.normalize()
	if err := req.validate(); err != nil {
		return ResolverDocument{}, 0, err
	}
	sections := parseSectionsForChunking(req.Title, req.Body, req.DocType, req.Chunking.Mode)
	if len(req.AcceptanceCriteria) == 0 {
		req.AcceptanceCriteria = extractSectionItems(sections, []string{"acceptance criteria", "acceptance"})
	}
	if len(req.ForbiddenPatterns) == 0 {
		req.ForbiddenPatterns = extractSectionItems(sections, []string{"forbidden patterns", "forbidden", "anti-pattern"})
	}
	if len(req.DefinitionOfDone) == 0 {
		req.DefinitionOfDone = extractSectionItems(sections, []string{"definition of done", "done"})
	}
	if len(req.Dependencies) == 0 {
		req.Dependencies = extractSectionItems(sections, []string{"dependency", "dependencies"})
	}
	if req.Summary == "" {
		req.Summary = buildDocumentSummary(req.Body)
	}
	if existing, err := s.getResolverDocument(ctx, req.DocID); err == nil {
		if compareVersions(req.Version, existing.Version) < 0 {
			return ResolverDocument{}, 0, fmt.Errorf("%w: version conflict", ErrConflict)
		}
	} else if err != ErrNotFound {
		return ResolverDocument{}, 0, err
	}
	doc := ResolverDocument{
		DocID:              req.DocID,
		DocType:            req.DocType,
		Title:              req.Title,
		SourcePath:         req.SourcePath,
		Version:            req.Version,
		UpdatedAt:          req.UpdatedAt,
		Summary:            req.Summary,
		Body:               req.Body,
		Tags:               req.Tags,
		FeatureKeys:        req.FeatureKeys,
		TaskIDs:            req.TaskIDs,
		AcceptanceCriteria: req.AcceptanceCriteria,
		ForbiddenPatterns:  req.ForbiddenPatterns,
		DefinitionOfDone:   req.DefinitionOfDone,
		Dependencies:       req.Dependencies,
		Importance:         defaultDocImportance(req.DocType),
	}
	chunks := buildResolverChunks(req.DocID, sections, req.Chunking)
	tx, err := s.Conn.DB.BeginTx(ctx, nil)
	if err != nil {
		return ResolverDocument{}, 0, err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `DELETE FROM resolver_document_links WHERE src_doc_id = ?;`, doc.DocID); err != nil {
		return ResolverDocument{}, 0, err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM resolver_chunks WHERE doc_id = ?;`, doc.DocID); err != nil {
		return ResolverDocument{}, 0, err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM resolver_documents WHERE doc_id = ?;`, doc.DocID); err != nil {
		return ResolverDocument{}, 0, err
	}
	if _, err := tx.ExecContext(ctx, `
INSERT INTO resolver_documents(
  doc_id, doc_type, title, source_path, version, updated_at,
  summary, body, tags_json, feature_keys_json, task_ids_json,
  acceptance_criteria_json, forbidden_patterns_json, definition_of_done_json,
  dependencies_json, importance
) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
`, doc.DocID, doc.DocType, doc.Title, doc.SourcePath, doc.Version, doc.UpdatedAt,
		doc.Summary, doc.Body, mustJSON(doc.Tags), mustJSON(doc.FeatureKeys), mustJSON(doc.TaskIDs),
		mustJSON(doc.AcceptanceCriteria), mustJSON(doc.ForbiddenPatterns), mustJSON(doc.DefinitionOfDone),
		mustJSON(doc.Dependencies), doc.Importance); err != nil {
		return ResolverDocument{}, 0, err
	}
	for _, chunk := range chunks {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO resolver_chunks(
  chunk_id, doc_id, heading, heading_path_json, ordinal, body, token_estimate, importance
) VALUES(?, ?, ?, ?, ?, ?, ?, ?);
`, chunk.ChunkID, chunk.DocID, chunk.Heading, mustJSON(chunk.HeadingPath), chunk.Ordinal, chunk.Body, chunk.TokenEstimate, chunk.Importance); err != nil {
			return ResolverDocument{}, 0, err
		}
	}
	for _, dep := range doc.Dependencies {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO resolver_document_links(src_doc_id, dst_doc_id, link_type) VALUES(?, ?, 'depends_on');
`, doc.DocID, dep); err != nil {
			return ResolverDocument{}, 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return ResolverDocument{}, 0, err
	}
	return doc, len(chunks), nil
}

func (s *Service) DocsResolve(ctx context.Context, req DocsResolveRequest) ([]ResolveEntry, []ResolveEntry, error) {
	req.Feature = strings.TrimSpace(req.Feature)
	req.TaskID = strings.TrimSpace(req.TaskID)
	req.Topic = strings.TrimSpace(req.Topic)
	if req.Feature == "" && req.TaskID == "" && req.Topic == "" {
		return nil, nil, fmt.Errorf("%w: feature/task_id/topic is required", ErrInvalidArgument)
	}
	docs, err := s.listResolverDocuments(ctx)
	if err != nil {
		return nil, nil, err
	}
	scored := scoreResolverDocuments(docs, req.Feature, req.TaskID, req.Topic)
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	required := make([]ResolveEntry, 0, limit)
	recommended := make([]ResolveEntry, 0, limit)
	count := 0
	query := firstNonEmpty(req.Topic, req.Feature, req.TaskID)
	for _, item := range scored {
		if count >= limit {
			break
		}
		entry, err := s.buildResolveEntry(ctx, item.Doc, item.Why, query)
		if err != nil {
			return nil, nil, err
		}
		if item.Doc.Importance == "required" {
			required = append(required, entry)
		} else {
			recommended = append(recommended, entry)
		}
		count++
	}
	return required, recommended, nil
}

func (s *Service) DocsSearch(ctx context.Context, req DocsSearchRequest) ([]ResolveEntry, error) {
	req.Query = strings.TrimSpace(req.Query)
	if req.Query == "" {
		return nil, fmt.Errorf("%w: query is required", ErrInvalidArgument)
	}
	required, recommended, err := s.DocsResolve(ctx, DocsResolveRequest{Topic: req.Query, Limit: req.Limit})
	if err != nil {
		return nil, err
	}
	return append(required, recommended...), nil
}

func (s *Service) ChunksGet(ctx context.Context, req ChunksGetRequest) (string, []ResolverChunk, error) {
	if len(req.ChunkIDs) > 0 {
		chunks, err := s.getResolverChunksByIDs(ctx, req.ChunkIDs)
		if err != nil {
			return "", nil, err
		}
		docID := ""
		if len(chunks) > 0 {
			docID = chunks[0].DocID
		}
		return docID, chunks, nil
	}
	req.DocID = strings.TrimSpace(req.DocID)
	if req.DocID == "" {
		return "", nil, fmt.Errorf("%w: doc_id or chunk_ids is required", ErrInvalidArgument)
	}
	chunks, err := s.getResolverChunksByDocID(ctx, req.DocID)
	if err != nil {
		return "", nil, err
	}
	return req.DocID, filterResolverChunks(chunks, strings.TrimSpace(req.Heading), strings.TrimSpace(req.Query), req.Limit), nil
}

func (s *Service) ReadsAck(ctx context.Context, req ReadsAckRequest) (ResolverReadReceipt, error) {
	req.TaskID = strings.TrimSpace(req.TaskID)
	req.DocID = strings.TrimSpace(req.DocID)
	req.Version = strings.TrimSpace(req.Version)
	req.Reader = strings.TrimSpace(req.Reader)
	if req.TaskID == "" || req.DocID == "" {
		return ResolverReadReceipt{}, fmt.Errorf("%w: task_id/doc_id is required", ErrInvalidArgument)
	}
	doc, err := s.getResolverDocument(ctx, req.DocID)
	if err != nil {
		return ResolverReadReceipt{}, err
	}
	if req.Version == "" {
		req.Version = doc.Version
	}
	if req.Reader == "" {
		req.Reader = "agent"
	}
	receipt := ResolverReadReceipt{
		TaskID:   req.TaskID,
		DocID:    req.DocID,
		Version:  req.Version,
		ChunkIDs: uniqueTrimmedStrings(req.ChunkIDs),
		Reader:   req.Reader,
		ReadAt:   time.Now().UTC().Format(time.RFC3339),
	}
	if _, err := s.Conn.DB.ExecContext(ctx, `
INSERT INTO resolver_read_receipts(task_id, doc_id, version, chunk_ids_json, reader, read_at)
VALUES(?, ?, ?, ?, ?, ?);
`, receipt.TaskID, receipt.DocID, receipt.Version, mustJSON(receipt.ChunkIDs), receipt.Reader, receipt.ReadAt); err != nil {
		return ResolverReadReceipt{}, err
	}
	return receipt, nil
}

func (s *Service) DocsStaleCheck(ctx context.Context, req DocsStaleCheckRequest) ([]ResolverStaleReason, error) {
	req.TaskID = strings.TrimSpace(req.TaskID)
	if req.TaskID == "" {
		return nil, fmt.Errorf("%w: task_id is required", ErrInvalidArgument)
	}
	rows, err := s.Conn.DB.QueryContext(ctx, `
SELECT rr.task_id, rr.doc_id, rr.version
FROM resolver_read_receipts rr
JOIN (
  SELECT doc_id, MAX(id) AS latest_id
  FROM resolver_read_receipts
  WHERE task_id = ?
  GROUP BY doc_id
) latest ON latest.latest_id = rr.id
WHERE rr.task_id = ?
ORDER BY rr.id ASC;
`, req.TaskID, req.TaskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stale []ResolverStaleReason
	for rows.Next() {
		var taskID, docID, version string
		if err := rows.Scan(&taskID, &docID, &version); err != nil {
			return nil, err
		}
		doc, err := s.getResolverDocument(ctx, docID)
		if err != nil {
			if err == ErrNotFound {
				stale = append(stale, ResolverStaleReason{TaskID: taskID, DocID: docID, PreviousVersion: version, Reason: "document_missing", DetectedAt: time.Now().UTC().Format(time.RFC3339)})
				continue
			}
			return nil, err
		}
		if doc.Version != version {
			stale = append(stale, ResolverStaleReason{TaskID: taskID, DocID: docID, PreviousVersion: version, CurrentVersion: doc.Version, Reason: "version_mismatch", DetectedAt: time.Now().UTC().Format(time.RFC3339)})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stale, nil
}

func (s *Service) ContractsResolve(ctx context.Context, req ContractsResolveRequest) ([]ResolveEntry, []string, []string, []string, []string, error) {
	req.Feature = strings.TrimSpace(req.Feature)
	req.TaskID = strings.TrimSpace(req.TaskID)
	if req.Feature == "" && req.TaskID == "" {
		return nil, nil, nil, nil, nil, fmt.Errorf("%w: feature or task_id is required", ErrInvalidArgument)
	}
	required, recommended, err := s.DocsResolve(ctx, DocsResolveRequest{Feature: req.Feature, TaskID: req.TaskID, Limit: 20})
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	all := append(append([]ResolveEntry{}, required...), recommended...)
	var acceptance, forbidden, done, dependencies []string
	for _, entry := range all {
		doc, err := s.getResolverDocument(ctx, entry.DocID)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}
		acceptance = append(acceptance, doc.AcceptanceCriteria...)
		forbidden = append(forbidden, doc.ForbiddenPatterns...)
		done = append(done, doc.DefinitionOfDone...)
		dependencies = append(dependencies, doc.Dependencies...)
	}
	return required, uniqueTrimmedStrings(acceptance), uniqueTrimmedStrings(forbidden), uniqueTrimmedStrings(done), uniqueTrimmedStrings(dependencies), nil
}

func (s *Service) getResolverDocument(ctx context.Context, docID string) (ResolverDocument, error) {
	var doc ResolverDocument
	var tagsJSON, featureJSON, taskJSON, acceptanceJSON, forbiddenJSON, doneJSON, dependenciesJSON string
	err := s.Conn.DB.QueryRowContext(ctx, `
SELECT doc_id, doc_type, title, source_path, version, updated_at, summary, body,
       tags_json, feature_keys_json, task_ids_json,
       acceptance_criteria_json, forbidden_patterns_json, definition_of_done_json,
       dependencies_json, importance
FROM resolver_documents
WHERE doc_id = ?;
`, docID).Scan(
		&doc.DocID, &doc.DocType, &doc.Title, &doc.SourcePath, &doc.Version, &doc.UpdatedAt, &doc.Summary, &doc.Body,
		&tagsJSON, &featureJSON, &taskJSON, &acceptanceJSON, &forbiddenJSON, &doneJSON, &dependenciesJSON, &doc.Importance,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ResolverDocument{}, ErrNotFound
		}
		return ResolverDocument{}, err
	}
	doc.Tags = decodeStringSlice(tagsJSON)
	doc.FeatureKeys = decodeStringSlice(featureJSON)
	doc.TaskIDs = decodeStringSlice(taskJSON)
	doc.AcceptanceCriteria = decodeStringSlice(acceptanceJSON)
	doc.ForbiddenPatterns = decodeStringSlice(forbiddenJSON)
	doc.DefinitionOfDone = decodeStringSlice(doneJSON)
	doc.Dependencies = decodeStringSlice(dependenciesJSON)
	return doc, nil
}
func (s *Service) listResolverDocuments(ctx context.Context) ([]ResolverDocument, error) {
	rows, err := s.Conn.DB.QueryContext(ctx, `
SELECT doc_id, doc_type, title, source_path, version, updated_at, summary, body,
       tags_json, feature_keys_json, task_ids_json,
       acceptance_criteria_json, forbidden_patterns_json, definition_of_done_json,
       dependencies_json, importance
FROM resolver_documents;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ResolverDocument
	for rows.Next() {
		var doc ResolverDocument
		var tagsJSON, featureJSON, taskJSON, acceptanceJSON, forbiddenJSON, doneJSON, dependenciesJSON string
		if err := rows.Scan(&doc.DocID, &doc.DocType, &doc.Title, &doc.SourcePath, &doc.Version, &doc.UpdatedAt, &doc.Summary, &doc.Body, &tagsJSON, &featureJSON, &taskJSON, &acceptanceJSON, &forbiddenJSON, &doneJSON, &dependenciesJSON, &doc.Importance); err != nil {
			return nil, err
		}
		doc.Tags = decodeStringSlice(tagsJSON)
		doc.FeatureKeys = decodeStringSlice(featureJSON)
		doc.TaskIDs = decodeStringSlice(taskJSON)
		doc.AcceptanceCriteria = decodeStringSlice(acceptanceJSON)
		doc.ForbiddenPatterns = decodeStringSlice(forbiddenJSON)
		doc.DefinitionOfDone = decodeStringSlice(doneJSON)
		doc.Dependencies = decodeStringSlice(dependenciesJSON)
		out = append(out, doc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) getResolverChunksByDocID(ctx context.Context, docID string) ([]ResolverChunk, error) {
	rows, err := s.Conn.DB.QueryContext(ctx, `
SELECT chunk_id, doc_id, heading, heading_path_json, ordinal, body, token_estimate, importance
FROM resolver_chunks
WHERE doc_id = ?
ORDER BY ordinal ASC;
`, docID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ResolverChunk
	for rows.Next() {
		var chunk ResolverChunk
		var headingPathJSON string
		if err := rows.Scan(&chunk.ChunkID, &chunk.DocID, &chunk.Heading, &headingPathJSON, &chunk.Ordinal, &chunk.Body, &chunk.TokenEstimate, &chunk.Importance); err != nil {
			return nil, err
		}
		chunk.HeadingPath = decodeStringSlice(headingPathJSON)
		out = append(out, chunk)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(out) == 0 {
		if _, err := s.getResolverDocument(ctx, docID); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (s *Service) getResolverChunksByIDs(ctx context.Context, chunkIDs []string) ([]ResolverChunk, error) {
	var out []ResolverChunk
	for _, chunkID := range uniqueTrimmedStrings(chunkIDs) {
		var chunk ResolverChunk
		var headingPathJSON string
		err := s.Conn.DB.QueryRowContext(ctx, `
SELECT chunk_id, doc_id, heading, heading_path_json, ordinal, body, token_estimate, importance
FROM resolver_chunks
WHERE chunk_id = ?;
`, chunkID).Scan(&chunk.ChunkID, &chunk.DocID, &chunk.Heading, &headingPathJSON, &chunk.Ordinal, &chunk.Body, &chunk.TokenEstimate, &chunk.Importance)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrNotFound
			}
			return nil, err
		}
		chunk.HeadingPath = decodeStringSlice(headingPathJSON)
		out = append(out, chunk)
	}
	return out, nil
}

func (s *Service) buildResolveEntry(ctx context.Context, doc ResolverDocument, reason string, query string) (ResolveEntry, error) {
	chunks, err := s.getResolverChunksByDocID(ctx, doc.DocID)
	if err != nil {
		return ResolveEntry{}, err
	}
	return ResolveEntry{DocID: doc.DocID, Title: doc.Title, Version: doc.Version, Importance: doc.Importance, Reason: reason, TopChunks: pickTopChunkIDs(chunks, query, 3)}, nil
}

func buildResolverChunks(docID string, sections []resolverSection, opts ChunkingOptions) []ResolverChunk {
	chunks := make([]ResolverChunk, 0, len(sections))
	ordinal := 1
	for _, section := range sections {
		for _, part := range splitChunkBody(section.Body, opts.MaxChars) {
			body := strings.TrimSpace(part)
			if body == "" {
				continue
			}
			chunks = append(chunks, ResolverChunk{ChunkID: fmt.Sprintf("chunk:%s:%03d", docID, ordinal), DocID: docID, Heading: section.Heading, HeadingPath: append([]string(nil), section.HeadingPath...), Ordinal: ordinal, Body: body, TokenEstimate: estimateTokens(body), Importance: section.Importance})
			ordinal++
		}
	}
	return chunks
}

func parseSectionsForChunking(title string, body string, docType string, mode string) []resolverSection {
	if mode == "fixed" {
		return []resolverSection{{
			Heading:     strings.TrimSpace(title),
			HeadingPath: []string{strings.TrimSpace(title)},
			Body:        strings.TrimSpace(body),
			Importance:  defaultDocImportance(docType),
		}}
	}
	return parseResolverSections(title, body, docType)
}
func parseResolverSections(title string, body string, docType string) []resolverSection {
	defaultImportance := defaultDocImportance(docType)
	lines := strings.Split(body, "\n")
	path := []string{strings.TrimSpace(title)}
	currentHeading := strings.TrimSpace(title)
	var currentBody strings.Builder
	sections := make([]resolverSection, 0)
	flush := func() {
		text := strings.TrimSpace(currentBody.String())
		if text == "" {
			return
		}
		sections = append(sections, resolverSection{Heading: currentHeading, HeadingPath: append([]string(nil), path...), Body: text, Importance: inferChunkImportance(defaultImportance, currentHeading, path)})
		currentBody.Reset()
	}
	for _, line := range lines {
		if matches := markdownHeadingPattern.FindStringSubmatch(line); matches != nil {
			flush()
			level := len(matches[1])
			heading := strings.TrimSpace(matches[2])
			if level <= 1 {
				path = []string{heading}
			} else {
				base := []string{title}
				if len(path) > 0 {
					base = append([]string(nil), path...)
				}
				if len(base) >= level {
					base = append([]string(nil), base[:level-1]...)
				}
				base = append(base, heading)
				path = base
			}
			currentHeading = heading
			continue
		}
		currentBody.WriteString(line)
		currentBody.WriteByte('\n')
	}
	flush()
	if len(sections) == 0 {
		sections = append(sections, resolverSection{Heading: title, HeadingPath: []string{title}, Body: strings.TrimSpace(body), Importance: defaultImportance})
	}
	return sections
}

func extractSectionItems(sections []resolverSection, keywords []string) []string {
	var out []string
	for _, section := range sections {
		// Match only the heading name itself, not the full path
		// This prevents false positives when document title contains a keyword
		headingLower := strings.ToLower(section.Heading)
		matched := false
		for _, keyword := range keywords {
			if strings.Contains(headingLower, keyword) {
				matched = true
				break
			}
		}
		if matched {
			out = append(out, parseListItems(section.Body)...)
		}
	}
	return uniqueTrimmedStrings(out)
}

func parseListItems(body string) []string {
	var out []string
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "- "):
			trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
		case strings.HasPrefix(trimmed, "* "):
			trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "* "))
		default:
			trimmed = trimNumericListPrefix(trimmed)
		}
		if trimmed != "" && trimmed != strings.TrimSpace(line) {
			out = append(out, trimmed)
		}
	}
	if len(out) > 0 {
		return uniqueTrimmedStrings(out)
	}
	text := strings.TrimSpace(body)
	if text == "" {
		return nil
	}
	return []string{text}
}

func trimNumericListPrefix(s string) string {
	idx := strings.Index(s, ". ")
	if idx <= 0 {
		return s
	}
	for _, r := range s[:idx] {
		if r < '0' || r > '9' {
			return s
		}
	}
	return strings.TrimSpace(s[idx+2:])
}
func splitChunkBody(body string, maxChars int) []string {
	text := strings.TrimSpace(body)
	if text == "" {
		return nil
	}
	if maxChars <= 0 || len([]rune(text)) <= maxChars {
		return []string{text}
	}
	paragraphs := strings.Split(text, "\n\n")
	parts := make([]string, 0)
	var current strings.Builder
	flush := func() {
		chunk := strings.TrimSpace(current.String())
		if chunk != "" {
			parts = append(parts, chunk)
		}
		current.Reset()
	}
	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}
		candidate := paragraph
		if current.Len() > 0 {
			candidate = current.String() + "\n\n" + paragraph
		}
		if len([]rune(candidate)) <= maxChars {
			if current.Len() > 0 {
				current.WriteString("\n\n")
			}
			current.WriteString(paragraph)
			continue
		}
		if current.Len() > 0 {
			flush()
		}
		for _, hard := range splitRunes(paragraph, maxChars) {
			parts = append(parts, hard)
		}
	}
	flush()
	return parts
}

func splitRunes(text string, maxChars int) []string {
	runes := []rune(text)
	parts := make([]string, 0)
	for start := 0; start < len(runes); start += maxChars {
		end := start + maxChars
		if end > len(runes) {
			end = len(runes)
		}
		parts = append(parts, strings.TrimSpace(string(runes[start:end])))
	}
	return parts
}

func inferChunkImportance(defaultImportance string, heading string, path []string) string {
	joined := strings.ToLower(heading + " " + strings.Join(path, " "))
	switch {
	case strings.Contains(joined, "acceptance"), strings.Contains(joined, "forbidden"), strings.Contains(joined, "definition of done"):
		return "required"
	case strings.Contains(joined, "overview"), strings.Contains(joined, "background"):
		if defaultImportance == "required" {
			return "recommended"
		}
		return defaultImportance
	default:
		return defaultImportance
	}
}

func defaultDocImportance(docType string) string {
	switch strings.ToLower(docType) {
	case "spec", "requirement", "requirements", "interfaces", "design":
		return "required"
	case "cookbook", "runbook", "birdseye", "task":
		return "recommended"
	default:
		return "reference"
	}
}

func buildDocumentSummary(body string) string {
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return ""
	}
	runes := []rune(trimmed)
	if len(runes) <= 120 {
		return trimmed
	}
	return strings.TrimSpace(string(runes[:120]))
}

func buildDocID(docType string, title string) string {
	slug := slugify(title)
	if slug == "" {
		slug = fmt.Sprintf("doc-%d", time.Now().Unix())
	}
	return fmt.Sprintf("doc:%s:%s", docType, slug)
}

func slugify(s string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastDash = false
		case r == '-' || r == '_' || unicode.IsSpace(r):
			if !lastDash && b.Len() > 0 {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

func scoreResolverDocuments(docs []ResolverDocument, feature string, taskID string, topic string) []scoredResolverDoc {
	scored := make([]scoredResolverDoc, 0)
	for _, doc := range docs {
		score := 0
		reason := ""
		if taskID != "" && containsFold(doc.TaskIDs, taskID) {
			score += 1000
			reason = "task dependency matched"
		}
		if feature != "" && containsFold(doc.FeatureKeys, feature) {
			score += 900
			if reason == "" {
				reason = "feature key matched"
			}
		}
		query := firstNonEmpty(topic, feature)
		if query != "" {
			if containsFold(doc.Tags, query) {
				score += 300
				if reason == "" {
					reason = "tag matched"
				}
			}
			if textContainsFold(doc.Title, query) || textContainsFold(doc.Summary, query) || textContainsFold(doc.Body, query) || textContainsFold(doc.SourcePath, query) {
				score += 200
				if reason == "" {
					reason = "topic matched"
				}
			}
		}
		if score == 0 {
			continue
		}
		if reason == "" {
			reason = "relevance matched"
		}
		scored = append(scored, scoredResolverDoc{Doc: doc, Score: score, Why: reason})
	}
	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].Score == scored[j].Score {
			return scored[i].Doc.UpdatedAt > scored[j].Doc.UpdatedAt
		}
		return scored[i].Score > scored[j].Score
	})
	return scored
}

func filterResolverChunks(chunks []ResolverChunk, heading string, query string, limit int) []ResolverChunk {
	filtered := make([]ResolverChunk, 0, len(chunks))
	for _, chunk := range chunks {
		// When heading is specified, match only the chunk's heading name, not the full path
		// This prevents false positives when document title contains the heading keyword
		if heading != "" && !textContainsFold(chunk.Heading, heading) {
			continue
		}
		if query != "" && !textContainsFold(chunk.Body, query) && !textContainsFold(chunk.Heading, query) && !textContainsFold(strings.Join(chunk.HeadingPath, " > "), query) {
			continue
		}
		filtered = append(filtered, chunk)
	}

	if limit <= 0 || limit > len(filtered) {
		limit = len(filtered)
	}
	return filtered[:limit]
}

func pickTopChunkIDs(chunks []ResolverChunk, query string, limit int) []string {
	selected := filterResolverChunks(chunks, "", query, 0)
	if len(selected) == 0 {
		selected = chunks
	}
	sort.SliceStable(selected, func(i, j int) bool {
		if selected[i].Importance == selected[j].Importance {
			return selected[i].Ordinal < selected[j].Ordinal
		}
		return chunkImportanceRank(selected[i].Importance) < chunkImportanceRank(selected[j].Importance)
	})
	if limit <= 0 || limit > len(selected) {
		limit = len(selected)
	}
	ids := make([]string, 0, limit)
	for _, chunk := range selected[:limit] {
		ids = append(ids, chunk.ChunkID)
	}
	return ids
}

func chunkImportanceRank(v string) int {
	switch v {
	case "required":
		return 0
	case "recommended":
		return 1
	default:
		return 2
	}
}

func estimateTokens(text string) int {
	return (len([]rune(text)) + 3) / 4
}

func containsFold(values []string, target string) bool {
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value), strings.TrimSpace(target)) {
			return true
		}
	}
	return false
}

func textContainsFold(text string, query string) bool {
	return strings.Contains(strings.ToLower(text), strings.ToLower(strings.TrimSpace(query)))
}

func uniqueTrimmedStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func mustJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func decodeStringSlice(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

func compareVersions(left string, right string) int {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)
	switch {
	case left == right:
		return 0
	case right == "":
		return 1
	case left == "":
		return -1
	case left > right:
		return 1
	default:
		return -1
	}
}
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
