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
	lines   []string
	display []string
	index   int
}

func newkubaprinter(items []string) *kubaprinter {
	return &kubaprinter{
		lines:   items,
		display: []string{},
		index:   -1,
	}
}

func (k *kubaprinter) next() bool {
	k.index++
	return k.index < len(k.lines)
}

func (k *kubaprinter) print() string {
	if len(k.display) == 5 {
		k.display = k.display[1:]
	}
	k.display = append(k.display, k.lines[k.index])
	return strings.Join(k.display, "\n")
}

type FightCommand struct {
	bot  *tg.Bot
	game *pvpgame.Game
	um   *usermanager.UserManager
}

func NewFightCommand(bot *tg.Bot, game *pvpgame.Game, um *usermanager.UserManager) *FightCommand {
	return &FightCommand{
		bot:  bot,
		game: game,
		um:   um,
	}
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

	if !f.game.CheckTime(yourId) {
		f.bot.Reply(m, fmt.Sprintf("Следующий бой через %v", fmtDuration(f.game.NextFight(yourId))))
		return
	}

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
		time.Sleep(5 * time.Second)
		fightlog := result.StringSliceFight()
		if len(fightlog) > 0 {
			printer := newkubaprinter(fightlog)
			for printer.next() {
				if msg == nil {
					break
				}
				msg, _ = f.bot.Edit(msg, printer.print())
				time.Sleep(2 * time.Second)
			}
			time.Sleep(2 * time.Second)
		}
		if msg != nil {
			f.bot.Edit(msg, result.StringPost())
		}
	}()
	f.game.UpdateTime(yourId)
}

type BossCommand struct {
	bot  *tg.Bot
	game *pvpgame.Game
	um   *usermanager.UserManager
}

func NewBossCommand(bot *tg.Bot, game *pvpgame.Game, um *usermanager.UserManager) *BossCommand {
	return &BossCommand{
		bot:  bot,
		game: game,
		um:   um,
	}
}

func (b *BossCommand) Execute(m *tg.Message) {
	names := []string{
		"mrekk", "WhiteCat", "aetrna", "Lifeline", "RyuK", "Akolibed", "NyanPotato", "Vaxei", "FlyingTuna", "A21", "Intercambing",
		"Andros", "Micca", "Mathi", "Dereban", "Paraqeet", "Karthy", "Aireu", "Rafis", "Bubbleman", "badeu", "Alumetri",
		"Freddie Benson", "InBefore", "im a fancy lad", "SWAGGYSWAGSTER", "ChocoPafe", "Spare", "xootynator", "Umbre",
		"flae",
	}

	yourId := m.Sender.ID
	if !b.game.CheckTime(yourId) {
		b.bot.Reply(m, fmt.Sprintf("Следующий бой через %v", fmtDuration(b.game.NextFight(yourId))))
		return
	}
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
		time.Sleep(5 * time.Second)
		fightlog := result.StringSliceFight()
		if len(fightlog) > 0 {
			printer := newkubaprinter(fightlog)
			for printer.next() {
				if msg == nil {
					break
				}
				msg, _ = b.bot.Edit(msg, printer.print())
				time.Sleep(2 * time.Second)
			}
			time.Sleep(2 * time.Second)
		}
		if msg != nil {
			b.bot.Edit(msg, result.StringPost())
		}
	}()
	b.game.UpdateTime(yourId)
}

type StatsCommand struct {
	bot  *tg.Bot
	game *pvpgame.Game
}

func NewStatsCommand(bot *tg.Bot, game *pvpgame.Game) *StatsCommand {
	return &StatsCommand{
		bot:  bot,
		game: game,
	}
}

func (sc *StatsCommand) Execute(m *tg.Message) {
	yourId := m.Sender.ID
	sc.bot.Reply(m, sc.game.GetStatsStr(yourId))
}

type ItemsCommand struct {
	bot  *tg.Bot
	game *pvpgame.Game
}

func NewItemsCommand(bot *tg.Bot, game *pvpgame.Game) *ItemsCommand {
	return &ItemsCommand{
		bot:  bot,
		game: game,
	}
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

func NewTopCommand(bot *tg.Bot, game *pvpgame.Game, um *usermanager.UserManager) *TopCommand {
	return &TopCommand{
		bot:  bot,
		game: game,
		um:   um,
	}
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

func NewTopEloCommand(bot *tg.Bot, game *pvpgame.Game, um *usermanager.UserManager) *TopEloCommand {
	return &TopEloCommand{
		bot:  bot,
		game: game,
		um:   um,
	}
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

func NewSwitchWeaponCommand(bot *tg.Bot, game *pvpgame.Game) *SwitchWeaponCommand {
	return &SwitchWeaponCommand{
		bot:  bot,
		game: game,
	}
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

func NewSwitchArmorCommand(bot *tg.Bot, game *pvpgame.Game) *SwitchArmorCommand {
	return &SwitchArmorCommand{
		bot:  bot,
		game: game,
	}
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

func NewUpCommand(bot *tg.Bot, game *pvpgame.Game) *UpCommand {
	return &UpCommand{
		bot:  bot,
		game: game,
	}
}

func (uc *UpCommand) Execute(m *tg.Message) {
	arg := strings.Split(m.Payload, " ")
	if len(arg) != 3 {
		uc.bot.Reply(m, "Что-то не так. Правильный формат: /up УРОН ЗАЩИТА ЗДОРОВЬЕ\nПример: /up 1 2 3")
		return
	}
	dmg, err1 := strconv.Atoi(arg[0])
	prt, err2 := strconv.Atoi(arg[1])
	hel, err3 := strconv.Atoi(arg[2])
	if err1 != nil || err2 != nil || err3 != nil ||
		dmg+prt+hel == 0 || dmg < 0 || prt < 0 || hel < 0 {
		uc.bot.Reply(m, "Что-то не так. Правильный формат: /up УРОН ЗАЩИТА ЗДОРОВЬЕ\nПример: /up 1 2 3")
		return
	}

	player := uc.game.GetPlayer(m.Sender.ID)
	if player.Stats.UnusedPoints < dmg+prt+hel {
		uc.bot.Reply(m, fmt.Sprintf("Недостаточно очков. Доступно очков: %v", player.Stats.UnusedPoints))
		return
	}

	player.Stats.Damage += dmg
	player.Stats.Protection += prt
	player.Stats.Health += hel
	player.Stats.UnusedPoints -= dmg + prt + hel
	uc.game.SavePlayer(player)

	uc.bot.Reply(m, fmt.Sprintf("Готово! Доступно очков: %v", player.Stats.UnusedPoints))
}
