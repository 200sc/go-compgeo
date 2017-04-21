package geom

// Dimensional values can be queried at a
// given dimension for a positional value
type Dimensional interface {
	Val(int) float64
	Set(int, float64) Dimensional
	// D returns the number of dimensions
	// a value has.
	D() int
	Eq(Dimensional) bool
}

// D1 Values are just xes.
type D1 interface {
	Dimensional
	X() float64
}

// D2 Values are dimensional values with
// at least X and Y values.
// Consider whether this should be "TwoD"
type D2 interface {
	D1
	Y() float64
}

// D3 values have Z values on top of D2.
type D3 interface {
	D2
	Z() float64
}

// Spanning types are a length of
// dimensionals.
type Spanning interface {
	D() int
	Len() int
	Set(int, Dimensional) Spanning
	At(int) Dimensional
	Low(int) Dimensional
	High(int) Dimensional
	Eq(Spanning) bool
}

// S1 identifies spans whose points are 1d.
type S1 interface {
	Spanning
	Left() D3
	Right() D3
}

// S2 identifies spans whose points are 2d.
type S2 interface {
	S1
	Top() D3
	Bottom() D3
}

// S3 identifies spans whose points are 3d
type S3 interface {
	S2
	Inner() D3
	Outer() D3
}
