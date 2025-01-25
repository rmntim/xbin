package routes

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/rmntim/xbin/internal/services/bins"
	binErr "github.com/rmntim/xbin/internal/services/bins/errors"
	"github.com/rmntim/xbin/internal/services/bins/models"
	"github.com/rmntim/xbin/internal/utils"
)

func Register(mux *http.ServeMux, log *slog.Logger, srv bins.Service) {
	mux.Handle("GET /bin/{id}", getBin(srv, log))
	mux.Handle("POST /bin", createBin(srv, log))
}

func getBin(srv bins.Service, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		id := r.PathValue("id")
		if id == "" {
			utils.MustRespondError(w, http.StatusBadRequest, "id not provided")
			return
		}

		log = log.With(slog.String("id", id))

		bin, err := srv.Get(ctx, id)
		if err != nil {
			if errors.Is(err, binErr.ErrNotFound) {
				utils.MustRespondError(w, http.StatusNotFound, err.Error())
			} else {
				log.Error("couldn't get bin", slog.String("err", err.Error()))
				utils.MustRespondError(w, http.StatusInternalServerError, "couldn't get bin")
			}
			return
		}

		utils.MustRespondJSON(w, http.StatusOK, bin)
	}
}

func createBin(srv bins.Service, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		newBin, err := utils.ReadJSON[models.NewBin](r)
		if err != nil {
			utils.MustRespondError(w, http.StatusBadRequest, "couldn't parse bin")
			return
		}

		bin, err := srv.Create(ctx, newBin)
		if err != nil {
			log.Error("couldn't create bin")
			utils.MustRespondError(w, http.StatusInternalServerError, "couldn't create bin")
			return
		}

		utils.MustRespondJSON(w, http.StatusCreated, bin)
	}
}
