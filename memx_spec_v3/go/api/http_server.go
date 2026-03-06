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

	mux.HandleFunc("/v1/notes:ingest", s.handleNotesIngest)
	mux.HandleFunc("/v1/notes:search", s.handleNotesSearch)
	mux.HandleFunc("/v1/notes:summarize", s.handleSummarize)
	mux.HandleFunc("/v1/notes:summarize-batch", s.handleSummarizeBatch)
	mux.HandleFunc("/v1/gc:run", s.handleGCRun)
	mux.HandleFunc("/v1/notes/", s.handleNotesGet)

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
