package konturtransferbot

import (
	"testing"
	"time"

	"github.com/ghodss/yaml"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRoute(t *testing.T) {
	Convey("Given route with several departures", t, func() {
		r := Route{
			Departure{time.Date(0, 0, 0, 10, 30, 0, 0, &time.Location{})},
			Departure{time.Date(0, 0, 0, 11, 0, 0, 0, &time.Location{})},
			Departure{time.Date(0, 0, 0, 12, 30, 0, 0, &time.Location{})},
		}
		Convey("Its string representation should contain all of them", func() {
			So(r.String(), ShouldEqual, "10:30\n11:00\n12:30\n")
		})
	})

	Convey("Given string representation of a route", t, func() {
		rYaml := []byte("- 7:30\n- 11:00\n- 12:55")

		Convey("It should parse correctly", func() {
			r := Route{}
			err := yaml.Unmarshal(rYaml, &r)
			So(err, ShouldBeNil)

			Convey("Its numerical representation should be correct and in the same order", func() {
				So(r[0].Hour(), ShouldEqual, 7)
				So(r[0].Minute(), ShouldEqual, 30)
				So(r[1].Hour(), ShouldEqual, 11)
				So(r[1].Minute(), ShouldEqual, 0)
				So(r[2].Hour(), ShouldEqual, 12)
				So(r[2].Minute(), ShouldEqual, 55)
			})

			Convey("It should correctly find no best trips", func() {
				now, _ := time.Parse("15:04", "17:00")
				bestTrip, nextBestTrip := r.findBestTripMatches(now)
				So(bestTrip, ShouldBeNil)
				So(nextBestTrip, ShouldBeNil)
			})

			Convey("It should correctly find two best trips", func() {
				now, _ := time.Parse("15:04", "07:00")
				bestTrip, nextBestTrip := r.findBestTripMatches(now)
				So(bestTrip.Format("15:04"), ShouldEqual, "07:30")
				So(nextBestTrip.Format("15:04"), ShouldEqual, "11:00")
			})

			Convey("It should correctly find one best trip", func() {
				now, _ := time.Parse("15:04", "12:15")
				bestTrip, nextBestTrip := r.findBestTripMatches(now)
				So(bestTrip.Format("15:04"), ShouldEqual, "12:55")
				So(nextBestTrip, ShouldBeNil)
			})

			Convey("It should not offer trips more than (roughly) 5 hours in advance", func() {
				now, _ := time.Parse("15:04", "01:25")
				bestTrip, nextBestTrip := r.findBestTripMatches(now)
				So(bestTrip, ShouldBeNil)
				So(nextBestTrip, ShouldBeNil)
			})

		})
	})

	Convey("Given incorrect string representation of a route", t, func() {
		rYaml := []byte("- NNNNNOOOOOOOOOOO\n- 11:00\n- 12:55")

		Convey("It should not parse", func() {
			r := Route{}
			err := yaml.Unmarshal(rYaml, &r)
			So(err, ShouldNotBeNil)
		})
	})
}
