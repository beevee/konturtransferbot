package konturtransferbot

import (
	"sort"
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

func (r Route) findBestTripMatches(now time.Time) (*Departure, *Departure) {
	bestDepartureMatch := sort.Search(len(r), func(i int) bool {
		return r[i].Hour() > now.Hour() || r[i].Hour() == now.Hour() && r[i].Minute() >= now.Minute()
	})
	var bestTrip, nextBestTrip *Departure
	if bestDepartureMatch < len(r) {
		bestTrip = &r[bestDepartureMatch]
		if bestTrip.Hour()-now.Hour() >= 6 {
			return nil, nil
		}
		if bestDepartureMatch < len(r)-1 {
			nextBestTrip = &r[bestDepartureMatch+1]
		}
	}
	return bestTrip, nextBestTrip
}
