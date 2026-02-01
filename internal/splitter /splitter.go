package splitter

func Split(lines []string, parts int) [][]string {
	if parts <= 1 || len(lines) == 0 {
		return [][]string{lines}
	}

	if parts > len(lines) {
		parts = len(lines)
	}

	chunks := make([][]string, 0, parts)

	base := len(lines) / parts
	extra := len(lines) % parts

	start := 0
	for i := 0; i < parts; i++ {
		size := base
		if i < extra {
			size++
		}

		end := start + size
		chunks = append(chunks, lines[start:end])
		start = end
	}

	return chunks
}
