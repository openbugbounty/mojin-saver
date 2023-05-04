package util

import "time"

func GetMidNightTime() time.Time {
	strNowDay := time.Now().Format("2006-01-02")
	nowDay, _ := time.ParseInLocation("2006-01-02", strNowDay, time.Local)
	return nowDay
}
