package triangulation

import (
	"fmt"
	"math/rand"

	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
)

// TrapezoidalMap converts a dcel into a version of itself split into
// trapezoids and a search structure to find a containing trapezoid in
// the map in response to a point location query.
func TrapezoidalMap(dc *dcel.DCEL) (*dcel.DCEL, map[*dcel.Face]*dcel.Face, *TrapezoidNode, error) {
	bounds := dc.Bounds()

	Search := NewRoot()
	Search.set(left, NewTrapNode(newTrapezoid(bounds)))

	fullEdges, faces, err := dc.FullEdges()
	if err != nil {
		return nil, nil, nil, err
	}
	// Get rid of bad edges
	i := 0
	for i < len(fullEdges) {
		fe := fullEdges[i]
		l := fe.Left()
		r := fe.Right()
		if geom.F64eq(l.X(), r.X()) && geom.F64eq(l.Y(), r.Y()) {
			fullEdges = append(fullEdges[0:i], fullEdges[i+1:]...)
			faces = append(faces[0:i], faces[i+1:]...)
			i--
		}
		i++
	}
	fmt.Println("FullEdges")
	// Scramble the edges
	for i := range fullEdges {
		fmt.Println(fullEdges[i])
		j := i + rand.Intn(len(fullEdges)-i)
		fullEdges[i], fullEdges[j] = fullEdges[j], fullEdges[i]
	}
	for k, fe := range fullEdges {
		lowFace := faces[k][0]
		highFace := faces[k][1]
		// 1: Find the trapezoids intersected by fe
		trs := Search.Query(fe)
		// 2: Remove those and replace them with what they become
		//    due to the intersection of halfEdges[i]

		fmt.Println(Search)

		// Case A: A fe is contained in a single trapezoid tr
		// Then we make (up to) four trapezoids out of tr.
		if len(trs) == 0 {
			fmt.Println(fe, "intersected nothing?")
		}
		if len(trs) == 1 {
			fmt.Println(fe, "intersected one trapezoid", trs[0])
			mapSingleCase(trs[0], fe, lowFace, highFace)
		} else {
			fmt.Println(fe, "Intersected multiple zoids", trs)
			mapMultipleCase(trs, fe, lowFace, highFace)
		}
	}
	fmt.Println("Search:\n", Search)
	dc, m := Search.DCEL()
	return dc, m, Search, nil
}

