package logs

import (
	"os"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func LoadLog() {
	log.SetOutput(os.Stdout)
}

func Info(message string) {
	log.Info(message)
}

