package conf

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
)

// The JieZhang Conf
var Conf *Configuration

type Configuration struct {
	// Environment: development, production, test
	Environment		string  	`json:"environment"`
	Host			string		`json:"host"`
	Port			string		`json:"port"`

	// Database Conf Struct
	DbConfig		DatabaseConf
	// Redis Conf Struct
	RedisConfig		RedisConf
}

type DatabaseConf struct {
	Host 			string 		`json:"host"`
	Username		string		`json:"username"`
	Password 	 	string		`json:"password"`
	Port 			string		`json:"port"`
	Database		string		`json:"database"`
}

type RedisConf struct {
	Host 			string 		`json:"host"`
	Password 	 	string		`json:"password"`
	Port 			string		`json:"port"`
	Db				int 		`json:"db"`
}

var (
	jzConfigPath 	   	= 	"conf/conf.d/jz.json"
	databaseConfigPath 	= 	"conf/conf.d/database.json"
	redisConfigPath	   	=	"conf/conf.d/redis.json"
)

func Development() bool {
	return Conf.Environment == "development"
}

func LoadConf() {
	loadCommonConf()
	loadDatabaseConf()
	loadRedisConf()
}

func loadCommonConf() {
	bytes, err := ioutil.ReadFile(jzConfigPath)
	if nil != err {
		panic(err)
	}

	if err = json.Unmarshal(bytes, &Conf); nil != err {
		panic(err)
	}
}

func loadDatabaseConf() {
	bytes, err := ioutil.ReadFile(databaseConfigPath)
	if nil != err {
		fmt.Println("配置文件加载失败")
	}

	conf := make(map[string]DatabaseConf, 0)
	json.Unmarshal(bytes, &conf)

	dbConf := conf[Conf.Environment]
	Conf.DbConfig = dbConf
}

func loadRedisConf() {
	bytes, err := ioutil.ReadFile(redisConfigPath)
	if nil != err {
		fmt.Println("redis 配置文件加载失败")
	}
	conf := &RedisConf{}
	json.Unmarshal(bytes, &conf)

	Conf.RedisConfig = *conf
}