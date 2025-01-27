package bins

import (
	"context"
)

type Repository interface {
	Get(ctx context.Context, id string) (BinStorage, error)
	Create(ctx context.Context, newBin BinStorage) (string, error)
	Close() error
}
