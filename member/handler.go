package member

import (
	"fast_order/db"
	order2 "fast_order/order"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type PostInfoReq struct {
	OrderID  int64  `json:"order_id"`
	UserID   int64  `json:"user_id"`
	Building int    `json:"building"`
	Room     string `json:"room"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Value    string `json:"value"`
}

func PostInfo(ctx *gin.Context) {
	var param PostInfoReq
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//如果过期则不能提交修改
	order := db.Order{
		ID: param.OrderID,
	}
	err := db.DBEngine.First(&order, param.OrderID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("db error %v", err)
		return
	}
	if order.Expired.Unix() < time.Now().Unix() {
		ctx.JSON(http.StatusOK, gin.H{"status": "200", "data": order2.NoticeExpired})
		return
	}

	member := db.Member{}
	err = db.DBEngine.First(&member, "order_id=? and user_id=?", param.OrderID, param.UserID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("db error")
		return
	}
	member.Building = param.Building
	member.Room = param.Room
	member.Name = param.Name
	member.Phone = param.Phone
	member.Value = param.Value
	err = db.DBEngine.Save(&member).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("db error %v", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "200"})
}
