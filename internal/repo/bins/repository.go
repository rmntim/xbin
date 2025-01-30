package bins

import (
	"context"
)

type Repository interface {
	GetBySlug(ctx context.Context, slug string) (BinStorage, error)
	Create(ctx context.Context, newBin BinStorage) (string, error)
	Close() error
}
