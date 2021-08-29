package pvpgame

import (
	"fmt"
	"strings"
	"time"

	"github.com/nyan2d/bolteo"
)

type Game struct {
	db       *bolteo.Bolteo
	players  map[int]*Player
	weapons  map[int]*Weapon
	armors   map[int]*Armor
	lastturn map[int]time.Time
}

func NewGame(db *bolteo.Bolteo) *Game {
	g := &Game{
		db:       db,
		players:  make(map[int]*Player),
		weapons:  make(map[int]*Weapon),
		armors:   make(map[int]*Armor),
		lastturn: make(map[int]time.Time),
	}
	initBuckets(g.db)
	g.loadItems()
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

func (g *Game) GetPlayer(id int) *Player {
	g.db.Bucket("pvp_users")
	player, ok := g.players[id]
	if !ok {
		if g.db.Get(id, &player) != nil {
			player = CreatePlayer(id, false)
			g.db.Put(id, player)
		}
		g.players[id] = player
	}
	return player
}

func (g *Game) SavePlayer(p *Player) {
	if p.IsBot {
		return
	}

	g.players[p.ID] = p

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
	weap, armr := g.weapons[p.Items.WeaponID], g.armors[p.Items.ArmorID]
	qweap, qarmr := g.weapons[p.Items.ArchivedWeapon], g.armors[p.Items.ArchivedArmor]
	lines := []string{
		fmt.Sprintf("Оружие: %v; Урон: %v; Шанс крита: %v%%; Урон крита: %v%%", weap.Name, weap.Damage, int(weap.CritChance*100.0), int(weap.CritMultiplier*100.0)),
		fmt.Sprintf("Броня: %v; Защита: %v; Здоровье: %v", armr.Name, armr.Protection, armr.BonusHealth),
		fmt.Sprintf("Доступное оружие: %v; Урон: %v; Шанс крита: %v%%; Урон крита: %v%%", qweap.Name, qweap.Damage, int(qweap.CritChance*100.0), int(qweap.CritMultiplier*100.0)),
		fmt.Sprintf("Доступная броня: %v; Защита: %v; Здоровье: %v", qarmr.Name, qarmr.Protection, qarmr.BonusHealth),
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

func (g *Game) loadItems() {
	// add default items
	g.weapons[-1] = &Weapon{
		ID:   -1,
		Name: "Отсутствует",
	}
	g.armors[-1] = &Armor{
		ID:   -1,
		Name: "Отсутствует",
	}

	// load from databasew
	g.db.Bucket("pvp_weapons")
	weaponsfromdb := g.db.As(Weapon{}).Collect().([]Weapon)
	for _, v := range weaponsfromdb {
		g.weapons[v.ID] = &v
	}
	g.db.Bucket("pvp_armors")
	armorsfromdb := g.db.As(Armor{}).Collect().([]Armor)
	for _, v := range armorsfromdb {
		g.armors[v.ID] = &v
	}
}

func initBuckets(db *bolteo.Bolteo) {
	db.Bucket("pvp_users")
	db.InitBucket()
	db.Bucket("pvp_weapons")
	db.InitBucket()
	db.Bucket("pvp_armors")
	db.InitBucket()
}
