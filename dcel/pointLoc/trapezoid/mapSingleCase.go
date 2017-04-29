package trapezoid

import (
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
)

func mapSingleCase(tr *Trapezoid, fe geom.FullEdge, faces [2]*dcel.Face) {

	var l, r *Trapezoid
	ur, br, ul, bl := tr.GetNeighbors()
	lp, rp := fe.BothPoints()

	u := tr.Copy()
	d := tr.Copy()

	u.faces = faces
	d.faces = faces

	// LP does not lie on the left edge of TR
	if !geom.F64eq(lp.X(), tr.left) {
		l = tr.Copy()
		NewTopRight, _ := l.TopEdge().PointAt(0, lp.X())
		NewBotRight, _ := l.BotEdge().PointAt(0, lp.X())
		l.right = lp.X()
		l.bot[right] = NewBotRight.Y()
		l.top[right] = NewTopRight.Y()
		d.bot[left] = NewBotRight.Y()
		u.top[left] = NewTopRight.Y()
		ul.replaceNeighbors(tr, l)
		bl.replaceNeighbors(tr, l)

		l.twoRights(u, d, lp.Y())
	} else {
		tr.replaceLeftPointers(u, d, lp.Y())
	}
	if !geom.F64eq(rp.X(), tr.right) {
		r = tr.Copy()
		NewTopLeft, _ := r.TopEdge().PointAt(0, rp.X())
		NewBotLeft, _ := r.BotEdge().PointAt(0, rp.X())
		r.left = rp.X()
		r.bot[left] = NewBotLeft.Y()
		r.top[left] = NewTopLeft.Y()
		d.bot[right] = NewBotLeft.Y()
		u.top[right] = NewTopLeft.Y()
		ur.replaceNeighbors(tr, r)
		br.replaceNeighbors(tr, r)

		r.twoLefts(u, d, rp.Y())
	} else {
		tr.replaceRightPointers(u, d, rp.Y())
	}

	// D and U are exactly below // above
	// the input edge.
	splitExactly(u, d, fe)

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

	annotatedVisualize([]string{"L", "R", "U", "D"}, []*Trapezoid{l, r, u, d})
}
