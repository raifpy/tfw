package src

import (
	"testing"

	"gopkg.in/telebot.v3"
)

func Test_Database(t *testing.T) {
	db, err := NewDatabase(DatabaseOptions{
		Path: "debug.db",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.New(GeneralDatabaseSchema{
		BotID:      "test value",
		ChatID:     111111,
		ChatRights: telebot.AdminRights(),
		BotRole:    "admin",
		ChatType:   "channel",
		ChatName:   "Zortchannel",
	}); err != nil {
		t.Fatal(err)
	}

	veri, err := db.FindAll("bot_id", "test value")
	if err != nil {
		t.Fatal(err)
	}
	_ = veri
}
