package grep

import (
	"regexp"
	"sync"

	"github.com/kstsm/wb-l4.2/internal/constants"
	"github.com/kstsm/wb-l4.2/internal/types"
)

func Search(
	lines []string,
	pattern string,
	caseSensitive bool,
	invertMatch bool,
	lineOffset int,
	workers int,
) []types.Match {
	re, err := compileRegexp(pattern, caseSensitive)
	if err != nil {
		return nil
	}

	if workers <= 1 || len(lines) < constants.MinLinesForConcurrency {
		return searchSequential(lines, re, invertMatch, lineOffset)
	}

	return searchConcurrent(lines, re, invertMatch, lineOffset, workers)
}

func compileRegexp(pattern string, caseSensitive bool) (*regexp.Regexp, error) {
	if caseSensitive {
		return regexp.Compile(pattern)
	}
	return regexp.Compile("(?i)" + pattern)
}

func searchSequential(
	lines []string,
	re *regexp.Regexp,
	invertMatch bool,
	lineOffset int,
) []types.Match {
	results := make([]types.Match, 0)

	for i, line := range lines {
		if isMatch(re, line, invertMatch) {
			results = append(results, types.Match{
				LineNumber: i + 1 + lineOffset,
				Content:    line,
			})
		}
	}

	return results
}

func searchConcurrent(
	lines []string,
	re *regexp.Regexp,
	invertMatch bool,
	lineOffset int,
	workers int,
) []types.Match {
	results := make(chan types.Match, workers*2)
	var wg sync.WaitGroup

	chunkSize := len(lines) / workers
	if chunkSize == 0 {
		chunkSize = 1
	}

	for w := 0; w < workers; w++ {
		start := w * chunkSize
		end := start + chunkSize

		if start >= len(lines) {
			break
		}
		if end > len(lines) {
			end = len(lines)
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()

			for i := start; i < end; i++ {
				if isMatch(re, lines[i], invertMatch) {
					results <- types.Match{
						LineNumber: i + 1 + lineOffset,
						Content:    lines[i],
					}
				}
			}
		}(start, end)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	allResults := make([]types.Match, 0)
	for match := range results {
		allResults = append(allResults, match)
	}

	return allResults
}

func isMatch(re *regexp.Regexp, line string, invert bool) bool {
	matched := re.MatchString(line)
	if invert {
		return !matched
	}
	return matched
}
