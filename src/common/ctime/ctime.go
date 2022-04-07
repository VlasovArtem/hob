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

func (t CTime) StartOfYear() *time.Time {
	now := time.Now()
	startOfYearDate := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	return &startOfYearDate
}

func (t CTime) StartOfYearAndCurrent() (from *time.Time, to *time.Time) {
	now := time.Now()
	startOfYearDate := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	from = &startOfYearDate
	to = &now
	return
}

func (t CTime) CurrentMonthName() string {
	now := time.Now()
	return now.Month().String()
}
