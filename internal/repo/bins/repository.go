package bins

import (
	"context"

	"github.com/rmntim/xbin/internal/services/bins/models"
)

type Repository interface {
	Get(ctx context.Context, id string) (models.Bin, error)
	Create(ctx context.Context, newBin models.NewBin) (models.Bin, error)
	Close() error
}
