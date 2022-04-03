package user

import "github.com/gin-gonic/gin"

func Route(r *gin.Engine) {
	r.GET("/user/my_order", GetMyOrder)
	r.GET("/user/my_member", GetMyJoinedOrder)
}
