package main

import (
	"github.com/yigger/JZ-back/controllers"
	"github.com/yigger/JZ-back/conf"
	"fmt"
)

func init() {
	conf.LoadConf()
}

func main() {
	e := controllers.EchoNew()
	path := fmt.Sprintf("%s:%s", conf.Conf.Host, conf.Conf.Port)
	e.Logger.Fatal(e.Start(path))
}
