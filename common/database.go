package common

import (
	"fmt"
	"ginpro/model"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driveName := "mysql"
	host := "127.0.0.1"
	port := "3306"
	database := "ginpro"
	username := "root"
	password := "0614"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(driveName, args)
	if err != nil {
		panic("failed to connect database,err:" + err.Error())
	}
	db.AutoMigrate(&model.User{})
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
