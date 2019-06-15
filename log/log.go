package log

import (
	"os"
	"github.com/b3log/gulu"
)

var Log = gulu.Log.NewLogger(os.Stdout)
