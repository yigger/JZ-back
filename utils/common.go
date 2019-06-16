package utils

import (
	"github.com/leekchan/accounting"
	"time"
)

var WeekMap = map[string]string{
	"Monday": 		"周一",
	"Tuesday": 		"周二",
	"Wednesday": 	"周三",
	"Thursday": 	"周四",
	"Friday":		"周五",
	"Saturday": 	"周六",
	"Sunday": 		"周日",
}

var (
	CurrentYear = time.Now().Format("2006")
	CurrentMonth = time.Now().Format("01")
	CurrentDay = time.Now().Format("02")
)

func FormatMoney(money float64) string {
	ac := accounting.Accounting{Symbol: "", Precision: 2}
	return ac.FormatMoney(money)
}
