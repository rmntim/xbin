package v1

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	repo "github.com/rmntim/xbin/internal/repo/bins"
	"github.com/rmntim/xbin/internal/services/bins"
	svcErr "github.com/rmntim/xbin/internal/services/bins/errors"
	"github.com/rmntim/xbin/internal/services/bins/models"
)

const slugLength = 8

type Service struct {
	log  *slog.Logger
	repo repo.Repository
}

func NewService(log *slog.Logger, repo repo.Repository) bins.Service {
	return &Service{log: log, repo: repo}
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (models.Bin, error) {
	log := s.log.With(slog.String("slug", slug))

	bin, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Bin{}, svcErr.ErrNotFound
		}

		log.Error("could not get bin", slog.String("err", err.Error()))
		return models.Bin{}, fmt.Errorf("could not get bin with slug %s", slug)
	}

	log.Debug("bin found")

	if bin.ExpiresAt.Before(time.Now()) {
		return models.Bin{}, svcErr.ErrExpired
	}

	return models.Bin{
		Id:        bin.Id,
		Content:   bin.Content,
		CreatedAt: bin.CreatedAt,
		ExpiresAt: bin.ExpiresAt,
	}, nil
}

func (s *Service) Create(ctx context.Context, newBin models.NewBinRequest) (models.NewBinResponse, error) {
	storageBin, err := svcToStorageBin(newBin)
	if err != nil {
		return models.NewBinResponse{}, fmt.Errorf("could not convert service to storage bin: %w", err)
	}

	slug, err := s.repo.Create(ctx, storageBin)
	if err != nil {
		s.log.Error("could not create bin", slog.String("err", err.Error()))
		return models.NewBinResponse{}, fmt.Errorf("could not create bin: %w", err)
	}

	bin := models.NewBinResponse{
		URL: fmt.Sprintf("/bin/%s", slug),
	}

	s.log.Debug("successfully created bin", slog.String("slug", slug))

	return bin, nil
}

func svcToStorageBin(bin models.NewBinRequest) (repo.BinStorage, error) {
	slug, err := generateRandomString(slugLength)
	if err != nil {
		return repo.BinStorage{}, fmt.Errorf("error creating slug: %w", err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return repo.BinStorage{}, fmt.Errorf("error creating uuid: %w", err)
	}

	now := time.Now()

	return repo.BinStorage{
		Id:        id.String(),
		Content:   bin.Content,
		CreatedAt: now,
		ExpiresAt: now.Add(bin.Expiration.Duration),
		Slug:      slug,
	}, nil
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
