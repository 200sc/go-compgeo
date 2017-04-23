package dcel

import "github.com/200sc/go-compgeo/geom"

// Bounds on a DCEL returns a Span calculated
// from every point in the DCEL.
func (dc *DCEL) Bounds() geom.Span {
	sp := geom.NewSpan()
	for _, v := range dc.Vertices {
		sp = sp.Expand(v)
	}
	return sp
}

// Bounds returns a Span calculated from
// every point on the Inner of this face
// because at time of writing we don't
// populate Outer
func (f *Face) Bounds() geom.Span {
	sp := geom.NewSpan()
	if f == nil {
		return sp
	}
	e := f.Outer
	sp = sp.Expand(e.Origin)
	for e.Next != f.Outer {
		e = e.Next
		sp = sp.Expand(e.Origin)
	}
	return sp
}

// Bounds on an Edge returns a Span
// on the edge's origin and it's twin's
// origin.
func (e *Edge) Bounds() geom.Span {
	sp := geom.NewSpan()
	sp = sp.Expand(e.Origin, e.Twin.Origin)
	return sp
}
