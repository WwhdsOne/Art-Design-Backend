package utils

import (
	"strconv"
	"strings"
	"time"
)

func ParseDuration(d string) time.Duration {
	d = strings.TrimSpace(d)
	dr, _ := time.ParseDuration(d)
	if dr != 0 {
		return dr
	}
	if strings.Contains(d, "d") {
		index := strings.Index(d, "d")
		hour, _ := strconv.Atoi(d[:index])
		dr = time.Hour * 24 * time.Duration(hour)
		ndr, _ := time.ParseDuration(d[index+1:])
		return dr + ndr
	}
	dv, _ := strconv.ParseInt(d, 10, 64)
	return time.Duration(dv)
}
