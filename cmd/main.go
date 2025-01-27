package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/rmntim/xbin/internal/httpserver"
	"github.com/rmntim/xbin/internal/repo/bins/sqlite"
	bins "github.com/rmntim/xbin/internal/services/bins/v1"
)

type Config struct {
	Port        uint16
	Env         loggerEnv
	StoragePath string
}

func configure() (Config, error) {
	var cfg Config

	port, ok := os.LookupEnv("XBIN_PORT")
	if !ok {
		cfg.Port = 8080
	} else {
		parsedPort, err := strconv.Atoi(port)
		if err != nil {
			return Config{}, fmt.Errorf("bad port value %s", port)
		}

		cfg.Port = uint16(parsedPort)
	}

	env, ok := os.LookupEnv("XBIN_ENV")
	if !ok {
		cfg.Env = envProd
	} else {
		if env != string(envProd) && env != string(envDev) {
			return Config{}, fmt.Errorf("invalid env")
		}

		cfg.Env = loggerEnv(env)
	}

	storagePath, ok := os.LookupEnv("XBIN_STORAGE_PATH")
	if !ok {
		cfg.StoragePath = "./bins.db"
	} else {
		cfg.StoragePath = storagePath
	}

	return cfg, nil
}

type loggerEnv string

const (
	envDev  loggerEnv = "dev"
	envProd loggerEnv = "prod"
)

func setupLogger(env loggerEnv) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}

func main() {
	cfg, err := configure()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	log := setupLogger(cfg.Env)

	db, err := sqlite.NewRepository(log, cfg.StoragePath)
	if err != nil {
		log.Error("error creating database connection", slog.String("err", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	binSrv := bins.NewService(log, db)

	srv := httpserver.NewServer(fmt.Sprintf(":%d", cfg.Port), log, binSrv)

	log.Info("listening", slog.String("address", srv.Addr))

	if err := srv.ListenAndServe(); err != nil {
		log.Error("error listening", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
