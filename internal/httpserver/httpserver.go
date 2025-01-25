package httpserver

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/rmntim/xbin/internal/httpserver/middleware"
	"github.com/rmntim/xbin/internal/services/bins"
	binRoutes "github.com/rmntim/xbin/internal/services/bins/routes"
)

const FIVE_SECONDS = 5 * time.Second

func NewServer(address string, log *slog.Logger, binService bins.Service) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./static")))
	binRoutes.Register(mux, log, binService)

	handler := middleware.NewLogMiddleware(log)(mux)

	srv := &http.Server{
		Addr:              address,
		Handler:           handler,
		ReadTimeout:       FIVE_SECONDS,
		ReadHeaderTimeout: FIVE_SECONDS,
		WriteTimeout:      FIVE_SECONDS,
		IdleTimeout:       FIVE_SECONDS,
	}

	return srv
}
