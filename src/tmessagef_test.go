package src

import (
	"fmt"
	"testing"
)

func TestTMessageF(t *testing.T) {
	obj, err := NewTMessageF(&Options{
		IdWhlitelist: []int64{821818337},
		TelegramOptions: []TelegramOptions{
			{
				Id:       "root",
				BotToken: "",
			},
		},
		Log: TMessageFLoggerOption{
			LogDirector: "debug_log",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("obj: %v\n", obj)
	t.Log("!DONE!")
	<-make(chan struct{})
}
