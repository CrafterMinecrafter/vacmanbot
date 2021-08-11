package commands

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/nyan2d/vacmanbot/app/models"
	"github.com/nyan2d/vacmanbot/app/usermanager"
	tg "gopkg.in/tucnak/telebot.v2"
)

var penisNames = []string{
	"пениса",
	"дружка",
	"малыша",
	"конца",
	"болта",
	"елдака",
	"дика",
	"мпх",
	"хуя",
	"хера",
	"хрена",
	"писюна",
}

type PenisCommand struct {
	bot     *tg.Bot
	penises map[int]models.Penis
	um      *usermanager.UserManager
}

func NewPenisCommand(bot *tg.Bot, u *usermanager.UserManager) *PenisCommand {
	return &PenisCommand{
		bot:     bot,
		penises: make(map[int]models.Penis),
		um:      u,
	}
}

func (pe *PenisCommand) Execute(m *tg.Message) {
	if m.Payload == "top" {
		items := make([]models.Penis, 0)
		for _, v := range pe.penises {
			items = append(items, v)
		}
		if len(items) < 1 {
			pe.bot.Reply(m, "Нет пенисов 😭")
			return
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].Length > items[j].Length
		})
		result := "Самая большая елда у " + pe.um.GetNames(items[0].UserID) + ": " + lengthToString(items[0].UserID)
		if len(items) > 1 {
			micro := items[len(items)-1]
			result += "\n Владелец микрописьки " + pe.um.GetNames(micro.UserID) + ": " + lengthToString(micro.Length)
		}
		pe.bot.Reply(m, result)
	} else {
		penis, ok := pe.penises[m.Sender.ID]
		if !ok {
			penis = models.Penis{
				UserID:  m.Sender.ID,
				Length:  penisRoll(),
				Expires: time.Now().Add(time.Hour),
			}
			pe.penises[penis.UserID] = penis
		}

		if time.Now().After(penis.Expires) {
			penis.Length = penisRoll()
			penis.Expires = time.Now().Add(time.Hour)
			pe.penises[penis.UserID] = penis
		}

		text := fmt.Sprintf(
			"Длина твоего %v: %v\nСледующий запрос через %v",
			penisNames[rand.Intn(len(penisNames))],
			lengthToString(penis.Length),
			fmtDuration(time.Until(penis.Expires)),
		)
		pe.bot.Reply(m, text)
	}
}
