package order

import (
	"fast_order/db"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type CreateOrderReq struct {
	UserID    int64
	Title     string
	Info      string
	Expired   time.Time
	Threshold int64
}

func Create(ctx *gin.Context) {
	var param CreateOrderReq
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newOrder := db.Order{
		Title:     param.Title,
		Info:      param.Info,
		Expired:   param.Expired,
		Threshold: param.Threshold,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	//order表当中新建记录
	err := db.DBEngine.Create(&newOrder).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatalf("db error")
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "200"})
}
