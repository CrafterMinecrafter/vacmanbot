package pvpgame

import (
	"math"
	"math/rand"
)

type Fight struct {
	game *Game
	P    []*Player
	N    []string
	log  *BattleLog
}

func (f *Fight) Execute() *BattleLog {
	log := NewBattleLog()

	// собираем нужные нам данные
	wea1, wea2 := &Weapon{Name: "Отсутствует"}, &Weapon{Name: "Отсутствует"}
	arm1, arm2 := &Armor{Name: "Отсутствует"}, &Armor{Name: "Отсутствует"}

	if !f.P[0].IsBot {
		f.game.db.Bucket("pvp_weapons")
		f.game.db.Get(f.P[0].Items.WeaponID, wea1)
		f.game.db.Bucket("pvp_armors")
		f.game.db.Get(f.P[0].Items.ArmorID, arm1)
	} else {
		if f.P[0].Stats.Level >= 3 {
			wea1 = BotWeapon
			arm1 = BotArmor
		}
	}
	if !f.P[1].IsBot {
		f.game.db.Bucket("pvp_weapons")
		f.game.db.Get(f.P[1].Items.WeaponID, wea2)
		f.game.db.Bucket("pvp_armors")
		f.game.db.Get(f.P[1].Items.ArmorID, arm2)
	} else {
		if f.P[1].Stats.Level >= 3 {
			wea2 = BotWeapon
			arm2 = BotArmor
		}
	}

	weapons := []*Weapon{wea1, wea2}
	armors := []*Armor{arm1, arm2}

	// аппендим  сообщение о начале боя
	log.AppendfPre(TextFightStart, f.N[0], f.P[0].Stats.Level, f.N[1], f.P[1].Stats.Level)
	log.AppendPre("")
	// аппендим информацию об игроках
	log.AppendfPre(TextPlayerInfo, f.N[0], f.P[0].Stats.Damage, f.P[0].Stats.Protection, f.P[0].Stats.Health,
		weapons[0].Name, weapons[0].Damage, int(weapons[0].CritChance*100.0), int(weapons[0].CritMultiplier*100.0),
		armors[0].Name, armors[0].Protection, armors[0].BonusHealth)
	log.AppendPre("")
	log.AppendfPre(TextPlayerInfo, f.N[1], f.P[1].Stats.Damage, f.P[1].Stats.Protection, f.P[1].Stats.Health,
		weapons[1].Name, weapons[1].Damage, int(weapons[1].CritChance*100.0), int(weapons[1].CritMultiplier*100.0),
		armors[1].Name, armors[1].Protection, armors[1].BonusHealth)

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

		// считаем вероятность блока
		blockChance := 0.7 * (float64(protection) / float64(damage))
		if blockChance > 0.7 {
			blockChance = 0.7
		}
		if rand.Float64() <= blockChance {
			log.AppendfFight(TextBlocked, f.N[turn], f.N[1-turn])
			turn = 1 - turn
			continue
		}

		// считаем урон, который прошел
		truedamage := int(float64(damage) * 0.2)
		rawdamage := (damage - truedamage) - protection
		if rawdamage < 0 {
			rawdamage = 0
		}

		// вычитаем хп
		health[1-turn] -= truedamage + rawdamage

		// генерируем сводку
		if iscrit {
			if health[1-turn] <= 0 {
				log.AppendfFight(TextCritKill, f.N[turn], f.N[1-turn], truedamage+rawdamage)
				break
			} else {
				log.AppendfFight(TextCrit, f.N[turn], truedamage+rawdamage, f.N[1-turn], health[1-turn])
			}
		} else {
			if health[1-turn] <= 0 {
				log.AppendfFight(TextKill, f.N[turn], f.N[1-turn], truedamage+rawdamage)
				break
			} else {
				log.AppendfFight(TextDamage, f.N[turn], truedamage+rawdamage, f.N[1-turn], health[1-turn])
			}
		}

		// меняем игрока
		turn = 1 - turn
	}

	// определяем победителя и проигравшего
	winner, loser := turn, 1-turn

	// генерируем статистику
	if winner == 0 && f.P[loser].IsBot {
		log.AppendfPost(TextInfoPlayerWinsBoss, f.N[winner], f.N[loser])
	} else {
		log.AppendfPost(TextInfoPlayerWins, f.N[winner], f.N[loser])
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
			log.AppendfPost(TextInfoEloChanges, f.N[winner], winnerElo, f.N[loser], loserElo)
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
			log.AppendfPost(TextInfoEloChanges, f.N[winner], winnerElo, f.N[loser], loserElo)
			f.P[winner].Statistics.Elo += winnerElo
			f.P[loser].Statistics.Elo += loserElo
			f.P[loser].Statistics.Fights++
			f.P[loser].Statistics.Loses++
		}
		expgain := calcPrizeExp(f.P[winner].Stats.Level, f.P[loser].Stats.Level)
		f.P[winner].Stats.Experience += expgain
		log.AppendfPost(TextInfoExpChanges, f.N[winner], expgain, f.P[winner].Stats.Experience, calcLevelToExp(f.P[winner].Stats.Level+1))
		oldlevel := f.P[winner].Stats.Level
		newlevel := calcExpToLevel(f.P[winner].Stats.Experience)
		if newlevel > oldlevel {
			f.P[winner].Stats.Level = newlevel
			for i := 0; i < newlevel-oldlevel; i++ {
				f.P[winner].Stats.UnusedPoints += 5
			}
			log.AppendfPost(TextInfoNewLevel, f.N[winner], newlevel)
			log.AppendfPost(TextInfoPointsAvalible, f.P[winner].Stats.UnusedPoints)
		}
		itemtype, itemid := f.game.GenerateItemDrop(f.P[winner].Stats.Level)
		switch itemtype {
		case 1:
			f.P[winner].Items.ArchivedWeapon = itemid
			nweap := &Weapon{}
			f.game.db.Bucket("pvp_weapons")
			f.game.db.Get(itemid, nweap)
			log.AppendfPost(TextNewWeapon, nweap.Name, nweap.Damage, int(nweap.CritChance*100.0), int(nweap.CritMultiplier*100.0))
			log.AppendfPost(TextSwitchToNew, "switchweapon", "новое оружие")
		case 2:
			f.P[winner].Items.ArchivedArmor = itemid
			narm := &Armor{}
			f.game.db.Bucket("pvp_armors")
			f.game.db.Get(itemid, narm)
			log.AppendfPost(TextNewArmor, narm.Name, narm.Protection, narm.BonusHealth)
			log.AppendfPost(TextSwitchToNew, "switcharmor", "новую броню")
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

	return log
}

func isCrit(critChance float64) bool {
	return rand.Float64() <= critChance
}

func calcCritDamage(damage int, critMultiplier float64) int {
	d := float64(damage) * critMultiplier
	return int(math.Round(d))
}
