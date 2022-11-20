package src

import (
	"fmt"
	"time"

	"github.com/jessevdk/go-flags"
	"gopkg.in/telebot.v3"
)

type ForwardCliOptions struct {
	//Client string `short:"c" long:"client" description:""`
	Slient bool `short:"s" long:"silent" description:"sets up message are silent or not"`
	//Max    int  `short:"m" long:"max" description:"max message number"`
	Delay int  `short:"d" long:"delay" required:"true" description:"delay time 1s=1000"`
	Pin   bool `short:"p" long:"pin"`
}

func (t *TMessageF) forward(group *telebot.Group) {
	group.Handle("/forward", func(ctx telebot.Context) error {
		if len(ctx.Args()) == 0 || ctx.Message() == nil || ctx.Message().ReplyTo == nil {
			var m string
			for _, v := range t.Bots {
				m += v.Id + "/"
			}
			return ctx.Reply(fmt.Sprintf("Usage: `<REPLY MESSAGE> /forward (all/%s) -h`", m), telebot.ModeMarkdown)
		}

		var fco ForwardCliOptions

		rest, err := flags.ParseArgs(&fco, ctx.Args())
		if err != nil {
			return ctx.Reply(err.Error())
		}

		if len(rest) == 0 {
			return ctx.Reply("please provide me something!")
		}

		ctx.Reply(fmt.Sprintf("Searching %s..", rest[0]))

		now := time.Now()

		var sandle any
		switch {
		case ctx.Message().ReplyTo.Photo != nil:

			ph := ctx.Message().ReplyTo.Photo
			/*if err := downloadFile(&ph.File, ctx.Bot()); err != nil {
				return ctx.Reply(err.Error())
			}*/

			//ph.File = getdownloadedFile()
			sandle = &telebot.Photo{
				File:    ph.File,
				Caption: ctx.Message().ReplyTo.Caption,
			}

		// case ctx.Message().ReplyTo.Video != nil:

		// 	sandle = &telebot.Video{
		// 		File:      ctx.Message().ReplyTo.Video.File,
		// 		Caption:   ctx.Message().ReplyTo.Caption,
		// 		Streaming: true,
		// 	}

		// case ctx.Message().ReplyTo.Document != nil:
		// 	sandle = ctx.Message().ReplyTo.Document
		// 	sandle.(*telebot.Document).Caption = ctx.Message().ReplyTo.Caption

		default:
			sandle = ctx.Message().ReplyTo.Text
			if sandle == "" {
				return ctx.Reply("Please give me only image or text!")
			}
		}

		//var keys = []any{telebot.ModeMarkdown}
		var keys = make([]any, 0)
		keys = append(keys, telebot.ModeMarkdown)
		if fco.Slient {
			keys = append(keys, telebot.Silent)
		}

		if ctx.Message().Entities != nil {
			keys = append(keys, ctx.Message().Entities)
		}

		var list []GeneralDatabaseSchema

		switch rest[0] {
		case "all":
			list, err = t.Database.GetAll()
		default:
			if _bot := t.Bots.FindbyId(rest[0]); _bot == nil {
				err = fmt.Errorf("bot id %s not exists", rest[0])
				break
			}

			list, err = t.Database.FindAll("bot_id", rest[0])
		}

		if err != nil {
			t.Log.Errorf("Message request occured %v", err)
			return ctx.Reply("Error: " + err.Error())
		}

		ctx.Reply(fmt.Sprintf("Processing %d value", len(list)))

		t, f := t.ForwardMessages(list, PinOptions{
			Pin:    fco.Pin,
			Notify: fco.Slient,
		}, time.Duration(fco.Delay)*time.Millisecond, sandle, keys...)

		future := time.Now()

		totalm := future.Sub(now).Minutes()

		return ctx.Reply(fmt.Sprintf("Done!\nTotal m %fm\nTotal success: %d\nTotal fail: %d\n\n You can watch logs.", totalm, t, f))
	})
}
