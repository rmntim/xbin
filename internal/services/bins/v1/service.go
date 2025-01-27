package v1

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	repo "github.com/rmntim/xbin/internal/repo/bins"
	"github.com/rmntim/xbin/internal/services/bins"
	svcErr "github.com/rmntim/xbin/internal/services/bins/errors"
	"github.com/rmntim/xbin/internal/services/bins/models"
)

type Service struct {
	log  *slog.Logger
	repo repo.Repository
}

func NewService(ctx context.Context, log *slog.Logger, repo repo.Repository) bins.Service {
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

	return bin, nil
}

func (s *Service) Create(ctx context.Context, newBin models.NewBin) (models.Bin, error) {
	bin, err := s.repo.Create(ctx, newBin)
	if err != nil {
		s.log.Error("could not create bin", slog.String("err", err.Error()))
		return models.Bin{}, fmt.Errorf("could not create bin")
	}

	s.log.Debug("successfully created bin", slog.String("id", bin.Id))

	return bin, nil
}
