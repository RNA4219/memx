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

// -------------------- Helpers --------------------

func fromServiceNote(n service.Note) Note {
	return Note{
		ID:             n.ID,
		Ref:            NewTypedRef(EntityTypeEvidence, n.ID),
		Title:          n.Title,
		Summary:        n.Summary,
		Body:           n.Body,
		CreatedAt:      n.CreatedAt,
		UpdatedAt:       n.UpdatedAt,
		LastAccessedAt: n.LastAccessedAt,
		AccessCount:    n.AccessCount,
		SourceType:     n.SourceType,
		Origin:         n.Origin,
		SourceTrust:    n.SourceTrust,
		Sensitivity:    n.Sensitivity,
	}
}

func fromServiceJournalNote(n service.JournalNote) JournalNote {
	return JournalNote{
		NoteBase: NoteBase{
			ID:             n.ID,
			Ref:            NewTypedRef(EntityTypeEvidence, n.ID),
			Title:          n.Title,
			Summary:        n.Summary,
			Body:           n.Body,
			CreatedAt:      n.CreatedAt,
			UpdatedAt:       n.UpdatedAt,
			LastAccessedAt:  n.LastAccessedAt,
			AccessCount:     n.AccessCount,
			SourceType:      n.SourceType,
			Origin:          n.Origin,
			SourceTrust:     n.SourceTrust,
			Sensitivity:     n.Sensitivity,
		},
		WorkingScope: n.WorkingScope,
		IsPinned:     n.IsPinned,
	}
}

func fromServiceKnowledgeNote(n service.KnowledgeNote) KnowledgeNote {
	return KnowledgeNote{
		NoteBase: NoteBase{
			ID:             n.ID,
			Ref:            NewTypedRef(EntityTypeKnowledge, n.ID),
			Title:          n.Title,
			Summary:        n.Summary,
			Body:           n.Body,
			CreatedAt:      n.CreatedAt,
			UpdatedAt:       n.UpdatedAt,
			LastAccessedAt:  n.LastAccessedAt,
			AccessCount:     n.AccessCount,
			SourceType:      n.SourceType,
			Origin:          n.Origin,
			SourceTrust:     n.SourceTrust,
			Sensitivity:     n.Sensitivity,
		},
		WorkingScope: n.WorkingScope,
		IsPinned:     n.IsPinned,
	}
}

func fromServiceArchiveNote(n service.ArchiveNote) ArchiveNote {
	return ArchiveNote{
		ID:             n.ID,
		Ref:            NewTypedRef(EntityTypeEvidence, n.ID),
		Title:          n.Title,
		Summary:        n.Summary,
		Body:           n.Body,
		CreatedAt:      n.CreatedAt,
		UpdatedAt:       n.UpdatedAt,
		LastAccessedAt:  n.LastAccessedAt,
		AccessCount:     n.AccessCount,
		SourceType:      n.SourceType,
		Origin:          n.Origin,
		SourceTrust:     n.SourceTrust,
		Sensitivity:     n.Sensitivity,
	}
}

// -------------------- Short Store --------------------

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
	return NotesSearchResponse{Notes: mapNotes(ns)}, nil
}

func (c *InProcClient) NotesGet(ctx context.Context, id string) (Note, *Error) {
	n, err := c.Svc.GetShort(ctx, id)
	if err != nil {
		return Note{}, mapError(err)
	}
	return fromServiceNote(n), nil
}

func (c *InProcClient) GCRun(ctx context.Context, req GCRunRequest) (GCRunResponse, *Error) {
	result, err := c.Svc.GCShort(ctx, service.GCRequest{
		Target:  req.Target,
		DryRun:  req.Options.DryRun,
		Enabled: true,
	})
	if err != nil {
		return GCRunResponse{}, mapError(err)
	}
	if result.DryRun && result.DryRunResult != nil {
		return GCRunResponse{Status: result.DryRunResult.ToJSON()}, nil
	}
	return GCRunResponse{Status: result.Status}, nil
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
	return SummarizeBatchResponse{Summary: result.Summary, NoteCount: result.NoteCount}, nil
}

// -------------------- Journal Store --------------------

func (c *InProcClient) JournalIngest(ctx context.Context, req JournalIngestRequest) (JournalIngestResponse, *Error) {
	n, err := c.Svc.IngestJournal(ctx, service.IngestJournalRequest{
		Title:        req.Title,
		Body:         req.Body,
		Summary:      req.Summary,
		SourceType:   req.SourceType,
		Origin:       req.Origin,
		SourceTrust:  req.SourceTrust,
		Sensitivity:  req.Sensitivity,
		Tags:         req.Tags,
		WorkingScope: req.WorkingScope,
		IsPinned:     req.IsPinned,
	})
	if err != nil {
		return JournalIngestResponse{}, mapError(err)
	}
	return JournalIngestResponse{Note: fromServiceJournalNote(n)}, nil
}

func (c *InProcClient) JournalSearch(ctx context.Context, req JournalSearchRequest) (JournalSearchResponse, *Error) {
	ns, err := c.Svc.SearchJournal(ctx, req.Query, req.TopK)
	if err != nil {
		return JournalSearchResponse{}, mapError(err)
	}
	return JournalSearchResponse{Notes: mapJournalNotes(ns)}, nil
}

func (c *InProcClient) JournalGet(ctx context.Context, id string) (JournalNote, *Error) {
	n, err := c.Svc.GetJournal(ctx, id)
	if err != nil {
		return JournalNote{}, mapError(err)
	}
	return fromServiceJournalNote(n), nil
}

