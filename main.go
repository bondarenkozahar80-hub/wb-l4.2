package main

import (
	"os"

	"github.com/kstsm/wb-l4.2/cmd/mygrep"
	"github.com/kstsm/wb-l4.2/pkg/logger"
)

func main() {
	lg := logger.NewSlogLogger()

	if err := mygrep.Run(lg); err != nil {
		lg.Errorf("application error: %v", err)
		os.Exit(1)
	}
}
