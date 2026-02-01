package cli

import (
	"errors"
	"flag"
	"os"

	"github.com/kstsm/wb-l4.2/internal/constants"
	"github.com/kstsm/wb-l4.2/internal/types"
)

func Parse(args []string) (types.Config, error) {
	var cfg types.Config

	fs := flag.NewFlagSet("mygrep", flag.ContinueOnError)

	fs.BoolVar(&cfg.ServerMode, "server", false, "run in server mode")
	fs.StringVar(&cfg.Addr, "addr", constants.DefaultPort, "server address")

	fs.StringVar(&cfg.File, "file", "", "file path")
	fs.StringVar(&cfg.Nodes, "nodes", "", "comma-separated servers")

	fs.BoolVar(&cfg.CaseSensitive, "case", false, "case sensitive")
	fs.BoolVar(&cfg.InvertMatch, "v", false, "invert match")
	fs.BoolVar(&cfg.LineNumbers, "n", false, "show line numbers")

	var patternArg string
	var flagArgs []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if len(arg) > 0 && arg[0] == '-' {
			flagArgs = append(flagArgs, arg)
			if i+1 < len(args) && len(args[i+1]) > 0 && args[i+1][0] != '-' {
				flagArgs = append(flagArgs, args[i+1])
				i++
			}
		} else if patternArg == "" {
			patternArg = arg
		}
	}

	if err := fs.Parse(flagArgs); err != nil {
		return cfg, err
	}

	if !cfg.ServerMode {
		if patternArg == "" {
			return cfg, errors.New("pattern not specified")
		}
		cfg.Pattern = patternArg

		if cfg.Nodes == "" {
			return cfg, errors.New("servers not specified")
		}
	}

	return cfg, nil
}

func ParseFlags() (types.Config, error) {
	return Parse(os.Args[1:])
}
