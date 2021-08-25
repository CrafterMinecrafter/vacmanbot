package pvpgame

import (
	"math"
	"math/rand"
)

type Fight struct {
	game *Game
	P    []*Player
	N    []string
	Log  *BattleLog
}

func (f *Fight) Execute() []string {
	log := NewBattleLog()

	// собираем нужные нам данные
	weapons := []*Weapon{f.game.weapons[f.P[0].Items.WeaponID], f.game.weapons[f.P[1].Items.WeaponID]}
	armors := []*Armor{f.game.armors[f.P[0].Items.ArmorID], f.game.armors[f.P[1].Items.ArmorID]}

	// аппендим  сообщение о начале боя
	log.Appendf(TextFightStart, f.N[0], f.P[0].Stats.Level, f.N[1], f.P[1].Stats.Level)
	// аппендим информацию об игроках
	log.Appendf(TextPlayerInfo, f.N[0], f.P[0].Stats.Damage, f.P[0].Stats.Protection, f.P[0].Stats.Health,
		weapons[0].Name, weapons[0].Damage, armors[0].Name, armors[0].Protection, armors[0].BonusHealth)
	log.Appendf(TextPlayerInfo, f.N[1], f.P[1].Stats.Damage, f.P[1].Stats.Protection, f.P[1].Stats.Health,
		weapons[1].Name, weapons[1].Damage, armors[1].Name, armors[1].Protection, armors[1].BonusHealth)
	// TODO: remove this
	log.Append("")

	// проверяем, возможен ли бой вообще
	if !f.canFight() {
		log.Append(TextCantFight)
		return log.Compile()
	}

	// сражаемся
	turn := rand.Intn(2)
	health := []int{
		f.P[0].Stats.Health + armors[0].BonusHealth,
		f.P[1].Stats.Health + armors[1].BonusHealth,
	}
	for {
		// считаем урон
		iscrit := isCrit(weapons[turn].CritChance)
		damage := f.P[turn].Stats.Damage
		if iscrit {
			damage += calcCritDamage(weapons[turn].Damage, weapons[turn].CritMultiplier)
		} else {
			damage += weapons[turn].Damage
		}

		// считаем защиту
		protection := f.P[1-turn].Stats.Protection + armors[1-turn].Protection

		// считаем урон, который прошел
		rawdamage := damage - protection
		if rawdamage < 0 {
			rawdamage = 0
		}

		// вычитаем хп
		health[1-turn] -= rawdamage

		// генерируем сводку
		if iscrit {
			if health[1-turn] <= 0 {
				log.Appendf(TextCritKill, f.N[turn], f.N[1-turn], rawdamage)
				break
			} else {
				log.Appendf(TextCrit, f.N[turn], rawdamage, f.N[1-turn], health[1-turn])
			}
		} else {
			if health[1-turn] <= 0 {
				log.Appendf(TextKill, f.N[turn], f.N[1-turn], rawdamage)
				break
			} else {
				log.Appendf(TextDamage, f.N[turn], rawdamage, f.N[1-turn], health[1-turn])
			}
		}

		// меняем игрока
		turn = 1 - turn
	}
	// TODO: REMOVE
	log.Append("")

	// определяем победителя и проигравшего
	winner, loser := turn, 1-turn

	// генерируем статистику
	if winner == 0 && f.P[loser].IsBot {
		log.Appendf(TextInfoPlayerWinsBoss, f.N[winner], f.N[loser])
	} else {
		log.Appendf(TextInfoPlayerWins, f.N[winner], f.N[loser])
	}
	if winner != 0 {
		// победил не зачинщик драки. Опыт не считаем. Только правим статистику и эло.
		if f.P[winner].IsBot {
			// победил бот
			f.P[loser].Statistics.Fights++
			f.P[loser].Statistics.Loses++
		} else {
			// победил человек
			winnerElo, loserElo := calcElo(f.P[winner].Statistics.Elo, f.P[loser].Statistics.Elo)
			log.Appendf(TextInfoEloChanges, f.N[winner], winnerElo, f.N[loser], loserElo)
			f.P[winner].Statistics.Elo += winnerElo
			f.P[loser].Statistics.Elo += loserElo
			f.P[winner].Statistics.Fights++
			f.P[winner].Statistics.Wins++
			f.P[loser].Statistics.Fights++
			f.P[loser].Statistics.Loses++
		}
	} else {
		// победил игрок, который начал драку. Считаем ему опыт и всё остальное.
		f.P[winner].Statistics.Fights++
		f.P[winner].Statistics.Wins++
		if !f.P[loser].IsBot {
			winnerElo, loserElo := calcElo(f.P[winner].Statistics.Elo, f.P[loser].Statistics.Elo)
			log.Appendf(TextInfoEloChanges, f.N[winner], winnerElo, f.N[loser], loserElo)
			f.P[winner].Statistics.Elo += winnerElo
			f.P[loser].Statistics.Elo += loserElo
			f.P[loser].Statistics.Fights++
			f.P[loser].Statistics.Loses++
		}
		expgain := calcPrizeExp(f.P[loser].Stats.Level)
		f.P[winner].Stats.Experience += expgain
		log.Appendf(TextInfoExpChanges, f.N[winner], expgain, f.P[winner].Stats.Experience, calcLevelToExp(f.P[winner].Stats.Level+1))
		oldlevel := f.P[winner].Stats.Level
		newlevel := calcExpToLevel(f.P[winner].Stats.Experience)
		if newlevel > oldlevel {
			f.P[winner].Stats.Level = newlevel
			for i := 0; i < newlevel-oldlevel; i++ {
				f.P[winner].Stats.UnusedPoints += 5
			}
			log.Appendf(TextInfoNewLevel, f.N[winner], newlevel)
			log.Appendf(TextInfoPointsAvalible, f.P[winner].Stats.UnusedPoints)
		}
	}

	if !f.P[winner].IsBot {
		f.game.db.Bucket("pvp_users")
		f.game.db.Put(f.P[winner].ID, f.P[winner])
	}
	if !f.P[loser].IsBot {
		f.game.db.Bucket("pvp_users")
		f.game.db.Put(f.P[loser].ID, f.P[loser])
	}

	return log.Compile()
}

func (f *Fight) canFight() bool {
	damage1 := f.P[0].Stats.Damage
	damage2 := f.P[1].Stats.Damage
	protection1 := f.P[0].Stats.Protection
	protection2 := f.P[1].Stats.Protection

	damage1 += calcCritDamage(f.game.weapons[f.P[0].Items.WeaponID].Damage, f.game.weapons[f.P[0].Items.WeaponID].CritMultiplier)
	damage2 += calcCritDamage(f.game.weapons[f.P[1].Items.WeaponID].Damage, f.game.weapons[f.P[1].Items.WeaponID].CritMultiplier)
	protection1 += f.game.armors[f.P[0].Items.ArmorID].Protection
	protection2 += f.game.armors[f.P[1].Items.ArmorID].Protection

	return damage1 > protection2 || damage2 > protection1
}

func isCrit(critChance float64) bool {
	return rand.Float64() <= critChance
}

func calcCritDamage(damage int, critMultiplier float64) int {
	d := float64(damage) * critMultiplier
	return int(math.Round(d))
}
