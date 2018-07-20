package cron

import (
	"errors"
	"time"
)

type test func(int64)

// AddScheduleBySec : Execute f every period second interval
func AddScheduleBySec(sec int64, f func()) error {
	if sec < 0 {
		return errors.New(" period can not be negative ")
	}
	ticker := time.NewTicker(time.Duration(1000 * 1000 * 1000 * sec))
	defer ticker.Stop()
	for {
		<-ticker.C
		f()
	}
	return nil
}

// AddScheduleByMin : Execute f every period minutes interval
func AddScheduleByMin(min int64, f func()) {
	AddScheduleBySec(min*60, f)
}

// AddScheduleByHours : Execute f every period hours interval
func AddScheduleByHours(hours int64, f func()) {
	AddScheduleByMin(hours*60, f)
}
