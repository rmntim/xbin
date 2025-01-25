package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/rmntim/xbin/internal/httpserver"
	bins "github.com/rmntim/xbin/internal/services/bins/v1"
)

func main() {
	ctx := context.Background()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	binSrv := bins.NewService(ctx, log)

	srv := httpserver.NewServer("localhost:8080", log, binSrv)

	log.Info("listening", slog.String("address", srv.Addr))

	if err := srv.ListenAndServe(); err != nil {
		log.Error("error listening", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
