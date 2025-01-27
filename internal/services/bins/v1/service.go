package v1

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"time"

	repo "github.com/rmntim/xbin/internal/repo/bins"
	"github.com/rmntim/xbin/internal/services/bins"
	svcErr "github.com/rmntim/xbin/internal/services/bins/errors"
	"github.com/rmntim/xbin/internal/services/bins/models"
)

type Service struct {
	log  *slog.Logger
	repo repo.Repository
}

func NewService(log *slog.Logger, repo repo.Repository) bins.Service {
	return &Service{log: log, repo: repo}
}

func (s *Service) Get(ctx context.Context, id string) (models.Bin, error) {
	log := s.log.With(slog.String("id", id))

	bin, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, svcErr.ErrNotFound) {
			return models.Bin{}, err
		}

		log.Error("could not get bin", slog.String("err", err.Error()))
		return models.Bin{}, fmt.Errorf("could not get bin with id %s", id)
	}

	log.Debug("bin found")

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

	id, err := s.repo.Create(ctx, storageBin)
	if err != nil {
		s.log.Error("could not create bin", slog.String("err", err.Error()))
		return models.NewBinResponse{}, fmt.Errorf("could not create bin: %w", err)
	}

	bin := models.NewBinResponse{
		Id:  id,
		URL: fmt.Sprintf("/bin/%s", id),
	}

	s.log.Debug("successfully created bin", slog.String("id", bin.Id))

	return bin, nil
}

func svcToStorageBin(bin models.NewBinRequest) (repo.BinStorage, error) {
	id, err := generateRandomString(8)
	if err != nil {
		return repo.BinStorage{}, fmt.Errorf("error creating uuid: %w", err)
	}

	now := time.Now()

	return repo.BinStorage{
		Id:        id,
		Content:   bin.Content,
		CreatedAt: now,
		ExpiresAt: now.Add(bin.Expiration.Duration),
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
