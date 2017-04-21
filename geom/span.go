package geom

import "math"

const (
	// SPAN_MIN refers to the point in a span
	// which it holds its minimum value in.
	SPAN_MIN = iota
	// SPAN_MAX returns to the index in a span
	// which it holds its maximum value in.
	SPAN_MAX
)

// A Span represents n-dimensions of
// span from one point to another
// for however many dimensions a point
// has
type Span struct {
	FullEdge
}

// NewSpan returns a span with its values
// set to appropriate infinities
func NewSpan(ds ...Dimensional) Span {
	sp := Span{}
	for i := 0; i < sp.D(); i++ {
		sp.At(SPAN_MIN).Set(i, math.MaxFloat64)
		sp.At(SPAN_MAX).Set(i, math.MaxFloat64*-1)
	}
	return sp.Expand(ds...)
}

// Low returns SPAN_MIN
func (sp Span) Low(i int) Dimensional {
	return sp.At(SPAN_MIN)
}

// High returns SPAN_MAX
func (sp Span) High(i int) Dimensional {
	return sp.At(SPAN_MAX)
}

// Eq returns whether Span is equivalent to
// the given spanning type
func (sp Span) Eq(s Spanning) bool {
	if sp.D() != s.D() {
		return false
	}
	for i := 0; i < sp.D(); i++ {
		if !sp.At(i).Eq(s.At(i)) {
			return false
		}
	}
	return true
}

// Set , which you should probably not use on a Span,
// sets the value of one of span's min or max to something
// which will no longer necessarily be the min or max of
// the points the span has been exposed to. Consider using
// a FullEdge.
func (sp Span) Set(i int, d Dimensional) Spanning {
	sp.FullEdge[i] = d.(Point)
	return sp
}

// Expand on a Span will reduce or increase
// a span's min and max values if the input points
// falls outside of the span on any dimension
func (sp Span) Expand(ps ...Dimensional) Span {
	for _, p := range ps {
		j := p.D()
		for i := 0; i < j; i++ {
			v := p.Val(i)
			if v < sp.At(SPAN_MIN).Val(i) {
				sp.FullEdge[0][i] = v
			}
			if v > sp.At(SPAN_MAX).Val(i) {
				sp.FullEdge[1][i] = v
			}
		}
	}
	return sp
}

// Lesser returns Span Min
func (sp Span) Lesser(d int) Dimensional {
	return sp.At(SPAN_MIN)
}

// Greater returns Span Max
func (sp Span) Greater(d int) Dimensional {
	return sp.At(SPAN_MAX)
}

// Left returns Span Min
func (sp Span) Left() D3 {
	return sp.Lesser(0).(D3)
}

// Right returns Span Max
func (sp Span) Right() D3 {
	return sp.Greater(0).(D3)
}

// Top returns Span Max
func (sp Span) Top() D3 {
	return sp.Greater(1).(D3)
}

// Bottom returns Span Min
func (sp Span) Bottom() D3 {
	return sp.Lesser(1).(D3)
}

// Outer returns Span Min
func (sp Span) Outer() D3 {
	return sp.Lesser(2).(D3)
}

// Inner returns Span Max
func (sp Span) Inner() D3 {
	return sp.Greater(2).(D3)
}
