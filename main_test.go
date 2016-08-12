package main

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTimeWithoutDate(t *testing.T) {
	Convey("Given time with single digits", t, func() {
		twd := timeWithoutDate{7, 3}
		Convey("Its string representation should be zero-padded", func() {
			So(twd.String(), ShouldEqual, "07:03")
		})
	})
}

func TestRoute(t *testing.T) {
	Convey("Given route with several departures", t, func() {
		r := route{
			timeWithoutDate{10, 30},
			timeWithoutDate{11, 0},
			timeWithoutDate{12, 30},
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
				So(r[0].hour, ShouldEqual, 7)
				So(r[0].minute, ShouldEqual, 30)
				So(r[1].hour, ShouldEqual, 11)
				So(r[1].minute, ShouldEqual, 0)
				So(r[2].hour, ShouldEqual, 12)
				So(r[2].minute, ShouldEqual, 55)
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
				So(s.workDayRouteToOffice[0].String(), ShouldEqual, "07:30")
				So(s.workDayRouteToOffice[1].String(), ShouldEqual, "08:00")
				So(s.workDayRouteToOffice[2].String(), ShouldEqual, "20:00")
				So(s.workDayRouteToOffice[3].String(), ShouldEqual, "20:30")

				So(s.holidayRouteToOffice[0].String(), ShouldEqual, "10:30")

				So(s.workDayRouteFromOffice[0].String(), ShouldEqual, "08:20")
				So(s.workDayRouteFromOffice[1].String(), ShouldEqual, "08:50")
				So(s.workDayRouteFromOffice[2].String(), ShouldEqual, "20:20")
				So(s.workDayRouteFromOffice[3].String(), ShouldEqual, "20:50")

				So(s.holidayRouteFromOffice[0].String(), ShouldEqual, "18:00")
			})

			Convey("It should correctly find two best trips", func() {
				now, _ := time.Parse("15:04", "07:00")
				bestTrip, nextBestTrip := findBestTripMatches(now, s.workDayRouteToOffice)
				So(bestTrip.String(), ShouldEqual, "07:30")
				So(nextBestTrip.String(), ShouldEqual, "08:00")
			})

			Convey("It should correctly find one best trip", func() {
				now, _ := time.Parse("15:04", "20:15")
				bestTrip, nextBestTrip := findBestTripMatches(now, s.workDayRouteToOffice)
				So(bestTrip.String(), ShouldEqual, "20:30")
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
