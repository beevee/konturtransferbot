package konturtransferbot

import "time"

// Schedule contains all information on transfer departure times
type Schedule struct {
	WorkDayRouteToOffice    Route
	WorkDayRouteFromOffice  Route
	SaturdayRouteToOffice   Route
	SaturdayRouteFromOffice Route
}

// GetToOfficeText returns text representation of full schedule to office
func (s Schedule) GetToOfficeText(now time.Time) string {
	var routeText string
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		routeText = s.WorkDayRouteToOffice.String()
	} else {
		routeText = s.WorkDayRouteToOffice.StringWithDivider(now)
	}
	return "*Геологическая → Офис*\n\n" + routeText + "\nСубботний рейс в " + s.SaturdayRouteToOffice.String()
}

// GetFromOfficeText returns text representation of full schedule from office
func (s Schedule) GetFromOfficeText(now time.Time) string {
	var routeText string
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		routeText = s.WorkDayRouteFromOffice.String()
	} else {
		routeText = s.WorkDayRouteFromOffice.StringWithDivider(now)
	}
	return "*Офис → Геологическая*\n\n" + routeText + "\nСубботний рейс в " + s.SaturdayRouteFromOffice.String()
}
