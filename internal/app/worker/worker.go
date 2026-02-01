package worker

import (
	"fmt"
	"net/http"

	"github.com/gookit/slog"
	"github.com/kstsm/wb-l4.2/internal/constants"
	httptransport "github.com/kstsm/wb-l4.2/internal/transport/http"
)

func Run(addr string, lg *slog.Logger) error {
	mux := http.NewServeMux()
	handler := httptransport.NewSearchHandler(lg)
	mux.HandleFunc(constants.SearchEndpoint, handler.HandleSearch)

	lg.Infof("Worker started on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		return fmt.Errorf("listening on %s: %w", addr, err)
	}

	return nil
}
