package pvpgame

import "github.com/nyan2d/bolteo"

type Game struct {
	db      *bolteo.Bolteo
	players map[int]*Player
	weapons map[int]*Weapon
	armors  map[int]*Armor
}

func NewGame(db *bolteo.Bolteo) *Game {
	g := &Game{
		db:      db,
		players: make(map[int]*Player),
		weapons: make(map[int]*Weapon),
		armors:  make(map[int]*Armor),
	}
	initBuckets(g.db)
	g.loadItems()
	return g
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

	// load from database
	g.db.Bucket("pvp_weapons")
	weaponsfromdb := g.db.As(Weapon{}).Collect().([]Weapon)
	g.db.Bucket("pvp_armors")
	armorsfromdb := g.db.As(Armor{}).Collect().([]Armor)
	for _, v := range weaponsfromdb {
		g.weapons[v.ID] = &v
	}
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
