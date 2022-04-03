package member

import "github.com/gin-gonic/gin"

func Route(r *gin.Engine) {
	r.POST("/member", PostInfo)
}
