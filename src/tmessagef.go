package src

import (
	"errors"

	"gopkg.in/telebot.v3"
)

type Options struct {
	IdWhlitelist    []int64
	TelegramOptions []TelegramOptions
	Log             TMessageFLoggerOption
	Database        DatabaseOptions
}

type TMessageF struct {
	Bots       TelegramArray
	Log        *TMessageFLogger
	Database   *Database
	Root       *Telegram
	RoortGroup *telebot.Group
	Options    *Options
}

func NewTMessageF(Options *Options) (t *TMessageF, err error) {

	t = &TMessageF{
		Bots:    TelegramArray{},
		Log:     NewTMessageFLogger(Options.Log),
		Options: Options,
	}
	if t.Database, err = NewDatabase(Options.Database); err != nil {
		return
	}

	for _, to := range Options.TelegramOptions {
		_bot, err := NewTelegram(to, t)
		if err != nil {
			return t, err
		}

		t.Bots = append(t.Bots, _bot)
		/*if t.Bots[index], err = NewTelegram(to, t); err != nil {
			return t, err
		}*/
	}

	if err = t.Log.Fabric(); err != nil {
		return
	}
	err = t.init()

	return

}

func (t *TMessageF) init() error {
	t.Root = t.Bots.FindbyId("root")
	if t.Root == nil {
		return errors.New("searching bot which has a root id")
	}

	t.Admin()
	t.HandleTelegramActions()

	return nil

}
