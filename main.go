package main

import (
	"fast_order/db"
	"fast_order/member"
	"fast_order/order"
	"github.com/gin-gonic/gin"
)

func main() {
	err := db.Init()
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	order.Route(r)
	member.Route(r)

	r.Run(":6789")
}
