package cron

import (
	"fmt"
	"testing"
)

func TestCrontabAddBySec(t *testing.T) {
	AddScheduleBySec(5, func() {
		fmt.Println("test")
	})
}

func TestCrontabAddByMin(t *testing.T) {
	AddScheduleByMin(1, func() {
		fmt.Println("test")
	})
}

func TestCrontabAddByHours(t *testing.T) {
	AddScheduleByHours(1, func() {
		fmt.Println("test")
	})
}
