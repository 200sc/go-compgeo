package dcel

import "math"

const (
	ε = 1.0e-4
)

func f64eq(f1, f2 float64) bool {
	return math.Abs(f1-f2) <= ε
}

// Credit to David Calhoun at StackOverflow
// http://stackoverflow.com/questions/18390266
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}