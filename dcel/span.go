package dcel

import "math"

type Span struct {
	Min Point
	Max Point
}

func NewSpan() Span {
	sp := Span{}
	for i := range sp.Min {
		sp.Min[i] = math.MaxFloat64
		sp.Max[i] = math.MaxFloat64 * -1
	}
	return sp
}

func (sp Span) Expand(p *Point) Span {
	for i := range sp.Min {
		if p[i] < sp.Min[i] {
			sp.Min[i] = p[i]
		}
		if p[i] > sp.Max[i] {
			sp.Max[i] = p[i]
		}
	}
	return sp
}
