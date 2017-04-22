package triangulation

import (
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
)

// MonotoneVertexType Const
const (
	MONO_START = iota
	MONO_END
	MONO_REGULAR
	MONO_SPLIT
	MONO_MERGE
)

func MonotoneVertexType(v *dcel.Vertex, dc *dcel.DCEL) int {
	p := v.OutEdge.Prev.Origin
	n := v.OutEdge.Next.Origin

	cp := geom.Cross2D(p, v, n)
	if v.OutEdge.Face == dc.Faces[dcel.OUTER_FACE] {
		cp *= -1
	}
	if geom.Lesser2D(p, v) == p {
		if geom.Lesser2D(n, v) == n {
			// This might need to be flipped.
			if cp < 0 {
				return MONO_START
			}
			return MONO_SPLIT
		}
	} else if geom.Greater2D(n, v) == n {
		if cp < 0 {
			return MONO_END
		}
		return MONO_MERGE
	}
	return MONO_REGULAR

}
