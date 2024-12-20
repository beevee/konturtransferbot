package konturtransferbot

import (
	"time"

	"github.com/rickar/cal/v2"
	"github.com/rickar/cal/v2/ru"
)

// Schedule contains all information on transfer departure times
type Schedule struct {
	WorkDayRouteToOffice    Route
	WorkDayRouteFromOffice  Route
	SaturdayRouteToOffice   Route
	SaturdayRouteFromOffice Route
}

// customCalendar создаёт календарь с российскими праздниками
func customCalendar() *cal.BusinessCalendar {
	c := cal.NewBusinessCalendar()
	// Добавляем российские праздники
	c.AddHoliday(
		ru.NewYear,      // Новый Год (1 января)
		ru.ChristmasDay, // Рождество (7 января)
		ru.DefenderOfFatherlandDay,
		ru.InternationalWomenDay,
		ru.SpringAndLabourDay,
		ru.VictoryDay,
		ru.RussiaDay,
		ru.NationalUnityDay,
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

	cal := customCalendar() // инициализация календаря

	timeAgnosticRoute := prefix + s.WorkDayRouteToOffice.String() + suffix
	if isHoliday(cal, now) {
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
	if isHoliday(cal, now) {
		return timeAgnosticRoute, ""
	}

	timeSensitiveRoute := prefix + s.WorkDayRouteFromOffice.StringWithDivider(now) + suffix
	if timeAgnosticRoute == timeSensitiveRoute {
		return timeAgnosticRoute, ""
	}

	return timeSensitiveRoute, timeAgnosticRoute
}
