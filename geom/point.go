package geom

import (
	"math"

	"github.com/200sc/go-compgeo/printutil"
)

const (
	// POINT_DIM is the number of dimensions a point holds.
	POINT_DIM = 3
)

// A Point is just a N-dimensional point.
type Point [POINT_DIM]float64

// NewPoint returns a Point initialized at the given position
func NewPoint(x, y, z float64) Point {
	return Point{x, y, z}
}

// String converts dp into a string.
func (dp Point) String() string {
	s := ""
	s += "("
	s += printutil.Stringf64(dp[0], dp[1], dp[2])
	s += ")"
	return s
}

// Set sets the value at the given dimension
// on the point.
func (dp Point) Set(i int, v float64) Dimensional {
	dp[i] = v
	return dp
}

// D returns the number of dimensions supported by
// a point.
func (dp Point) D() int {
	return POINT_DIM
}

// Val is equivalent to array access on a point.
func (dp Point) Val(d int) float64 {
	return dp[d]
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

// Eq returns whether two points are equivalent.
func (dp Point) Eq(p2 Dimensional) bool {
	if dp.D() != p2.D() {
		return false
	}
	for i := range dp {
		if dp[i] != p2.Val(i) {
			return false
		}
	}
	return true
}

// Mid2D returns the point in the middle of
// this point and p2.
func (dp Point) Mid2D(p2 D2) Point {
	p3 := Point{}
	for i := range dp {
		p3[i] = (dp[i] + p2.Val(i)) / 2
	}
	return p3
}

// Dot2D performs Dot multiplication on the two
// points, in a two-dimensional context.
func (dp Point) Dot2D(p2 D2) float64 {
	return Dot2D(dp, p2)
}

// Cross2D performs the Cross Product on the three
// points, in a two-dimensional context.
func (dp Point) Cross2D(p2, p3 D2) float64 {
	return Cross2D(dp, p2, p3)
}

// Lesser2D reports the lower point by y value,
// or by x value given equal y values. If the
// two points are equal the latter point is returned.
func (dp Point) Lesser2D(p2 D2) D2 {
	return Lesser2D(dp, p2)
}

// Greater2D reports the higher point by y value,
// or by x value given equal y values. If the
// two points are equal the latter point is returned.
func (dp Point) Greater2D(p2 D2) D2 {
	return Greater2D(dp, p2)
}

// Magnitude2D reports the magnitude of the point
// interpreted as a vector in a two-dimensional context.
func (dp Point) Magnitude2D() float64 {
	return math.Sqrt((dp[0] * dp[0]) + (dp[1] * dp[1]))
}

// Bounds on a Point will return
// the point itself.
func (dp Point) Bounds() S3 {
	return NewSpan(dp)
}

// Lesser2D returns which D2 is
// lexiographically lesser, preferring y to x.
func Lesser2D(p1, p2 D2) D2 {
	if p1.Val(1) < p2.Val(1) {
		return p1
	} else if p1.Val(1) > p2.Val(1) {
		return p2
	}
	if p1.Val(0) < p2.Val(0) {
		return p1
	}
	return p2
}

// Greater2D acts as Lesser2D, but reversed.
func Greater2D(p1, p2 D2) D2 {
	if p1.Val(1) < p2.Val(1) {
		return p2
	} else if p1.Val(1) > p2.Val(1) {
		return p1
	}
	if p1.Val(0) < p2.Val(0) {
		return p2
	}
	return p1

}
