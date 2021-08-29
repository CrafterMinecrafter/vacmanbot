package app

import (
	"log"

	"github.com/nyan2d/bolteo"
	"github.com/nyan2d/vacmanbot/app/commands"
	"github.com/nyan2d/vacmanbot/app/pvpgame"
	"github.com/nyan2d/vacmanbot/app/usermanager"
	tg "gopkg.in/tucnak/telebot.v2"
)

type App struct {
	bot         *tg.Bot
	database    *bolteo.Bolteo
	usermanager *usermanager.UserManager
	game        *pvpgame.Game
	commands    map[string]commands.Command
}

func NewApp(token, endpoint, certificate, key, databasePath string) *App {
	db := bolteo.MustOpen(databasePath)
	um := usermanager.NewUserManager(db)
	bt, err := tg.NewBot(tg.Settings{
		Token:  token,
		Poller: createMiddlePoller(createMainPoller(endpoint, certificate, key), um),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("success:", bt.Me.FirstName, bt.Me.LastName)

	a := &App{
		bot:         bt,
		database:    db,
		usermanager: um,
		game:        pvpgame.NewGame(db),
	}

	a.fillCommands()
	a.bindHandlers()

	return a
}

func (ap *App) Start() {
	ap.bot.Start()
}

func (ap *App) fillCommands() {
	ap.commands = make(map[string]commands.Command)
	ap.commands["/penis"] = commands.NewPenisCommand(ap.bot, ap.usermanager)
	ap.commands["/admin"] = commands.NewAdminCommand(ap.bot, ap.database, ap.usermanager)

	ap.commands["/fight"] = commands.NewFightCommand(ap.bot, ap.game, ap.usermanager)
	ap.commands["/boss"] = commands.NewBossCommand(ap.bot, ap.game, ap.usermanager)
	ap.commands["/stats"] = commands.NewStatsCommand(ap.bot, ap.game)
	ap.commands["/items"] = commands.NewItemsCommand(ap.bot, ap.game)
	ap.commands["/top"] = commands.NewTopCommand(ap.bot, ap.game, ap.usermanager)
	ap.commands["/elo"] = commands.NewTopEloCommand(ap.bot, ap.game, ap.usermanager)
	ap.commands["/switchweapon"] = commands.NewSwitchWeaponCommand(ap.bot, ap.game)
	ap.commands["/switcharmor"] = commands.NewSwitchArmorCommand(ap.bot, ap.game)
	ap.commands["/up"] = commands.NewUpCommand(ap.bot, ap.game)
}

func (ap *App) bindHandlers() {
	for k, v := range ap.commands {
		ap.bot.Handle(k, v.Execute)
	}
}

func createMainPoller(endpoint, certificate, key string) tg.Poller {
	return &tg.Webhook{
		Listen: ":443",
		TLS: &tg.WebhookTLS{
			Cert: certificate,
			Key:  key,
		},
		Endpoint: &tg.WebhookEndpoint{
			PublicURL: endpoint,
			Cert:      certificate,
		},
	}
}

func createMiddlePoller(p tg.Poller, um *usermanager.UserManager) tg.Poller {
	filter := func(u *tg.Update) bool {
		if u.Message != nil && u.Message.Sender != nil {
			um.Check(*u.Message.Sender)
			usr := um.GetUser(u.Message.Sender.ID)
			return !usr.IsIgnored || usr.IsAdmin
		}
		return true
	}
	return tg.NewMiddlewarePoller(p, filter)
}
