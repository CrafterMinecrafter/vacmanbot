package pvpgame

import (
	"fmt"
	"strings"
)

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

func (bl *BattleLog) StringSlice() []string {
	return bl.items
}

func (bl *BattleLog) String() string {
	return strings.Join(bl.items, "\n")
}
