package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBEngine *gorm.DB

func Init() error {
	var err error
	dsn := "wfs:daaitangyan0219.@tcp(127.0.0.1:3306)/fast_order?charset=utf8mb4&parseTime=True&loc=Local"
	DBEngine, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}
