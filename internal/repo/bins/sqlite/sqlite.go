package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	repo "github.com/rmntim/xbin/internal/repo/bins"
	"github.com/tursodatabase/go-libsql"
)

type repository struct {
	log *slog.Logger
	db  *sql.DB
}

type TursoReplicaConfig struct {
	URL       string
	AuthToken string
}

func NewRepository(log *slog.Logger, url string, tursoConfig *TursoReplicaConfig) (repo.Repository, error) {
	var db *sql.DB

	if tursoConfig != nil {
		connector, err := libsql.NewEmbeddedReplicaConnector(url, tursoConfig.URL, libsql.WithAuthToken(tursoConfig.AuthToken))
		if err != nil {
			return nil, fmt.Errorf("error creating connector: %w", err)
		}

		log.Debug("connecting to turso")

		db = sql.OpenDB(connector)
	} else {
		log.Debug("connecting to local db")

		dbLocal, err := sql.Open("libsql", "file:"+url)
		if err != nil {
			return nil, fmt.Errorf("error creating local db: %w", err)
		}

		db = dbLocal
	}

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return nil, fmt.Errorf("error creating migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "sqlite", driver)
	if err != nil {
		return nil, fmt.Errorf("error creating migration: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("error running migrations: %w", err)
	}

	r := &repository{
		log: log,
		db:  db,
	}

	return r, nil
}

func (r *repository) Close() error {
	return r.db.Close()
}

func (r *repository) GetBySlug(ctx context.Context, slug string) (repo.BinStorage, error) {
	var bin repo.BinStorage
	err := r.db.QueryRowContext(ctx, "SELECT id, content, created_at, expires_at, slug FROM bins WHERE slug = ?", slug).Scan(&bin.Id, &bin.Content, &bin.CreatedAt, &bin.ExpiresAt, &bin.Slug)
	if err != nil {
		return repo.BinStorage{}, fmt.Errorf("could not find bin with slug %s: %w", slug, err)
	}

	return bin, nil
}

func (r *repository) Create(ctx context.Context, bin repo.BinStorage) (string, error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO bins (id, content, created_at, expires_at, slug) VALUES (?, ?, ?, ?, ?) RETURNING slug")
	if err != nil {
		return "", fmt.Errorf("could not prepare statement: %w", err)
	}
	defer stmt.Close()

	var insertedSlug string
	err = stmt.QueryRowContext(ctx, bin.Id, bin.Content, bin.CreatedAt, bin.ExpiresAt, bin.Slug).Scan(&insertedSlug)
	if err != nil {
		return "", fmt.Errorf("could not create bin: %w", err)
	}

	return insertedSlug, nil
}
