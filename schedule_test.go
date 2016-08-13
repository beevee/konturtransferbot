package main

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSchedule(t *testing.T) {
	Convey("Given correct YAML schedule", t, func() {
		sYaml := `WorkDayRouteToOffice:
  - "07:30"
  - "08:00"
  - "20:00"
  - "20:30"
HolidayRouteToOffice:
  - "10:30"
WorkDayRouteFromOffice:
  - "08:20"
  - "08:50"
  - "20:20"
  - "20:50"
HolidayRouteFromOffice:
  - "18:00"`
		Convey("It should parse into a schedule structure", func() {
			s := buildSchedule([]byte(sYaml))
			Convey("Its entries should be correct and in the same order", func() {
				So(s.workDayRouteToOffice[0].Format("15:04"), ShouldEqual, "07:30")
				So(s.workDayRouteToOffice[1].Format("15:04"), ShouldEqual, "08:00")
				So(s.workDayRouteToOffice[2].Format("15:04"), ShouldEqual, "20:00")
				So(s.workDayRouteToOffice[3].Format("15:04"), ShouldEqual, "20:30")

				So(s.holidayRouteToOffice[0].Format("15:04"), ShouldEqual, "10:30")

				So(s.workDayRouteFromOffice[0].Format("15:04"), ShouldEqual, "08:20")
				So(s.workDayRouteFromOffice[1].Format("15:04"), ShouldEqual, "08:50")
				So(s.workDayRouteFromOffice[2].Format("15:04"), ShouldEqual, "20:20")
				So(s.workDayRouteFromOffice[3].Format("15:04"), ShouldEqual, "20:50")

				So(s.holidayRouteFromOffice[0].Format("15:04"), ShouldEqual, "18:00")
			})

			Convey("It should correctly identify Friday as a workday", func() {
				now, _ := time.Parse("02.01.2006 15:04", "12.08.2016 07:00")
				So(s.findCorrectRoute(now, true).String(), ShouldEqual, s.workDayRouteToOffice.String())
				So(s.findCorrectRoute(now, false).String(), ShouldEqual, s.workDayRouteFromOffice.String())
			})

			Convey("It should correctly identify Sunday as a holiday", func() {
				now, _ := time.Parse("02.01.2006 15:04", "14.08.2016 07:00")
				So(s.findCorrectRoute(now, true).String(), ShouldEqual, s.holidayRouteToOffice.String())
				So(s.findCorrectRoute(now, false).String(), ShouldEqual, s.holidayRouteFromOffice.String())
			})

			Convey("It should correctly return whole schedule to office", func() {
				texts := s.getFullToOfficeTexts()
				So(texts[0], ShouldEqual, "Дежурные рейсы от Геологической в будни:\n07:30\n08:00\n20:00\n20:30\n")
				So(texts[1], ShouldEqual, "Дежурные рейсы от Геологической в выходные:\n10:30\n")
			})

			Convey("It should correctly return whole schedule from office", func() {
				texts := s.getFullFromOfficeTexts()
				So(texts[0], ShouldEqual, "Дежурные рейсы от офиса в будни:\n08:20\n08:50\n20:20\n20:50\n")
				So(texts[1], ShouldEqual, "Дежурные рейсы от офиса в выходные:\n18:00\n")
			})
		})
	})
}
