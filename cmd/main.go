package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/rmntim/xbin/internal/httpserver"
	"github.com/rmntim/xbin/internal/repo/bins/sqlite"
	bins "github.com/rmntim/xbin/internal/services/bins/v1"
)

func main() {
	ctx := context.Background()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	db, err := sqlite.NewRepository(log, "./bins.db")
	if err != nil {
		log.Error("error creating database connection", slog.String("err", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	binSrv := bins.NewService(ctx, log, db)

	srv := httpserver.NewServer("localhost:8080", log, binSrv)

	log.Info("listening", slog.String("address", srv.Addr))

	if err := srv.ListenAndServe(); err != nil {
		log.Error("error listening", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
