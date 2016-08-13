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
