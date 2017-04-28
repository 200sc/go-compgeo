package trapezoid

import (
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
)

func mapSingleCase(tr *Trapezoid, fe geom.FullEdge, faces [2]*dcel.Face) {
	lp := fe.Left()
	rp := fe.Right()

	var l, r *Trapezoid
	ul := tr.Neighbors[upleft]
	bl := tr.Neighbors[botleft]
	ur := tr.Neighbors[upright]
	br := tr.Neighbors[botright]

	// Case 2A.2
	// If fe.left or fe.right lies ON tr's left and right
	// edges, we don't make new trapezoids for them.
	if !geom.F64eq(lp.X(), tr.left) {
		l = tr.Copy()
		l.right = lp.X()
	}
	if !geom.F64eq(rp.X(), tr.right) {
		r = tr.Copy()
		r.left = rp.X()
	}

	u := tr.Copy()
	u.faces = faces
	d := tr.Copy()
	d.faces = faces

	u.Neighbors[upleft] = ul
	u.Neighbors[botleft] = ul
	u.Neighbors[upright] = ur
	u.Neighbors[botright] = ur

	d.Neighbors[upleft] = bl
	d.Neighbors[botleft] = bl
	d.Neighbors[upright] = br
	d.Neighbors[botright] = br

	if l != nil {
		// L gets in the way of needing
		// to think about ul,bl issues
		l.Neighbors[upright] = u
		l.Neighbors[botright] = d
		u.Lefts(l)
		d.Lefts(l)
	} else if ul != nil && !geom.F64eq(ul.bot[right], lp.Y()) {
		// If these values are equal, our original assignments
		// above are fine.
		if ul.bot[right] > lp.Y() {
			u.Neighbors[botleft] = bl
		} else {
			d.Neighbors[upleft] = ul
		}
	} else if ul == nil && bl != nil {
		u.Lefts(bl)
	}
	if r != nil {
		r.Neighbors[upleft] = u
		r.Neighbors[botleft] = d
		u.Rights(r)
		d.Rights(r)
	} else if ur != nil && !geom.F64eq(ur.bot[right], rp.Y()) {
		if ur.bot[right] > rp.Y() {
			u.Neighbors[botright] = br
		} else {
			d.Neighbors[upright] = ur
		}
	} else if ur == nil && br != nil {
		u.Rights(br)
	}

	// The border between these two
	// is explicitly defined by fe.
	d.top[left] = lp.Y()
	d.top[right] = rp.Y()
	u.bot[left] = lp.Y()
	u.bot[right] = rp.Y()

	u.left = lp.X()
	d.left = lp.X()
	u.right = rp.X()
	d.right = rp.X()
	// 3: From the query structure, remove the leaves of the
	//    removed trapezoids and add new leaves for the new
	//    trapezoids, with additional inner nodes as necessary.

	// This part seems to work fine!

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

	l.visualize()
	r.visualize()
	u.visualize()
	d.visualize()
}
