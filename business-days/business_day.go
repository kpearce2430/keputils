package business_days

import (
	"github.com/sirupsen/logrus"
	"time"
)

func IsBeforeMarketOpen(current time.Time) bool {
	// If hour, minutes, and seconds are 0, then it's just the date and we don't know the time.
	if current.Hour() == 0 && current.Minute() == 0 && current.Second() == 0 {
		return false
	}

	if current.Hour() < 9 {
		return true
	}

	if current.Hour() < 10 && current.Minute() <= 30 {
		return true
	}

	return false
}

func GetBusinessDay(start time.Time) time.Time {
	logrus.Debug(start.Weekday())
	var reqDate time.Time

	switch start.Weekday() {
	case time.Saturday:
		reqDate = time.Date(start.Year(), start.Month(), start.Day()-1, 00, 00, 00, 00, time.UTC)
	case time.Sunday:
		reqDate = time.Date(start.Year(), start.Month(), start.Day()-2, 00, 00, 00, 00, time.UTC)
	default:
		reqDate = time.Date(start.Year(), start.Month(), start.Day(), 00, 00, 00, 00, time.UTC)
		if IsBeforeMarketOpen(start) == true {
			reqDate = time.Date(start.Year(), start.Month(), start.Day()-1, 00, 00, 00, 00, time.UTC)
		}
	}

	logrus.Debug("reqDate>", reqDate)
	if IsHoliday(reqDate) {
		reqDate = time.Date(reqDate.Year(), reqDate.Month(), reqDate.Day()-1, 00, 00, 00, 00, time.UTC)
		reqDate = GetBusinessDay(reqDate)
	}
	return reqDate
}
