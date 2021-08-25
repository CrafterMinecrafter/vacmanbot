package pvpgame

import "fmt"

type BattleLog struct {
	items []string
}

func NewBattleLog() *BattleLog {
	return &BattleLog{
		items: make([]string, 0),
	}
}

func (bl *BattleLog) Append(line string) {
	bl.items = append(bl.items, line)
}

func (bl *BattleLog) Appendf(format string, items ...interface{}) {
	s := fmt.Sprintf(format, items...)
	bl.Append(s)
}

func (bl *BattleLog) Compile() []string {
	return bl.items
}
