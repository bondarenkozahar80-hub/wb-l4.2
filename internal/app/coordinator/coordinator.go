package coordinator

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/gookit/slog"
	"github.com/kstsm/wb-l4.2/internal/constants"
	"github.com/kstsm/wb-l4.2/internal/quorum"
	"github.com/kstsm/wb-l4.2/internal/reader"
	"github.com/kstsm/wb-l4.2/internal/splitter"
	httptransport "github.com/kstsm/wb-l4.2/internal/transport/http"
	"github.com/kstsm/wb-l4.2/internal/types"
)

type workerResult struct {
	response types.SearchResponse
	success  bool
	nodeURL  string
}

func RunWithOptions(pattern, filePath, nodesStr string, opts types.GrepOptions, lg *slog.Logger) error {
	lines, err := reader.ReadLines(filePath)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	nodes := parseNodes(nodesStr)
	if len(nodes) == 0 {
		return fmt.Errorf("servers not specified")
	}

	quorumSize := quorum.Calculate(len(nodes))
	source := filePath
	if source == "" {
		source = "stdin"
	}
	lg.Infof("Starting search: pattern=%s, source=%s, servers=%d, quorum=%d",
		pattern,
		source,
		len(nodes),
		quorumSize,
	)

	chunks := splitter.Split(lines, len(nodes))
	if len(chunks) > len(nodes) {
		chunks = chunks[:len(nodes)]
	}

	resultsChan := make(chan workerResult, len(nodes))
	var wg sync.WaitGroup

	for i, node := range nodes {
		if i >= len(chunks) {
			break
		}

		nodeURL := node
		chunk := chunks[i]

		lineOffset := 0
		for j := 0; j < i; j++ {
			lineOffset += len(chunks[j])
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			url := fmt.Sprintf("http://%s%s", nodeURL, constants.SearchEndpoint)
			req := types.SearchRequest{
				Lines:         chunk,
				Pattern:       pattern,
				CaseSensitive: opts.CaseSensitive,
				InvertMatch:   opts.InvertMatch,
				LineOffset:    lineOffset,
			}

			resp, err := httptransport.SendRequest(url, req)
			if err != nil {
				lg.Warnf("Server %s failed: %v", nodeURL, err)
				resultsChan <- workerResult{
					success: false,
					nodeURL: nodeURL,
				}
				return
			}

			lg.Infof("Server %s responded successfully with %d matches", nodeURL, len(resp.Matches))
			resultsChan <- workerResult{
				response: resp,
				success:  true,
				nodeURL:  nodeURL,
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var successfulResponses []types.SearchResponse
	var failedCount int

	for result := range resultsChan {
		if result.success {
			successfulResponses = append(successfulResponses, result.response)
		} else {
			failedCount++
		}
	}

	if !quorum.IsReached(len(successfulResponses), len(nodes)) {
		return fmt.Errorf("quorum not reached: got %d successful responses out of %d, required %d (failed: %d)",
			len(successfulResponses), len(nodes), quorumSize, failedCount)
	}

	lg.Infof("Quorum reached: %d/%d servers responded successfully", len(successfulResponses), len(nodes))

	allMatches := collectMatches(successfulResponses)
	sortMatches(allMatches)

	for _, match := range allMatches {
		if opts.LineNumbers {
			fmt.Printf("%d:%s\n", match.LineNumber, match.Content)
		} else {
			fmt.Printf("%s\n", match.Content)
		}
	}

	lg.Infof("Found matches: %d", len(allMatches))
	return nil
}

func parseNodes(nodesStr string) []string {
	parts := strings.Split(nodesStr, ",")
	var nodes []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			if !strings.Contains(part, ":") {
				part = part + constants.DefaultPort
			}
			nodes = append(nodes, part)
		}
	}
	return nodes
}

func collectMatches(responses []types.SearchResponse) []types.Match {
	seen := make(map[string]bool)
	var allMatches []types.Match

	for _, resp := range responses {
		for _, match := range resp.Matches {
			key := fmt.Sprintf("%d:%s", match.LineNumber, match.Content)
			if !seen[key] {
				seen[key] = true
				allMatches = append(allMatches, match)
			}
		}
	}

	return allMatches
}

func sortMatches(matches []types.Match) {
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].LineNumber < matches[j].LineNumber
	})
}
