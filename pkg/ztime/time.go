package ztime

import "time"

func GetCurrentYM() string {
	t := time.Now()
	return t.Format("2006-01")
}
