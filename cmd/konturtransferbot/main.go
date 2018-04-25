package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/beevee/konturtransferbot"
	"github.com/beevee/konturtransferbot/telegram"

	"github.com/ghodss/yaml"
	"github.com/go-kit/kit/log"
	"github.com/jessevdk/go-flags"
)

func main() {
	var opts struct {
		TelegramToken string `short:"t" long:"telegram-token" description:"@KonturTransferBot Telegram token" env:"KONTUR_TRANSFER_BOT_TOKEN"`
		ScheduleYaml  string `short:"s" long:"schedule-yaml" default:"schedule.yml" description:"YAML file with schedule" env:"KONTUR_TRANSFER_SCHEDULE_YAML"`
		Timezone      string `short:"z" long:"timezone" default:"Asia/Yekaterinburg" description:"Local timezone" env:"KONTUR_TRANSFER_BOT_TIMEZONE"`
	}

	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(0)
	}

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	yamlSchedule, err := ioutil.ReadFile(opts.ScheduleYaml)
	if err != nil {
		logger.Log("msg", "failed to load schedule YAML file", "error", err)
		os.Exit(1)
	}
	schedule := konturtransferbot.Schedule{}
	if err = yaml.Unmarshal([]byte(yamlSchedule), &schedule); err != nil {
		logger.Log("msg", "failed to build schedule from YAML", "error", err)
		os.Exit(1)
	}

	tz, err := time.LoadLocation(opts.Timezone)
	if err != nil {
		logger.Log("msg", "failed to recognize timezone", "error", err)
		os.Exit(1)
	}
	bot := &telegram.Bot{
		Schedule:      schedule,
		TelegramToken: opts.TelegramToken,
		Timezone:      tz,
		Logger:        log.With(logger, "component", "telegram"),
	}

	logger.Log("msg", "starting Telegram bot")
	if err := bot.Start(); err != nil {
		logger.Log("msg", "error starting Telegram bot", "error", err)
		os.Exit(1)
	}
	logger.Log("msg", "started Telegram bot")

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	logger.Log("msg", "received signal", "signal", <-signalChannel)

	logger.Log("msg", "stopping Telegram bot")
	if err := bot.Stop(); err != nil {
		logger.Log("msg", "error stopping Telegram bot", "error", err)
	}
	logger.Log("msg", "stopped Telegram bot")
}
