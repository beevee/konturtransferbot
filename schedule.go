package main

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v1"
)

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
