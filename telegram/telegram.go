package telegram

import (
	"time"

	"github.com/beevee/konturtransferbot"

	"github.com/vlad-lukyanov/telebot"
	"gopkg.in/tomb.v2"
)

const (
	buttonToOfficeLabel   = "🛠 В офис"
	buttonFromOfficeLabel = "🍻 Из офиса"
)

// Bot handles communication with Telegram users
type Bot struct {
	Schedule      konturtransferbot.Schedule
	TelegramToken string
	Timezone      *time.Location
	Logger        konturtransferbot.Logger
	telebot       *telebot.Bot
	tomb          tomb.Tomb
}

// Start initializes Telegram request loop
func (b *Bot) Start() error {
	var err error
	b.telebot, err = telebot.NewBot(b.TelegramToken)
	if err != nil {
		return err
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
		messageNow, messageLater := b.Schedule.GetToOfficeText(now)
		return b.sendAndDelayReply(message.Chat, messageNow, messageLater)

	case buttonFromOfficeLabel:
		messageNow, messageLater := b.Schedule.GetFromOfficeText(now)
		return b.sendAndDelayReply(message.Chat, messageNow, messageLater)
	}

	_, err := b.telebot.SendMessage(
		message.Chat,
		"Привет! Я могу подсказать расписание трансфера по дежурному маршруту.",
		&telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				CustomKeyboard: [][]string{
					[]string{buttonToOfficeLabel, buttonFromOfficeLabel},
				},
				ResizeKeyboard: true,
			},
			ParseMode: telebot.ModeMarkdown,
		})
	return err
}

func (b *Bot) sendAndDelayReply(chat telebot.Chat, messageNow string, messageLater string) error {
	message, err := b.telebot.SendMessage(chat, messageNow, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	if err != nil {
		b.Logger.Log("msg", "error sending message", "chatid", chat.ID, "messageid", message.ID, "text", messageNow, "error", err)
		return err
	}
	b.Logger.Log("msg", "message sent", "chatid", chat.ID, "messageid", message.ID, "text", messageNow)

	if messageLater != "" {
		b.tomb.Go(func() error {
			timer := time.NewTimer(5 * time.Minute)
			select {
			case <-timer.C:
				break
			case <-b.tomb.Dying():
				break
			}
			_, errEdit := b.telebot.EditMessageText(chat, message.ID, messageLater, nil)
			if errEdit != nil {
				b.Logger.Log("msg", "error editing message", "chatid", chat.ID, "messageid", message.ID, "text", messageLater, "error", err)
				return errEdit
			}
			b.Logger.Log("msg", "message edited", "chatid", chat.ID, "messageid", message.ID, "text", messageLater)
			return nil
		})
	}

	return nil
}
