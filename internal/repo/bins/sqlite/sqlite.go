package sqlite

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	repo "github.com/rmntim/xbin/internal/repo/bins"
	svcErr "github.com/rmntim/xbin/internal/services/bins/errors"
	"github.com/rmntim/xbin/internal/services/bins/models"
)

type repository struct {
	log *slog.Logger
	db  *sql.DB
}

func NewRepository(log *slog.Logger, url string) (repo.Repository, error) {
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		return nil, err
	}

	return &repository{
		log: log,
		db:  db,
	}, nil
}

func (r *repository) Close() error {
	return r.db.Close()
}

func (r *repository) Get(ctx context.Context, id string) (models.Bin, error) {
	var bin models.Bin
	err := r.db.QueryRowContext(ctx, "SELECT id, content, created_at FROM bins WHERE id = ?", id).Scan(&bin.Id, &bin.Content, &bin.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Bin{}, svcErr.ErrNotFound
		}
		return models.Bin{}, fmt.Errorf("could not find bin with id %s: %w", id, err)
	}

	return bin, nil
}

func (r *repository) Create(ctx context.Context, newBin models.NewBin) (models.Bin, error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO bins (id, content) VALUES (?, ?) RETURNING id, content, created_at")
	if err != nil {
		return models.Bin{}, fmt.Errorf("could not prepare statement: %w", err)
	}
	defer stmt.Close()

	id, err := generateRandomString(8)
	if err != nil {
		return models.Bin{}, fmt.Errorf("error creating uuid: %w", err)
	}

	var bin models.Bin
	err = stmt.QueryRowContext(ctx, id, newBin.Content).Scan(&bin.Id, &bin.Content, &bin.CreatedAt)
	if err != nil {
		return models.Bin{}, fmt.Errorf("could not create bin: %w", err)
	}

	return bin, nil
}

func generateRandomString(length uint) (string, error) {
	// Calculate the number of bytes needed
	byteLength := (length * 6) / 8
	if (length*6)%8 != 0 {
		byteLength++
	}

	// Generate random bytes
	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode the bytes to a base64 URL string
	randomString := base64.URLEncoding.EncodeToString(randomBytes)

	// Trim the string to the desired length
	return randomString[:length], nil
}