func mapSingleCase(tr *Trapezoid, fe geom.FullEdge, lowFace, highFace *dcel.Face) {
	lp := fe.Left()
	rp := fe.Right()
	// Most pointers on the following trapezoids are
	// the same as the pointers on the trapezoid they
	// were split from.

	var l, r *Trapezoid
	ul := tr.Neighbors[upleft]
	bl := tr.Neighbors[botleft]
	ur := tr.Neighbors[upright]
	br := tr.Neighbors[botright]

	// Case 2A.2
	// If fe.left or fe.right lies ON tr's left and right
	// edges, we don't make new trapezoids for them.
	fmt.Println(lp, tr.left, rp, tr.right)
	if !geom.F64eq(lp.X(), tr.left) {
		l = tr.Copy()
		l.right = lp.X()
		tl := l.Edges[top].Left()
		tr := l.Edges[top].Right()
		tr = tr.Set(0, l.right).(geom.D3)
		l.Edges[top] = geom.NewFullEdge(tl, tr)
		bl := l.Edges[bot].Left()
		br := l.Edges[bot].Right()
		br = br.Set(0, l.right).(geom.D3)
		l.Edges[bot] = geom.NewFullEdge(bl, br)
	}
	if !geom.F64eq(rp.X(), tr.right) {
		r = tr.Copy()
		r.left = rp.X()
		tl := r.Edges[top].Left()
		tl = tl.Set(0, r.left).(geom.D3)
		tr := r.Edges[top].Right()
		r.Edges[top] = geom.NewFullEdge(tl, tr)
		bl := r.Edges[bot].Left()
		bl = bl.Set(0, r.left).(geom.D3)
		br := r.Edges[bot].Right()
		r.Edges[bot] = geom.NewFullEdge(bl, br)
	}

	u := tr.Copy()
	u.face = highFace
	d := tr.Copy()
	d.face = lowFace

	u.Neighbors[upleft] = ul
	u.Neighbors[botleft] = ul
	u.Neighbors[upright] = ur
	u.Neighbors[botright] = ur

	d.Neighbors[upleft] = bl
	d.Neighbors[botleft] = bl
	d.Neighbors[upright] = br
	d.Neighbors[botright] = br

	if l != nil {
		l.Neighbors[upright] = u
		l.Neighbors[botright] = d
		u.Neighbors[upleft] = l
		u.Neighbors[botleft] = l
		d.Neighbors[upleft] = l
		d.Neighbors[botleft] = l
	} else if ul != nil && !geom.F64eq(ul.Edges[bot].Right().Y(), lp.Y()) {
		// If these values are equal, our original assignments
		// above are fine.
		if ul.Edges[bot].Right().Y() > lp.Y() {
			u.Neighbors[botleft] = bl
		} else {
			d.Neighbors[upleft] = ul
		}
	} else if ul == nil && bl != nil {
		u.Neighbors[upleft] = bl
		u.Neighbors[botleft] = bl
	}
	if r != nil {
		r.Neighbors[upleft] = u
		r.Neighbors[botleft] = d
		u.Neighbors[upright] = r
		u.Neighbors[botright] = r
		d.Neighbors[upright] = r
		d.Neighbors[botright] = r
	} else if ur != nil && !geom.F64eq(ur.Edges[bot].Right().Y(), rp.Y()) {
		if ur.Edges[bot].Right().Y() > rp.Y() {
			u.Neighbors[botright] = br
		} else {
			d.Neighbors[upright] = ur
		}
	} else if ur == nil && br != nil {
		u.Neighbors[upright] = br
		u.Neighbors[botright] = br
	}

	d.Edges[top] = fe
	u.Edges[bot] = fe

	u.left = lp.X()
	d.left = lp.X()
	u.right = rp.X()
	d.right = rp.X()
	// 3: From the query structure, remove the leaves of the
	//    removed trapezoids and add new leaves for the new
	//    trapezoids, with additional inner nodes as necessary.

	a := NewX(lp)
	b := NewX(rp)
	c := NewY(fe)

	// Our structure should have tr's parent point to a,
	// a point to l and b, b point to r and c, and c
	// point to u and d

	tr.node.discard(a)

	if l != nil {
		a.set(left, NewTrapNode(l))
	}
	a.set(right, b)
	b.set(left, c)
	if r != nil {
		b.set(right, NewTrapNode(r))
	}
	c.set(left, NewTrapNode(u))
	c.set(right, NewTrapNode(d))
}

