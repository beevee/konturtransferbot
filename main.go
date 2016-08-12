package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"

	"github.com/tucnak/telebot"
	"gopkg.in/yaml.v2"
)

const (
	buttonToOfficeLabel           = "Хочу на работу"
	buttonFromOfficeLabel         = "Хочу домой"
	buttonScheduleToOfficeLabel   = "Все рейсы от Геологической"
	buttonScheduleFromOfficeLabel = "Все рейсы от офиса"
)

type timeWithoutDate struct {
	hour   int
	minute int
}

func (twd timeWithoutDate) String() string {
	return fmt.Sprintf("%02d:%02d", twd.hour, twd.minute)
}

type route []timeWithoutDate

func (r route) String() string {
	var result string
	for _, departure := range r {
		result += fmt.Sprintf("%s\n", departure)
	}
	return result
}

type schedule struct {
	workDayRouteToOffice   route
	workDayRouteFromOffice route
	holidayRouteToOffice   route
	holidayRouteFromOffice route
}

// ScheduleYaml - модель расписания для конфига
type ScheduleYaml struct {
	WorkDayRouteToOffice   []string `yaml:"WorkDayRouteToOffice"`
	WorkDayRouteFromOffice []string `yaml:"WorkDayRouteFromOffice"`
	HolidayRouteToOffice   []string `yaml:"HolidayRouteToOffice"`
	HolidayRouteFromOffice []string `yaml:"HolidayRouteFromOffice"`
}

func buildRoute(departures []string) route {
	result := make([]timeWithoutDate, len(departures))
	for index, departure := range departures {
		twd := timeWithoutDate{}
		fmt.Sscanf(departure, "%d:%d", &twd.hour, &twd.minute)
		result[index] = twd
	}
	return result
}

func buildSchedule() schedule {
	data, _ := ioutil.ReadFile("schedule.yml")
	scheduleYaml := ScheduleYaml{}
	err := yaml.Unmarshal([]byte(data), &scheduleYaml)
	if err != nil {
		log.Print(err)
	}
	result := schedule{
		workDayRouteToOffice:   buildRoute(scheduleYaml.WorkDayRouteToOffice),
		workDayRouteFromOffice: buildRoute(scheduleYaml.WorkDayRouteFromOffice),
		holidayRouteToOffice:   buildRoute(scheduleYaml.HolidayRouteToOffice),
		holidayRouteFromOffice: buildRoute(scheduleYaml.HolidayRouteFromOffice),
	}
	return result
}

func findBestTripMatches(now time.Time, r route) (*timeWithoutDate, *timeWithoutDate) {
	bestDepartureMatch := sort.Search(len(r), func(i int) bool {
		return r[i].hour > now.Hour() || r[i].hour == now.Hour() && r[i].minute >= now.Minute()
	})
	var bestTrip, nextBestTrip *timeWithoutDate
	if bestDepartureMatch < len(r) {
		bestTrip = &r[bestDepartureMatch]
		if bestDepartureMatch < len(r)-1 {
			nextBestTrip = &r[bestDepartureMatch+1]
		}
	}
	return bestTrip, nextBestTrip
}

func main() {
	theSchedule := buildSchedule()

	bot, err := telebot.NewBot(os.Getenv("KONTUR_TRANSFER_BOT_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	var defaultMessageOptions = &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{buttonToOfficeLabel, buttonFromOfficeLabel},
				[]string{buttonScheduleToOfficeLabel, buttonScheduleFromOfficeLabel},
			},
		},
	}
	var monetizationMessage = "Монетизация! Промокод Gett на первую поездку - GTFUNKP, Яндекс.Такси - daf3qsau, Uber - ykt6m, Wheely - MHPRL."

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)

	for message := range messages {
		log.Printf("%s %s (username %s) said %s", message.Sender.FirstName, message.Sender.LastName, message.Sender.Username, message.Text)

		var reply string
		var currentRoute route

		ekbTimezone, _ := time.LoadLocation("Asia/Yekaterinburg")
		now := time.Now().In(ekbTimezone)
		isWorkDay := now.Weekday() != time.Sunday && now.Weekday() != time.Saturday

		switch message.Text {
		case buttonToOfficeLabel:
			if isWorkDay {
				currentRoute = theSchedule.workDayRouteToOffice
			} else {
				currentRoute = theSchedule.holidayRouteToOffice
			}
			bestTrip, nextBestTrip := findBestTripMatches(now, currentRoute)
			if bestTrip != nil {
				reply = fmt.Sprintf("Ближайший дежурный рейс от Геологической будет в %s.", bestTrip)
				if nextBestTrip != nil {
					reply += fmt.Sprintf(" Следующий - в %s.", nextBestTrip)
				}
			} else {
				reply = "Сегодня уехать на работу уже не получится :( Удачи завтра!"
			}
			bot.SendMessage(message.Chat, reply, defaultMessageOptions)
			continue

		case buttonFromOfficeLabel:
			if isWorkDay {
				currentRoute = theSchedule.workDayRouteFromOffice
			} else {
				currentRoute = theSchedule.holidayRouteFromOffice
			}
			bestTrip, nextBestTrip := findBestTripMatches(now, currentRoute)
			if bestTrip != nil {
				reply = fmt.Sprintf("Ближайший дежурный рейс от офиса будет в %s.", bestTrip)
				if nextBestTrip != nil {
					reply += fmt.Sprintf(" Следующий - в %s.", nextBestTrip)
				} else {
					reply += " Это последний на сегодня рейс, дальше - только на такси. " + monetizationMessage
				}
			} else {
				reply = "Сегодня уехать домой уже не получится :( Придется остаться на ночь или ехать на такси. " + monetizationMessage
			}
			bot.SendMessage(message.Chat, reply, defaultMessageOptions)
			continue

		case buttonScheduleToOfficeLabel:
			bot.SendMessage(message.Chat, fmt.Sprintf("Дежурные рейсы в будни:\n%s", theSchedule.workDayRouteToOffice), defaultMessageOptions)
			bot.SendMessage(message.Chat, fmt.Sprintf("Дежурные рейсы в выходные:\n%s", theSchedule.holidayRouteToOffice), defaultMessageOptions)
			continue

		case buttonScheduleFromOfficeLabel:
			bot.SendMessage(message.Chat, fmt.Sprintf("Дежурные рейсы в будни:\n%s", theSchedule.workDayRouteFromOffice), defaultMessageOptions)
			bot.SendMessage(message.Chat, fmt.Sprintf("Дежурные рейсы в выходные:\n%s", theSchedule.holidayRouteFromOffice), defaultMessageOptions)
			continue
		}
		bot.SendMessage(message.Chat, "Привет! Я могу подсказать расписание трансфера по дежурному маршруту.", defaultMessageOptions)
	}
}
