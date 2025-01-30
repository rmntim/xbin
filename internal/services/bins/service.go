package bins

import (
	"context"

	"github.com/rmntim/xbin/internal/services/bins/models"
)

type Service interface {
	GetBySlug(ctx context.Context, slug string) (models.Bin, error)
	Create(ctx context.Context, content models.NewBinRequest) (models.NewBinResponse, error)
}
