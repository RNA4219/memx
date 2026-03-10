package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"memx/service"
)

// HTTPServer は /v1/* の JSON API を提供する。
type HTTPServer struct {
	InProc *InProcClient
}

func NewHTTPServer(svc *service.Service) *HTTPServer {
	return &HTTPServer{InProc: NewInProcClient(svc)}
}

func (s *HTTPServer) Handler() http.Handler {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/healthz", s.handleHealthz)

	// Short store
	mux.HandleFunc("/v1/notes:ingest", s.handleNotesIngest)
	mux.HandleFunc("/v1/notes:search", s.handleNotesSearch)
	mux.HandleFunc("/v1/notes:recall", s.handleRecall)
	mux.HandleFunc("/v1/notes:summarize", s.handleSummarize)
	mux.HandleFunc("/v1/notes:summarize-batch", s.handleSummarizeBatch)
	mux.HandleFunc("/v1/gc:run", s.handleGCRun)
	mux.HandleFunc("/v1/notes/", s.handleNotesGet)

	// Journal store
	mux.HandleFunc("/v1/journal:ingest", s.handleJournalIngest)
	mux.HandleFunc("/v1/journal:search", s.handleJournalSearch)
	mux.HandleFunc("/v1/journal:list-by-scope", s.handleJournalListByScope)
	mux.HandleFunc("/v1/journal/", s.handleJournalGet)

	// Knowledge store
	mux.HandleFunc("/v1/knowledge:ingest", s.handleKnowledgeIngest)
	mux.HandleFunc("/v1/knowledge:search", s.handleKnowledgeSearch)
	mux.HandleFunc("/v1/knowledge:list-by-scope", s.handleKnowledgeListByScope)
	mux.HandleFunc("/v1/knowledge:list-pinned", s.handleKnowledgeListPinned)
	mux.HandleFunc("/v1/knowledge/", s.handleKnowledgeGet)

	// Archive store
	mux.HandleFunc("/v1/archive", s.handleArchiveList)
	mux.HandleFunc("/v1/archive/", s.handleArchiveGetOrRestore)

	// Resolver API (P4)
	mux.HandleFunc("/v1/resolve", s.handleResolveRef)
	mux.HandleFunc("/v1/resolve-many", s.handleResolveMany)
	mux.HandleFunc("/v1/summary", s.handleLoadSummary)
	mux.HandleFunc("/v1/raw", s.handleLoadSelectedRaw)
	mux.HandleFunc("/v1/bundle:build", s.handleBuildBundle)

	// Cookbook resolver
	mux.HandleFunc("/v1/docs:ingest", s.handleDocsIngest)
	mux.HandleFunc("/v1/docs:resolve", s.handleDocsResolve)
	mux.HandleFunc("/v1/chunks:get", s.handleChunksGet)
	mux.HandleFunc("/v1/docs:search", s.handleDocsSearch)
	mux.HandleFunc("/v1/reads:ack", s.handleReadsAck)
	mux.HandleFunc("/v1/docs:stale-check", s.handleDocsStaleCheck)
	mux.HandleFunc("/v1/contracts:resolve", s.handleContractsResolve)

	return mux
}

// -------------------- Common Helpers --------------------

func (s *HTTPServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func writeOK(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, e *Error) {
	status := http.StatusInternalServerError
	if e != nil {
		switch e.Code {
		case CodeInvalidArgument:
			status = http.StatusBadRequest
		case CodeNotFound:
			status = http.StatusNotFound
		case CodeConflict:
			status = http.StatusConflict
		case CodeGatekeepDeny:
			status = http.StatusForbidden
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(e)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) *Error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return &Error{Code: CodeInvalidArgument, Message: "invalid json"}
	}
	return nil
}

func extractID(w http.ResponseWriter, r *http.Request, prefix string) (string, bool) {
	id := strings.TrimPrefix(r.URL.Path, prefix)
	id = strings.TrimSpace(id)
	if id == "" {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "id is required"})
		return "", false
	}
	return id, true
}

