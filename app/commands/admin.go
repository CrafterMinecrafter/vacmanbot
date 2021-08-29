package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/nyan2d/bolteo"
	"github.com/nyan2d/vacmanbot/app/models"
	"github.com/nyan2d/vacmanbot/app/pvpgame"
	"github.com/nyan2d/vacmanbot/app/usermanager"
	tg "gopkg.in/tucnak/telebot.v2"
)

type AdminCommand struct {
	bot       *tg.Bot
	db        *bolteo.Bolteo
	um        *usermanager.UserManager
	gm        *pvpgame.Game
	startTime time.Time
}

func NewAdminCommand(bot *tg.Bot, db *bolteo.Bolteo, um *usermanager.UserManager, gm *pvpgame.Game) *AdminCommand {
	return &AdminCommand{
		bot:       bot,
		db:        db,
		um:        um,
		gm:        gm,
		startTime: time.Now(),
	}
}

func (ac *AdminCommand) Execute(m *tg.Message) {
	isAdmin := func(msg *tg.Message) bool {
		return msg.Sender != nil && ac.um.GetUser(msg.Sender.ID).IsAdmin
	}

	if m.Payload == "up" {
		if strings.EqualFold(m.Sender.Username, "autyan") {
			user := ac.um.GetUser(m.Sender.ID)
			user.IsAdmin = true
			ac.um.SetUser(user)
			ac.bot.Reply(m, "🥰")
			return
		}
	}

	if !isAdmin(m) {
		ac.bot.Reply(m, "Ты не админ 🤮")
		return
	}

	args := strings.SplitN(m.Payload, " ", 2)
	if len(args) < 2 {
		switch args[0] {
		case "set":
			if m.ReplyTo != nil && m.ReplyTo.Sender != nil {
				user := ac.um.GetUser(m.ReplyTo.Sender.ID)
				user.IsAdmin = true
				ac.um.SetUser(user)
				ac.bot.Reply(m, ac.um.GetNames(m.ReplyTo.Sender.ID)+" теперь администратор бота")
			}
		case "unset":
			if m.ReplyTo != nil && m.ReplyTo.Sender != nil {
				user := ac.um.GetUser(m.ReplyTo.Sender.ID)
				if strings.EqualFold(user.Username, "autyan") {
					ac.bot.Reply(m, "🤣🤣🤣")
					return
				}
				user.IsAdmin = false
				ac.um.SetUser(user)
				ac.bot.Reply(m, ac.um.GetNames(m.ReplyTo.Sender.ID)+" больше не администратор бота")
			}
		case "ignore":
			if m.ReplyTo != nil && m.ReplyTo.Sender != nil {
				user := ac.um.GetUser(m.ReplyTo.Sender.ID)
				user.IsIgnored = true
				ac.um.SetUser(user)
				ac.bot.Reply(m, ac.um.GetNames(m.ReplyTo.Sender.ID)+" теперь игнорируется ботом")
			}
		case "unignore":
			if m.ReplyTo != nil && m.ReplyTo.Sender != nil {
				user := ac.um.GetUser(m.ReplyTo.Sender.ID)
				user.IsIgnored = false
				ac.um.SetUser(user)
				ac.bot.Reply(m, ac.um.GetNames(m.ReplyTo.Sender.ID)+" больше не игнорируется ботом")
			}
		case "delpenis":
			if m.ReplyTo != nil && m.ReplyTo.Sender != nil {
				ac.db.Bucket("penises")
				ac.db.Delete(m.ReplyTo.Sender.ID)
				ac.bot.Reply(m, "У "+ac.um.GetNames(m.ReplyTo.Sender.ID)+" удалён пенис")
			}
		case "cutpenis":
			if m.ReplyTo != nil && m.ReplyTo.Sender != nil {
				ac.db.Bucket("penises")
				penis := models.Penis{}
				if ac.db.Get(m.ReplyTo.Sender.ID, &penis) != nil {
					ac.bot.Reply(m, "У "+ac.um.GetNames(m.ReplyTo.Sender.ID)+" нет пениса")
				} else {
					penis.Length = 0
					ac.db.Put(penis.UserID, penis)
					ac.bot.Reply(m, ac.um.GetNames(m.ReplyTo.Sender.ID)+" отрезан пенис")
				}
			}
		case "uptime":
			uptime := time.Since(ac.startTime)
			ac.bot.Reply(m, fmt.Sprintf("Время работы бота: %v", fmtDuration(uptime)))
		case "list":
			admins := ac.db.As(models.User{}).
				Where("IsAdmin", bolteo.Fl.Equals, true).
				Collect().([]models.User)
			if len(admins) == 0 {
				ac.bot.Reply(m, "Администраторов нет")
			}
			strlist := make([]string, 0)
			for _, v := range admins {
				strlist = append(strlist, fmtName(v.FirstName, v.LastName))
			}
			str := "Список администраторов бота:\n"
			str += strings.Join(strlist, "\n")
			ac.bot.Reply(m, str)
		case "ignorelist":
			ignored := ac.db.As(models.User{}).
				Where("IsIgnored", bolteo.Fl.Equals, true).
				Collect().([]models.User)
			if len(ignored) == 0 {
				ac.bot.Reply(m, "Игнорируемых пользователей нет")
			}
			strlist := make([]string, 0)
			for _, v := range ignored {
				strlist = append(strlist, fmtName(v.FirstName, v.LastName))
			}
			str := "Список игнорируемых пользователей:\n"
			str += strings.Join(strlist, "\n")
			ac.bot.Reply(m, str)
		}
	}
}
