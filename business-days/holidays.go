package business_days

import (
	"github.com/sirupsen/logrus"
	"time"
)

type HolidayDates struct {
	Description string
	Year        int
	Month       int
	Day         int
}

func IsHoliday(current time.Time) bool {

	holidays2024 := []HolidayDates{
		{Year: 2024, Month: 12, Day: 25, Description: "Christmas"},
		{Year: 2024, Month: 11, Day: 28, Description: "Thanksgiving"},
		{Year: 2024, Month: 9, Day: 02, Description: "Labor Day"},
		{Year: 2024, Month: 07, Day: 04, Description: "Independence Day"},
		{Year: 2024, Month: 06, Day: 19, Description: "Juneteenth"},
		{Year: 2024, Month: 05, Day: 27, Description: "Memorial Day"},
		{Year: 2024, Month: 03, Day: 29, Description: "Good Friday"},
		{Year: 2024, Month: 02, Day: 19, Description: "Washington's Birthday"},
		{Year: 2024, Month: 01, Day: 15, Description: "Martin Luther King, Jr. Day"},
		{Year: 2024, Month: 01, Day: 01, Description: "New Years Day"},
	}

	holidays2023 := []HolidayDates{
		{Year: 2023, Month: 12, Day: 25, Description: "Christmas"},
		{Year: 2023, Month: 11, Day: 23, Description: "Thanksgiving"},
		{Year: 2023, Month: 9, Day: 04, Description: "Labor Day"},
		{Year: 2023, Month: 07, Day: 04, Description: "July 4th"},
		{Year: 2023, Month: 06, Day: 19, Description: "Junteenth"},
		{Year: 2023, Month: 05, Day: 29, Description: "Memorial Day"},
		{Year: 2023, Month: 04, Day: 07, Description: "Good Friday"},
		{Year: 2023, Month: 02, Day: 20, Description: "Presidents Day"},
		{Year: 2023, Month: 01, Day: 16, Description: "Martin Luther King, Jr. Day"},
		{Year: 2023, Month: 01, Day: 02, Description: "New Years Day Observed"},
	}

	holidays2022 := []HolidayDates{
		{Year: 2022, Month: 12, Day: 26, Description: "Christmas Observed"},
		{Year: 2022, Month: 11, Day: 24, Description: "Thanksgiving"},
		{Year: 2022, Month: 9, Day: 05, Description: "Labor Day"},
		{Year: 2022, Month: 07, Day: 04, Description: "July 4th"},
		{Year: 2022, Month: 06, Day: 20, Description: "Junteenth"},
		{Year: 2022, Month: 05, Day: 30, Description: "Memorial Day"},
		{Year: 2022, Month: 04, Day: 15, Description: "Good Friday"},
		{Year: 2022, Month: 02, Day: 21, Description: "Presidents Day"},
		{Year: 2022, Month: 01, Day: 17, Description: "MLK Jr Day"},
	}

	holidays2021 := []HolidayDates{
		{Year: 2021, Month: 9, Day: 06, Description: "Labor Day"},
		{Year: 2021, Month: 12, Day: 24, Description: "Christmas Observed"},
		{Year: 2021, Month: 11, Day: 25, Description: "Thanksgiving"},
		{Year: 2021, Month: 07, Day: 05, Description: "July 4th Observed"},
		{Year: 2021, Month: 05, Day: 31, Description: "Memorial Day"},
		{Year: 2021, Month: 04, Day: 02, Description: "Good Friday"},
		{Year: 2021, Month: 02, Day: 15, Description: "Presidents Day"},
		{Year: 2021, Month: 01, Day: 18, Description: "MLK Jr Day"},
		{Year: 2021, Month: 01, Day: 01, Description: "New Years Day"},
	}

	holidays2020 := []HolidayDates{
		{Year: 2020, Month: 9, Day: 07, Description: "Labor Day"},
		{Year: 2020, Month: 12, Day: 25, Description: "Christmas"},
		{Year: 2020, Month: 11, Day: 26, Description: "Thanksgiving"},
		{Year: 2020, Month: 07, Day: 03, Description: "July 4th Observed"},
		{Year: 2020, Month: 05, Day: 25, Description: "Memorial Day"},
		{Year: 2020, Month: 04, Day: 10, Description: "Good Friday"},
		{Year: 2020, Month: 02, Day: 17, Description: "Presidents Day"},
		{Year: 2020, Month: 01, Day: 20, Description: "MLK Jr Day"},
		{Year: 2020, Month: 01, Day: 01, Description: "New Years Day"},
	}

	holidays2019 := []HolidayDates{
		{Year: 2019, Month: 9, Day: 02, Description: "Labor Day"},
		{Year: 2019, Month: 12, Day: 25, Description: "Christmas"},
		{Year: 2019, Month: 11, Day: 28, Description: "Thanksgiving"},
		{Year: 2019, Month: 07, Day: 04, Description: "July 4th"},
		{Year: 2019, Month: 05, Day: 27, Description: "Memorial Day"},
		{Year: 2019, Month: 04, Day: 19, Description: "Good Friday"},
		{Year: 2019, Month: 02, Day: 18, Description: "Presidents Day"},
		{Year: 2019, Month: 01, Day: 21, Description: "MLK Jr Day"},
		{Year: 2019, Month: 01, Day: 01, Description: "New Years Day"},
	}

	if current.Month() == time.January && current.Day() == 1 {
		return true
	}

	if current.Month() == time.July && current.Day() == 4 {
		return true
	}

	if current.Month() == time.December && current.Day() == 25 {
		return true
	}

	var holidays []HolidayDates
	switch current.Year() {
	case 2024:
		holidays = holidays2024
	case 2023:
		holidays = holidays2023
	case 2022:
		holidays = holidays2022
	case 2021:
		holidays = holidays2021
	case 2020:
		holidays = holidays2020
	case 2019:
		holidays = holidays2019
	default:
		// TODO:  Make more dynamic
		logrus.Error("Invalid Year:", current.Year())
		panic("invalid year")
	}

	for _, h := range holidays {
		if current.Year() == h.Year && int(current.Month()) == h.Month && current.Day() == h.Day {
			return true
		}
	}
	return false
}
