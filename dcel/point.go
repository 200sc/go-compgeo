package dcel

import "github.com/200sc/go-compgeo/printutil"
import "github.com/200sc/go-compgeo/search"

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

func (dp Point) Mid(p2 *Point) *Point {
	p3 := new(Point)
	for i := range dp {
		p3[i] = (dp[i] + (*p2)[i]) / 2
	}
	return p3
}

func (dp Point) VerticalCompare(e *Edge) search.CompareResult {
	p1 := e.Origin
	p2 := e.Twin.Origin
	if p1[0] < p2[0] {
		p1, p2 = p2, p1
	}
	s := (p2[0]-p1[0])*(dp[1]-p1[1]) - (p2[1]-p1[1])*(dp[0]-p1[0])
	if s == 0 {
		return search.Equal
	} else if s < 0 {
		return search.Less
	}
	return search.Greater
}
