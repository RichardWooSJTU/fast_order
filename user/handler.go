package user

import (
	"fast_order/db"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type UserReq struct {
	UserID int64 `form:"user_id"`
}

func GetMyOrder(ctx *gin.Context) {
	var param UserReq
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var orders []db.Order
	err := db.DBEngine.Find(&orders, "owner_id=?", param.UserID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("db error %v", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "200", "data": orders})
}

func GetMyJoinedOrder(ctx *gin.Context) {
	var param UserReq
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
