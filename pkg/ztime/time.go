package ztime

import "time"

// GetCurrentYM 获取当前月份的年和月
func GetCurrentYM() string {
	t := time.Now()
	return t.Format("2006-01")
}

// GetLastMonthYM 获取当前月份的上一个月份的年和月。
func GetLastMonthYM() string {
	t := time.Now()
	lastMonth := t.AddDate(0, -1, 0)
	return lastMonth.Format("2006-01")
}
