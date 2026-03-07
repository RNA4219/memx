package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"memx/service"
)

// HTTPServer は /v1/* の JSON API を提供する。
// ローカル利用が前提（127.0.0.1 / unix socket）。
type HTTPServer struct {
	InProc *InProcClient
}

func NewHTTPServer(svc *service.Service) *HTTPServer {
	return &HTTPServer{InProc: NewInProcClient(svc)}
}

func (s *HTTPServer) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Short store
	mux.HandleFunc("/v1/notes:ingest", s.handleNotesIngest)
	mux.HandleFunc("/v1/notes:search", s.handleNotesSearch)
	mux.HandleFunc("/v1/notes:summarize", s.handleSummarize)
	mux.HandleFunc("/v1/notes:summarize-batch", s.handleSummarizeBatch)
	mux.HandleFunc("/v1/gc:run", s.handleGCRun)
	mux.HandleFunc("/v1/notes/", s.handleNotesGet)

	// Chronicle store
	mux.HandleFunc("/v1/chronicle:ingest", s.handleChronicleIngest)
	mux.HandleFunc("/v1/chronicle:search", s.handleChronicleSearch)
	mux.HandleFunc("/v1/chronicle:list-by-scope", s.handleChronicleListByScope)
	mux.HandleFunc("/v1/chronicle/", s.handleChronicleGet)

	// Memopedia store
	mux.HandleFunc("/v1/memopedia:ingest", s.handleMemopediaIngest)
	mux.HandleFunc("/v1/memopedia:search", s.handleMemopediaSearch)
	mux.HandleFunc("/v1/memopedia:list-by-scope", s.handleMemopediaListByScope)
	mux.HandleFunc("/v1/memopedia:list-pinned", s.handleMemopediaListPinned)
	mux.HandleFunc("/v1/memopedia/", s.handleMemopediaGet)

	// Archive store
	mux.HandleFunc("/v1/archive", s.handleArchiveList)
	mux.HandleFunc("/v1/archive/", s.handleArchiveGetOrRestore)

	return mux
}

func (s *HTTPServer) handleNotesIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req NotesIngestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "invalid json"})
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "invalid json"})
		return
	}
	resp, apiErr := s.InProc.NotesSearch(r.Context(), req)
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
	id := strings.TrimPrefix(r.URL.Path, "/v1/notes/")
	id = strings.TrimSpace(id)
	if id == "" {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "id is required"})
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "invalid json"})
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "invalid json"})
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "invalid json"})
		return
	}
	resp, apiErr := s.InProc.SummarizeBatch(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
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
		default:
			status = http.StatusInternalServerError
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(e)
}

// -------------------- Chronicle --------------------

func (s *HTTPServer) handleChronicleIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req ChronicleIngestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "invalid json"})
		return
	}
	resp, apiErr := s.InProc.ChronicleIngest(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleChronicleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req ChronicleSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "invalid json"})
		return
	}
	resp, apiErr := s.InProc.ChronicleSearch(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleChronicleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/v1/chronicle/")
	id = strings.TrimSpace(id)
	if id == "" {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "id is required"})
		return
	}
	n, apiErr := s.InProc.ChronicleGet(r.Context(), id)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, n)
}

func (s *HTTPServer) handleChronicleListByScope(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req ChronicleListByScopeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "invalid json"})
		return
	}
	resp, apiErr := s.InProc.ChronicleListByScope(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

// -------------------- Memopedia --------------------

func (s *HTTPServer) handleMemopediaIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req MemopediaIngestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "invalid json"})
		return
	}
	resp, apiErr := s.InProc.MemopediaIngest(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleMemopediaSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req MemopediaSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "invalid json"})
		return
	}
	resp, apiErr := s.InProc.MemopediaSearch(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleMemopediaGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/v1/memopedia/")
	id = strings.TrimSpace(id)
	if id == "" {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "id is required"})
		return
	}
	// Check for pin/unpin actions
	if strings.HasSuffix(id, ":pin") {
		s.handleMemopediaPinAction(w, r, strings.TrimSuffix(id, ":pin"))
		return
	}
	if strings.HasSuffix(id, ":unpin") {
		s.handleMemopediaUnpinAction(w, r, strings.TrimSuffix(id, ":unpin"))
		return
	}
	n, apiErr := s.InProc.MemopediaGet(r.Context(), id)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, n)
}

func (s *HTTPServer) handleMemopediaPinAction(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	resp, apiErr := s.InProc.MemopediaPin(r.Context(), id)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleMemopediaUnpinAction(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	resp, apiErr := s.InProc.MemopediaUnpin(r.Context(), id)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleMemopediaListByScope(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req MemopediaListByScopeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "invalid json"})
		return
	}
	resp, apiErr := s.InProc.MemopediaListByScope(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleMemopediaListPinned(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req MemopediaListPinnedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "invalid json"})
		return
	}
	resp, apiErr := s.InProc.MemopediaListPinned(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

// -------------------- Archive --------------------

func (s *HTTPServer) handleArchiveList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		var lim int
		if val, err := json.Number(l).Int64(); err == nil {
			lim = int(val)
			if lim > 0 {
				limit = lim
			}
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
	id := strings.TrimPrefix(r.URL.Path, "/v1/archive/")
	id = strings.TrimSpace(id)
	if id == "" {
		writeErr(w, &Error{Code: CodeInvalidArgument, Message: "id is required"})
		return
	}

	// Check for restore action
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
