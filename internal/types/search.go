package types

type Match struct {
	LineNumber int
	Content    string
}

type SearchRequest struct {
	Lines         []string `json:"lines"`
	Pattern       string   `json:"pattern"`
	CaseSensitive bool     `json:"case_sensitive"`
	InvertMatch   bool     `json:"invert_match"`
	LineOffset    int      `json:"line_offset"`
}

type SearchResponse struct {
	Matches []Match `json:"matches"`
	Error   string  `json:"error,omitempty"`
}
