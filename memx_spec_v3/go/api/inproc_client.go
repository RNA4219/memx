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
	result, err := c.Svc.GCShort(ctx, service.GCRequest{
		Target:  req.Target,
		DryRun:  req.Options.DryRun,
		Enabled: true, // in-proc は常に有効（CLI側で制御）
	})
	if err != nil {
		return GCRunResponse{}, mapError(err)
	}

	// dry-run の場合は詳細を含める
	if result.DryRun && result.DryRunResult != nil {
		return GCRunResponse{
			Status: result.DryRunResult.ToJSON(),
		}, nil
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

// -------------------- Chronicle --------------------

func (c *InProcClient) ChronicleIngest(ctx context.Context, req ChronicleIngestRequest) (ChronicleIngestResponse, *Error) {
	n, err := c.Svc.IngestChronicle(ctx, service.IngestChronicleRequest{
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
		return ChronicleIngestResponse{}, mapError(err)
	}
	return ChronicleIngestResponse{Note: fromServiceChronicleNote(n)}, nil
}

func (c *InProcClient) ChronicleSearch(ctx context.Context, req ChronicleSearchRequest) (ChronicleSearchResponse, *Error) {
	ns, err := c.Svc.SearchChronicle(ctx, req.Query, req.TopK)
	if err != nil {
		return ChronicleSearchResponse{}, mapError(err)
	}
	out := make([]ChronicleNote, 0, len(ns))
	for _, n := range ns {
		out = append(out, fromServiceChronicleNote(n))
	}
	return ChronicleSearchResponse{Notes: out}, nil
}

func (c *InProcClient) ChronicleGet(ctx context.Context, id string) (ChronicleNote, *Error) {
	n, err := c.Svc.GetChronicle(ctx, id)
	if err != nil {
		return ChronicleNote{}, mapError(err)
	}
	return fromServiceChronicleNote(n), nil
}

func (c *InProcClient) ChronicleListByScope(ctx context.Context, req ChronicleListByScopeRequest) (ChronicleListByScopeResponse, *Error) {
	ns, err := c.Svc.ListChronicleByScope(ctx, req.WorkingScope, req.Limit)
	if err != nil {
		return ChronicleListByScopeResponse{}, mapError(err)
	}
	out := make([]ChronicleNote, 0, len(ns))
	for _, n := range ns {
		out = append(out, fromServiceChronicleNote(n))
	}
	return ChronicleListByScopeResponse{Notes: out}, nil
}

func fromServiceChronicleNote(n service.ChronicleNote) ChronicleNote {
	return ChronicleNote{
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
		WorkingScope:   n.WorkingScope,
		IsPinned:       n.IsPinned,
	}
}

// -------------------- Memopedia --------------------

func (c *InProcClient) MemopediaIngest(ctx context.Context, req MemopediaIngestRequest) (MemopediaIngestResponse, *Error) {
	n, err := c.Svc.IngestMemopedia(ctx, service.IngestMemopediaRequest{
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
		return MemopediaIngestResponse{}, mapError(err)
	}
	return MemopediaIngestResponse{Note: fromServiceMemopediaNote(n)}, nil
}

func (c *InProcClient) MemopediaSearch(ctx context.Context, req MemopediaSearchRequest) (MemopediaSearchResponse, *Error) {
	ns, err := c.Svc.SearchMemopedia(ctx, req.Query, req.TopK)
	if err != nil {
		return MemopediaSearchResponse{}, mapError(err)
	}
	out := make([]MemopediaNote, 0, len(ns))
	for _, n := range ns {
		out = append(out, fromServiceMemopediaNote(n))
	}
	return MemopediaSearchResponse{Notes: out}, nil
}

func (c *InProcClient) MemopediaGet(ctx context.Context, id string) (MemopediaNote, *Error) {
	n, err := c.Svc.GetMemopedia(ctx, id)
	if err != nil {
		return MemopediaNote{}, mapError(err)
	}
	return fromServiceMemopediaNote(n), nil
}

func (c *InProcClient) MemopediaListByScope(ctx context.Context, req MemopediaListByScopeRequest) (MemopediaListByScopeResponse, *Error) {
	ns, err := c.Svc.ListMemopediaByScope(ctx, req.WorkingScope, req.Limit)
	if err != nil {
		return MemopediaListByScopeResponse{}, mapError(err)
	}
	out := make([]MemopediaNote, 0, len(ns))
	for _, n := range ns {
		out = append(out, fromServiceMemopediaNote(n))
	}
	return MemopediaListByScopeResponse{Notes: out}, nil
}

func (c *InProcClient) MemopediaListPinned(ctx context.Context, req MemopediaListPinnedRequest) (MemopediaListPinnedResponse, *Error) {
	ns, err := c.Svc.ListPinnedMemopedia(ctx, req.WorkingScope, req.Limit)
	if err != nil {
		return MemopediaListPinnedResponse{}, mapError(err)
	}
	out := make([]MemopediaNote, 0, len(ns))
	for _, n := range ns {
		out = append(out, fromServiceMemopediaNote(n))
	}
	return MemopediaListPinnedResponse{Notes: out}, nil
}

func (c *InProcClient) MemopediaPin(ctx context.Context, id string) (PinResponse, *Error) {
	if err := c.Svc.PinMemopedia(ctx, id); err != nil {
		return PinResponse{}, mapError(err)
	}
	return PinResponse{Success: true}, nil
}

func (c *InProcClient) MemopediaUnpin(ctx context.Context, id string) (UnpinResponse, *Error) {
	if err := c.Svc.UnpinMemopedia(ctx, id); err != nil {
		return UnpinResponse{}, mapError(err)
	}
	return UnpinResponse{Success: true}, nil
}

func fromServiceMemopediaNote(n service.MemopediaNote) MemopediaNote {
	return MemopediaNote{
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
		WorkingScope:   n.WorkingScope,
		IsPinned:       n.IsPinned,
	}
}

// -------------------- Archive --------------------

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
	out := make([]ArchiveNote, 0, len(ns))
	for _, n := range ns {
		out = append(out, fromServiceArchiveNote(n))
	}
	return ArchiveListResponse{Notes: out}, nil
}

func (c *InProcClient) ArchiveRestore(ctx context.Context, id string) (ArchiveRestoreResponse, *Error) {
	n, err := c.Svc.RestoreFromArchive(ctx, id)
	if err != nil {
		return ArchiveRestoreResponse{}, mapError(err)
	}
	return ArchiveRestoreResponse{Note: fromServiceNote(n)}, nil
}

func fromServiceArchiveNote(n service.ArchiveNote) ArchiveNote {
	return ArchiveNote{
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
