package trapezoid

import (
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
)

func mapMultipleCase(trs []*Trapezoid, fe geom.FullEdge, faces [2]*dcel.Face) {

	lp := fe.Left()
	// Case B: fe is contained by more than one trapezoid
	// Step 1: if either fe.Left() or fe.Right() is not already
	// in the search structure, we define three new trapezoids by
	// drawing rays up and down from each new point.
	var ul, bl, l, u, b *Trapezoid
	var ln, un, bn *Node
	var x *Node
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
	u.faces = faces
	b.faces = faces
	u.left = lp.X()
	b.left = lp.X()

	edge, _ := fe.SubEdge(0, u.left, u.right)
	u.bot[left] = edge.Left().Y()
	u.bot[right] = edge.Right().Y()
	u.Neighbors[upleft] = ul
	u.Neighbors[botleft] = bl

	edge, _ = fe.SubEdge(0, b.left, b.right)
	b.top[left] = edge.Left().Y()
	b.top[right] = edge.Right().Y()
	b.Neighbors[upleft] = ul
	b.Neighbors[botleft] = bl

	if l != nil {
		l.Neighbors[upright] = u
		l.Neighbors[botright] = b
	}

	// If fe.Left() is equal to the leftmost point
	// on the top or bottom edges of u or b respectively,
	// then we need to drop one set of left pointers
	if l != nil && lp.Eq(b.BotEdge().Left()) {
		b.Neighbors[botleft] = nil
		b.Neighbors[upleft] = nil
		l.Neighbors[botright] = u
	} else if l != nil && lp.Eq(u.TopEdge().Left()) {
		u.Neighbors[botleft] = nil
		u.Neighbors[upleft] = nil
		l.Neighbors[upright] = b
		// Does this make any sense
	} else if l == nil && ul != nil && geom.F64eq(lp.X(), ul.right) {
		// If we added a point at the same horizontal value
		// of the edge connecting tr0 to its left neighbors,
		// but we weren't on an existing vertex, we need to
		// split up the definitions for ul and bl. Specifically,
		if lp.Y() > ul.bot[right] {
			b.Neighbors[upleft] = ul
		} else if lp.Y() < ul.bot[right] {
			u.Neighbors[botleft] = bl
		}
		// We also need to change our pointer setup.
		// The X split is still on the right value, but now needs
		// to point to a y node that points to ul and bl.
		x = NewX(lp)
		y.discard(x)
		x.set(right, y)
		y2 := NewY(ul.BotEdge())
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
		p1 := u.TopEdge().Left()
		p2 := tr.TopEdge().Right()
		if geom.IsColinear(p1, u.TopEdge().Right(), p2) {
			// This can't be right?
			u.top[right] = p2.Y()
			u.bot[right] = tr.bot[right]
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
			u2tl := u2.TopEdge().Left()
			utr := u.TopEdge().Right()
			if u2tl.X() == utr.X() && u2tl.Y() == utr.Y() {
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
			u.right = u2.left
			// the bottom edge is now this shard of fe.
			// if this trapezoid is merged with another,
			// this may change.
			edge, _ = fe.SubEdge(0, u2.left, u2.right)
			u2.bot[left] = edge.Left().Y()
			u2.bot[right] = edge.Right().Y()

			// y points to a new trapezoid node holding u2
			un = NewTrapNode(u2)
			u = u2
		}

		y.set(left, un)

		p1 = b.BotEdge().Left()
		p2 = tr.BotEdge().Right()
		// Behavior is similar for the lower trapezoid
		if geom.IsColinear(p1, b.BotEdge().Right(), p2) {
			b.bot[right] = p2.Y()
			b.top[right] = tr.top[right]
		} else {
			b2 := tr.Copy()
			b.Neighbors[upright] = b2
			b2.Neighbors[upleft] = b

			b2bl := b2.BotEdge().Left()
			bbr := b.BotEdge().Right()
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
			b.right = b2.left

			edge, _ = fe.SubEdge(0, b2.left, b2.right)
			b2.top[left] = edge.Left().Y()
			b2.top[right] = edge.Right().Y()

			bn = NewTrapNode(b2)
			b = b2
		}
		u.faces = faces
		b.faces = faces
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
	if r != nil && rp.X() == b.right && rp.Y() == b.bot[right] {
		u.Rights(nil)
		r.Neighbors[upleft] = b
	} else if r != nil && rp.X() == u.right && rp.Y() == u.top[right] {
		b.Rights(nil)
		r.Neighbors[botleft] = u
	} else if r == nil && ur != nil && geom.F64eq(rp.X(), ur.left) {
		if rp.Y() > ur.bot[right] {
			b.Neighbors[upright] = ur
		} else if rp.Y() < ur.bot[right] {
			u.Neighbors[botright] = br
		}
		x = NewX(lp)
		y.discard(x)
		x.set(left, y)
		y2 := NewY(ur.BotEdge())
		x.set(right, y2)
		y2.set(left, ur.node)
		y2.set(right, br.node)
	} else if r == nil && ur == nil {
		u.Neighbors[upright] = br
		b.Neighbors[upright] = br
	}
}
