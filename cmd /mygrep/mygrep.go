package mygrep

import (
	"fmt"

	"github.com/gookit/slog"
	"github.com/kstsm/wb-l4.2/internal/app/coordinator"
	"github.com/kstsm/wb-l4.2/internal/app/worker"
	"github.com/kstsm/wb-l4.2/internal/cli"
	"github.com/kstsm/wb-l4.2/internal/types"
	_ "github.com/kstsm/wb-l4.2/pkg/logger"
)

func Run(lg *slog.Logger) error {
	cfg, err := cli.ParseFlags()
	if err != nil {
		return fmt.Errorf("parsing arguments: %w", err)
	}

	if cfg.ServerMode {
		if err := worker.Run(cfg.Addr, lg); err != nil {
			return fmt.Errorf("starting worker: %w", err)
		}
		return nil
	}

	opts := types.GrepOptions{
		CaseSensitive: cfg.CaseSensitive,
		InvertMatch:   cfg.InvertMatch,
		LineNumbers:   cfg.LineNumbers,
	}
	if err := coordinator.RunWithOptions(cfg.Pattern, cfg.File, cfg.Nodes, opts, lg); err != nil {
		return fmt.Errorf("coordinator: %w", err)
	}

	return nil
}
