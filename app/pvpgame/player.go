package pvpgame

import "math/rand"

type Player struct {
	ID         int              `json:"id"`
	IsBot      bool             `json:"is_bot"`
	Stats      PlayerStats      `json:"stats"`
	Statistics PlayerStatistics `json:"statistics"`
	Items      Items            `json:"items"`
}

func CreatePlayer(id int, isbot bool) *Player {
	return &Player{
		ID:    id,
		IsBot: isbot,
		Stats: PlayerStats{
			Level:        0,
			Experience:   0,
			Damage:       5,
			Protection:   1,
			Health:       10,
			UnusedPoints: 0,
			Gold:         0,
		},
		Statistics: PlayerStatistics{
			Fights: 0,
			Wins:   0,
			Loses:  0,
			Elo:    1500,
		},
		Items: Items{
			WeaponID:       -1,
			ArmorID:        -1,
			ArchivedWeapon: -1,
			ArchivedArmor:  -1,
		},
	}
}

func (g *Game) CreateBossFor(player *Player) *Player {
	lvl := float64(player.Stats.Level)
	winrate := 0.3
	if lvl < 3 {
		winrate = 0.0
	} else if lvl < 10 {
		winrate = 0.7
	} else if lvl < 50 {
		winrate = 0.5
	} else if lvl < 100 {
		winrate = 0.4
	}

	vac := !(rand.Float64() <= winrate)

	damage := is(vac, scale(player.Stats.Damage, 1.1), scale(player.Stats.Damage, 0.9))
	protection := is(vac, scale(player.Stats.Protection, 1.1), scale(player.Stats.Protection, 0.9))
	health := is(vac, scale(player.Stats.Health, 1.3), scale(player.Stats.Health, 0.8))
	level := ((damage + protection + health) - 16) / 5
	experience := calcLevelToExp(level)

	result := &Player{
		ID:    -1,
		IsBot: true,
		Stats: PlayerStats{
			Level:        level,
			Experience:   experience,
			Damage:       damage,
			Protection:   protection,
			Health:       health,
			UnusedPoints: 0,
			Gold:         0,
		},
		Statistics: PlayerStatistics{
			Fights: 0,
			Wins:   0,
			Loses:  0,
			Elo:    1500,
		},
		Items: Items{
			WeaponID:       -5,
			ArmorID:        -5,
			ArchivedWeapon: -1,
			ArchivedArmor:  -1,
		},
	}

	playerWeapon := Weapon{}
	if player.Items.WeaponID > -1 {
		g.db.Bucket("pvp_weapons")
		g.db.Get(player.Items.WeaponID, &playerWeapon)
	}
	playerArmor := Armor{}
	if player.Items.ArmorID > -1 {
		g.db.Bucket("pvp_armors")
		g.db.Get(player.Items.ArmorID, &playerArmor)
	}

	weapDamage := is(vac, scale(playerWeapon.Damage, 1.2), scale(playerWeapon.Damage, 0.9))
	weapCritChance := playerWeapon.CritChance
	weapCritDamage := isf(rand.Float64() >= 0.5, playerWeapon.CritMultiplier*1.2, playerWeapon.CritMultiplier*0.9)

	armProtection := is(vac, scale(playerArmor.Protection, 1.2), scale(playerArmor.Protection, 0.9))
	armHealth := is(vac, scale(playerArmor.BonusHealth, 1.2), scale(playerArmor.BonusHealth, 0.9))

	BotWeapon = &Weapon{
		ID:             -5,
		Name:           generateWeaponName(false),
		Damage:         weapDamage,
		CritChance:     weapCritChance,
		CritMultiplier: weapCritDamage,
	}

	BotArmor = &Armor{
		ID:          -5,
		Name:        generateArmorName(false),
		BonusHealth: armHealth,
		Protection:  armProtection,
	}

	return result
}

func is(x bool, a, b int) int {
	if x {
		return a
	}
	return b
}

func isf(x bool, a, b float64) float64 {
	if x {
		return a
	}
	return b
}

func scale(x int, factor float64) int {
	return int(float64(x) * factor)
}
