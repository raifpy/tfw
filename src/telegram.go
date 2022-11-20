package src

import (
	"time"

	"gopkg.in/telebot.v3"
)

type TelegramArray []*Telegram

func (ta TelegramArray) FindbyId(id string) (t *Telegram) {
	for _, v := range ta {
		if v.Id == id {
			return v
		}

	}
	return nil
}

type Telegram struct {
	*TMessageF
	*telebot.Bot
	Id string
}

type TelegramOptions struct {
	Id       string
	BotToken string
	BotAPI   string
}

func NewTelegram(o TelegramOptions, feed *TMessageF) (t *Telegram, err error) {
	t = &Telegram{
		Id:        o.Id,
		TMessageF: feed,
	}
	t.Bot, err = telebot.NewBot(telebot.Settings{
		URL:   o.BotAPI,
		Token: o.BotToken,
		//Verbose: true,
		OnError: func(err error, ctx telebot.Context) {
			t.Log.Errorf("telegram bot %s:%s error: %v chat: %v context: %v", t.Id, t.Me.Username, err, ctx.Chat(), ctx)
		},
		Poller: &telebot.LongPoller{Timeout: time.Second * 5},
	})

	go t.Bot.Start()
	return
}
