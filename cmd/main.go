package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/rmntim/xbin/internal/httpserver"
	"github.com/rmntim/xbin/internal/repo/bins/sqlite"
	bins "github.com/rmntim/xbin/internal/services/bins/v1"
)

type Config struct {
	Port        uint16
	Env         loggerEnv
	StoragePath string
	TursoCfg    *sqlite.TursoReplicaConfig
}

func configure() (Config, error) {
	var cfg Config

	flag.StringVar(&cfg.StoragePath, "storagePath", "./data/bins.db", "path to storage")
	port := flag.Uint64("port", 8080, "port to listen on")
	env := flag.String("env", string(envProd), "program environment (values: dev, prod)")

	flag.Parse()

	if *env != string(envDev) && *env != string(envProd) {
		return Config{}, fmt.Errorf("invalid env value: %+v", *env)
	}

	portEnv, okPort := os.LookupEnv("PORT")
	if okPort {
		var err error

		*port, err = strconv.ParseUint(portEnv, 10, 0)
		if err != nil {
			return Config{}, fmt.Errorf("invalid port value: %w", err)
		}
	}

	if *port > math.MaxUint16 {
		return Config{}, fmt.Errorf("invalid port value: %+v", *port)
	}

	cfg.Port = uint16(*port)
	cfg.Env = loggerEnv(*env)

	tursoUrl, okUrl := os.LookupEnv("TURSO_URL")
	token, okToken := os.LookupEnv("TURSO_TOKEN")
	ok := okUrl && okToken

	if ok {
		cfg.TursoCfg = &sqlite.TursoReplicaConfig{
			AuthToken: token,
			URL:       tursoUrl,
		}
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

	db, err := sqlite.NewRepository(log, cfg.StoragePath, cfg.TursoCfg)
	if err != nil {
		log.Error("error creating database connection", slog.String("err", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	binSrv := bins.NewService(log, db)

	srv := httpserver.NewServer(fmt.Sprintf(":%d", cfg.Port), log, binSrv)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info("listening", slog.String("address", srv.Addr))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("error listening", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("got interruption signal")
	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Error("server shutdown returned an err", slog.String("err", err.Error()))
	}

	log.Info("final")
}
