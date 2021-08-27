package pvpgame

import (
	"math/rand"
)

type Player struct {
	ID         int
	IsBot      bool
	Stats      PlayerStats
	Statistics PlayerStatistics
	Items      Items
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
			WeaponID:         -1,
			ArmorID:          -1,
			AvailableWeapons: []int{},
			AvailableArmors:  []int{},
		},
	}
}

func (g *Game) CreateBossFor(player *Player) *Player {
	bossLevel := generateBossLevel(player.Stats.Level)
	bossDamage, bossArmor, bossHealth := generateBossPoints(bossLevel)

	p := Player{
		ID:    -1,
		IsBot: true,
		Stats: PlayerStats{
			Level:        bossLevel,
			Experience:   calcLevelToExp(bossLevel),
			Damage:       bossDamage + g.weapons[player.Items.WeaponID].Damage,
			Protection:   bossArmor + g.armors[player.Items.ArmorID].Protection,
			Health:       bossHealth + g.armors[player.Items.ArmorID].BonusHealth,
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
			WeaponID:         -5,
			ArmorID:          -5,
			AvailableWeapons: []int{},
			AvailableArmors:  []int{},
		},
	}
	return &p
}

func generateBossLevel(playerLevel int) int {
	if playerLevel < 10 {
		return playerLevel + 1
	}
	if playerLevel < 20 {
		return playerLevel + 3
	}
	if playerLevel < 100 {
		return playerLevel + rand.Intn(5) + 1
	}
	return playerLevel + rand.Intn(10) + 1
}

func generateBossPoints(level int) (damage, armor, health int) {
	freePoints := level * 5

	dmg, arm, hp := 5, 1, 10
	rand30p := int(float64(freePoints) * 0.3)
	rand20p := int(float64(freePoints) * 0.2)
	freePoints -= rand30p + rand20p + rand20p
	dmg += rand20p
	arm += rand20p
	hp += rand30p * 2

	if freePoints > 0 {
		rnddmg := rand.Intn(freePoints)
		freePoints -= rnddmg
		dmg += rnddmg
	}

	if freePoints > 0 {
		rndarm := rand.Intn(freePoints)
		freePoints -= rndarm
		arm += rndarm
	}

	if freePoints > 0 {
		hp += freePoints * 2
	}

	return dmg, arm, hp
}
