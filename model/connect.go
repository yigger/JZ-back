package model

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"github.com/go-redis/redis"
	"github.com/yigger/JZ-back/conf"
	"fmt"
)

var db *gorm.DB
var redisCli *redis.Client

func DB() *gorm.DB {
	if db == nil {
		path := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", conf.Conf.DbConfig.Username, conf.Conf.DbConfig.Password, conf.Conf.DbConfig.Database)
		var err error
		db, err = gorm.Open("mysql", path)
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
			Addr:     conf.Conf.RedisConfig.Host,
			Password: conf.Conf.RedisConfig.Password, // no password set
			DB:       conf.Conf.RedisConfig.Db,  // use default DB
		})
	}

	return redisCli
}
