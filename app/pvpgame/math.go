package pvpgame

import (
	"math"
)

func calcLevelToExp(lvl int) int {
	w := float64(lvl)
	x := w * 123.0
	return int(math.Round(x))
}

func calcExpToLevel(exp int) int {
	e := float64(exp)
	x := e / 123
	return int(x)
}

func calcPrizeExp(yourLevel, enemyLevel int) int {
	a := float64(yourLevel)
	b := float64(enemyLevel)
	x := 123.0
	exp := ((b + 1) / (a + 1)) * (x / 3)

	return int(exp)
}

func calcElo(winner, loser int) (w, l int) {
	rA, rB := float64(winner), float64(loser)

	eA := 1 / (1 + math.Pow(10, ((rB-rA)/400)))
	eB := 1 / (1 + math.Pow(10, ((rA-rB)/400)))

	nrA := 64 * (1.0 - eA)
	nrB := 64 * (0.0 - eB)

	return int(nrA), int(nrB)
}
