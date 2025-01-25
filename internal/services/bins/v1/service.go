package v1

import (
	"context"
	"log/slog"

	"github.com/rmntim/xbin/internal/services/bins"
	"github.com/rmntim/xbin/internal/services/bins/models"
)

type Service struct {
	log *slog.Logger
}

func NewService(ctx context.Context, log *slog.Logger) bins.Service {
	return &Service{log: log}
}

func (s *Service) Get(ctx context.Context, id string) (models.Bin, error) {
	s.log.Debug("bins.Get", slog.String("id", id))
	return models.Bin{}, nil
}

func (s *Service) Create(ctx context.Context, content string) (models.Bin, error) {
	s.log.Debug("bins.Create", slog.String("content", content))
	return models.Bin{}, nil
}
