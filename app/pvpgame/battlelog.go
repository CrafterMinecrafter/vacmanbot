package pvpgame

import (
	"fmt"
	"strings"
)

type BattleLog struct {
	prelog   []string
	fightlog []string
	postlog  []string
}

func NewBattleLog() *BattleLog {
	return &BattleLog{
		prelog:   make([]string, 0),
		fightlog: make([]string, 0),
		postlog:  make([]string, 0),
	}
}

func (bl *BattleLog) AppendPre(line string) {
	bl.prelog = append(bl.prelog, line)
}

func (bl *BattleLog) AppendFight(line string) {
	bl.fightlog = append(bl.fightlog, line)
}

func (bl *BattleLog) AppendPost(line string) {
	bl.postlog = append(bl.postlog, line)
}

func (bl *BattleLog) AppendfPre(format string, items ...interface{}) {
	s := fmt.Sprintf(format, items...)
	bl.AppendPre(s)
}

func (bl *BattleLog) AppendfFight(format string, items ...interface{}) {
	s := fmt.Sprintf(format, items...)
	bl.AppendFight(s)
}

func (bl *BattleLog) AppendfPost(format string, items ...interface{}) {
	s := fmt.Sprintf(format, items...)
	bl.AppendPost(s)
}

func (bl *BattleLog) StringSlicePre() []string {
	return bl.prelog
}

func (bl *BattleLog) StringSliceFight() []string {
	return bl.fightlog
}

func (bl *BattleLog) StringSlicePost() []string {
	return bl.postlog
}

func (bl *BattleLog) StringPre() string {
	return strings.Join(bl.prelog, "\n")
}

func (bl *BattleLog) StringFight() string {
	return strings.Join(bl.fightlog, "\n")
}
func (bl *BattleLog) StringPost() string {
	return strings.Join(bl.postlog, "\n")
}
