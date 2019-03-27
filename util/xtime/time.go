package xtime

import (
	"time"
)

const DATE = `2006-01-02`
const DATE_TIME = `2006-01-02 15:04:05`
const DATE_TIME_START_SECOND = `2006-01-02 15:04:00`
const DATE_TIME_ZERO_HOUR = `2006-01-02 00:00:00`

func ParseInLocal(layout, value string) (time.Time, error) {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return time.Time{}, err
	}

	return time.ParseInLocation(layout, value, loc)
}

//获取月份开始、结束
func GetMonthStartEnd(dd time.Time) (start time.Time, end time.Time) {
	year, month, _ := dd.Date()
	loc := dd.Location()

	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	return startOfMonth, endOfMonth
}
