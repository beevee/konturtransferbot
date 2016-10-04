package telegram

import (
	"time"

	"github.com/beevee/konturtransferbot"

	"github.com/tucnak/telebot"
	"gopkg.in/tomb.v2"
)

const (
	buttonToOfficeLabel           = "Хочу на работу"
	buttonFromOfficeLabel         = "Хочу домой"
	buttonScheduleToOfficeLabel   = "Все рейсы на работу"
	buttonScheduleFromOfficeLabel = "Все рейсы домой"
)

// Bot handles communication with Telegram users
type Bot struct {
	Schedule              konturtransferbot.Schedule
	TelegramToken         string
	Timezone              *time.Location
	Logger                konturtransferbot.Logger
	defaultMessageOptions *telebot.SendOptions
	telebot               *telebot.Bot
	tomb                  tomb.Tomb
}

// Start initializes Telegram request loop
func (b *Bot) Start() error {
	var err error
	b.telebot, err = telebot.NewBot(b.TelegramToken)
	if err != nil {
		return err
	}

	b.defaultMessageOptions = &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{buttonToOfficeLabel, buttonFromOfficeLabel},
				[]string{buttonScheduleToOfficeLabel, buttonScheduleFromOfficeLabel},
			},
			ResizeKeyboard: true,
		},
	}

	messages := make(chan telebot.Message)
	b.telebot.Listen(messages, 1*time.Second)

	b.tomb.Go(func() error {
		for {
			select {
			case message := <-messages:
				if err := b.handleMessage(message); err != nil {
					b.Logger.Log("msg", "error sending message", "error", err)
				}
			case <-b.tomb.Dying():
				return nil
			}
		}
	})

	return nil
}

// Stop gracefully finishes request loop
func (b *Bot) Stop() error {
	b.tomb.Kill(nil)
	return b.tomb.Wait()
}

func (b *Bot) handleMessage(message telebot.Message) error {
	now := time.Now().In(b.Timezone)

	b.Logger.Log("msg", "message received", "firstname", message.Sender.FirstName, "lastname", message.Sender.LastName, "username", message.Sender.Username, "chatid", message.Chat.ID, "text", message.Text)

	switch message.Text {
	case buttonToOfficeLabel:
		return b.telebot.SendMessage(message.Chat, b.Schedule.GetBestTripToOfficeText(now), b.defaultMessageOptions)

	case buttonFromOfficeLabel:
		return b.telebot.SendMessage(message.Chat, b.Schedule.GetBestTripFromOfficeText(now), b.defaultMessageOptions)

	case buttonScheduleToOfficeLabel:
		for _, text := range b.Schedule.GetFullToOfficeTexts() {
			if err := b.telebot.SendMessage(message.Chat, text, b.defaultMessageOptions); err != nil {
				return err
			}
		}
		return nil

	case buttonScheduleFromOfficeLabel:
		for _, text := range b.Schedule.GetFullFromOfficeTexts() {
			if err := b.telebot.SendMessage(message.Chat, text, b.defaultMessageOptions); err != nil {
				return err
			}
		}
		return nil
	}
	return b.telebot.SendMessage(message.Chat, "Привет! Я могу подсказать расписание трансфера по дежурному маршруту.", b.defaultMessageOptions)
}
