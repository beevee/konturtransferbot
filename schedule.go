package konturtransferbot

import (
	"fmt"
	"time"
)

const monetizationMessage = "Промокод на первую бесплатную поездку: Gett — GTFUNKP, Яндекс.Такси — daf3qsau, Uber — ykt6m, Wheely — MHPRL."

// Schedule contains all information on transfer departure times
type Schedule struct {
	WorkDayRouteToOffice   Route
	WorkDayRouteFromOffice Route
	HolidayRouteToOffice   Route
	HolidayRouteFromOffice Route
}

func (s Schedule) findCorrectRoute(now time.Time, toOffice bool) Route {
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		if toOffice {
			return s.HolidayRouteToOffice
		}
		return s.HolidayRouteFromOffice
	}

	if toOffice {
		return s.WorkDayRouteToOffice
	}
	return s.WorkDayRouteFromOffice
}

// GetBestTripToOfficeText returns text recommendation for two (or maybe less) closest departures to office
func (s Schedule) GetBestTripToOfficeText(now time.Time) string {
	var reply string

	bestTrip, nextBestTrip := s.findCorrectRoute(now, true).findBestTripMatches(now)
	if bestTrip != nil {
		reply = fmt.Sprintf("Ближайший дежурный рейс от Геологической будет в %s.", bestTrip.Format("15:04"))
		if nextBestTrip != nil {
			reply += fmt.Sprintf(" Следующий - в %s.", nextBestTrip.Format("15:04"))
		} else {
			reply += " Это последний на сегодня рейс."
		}
	} else {
		reply = "В ближайшие несколько часов уехать на работу на трансфере не получится. Лучше лечь поспать и поехать с утра. Первые рейсы от Геологической: "
		nextDay := now.Add(24 * time.Hour)
		currentRoute := s.findCorrectRoute(nextDay, true)
		for index, trip := range currentRoute {
			if trip.Hour() >= 12 || index >= len(currentRoute)-1 {
				reply += fmt.Sprintf("%s.", trip.Format("15:04"))
				break
			}
			reply += fmt.Sprintf("%s, ", trip.Format("15:04"))
		}
	}

	return reply
}

// GetBestTripFromOfficeText returns text recommendation for two (or maybe less) closest departures from office
func (s Schedule) GetBestTripFromOfficeText(now time.Time) string {
	var reply string

	bestTrip, nextBestTrip := s.findCorrectRoute(now, false).findBestTripMatches(now)
	if bestTrip != nil {
		reply = fmt.Sprintf("Ближайший дежурный рейс от офиса будет в %s.", bestTrip.Format("15:04"))
		if nextBestTrip != nil {
			reply += fmt.Sprintf(" Следующий - в %s.", nextBestTrip.Format("15:04"))
		} else {
			reply += " Это последний на сегодня рейс, дальше - только на такси. " + monetizationMessage
		}
	} else {
		reply = "В ближайшие несколько часов уехать домой на трансфере не получится :( Придется остаться в офисе или ехать на такси. " + monetizationMessage
	}

	return reply
}

// GetFullToOfficeTexts returns text representation of full schedule to office
func (s Schedule) GetFullToOfficeTexts() []string {
	return []string{
		fmt.Sprintf("Дежурные рейсы от Геологической в будни:\n%s", s.WorkDayRouteToOffice),
		fmt.Sprintf("Дежурные рейсы от Геологической в выходные:\n%s", s.HolidayRouteToOffice),
	}
}

// GetFullFromOfficeTexts returns text representation of full schedule from office
func (s Schedule) GetFullFromOfficeTexts() []string {
	return []string{
		fmt.Sprintf("Дежурные рейсы от офиса в будни:\n%s", s.WorkDayRouteFromOffice),
		fmt.Sprintf("Дежурные рейсы от офиса в выходные:\n%s", s.HolidayRouteFromOffice),
	}
}
