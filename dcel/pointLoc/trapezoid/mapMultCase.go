package trapezoid

import (
	"fmt"

	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
)

func mapMultipleCase(trs []*Trapezoid, fe geom.FullEdge, faces [2]*dcel.Face) {

	lp, rp := fe.BothPoints()
	// Case B: fe is contained by more than one trapezoid
	// Step 1: if either fe.Left() or fe.Right() is not already
	// in the search structure, we define three new trapezoids by
	// drawing rays up and down from each new point.
	var ln, un, bn, x *Node

	y := NewY(fe)

	u := trs[0].Copy()
	b := trs[0].Copy()
	u.faces = faces
	b.faces = faces

	u.left = lp.X()
	u.setBotleft(fe)

	b.left = lp.X()
	b.setTopleft(fe)

	// At this point we have u and b defined as
	//       . .
	//      |   u
	//      |
	//   l? |-- . . .
	//      |
	//      |   b
	//       . .
	// with no neighbors defined

	if !geom.F64eq(lp.X(), trs[0].left) {
		fmt.Println("L exists")
		// The three trapezoids are split into
		// one to the left of an x node
		// and two below the previous y node
		x = NewX(lp)
		l := trs[0].Copy()
		NewTopRight, _ := l.TopEdge().PointAt(0, lp.X())
		NewBotRight, _ := l.BotEdge().PointAt(0, lp.X())
		l.right = lp.X()
		l.bot[right] = NewBotRight.Y()
		l.top[right] = NewTopRight.Y()
		b.bot[left] = NewBotRight.Y()
		u.top[left] = NewTopRight.Y()
		l.Neighbors[upleft].replaceNeighbors(trs[0], l)
		l.Neighbors[botleft].replaceNeighbors(trs[0], l)

		ln = NewTrapNode(l)

		trs[0].node.discard(x)
		x.set(left, ln)
		x.set(right, y)

		l.twoRights(u, b, lp.Y())

		annotatedVisualize([]string{"L"}, []*Trapezoid{l})

	} else {
		fmt.Println("No L")
		// Otherwise we just split trs[0] into two trapezoids.
		trs[0].node.discard(y)
		trs[0].replaceLeftPointers(u, b, lp.Y())
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

		y = NewY(fe)

		// It is possible that one or both trapezoids
		// we make are mergeable into the previous upper
		// and lower trapezoids we made.

		p1 := u.TopEdge().Left()
		p2 := tr.TopEdge().Right()
		if geom.IsColinear(p1, u.TopEdge().Right(), p2) &&
			geom.IsColinear(p1, tr.TopEdge().Left(), p2) {
			fmt.Println("Merge U", u)
			u.top[right] = p2.Y()
			u.right = tr.right // temporary, will be replaced by later loops
			u.setBotleft(fe)
			fmt.Println("Merged:", u)
		} else {
			fmt.Println("Non merge U")
			u2 := tr.Copy()
			//
			u.Neighbors[botright] = u2
			u2.Neighbors[botleft] = u
			// There are three reasons we might not be
			// able to merge.
			//
			// A: this trapezoid's upper edge
			// shares a vertex with the previous trapezoid, but
			// is at a different angle.
			u2tl := u2.TopEdge().Left()
			utr := u.TopEdge().Right()
			if u2tl.X() == utr.X() && u2tl.Y() == utr.Y() {
				fmt.Println("U2tl and utr equal")
				// In this case, u2's top left and bot left are both u.
				// u's bot right and bot left are similarly both u2.
				u.Neighbors[upright] = u2
				u2.Neighbors[upleft] = u
				// The top edges of u and u2 do not need to be updated.
				// The top edge of u2 is still accurate from the copy.
				// the search structure is updated later.
				//
				// tr's left neighbors('s neighbors) do not need to be updated,
				// because both left neighbors were consumed by u.
			} else if u2tl.Y() > utr.Y() {
				fmt.Println("U2tl above utr")
				// B: this trapezoid's left endpoint is above
				// the left endpoint of the previous trapezoid.
				u.Neighbors[upright] = u2
				u2.Neighbors[upleft].replaceNeighbors(tr, u2)
			} else {
				fmt.Println("U2tl below utr")
				// C: this trapezoid's left endpoint is below
				// the left endpoint of the previous trapezoid.
				u2.Neighbors[upleft] = u
				// U's upright neighbor better point to u instead
				// of the former trapezoid by now.
				u.Neighbors[upright].replaceNeighbors(trs[i-1], u)
			}
			u.right = u2.left
			// the bottom edge is now this shard of fe.
			// if this trapezoid is merged with another,
			// this may change.
			u2.setBotleft(fe)

			// y points to a new trapezoid node holding u2
			un = NewTrapNode(u2)
			annotatedVisualize([]string{"U"}, []*Trapezoid{u})
			u = u2
		}

		p1 = b.BotEdge().Left()
		p2 = tr.BotEdge().Right()
		if geom.IsColinear(p1, b.BotEdge().Right(), p2) &&
			geom.IsColinear(p1, tr.BotEdge().Left(), p2) {
			fmt.Println("Merge B", b)
			b.bot[right] = p2.Y()
			b.right = tr.right
			b.setTopleft(fe)
			fmt.Println("Merged:", b)
		} else {
			fmt.Println("Did not merge B")
			b2 := tr.Copy()
			b.Neighbors[upright] = b2
			b2.Neighbors[upleft] = b

			b2bl := b2.BotEdge().Left()
			bbr := b.BotEdge().Right()
			if b2bl.X() == bbr.X() && b2bl.Y() == bbr.Y() {
				fmt.Println("Equal b2bl and bbr")
				b.Neighbors[botright] = b2
				b2.Neighbors[botleft] = b
			} else if b2bl.Y() < bbr.Y() {
				fmt.Println("b2bl below bbr")
				b.Neighbors[botright] = b2
				b2.Neighbors[botleft].replaceNeighbors(tr, b2)
			} else {
				fmt.Println("b2bl above bbr")
				b2.Neighbors[botleft] = b
				b.Neighbors[botright].replaceNeighbors(trs[i-1], b)
			}
			b.right = b2.left

			b2.setTopleft(fe)

			bn = NewTrapNode(b2)
			annotatedVisualize([]string{"B"}, []*Trapezoid{b})
			b = b2
		}
		u.faces = faces
		b.faces = faces
		tr.node.discard(y)
		y.set(left, un)
		y.set(right, bn)
	}
	// If fe.right is on some edge,
	// then we're done, except we need
	// to give u and b right edges, which
	// will be the same in the next case

	var r *Trapezoid
	u.right = rp.X()
	b.right = rp.X()

	trn := trs[len(trs)-1]

	if !geom.F64eq(rp.X(), trn.right) {
		fmt.Println("RP Not On Right Edge, TRN")
		r = trn.Copy()

		NewTopLeft, _ := r.TopEdge().PointAt(0, rp.X())
		NewBotLeft, _ := r.BotEdge().PointAt(0, rp.X())
		r.left = rp.X()
		r.bot[left] = NewBotLeft.Y()
		r.top[left] = NewTopLeft.Y()
		b.bot[right] = NewBotLeft.Y()
		u.top[right] = NewTopLeft.Y()
		r.Neighbors[upright].replaceNeighbors(trn, r)
		r.Neighbors[botright].replaceNeighbors(trn, r)

		r.twoLefts(u, b, rp.Y())

		x = NewX(rp)
		// X needs to be put between y's
		// parents and y
		y.discard(x)
		x.set(left, y)
		x.set(right, NewTrapNode(r))

	} else {
		fmt.Println("RP On Right edge, TRN")
		trn.replaceRightPointers(u, b, rp.Y())
	}

	fmt.Println("B", b)
	annotatedVisualize([]string{"U", "B", "R"}, []*Trapezoid{u, b, r})
}
