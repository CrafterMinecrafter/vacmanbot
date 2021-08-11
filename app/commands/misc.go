package commands

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const Day = 24 * time.Hour

func lengthToString(l int) string {
	fmtfloat := func(f float64) string {
		s := fmt.Sprintf("%.2f", f)
		return strings.TrimRight(strings.TrimRight(s, "0"), ".")
	}
	if l < 100 {
		return fmt.Sprintf("%v см", l)
	} else if l < 100*1000 {
		return fmt.Sprintf("%v м", fmtfloat(float64(l)/100.0))
	} else {
		return fmt.Sprintf("%v км", fmtfloat(float64(l)/100.0/1000.0))
	}
}

func penisRoll() int {
	base := rand.Float64()
	units := rand.Float64()

	if units > 0.8 {
		return int(base*100.0*1000.0*10.0) + 1
	} else if units > 0.5 {
		return int(base*100.0*999.0) + 1
	} else {
		return int(base*99.0) + 1
	}
}

func fmtName(first, last string) string {
	if len(last) == 0 {
		return first
	}
	return first + " " + last
}

func fmtDuration(tm time.Duration) string {
	tm = tm.Round(time.Second)
	d := tm / Day
	tm -= d * Day
	h := tm / time.Hour
	tm -= h * time.Hour
	m := tm / time.Minute
	tm -= m * time.Minute
	s := tm / time.Second
	str := ""
	if d > 0 {
		str += fmt.Sprintf("%02d д ", d)
	}
	if h > 0 {
		str += fmt.Sprintf("%02d ч ", h)
	}
	if h > 0 || m > 0 {
		str += fmt.Sprintf("%02d мин ", m)
	}
	str += fmt.Sprintf("%02d сек", s)

	return str
}
