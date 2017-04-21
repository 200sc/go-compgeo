package geom

import compgeo "github.com/200sc/go-compgeo"

// FullEdge is (eventually) a span
// across two arbitrary points
type FullEdge [2]Point

// NewFullEdge type casts interfaces to points to hold in
// a FullEdge
func NewFullEdge(p1, p2 D3) FullEdge {
	return FullEdge{p1.(Point), p2.(Point)}
}

// Slope returns the slope on this
// FullEdge. This slope will always
// have a positive X value.
func (fe FullEdge) Slope() float64 {
	p1 := fe.Left().(D2)
	p2 := fe.Right().(D2)
	if p1.X() == p2.X() {
		return Inf
	}
	return (p2.Y() - p1.Y()) / (p2.X() - p2.X())
}

// PointAt returns the value along this edge
// at v in the dth dimension.
func (fe FullEdge) PointAt(d int, v float64) (Point, error) {
	e1 := fe.Low(d)
	e2 := fe.High(d)
	if v < e1.Val(d) || v > e2.Val(d) {
		return Point{}, compgeo.RangeError{}
	}
	v -= e1.Val(d)
	span := e2.Val(d) - e1.Val(d)
	if span == 0 {
		return Point{}, compgeo.DivideByZero{}
	}
	portion := v / span
	p := Point{}
	for i := 0; i < len(p); i++ {
		if i == d {
			p[i] = v
		} else {
			p[i] = e1.Val(i) + (portion * (e2.Val(i) - e1.Val(i)))
		}
	}
	return p, nil
}

// SubEdge returns the portion of this edge
// from pointAt(d,v1) to pointAt(d,v2)
func (fe FullEdge) SubEdge(d int, v1, v2 float64) (FullEdge, error) {
	p1, err := fe.PointAt(d, v1)
	if err != nil {
		return FullEdge{}, err
	}
	p2, err := fe.PointAt(d, v2)
	if err != nil {
		return FullEdge{}, err
	}
	return FullEdge{p1, p2}, nil
}

// D returns the number of dimensions in this edge.
func (fe FullEdge) D() int {
	return fe[0].D()
}

// Len returns the number of elements on this edge (2)
func (fe FullEdge) Len() int {
	return 2
}

// At returns the value at a given index on this edge.
func (fe FullEdge) At(i int) Dimensional {
	return fe[i]
}

// Low returns whatever value is lower at dimension d
// on this edge
func (fe FullEdge) Low(d int) Dimensional {
	if fe[0].Val(d) <= fe[1].Val(d) {
		return fe[0]
	}
	return fe[1]
}

// High returns whatever value is higher at dimension d
// on this edge
func (fe FullEdge) High(d int) Dimensional {
	if fe[0].Val(d) > fe[1].Val(d) {
		return fe[0]
	}
	return fe[1]
}

// Left renames Low(0)
func (fe FullEdge) Left() D3 {
	return fe.Low(0).(D3)
}

// Right renames High(0)
func (fe FullEdge) Right() D3 {
	return fe.High(0).(D3)
}

// Top renames High(1)
func (fe FullEdge) Top() D3 {
	return fe.High(1).(D3)
}

// Bottom renames Low(1)
func (fe FullEdge) Bottom() D3 {
	return fe.Low(1).(D3)
}

// Inner renames High(1)
func (fe FullEdge) Inner() D3 {
	return fe.High(2).(D3)
}

// Outer renames Low(2)
func (fe FullEdge) Outer() D3 {
	return fe.Low(2).(D3)
}

// Eq returns whether this edge is equivalent
// to another spanning type
func (fe FullEdge) Eq(s Spanning) bool {
	if fe.D() != s.D() {
		return false
	}
	for i := 0; i < fe.D(); i++ {
		if !fe.At(i).Eq(s.At(i)) {
			return false
		}
	}
	return true
}

// Set returns this edge with a value at i
// changed to d
func (fe FullEdge) Set(i int, d Dimensional) Spanning {
	fe[i] = d.(Point)
	return fe
}
