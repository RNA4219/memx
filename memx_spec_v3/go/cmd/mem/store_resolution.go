package main

import (
	"context"
	"fmt"

	"memx/api"
)

func searchAcrossStores(ctx context.Context, client api.Client, query string, topK int) (api.NotesSearchResponse, *api.Error) {
	if topK <= 0 {
		topK = 20
	}

	out := make([]api.Note, 0, topK)
	seen := make(map[string]struct{}, topK)
	appendNote := func(note api.Note) {
		if len(out) >= topK {
			return
		}
		if _, ok := seen[note.ID]; ok {
			return
		}
		seen[note.ID] = struct{}{}
		out = append(out, note)
	}

	shortResp, apiErr := client.NotesSearch(ctx, api.NotesSearchRequest{Query: query, TopK: topK})
	if apiErr != nil {
		return api.NotesSearchResponse{}, apiErr
	}
	for _, note := range shortResp.Notes {
		appendNote(note)
	}

	journalResp, apiErr := client.JournalSearch(ctx, api.JournalSearchRequest{Query: query, TopK: topK})
	if apiErr != nil {
		return api.NotesSearchResponse{}, apiErr
	}
	for _, note := range journalResp.Notes {
		appendNote(noteFromJournal(note))
	}

	knowledgeResp, apiErr := client.KnowledgeSearch(ctx, api.KnowledgeSearchRequest{Query: query, TopK: topK})
	if apiErr != nil {
		return api.NotesSearchResponse{}, apiErr
	}
	for _, note := range knowledgeResp.Notes {
		appendNote(noteFromKnowledge(note))
	}

	return api.NotesSearchResponse{Notes: out}, nil
}

func resolveNoteAcrossStores(ctx context.Context, client api.Client, id string) (interface{}, *api.Error) {
	if note, apiErr := client.NotesGet(ctx, id); apiErr == nil {
		return note, nil
	} else if apiErr.Code != api.CodeNotFound {
		return nil, apiErr
	}

	if note, apiErr := client.JournalGet(ctx, id); apiErr == nil {
		return note, nil
	} else if apiErr.Code != api.CodeNotFound {
		return nil, apiErr
	}

	if note, apiErr := client.KnowledgeGet(ctx, id); apiErr == nil {
		return note, nil
	} else if apiErr.Code != api.CodeNotFound {
		return nil, apiErr
	}

	if note, apiErr := client.ArchiveGet(ctx, id); apiErr == nil {
		return note, nil
	} else if apiErr.Code != api.CodeNotFound {
		return nil, apiErr
	}

	return nil, &api.Error{Code: api.CodeNotFound, Message: "not found"}
}

func noteFromJournal(n api.JournalNote) api.Note {
	return api.Note{
		ID:             n.ID,
		Ref:            n.Ref,
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

func noteFromKnowledge(n api.KnowledgeNote) api.Note {
	return api.Note{
		ID:             n.ID,
		Ref:            n.Ref,
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

func printResolvedNote(note interface{}) {
	switch n := note.(type) {
	case api.Note:
		fmt.Printf("# %s\n\n%s\n", n.Title, n.Body)
	case api.JournalNote:
		fmt.Printf("# %s\n\nScope: %s\n\n%s\n", n.Title, n.WorkingScope, n.Body)
	case api.KnowledgeNote:
		title := n.Title
		if n.IsPinned {
			title += " (pinned)"
		}
		fmt.Printf("# %s\n\nScope: %s\n\n%s\n", title, n.WorkingScope, n.Body)
	case api.ArchiveNote:
		fmt.Printf("# %s\n\n%s\n", n.Title, n.Body)
	default:
		fmt.Printf("%v\n", note)
	}
}
