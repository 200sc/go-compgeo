package triangulation

import (
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
	Search.left = NewTrapNode(newTrapezoid(bounds))

	fullEdges, faces, err := dc.FullEdges()
	if err != nil {
		return nil, nil, nil, err
	}

	// Scramble the edges
	for i := range fullEdges {
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

		// Case A: A fe is contained in a single trapezoid tr
		// Then we make (up to) four trapezoids out of tr.
		if len(trs) == 1 {
			mapSingleCase(trs[0], fe)
		} else {
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

			p1, _ := tr0.Edges[top].PointAt(0, lp.X())
			p2, _ := tr0.Edges[bot].PointAt(0, lp.X())
			if !tr0.HasDefinedPoint(lp) {
				// The three trapezoids are split into
				// one to the left of an x node
				x = NewX(lp)
				// and two below the previous y node
				l = tr0.Copy()
				l.Edges[right] = geom.NewFullEdge(p1, p2)

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
			u.Edges[left] = geom.NewFullEdge(p1, lp)
			b.Edges[left] = geom.NewFullEdge(p2, lp)

			u.Edges[bot], _ = fe.SubEdge(0, u.Edges[left].Left().X(),
				u.Edges[right].Right().X())
			u.Neighbors[upleft] = ul
			u.Neighbors[botleft] = bl

			b.Edges[top], _ = fe.SubEdge(0, b.Edges[left].Left().X(),
				b.Edges[right].Right().X())
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
			} else if l == nil && lp.X() == ul.Edges[bot].Right().X() {
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
						// The left edge of u2, and the right edge of u, are
						// the ray cast down from u.Edges[top].Right() to
						// fe.
						p1 = utr
						p2, _ = fe.PointAt(0, p1.X())
						e := geom.NewFullEdge(p1, p2)
						u2.Edges[left] = e
						u.Edges[right] = e
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
						// u's right edge is the segment from fe to
						// u's top.
						p1 = utr
						p2, _ = fe.PointAt(0, p1.X())
						u.Edges[right] = geom.NewFullEdge(p1, p2)
						// u2's left edge is the segment from fe to u2's
						// top.
						p3, _ := u.Edges[top].PointAt(0, p1.X())
						u2.Edges[left] = geom.NewFullEdge(p2, p3)
					}
					// the bottom edge is now this shard of fe.
					// if this trapezoid is merged with another,
					// this may change.
					u2.Edges[bot], _ = fe.SubEdge(0, u2.Edges[left].Left().X(),
						u2.Edges[right].Right().X())
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
						p1 = bbr
						p2, _ = fe.PointAt(0, p1.X())
						e := geom.NewFullEdge(p1, p2)
						b2.Edges[left] = e
						b.Edges[right] = e
					} else {
						if b2bl.Y() < bbr.Y() {
							b.Neighbors[botright] = b2
						} else {
							b2.Neighbors[botleft] = b
						}
						p1 = bbr
						p2, _ = fe.PointAt(0, p1.X())
						b.Edges[right] = geom.NewFullEdge(p1, p2)
						p3, _ := u.Edges[bot].PointAt(0, p1.X())
						b2.Edges[left] = geom.NewFullEdge(p2, p3)
					}
					b2.Edges[top], _ = fe.SubEdge(0, b2.Edges[left].Left().X(),
						b2.Edges[right].Right().X())
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
			p2, _ = u.Edges[top].PointAt(0, rp.X())
			p3, _ := b.Edges[bot].PointAt(0, rp.X())
			u.Edges[right] = geom.NewFullEdge(rp, p2)
			b.Edges[right] = geom.NewFullEdge(rp, p3)

			var r, ur, br *Trapezoid
			trn := trs[len(trs)-1]
			if !trn.HasDefinedPoint(fe.Right()) {
				r = trn.Copy()
				r.Edges[left] = geom.NewFullEdge(p2, p3)
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
			} else if r == nil && rp.X() == ur.Edges[bot].Left().X() {
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
			}
		}
	}
	dc, m := Search.DCEL()
	return dc, m, Search, nil
}

func mapSingleCase(tr *Trapezoid, fe geom.FullEdge) {
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
	if !geom.IsColinear(lp, tr.Edges[left].Left(), tr.Edges[left].Right()) {
		l = tr.Copy()
		p1, _ := tr.Edges[top].PointAt(0, lp.X())
		p2, _ := tr.Edges[bot].PointAt(0, lp.X())
		l.Edges[right] = geom.NewFullEdge(p1, p2)
	}
	if !geom.IsColinear(rp, tr.Edges[right].Left(), tr.Edges[right].Right()) {
		r = tr.Copy()
		p1, _ := tr.Edges[top].PointAt(0, rp.X())
		p2, _ := tr.Edges[bot].PointAt(0, rp.X())
		r.Edges[left] = geom.NewFullEdge(p1, p2)
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
	u.Edges[left] = geom.NewFullEdge(lp, p)
	p, _ = tr.Edges[bot].PointAt(0, lp.X())
	d.Edges[left] = geom.NewFullEdge(lp, p)
	p, _ = tr.Edges[top].PointAt(0, rp.X())
	u.Edges[right] = geom.NewFullEdge(rp, p)
	p, _ = tr.Edges[bot].PointAt(0, rp.X())
	d.Edges[right] = geom.NewFullEdge(rp, p)
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

	a.set(left, NewTrapNode(l))
	a.set(right, b)
	b.set(left, c)
	b.set(right, NewTrapNode(r))
	c.set(left, NewTrapNode(u))
	c.set(right, NewTrapNode(d))
}
