package models

import "time"

type Bin struct {
	Id        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type NewBin struct {
	Content string `json:"content"`
}
