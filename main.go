package main

import (
	"github.com/yigger/JZ-back/controller"
	"github.com/yigger/JZ-back/model"
	"github.com/yigger/JZ-back/conf"
	"fmt"
)

func init() {
	conf.LoadConf()
}

func main() {
	model.ConnectDB()
	model.ConnectRedis()
	e := controller.EchoNew()
	path := fmt.Sprintf("%s:%s", conf.Conf.Host, conf.Conf.Port)
	e.Logger.Fatal(e.Start(path))
}
