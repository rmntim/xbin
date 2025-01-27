package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/mattn/go-sqlite3"
	repo "github.com/rmntim/xbin/internal/repo/bins"
	svcErr "github.com/rmntim/xbin/internal/services/bins/errors"
	cr "github.com/robfig/cron/v3"
)

type repository struct {
	log  *slog.Logger
	db   *sql.DB
	cron *cr.Cron
}

func NewRepository(log *slog.Logger, url string) (repo.Repository, error) {
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		return nil, err
	}

	r := &repository{
		log:  log,
		db:   db,
		cron: cr.New(),
	}

	_, err = r.cron.AddFunc("*/10 * * * *", r.deleteOldRows)
	if err != nil {
		return nil, err
	}

	r.deleteOldRows()

	return r, nil
}

func (r *repository) deleteOldRows() {
	_, err := r.db.Exec("DELETE FROM bins WHERE expires_at < ?", time.Now())
	if err != nil {
		r.log.Error("could not delete old bins", slog.String("err", err.Error()))
	}
}

func (r *repository) Close() error {
	return r.db.Close()
}

func (r *repository) Get(ctx context.Context, id string) (repo.BinStorage, error) {
	var bin repo.BinStorage
	err := r.db.QueryRowContext(ctx, "SELECT id, content, created_at, expires_at FROM bins WHERE id = ?", id).Scan(&bin.Id, &bin.Content, &bin.CreatedAt, &bin.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repo.BinStorage{}, svcErr.ErrNotFound
		}
		return repo.BinStorage{}, fmt.Errorf("could not find bin with id %s: %w", id, err)
	}

	return bin, nil
}

func (r *repository) Create(ctx context.Context, bin repo.BinStorage) (string, error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO bins (id, content, created_at, expires_at) VALUES (?, ?, ?, ?) RETURNING id")
	if err != nil {
		return "", fmt.Errorf("could not prepare statement: %w", err)
	}
	defer stmt.Close()

	var insertedId string
	err = stmt.QueryRowContext(ctx, bin.Id, bin.Content, bin.CreatedAt, bin.ExpiresAt).Scan(&insertedId)
	if err != nil {
		return "", fmt.Errorf("could not create bin: %w", err)
	}

	return insertedId, nil
}
