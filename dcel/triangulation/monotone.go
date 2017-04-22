package triangulation

import (
	"errors"

	"github.com/200sc/go-compgeo/dcel"
)

type helper struct {
	*dcel.Vertex
	typ int
}

// YMonotoneSplit converts a dcel into another dcel of its
// triangles, along with a mapping of faces in the new set
// to faces in the input set.
func YMonotoneSplit(dc *dcel.DCEL) (*dcel.DCEL, map[*dcel.Face]*dcel.Face, error) {

	dc = dc.Copy()

	// dc.Faces is modified through this algorithm,
	// so we need to iterate over a copy of it.
	faces := make([]*dcel.Face, len(dc.Faces))
	copy(faces, dc.Faces)

	for _, f := range faces {
		ypts := f.VerticesSorted(1, 0)
		helpers := make(map[*dcel.Edge]helper)
		for _, v := range ypts {
			switch MonotoneVertexType(v, dc) {
			case MONO_START:
				// Insert v's edge, map that edge to v
				e := CounterClockwiseEdge(v, dc)
				// insert e 
				helpers[e] = helper{v, MONO_START}
			case MONO_END:
				e := CounterClockwiseEdge(v, dc).Prev
				// If the previous edge's helper is a merge vertex,
				// we want to insert a diagonal in the dcel from v
				// to that helper.
				err := MergeInsert(helpers, e, v, dc)
				if err != nil {
					return nil, nil, err
				}
				delete(helpers, e)
			case MONO_REGULAR:
				if IsLeftOf(v, f, dc) {
					e := CounterClockwiseEdge(v, dc)
					prev := e.Prev
					err := MergeInsert(helpers, prev, v, dc)
					if err != nil {
						return nil, nil, err
					}
					delete(helpers, prev)
					// insert e
					helpers[e] = helper{v, MONO_REGULAR}
				} else {
					e := //edge left of v
					err := MergeInsert(helpers, e, v, dc)
					if err != nil {
						return nil, nil, err
					}
					helpers[e] = helper{v, MONO_REGULAR
				}
			case MONO_MERGE:
				e := CounterClockwiseEdge(v, dc).Prev
				err := MergeInsert(helpers, e, v, dc)
				if err != nil {
					return nil, nil, err
				}
				delete(helpers, e)
				e = // edge left of v
				err := MergeInsert(helpers, e, v, dc)
				if err != nil {
					return nil, nil, err
				}
				helpers[e] = v

			case MONO_SPLIT:
				e := //edge left of v
				err := MergeInsert(helpers, e, v, dc)
				if err != nil {
					return nil, nil, err
				}
				// insert counterClockwise(v)
				helpers[e] = helper{v, MONO_REGULAR

			}
		}
	}
	return nil, nil, nil
}

func CounterClockwiseEdge(v *dcel.Vertex, dc *dcel.DCEL) *dcel.Edge {
	if v.OutEdge.Face == dc.Faces[dcel.OUTER_FACE] {
		return v.OutEdge
	}
	return v.OutEdge.Twin.Prev
}

func MergeInsert(helpers map[*dcel.Edge]helper, e *dcel.Edge, v *dcel.Vertex, dc *dcel.DCEL) error {
	if help, ok := helpers[e]; ok {
		if help.typ == MONO_MERGE {
			return dc.ConnectVerts(help.Vertex, v)
		}
		return nil
	}
	return errors.New("Malformed helpers")
}

// IsLeftOf returns whether v is to the left or to the right
// of --the face which neighbors v--, f
func IsLeftOf(v *dcel.Vertex, f *dcel.Face, dc *dcel.DCEL) bool {
	// All internal faces have their points oriented counter-clockwise
	// We know that one of v's previous or next points is above,
	// and one is below, v.
	// If the previous is the one that is below v, then the face is
	// to the left.
	result := false
	if v.OutEdge.Prev.Origin.Y() < v.Y() {
		result = true
	}
	// If this is an external face our answer is reversed
	if f == dc.Faces[dcel.OUTER_FACE] {
		result = !result
	}
	return result
}
