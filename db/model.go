package db

import "time"

type Order struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Info      string    `json:"info"`
	Expired   time.Time `json:"expired"`
	Threshold int64     `json:"threshold"`
	OwnerID   int64     `json:"owner_id"`
	CreateAt  time.Time `json:"create_at"`
	UpdateAt  time.Time `json:"update_at"`
}

type Member struct {
	ID       int64     `json:"id"`
	OrderID  int64     `json:"order_id"`
	UserID   int64     `json:"user_id"`
	Building int       `json:"building"`
	Room     string    `json:"room"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Value    string    `json:"value"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}
