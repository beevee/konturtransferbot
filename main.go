package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/tucnak/telebot"
)

type busStop int

const (
	city busStop = iota
	office
)

type timeWithoutDate struct {
	hour   int
	minute int
}

type trip struct {
	startingPoint busStop
	departureTime timeWithoutDate
	isWeekend     bool
}

var trips = [...]trip{
	trip{city, timeWithoutDate{7, 30}, false},
	trip{city, timeWithoutDate{8, 0}, false},
	trip{city, timeWithoutDate{8, 30}, false},
	trip{city, timeWithoutDate{9, 0}, false},
	trip{city, timeWithoutDate{9, 30}, false},
	trip{city, timeWithoutDate{10, 0}, false},
	trip{city, timeWithoutDate{10, 40}, false},
	trip{city, timeWithoutDate{11, 10}, false},
	trip{city, timeWithoutDate{11, 40}, false},
	trip{city, timeWithoutDate{12, 20}, false},
	trip{city, timeWithoutDate{12, 50}, false},
	trip{city, timeWithoutDate{13, 20}, false},
	trip{city, timeWithoutDate{14, 0}, false},
	trip{city, timeWithoutDate{14, 30}, false},
	trip{city, timeWithoutDate{15, 0}, false},
	trip{city, timeWithoutDate{15, 40}, false},
	trip{city, timeWithoutDate{16, 10}, false},
	trip{city, timeWithoutDate{16, 40}, false},
	trip{city, timeWithoutDate{17, 20}, false},
	trip{city, timeWithoutDate{17, 50}, false},
	trip{city, timeWithoutDate{18, 20}, false},
	trip{city, timeWithoutDate{19, 0}, false},
	trip{city, timeWithoutDate{19, 30}, false},
	trip{city, timeWithoutDate{20, 0}, false},
	trip{city, timeWithoutDate{20, 30}, false},

	trip{city, timeWithoutDate{10, 30}, true},

	trip{office, timeWithoutDate{8, 20}, false},
	trip{office, timeWithoutDate{8, 50}, false},
	trip{office, timeWithoutDate{9, 20}, false},
	trip{office, timeWithoutDate{9, 50}, false},
	trip{office, timeWithoutDate{10, 20}, false},
	trip{office, timeWithoutDate{10, 50}, false},
	trip{office, timeWithoutDate{11, 30}, false},
	trip{office, timeWithoutDate{12, 0}, false},
	trip{office, timeWithoutDate{12, 30}, false},
	trip{office, timeWithoutDate{13, 10}, false},
	trip{office, timeWithoutDate{13, 40}, false},
	trip{office, timeWithoutDate{14, 10}, false},
	trip{office, timeWithoutDate{14, 50}, false},
	trip{office, timeWithoutDate{15, 20}, false},
	trip{office, timeWithoutDate{15, 50}, false},
	trip{office, timeWithoutDate{16, 30}, false},
	trip{office, timeWithoutDate{17, 0}, false},
	trip{office, timeWithoutDate{17, 30}, false},
	trip{office, timeWithoutDate{18, 10}, false},
	trip{office, timeWithoutDate{18, 40}, false},
	trip{office, timeWithoutDate{19, 10}, false},
	trip{office, timeWithoutDate{19, 50}, false},
	trip{office, timeWithoutDate{20, 20}, false},
	trip{office, timeWithoutDate{20, 50}, false},

	trip{office, timeWithoutDate{18, 0}, true},
}

func findBestTripMatches(now time.Time, departurePoint busStop) (*trip, *trip) {
	isWeekend := now.Weekday() == time.Sunday || now.Weekday() == time.Saturday

	filteredTrips := make([]trip, 0, len(trips))
	for _, t := range trips {
		if t.isWeekend == isWeekend && t.startingPoint == departurePoint {
			filteredTrips = append(filteredTrips, t)
		}
	}

	bestDepartureMatch := sort.Search(len(filteredTrips), func(i int) bool {
		return filteredTrips[i].departureTime.hour > now.Hour() || filteredTrips[i].departureTime.hour == now.Hour() && filteredTrips[i].departureTime.minute >= now.Minute()
	})
	var bestTrip, nextBestTrip *trip
	if bestDepartureMatch < len(filteredTrips) {
		bestTrip = &filteredTrips[bestDepartureMatch]
		if bestDepartureMatch < len(filteredTrips)-1 {
			nextBestTrip = &filteredTrips[bestDepartureMatch+1]
		}
	}
	return bestTrip, nextBestTrip
}

func main() {
	bot, err := telebot.NewBot(os.Getenv("KONTUR_TRANSFER_BOT_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	var defaultMessageOptions = &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{"Хочу на работу", "Хочу домой"},
			},
		},
	}

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)

	for message := range messages {
		log.Printf("%s %s (username %s) said %s", message.Sender.FirstName, message.Sender.LastName, message.Sender.Username, message.Text)

		var reply string
		ekbTimezone, _ := time.LoadLocation("Asia/Yekaterinburg")
		now := time.Now().In(ekbTimezone)

		switch message.Text {
		case "Хочу на работу":
			trip1, trip2 := findBestTripMatches(now, city)
			if trip1 != nil {
				reply = fmt.Sprintf("Ближайший дежурный рейс от Геологической будет в %02d:%02d.", trip1.departureTime.hour, trip1.departureTime.minute)
				if trip2 != nil {
					reply += fmt.Sprintf(" Следующий - в %02d:%02d.", trip2.departureTime.hour, trip2.departureTime.minute)
				}
			} else {
				reply = "Сегодня уехать на работу уже не получится :( Удачи завтра!"
			}
			bot.SendMessage(message.Chat, reply, defaultMessageOptions)
			continue
		case "Хочу домой":
			bestTrip, nextBestTrip := findBestTripMatches(now, office)
			if bestTrip != nil {
				reply = fmt.Sprintf("Ближайший дежурный рейс от офиса будет в %02d:%02d.", bestTrip.departureTime.hour, bestTrip.departureTime.minute)
				if nextBestTrip != nil {
					reply += fmt.Sprintf(" Следующий - в %02d:%02d.", nextBestTrip.departureTime.hour, nextBestTrip.departureTime.minute)
				}
			} else {
				reply = "Сегодня уехать домой уже не получится :( Придется остаться на ночь."
			}
			bot.SendMessage(message.Chat, reply, defaultMessageOptions)
			continue
		}
		bot.SendMessage(message.Chat, "Привет! Я понимаю только две команды.", defaultMessageOptions)
	}
}
