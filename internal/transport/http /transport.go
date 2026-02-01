package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gookit/slog"
	"github.com/kstsm/wb-l4.2/internal/constants"
	"github.com/kstsm/wb-l4.2/internal/grep"
	"github.com/kstsm/wb-l4.2/internal/types"
)

type SearchHandler struct {
	logger *slog.Logger
}

func NewSearchHandler(logger *slog.Logger) *SearchHandler {
	return &SearchHandler{
		logger: logger,
	}
}

func (h *SearchHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method must be POST", http.StatusMethodNotAllowed)
		return
	}

	var req types.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("Error decoding request: %v", err)
		http.Error(w, fmt.Sprintf("error decoding: %v", err), http.StatusBadRequest)
		return
	}

	matches := grep.Search(req.Lines,
		req.Pattern,
		req.CaseSensitive,
		req.InvertMatch,
		req.LineOffset,
		constants.DefaultWorkers,
	)

	resp := types.SearchResponse{
		Matches: matches,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorf("Error encoding response: %v", err)
		http.Error(w, fmt.Sprintf("error encoding: %v", err), http.StatusInternalServerError)
		return
	}
}

func SendRequest(url string, req types.SearchRequest) (types.SearchResponse, error) {
	var resp types.SearchResponse

	body, err := json.Marshal(req)
	if err != nil {
		return resp, fmt.Errorf("encoding request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return resp, fmt.Errorf("creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: constants.DefaultTimeout * time.Second,
	}

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return resp, fmt.Errorf("sending request to %s: %w", url, err)
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		return resp, fmt.Errorf("server error %s: %s", httpResp.Status, string(bodyBytes))
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return resp, fmt.Errorf("decoding response: %w", err)
	}

	return resp, nil
}
