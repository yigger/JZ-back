package routers

import (
	"net/http"
	"github.com/labstack/echo"
)

func UserLogin(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, test!")
}
