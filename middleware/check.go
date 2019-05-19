package middleware

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/yigger/JZ-back/model"
)

func CheckOpenId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		sessionKey := req.Header.Get("X-WX-Skey")

		if user := new(model.User); user.IsLogin(sessionKey) {
			return next(c)
		} else {
			return c.String(http.StatusOK, "error third session key")
		}
	}
}
