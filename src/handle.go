package src

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

func (t *TMessageF) HandleTelegramActions() {

	for i, _ := range t.Bots {
		bot := t.Bots[i]

		// bot.Handle("/start", func(ctx telebot.Context) error {
		// 	return ctx.Reply(fmt.Sprintf("Selam: %s %+v", bot.Id, t.Bots))
		// })

		bot.Handle(telebot.OnMyChatMember, func(ctx telebot.Context) error {

			if ctx.ChatMember() == nil || ctx.ChatMember().NewChatMember == nil {
				return fmt.Errorf("unexcepted value:handle.go -> %v", ctx)
			}

			chat := ctx.ChatMember().Chat
			role := ctx.ChatMember().NewChatMember.Role

			t.Log.Infof("bot %s chat (%d:%s) action handled %s", bot.Id, chat.ID, chat.Title, role)

			dbsc := GeneralDatabaseSchema{
				BotID:           bot.Id,
				ChatID:          chat.ID,
				ChatName:        chat.Title,
				BotRole:         "member",
				ChatType:        string(chat.Type),
				ErrorCounter:    0,
				LastMessageID:   0,
				AddedById:       ctx.Sender().ID,
				AddedByName:     ctx.Sender().FirstName,
				AddedByUsername: ctx.Sender().Username,
				ChatRights:      telebot.NoRights(),
			}
			var upd bool
			// fmt.Printf("exists: %v\n", exists)
			// fmt.Printf("existserror: %v\n", existserror)

			switch role {
			case telebot.Kicked:
				{
					t.Log.Infof("bot %s kicked from chat (%d:%s). id removing from list", bot.Id, chat.ID, chat.Title)

					if err := t.DeleteFromDatabase(dbsc); err != nil {
						t.Log.Errorf("bot %s chat (%d:%s) couldn't remove from the list= %v", bot.Id, chat.ID, chat.Title, err)
						break
					}

				}
			case telebot.Member:
				{
					t.Log.Infof("bot %s added to (%d:%s) chat. Adding to the list", bot.Id, chat.ID, chat.Title)

					upd = true
					break

				}

			case telebot.Administrator:
				{

					t.Log.Infof("bot %s is assigned to chat (%d:%s). updating role", bot.Id, chat.ID, chat.Title)

					dbsc.BotRole = "administrator"
					dbsc.ChatRights = ctx.ChatMember().NewChatMember.Rights
					upd = true
					break

				}
			}

			if upd {
				if err := t.AddSetUpdate(dbsc); err != nil {
					t.Log.Errorf("bot %s chat (%d:%s) couldn't setupdate to the list= %v", bot.Id, chat.ID, chat.Title, err)
				}
			}

			return nil
		})

	}
}
