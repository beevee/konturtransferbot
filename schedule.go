package main

import (
	"time"

	"gopkg.in/yaml.v1"
)

type schedule struct {
	workDayRouteToOffice   route
	workDayRouteFromOffice route
	holidayRouteToOffice   route
	holidayRouteFromOffice route
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

// ScheduleYaml - модель расписания для конфига
type ScheduleYaml struct {
	WorkDayRouteToOffice   []string `yaml:"WorkDayRouteToOffice"`
	WorkDayRouteFromOffice []string `yaml:"WorkDayRouteFromOffice"`
	HolidayRouteToOffice   []string `yaml:"HolidayRouteToOffice"`
	HolidayRouteFromOffice []string `yaml:"HolidayRouteFromOffice"`
}

func buildSchedule(data []byte) schedule {
	scheduleYaml := ScheduleYaml{}
	err := yaml.Unmarshal([]byte(data), &scheduleYaml)
	if err != nil {
		//log.Fatal(err)
	}
	result := schedule{
		workDayRouteToOffice:   buildRoute(scheduleYaml.WorkDayRouteToOffice),
		workDayRouteFromOffice: buildRoute(scheduleYaml.WorkDayRouteFromOffice),
		holidayRouteToOffice:   buildRoute(scheduleYaml.HolidayRouteToOffice),
		holidayRouteFromOffice: buildRoute(scheduleYaml.HolidayRouteFromOffice),
	}
	return result
}
