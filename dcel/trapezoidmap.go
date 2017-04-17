package dcel

import "math/rand"

func (dc *DCEL) TrapezoidalMap() (*DCEL, error) {
	bounds := dc.Bounds()
	// Scramble the edges of a new DCEL
	Traps := []*Trapezoid{}
	Traps = append(Traps, bounds.ToTrapezoid())

	Search := Traps[0]

	dc = dc.Copy()
	fullEdges, err := dc.FullEdges()
	if err != nil {
		return nil, err
	}

	// Scramble the edges
	for i := range fullEdges {
		j := i + rand.Intn(len(fullEdges)-i)
		fullEdges[i], fullEdges[j] = fullEdges[j], fullEdges[i]
	}
	for i, fe := range fullEdges {
		// 1: Find the trapezoids intersected by fe
		trs := Search.Query(fe)
		// 2: Remove those and replace them with what they become
		//    due to the intersection of halfEdges[i]

		// Case 2A: A fe is contained in a single trapezoid tr
		// Then we make four trapezoids out of tr.
		if len(trs) == 1 {
			tr := trs[0]
			lp := fe.Left()
			rp := fe.Right()
			// Most pointers on the following trapezoids are
			// the same as the pointers on the trapezoid they
			// were split from.
			
			l := tr.Neighbors[left]
			r := tr.Neighbors[right]

			// If fe.left or fe.right lies ON tr's left and right
			// edges, we don't make new trapezoids for them.
			if !lp.LiesOn(tr.Edges[left]) {
				l = tr.Copy()
				l.Edges[right] = FullEdge{
					tr.Edges[top].PointAt(0, lp.X()),
					tr.Edges[bot].PointAt(0, lp.X()),
				}
			}
			if !rp.LiesOn(tr.Edges[right]) {
				r = tr.Copy()
				r.Edges[left] = FullEdge {
					tr.Edges[top].PointAt(0, rp.X()),
					tr.Edges[bot].PoinAt(0, rp.X()),
				}
			}

			u := tr.Copy()
			d := tr.Copy()

			l.Neighbors[upright] = u
			l.Neighbors[botright] = d

			r.Neighbors[upleft] = u
			r.Neighbors[botleft] = d

			u.Neighbors[upleft] = l
			u.Neighbors[botleft] = l
			u.Neighbors[upright] = r
			u.Neighbors[botright] = r

			d.Neighbors[upleft] = l
			d.Neighbors[botleft] = l
			d.Neighbors[upright] = r
			d.Neighbors[botright] = r

			d.Edges[top] = fe
			u.Edges[bot] = fe

			u.Edges[left] = FullEdge {
				lp, tr.Edges[top].PointAt(0, lp.X())
			}
			b.Edges[left] = FullEdge {
				lp, tr.Edges[bot].PoinAt(0, lp.X())
			}
			u.Edges[left] = FullEdge {
				rp, tr.Edges[top].PointAt(0, rp.X())
			}
			b.Edges[left] = FullEdge {
				rp, tr.Edges[bot].PoinAt(0, rp.X())
			}
		}
		// 3: From the query structure, remove the leaves of the
		//    removed trapezoids and add new leaves for the new
		//    trapezoids, with additional inner nodes as necessary.
	}
	return dc, nil
}

const (
	top = iota
	bot
	left
	right
)

const (
	upright = iota
	botright
	upleft
	botleft
)

type TrapezoidNode interface {
	Query(FullEdge) []*Trapezoid
}

type Trapezoid struct {
	// See above indices
	Edges     [4]FullEdge
	Neighbors [4]*Trapezoid
}

func (tr *Trapezoid) Copy() *Trapezoid {
	tr2 := new(Trapezoid)
	tr2.Edges = tr.Edges
	tr2.Neighbors = tr.Neighbors
	return tr2
}

func (tr *Trapezoid) Query(fe FullEdge) []*Trapezoid {
	traps := []*Trapezoid{tr}
	r := fe.Right()
	for tr != nil && tr.Edges[right] != nil &&
		r.IsRightOf(tr.Edges[right]) {

		// Case 1: tr.Upright = tr.UpLeft
		// Caught by the next case, but it doesn't
		// matter which we add.
		if tr.Neighbors[upright].Edges[bot].IsAbove(fe) {
			tr = tr.Neighbors[botright]
		} else {
			tr = tr.Neighbors[botleft]
		}

		traps = append(traps, t)
	}
	return traps
}

type XNode struct {
	*Vertex
	left,right,parent TrapezoidNode
}

func (xn *XNode) Query(fe FullEdge) []*Trapezoid {
	//...
}

type YNode struct {
	FullEdge
	left,right,parent Trapezoid
}

func (yn *YNode) Query(fe FullEdge) []*Trapezoid {
	//...
}
