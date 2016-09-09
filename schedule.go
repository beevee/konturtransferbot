package main

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v1"
)

const monetizationMessage = "Скидка 50% на 10 поездок до 18 сентября в Uber — EKB50. Промокод на первую бесплатную поездку: Gett — GTFUNKP, Яндекс.Такси — daf3qsau, Uber — ykt6m, Wheely — MHPRL."

// ScheduleYaml - модель расписания для конфига
type ScheduleYaml struct {
	WorkDayRouteToOffice   []string `yaml:"WorkDayRouteToOffice"`
	WorkDayRouteFromOffice []string `yaml:"WorkDayRouteFromOffice"`
	HolidayRouteToOffice   []string `yaml:"HolidayRouteToOffice"`
	HolidayRouteFromOffice []string `yaml:"HolidayRouteFromOffice"`
}

type schedule struct {
	workDayRouteToOffice   route
	workDayRouteFromOffice route
	holidayRouteToOffice   route
	holidayRouteFromOffice route
}

func buildSchedule(data []byte) (*schedule, error) {
	scheduleYaml := ScheduleYaml{}
	if err := yaml.Unmarshal([]byte(data), &scheduleYaml); err != nil {
		return nil, err
	}
	workDayRouteToOffice, err := buildRoute(scheduleYaml.WorkDayRouteToOffice)
	if err != nil {
		return nil, err
	}
	workDayRouteFromOffice, err := buildRoute(scheduleYaml.WorkDayRouteFromOffice)
	if err != nil {
		return nil, err
	}
	holidayRouteToOffice, err := buildRoute(scheduleYaml.HolidayRouteToOffice)
	if err != nil {
		return nil, err
	}
	holidayRouteFromOffice, err := buildRoute(scheduleYaml.HolidayRouteFromOffice)
	if err != nil {
		return nil, err
	}
	return &schedule{
		workDayRouteToOffice:   workDayRouteToOffice,
		workDayRouteFromOffice: workDayRouteFromOffice,
		holidayRouteToOffice:   holidayRouteToOffice,
		holidayRouteFromOffice: holidayRouteFromOffice,
	}, nil
}

func (s schedule) findCorrectRoute(now time.Time, toOffice bool) route {
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		if toOffice {
			return s.holidayRouteToOffice
		}
		return s.holidayRouteFromOffice
	}

	if toOffice {
		return s.workDayRouteToOffice
	}
	return s.workDayRouteFromOffice
}

func (s schedule) getBestTripToOfficeText(now time.Time) string {
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

func (s schedule) getBestTripFromOfficeText(now time.Time) string {
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

func (s schedule) getFullToOfficeTexts() []string {
	return []string{
		fmt.Sprintf("Дежурные рейсы от Геологической в будни:\n%s", s.workDayRouteToOffice),
		fmt.Sprintf("Дежурные рейсы от Геологической в выходные:\n%s", s.holidayRouteToOffice),
	}
}

func (s schedule) getFullFromOfficeTexts() []string {
	return []string{
		fmt.Sprintf("Дежурные рейсы от офиса в будни:\n%s", s.workDayRouteFromOffice),
		fmt.Sprintf("Дежурные рейсы от офиса в выходные:\n%s", s.holidayRouteFromOffice),
	}
}
