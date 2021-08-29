package pvpgame

import (
	"fmt"
	"strings"
	"time"

	"github.com/nyan2d/bolteo"
)

type Game struct {
	db       *bolteo.Bolteo
	lastturn map[int]time.Time
}

func NewGame(db *bolteo.Bolteo) *Game {
	g := &Game{
		db:       db,
		lastturn: make(map[int]time.Time),
	}
	initBuckets(g.db)
	return g
}

func (g *Game) CheckTime(id int) bool {
	t, ok := g.lastturn[id]
	if !ok {
		g.UpdateTime(id)
		return true
	}
	return t.Add(time.Minute).Before(time.Now())
}

func (g *Game) UpdateTime(id int) {
	g.lastturn[id] = time.Now()
}

func (g *Game) NextFight(id int) time.Duration {
	return time.Until(g.lastturn[id].Add(time.Minute))
}

func (g *Game) GetPlayer(id int) *Player {
	g.db.Bucket("pvp_users")
	player := &Player{}
	if g.db.Get(id, player) != nil {
		player = CreatePlayer(id, false)
	}
	return player
}

func (g *Game) SavePlayer(p *Player) {
	if p.IsBot {
		return
	}
	g.db.Bucket("pvp_users")
	g.db.Put(p.ID, p)
}

func (g *Game) MakeFight(firstPlayer, secondPlayer *Player, firstPlayerName, secondPlayerName string) *Fight {
	return &Fight{
		game: g,
		P: []*Player{
			firstPlayer,
			secondPlayer,
		},
		N: []string{
			firstPlayerName,
			secondPlayerName,
		},
		log: NewBattleLog(),
	}
}

func (g *Game) GetStatsStr(id int) string {
	p := g.GetPlayer(id)
	lines := []string{
		fmt.Sprintf("Уровень: %v", p.Stats.Level),
		fmt.Sprintf("Опыт: %v/%v", p.Stats.Experience, calcLevelToExp(p.Stats.Level+1)),
		fmt.Sprintf("Урон: %v", p.Stats.Damage),
		fmt.Sprintf("Защита: %v", p.Stats.Protection),
		fmt.Sprintf("Здоровье: %v", p.Stats.Health),
		fmt.Sprintf("Очки умений: %v", p.Stats.UnusedPoints),
		fmt.Sprintf("Доллары: %v", p.Stats.Gold),
		fmt.Sprintf("Боёв: %v [%v/%v]", p.Statistics.Fights, p.Statistics.Wins, p.Statistics.Loses),
		fmt.Sprintf("Эло: %v", p.Statistics.Elo),
	}
	return strings.Join(lines, "\n")
}

func (g *Game) GetItemsStr(id int) string {
	p := g.GetPlayer(id)

	weapon := &Weapon{
		ID:   -1,
		Name: "Отсутствует",
	}
	archivedWeapon := &Weapon{
		ID:   -1,
		Name: "Отсутствует",
	}
	armor := &Armor{
		ID:   -1,
		Name: "Отсутствует",
	}
	archivedArmor := &Armor{
		ID:   -1,
		Name: "Отсутствует",
	}
	g.db.Bucket("pvp_weapons")
	g.db.Get(p.Items.WeaponID, weapon)
	g.db.Get(p.Items.ArchivedWeapon, archivedWeapon)
	g.db.Bucket("pvp_armors")
	g.db.Get(p.Items.ArmorID, armor)
	g.db.Get(p.Items.ArchivedArmor, archivedArmor)

	lines := []string{
		fmt.Sprintf("Оружие: %v\n•‎ Урон: %v\n•‎ Шанс крита: %v%%\n•‎ Урон крита: %v%%", weapon.Name, weapon.Damage, int(weapon.CritChance*100.0), int(weapon.CritMultiplier*100.0)),
		fmt.Sprintf("Броня: %v\n•‎ Защита: %v\n•‎ Здоровье: %v\n", armor.Name, armor.Protection, armor.BonusHealth),
		"",
		fmt.Sprintf("Доступное оружие: %v\n•‎ Урон: %v\n•‎ Шанс крита: %v%%\n•‎ Урон крита: %v%%", archivedWeapon.Name, archivedWeapon.Damage, int(archivedWeapon.CritChance*100.0), int(archivedWeapon.CritMultiplier*100.0)),
		fmt.Sprintf("Доступная броня: %v\n•‎ Защита: %v\n•‎ Здоровье: %v", archivedArmor.Name, archivedArmor.Protection, archivedArmor.BonusHealth),
	}
	return strings.Join(lines, "\n")
}

func (g *Game) GetTopPlayersByExpStr(getname func(id int) string) string {
	g.db.Bucket("pvp_users")
	top := g.db.As(Player{}).OrderBy("Stats.Experience").Reverse().Take(5).Collect().([]Player)
	lines := []string{}
	i := 1
	for _, v := range top {
		lines = append(lines, fmt.Sprintf("%v) %v - %v ур", i, getname(v.ID), v.Stats.Level))
		i++
	}
	return strings.Join(lines, "\n")
}

func (g *Game) GetTopPlayersByEloStr(getname func(id int) string) string {
	g.db.Bucket("pvp_users")
	top := g.db.As(Player{}).OrderBy("Statistics.Elo").Reverse().Take(5).Collect().([]Player)
	lines := []string{}
	i := 1
	for _, v := range top {
		if v.ID == 1862102456 {
			continue
		}
		lines = append(lines, fmt.Sprintf("%v) %v - %v", i, getname(v.ID), v.Statistics.Elo))
		i++
	}
	return strings.Join(lines, "\n")
}

func initBuckets(db *bolteo.Bolteo) {
	db.Bucket("pvp_users")
	db.InitBucket()
	db.Bucket("pvp_weapons")
	db.InitBucket()
	db.Bucket("pvp_armors")
	db.InitBucket()
}
