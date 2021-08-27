package main

import (
	"fmt"

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
	result := fight.Execute().String()

	fmt.Println(result)
}
