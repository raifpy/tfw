package src

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func (t *TMessageF) Admin() {
	group := t.Root.Group()
	t.RoortGroup = group

	group.Use(func(hf telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			id := ctx.Sender().ID
			for _, pr := range t.Options.IdWhlitelist {
				if pr == id {
					return hf(ctx)

				}
			}

			return nil
		}

	})

	group.Use(middleware.Whitelist(t.Options.IdWhlitelist...))
	group.Use(func(hf telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			t.Log.Infof("Admin command usage: %d:%s -> %s", ctx.Sender().ID, ctx.Sender().Username, ctx.Text())
			return hf(ctx)
		}
	})
	group.Handle("/start", func(ctx telebot.Context) error {
		m, err := t.Root.Bot.Reply(ctx.Message(), "I am here. Counting total group and channels")
		if err != nil {
			return err
		}

		_, err = t.Root.Bot.Edit(m, fmt.Sprintf("I am here. I have total %d group/channels.", t.Database.CountAll()))
		return err
	})

	group.Handle("/count <key value> ...", func(ctx telebot.Context) error {
		if len(ctx.Args()) < 2 {
			return ctx.Reply("Please give least 2 argumants")
		}
		dig, err := t.Database.CountCustom(ctx.Args())
		if err != nil {
			return ctx.Reply(fmt.Sprintf("Error= %s", err))
		}
		return ctx.Reply(fmt.Sprintf("Total %d record", dig))
	})

	group.Handle("/log", func(ctx telebot.Context) error {

		if len(ctx.Args()) == 0 {
			logdir, err := os.ReadDir(t.Options.Log.LogDirector)
			if err != nil {
				return ctx.Reply(err.Error())

			}
			var buf = &strings.Builder{}
			if len(logdir) > 30 {
				logdir = logdir[len(logdir)-30:]
			}
			for _, i := range logdir {
				buf.WriteString(fmt.Sprintf("log: `%s`\n", i.Name()))
			}

			return ctx.Reply(buf.String(), telebot.ModeMarkdown)

		}
		var sendpath string
		var filename string

		switch ctx.Args()[0] {
		case "error":
			{
				tf, err := os.Create(fmt.Sprintf("error_log_%s", time.Now().Format("2006-01-02")))
				if err != nil {
					t.Log.Errorf("couldn't create error_log_file %v", err)
					return ctx.Reply(err.Error())
				}
				defer tf.Close()

				read, err := os.OpenFile(t.Log.GetCurrentLogPath(), os.O_RDONLY, 0)
				if err != nil {
					t.Log.Errorf("couldn't read ordinary log file %v", err)
					return ctx.Reply(err.Error())
				}
				defer read.Close()

				bufr := bufio.NewScanner(read)
				for bufr.Scan() {
					if bufr.Bytes()[10] != 'e' {
						continue
					}
					fmt.Fprintf(tf, "%s\n", bufr.Bytes())
				}
				filename = fmt.Sprintf("Error log %s", time.Now().Format("2006-01-02"))
				tf.Close()
				read.Close()

				sendpath = tf.Name()
			}
		case "last":
			sendpath = t.Log.GetCurrentLogPath()
			filename = sendpath

		default:
			logdir, err := os.ReadDir(t.Options.Log.LogDirector)
			if err != nil {
				return ctx.Reply(err.Error())
			}

			for _, i := range logdir {
				if i.Name() == ctx.Args()[0] {
					filename = i.Name()
					sendpath = filepath.Join(t.Options.Log.LogDirector, filename)
					break
				}
			}
			if sendpath == "" {
				return ctx.Reply("There is no such log file.")
			}

		}

		return ctx.Reply(&telebot.Document{
			Caption:  "LOG: " + ctx.Args()[0],
			File:     telebot.FromDisk(sendpath),
			MIME:     "text/plain",
			FileName: fmt.Sprintf("Log file %s", filename),
		})
	})

	group.Handle("/db", func(ctx telebot.Context) error {
		return ctx.Reply(&telebot.Document{
			File:     telebot.FromDisk(t.Options.Database.Path),
			Caption:  "sqlite3 database",
			FileName: "database.sqlite3",
			MIME:     "application/sqlite3",
		})
	})

	group.Handle("/shutdown", func(ctx telebot.Context) error {
		t.Log.Warnf("!! EMERGENCY EXIT TRIGGERED !!")
		ctx.Reply("Done.")
		os.Exit(0)
		return nil
	})

	t.forward(group)
}
