package controller

import (
	"os"
	"github.com/b3log/gulu"
)

// Result represents HTTP response body.
type Result struct {
	Status  int         `json:"status"` // return code, 0 for succ
	Msg  	string      `json:"msg"`  // message
	Data    interface{} `json:"data"` // data object
	Error   interface{}	`json:"error"`
}

// NewResult creates a result with Code=0, Msg="", Data=nil.
func RenderJson() *Result {
	return &Result{
		Status: 200,
		Msg:  "",
		Data: nil,
		Error: nil,
	}
}

// Result codes.
const (
	CodeOk      = 200  // OK
	CodeErr     = -1 // general error
	CodeAuthErr = 2  // unauthenticated request
)

var logger = gulu.Log.NewLogger(os.Stdout)
