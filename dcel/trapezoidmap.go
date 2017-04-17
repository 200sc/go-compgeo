package dcel

// func (dc *DCEL) TrapezoidalMap() *DCEL {
// 	bounds := dc.Bounds()
// 	// Scramble the edges of a new DCEL
// 	Traps := []*Trapezoid{}
// 	Traps = append(Traps, bounds.ToTrapezoid())

// 	Search := &Node{
// 		nil,
// 		Traps[0],
// 		0,
// 		Point{0, 0, 0},
// 		nil, nil, nil,
// }

// 	dc = dc.Copy()
// 	for i := 0; i < len(dc.HalfEdges); i++ {
// 		j := i + rand.Intn(len(dc.HalfEdges)-i)
// 		dc.HalfEdges[i], dc.HalfEdges[j] = dc.HalfEdges[j], dc.HalfEdges[i]
// 	}
// 	for i := 0; i < len(dc.HalfEdges); i++ {
// 		// 1: Find the trapezoids intersected by halfEdges[i]
// 		// 2: Remove those and replace them with what they become
// 		//    due to the intersection of halfEdges[i]
// 		// 3: From the query structure, remove the leaves of the
// 		//    removed trapezoids and add new leaves for the new
// 		//    trapezoids, with additional inner nodes as necessary.
// 	}
// 	return dc
// }

type Trapezoid struct {
	Left, Right          *Edge
	MinY, MaxY           float64
	Valid                bool
	u0, u1, uSave, uSide int
	d0, d1               int
	sink                 *Node
}

type NodeType int

// NodeType const
const (
	T_X NodeType = iota
	T_Y
	T_SINK
)

type Node struct {
	e                   *Edge
	t                   *Trapezoid
	typ                 NodeType
	yval                Point
	parent, left, right *Node
}

type TrapEdge struct {
	*Edge
	r0, r1 *Node
}
