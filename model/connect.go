package model

import (
	"fmt"
	"time"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"github.com/go-redis/redis"

	"github.com/yigger/JZ-back/conf"	
)

var db *gorm.DB
var redisCli *redis.Client

type CommonModel struct {
	ID        	uint	 		`gorm:"primary_key" json:"id"`
	CreatedAt 	time.Time		`json:"created_at"`
	UpdatedAt 	time.Time		`json:"updated_at"`
}

func ConnectDB() *gorm.DB {
	if db == nil {
		path := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=true&loc=Local", conf.Conf.DbConfig.Username, conf.Conf.DbConfig.Password, conf.Conf.DbConfig.Database)
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

func ConnectRedis() (*redis.Client) {
	if redisCli == nil {
		redisCli = redis.NewClient(&redis.Options{
			Addr:     conf.Conf.RedisConfig.Host,
			Password: conf.Conf.RedisConfig.Password, // no password set
			DB:       conf.Conf.RedisConfig.Db,  // use default DB
		})
	}

	return redisCli
}
