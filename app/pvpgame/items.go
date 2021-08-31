package pvpgame

import (
	"math/rand"
	"strings"
)

type Items struct {
	WeaponID       int `json:"weapon_id"`
	ArmorID        int `json:"armor_id"`
	ArchivedWeapon int `json:"arch_weapon"`
	ArchivedArmor  int `json:"arch_armor"`
}

type Weapon struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Damage         int     `json:"damage"`
	CritChance     float64 `json:"crit_chance"`
	CritMultiplier float64 `json:"crit_multiplier"`
}

type Armor struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	BonusHealth int    `json:"bonus_health"`
	Protection  int    `json:"protection"`
}

var (
	BotWeapon *Weapon = &Weapon{ID: -5}
	BotArmor  *Armor  = &Armor{ID: -5}
)

func (g *Game) GenerateItemDrop(level int) (itemType int, itemID int) {
	rnd := rand.Float64()
	if rnd > 0.9 {
		return 1, g.generateWeapon(level)
	} else if rnd > 0.8 {
		return 2, g.generateArmor(level)
	}
	return 0, -1
}

func (g *Game) generateWeapon(level int) int {
	lvl := float64(level)

	damage := int((lvl * 0.4) + ((rand.Float64() * 0.3) * lvl))
	crit := rand.Float64() * 0.3
	critmul := (rand.Float64() * 2.0) + 1
	israre := rand.Float64() >= 0.9

	if israre {
		damage = int((lvl * 0.6) + ((rand.Float64() * 0.5) * lvl))
		crit = rand.Float64()*0.5 + 0.1
		critmul = 2 + rand.Float64()*2.0
	}

	w := Weapon{
		ID:             -5,
		Name:           generateWeaponName(israre),
		Damage:         damage,
		CritChance:     crit,
		CritMultiplier: critmul,
	}

	g.db.Bucket("pvp_weapons")
	w.ID = g.db.NextID()
	g.db.Put(w.ID, w)

	return w.ID
}

func (g *Game) generateArmor(level int) int {
	lvl := float64(level)

	protection := int((lvl * 0.4) + ((rand.Float64() * 0.3) * lvl))
	health := int(lvl + (rand.Float64() * (lvl * 1.5)))
	israre := rand.Float64() >= 0.9

	if israre {
		protection = int(float64(protection) * 1.3)
		health = int(float64(health) * 1.3)
	}

	a := Armor{
		ID:          -5,
		Name:        generateArmorName(israre),
		BonusHealth: health,
		Protection:  protection,
	}

	g.db.Bucket("pvp_armors")
	a.ID = g.db.NextID()
	g.db.Put(a.ID, a)

	return a.ID
}

func generateWeaponName(isRare bool) string {
	a1 := []string{"пунцовый", "малиновый", "синий", "лиловый", "сизый", "красный", "лазурный", "светлый",
		"черный", "бледный", "большой", "маленький", "узкий", "длинный", "низкий", "широкий", "высокий", "треугольный",
		"круглый", "овальный", "квадратный", "прямой", "извилистый", "бодрый", "слабый", "сильный", "здоровый", "вкусный",
		"холодный", "теплый", "горячий", "сырой", "сухой", "молодой", "старый", "древний", "дряхлый", "пожилой", "крепкий",
		"юный", "слепой", "хромой", "лысый", "добрый", "злой", "сердитый", "жестокий", "верный", "милый", "умный", "глупый",
		"честный", "строгий", "скромный", "хитрый"}
	a2 := []string{"нож", "меч", "клинок", "кастет", "автомат", "гранатомёт", "тесак", "пулемёт", "лук", "бумеранг", "бластер", "посох"}
	a3 := []string{"истребления", "гибели", "вымирания", "ликвидации", "изъятия", "подавления", "нейтрализации", "подрыва",
		"диверсии", "разрушения", "ослабления", "осквернения", "сокрушения", "аннулирования", "удаления"}

	r1 := a1[rand.Intn(len(a1))]
	r2 := a2[rand.Intn(len(a2))]
	r3 := a3[rand.Intn(len(a3))]

	if isRare {
		r1 = "легендарный"
	}

	return strings.Title(r1 + " " + r2 + " " + r3)
}

func generateArmorName(isRare bool) string {
	a1 := []string{"броня", "куртка", "обшивка", "оболочка", "кираса", "латы", "кольчуга", "рубаха"}
	a2 := []string{"силы", "здоровья", "крепости", "могущества", "мощи", "власти",
		"вескости", "повелительства", "всесилия", "матёрости", "плотности", "слабости", "бессилия", "малосильности",
		"расслабления", "истощения", "разбитости", "дряблости", "бесцветности"}

	r1 := a1[rand.Intn(len(a1))]
	r2 := a2[rand.Intn(len(a2))]

	if isRare {
		r2 = "легендарности"
	}

	return strings.Title(r1 + " " + r2)
}
