package utils

import (
	"fmt"
	"time"
)

func JulDate() string {
	return JulDateFromTime(time.Now())
}

func JulDateFromTime(tm time.Time) string {
	return fmt.Sprintf("%d%03d", tm.Year(), tm.YearDay())
}
