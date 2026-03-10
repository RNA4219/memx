package api

import "net/http"

func (s *HTTPServer) handleDocsIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req DocsIngestRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.DocsIngest(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleDocsResolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req DocsResolveRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.DocsResolve(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleChunksGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req ChunksGetRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.ChunksGet(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleDocsSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req DocsSearchRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.DocsSearch(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleReadsAck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req ReadsAckRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.ReadsAck(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleDocsStaleCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req DocsStaleCheckRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.DocsStaleCheck(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}

func (s *HTTPServer) handleContractsResolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req ContractsResolveRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	resp, apiErr := s.InProc.ContractsResolve(r.Context(), req)
	if apiErr != nil {
		writeErr(w, apiErr)
		return
	}
	writeOK(w, resp)
}
