package model

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"github.com/go-redis/redis"
)

var db *gorm.DB
var redisCli *redis.Client

func DB() *gorm.DB {
	if db == nil {
		var err error
		db, err = gorm.Open("mysql", "root:root@/ljt_development?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			panic(err)
		}

		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
	}
	return db
}

func Redis() (*redis.Client) {
	if redisCli == nil {
		redisCli = redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}

	return redisCli
}
