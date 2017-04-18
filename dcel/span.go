package dcel

import "math"

// A Span represents n-dimensions of
// span from one point to another
// for however many dimensions a point
// has (3 at time or writing)
type Span struct {
	Min Point
	Max Point
}

// NewSpan returns a span with its values
// set to appropriate infinities
func NewSpan() Span {
	sp := Span{}
	for i := range sp.Min {
		sp.Min[i] = math.MaxFloat64
		sp.Max[i] = math.MaxFloat64 * -1
	}
	return sp
}

// Expand on a Span will reduce or increase
// a span's min and max values if the input points
// falls outside of the span on any dimension
func (sp Span) Expand(ps ...Dimensional) Span {
	for _, p := range ps {
		for i := range sp.Min {
			v := p.Val(i)
			if v < sp.Min[i] {
				sp.Min[i] = v
			}
			if v > sp.Max[i] {
				sp.Max[i] = v
			}
		}
	}
	return sp
}

func (sp Span) Trapezoid() *Trapezoid {
	return nil
}

func (sp Span) Lesser(d int) Point {
	return sp.Min
}

func (sp Span) Greater(d int) Point {
	return sp.Max
}

func (sp Span) Left() Point {
	return sp.Lesser(0)
}

func (sp Span) Right() Point {
	return sp.Greater(0)
}

func (sp Span) Top() Point {
	return sp.Greater(1)
}

func (sp Span) Bottom() Point {
	return sp.Lesser(1)
}

func (sp Span) Outer() Point {
	return sp.Lesser(2)
}

func (sp Span) Inner() Point {
	return sp.Greater(2)
}

// Bounds on a DCEL returns a Span calculated
// from every point in the DCEL.
func (dc *DCEL) Bounds() Span {
	sp := NewSpan()
	for _, v := range dc.Vertices {
		sp.Expand(v)
	}
	return sp
}

// Bounds returns a Span calculated from
// every point on the Inner of this face
// because at time of writing we don't
// populate Outer
func (f *Face) Bounds() Span {
	sp := NewSpan()
	e := f.Inner
	sp = sp.Expand(e.Origin)
	for e.Next != f.Inner {
		e = e.Next
		sp = sp.Expand(e.Origin)
	}
	return sp
}

// Bounds on a Point will return
// the point itself.
func (p *Point) Bounds() Span {
	return Span{*p, *p}
}

// Bounds on an Edge returns a Span
// on the edge's origin and it's twin's
// origin.
func (e *Edge) Bounds() Span {
	sp := NewSpan()
	sp.Expand(e.Origin, e.Twin.Origin)
	return sp
}
