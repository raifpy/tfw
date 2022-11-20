package src

import (
	"strconv"
	"time"

	"gopkg.in/telebot.v3"
)

type PinOptions struct {
	Pin    bool
	Notify bool
}

func (t *TMessageF) ForwardMessages(vl []GeneralDatabaseSchema, pin PinOptions, sleep time.Duration, sendable any, opts ...any) (success int, failed int) {

	for _, v := range vl {
		time.Sleep(sleep)

		b := t.Bots.FindbyId(v.BotID)
		if b == nil {
			t.Log.Errorf("bot %s not registered!", v.BotID)
			continue
		}

		t.Log.Infof("bot %s sending message to chat %d", v.BotID, v.ChatID)

		msg, err := b.Send(&telebot.Chat{ID: v.ChatID}, sendable, opts...)
		if err != nil {
			failed++
			t.Log.Errorf("bot %s chat %d:%s error while sending message= %s", v.BotID, v.ChatID, v.ChatName, err)
			continue
		}

		success++

		if v.LastMessageID != 0 {
			t.Log.Infof("bot %s chat %d deleting the old message", v.BotID, v.ChatID)

			b.Delete(&telebot.Message{
				ID: int(v.LastMessageID),
				Chat: &telebot.Chat{
					ID: v.ChatID,
				},
			})
		}

		v.Update("last_message_id", strconv.Itoa(msg.ID))

		if pin.Pin {
			if !v.ChatRights.CanPinMessages {
				t.Log.Warnf("bot %s chat %d pin mode is on but admin required", v.BotID, v.ChatID)
			} else {
				var opts []any
				if !pin.Notify {
					opts = append(opts, telebot.Silent)
				}
				if err := b.Pin(msg, opts...); err != nil {
					t.Log.Errorf("bot %s chat %d pin err %s", v.BotID, v.ChatID, err)
				} else {
					t.Log.Infof("bot %s chat %d pinned", v.BotID, v.ChatID)
				}
			}
		}

		time.Sleep(sleep)

	}

	return
}
