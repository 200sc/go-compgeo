// package geom holds geometric primitives and interfaces.

package geom

import "math"

const (
	// Epsilon could probably be smaller than this without
	// causing problems, but we're being overly cautious.
	ε = 1.0e-7
	// Inf is shorthand for math.MaxFloat64
	Inf = math.MaxFloat64
	// NegInf is shorthand for math.MaxFloat64 * -1
	NegInf = -math.MaxFloat64
)

// F64eq returns whether two input float64s have
// a lesser difference than epsilon.
func F64eq(f1, f2 float64) bool {
	return math.Abs(f1-f2) <= ε
}

// Credit to David Calhoun at StackOverflow
// http://stackoverflow.com/questions/18390266

// Round converts a float64 to an integer.
func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// ToFixed converts a float64 to a float64
// with a limited number of decimal places
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

// End Credit
