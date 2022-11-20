package src

import (
	"encoding/json"
	"time"

	"gopkg.in/telebot.v3"
)

func howLongToTomorrow() time.Duration {
	now := time.Now()

	hour, min, second := now.Clock()
	if 24-hour != 0 {
		hour = 24 - hour - 1
	}
	if 60-min != 0 {
		min = 60 - min - 1
	}
	if 60-second != 0 {
		second = 60 - second
	}
	return now.Add(time.Second*time.Duration(second) + time.Minute*time.Duration(min) + time.Hour*time.Duration(hour)).Sub(now)
}

func tbRightToString(r telebot.Rights) string {
	b, _ := json.Marshal(r)
	return string(b)
}

func stringToTbRights(value string) (r telebot.Rights, err error) {
	err = json.Unmarshal([]byte(value), &r)
	return
}

func downloadFile(file *telebot.File, bot *telebot.Bot) error {
	return bot.Download(file, "/tmp/dd")
}

func getdownloadedFilePath() string {
	return "/tmp/dd"
}

func getdownloadedFile() telebot.File {
	return telebot.FromDisk(getdownloadedFilePath())
}

func PublicTbRightsToString(r telebot.Rights) string {
	return tbRightToString(r)
}
