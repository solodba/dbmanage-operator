package controller

import (
	"strconv"
	"strings"
	"time"
)

// 获取任务开始时间
func (r *DbManageReconciler) GetTaskDelaySeconds(startTime string) time.Duration {
	times := strings.Split(startTime, ":")
	expectedHour, _ := strconv.Atoi(times[0])
	expectedMin, _ := strconv.Atoi(times[1])
	now := time.Now().Truncate(time.Second)
	todayDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrowDate := todayDate.Add(time.Hour * 24)
	curDuration := time.Hour*time.Duration(now.Hour()) + time.Minute*time.Duration(now.Minute())
	expectedDuration := time.Hour*time.Duration(expectedHour) + time.Minute*time.Duration(expectedMin)
	var seconds int
	if curDuration >= expectedDuration {
		seconds = int(tomorrowDate.Add(expectedDuration).Sub(now).Seconds())
	} else {
		seconds = int(todayDate.Add(expectedDuration).Sub(now).Seconds())
	}
	return time.Second * time.Duration(seconds)
}

// 获取任务下次执行时间
func (r *DbManageReconciler) GetTaskNextTime(seconds float64) time.Time {
	return time.Now().Add(time.Second * time.Duration(seconds))
}
