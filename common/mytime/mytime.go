package mytime

import (
	"strconv"
	"time"
)

func TomorrowZero() time.Time {
	now := time.Now()

	h, _ := time.ParseDuration("-" + strconv.Itoa(now.Hour()) + "h")
	m, _ := time.ParseDuration("-" + strconv.Itoa(now.Minute()) + "m")
	s, _ := time.ParseDuration("-" + strconv.Itoa(now.Second()) + "s")
	d, _ := time.ParseDuration("24h")

	return now.Add(h).Add(m).Add(s).Add(d)
}
