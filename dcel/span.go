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
// a span's min and max values if the input point
// falls outside of the span on any dimension
func (sp Span) Expand(p Dimensional) Span {
	for i := range sp.Min {
		v := p.Val(i)
		if v < sp.Min[i] {
			sp.Min[i] = v
		}
		if v > sp.Max[i] {
			sp.Max[i] = v
		}
	}
	return sp
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
	sp := Span{e.Origin.Point, e.Origin.Point}
	sp.Expand(e.Twin.Origin)
	return sp
}
