package main

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/yigger/JZ-back/middleware"
	// routers "github.com/yigger/JZ-back/routers"
	"github.com/yigger/JZ-back/conf"
	"fmt"
)

func init() {
	conf.LoadConf()
}

func main() {
	e := echo.New()
	e.Use(middleware.CheckOpenId)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	server := fmt.Sprintf("%s:%s", conf.Conf.Host, conf.Conf.Port)
	e.Logger.Fatal(e.Start(server))
}
