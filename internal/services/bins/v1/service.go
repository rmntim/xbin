package v1

import (
	"context"
	"log/slog"
	"time"

	repo "github.com/rmntim/xbin/internal/repo/bins"
	"github.com/rmntim/xbin/internal/services/bins"
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
	s.log.Debug("bins.Get", slog.String("id", id))
	return models.Bin{
		Id:        id,
		Content:   "foo",
		CreatedAt: time.Now(),
	}, nil
}

func (s *Service) Create(ctx context.Context, newBin models.NewBin) (models.Bin, error) {
	s.log.Debug("bins.Create", slog.String("content", newBin.Content))
	return models.Bin{
		Id:        "id",
		Content:   newBin.Content,
		CreatedAt: time.Now(),
	}, nil
}
