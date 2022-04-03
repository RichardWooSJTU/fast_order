package order

import (
	"fast_order/db"
	"fast_order/lock"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

const (
	NoticeFull    = 0
	NoticeExpired = 1
)

func Create(ctx *gin.Context) {
	var param CreateOrderReq
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newOrder := db.Order{
		Title:     param.Title,
		Info:      param.Info,
		OwnerID:   param.UserID,
		Expired:   time.Unix(param.Expired, 0),
		Threshold: param.Threshold,
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
	}
	//order表当中新建记录
	err := db.DBEngine.Create(&newOrder).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("db error %v", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "200"})
}

func Get(ctx *gin.Context) {
	//仅发单人可见
	var param GetOrderReq
	if err := ctx.ShouldBindQuery(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order := db.Order{}
	err := db.DBEngine.First(&order, param.OrderID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("db error %v", err)
		return
	}
	if order.OwnerID != param.UserID {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "you are not the owner of this order"})
		log.Printf("owner error")
		return
	}

	var members []db.Member
	err = db.DBEngine.Find(&members, "order_id=?", order.ID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("db error %v", err)
		return
	}
	response := GetOrderRes{
		Order:   order,
		Members: members,
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{"status": "200", "data": response})
}

func Update(ctx *gin.Context) {
	//仅发单人可见
	var param UpdateOrderReq
	if err := ctx.ShouldBindQuery(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order := db.Order{}
	err := db.DBEngine.First(&order, param.OrderID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("db error %v", err)
		return
	}
	if order.OwnerID != param.UserID {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "you are not the owner of this order"})
		log.Printf("owner error")
		return
	}
	order.Title = param.Title
	order.Info = param.Info
	order.Expired = time.Unix(param.Expired, 0)
	order.Threshold = param.Threshold
	order.UpdateAt = time.Now()

	err = db.DBEngine.Save(&order).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("db error %v", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "200"})
}

func GetSummaryAndCreateMember(ctx *gin.Context) {
	//超过截至时间不能再加入 纵然还没到阈值
	var param GetOrderReq
	if err := ctx.ShouldBindQuery(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(param)
	//先看自己在不在里面 在的话直接返回 否则 当前order的人数如果没超 把自己加进去 否则返回人数已满
	order := db.Order{
		ID: param.OrderID,
	}
	err := db.DBEngine.First(&order, param.OrderID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("db error %v", err)
		return
	}

	var memberId int64
	err = db.DBEngine.Model(&db.Member{}).Select("id").Where("order_id=? and user_id=?", param.OrderID, param.UserID).First(&memberId).Error
	if err == gorm.ErrRecordNotFound {
		//不在
		//先判断是否过期 如果过期 不再新建
		if order.Expired.Unix() < time.Now().Unix() {
			ctx.JSON(http.StatusOK, gin.H{"status": "200", "data": NoticeExpired})
			return
		}
		//看看是否可以新建
		//上锁
		lock.Lock.Lock()
		var cnt int64
		err = db.DBEngine.Model(&db.Member{}).Select("count(id)").Where("order_id=?", param.OrderID).Find(&cnt).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Printf("db error %v", err)
			return
		}

		if cnt < order.Threshold {
			//可以新建
			fmt.Printf("cnt %v\n", cnt)
			member := db.Member{
				OrderID:  param.OrderID,
				UserID:   param.UserID,
				Building: 0,
				Room:     "",
				Name:     "",
				Phone:    "",
				Value:    "",
				CreateAt: time.Now(),
				UpdateAt: time.Now(),
			}
			err = db.DBEngine.Create(&member).Error
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				log.Printf("db error %v", err)
				return
			}
			//解锁
			lock.Lock.Unlock()
			//返回信息
			res := GetSummaryRes{
				Order:  order,
				Member: member,
			}
			ctx.JSON(http.StatusOK, gin.H{"status": "200", "data": res})
			return
		} else {
			//不行 下次再来
			ctx.JSON(http.StatusOK, gin.H{"status": "200", "data": NoticeFull})
			return
		}
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("db error %v", err)
		return
	} else {
		//数据库里面有 直接返回
		member := db.Member{
			ID: memberId,
		}
		err = db.DBEngine.First(&member).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Printf("db error %v", err)
			return
		}
		//返回信息
		res := GetSummaryRes{
			Order:  order,
			Member: member,
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "200", "data": res})
		return
	}
}
