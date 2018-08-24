package comicdao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func conectDB() {
	if db != nil {
		return
	}
	var err error = nil
	db, err = gorm.Open("mysql", "root:1124@(127.0.0.1)/comic?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}
