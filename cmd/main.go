package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"tmessagefb/src"

	"gopkg.in/telebot.v3"
	"gopkg.in/yaml.v3"
)

var config *string = flag.String("config", "config.yaml", "Config file")
var recovergroup = flag.String("recovergroup", "", "Recover group from log file")

func init() {
	flag.Parse()
}

func writedefaultconfig() {

	f, err := os.Create(*config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while creating default config file %s %d\n", *config, err)
		os.Exit(1)
	}
	yaml.NewEncoder(f).Encode(src.Options{
		Log: src.TMessageFLoggerOption{},
		TelegramOptions: []src.TelegramOptions{
			{},
		},
		Database: src.DatabaseOptions{},
	})
	f.Close()
}

func writenewconfig(s *src.Options) error {
	f, err := os.Create(*config)
	if err != nil {
		return err
	}
	defer f.Close()

	return yaml.NewEncoder(f).Encode(*s)
}

func main() {
	if *config == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(*config)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error while opening file %s %v\n", *config, err)
		fmt.Println("creating default config file")
		writedefaultconfig()
		os.Exit(1)
	}

	var s *src.Options = &src.Options{}
	if err := yaml.NewDecoder(file).Decode(s); err != nil {
		fmt.Fprintf(os.Stderr, "error while decoding yaml config %v\n", err)
		os.Exit(1)
	}
	tm, err := src.NewTMessageF(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unexcepted error %v", err)
		os.Exit(1)
	}

	if *recovergroup != "" {
		file, err := os.Open(*recovergroup)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while reading recovergroup file %s %v\n", *recovergroup, err)
			os.Exit(1)
		}
		defer file.Close()
		buf := bufio.NewReader(file)
		for {
			s, err := buf.ReadString('\n')
			if err != nil {
				break
			}

			if s[10] != 'e' || !strings.Contains(s, "not exists in the list") {
				continue
			}
			botindex := strings.Index(s, "bot ")
			if botindex == -1 {
				continue
			}
			spaceindex := strings.Index(s[botindex+5:], " ")
			if spaceindex == -1 {
				continue
			}

			bot := s[botindex+5 : botindex+5+spaceindex]
			chatindex := strings.Index(s, "chat (") + 6
			if chatindex == -1 {
				continue
			}
			if spaceindex = strings.Index(s[chatindex:], ":"); spaceindex == -1 {
				continue
			}
			chat := s[chatindex : chatindex+spaceindex]
			fmt.Printf("chat: %v bot: %v\n", chat, bot)
			chatid, err := strconv.Atoi(chat)
			if err != nil {
				fmt.Fprintf(os.Stderr, "chat %v: %v\n", chat, err)
				continue
			}
			botm := tm.Bots.FindbyId(bot)
			if botm == nil {
				fmt.Fprintf(os.Stderr, "bot %s not found\n", bot)
				continue
			}
			time.Sleep(time.Second / 2)
			chatm, err := botm.ChatByID(int64(chatid))
			if err != nil {
				fmt.Fprintf(os.Stderr, "chat:%v bot:%v err:%v\n", chat, bot, err)
				continue
			}
			memb, err := botm.ChatMemberOf(telebot.ChatID(chatid), &telebot.User{
				ID: botm.Me.ID,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "chat:%v bot:%v err:%v\n", chat, bot, err)
				continue
			}
			role := "user"
			if memb.CanPinMessages {
				role = "administrator"
			}
			if err := tm.AddSetUpdate(src.GeneralDatabaseSchema{
				BotID:      bot,
				ChatID:     int64(chatid),
				ChatName:   chatm.Title,
				ChatType:   string(chatm.Type),
				ChatRights: memb.Rights,
				BotRole:    role,
			}); err != nil {
				fmt.Fprintf(os.Stderr, "error while adding set update %v\n", err)
				continue
			}

			fmt.Printf("bot %v chat %v successfully added \n", bot, chat)

		}

		os.Exit(0)
	}

	log.Println("init")

	tm.RoortGroup.Handle("/admin", func(ctx telebot.Context) error {
		if len(ctx.Args()) < 2 {
			return ctx.Reply(fmt.Sprintf("Admin list: `%+v`\n\nUsage: `/admin add/delete <id>`", s.IdWhlitelist), telebot.ModeMarkdown)
		}

		id, err := strconv.Atoi(ctx.Args()[1])
		if err != nil {
			return ctx.Reply(err.Error())
		}

		tm.Log.Warnf("NEW ADMIN OPERATION %d-%s by %d %s", id, ctx.Args()[0], ctx.Sender().ID, ctx.Sender().Username)

		switch ctx.Args()[0] {
		case "add":
			s.IdWhlitelist = append(s.IdWhlitelist, int64(id))
		case "delete":
			var newlist = []int64{}
			for _, v := range s.IdWhlitelist {
				if v != int64(id) {
					newlist = append(newlist, v)
				}
			}
			s.IdWhlitelist = newlist
		default:
			return ctx.Reply("unknown operation. ")
		}
		tm.Log.Warnf("NEW ADMIN OPERATION %d-%s by %d WRITED", id, ctx.Args()[0], ctx.Sender().ID)

		if err := writenewconfig(s); err != nil {
			ctx.Reply(fmt.Sprintf("New admin seted but not couldn't writed into file: %s", err))
		}
		return ctx.Reply(fmt.Sprintf("New Admin list: `%+v`", s.IdWhlitelist), telebot.ModeMarkdown)
	})

	tm.CloseWait()
}
