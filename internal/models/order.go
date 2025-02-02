package models

import "time"

type Position struct {
	ID        int64     `json:"id" db:"id"`
	Price     int64     `json:"price" db:"price"`
	Name      string    `json:"name" db:"name"`
	UpdatedAt time.Time `json:"updated_at" db:"updatedAt"`
}
