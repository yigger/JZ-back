package main

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/yigger/JZ-back/middleware"
	routers "github.com/yigger/JZ-back/routers"
)

func main() {
	e := echo.New()
	g := e.Group("/", middleware.CheckOpenId)

	g.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	g.GET("/test", routers.UserLogin)

	e.Logger.Fatal(e.Start(":1323"))
}
