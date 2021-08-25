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

func calcPrizeExp(enemyLevel int) int {
	nl := float64(calcLevelToExp(enemyLevel + 1))
	cl := float64(calcLevelToExp(enemyLevel))
	x := math.Round((nl - cl) / 3.0)
	return int(x)
}

func calcElo(winner, loser int) (w, l int) {
	rA, rB := float64(winner), float64(loser)

	eA := 1 / (1 + math.Pow(10, ((rB-rA)/400)))
	eB := 1 / (1 + math.Pow(10, ((rA-rB)/400)))

	nrA := 64 * (1.0 - eA)
	nrB := 64 * (0.0 - eB)

	return int(nrA), int(nrB)
}
