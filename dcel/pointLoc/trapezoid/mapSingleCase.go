package trapezoid

import (
	"fmt"

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
	fmt.Println(lp, tr.left, rp, tr.right)
	if !geom.F64eq(lp.X(), tr.left) {
		fmt.Println("Not on left")
		l = tr.Copy()
		l.setRight(lp.X())
	}
	if !geom.F64eq(rp.X(), tr.right) {
		fmt.Println("Not on right")
		r = tr.Copy()
		r.setLeft(rp.X())
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
		fmt.Println("l != nil")
		// L gets in the way of needing
		// to think about ul,bl issues
		u.Lefts(l)
		d.Lefts(l)
	} else if ul != nil && !geom.F64eq(ul.bot[right], lp.Y()) {
		// If these values are equal, our original assignments
		// above are fine.
		if ul.bot[right] > lp.Y() {
			fmt.Println("Above lp")
			u.Neighbors[botleft] = bl
		} else {
			fmt.Println("Below lp")
			d.Neighbors[upleft] = ul
		}
	} else if ul == nil && bl != nil {
		fmt.Println("Not inline")
		u.Lefts(bl)
	}
	if r != nil {
		fmt.Println("r != nil")
		u.Rights(r)
		d.Rights(r)
	} else if ur != nil && !geom.F64eq(ur.bot[right], rp.Y()) {
		if ur.bot[right] > rp.Y() {
			u.Neighbors[botright] = br
			fmt.Println("Above rp")
		} else {
			fmt.Println("Below rp")
			d.Neighbors[upright] = ur
		}
	} else if ur == nil && br != nil {
		fmt.Println("right not inline")
		u.Rights(br)
	}

	// The border between these two
	// is explicitly defined by fe.
	d.top[left] = lp.Y()
	d.top[right] = rp.Y()
	u.bot[left] = lp.Y()
	u.bot[right] = rp.Y()

	u.setLeft(lp.X())
	d.setLeft(lp.X())
	u.setRight(rp.X())
	d.setRight(rp.X())
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
}
