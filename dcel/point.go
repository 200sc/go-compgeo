package dcel

import "github.com/200sc/go-compgeo/printutil"

// A Point is just a 3-dimensional point.
type Point [3]float64

// NewPoint returns a Point initialized at the given position
func NewPoint(x, y, z float64) *Point {
	return &Point{x, y, z}
}

// String converts dp into a string.
func (dp Point) String() string {
	s := ""
	s += "("
	s += printutil.Stringf64(dp[0], dp[1], dp[2])
	s += ")"
	return s
}

// X :
// Get the value of this point on the x axis
func (dp Point) X() float64 {
	return dp[0]
}

// Y :
// Get the value of this point on the y axis
func (dp Point) Y() float64 {
	return dp[1]
}

// Z :
// Get the value of this point on the z axis
func (dp Point) Z() float64 {
	return dp[2]
}
