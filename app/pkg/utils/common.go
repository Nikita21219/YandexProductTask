package utils

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return nil
	}
	return err
}

func ValidTime(time string) bool {
	re := regexp.MustCompile(`^\d{2}:\d{2}-\d{2}:\d{2}$`)
	if !re.MatchString(time) {
		return false
	}

	parts := strings.Split(time, "-")
	startTime := parts[0]
	endTime := parts[1]

	startHour, err := strconv.Atoi(startTime[:2])
	if err != nil {
		return false
	}
	startMinute, err := strconv.Atoi(startTime[3:5])
	if err != nil {
		return false
	}

	endHour, err := strconv.Atoi(endTime[:2])
	if err != nil {
		return false
	}
	endMinute, err := strconv.Atoi(endTime[3:5])
	if err != nil {
		return false
	}

	if startHour < 0 || startHour > 23 || startMinute < 0 || startMinute > 59 ||
		endHour < 0 || endHour > 23 || endMinute < 0 || endMinute > 59 {
		return false
	}

	if startHour > endHour || (startHour == endHour && startMinute > endMinute) {
		return false
	}

	return true
}