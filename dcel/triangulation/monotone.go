package triangulation

import (
	"errors"
	"fmt"

	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/visualize"
	"github.com/200sc/go-compgeo/geom"
	"github.com/200sc/go-compgeo/search"
	"github.com/200sc/go-compgeo/search/tree"
)

type helper struct {
	*dcel.Vertex
	typ int
}

func YMonotoneTriangulate(inDc *dcel.DCEL) (*dcel.DCEL, map[*dcel.Face]*dcel.Face, error) {
	monotonized, faceMap, err := YMonotoneSplit(inDc)
	if err != nil {
		return monotonized, faceMap, err
	}
	// ...
	// Triangulate each monotone polygon
	return monotonized, faceMap, nil
}

// YMonotoneSplit converts a dcel into another dcel of y
// monotone shapes, along with a mapping of faces in the new set
// to faces in the input set.
func YMonotoneSplit(inDc *dcel.DCEL) (*dcel.DCEL, map[*dcel.Face]*dcel.Face, error) {

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

	for i := dcel.OUTER_FACE + 1; i < faceLen; i++ {
		f := dc.Faces[i]
		ypts := f.VerticesSorted(1, 0)
		helpers := make(map[*dcel.Edge]helper)
		for _, v := range ypts {
			switch MonotoneVertexType(v, dc) {
			case MONO_START:
				// Insert v's edge, map that edge to v
				e := CounterClockwiseEdge(v, dc)
				edgeTree.Insert(edgeNode{e})
				helpers[e] = helper{v, MONO_START}
			case MONO_END:
				e := CounterClockwiseEdge(v, dc).Prev
				// If the previous edge's helper is a merge vertex,
				// we want to insert a diagonal in the dcel from v
				// to that helper.
				err := MergeInsert(helpers, faceMap, f, e, v, dc)
				if err != nil {
					return nil, nil, err
				}
				delete(helpers, e)
			case MONO_REGULAR:
				if IsLeftOf(v, f, dc) {
					e := CounterClockwiseEdge(v, dc)
					prev := e.Prev
					err := MergeInsert(helpers, faceMap, f, prev, v, dc)
					if err != nil {
						return nil, nil, err
					}
					delete(helpers, prev)
					edgeTree.Insert(edgeNode{e})
					helpers[e] = helper{v, MONO_REGULAR}
				} else {
					e2 := CounterClockwiseEdge(v, dc)
					c, _ := edgeTree.SearchDown(compEdge{e2}, 1)
					e := c.(compEdge).Edge
					err = MergeInsert(helpers, faceMap, f, e, v, dc)
					if err != nil {
						return nil, nil, err
					}
					helpers[e] = helper{v, MONO_REGULAR}
				}
			case MONO_MERGE:
				e2 := CounterClockwiseEdge(v, dc)
				e := e2.Prev
				err := MergeInsert(helpers, faceMap, f, e, v, dc)
				if err != nil {
					return nil, nil, err
				}
				delete(helpers, e)
				c, _ := edgeTree.SearchDown(compEdge{e2}, 1)
				e = c.(compEdge).Edge
				err = MergeInsert(helpers, faceMap, f, e, v, dc)
				if err != nil {
					return nil, nil, err
				}
				helpers[e] = helper{v, MONO_MERGE}

			case MONO_SPLIT:
				e2 := CounterClockwiseEdge(v, dc)
				c, _ := edgeTree.SearchDown(compEdge{e2}, 1)
				e := c.(compEdge).Edge
				err := MergeInsert(helpers, faceMap, f, e, v, dc)
				if err != nil {
					return nil, nil, err
				}
				edgeTree.Insert(edgeNode{e2})
				helpers[e] = helper{v, MONO_SPLIT}
				helpers[e2] = helper{v, MONO_SPLIT}
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

func MergeInsert(helpers map[*dcel.Edge]helper, faces map[*dcel.Face]*dcel.Face,
	f *dcel.Face, e *dcel.Edge, v *dcel.Vertex, dc *dcel.DCEL) error {
	if help, ok := helpers[e]; ok {
		if help.typ == MONO_MERGE {
			err := dc.ConnectVerts(help.Vertex, v)
			// Add to face map
			if err == nil {
				newFace := dc.Faces[len(dc.Faces)-1]
				// Whatever the old face used to point to in the original
				// dcel, the new face also does, as it was split off of that face.
				faces[newFace] = faces[f]
			}
			return err
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

type edgeNode struct {
	v *dcel.Edge
}

func (en edgeNode) Key() search.Comparable {
	return compEdge{en.v}
}

func (en edgeNode) Val() search.Equalable {
	return valEdge{en.v}
}

type valEdge struct {
	*dcel.Edge
}

func (ve valEdge) Equals(e search.Equalable) bool {
	switch ve2 := e.(type) {
	case valEdge:
		return ve.Edge == ve2.Edge
	}
	return false
}

// We need to have our keys be CompEdges so
// they are comparable within a certain y range.
type compEdge struct {
	*dcel.Edge
}

func (ce compEdge) Compare(i interface{}) search.CompareResult {
	switch c := i.(type) {
	case compEdge:
		if visualize.VisualCh != nil {
			visualize.DrawLine(ce.Edge.Origin, ce.Edge.Twin.Origin)
			visualize.DrawLine(c.Edge.Origin, c.Edge.Twin.Origin)
		}
		fmt.Println("Comparing", ce, c)
		if ce.Edge == c.Edge {
			fmt.Println("Equal1!")
			return search.Equal
		}

		if geom.F64eq(ce.X(), c.X()) && geom.F64eq(ce.Y(), c.Y()) &&
			geom.F64eq(ce.Twin.X(), c.Twin.X()) && geom.F64eq(ce.Twin.Y(), c.Twin.Y()) {
			fmt.Println("Equal2!")
			return search.Equal
		}
		y, err := ce.FindSharedPoint(c.Edge, 1)
		if err != nil {
			fmt.Println("Edges share no y point")
		}
		p1, _ := ce.PointAt(1, y)
		p2, _ := c.PointAt(1, y)
		if p1[0] < p2[0] {
			fmt.Println("Less!")
			return search.Less
		}
		fmt.Println("Greater!")
		return search.Greater
	}
	return ce.Edge.Compare(i)
}
