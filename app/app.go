package app

import (
	"log"

	"github.com/nyan2d/bolteo"
	"github.com/nyan2d/vacmanbot/app/commands"
	"github.com/nyan2d/vacmanbot/app/usermanager"
	tg "gopkg.in/tucnak/telebot.v2"
)

type App struct {
	bot         *tg.Bot
	database    *bolteo.Bolteo
	usermanager *usermanager.UserManager
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
	hook, err := bt.GetWebhook()
	if err != nil {
		log.Println("webhookerr", err)
	} else {
		log.Printf("%+v", hook)
	}

	a := &App{
		bot:         bt,
		database:    db,
		usermanager: um,
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
}

func (ap *App) bindHandlers() {
	for k, v := range ap.commands {
		ap.bot.Handle(k, v.Execute)
	}
}

func createMainPoller(endpoint, certificate, key string) tg.Poller {
	return &tg.Webhook{
		Listen: ":80",
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
