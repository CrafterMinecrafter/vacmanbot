package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/nyan2d/vacmanbot/app/pvpgame"
	"github.com/nyan2d/vacmanbot/app/usermanager"
	tg "gopkg.in/tucnak/telebot.v2"
)

type kubaprinter struct {
	lines []string
	index int
}

func newkubaprinter(items []string) *kubaprinter {
	return &kubaprinter{
		lines: items,
		index: 0,
	}
}

func (k *kubaprinter) next() bool {
	return k.index < len(k.lines)
}

func (k *kubaprinter) print() string {
	x := []string{}
	z := 0
	for k.next() && z < 4 {
		x = append(x, k.lines[k.index])
		k.index++
		z++
	}
	return strings.Join(x, "\n")
}

type FightCommand struct {
	bot  *tg.Bot
	game *pvpgame.Game
	um   *usermanager.UserManager
}

func (f FightCommand) Execute(m *tg.Message) {
	if m.ReplyTo == nil || m.ReplyTo.Sender == nil {
		f.bot.Reply(m, "Не с кем драться")
		return
	}

	if m.ReplyTo.Sender != nil && m.ReplyTo.Sender.ID == m.Sender.ID {
		f.bot.Reply(m, "Нельзя драться с самим собой")
		return
	}

	yourId := m.Sender.ID
	victimId := m.ReplyTo.Sender.ID
	isBot := m.ReplyTo.Sender.IsBot

	you := f.game.GetPlayer(yourId)
	victim := f.game.GetPlayer(victimId)
	if isBot {
		victim.IsBot = isBot
	}

	fight := f.game.MakeFight(you, victim, f.um.GetNames(yourId), f.um.GetNames(victimId))
	result := fight.Execute()

	// 900 IQ
	go func() {
		msg, _ := f.bot.Reply(m, result.StringPre())
		time.Sleep(3 * time.Second)
		fightlog := result.StringSliceFight()
		if len(fightlog) > 0 {
			printer := newkubaprinter(fightlog)
			for printer.next() {
				msg, _ = f.bot.Edit(msg, msg.Text, printer.print())
				time.Sleep(1500 * time.Millisecond)
			}
		}
		f.bot.Edit(msg, msg.Text, result.StringPost())
	}()
}

type BossCommand struct {
	bot  *tg.Bot
	game *pvpgame.Game
	um   *usermanager.UserManager
}

func (b *BossCommand) Execute(m *tg.Message) {
	names := []string{
		"mrekk", "WhiteCat", "aetrna", "Lifeline", "RyuK", "Akolibed", "NyanPotato", "Vaxei", "FlyingTuna", "A21", "Intercambing",
		"Andros", "Micca", "Mathi", "Dereban", "Paraqeet", "Karthy", "Aireu", "Rafis", "Bubbleman", "badeu", "Alumetri",
		"Freddie Benson", "InBefore", "im a fancy lad", "SWAGGYSWAGSTER", "ChocoPafe", "Spare", "xootynator", "Umbre",
	}

	yourId := m.Sender.ID
	you := b.game.GetPlayer(yourId)

	if you.Stats.Level < 3 {
		b.bot.Reply(m, "Приходи, когда будешь 3 уровня, малыш!")
		return
	}

	boss := b.game.CreateBossFor(you)
	fight := b.game.MakeFight(you, boss, b.um.GetNames(yourId), names[rand.Intn(len(names))])
	result := fight.Execute()

	go func() {
		msg, _ := b.bot.Reply(m, result.StringPre())
		time.Sleep(3 * time.Second)
		fightlog := result.StringSliceFight()
		if len(fightlog) > 0 {
			printer := newkubaprinter(fightlog)
			for printer.next() {
				msg, _ = b.bot.Edit(msg, msg.Text, printer.print())
				time.Sleep(1500 * time.Millisecond)
			}
		}
		b.bot.Edit(msg, msg.Text, result.StringPost())
	}()
}

type StatsCommand struct {
	bot  *tg.Bot
	game *pvpgame.Game
}

func (sc *StatsCommand) Execute(m *tg.Message) {
	yourId := m.Sender.ID
	sc.bot.Reply(m, sc.game.GetStatsStr(yourId))
}

type ItemsCommand struct {
	bot  *tg.Bot
	game *pvpgame.Game
}

func (ic *ItemsCommand) Execute(m *tg.Message) {
	yourId := m.Sender.ID
	ic.bot.Reply(m, ic.game.GetItemsStr(yourId))
}

type TopCommand struct {
	bot  *tg.Bot
	game *pvpgame.Game
	um   *usermanager.UserManager
}

func (tc *TopCommand) Execute(m *tg.Message) {
	top := tc.game.GetTopPlayersByExpStr(func(id int) string {
		return tc.um.GetNames(id)
	})
	tc.bot.Reply(m, top)
}

type TopEloCommand struct {
	bot  *tg.Bot
	game *pvpgame.Game
	um   *usermanager.UserManager
}

func (te *TopEloCommand) Execute(m *tg.Message) {
	top := te.game.GetTopPlayersByEloStr(func(id int) string {
		return te.um.GetNames(id)
	})
	te.bot.Reply(m, top)
}

type SwitchWeaponCommand struct {
	bot  *tg.Bot
	game *pvpgame.Game
}

func (swc *SwitchWeaponCommand) Execute(m *tg.Message) {
	yourId := m.Sender.ID
	p := swc.game.GetPlayer(yourId)
	aw := p.Items.WeaponID
	p.Items.WeaponID = p.Items.ArchivedWeapon
	p.Items.ArchivedWeapon = aw
	swc.game.SavePlayer(p)
	swc.bot.Reply(m, "Оружие заменено!")
}

type SwitchArmorCommand struct {
	bot  *tg.Bot
	game *pvpgame.Game
}

func (sac *SwitchArmorCommand) Execute(m *tg.Message) {
	yourId := m.Sender.ID
	p := sac.game.GetPlayer(yourId)
	aa := p.Items.ArmorID
	p.Items.ArmorID = p.Items.ArchivedArmor
	p.Items.ArchivedArmor = aa
	sac.game.SavePlayer(p)
	sac.bot.Reply(m, "Броня заменена!")
}

type UpCommand struct {
	bot  *tg.Bot
	game *pvpgame.Game
}

func (uc *UpCommand) Execute(m *tg.Message) {
	arg := strings.Split(m.Payload, " ")
	if len(arg) != 3 {
		uc.bot.Reply(m, "Что-то не так. Правильный формат: /up УРОН ЗАЩИТА ЗДОРОВЬЕ")
		return
	}
	dmg, err1 := strconv.Atoi(arg[0])
	prt, err2 := strconv.Atoi(arg[1])
	hel, err3 := strconv.Atoi(arg[2])
	if err1 != nil || err2 != nil || err3 != nil ||
		dmg+prt+hel == 0 || dmg < 0 || prt < 0 || hel < 0 {
		uc.bot.Reply(m, "Что-то не так. Правильный формат: /up УРОН ЗАЩИТА ЗДОРОВЬЕ")
		return
	}

	player := uc.game.GetPlayer(m.Sender.ID)
	if player.Stats.UnusedPoints < dmg+prt+hel {
		uc.bot.Reply(m, fmt.Sprintf("Недостаточно очков. Доступно очков: %v", player.Stats.UnusedPoints))
		return
	}

	player.Stats.Damage += dmg
	player.Stats.Protection += prt
	player.Stats.Health += hel * 2
	uc.game.SavePlayer(player)

	uc.bot.Reply(m, "Готово!")
}
