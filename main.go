package main

import (
	"fmt"

	"github.com/yigger/JZ-back/controller"
	"github.com/yigger/JZ-back/model"
	"github.com/yigger/JZ-back/conf"
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
