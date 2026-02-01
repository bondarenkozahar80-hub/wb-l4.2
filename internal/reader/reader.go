package reader

import (
	"bufio"
	"os"
)

func ReadLines(filename string) ([]string, error) {
	var file *os.File
	var err error

	if filename == "" {
		file = os.Stdin
	} else {
		file, err = os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = file.Close()
		}()
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