func (c *InProcClient) JournalListByScope(ctx context.Context, req JournalListByScopeRequest) (JournalListByScopeResponse, *Error) {
	ns, err := c.Svc.ListJournalByScope(ctx, req.WorkingScope, req.Limit)
	if err != nil {
		return JournalListByScopeResponse{}, mapError(err)
	}
	return JournalListByScopeResponse{Notes: mapJournalNotes(ns)}, nil
}

// -------------------- Knowledge Store --------------------

func (c *InProcClient) KnowledgeIngest(ctx context.Context, req KnowledgeIngestRequest) (KnowledgeIngestResponse, *Error) {
	n, err := c.Svc.IngestKnowledge(ctx, service.IngestKnowledgeRequest{
		Title:        req.Title,
		Body:         req.Body,
		Summary:      req.Summary,
		SourceType:   req.SourceType,
		Origin:       req.Origin,
		SourceTrust:  req.SourceTrust,
		Sensitivity:  req.Sensitivity,
		Tags:         req.Tags,
		WorkingScope: req.WorkingScope,
		IsPinned:     req.IsPinned,
	})
	if err != nil {
		return KnowledgeIngestResponse{}, mapError(err)
	}
	return KnowledgeIngestResponse{Note: fromServiceKnowledgeNote(n)}, nil
}

func (c *InProcClient) KnowledgeSearch(ctx context.Context, req KnowledgeSearchRequest) (KnowledgeSearchResponse, *Error) {
	ns, err := c.Svc.SearchKnowledge(ctx, req.Query, req.TopK)
	if err != nil {
		return KnowledgeSearchResponse{}, mapError(err)
	}
	return KnowledgeSearchResponse{Notes: mapKnowledgeNotes(ns)}, nil
}

func (c *InProcClient) KnowledgeGet(ctx context.Context, id string) (KnowledgeNote, *Error) {
	n, err := c.Svc.GetKnowledge(ctx, id)
	if err != nil {
		return KnowledgeNote{}, mapError(err)
	}
	return fromServiceKnowledgeNote(n), nil
}

func (c *InProcClient) KnowledgeListByScope(ctx context.Context, req KnowledgeListByScopeRequest) (KnowledgeListByScopeResponse, *Error) {
	ns, err := c.Svc.ListKnowledgeByScope(ctx, req.WorkingScope, req.Limit)
	if err != nil {
		return KnowledgeListByScopeResponse{}, mapError(err)
	}
	return KnowledgeListByScopeResponse{Notes: mapKnowledgeNotes(ns)}, nil
}

func (c *InProcClient) KnowledgeListPinned(ctx context.Context, req KnowledgeListPinnedRequest) (KnowledgeListPinnedResponse, *Error) {
	ns, err := c.Svc.ListPinnedKnowledge(ctx, req.WorkingScope, req.Limit)
	if err != nil {
		return KnowledgeListPinnedResponse{}, mapError(err)
	}
	return KnowledgeListPinnedResponse{Notes: mapKnowledgeNotes(ns)}, nil
}

func (c *InProcClient) KnowledgePin(ctx context.Context, id string) (PinResponse, *Error) {
	if err := c.Svc.PinKnowledge(ctx, id); err != nil {
		return PinResponse{}, mapError(err)
	}
	return PinResponse{Success: true}, nil
}

func (c *InProcClient) KnowledgeUnpin(ctx context.Context, id string) (UnpinResponse, *Error) {
	if err := c.Svc.UnpinKnowledge(ctx, id); err != nil {
		return UnpinResponse{}, mapError(err)
	}
	return UnpinResponse{Success: true}, nil
}

// -------------------- Archive Store --------------------

func (c *InProcClient) ArchiveGet(ctx context.Context, id string) (ArchiveNote, *Error) {
	n, err := c.Svc.GetArchive(ctx, id)
	if err != nil {
		return ArchiveNote{}, mapError(err)
	}
	return fromServiceArchiveNote(n), nil
}

func (c *InProcClient) ArchiveList(ctx context.Context, req ArchiveListRequest) (ArchiveListResponse, *Error) {
	ns, err := c.Svc.ListArchive(ctx, req.Limit)
	if err != nil {
		return ArchiveListResponse{}, mapError(err)
	}
	return ArchiveListResponse{Notes: mapArchiveNotes(ns)}, nil
}

func (c *InProcClient) ArchiveRestore(ctx context.Context, id string) (ArchiveRestoreResponse, *Error) {
	n, err := c.Svc.RestoreFromArchive(ctx, id)
	if err != nil {
		return ArchiveRestoreResponse{}, mapError(err)
	}
	return ArchiveRestoreResponse{Note: fromServiceNote(n)}, nil
}

// -------------------- Slice Mappers --------------------

func mapNotes(ns []service.Note) []Note {
	out := make([]Note, 0, len(ns))
	for _, n := range ns {
		out = append(out, fromServiceNote(n))
	}
	return out
}

func mapJournalNotes(ns []service.JournalNote) []JournalNote {
	out := make([]JournalNote, 0, len(ns))
	for _, n := range ns {
		out = append(out, fromServiceJournalNote(n))
	}
	return out
}

func mapKnowledgeNotes(ns []service.KnowledgeNote) []KnowledgeNote {
	out := make([]KnowledgeNote, 0, len(ns))
	for _, n := range ns {
		out = append(out, fromServiceKnowledgeNote(n))
	}
	return out
}

func mapArchiveNotes(ns []service.ArchiveNote) []ArchiveNote {
	out := make([]ArchiveNote, 0, len(ns))
	for _, n := range ns {
		out = append(out, fromServiceArchiveNote(n))
	}
	return out
}