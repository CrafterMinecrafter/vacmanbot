package pvpgame

type Items struct {
	WeaponID         int
	ArmorID          int
	AvailableWeapons []int
	AvailableArmors  []int
}

type Weapon struct {
	ID             int
	Name           string
	Damage         int
	CritChance     float64
	CritMultiplier float64
}

type Armor struct {
	ID          int
	Name        string
	BonusHealth int
	Protection  int
}

func (g *Game) GenerateWeapon(level int) int {

	return 0
}
