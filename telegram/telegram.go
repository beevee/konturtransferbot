package telegram

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/beevee/konturtransferbot"

	"golang.org/x/sync/errgroup"
	"gopkg.in/tucnak/telebot.v2"
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
	ProxyURL      *url.URL
	Logger        konturtransferbot.Logger
	telebot       *telebot.Bot
	ctx           context.Context
	ctxCancel     context.CancelFunc
	ctxGroup      *errgroup.Group
}

// Start initializes Telegram request loop
func (b *Bot) Start() error {
	transport := &http.Transport{}
	if b.ProxyURL != nil {
		transport.Proxy = http.ProxyURL(b.ProxyURL)
		b.Logger.Log("msg", "working via proxy", "proxy_host", b.ProxyURL.Host,
			"proxy_user", b.ProxyURL.User.Username(), "proxy_proto", b.ProxyURL.Scheme)
	}
	client := http.DefaultClient
	client.Transport = transport

	var err error
	b.telebot, err = telebot.NewBot(telebot.Settings{
		Token:  b.TelegramToken,
		Poller: &telebot.LongPoller{Timeout: 1 * time.Second},
		Client: client,
		Reporter: func(err error) {
			b.Logger.Log("msg", "telebot error", "error", err)
		},
	})
	if err != nil {
		return err
	}

	b.telebot.Handle(telebot.OnText, b.handleMessage)

	b.ctx, b.ctxCancel = context.WithCancel(context.Background())
	b.ctxGroup, b.ctx = errgroup.WithContext(b.ctx)

	go b.telebot.Start()

	return nil
}

// Stop gracefully finishes request loop
func (b *Bot) Stop() error {
	b.ctxCancel()
	err := b.ctxGroup.Wait()

	b.telebot.Stop()

	return err
}

func (b *Bot) handleMessage(message *telebot.Message) {
	now := time.Now().In(b.Timezone)

	b.Logger.Log("msg", "message received", "firstname", message.Sender.FirstName, "lastname",
		message.Sender.LastName, "username", message.Sender.Username, "chatid", message.Chat.ID, "text", message.Text)

	switch message.Text {
	case buttonToOfficeLabel:
		messageNow, messageLater := b.Schedule.GetToOfficeText(now)
		b.sendAndDelayReply(message.Chat, messageNow, messageLater)
		return

	case buttonFromOfficeLabel:
		messageNow, messageLater := b.Schedule.GetFromOfficeText(now)
		b.sendAndDelayReply(message.Chat, messageNow, messageLater)
		return
	}

	_, err := b.telebot.Send(
		message.Chat,
		"Привет! Я могу подсказать расписание трансфера по дежурному маршруту.",
		&telebot.SendOptions{
			ReplyMarkup: &telebot.ReplyMarkup{
				ReplyKeyboard: [][]telebot.ReplyButton{
					{
						telebot.ReplyButton{Text: buttonToOfficeLabel},
						telebot.ReplyButton{Text: buttonFromOfficeLabel},
					},
				},
				ResizeReplyKeyboard: true,
			},
			ParseMode: telebot.ModeMarkdown,
		})

	if err != nil {
		b.Logger.Log("msg", "error sending message", "error", err)
	}
}

func (b *Bot) sendAndDelayReply(chat *telebot.Chat, messageNow string, messageLater string) {
	message, err := b.telebot.Send(chat, messageNow, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	if err != nil {
		b.Logger.Log("msg", "error sending message", "chatid", chat.ID, "text", messageNow, "error", err)
		return
	}
	b.Logger.Log("msg", "message sent", "chatid", chat.ID, "messageid", message.ID, "text", messageNow)

	if messageLater != "" {
		b.ctxGroup.Go(func() error {
			timer := time.NewTimer(5 * time.Minute)
			select {
			case <-timer.C:
				break
			case <-b.ctx.Done():
				break
			}
			_, errEdit := b.telebot.Edit(message, messageLater, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
			if errEdit != nil {
				b.Logger.Log("msg", "error editing message", "chatid", chat.ID, "messageid", message.ID, "text", messageLater, "error", err)
				return errEdit
			}
			b.Logger.Log("msg", "message edited", "chatid", chat.ID, "messageid", message.ID, "text", messageLater)
			return nil
		})
	}
}
