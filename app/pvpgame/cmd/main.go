package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nyan2d/bolteo"
	"github.com/nyan2d/vacmanbot/app/pvpgame"
)

func main() {
	db := bolteo.MustOpen("test.db")
	game := pvpgame.NewGame(db)

	p1 := game.GetPlayer(0)
	p2 := game.CreateBossFor(p1)
	p1 = game.CreateBossFor(p2)

	fight := game.MakeFight(p1, p2, "player", "bot")
	result := fight.Execute()

	fmt.Println(result.StringPre())
	fmt.Println(result.StringFight())
	fmt.Println(result.StringPost())
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
