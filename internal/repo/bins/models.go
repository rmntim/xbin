package bins

import "time"

type BinStorage struct {
	Id        string
	Content   string
	CreatedAt time.Time
	ExpiresAt time.Time
	Slug      string
}
