package model

import (
	"fmt"
	"strings"
)

// 顶部标题栏 - 位置 1
var POSITION_ONE_HUMAN = map[string]string{
	"today_expend": "今日支出",
	"today_income": "今日收入",
	"today_surplus": "今日结余",
	"month_expend": "本月支出",
	"month_income": "本月收入",
	"month_surplus": "本月结余",
}

// 顶部标题栏 - 位置 2
var POSITION_TWO_HUMAN = map[string]string{
	"week_expend": "本周支出",
	"week_surplus": "本周结余",
	"month_expend": "本月支出",
	"month_surplus": "本月结余",
	"year_expend": "本年支出",
	"year_surplus": "本年结余",
}

// 顶部标题栏 - 位置 3
var POSITION_THREE_HUMAN = map[string]string{
	"month_expend": "本月支出",
	"month_surplus": "本月结余",
	"year_expend": "本年支出",
	"year_income": "本年收入",
	"year_surplus": "本年结余",
	"month_budget": "预算剩余",
}

func (user *User) GetHeaderAmount(headerName string) float64 {
	if headerName == "month_budget" {
		return user.monthBudget()
	}

	s := strings.Split(headerName, "_")
	name, countType := s[0], s[1]

	scope := StatementInDay
	switch name {
		case "month": 
			scope = StatementInMonth
		case "year":
			scope = StatementInYear
		case "week":
			scope = StatementInWeek
	}
	
	if countType == "income" || countType == "expend" {
		// 收入或支出
		rows, err := db.Table("statements").Scopes(scope).Where("user_id = ? AND type = ?", user.ID, countType).Select("sum(amount)").Rows()
		if err != nil {
			fmt.Println(err)
			return 0
		}

		if rows.Next() {
			var amount float64
			if err := rows.Scan(&amount); err != nil {
				return 0
			}
			return amount
		}
	} else {
		// 结余
		return user.GetHeaderAmount(name + "_income") - user.GetHeaderAmount(name + "_expend")
	}

	return 0
}

func (user *User) monthBudget() float64 {
	var amount float64
	rows, err := db.Table("statements").Scopes(StatementInMonth).Where("user_id = ? AND type = 'expend'", user.ID).Select("sum(amount)").Rows()
	if err != nil {
		fmt.Println(err)
		return 0
	}

	if rows.Next() {
		if err := rows.Scan(&amount); err != nil {
			return 0
		}
	}

	return user.Budget - amount
}
