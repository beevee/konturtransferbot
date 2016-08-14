//+build !test

package main

import (
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
	theSchedule, err := buildSchedule(yamlData)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := telebot.NewBot(os.Getenv("KONTUR_TRANSFER_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	ekbTimezone, _ := time.LoadLocation("Asia/Yekaterinburg")
	var defaultMessageOptions = &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{buttonToOfficeLabel, buttonFromOfficeLabel},
				[]string{buttonScheduleToOfficeLabel, buttonScheduleFromOfficeLabel},
			},
			ResizeKeyboard: true,
		},
	}

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)

	for message := range messages {
		log.Infof("%s %s (username %s, chat id %d) said %s", message.Sender.FirstName, message.Sender.LastName, message.Sender.Username, message.Chat.ID, message.Text)

		now := time.Now().In(ekbTimezone)

		switch message.Text {
		case buttonToOfficeLabel:
			bot.SendMessage(message.Chat, theSchedule.getBestTripToOfficeText(now), defaultMessageOptions)
			continue

		case buttonFromOfficeLabel:
			bot.SendMessage(message.Chat, theSchedule.getBestTripFromOfficeText(now), defaultMessageOptions)
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
