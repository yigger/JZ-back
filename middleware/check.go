package middleware

import (
	"net/http"
	"github.com/labstack/echo"
)

func CheckOpenId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		header := req.Header.Get("WX-KEY")
		if (header == "abc123") {
			return next(c)
		} else {
			// 没能校验通过
			return c.String(http.StatusOK, "wrong")
		}
	}
}