package dcel

// func (dc *DCEL) TrapezoidalMap() (*DCEL, error) {
// 	bounds := dc.Bounds()
// 	// Scramble the edges of a new DCEL
// 	Traps := []*Trapezoid{}
// 	Traps = append(Traps, bounds.Trapezoid())

// 	Search := NewRoot()
// 	Search.left = NewTrapNode(Traps[0])

// 	dc = dc.Copy()
// 	fullEdges, err := dc.FullEdges()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Scramble the edges
// 	for i := range fullEdges {
// 		j := i + rand.Intn(len(fullEdges)-i)
// 		fullEdges[i], fullEdges[j] = fullEdges[j], fullEdges[i]
// 	}
// 	for _, fe := range fullEdges {
// 		// 1: Find the trapezoids intersected by fe
// 		trs := Search.Query(fe)
// 		// 2: Remove those and replace them with what they become
// 		//    due to the intersection of halfEdges[i]

// 		// Case A: A fe is contained in a single trapezoid tr
// 		// Then we make (up to) four trapezoids out of tr.
// 		if len(trs) == 1 {
// 			mapSingleCase(trs[0], fe)
// 		} else {
// 			lp := fe.Left()
// 			rp := fe.Right()
// 			// Case B: fe is contained by more than one trapezoid
// 			// Step 1: if either fe.Left() or fe.Right() is not already
// 			// in the search structure, we define three new trapezoids by
// 			// drawing rays up and down from each new point.
// 			var l, u, b, r *Trapezoid
// 			var ln, un, bn, rn *TrapezoidNode
// 			var y *TrapezoidNode
// 			tr0 := trs[0]
// 			if !tr0.HasDefinedPoint(lp) {
// 				x := NewX(lp)
// 				y = NewY(fe)
// 				l = tr0.Copy()
// 				p1, _ := tr0.Edges[top].PointAt(0, lp.X())
// 				p2, _ := tr0.Edges[bot].PointAt(0, lp.X())
// 				l.Edges[right] = FullEdge{*p1, *p2}
// 				ln = NewTrapNode(l)
// 				tr0.Discard(x)
// 				x.left = ln
// 				x.right = y

// 			}
// 			trn := trs[len(trs)-1]
// 			if !trn.HasDefinedPoint(fe.Right()) {

// 			}
// 		}
// 	}
// 	return dc, nil
// }

func mapSingleCase(tr *Trapezoid, fe FullEdge) {
	lp := fe.Left()
	rp := fe.Right()
	// Most pointers on the following trapezoids are
	// the same as the pointers on the trapezoid they
	// were split from.

	l := tr.Neighbors[left]
	r := tr.Neighbors[right]

	// Case 2A.2
	// If fe.left or fe.right lies ON tr's left and right
	// edges, we don't make new trapezoids for them.
	if !IsColinear(lp, tr.Edges[left].Left(), tr.Edges[left].Right()) {
		l = tr.Copy()
		p1, _ := tr.Edges[top].PointAt(0, lp.X())
		p2, _ := tr.Edges[bot].PointAt(0, lp.X())
		l.Edges[right] = FullEdge{*p1, *p2}
	}
	if !IsColinear(rp, tr.Edges[right].Left(), tr.Edges[right].Right()) {
		r = tr.Copy()
		p1, _ := tr.Edges[top].PointAt(0, rp.X())
		p2, _ := tr.Edges[bot].PointAt(0, rp.X())
		r.Edges[left] = FullEdge{*p1, *p2}
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

	p, _ := tr.Edges[top].PointAt(0, lp.X())
	u.Edges[left] = FullEdge{lp, *p}
	p, _ = tr.Edges[bot].PointAt(0, lp.X())
	d.Edges[left] = FullEdge{lp, *p}
	p, _ = tr.Edges[top].PointAt(0, rp.X())
	u.Edges[left] = FullEdge{rp, *p}
	p, _ = tr.Edges[bot].PointAt(0, rp.X())
	d.Edges[left] = FullEdge{rp, *p}
	// 3: From the query structure, remove the leaves of the
	//    removed trapezoids and add new leaves for the new
	//    trapezoids, with additional inner nodes as necessary.

	a := NewX(lp)
	b := NewX(rp)
	c := NewY(fe)

	// Our structure should have tr's parent point to a,
	// a point to l and b, b point to r and c, and c
	// point to u and d

	tr.Discard(a)

	a.left = NewTrapNode(l)
	l.parents = []*TrapezoidNode{a}
	a.right = b
	b.parents = []*TrapezoidNode{a}
	b.left = c
	c.parents = []*TrapezoidNode{b}
	b.right = NewTrapNode(r)
	r.parents = []*TrapezoidNode{b}
	c.left = NewTrapNode(u)
	u.parents = []*TrapezoidNode{c}
	c.right = NewTrapNode(d)
	d.parents = []*TrapezoidNode{c}
}
