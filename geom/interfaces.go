package geom

// Dimensional values can be queried at a
// given dimension for a positional value
type Dimensional interface {
	Val(int) float64
	// D returns the number of dimensions
	// a value has.
	D() int
	Eq(Dimensional) bool
}

// D2 Values are dimensional values with
// at least X and Y values.
// Consider whether this should be "TwoD"
type D2 interface {
	Dimensional
	X() float64
	Y() float64
}

// D3 values have Z values on top of D2.
type D3 interface {
	D2
	Z() float64
}
