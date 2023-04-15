package model

import (
	"fmt"
	"jiuxia/chatDemo/conf"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func Database() {
	path := conf.Path
	fmt.Println(path)
	db, err := gorm.Open("mysql", path)
	if err != nil {
		panic(err)
	}
	fmt.Println("Mysql数据库连接成功")
	db.SingularTable(true)
	//使表名不加s
	db.DB().SetMaxIdleConns(20)
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	db.DB().SetMaxOpenConns(100)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	db.DB().SetConnMaxLifetime(time.Second * 30)
	// SetConnMaxLifetime 设置了连接可复用的最大时间
	DB = db
	DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&User{}).AutoMigrate(&Friend{})
}