func mapMultipleCase(trs []*Trapezoid, fe geom.FullEdge, lowFace, highFace *dcel.Face) {

	lp := fe.Left()
	// Case B: fe is contained by more than one trapezoid
	// Step 1: if either fe.Left() or fe.Right() is not already
	// in the search structure, we define three new trapezoids by
	// drawing rays up and down from each new point.
	var ul, bl, l, u, b *Trapezoid
	var ln, un, bn *TrapezoidNode
	var x *TrapezoidNode
	y := NewY(fe)
	tr0 := trs[0]
	if !tr0.HasDefinedPoint(lp) {
		// The three trapezoids are split into
		// one to the left of an x node
		x = NewX(lp)
		// and two below the previous y node
		l = tr0.Copy()
		l.right = lp.X()

		ln = NewTrapNode(l)
		tr0.node.discard(x)
		x.left = ln
		x.right = y

		ul = l
		bl = l

	} else {
		// Otherwise we just split tr0 into two trapezoids.
		tr0.node.discard(y)
		ul = tr0.Neighbors[upleft]
		bl = tr0.Neighbors[botleft]
	}
	u = tr0.Copy()
	b = tr0.Copy()
	u.face = highFace
	b.face = lowFace
	u.left = lp.X()
	b.left = lp.X()

	u.Edges[bot], _ = fe.SubEdge(0, u.left, u.right)
	u.Neighbors[upleft] = ul
	u.Neighbors[botleft] = bl

	b.Edges[top], _ = fe.SubEdge(0, b.left, b.right)
	b.Neighbors[upleft] = ul
	b.Neighbors[botleft] = bl

	if l != nil {
		l.Neighbors[upright] = u
		l.Neighbors[botright] = b
	}

	// If fe.Left() is equal to the leftmost point
	// on the top or bottom edges of u or b respectively,
	// then we need to drop one set of left pointers
	if l != nil && lp == b.Edges[bot].Left() {
		b.Neighbors[botleft] = nil
		b.Neighbors[upleft] = nil
		l.Neighbors[botright] = u
	} else if l != nil && lp == u.Edges[top].Left() {
		u.Neighbors[botleft] = nil
		u.Neighbors[upleft] = nil
		l.Neighbors[upright] = b
	} else if l == nil && ul != nil && lp.X() == ul.Edges[bot].Right().X() {
		// If we added a point at the same horizontal value
		// of the edge connecting tr0 to its left neighbors,
		// but we weren't on an existing vertex, we need to
		// split up the definitions for ul and bl. Specifically,
		split := ul.Edges[bot]
		if lp.Y() > split.Right().Y() {
			b.Neighbors[upleft] = ul
		} else if lp.Y() < split.Right().Y() {
			u.Neighbors[botleft] = bl
		}
		// We also need to change our pointer setup.
		// The X split is still on the right value, but now needs
		// to point to a y node that points to ul and bl.
		x = NewX(lp)
		y.discard(x)
		x.set(right, y)
		y2 := NewY(split)
		x.set(left, y2)
		y2.set(left, ul.node)
		y2.set(right, bl.node)
	} else if l == nil && ul == nil && bl != nil {
		u.Neighbors[upleft] = bl
		b.Neighbors[upleft] = bl
	}

	un = NewTrapNode(u)
	bn = NewTrapNode(b)

	y.set(left, un)
	y.set(right, bn)
	// len(trs)-1 as the nth element is a special case,
	// just like the first, but it is initially handled
	// as if it is not.
	for i := 1; i < len(trs); i++ {
		tr := trs[i]
		// We are going to split this trapezoid into
		// an upper and lower trapezoid.
		// It is possible that one or both trapezoids
		// we make are mergeable into the previous upper
		// and lower trapezoids we made.

		y = NewY(fe)

		// If the upper edge of both this trapezoid
		// and the previous upper trapezoid, u, are
		// colinear, that is a mergable case.
		//
		// We use tr.Edges[top].Right() here instead of Left
		// as the two edges will often share a right and left
		// endpoint.
		p1 := u.Edges[top].Left()
		p2 := tr.Edges[top].Right()
		if geom.IsColinear(p1, u.Edges[top].Right(), p2) {
			u.Edges[top] = geom.NewFullEdge(p1, p2)
			u.Edges[bot] = geom.NewFullEdge(u.Edges[bot].Left(),
				tr.Edges[bot].Right())
		} else {
			u2 := tr.Copy()
			//
			u.Neighbors[botright] = u2
			u2.Neighbors[botleft] = u
			// Otherwise there are three reasons we might not be
			// able to merge.
			//
			// A: this trapezoid's upper edge
			// shares a vertex with the previous trapezoid, but
			// is at a different angle.
			u2tl := u2.Edges[top].Left()
			utr := u.Edges[top].Right()
			if u2tl == utr {
				// In this case, top left and bot left are both u.
				// u's bot right and bot left are similarly both u2.
				u.Neighbors[upright] = u2
				u2.Neighbors[upleft] = u
				// The top edges of u and u2 do not need to be updated.
				// The top edge of u2 is still accurate from the copy.
				// the search structure is updated later.
			} else {
				if u2tl.Y() > utr.Y() {
					// B: this trapezoid's left endpoint is above
					// the left endpoint of the previous trapezoid.
					u.Neighbors[upright] = u2
				} else {
					// C: this trapezoid's left endpoint is below
					// the left endpoint of the previous trapezoid.
					u2.Neighbors[upleft] = u
				}
			}
			u2.left = p1.X()
			u.right = p1.X()
			// the bottom edge is now this shard of fe.
			// if this trapezoid is merged with another,
			// this may change.
			u2.Edges[bot], _ = fe.SubEdge(0, u2.left, u2.right)
			// y points to a new trapezoid node holding u2
			un = NewTrapNode(u2)
			u = u2
		}

		y.set(left, un)

		p1 = b.Edges[bot].Left()
		p2 = tr.Edges[bot].Right()
		// Behavior is similar for the lower trapezoid
		if geom.IsColinear(p1, b.Edges[bot].Right(), p2) {
			b.Edges[bot] = geom.NewFullEdge(p1, p2)
			b.Edges[top] = geom.NewFullEdge(b.Edges[top].Left(),
				tr.Edges[top].Right())
		} else {
			b2 := tr.Copy()
			b.Neighbors[upright] = b2
			b2.Neighbors[upleft] = b
			b2bl := b2.Edges[bot].Left()
			bbr := b.Edges[bot].Right()
			if b2bl == bbr {
				b.Neighbors[botright] = b2
				b2.Neighbors[botleft] = b
			} else {
				if b2bl.Y() < bbr.Y() {
					b.Neighbors[botright] = b2
				} else {
					b2.Neighbors[botleft] = b
				}
			}
			b.right = p1.X()
			b2.left = p1.X()
			b2.Edges[top], _ = fe.SubEdge(0, b2.left, b2.right)
			bn = NewTrapNode(b2)
			b = b2
		}
		u.face = highFace
		b.face = lowFace
		y.set(right, bn)
	}
	// If fe.right is on some edge,
	// then we're done, except we need
	// to give u and b right edges, which
	// will be the same in the next case

	rp := fe.Right()
	u.right = rp.X()
	b.right = rp.X()

	var r, ur, br *Trapezoid
	trn := trs[len(trs)-1]
	if !trn.HasDefinedPoint(fe.Right()) {
		r = trn.Copy()
		r.left = rp.X()
		// other edges don't change
		//
		r.Neighbors[upleft] = u
		r.Neighbors[botleft] = b
		// r's other neighbors are unchanged.

		x = NewX(lp)
		// X needs to be put between y's
		// parents and x
		y.discard(x)
		x.set(left, y)

		rn := NewTrapNode(r)
		y.set(right, rn)
		ur = r
		br = r
	} else {
		trn.node.discard(y)
		ur = trn.Neighbors[upright]
		br = trn.Neighbors[botright]
	}

	u.Neighbors[botright] = br
	u.Neighbors[upright] = ur
	b.Neighbors[botright] = br
	b.Neighbors[upright] = ur
	if r != nil && rp == b.Edges[bot].Right() {
		u.Neighbors[botright] = nil
		u.Neighbors[upright] = nil
		r.Neighbors[upleft] = b
	} else if r != nil && rp == u.Edges[top].Right() {
		b.Neighbors[botright] = nil
		b.Neighbors[upright] = nil
		r.Neighbors[botleft] = u
	} else if r == nil && ur != nil && rp.X() == ur.Edges[bot].Left().X() {
		split := ur.Edges[bot]
		if rp.Y() > split.Right().Y() {
			b.Neighbors[upright] = ur
		} else if rp.Y() < split.Right().Y() {
			u.Neighbors[botright] = br
		}
		x = NewX(lp)
		y.discard(x)
		x.set(left, y)
		y2 := NewY(split)
		x.set(right, y2)
		y2.set(left, ur.node)
		y2.set(right, br.node)
	} else if r == nil && ur == nil {
		u.Neighbors[upright] = br
		b.Neighbors[upright] = br
	}
}
