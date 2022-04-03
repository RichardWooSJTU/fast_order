package order

import (
	"fast_order/db"
)

type CreateOrderReq struct {
	UserID    int64  `json:"user_id"`
	Title     string `json:"title"`
	Info      string `json:"info"`
	Expired   int64  `json:"expired"`
	Threshold int64  `json:"threshold"`
}

type GetOrderReq struct {
	OrderID int64 `form:"order_id"`
	UserID  int64 `form:"user_id"`
}

type GetOrderRes struct {
	db.Order
	Members []db.Member `json:"members"`
}

type GetSummaryReq struct {
	OrderID int64 `form:"order_id"`
	UserID  int64 `form:"user_id"`
}

type GetSummaryRes struct {
	db.Order
	Member db.Member `json:"member"`
}

type UpdateOrderReq struct {
	OrderID   int64  `form:"order_id"`
	UserID    int64  `form:"user_id"`
	Title     string `form:"title"`
	Info      string `form:"info"`
	Expired   int64  `form:"expired"`
	Threshold int64  `form:"threshold"`
}
