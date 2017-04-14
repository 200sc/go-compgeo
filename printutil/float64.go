package printutil

import (
	"math"
	"strconv"
)

// Stringf64 is a helper function which will take in
// any number of float65s and return a standard-formatted
// string.
func Stringf64(ks ...float64) string {
	s := ""
	for i, k := range ks {
		if k == math.MaxFloat64*-1 {
			s += "-∞"
		} else if k == math.MaxFloat64 {
			s += "∞"
		} else {
			s += strconv.FormatFloat(k, 'f', 5, 64)
		}
		if i != len(ks)-1 {
			s += ", "
		}
	}
	return s
}
