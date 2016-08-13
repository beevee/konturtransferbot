package main

import (
	"sort"
	"time"
)

type route []time.Time

func buildRoute(departures []string) (route, error) {
	result := make([]time.Time, len(departures))
	for index, departure := range departures {
		var err error
		result[index], err = time.Parse("15:04", departure)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (r route) String() string {
	var result string
	for _, departure := range r {
		result += departure.Format("15:04\n")
	}
	return result
}

func (r route) findBestTripMatches(now time.Time) (*time.Time, *time.Time) {
	bestDepartureMatch := sort.Search(len(r), func(i int) bool {
		return r[i].Hour() > now.Hour() || r[i].Hour() == now.Hour() && r[i].Minute() >= now.Minute()
	})
	var bestTrip, nextBestTrip *time.Time
	if bestDepartureMatch < len(r) {
		bestTrip = &r[bestDepartureMatch]
		if bestTrip.Hour()-now.Hour() >= 5 {
			return nil, nil
		}
		if bestDepartureMatch < len(r)-1 {
			nextBestTrip = &r[bestDepartureMatch+1]
		}
	}
	return bestTrip, nextBestTrip
}
