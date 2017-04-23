package kirkpatrick

import (
	compgeo "github.com/200sc/go-compgeo"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/pointLoc"
	"github.com/200sc/go-compgeo/dcel/pointLoc/monotone"
	"github.com/200sc/go-compgeo/dcel/pointLoc/trapezoid"
	"github.com/200sc/go-compgeo/geom"
)

//Triangulation method constant
type Method int

const (
	MONOTONE Method = iota
	TRAPEZOID
)

func TriangleTree(dc *dcel.DCEL, m Method) (pointLoc.LocatesPoints, error) {
	var tri *dcel.DCEL
	var mp map[*dcel.Face]*dcel.Face
	var err error

	// We'll remove these vertices later
	//verts := dc.Vertices
	switch m {
	case MONOTONE:
		// We need to wrap this dcel in some bounding polygon.
		dc2 := dc.Copy()
		// We add the bounds of dc2, expanded by 1, to create
		// a large face out of each edge whose existing face was
		// the outer face.
		bounds := dc2.Bounds()
		p1 := bounds.At(geom.SPAN_MIN)
		p1 = p1.Set(0, p1.Val(0)-1)
		p1 = p1.Set(1, p1.Val(1)-1)
		p3 := bounds.At(geom.SPAN_MAX)
		p3 = p3.Set(0, p3.Val(0)+1)
		p3 = p3.Set(1, p3.Val(1)+1)
		p2 := geom.NewPoint(p1.Val(0), p3.Val(1), 0)
		p4 := geom.NewPoint(p3.Val(0), p1.Val(1), 0)

		boundDc := dcel.FourPoint(p1.(geom.D3), p2, p3.(geom.D3), p4)

		// Correct face pointers
		for _, e := range dc2.HalfEdges {
			if e.Face == dc2.Faces[dcel.OUTER_FACE] {
				e.Face = boundDc.Faces[1]
			}
		}

		// Find a point on dc2 to connect to
		minV := dc2.Vertices[0]
		for _, v := range dc2.Vertices {
			if v.X() < minV.X() {
				minV = v
			}
		}

		// Combine boundDc into dc2
		dc2.Faces[0] = boundDc.Faces[0]
		dc2.Faces = append(dc2.Faces, boundDc.Faces[1])
		dc2.Vertices = append(dc2.Vertices, boundDc.Vertices...)
		dc2.HalfEdges = append(dc2.HalfEdges, boundDc.HalfEdges...)
		v1 := dc2.Vertices[len(dc2.Vertices)-4] // Span_min
		// add an edge from boundDc to dc2's outer edge
		dc2.ConnectVerts(v1, minV, dc2.Faces[len(dc2.Faces)-1])

		tri, mp, err = monotone.Triangulate(dc2)
	case TRAPEZOID:
		// The trapezoidal map method requires that we add to our
		// dcel a wrapping square, so we already have our outer polygon.
		tri, mp, _, err = trapezoid.TrapezoidalMap(dc)
		if err != nil {
			return nil, err
		}
		tri, mp, err = monotone.TriangulateSplit(tri, mp)
	}
	if err != nil {
		return nil, err
	}

	// At this point we have a base level of triangles
	nodes := make([]*Tree, len(tri.Faces))
	for i, t := range tri.Faces {
		nodes[i] = &Tree{tri: t}
	}

	// Hit a wall here.
	// We need a structure that stores every vertex and every face that
	// each vertex borders, so when we remove a vertex we can somehow
	// indicate which facecs are being removed and what faces are being added.
	// this structure needs to somehow inform other vertices when a face they
	// are next to has been removed, or it needs to be regenerated at every
	// tier of the triangle stack.
	return nil, compgeo.UnsupportedError{}
}

type Tree struct {
	tri      *dcel.Face
	children []*Tree
}
