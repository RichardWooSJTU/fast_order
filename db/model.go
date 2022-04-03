package db

import "time"

type Order struct {
	ID        int64
	Title     string
	Info      string
	Expired   time.Time
	Threshold int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