// -------------------- Short Store --------------------

func (s *HTTPServer) handleNotesIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req NotesIngestRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.NotesIngest(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleNotesSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req NotesSearchRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.NotesSearch(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleRecall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req RecallRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.Recall(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleNotesGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id, ok := extractID(w, r, "/v1/notes/")
	if !ok {
		return
	}
	n, apiErr := s.InProc.NotesGet(r.Context(), id)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, n)
}

func (s *HTTPServer) handleGCRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req GCRunRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.GCRun(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleSummarize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req SummarizeRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.Summarize(r.Context(), req.ID)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleSummarizeBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req SummarizeBatchRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.SummarizeBatch(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

// -------------------- Journal Store --------------------

func (s *HTTPServer) handleJournalIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req JournalIngestRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.JournalIngest(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleJournalSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req JournalSearchRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.JournalSearch(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleJournalGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id, ok := extractID(w, r, "/v1/journal/")
	if !ok {
		return
	}
	n, apiErr := s.InProc.JournalGet(r.Context(), id)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, n)
}

func (s *HTTPServer) handleJournalListByScope(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req JournalListByScopeRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.JournalListByScope(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

// -------------------- Knowledge Store --------------------

func (s *HTTPServer) handleKnowledgeIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req KnowledgeIngestRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.KnowledgeIngest(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleKnowledgeSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req KnowledgeSearchRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.KnowledgeSearch(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleKnowledgeGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id, ok := extractID(w, r, "/v1/knowledge/")
	if !ok {
		return
	}
	// Check for pin/unpin actions
	if strings.HasSuffix(id, ":pin") {
		s.handleKnowledgePinAction(w, r, strings.TrimSuffix(id, ":pin"))
		return
	}
	if strings.HasSuffix(id, ":unpin") {
		s.handleKnowledgeUnpinAction(w, r, strings.TrimSuffix(id, ":unpin"))
		return
	}
	n, apiErr := s.InProc.KnowledgeGet(r.Context(), id)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, n)
}

func (s *HTTPServer) handleKnowledgePinAction(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	resp, apiErr := s.InProc.KnowledgePin(r.Context(), id)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleKnowledgeUnpinAction(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	resp, apiErr := s.InProc.KnowledgeUnpin(r.Context(), id)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleKnowledgeListByScope(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req KnowledgeListByScopeRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.KnowledgeListByScope(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleKnowledgeListPinned(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req KnowledgeListPinnedRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.KnowledgeListPinned(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

// -------------------- Archive Store --------------------

func (s *HTTPServer) handleArchiveList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := json.Number(l).Int64(); err == nil && val > 0 {
			limit = int(val)
		}
	}
	resp, apiErr := s.InProc.ArchiveList(r.Context(), ArchiveListRequest{Limit: limit})
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleArchiveGetOrRestore(w http.ResponseWriter, r *http.Request) {
	id, ok := extractID(w, r, "/v1/archive/")
	if !ok {
		return
	}
	if strings.HasSuffix(id, ":restore") {
		s.handleArchiveRestore(w, r, strings.TrimSuffix(id, ":restore"))
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	n, apiErr := s.InProc.ArchiveGet(r.Context(), id)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, n)
}

func (s *HTTPServer) handleArchiveRestore(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	resp, apiErr := s.InProc.ArchiveRestore(r.Context(), id)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

// -------------------- Resolver API (P4) --------------------

func (s *HTTPServer) handleResolveRef(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req ResolveRefRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.ResolveRef(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleResolveMany(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req ResolveManyRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.ResolveMany(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleLoadSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req LoadSummaryRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.LoadSummary(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleLoadSelectedRaw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req LoadSelectedRawRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.LoadSelectedRaw(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleBuildBundle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req BuildBundleRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.BuildBundle(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}
