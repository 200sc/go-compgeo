// monotone holds functions for converting dcels into monotone components

package monotone

import (
	"errors"

	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/search/tree"
)

type helper struct {
	*dcel.Vertex
	typ int
}

// Split converts a dcel into another dcel of y
// monotone shapes, along with a mapping of faces in the new set
// to faces in the input set.
func Split(inDc *dcel.DCEL) (*dcel.DCEL, map[*dcel.Face]*dcel.Face, error) {

	dc := inDc.Copy()

	faceMap := make(map[*dcel.Face]*dcel.Face)
	// Initialize the map to map 1-to-1
	for i, f := range dc.Faces {
		faceMap[f] = inDc.Faces[i]
	}

	// dc.Faces is modified through this algorithm,
	// so we need to iterate it's current length (ignoring OUTER_FACE)
	faceLen := len(dc.Faces)
	edgeTree := tree.New(tree.RedBlack)
	edgeLen := len(dc.HalfEdges)

	for i := dcel.OUTER_FACE + 1; i < faceLen; i++ {
		f := dc.Faces[i]
		ypts := f.VerticesSorted(1, 0)
		helpers := make(map[*dcel.Edge]helper)
		for _, v := range ypts {
			switch VertexType(v, dc) {
			case START:
				// Insert v's edge, map that edge to v
				e := CounterClockwiseEdge(v, dc)
				edgeTree.Insert(edgeNode{e})
				helpers[e] = helper{v, START}
			case END:
				e := CounterClockwiseEdge(v, dc).Prev
				// If the previous edge's helper is a merge vertex,
				// we want to insert a diagonal in the dcel from v
				// to that helper.
				err := MergeInsert(helpers, f, e, v, dc)
				if err != nil {
					return nil, nil, err
				}
				delete(helpers, e)
			case REGULAR:
				if IsLeftOf(v, f, dc) {
					e := CounterClockwiseEdge(v, dc)
					prev := e.Prev
					err := MergeInsert(helpers, f, prev, v, dc)
					if err != nil {
						return nil, nil, err
					}
					delete(helpers, prev)
					edgeTree.Insert(edgeNode{e})
					helpers[e] = helper{v, REGULAR}
				} else {
					e2 := CounterClockwiseEdge(v, dc)
					c, _ := edgeTree.SearchDown(compEdge{e2}, 1)
					e := c.(compEdge).Edge
					err := MergeInsert(helpers, f, e, v, dc)
					if err != nil {
						return nil, nil, err
					}
					helpers[e] = helper{v, REGULAR}
				}
			case MERGE:
				e2 := CounterClockwiseEdge(v, dc)
				e := e2.Prev
				err := MergeInsert(helpers, f, e, v, dc)
				if err != nil {
					return nil, nil, err
				}
				delete(helpers, e)
				c, _ := edgeTree.SearchDown(compEdge{e2}, 1)
				e = c.(compEdge).Edge
				err = MergeInsert(helpers, f, e, v, dc)
				if err != nil {
					return nil, nil, err
				}
				helpers[e] = helper{v, MERGE}

			case SPLIT:
				e2 := CounterClockwiseEdge(v, dc)
				c, _ := edgeTree.SearchDown(compEdge{e2}, 1)
				e := c.(compEdge).Edge
				err := MergeInsert(helpers, f, e, v, dc)
				if err != nil {
					return nil, nil, err
				}
				edgeTree.Insert(edgeNode{e2})
				helpers[e] = helper{v, SPLIT}
				helpers[e2] = helper{v, SPLIT}
			}
		}
		// Walk each new edge to create new faces.
		for i := edgeLen; i < len(dc.HalfEdges); i++ {
			e := dc.HalfEdges[i]
			if e.Face == f {
				newFace := dcel.NewFace()
				newFace.Outer = e
				e.Face = newFace
				for e = e.Next; e != newFace.Outer; e = e.Next {
					e.Face = newFace
				}
				faceMap[newFace] = f
			} // else we've already walked this edge
		}
		edgeLen = len(dc.HalfEdges)
	}
	return nil, nil, nil
}

func CounterClockwiseEdge(v *dcel.Vertex, dc *dcel.DCEL) *dcel.Edge {
	if v.OutEdge.Face == dc.Faces[dcel.OUTER_FACE] {
		return v.OutEdge
	}
	return v.OutEdge.Twin.Prev
}

func MergeInsert(helpers map[*dcel.Edge]helper, f *dcel.Face, e *dcel.Edge,
	v *dcel.Vertex, dc *dcel.DCEL) error {
	if help, ok := helpers[e]; ok {
		if help.typ == MERGE {
			dc.ConnectVerts(help.Vertex, v, f)
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
