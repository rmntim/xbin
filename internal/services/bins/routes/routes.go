package routes

import (
	"log/slog"
	"net/http"

	"github.com/rmntim/xbin/internal/services/bins"
)

func Register(mux *http.ServeMux, log *slog.Logger, srv bins.Service) {
	mux.Handle("GET /bin", getBin(srv, log))
	mux.Handle("POST /bin", createBin(srv, log))
}

func getBin(srv bins.Service, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("getBin")
		_ = srv

		w.WriteHeader(http.StatusOK)
	}
}

func createBin(srv bins.Service, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("createBin")
		_ = srv

		w.WriteHeader(http.StatusOK)
	}
}
