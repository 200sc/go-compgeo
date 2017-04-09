package printutil

import (
	"math"
	"strconv"
)

func Stringf64(ks ...float64) string {
	s := ""
	for i, k := range ks {
		if k == math.MaxFloat64*-1 {
			s += "-∞"
		} else if k == math.MaxFloat64 {
			s += "∞"
		} else {
			s += strconv.FormatFloat(k, 'f', -1, 64)
		}
		if i != len(ks)-1 {
			s += ", "
		}
	}
	return s
}
