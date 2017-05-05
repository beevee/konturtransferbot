package konturtransferbot

import (
	"fmt"
	"time"
)

// Route is a sorted sequence of departure times
type Route []Departure

// Departure is a single departure time
type Departure struct {
	time.Time
}

// UnmarshalJSON is a custom unmarshaler function for time, which works for both JSON and YAML
func (d *Departure) UnmarshalJSON(departure []byte) error {
	var err error
	d.Time, err = time.Parse("\"15:04\"", string(departure))
	if err != nil {
		return err
	}
	return nil
}

func (r Route) String() string {
	var result string
	for _, departure := range r {
		result += departure.Format("15:04\n")
	}
	return result
}

// StringWithDivider prints current time inside a route schedule
func (r Route) StringWithDivider(now time.Time) string {
	nowReset := time.Date(0, time.January, 1, now.Hour(), now.Minute(), 0, 0, &time.Location{})

	var result string
	for i := range r {
		if i > 0 && (r[i].After(nowReset) || r[i].Equal(nowReset)) && r[i-1].Before(nowReset) {
			result += fmt.Sprintf("———— сейчас %s ————\n", nowReset.Format("15:04"))
		}
		result += r[i].Format("15:04\n")
	}
	return result
}
