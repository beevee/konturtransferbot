package main

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRoute(t *testing.T) {
	Convey("Given route with several departures", t, func() {
		r := route{
			time.Date(0, 0, 0, 10, 30, 0, 0, &time.Location{}),
			time.Date(0, 0, 0, 11, 0, 0, 0, &time.Location{}),
			time.Date(0, 0, 0, 12, 30, 0, 0, &time.Location{}),
		}
		Convey("Its string representation should contain all of them", func() {
			So(r.String(), ShouldEqual, "10:30\n11:00\n12:30\n")
		})
	})

	Convey("Given string representation of a route", t, func() {
		rString := []string{
			"7:30",
			"11:00",
			"12:55",
		}
		Convey("It should parse correctly", func() {
			r := buildRoute(rString)
			Convey("Its numerical representation should be correct and in the same order", func() {
				So(r[0].Hour(), ShouldEqual, 7)
				So(r[0].Minute(), ShouldEqual, 30)
				So(r[1].Hour(), ShouldEqual, 11)
				So(r[1].Minute(), ShouldEqual, 0)
				So(r[2].Hour(), ShouldEqual, 12)
				So(r[2].Minute(), ShouldEqual, 55)
			})
		})
	})
}

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

			Convey("It should correctly find two best trips", func() {
				now, _ := time.Parse("15:04", "07:00")
				bestTrip, nextBestTrip := findBestTripMatches(now, s.workDayRouteToOffice)
				So(bestTrip.Format("15:04"), ShouldEqual, "07:30")
				So(nextBestTrip.Format("15:04"), ShouldEqual, "08:00")
			})

			Convey("It should correctly find one best trip", func() {
				now, _ := time.Parse("15:04", "20:15")
				bestTrip, nextBestTrip := findBestTripMatches(now, s.workDayRouteToOffice)
				So(bestTrip.Format("15:04"), ShouldEqual, "20:30")
				So(nextBestTrip, ShouldBeNil)
			})

			Convey("It should not offer trips more than (roughly) 5 hours in advance", func() {
				now, _ := time.Parse("15:04", "02:25")
				bestTrip, nextBestTrip := findBestTripMatches(now, s.workDayRouteToOffice)
				So(bestTrip, ShouldBeNil)
				So(nextBestTrip, ShouldBeNil)
			})
		})
	})
}
