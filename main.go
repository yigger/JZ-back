package main

import (
	"net/http"
	"github.com/labstack/echo"
	. "github.com/yigger/JZ-back/middleware"
	// "fmt"
)

func main() {
	e := echo.New()
	g := e.Group("/api", CheckOpenId)

	g.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	g.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, test!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
