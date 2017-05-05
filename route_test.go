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
			Departure{time.Date(0, 1, 1, 10, 30, 0, 0, &time.Location{})},
			Departure{time.Date(0, 1, 1, 11, 0, 0, 0, &time.Location{})},
			Departure{time.Date(0, 1, 1, 12, 30, 0, 0, &time.Location{})},
		}

		Convey("Its string representation should contain all of them", func() {
			So(r.String(), ShouldEqual, "10:30\n11:00\n12:30\n")
		})

		Convey("Divider should be placed correctly inside a schedule", func() {
			now := time.Date(0, 1, 1, 10, 45, 0, 0, &time.Location{})
			So(r.StringWithDivider(now), ShouldEqual, "10:30\n———— сейчас 10:45 ————\n11:00\n12:30\n")
		})

		Convey("Divider should be placed before departure time if current time is equal to this departure time", func() {
			now := time.Date(0, 1, 1, 11, 0, 0, 0, &time.Location{})
			So(r.StringWithDivider(now), ShouldEqual, "10:30\n———— сейчас 11:00 ————\n11:00\n12:30\n")
		})

		Convey("Divider calculation should not consider date, only time", func() {
			now := time.Date(2017, 5, 1, 10, 45, 0, 0, &time.Location{})
			So(r.StringWithDivider(now), ShouldEqual, "10:30\n———— сейчас 10:45 ————\n11:00\n12:30\n")
		})

		Convey("No divider should be placed before first departure", func() {
			now := time.Date(0, 1, 1, 9, 45, 0, 0, &time.Location{})
			So(r.StringWithDivider(now), ShouldEqual, "10:30\n11:00\n12:30\n")
		})

		Convey("No divider should be placed after last departure", func() {
			now := time.Date(0, 1, 1, 19, 45, 0, 0, &time.Location{})
			So(r.StringWithDivider(now), ShouldEqual, "10:30\n11:00\n12:30\n")
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
