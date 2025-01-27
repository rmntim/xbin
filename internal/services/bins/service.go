package bins

import (
	"context"

	"github.com/rmntim/xbin/internal/services/bins/models"
)

type Service interface {
	Get(ctx context.Context, id string) (models.Bin, error)
	Create(ctx context.Context, content models.NewBinRequest) (models.NewBinResponse, error)
}
