package ctime

import "time"

type CTime struct {
	time.Time
}

func Now() CTime {
	return CTime{time.Now()}
}

func (t CTime) StartOfMonth() *time.Time {
	now := time.Now()
	startOfMonthDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	return &startOfMonthDate
}
