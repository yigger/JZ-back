package tests

import (
	"testing"
	// "net/http"
	// "net/http/httptest"
	// "github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/yigger/JZ-back/service"
)

func TestLogin(t *testing.T) {
	Convey("Exist user Given code to login", t, func() {
		code := "abc"
		// user, err := service.User.Login(code)
		
	})

	Convey("Non-Exist user Given code to login", t, func() {
		
	})
}