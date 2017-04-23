package trapezoid

import (
	compgeo "github.com/200sc/go-compgeo"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
	"github.com/200sc/go-compgeo/printutil"
)

// A Node is a node in a tree structure
// for trapezoid map queries. It is structured so that
// each variety of Node is the same struct, but
// each has a different payload and query function.
type Node struct {
	left, right *Node
	parents     []*Node
	query       func(geom.FullEdge, *Node) []*Trapezoid
	payload     interface{}
}

// DCEL converts the trapezoids in the node search structure
// into a DCEL.
// This algorithm does not work?
// Consider scrapping this and replacing it with encoding the
// trapezoids as an OFF structure and decoding the resulting
// OFF structure as a DCEL.
func (tn *Node) DCEL() (*dcel.DCEL, map[*dcel.Face]*dcel.Face) {
	dc := new(dcel.DCEL)
	trs := tn.inOrder()
	// This maps from faces in the output of this algorithm
	// to faces in the input of the TrapezoidMap.
	fMap := make(map[*dcel.Face]*dcel.Face)
	vMap := make(map[geom.Point]*dcel.Edge)
	dc.Faces = make([]*dcel.Face, len(trs))
	// Todo: benchmark if it is faster to initially set this
	// at len(trs) because we know that's a minimum,
	// and then make if checks down the line for whether we append
	// or set.
	dc.Vertices = make([]*dcel.Vertex, 0)
	for i, tr := range trs {
		// Each trapezoid becomes a face
		dc.Faces[i] = dcel.NewFace()
		fMap[dc.Faces[i]] = tr.faces[0]
		// each of the up to four edges in a trapezoid
		// has an edge associated with it, the first of which is the face's
		// inner (going with the broken convention of always using inner)
		edges := tr.DCELEdges()
		dc.Faces[i].Outer = edges[0]
		dc.HalfEdges = append(dc.HalfEdges, edges...)
		// each vertex in each trapezoid, if it has not been seen before,
		// is added to the dcel and added to a map connected to it's edge
		// but if it is in the map,
		// it defines that the edge it is mapped to is the twin of the
		// current edge's previous, and vice versa.
		for _, e := range edges {
			e.Face = dc.Faces[i]
			if e2, ok := vMap[e.Origin.Point]; ok {
				e.SetTwin(e2.Prev)
				delete(vMap, e.Origin.Point)
			} else {
				vMap[e.Origin.Point] = e
				dc.Vertices = append(dc.Vertices, e.Origin)
			}
		}
	}
	//dc.CorrectTwins()
	return dc, fMap
}

func (tn *Node) inOrder() []*Trapezoid {
	if tn == nil {
		// error, unless this is root,
		// I think
		return []*Trapezoid{}
	}
	if tn.left == nil && tn.right == nil {
		// This is a trapezoid (or should be)
		return []*Trapezoid{tn.payload.(*Trapezoid)}
	}
	trs := tn.left.inOrder()
	return append(trs, tn.right.inOrder()...)
}

// PointLocate returns, from a given complex structure,
// which substructure that point falls into, if any.
// In the trapezoidal map, the query structure can
// be point-located on.
func (tn *Node) PointLocate(vs ...float64) (*dcel.Face, error) {
	if len(vs) < 2 {
		return nil, compgeo.InsufficientDimensionsError{}
	}
	// A point query on the structure is equivalent to an
	// edge query where both edges are the same.
	trs := tn.Query(geom.FullEdge{geom.Point{vs[0], vs[1], 0}, geom.Point{vs[0], vs[1], 0}})
	if len(trs) == 0 {
		return nil, nil
	}
	faces := trs[0].faces
	outerFace := tn.payload.(*dcel.Face)
	if faces[0] != outerFace && faces[0].Contains(geom.Point{vs[0], vs[1]}) {
		return faces[0], nil
	}
	if faces[1] != outerFace && faces[1].Contains(geom.Point{vs[0], vs[1]}) {
		return faces[1], nil
	}
	return nil, nil
}

// Query is shorthand for tn.query(fe, tn)
func (tn *Node) Query(fe geom.FullEdge) []*Trapezoid {
	if tn == nil {
		return []*Trapezoid{}
	}
	return tn.query(fe, tn)
}

func (tn *Node) discard(n *Node) {
	for _, p := range tn.parents {
		if p.left == tn {
			p.left = n
		} else {
			p.right = n
		}
	}
	n.parents = tn.parents
}

func (tn *Node) set(v int, n *Node) {
	switch v {
	case left:
		tn.left = n
	case right:
		tn.right = n
	}
	n.parents = append(n.parents, tn)
}

func (tn *Node) String() string {
	return tn.string("", true)
}

// Todo: can this string structure be generalized?
func (tn *Node) string(prefix string, isTail bool) string {
	if tn == nil || len(prefix) > 64 {
		return ""
	}
	s := prefix
	if isTail {
		s += "└──"
		prefix += "    "
	} else {
		s += "├──"
		prefix += "│   "
	}
	s += tn.payloadString() + "\n"
	s += tn.right.string(prefix, false)
	s += tn.left.string(prefix, true)

	return s
}

func (tn *Node) payloadString() string {
	switch v := tn.payload.(type) {
	case *Trapezoid:
		return "T" + v.String()
	case geom.D3:
		return "X" + "(" + printutil.Stringf64(v.Val(0)) + ")"
	case geom.FullEdge:
		return "Y" + "(" + printutil.Stringf64(v[0][0], v[0][1], v[1][0], v[1][1]) + ")"
	}
	return "Root"
}

// NewRoot returns a root node.
// There is only one root node. We use a root node to avoid
// having "if parent == nil, set this to root" checks in
// our map code.
func NewRoot() *Node {
	return &Node{
		query: rootQuery,
	}
}

func rootQuery(fe geom.FullEdge, n *Node) []*Trapezoid {
	return n.left.Query(fe)
}
