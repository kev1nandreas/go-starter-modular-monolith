package dto

import (
	"time"

	"github.com/google/uuid"
)

type BookResponse struct {
	ID        uuid.UUID `json:"uuid"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
