package create_order

import "github.com/gin-gonic/gin"

func Route(r *gin.Engine) {
	r.POST("/create_order", Create)
}
