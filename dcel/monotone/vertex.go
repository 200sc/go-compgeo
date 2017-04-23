package monotone

import (
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
)

// VertexType Const
const (
	START = iota
	END
	REGULAR
	SPLIT
	MERGE
)

func VertexType(v *dcel.Vertex, dc *dcel.DCEL) int {
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
				return START
			}
			return SPLIT
		}
	} else if geom.Greater2D(n, v) == n {
		if cp < 0 {
			return END
		}
		return MERGE
	}
	return REGULAR

}
