package konturtransferbot

import (
	"testing"
	"time"

	"github.com/ghodss/yaml"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSchedule(t *testing.T) {
	Convey("Given correct YAML schedule", t, func() {
		sYaml := []byte(`WorkDayRouteToOffice:
  - "07:30"
  - "08:00"
  - "20:00"
  - "20:30"
SaturdayRouteToOffice:
  - "10:30"
WorkDayRouteFromOffice:
  - "08:20"
  - "08:50"
  - "20:20"
  - "20:50"
SaturdayRouteFromOffice:
  - "18:00"`)

		Convey("It should parse into a schedule structure", func() {
			s := Schedule{}
			err := yaml.Unmarshal(sYaml, &s)
			So(err, ShouldBeNil)

			Convey("Its entries should be correct and in the same order", func() {
				So(s.WorkDayRouteToOffice[0].Format("15:04"), ShouldEqual, "07:30")
				So(s.WorkDayRouteToOffice[1].Format("15:04"), ShouldEqual, "08:00")
				So(s.WorkDayRouteToOffice[2].Format("15:04"), ShouldEqual, "20:00")
				So(s.WorkDayRouteToOffice[3].Format("15:04"), ShouldEqual, "20:30")

				So(s.SaturdayRouteToOffice[0].Format("15:04"), ShouldEqual, "10:30")

				So(s.WorkDayRouteFromOffice[0].Format("15:04"), ShouldEqual, "08:20")
				So(s.WorkDayRouteFromOffice[1].Format("15:04"), ShouldEqual, "08:50")
				So(s.WorkDayRouteFromOffice[2].Format("15:04"), ShouldEqual, "20:20")
				So(s.WorkDayRouteFromOffice[3].Format("15:04"), ShouldEqual, "20:50")

				So(s.SaturdayRouteFromOffice[0].Format("15:04"), ShouldEqual, "18:00")
			})

			Convey("It should draw divider on weekdays", func() {
				So(s.GetFromOfficeText(time.Date(2017, 5, 2, 19, 45, 0, 0, &time.Location{})),
					ShouldEqual,
					"*Офис → Геологическая*\n\n08:20\n08:50\n———— сейчас 19:45 ————\n20:20\n20:50\n\nСубботний рейс в 18:00\n")
				So(s.GetToOfficeText(time.Date(2017, 5, 2, 19, 45, 0, 0, &time.Location{})),
					ShouldEqual,
					"*Геологическая → Офис*\n\n07:30\n08:00\n———— сейчас 19:45 ————\n20:00\n20:30\n\nСубботний рейс в 10:30\n")
			})

			Convey("It should not draw divider on weekends", func() {
				So(s.GetFromOfficeText(time.Date(2017, 5, 6, 19, 45, 0, 0, &time.Location{})),
					ShouldEqual,
					"*Офис → Геологическая*\n\n08:20\n08:50\n20:20\n20:50\n\nСубботний рейс в 18:00\n")
				So(s.GetToOfficeText(time.Date(2017, 5, 6, 19, 45, 0, 0, &time.Location{})),
					ShouldEqual,
					"*Геологическая → Офис*\n\n07:30\n08:00\n20:00\n20:30\n\nСубботний рейс в 10:30\n")
			})
		})
	})

	Convey("Given totally invalid YAML schedule", t, func() {
		sYaml := []byte(`	1123123`)
		Convey("It should not parse", func() {
			s := Schedule{}
			err := yaml.Unmarshal(sYaml, &s)
			So(err, ShouldNotBeNil)
		})
	})
}
