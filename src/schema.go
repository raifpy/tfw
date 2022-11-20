package src

import (
	"gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

type GeneralDatabaseSchema struct {
	d *Database `gorm:"-"`
	gorm.Model
	BotID         string
	ChatID        int64
	ChatName      string
	ChatType      string // channel / group / private
	BotRole       string // admin/user
	ErrorCounter  int
	LastMessageID int64

	AddedById       int64
	AddedByName     string
	AddedByUsername string
	ReplyIds        []int64        `gorm:"serializer:json"`
	ChatRights      telebot.Rights `gorm:"serializer:json"`
}

func (g GeneralDatabaseSchema) Update(column string, value string) error {
	return g.d.DB.Model(&g).Update(column, value).Error
}

func (g GeneralDatabaseSchema) Delete() error {
	return g.d.DB.Model(&g).Delete(&g).Error
}
