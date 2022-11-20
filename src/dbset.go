package src

import (
	"fmt"
)

func (t *TMessageF) AddSetUpdate(v GeneralDatabaseSchema) error {
	old, err := t.Database.FindFirstMulti("chat_id", fmt.Sprint(v.ChatID), "bot_id", v.BotID)

	if err != nil {
		t.Log.Infof("%s-%s not exists on the database:%s adding", v.ChatID, v.BotID, err)
		return t.Database.New(v)
	}

	t.Log.Infof("%s-%s exists on the database:%s updating", v.BotID, v.ChatID)
	if err := old.Update("chat_rights", tbRightToString(v.ChatRights)); err != nil {
		return err
	}
	return old.Update("bot_role", v.BotRole)
}

func (t *TMessageF) DeleteFromDatabase(v GeneralDatabaseSchema) error {
	old, err := t.Database.FindFirstMulti("chat_id", fmt.Sprint(v.ChatID), "bot_id", v.BotID)

	if err != nil {
		return err
	}

	return old.Delete()
}
