package konturtransferbot

import (
	"time"

	"github.com/rickar/cal/v2"
	"github.com/rickar/cal/v2/ru"
)

// Schedule contains all information on transfer departure times
type Schedule struct {
	TransferRoutes
	HolidaysCalendar *cal.BusinessCalendar
}

type TransferRoutes struct {
	WorkDayRouteToOffice    Route
	WorkDayRouteFromOffice  Route
	SaturdayRouteToOffice   Route
	SaturdayRouteFromOffice Route
}

// NewSchedule is construct for schedule.
func NewSchedule(transferSchedule TransferRoutes) Schedule {
	return Schedule{
		TransferRoutes:   transferSchedule,
		HolidaysCalendar: newCustomCalendar(),
	}
}

// newCustomCalendar создаёт календарь с российскими праздниками
func newCustomCalendar() *cal.BusinessCalendar {
	c := cal.NewBusinessCalendar()

	// Добавляем российские праздники
	c.AddHoliday(
		ru.NewYear,           // Новый Год (1 января)
		ru.OrthodoxChristmas, // Рождество (7 января)
		ru.LabourDay,
		ru.MilitaryDay,
		ru.UnionDay,
		ru.VictoryDay,
		ru.RussiasDay,
		ru.WomensDay,
	)
	return c
}

// isHoliday проверяет, является ли текущая дата выходным или праздничным днём
func isHoliday(cal *cal.BusinessCalendar, date time.Time) bool {
	return !cal.IsWorkday(date)
}

// GetToOfficeText returns text representation of full schedule to office
func (s Schedule) GetToOfficeText(now time.Time) (string, string) {
	prefix := "*Рейсы в офис*\n\n"
	suffix := "\nВ выходные дни трансфера нет"

	timeAgnosticRoute := prefix + s.WorkDayRouteToOffice.String() + suffix
	if isHoliday(s.HolidaysCalendar, now) {
		return timeAgnosticRoute, ""
	}

	timeSensitiveRoute := prefix + s.WorkDayRouteToOffice.StringWithDivider(now) + suffix
	if timeAgnosticRoute == timeSensitiveRoute {
		return timeAgnosticRoute, ""
	}

	return timeSensitiveRoute, timeAgnosticRoute
}

// GetFromOfficeText returns text representation of full schedule from office
func (s Schedule) GetFromOfficeText(now time.Time) (string, string) {
	prefix := "*Рейсы из офиса*\n\n"
	suffix := "\nВ выходные дни трансфера нет"

	timeAgnosticRoute := prefix + s.WorkDayRouteFromOffice.String() + suffix
	if isHoliday(s.HolidaysCalendar, now) {
		return timeAgnosticRoute, ""
	}

	timeSensitiveRoute := prefix + s.WorkDayRouteFromOffice.StringWithDivider(now) + suffix
	if timeAgnosticRoute == timeSensitiveRoute {
		return timeAgnosticRoute, ""
	}

	return timeSensitiveRoute, timeAgnosticRoute
}
