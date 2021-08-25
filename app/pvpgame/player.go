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
		ID:    0,
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
		return playerLevel + rand.Intn(6)
	}
	return playerLevel + rand.Intn(11)
}

func generateBossPoints(level int) (damage, armor, health int) {
	freePoints := level * 5

	pe20 := int(float64(freePoints) * 0.2)
	pe25 := int(float64(freePoints) * 0.25)
	freePoints -= pe20 + pe25 + pe25

	genDmg, genArm, genHp := 0, 0, 0
	genDmg = rand.Intn(freePoints) + 1
	freePoints -= genDmg
	if freePoints > 0 {
		genArm = rand.Intn(freePoints) + 1
	}
	freePoints -= genArm
	genHp = freePoints

	return pe25 + genDmg, pe25 + genArm, pe20*2 + genHp
}
