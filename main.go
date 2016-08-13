package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/op/go-logging"
	"github.com/tucnak/telebot"
)

var log *logging.Logger

const (
	buttonToOfficeLabel           = "Хочу на работу"
	buttonFromOfficeLabel         = "Хочу домой"
	buttonScheduleToOfficeLabel   = "Все рейсы на работу"
	buttonScheduleFromOfficeLabel = "Все рейсы домой"
)

func main() {
	log = logging.MustGetLogger("kontur_transfer_bot")
	logFile, _ := os.OpenFile("kontur_transfer_bot.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	logBackend := logging.NewLogBackend(logFile, "", 0)
	logging.SetFormatter(logging.MustStringFormatter("%{time:2006-01-02 15:04:05}\t%{level}\t%{message}"))
	logging.SetBackend(logBackend)

	yamlData, err := ioutil.ReadFile("schedule.yml")
	if err != nil {
		log.Fatal(err)
	}
	theSchedule := buildSchedule(yamlData)

	bot, err := telebot.NewBot(os.Getenv("KONTUR_TRANSFER_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	var defaultMessageOptions = &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{buttonToOfficeLabel, buttonFromOfficeLabel},
				[]string{buttonScheduleToOfficeLabel, buttonScheduleFromOfficeLabel},
			},
			ResizeKeyboard: true,
		},
	}
	var monetizationMessage = "Монетизация! Промокод Gett на первую поездку - GTFUNKP, Яндекс.Такси - daf3qsau, Uber - ykt6m, Wheely - MHPRL."

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)

	for message := range messages {
		log.Infof("%s %s (username %s, chat id %d) said %s", message.Sender.FirstName, message.Sender.LastName, message.Sender.Username, message.Chat.ID, message.Text)

		var reply string
		var currentRoute route

		ekbTimezone, _ := time.LoadLocation("Asia/Yekaterinburg")
		now := time.Now().In(ekbTimezone)

		switch message.Text {
		case buttonToOfficeLabel:
			bestTrip, nextBestTrip := theSchedule.findCorrectRoute(now, true).findBestTripMatches(now)
			if bestTrip != nil {
				reply = fmt.Sprintf("Ближайший дежурный рейс от Геологической будет в %s.", bestTrip.Format("15:04"))
				if nextBestTrip != nil {
					reply += fmt.Sprintf(" Следующий - в %s.", nextBestTrip.Format("15:04"))
				}
			} else {
				reply = "В ближайшие несколько часов уехать на работу на трансфере не получится. Лучше лечь поспать и поехать с утра. Первые рейсы от Геологической: "
				nextDay := now.Add(24 * time.Hour)
				currentRoute = theSchedule.findCorrectRoute(nextDay, true)
				for index, trip := range currentRoute {
					if trip.Hour() >= 12 || index >= len(currentRoute)-1 {
						reply += fmt.Sprintf("%s.", trip.Format("15:04"))
						break
					}
					reply += fmt.Sprintf("%s, ", trip.Format("15:04"))
				}
			}
			bot.SendMessage(message.Chat, reply, defaultMessageOptions)
			continue

		case buttonFromOfficeLabel:
			bestTrip, nextBestTrip := theSchedule.findCorrectRoute(now, false).findBestTripMatches(now)
			if bestTrip != nil {
				reply = fmt.Sprintf("Ближайший дежурный рейс от офиса будет в %s.", bestTrip.Format("15:04"))
				if nextBestTrip != nil {
					reply += fmt.Sprintf(" Следующий - в %s.", nextBestTrip.Format("15:04"))
				} else {
					reply += " Это последний на сегодня рейс, дальше - только на такси. " + monetizationMessage
				}
			} else {
				reply = "В ближайшие несколько часов уехать домой на трансфере не получится :( Придется остаться в офисе или ехать на такси. " + monetizationMessage
			}
			bot.SendMessage(message.Chat, reply, defaultMessageOptions)
			continue

		case buttonScheduleToOfficeLabel:
			for _, text := range theSchedule.getFullToOfficeTexts() {
				bot.SendMessage(message.Chat, text, defaultMessageOptions)
			}
			continue

		case buttonScheduleFromOfficeLabel:
			for _, text := range theSchedule.getFullFromOfficeTexts() {
				bot.SendMessage(message.Chat, text, defaultMessageOptions)
			}
			continue
		}
		bot.SendMessage(message.Chat, "Привет! Я могу подсказать расписание трансфера по дежурному маршруту.", defaultMessageOptions)
	}
}
