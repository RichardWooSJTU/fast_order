package order

import "github.com/gin-gonic/gin"

func Route(r *gin.Engine) {
	r.POST("/order", Create)
	r.GET("/order", Get)
	r.GET("/order/summary", GetSummaryAndCreateMember)
	r.PUT("/order", Update)
}
