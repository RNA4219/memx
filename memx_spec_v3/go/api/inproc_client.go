package api

import (
	"context"

	"memx/service"
)

type InProcClient struct {
	Svc *service.Service
}

func NewInProcClient(svc *service.Service) *InProcClient {
	return &InProcClient{Svc: svc}
}

func (c *InProcClient) NotesIngest(ctx context.Context, req NotesIngestRequest) (NotesIngestResponse, *Error) {
	n, err := c.Svc.IngestShort(ctx, service.IngestNoteRequest{
		Title:       req.Title,
		Body:        req.Body,
		Summary:     req.Summary,
		SourceType:  req.SourceType,
		Origin:      req.Origin,
		SourceTrust: req.SourceTrust,
		Sensitivity: req.Sensitivity,
		Tags:        req.Tags,
	})
	if err != nil {
		return NotesIngestResponse{}, mapError(err)
	}
	return NotesIngestResponse{Note: fromServiceNote(n)}, nil
}

func (c *InProcClient) NotesSearch(ctx context.Context, req NotesSearchRequest) (NotesSearchResponse, *Error) {
	ns, err := c.Svc.SearchShort(ctx, req.Query, req.TopK)
	if err != nil {
		return NotesSearchResponse{}, mapError(err)
	}
	out := make([]Note, 0, len(ns))
	for _, n := range ns {
		out = append(out, fromServiceNote(n))
	}
	return NotesSearchResponse{Notes: out}, nil
}

func (c *InProcClient) NotesGet(ctx context.Context, id string) (Note, *Error) {
	n, err := c.Svc.GetShort(ctx, id)
	if err != nil {
		return Note{}, mapError(err)
	}
	return fromServiceNote(n), nil
}

func (c *InProcClient) GCRun(ctx context.Context, req GCRunRequest) (GCRunResponse, *Error) {
	// v1.3: GC はまだスタブ。
	// ただし API の形は固定する。
	_ = ctx
	_ = req
	return GCRunResponse{Status: "ok"}, nil
}

func (c *InProcClient) Summarize(ctx context.Context, id string) (SummarizeResponse, *Error) {
	n, err := c.Svc.SummarizeNote(ctx, id)
	if err != nil {
		return SummarizeResponse{}, mapError(err)
	}
	return SummarizeResponse{Note: fromServiceNote(n)}, nil
}

func (c *InProcClient) SummarizeBatch(ctx context.Context, req SummarizeBatchRequest) (SummarizeBatchResponse, *Error) {
	result, err := c.Svc.SummarizeNotes(ctx, req.IDs)
	if err != nil {
		return SummarizeBatchResponse{}, mapError(err)
	}
	return SummarizeBatchResponse{
		Summary:   result.Summary,
		NoteCount: result.NoteCount,
	}, nil
}

func fromServiceNote(n service.Note) Note {
	return Note{
		ID:             n.ID,
		Title:          n.Title,
		Summary:        n.Summary,
		Body:           n.Body,
		CreatedAt:      n.CreatedAt,
		UpdatedAt:      n.UpdatedAt,
		LastAccessedAt: n.LastAccessedAt,
		AccessCount:    n.AccessCount,
		SourceType:     n.SourceType,
		Origin:         n.Origin,
		SourceTrust:    n.SourceTrust,
		Sensitivity:    n.Sensitivity,
	}
}
