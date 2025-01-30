package routes

import (
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/rmntim/xbin/internal/services/bins"
	binErr "github.com/rmntim/xbin/internal/services/bins/errors"
	"github.com/rmntim/xbin/internal/services/bins/models"
	"github.com/rmntim/xbin/internal/utils"
)

var funcs = template.FuncMap{
	"unix": func(t time.Time) int64 {
		return t.UnixMilli()
	},
}

func Register(mux *http.ServeMux, log *slog.Logger, srv bins.Service) {
	mux.Handle("GET /bin/{id}", getBin(srv, log))
	mux.Handle("POST /bin", createBin(srv, log))
}

func getBin(srv bins.Service, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := r.PathValue("id")
		if id == "" {
			utils.MustRespondError(w, http.StatusBadRequest, "id not provided")
			return
		}

		log = log.With(slog.String("id", id))

		bin, err := srv.Get(ctx, id)
		if err != nil {
			if errors.Is(err, binErr.ErrNotFound) || errors.Is(err, binErr.ErrExpired) {
				utils.MustRespondError(w, http.StatusNotFound, "bin not found")
			} else {
				log.Error("couldn't get bin", slog.String("err", err.Error()))
				utils.MustRespondError(w, http.StatusInternalServerError, "couldn't get bin")
			}
			return
		}

		tmpl, err := template.New("bin.tmpl.html").Funcs(funcs).ParseFiles("./static/bin.tmpl.html")
		if err != nil {
			log.Error("couldn't parse template", slog.String("err", err.Error()))
			utils.MustRespondError(w, http.StatusInternalServerError, "internal error")
			return
		}

		if err = tmpl.Execute(w, bin); err != nil {
			log.Error("couldn't execute template", slog.String("err", err.Error()))
			utils.MustRespondError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func createBin(srv bins.Service, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		newBin, err := utils.ReadJSON[models.NewBinRequest](r)
		if err != nil {
			utils.MustRespondError(w, http.StatusBadRequest, "couldn't parse bin")
			return
		}

		if newBin.Expiration.Duration == 0*time.Second {
			utils.MustRespondError(w, http.StatusBadRequest, "duration cannot be 0")
			return
		}

		bin, err := srv.Create(ctx, newBin)
		if err != nil {
			log.Error("couldn't create bin", slog.String("err", err.Error()))
			utils.MustRespondError(w, http.StatusInternalServerError, "couldn't create bin")
			return
		}

		utils.MustRespondJSON(w, http.StatusCreated, bin)
	}
}
