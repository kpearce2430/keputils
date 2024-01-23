package business_days_test

import (
	businessdays "github.com/kpearce2430/keputils/business-days"
	"testing"
	"time"
)

type testDates struct {
	Name          string
	Year          int
	Month         int
	Day           int
	Hour          int
	Minute        int
	want          bool
	expectedDay   int
	expectedMonth int
}

func TestGetBusinessDay_IsHoliday(t *testing.T) {
	t.Parallel()
	holidayTests := []testDates{
		{Name: "Christmas", Year: 2023, Month: 12, Day: 25, want: true},
		{Name: "July 4th", Year: 2024, Month: 7, Day: 4, want: true},
		{Name: "November 1st", Year: 2019, Month: 11, Day: 1, want: false},
	}

	for _, tc := range holidayTests {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			testDate := time.Date(tc.Year, time.Month(tc.Month), tc.Day, 00, 00, 00, 00, time.UTC)
			if businessdays.IsHoliday(testDate) != tc.want {
				t.Log("Test:", tc.Name, "Failed", tc.want)
				t.Fail()
			}
		})
	}
}

func TestGetBusinessDay_IsBeforeMarketOpen(t *testing.T) {
	t.Parallel()
	today := time.Now()
	openTests := []testDates{
		{Name: "8:00am", Hour: 8, Minute: 0, want: true},
		{Name: "9:29 am", Hour: 9, Minute: 29, want: true},
		{Name: "12:29 am", Hour: 12, Minute: 29, want: false},
		{Name: "00:00 am", Hour: 0, Minute: 0, want: false},
	}

	for _, tc := range openTests {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			testDate := time.Date(today.Year(), today.Month(), today.Day(), tc.Hour, tc.Minute, 00, 00, time.UTC)
			if businessdays.IsBeforeMarketOpen(testDate) != tc.want {
				t.Log("Test:", tc.Name, "Failed", tc.want)
				t.Fail()
			}
		})
	}
}

func TestGetBusinessDay_GetBusinessDay(t *testing.T) {
	t.Parallel()
	openTests := []testDates{
		{Name: "Dec 22 Friday", Year: 2023, Month: 12, Day: 22, expectedDay: 22, expectedMonth: 12},
		{Name: "Dec 23 Saturday", Year: 2023, Month: 12, Day: 23, expectedDay: 22, expectedMonth: 12},
		{Name: "Dec 24 Sunday", Year: 2023, Month: 12, Day: 24, expectedDay: 22, expectedMonth: 12},
		{Name: "Christmas", Year: 2023, Month: 12, Day: 25, expectedDay: 22, expectedMonth: 12},
		{Name: "Boxing Day", Year: 2023, Month: 12, Day: 26, expectedDay: 26, expectedMonth: 12},
		{Name: "Dec 27 Wednesday", Year: 2023, Month: 12, Day: 27, expectedDay: 27, expectedMonth: 12},
		{Name: "Dec 28 Thursday", Year: 2023, Month: 12, Day: 28, expectedDay: 28, expectedMonth: 12},
		{Name: "Dec 29 Friday", Year: 2023, Month: 12, Day: 29, expectedDay: 29, expectedMonth: 12},
		{Name: "Dec 30 Saturday", Year: 2023, Month: 12, Day: 30, expectedDay: 29, expectedMonth: 12},
		{Name: "New Years Eve", Year: 2023, Month: 12, Day: 31, expectedDay: 29, expectedMonth: 12},
		{Name: "New Years Day", Year: 2024, Month: 1, Day: 1, expectedDay: 29, expectedMonth: 12},
		{Name: "Jan 2nd", Year: 2024, Month: 1, Day: 2, expectedDay: 2, expectedMonth: 1},
		{Name: "Junteenth 2024", Year: 2024, Month: 6, Day: 19, expectedDay: 18, expectedMonth: 6},
	}

	for _, tc := range openTests {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			testDate := time.Date(tc.Year, time.Month(tc.Month), tc.Day, tc.Hour, tc.Minute, 00, 00, time.UTC)
			result := businessdays.GetBusinessDay(testDate)
			if result.Month() != time.Month(tc.expectedMonth) && result.Day() != tc.expectedDay {
				t.Log("Test:", tc.Name, "Failed:", result)
				t.Fail()
			}
		})
	}
}
