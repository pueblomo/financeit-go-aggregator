package util

import "time"

func GetDateWithFirstDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
}

func GetDateWithLastDay(t time.Time) time.Time {
	firstday := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	return firstday.AddDate(0, 1, 0).Add(time.Nanosecond * -1)
}