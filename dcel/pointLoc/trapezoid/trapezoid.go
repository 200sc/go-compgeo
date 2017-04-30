// trapezoid holds structures and functions for point location through
// mapping a dcel to a trapezoidal map.

package trapezoid

import (
	"fmt"
	"image/color"

	"bitbucket.org/oakmoundstudio/oak/physics"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/pointLoc/visualize"
	"github.com/200sc/go-compgeo/geom"
	"github.com/200sc/go-compgeo/printutil"
)

// These constants refer to indices
// within trapezoids' Edges
const (
	top = iota
	bot
)
const (
	left = iota
	right
)

// These constants refer to indices
// within trapezoids' Neighbors
const (
	upright = iota
	botright
	upleft
	botleft
)

// A Trapezoid is used when contstructing a Trapezoid map,
// and contains references to its neighbor trapezoids and
// the edges that border it.
type Trapezoid struct {
	// See above indices
	top         [2]float64 // y values
	bot         [2]float64 // y values
	left, right float64    // x values
	Neighbors   [4]*Trapezoid
	node        *Node
	faces       [2]*dcel.Face
}

func (tr *Trapezoid) GetNeighbors() (*Trapezoid, *Trapezoid, *Trapezoid, *Trapezoid) {
	return tr.Neighbors[0], tr.Neighbors[1], tr.Neighbors[2], tr.Neighbors[3]
}

// DCELEdges evaluates and returns the edges of
// a trapezoid as DCElEdges with initialized origins,
// prevs, and nexts.
func (tr *Trapezoid) DCELEdges() []*dcel.Edge {
	edges := make([]*dcel.Edge, 1)
	i := 0
	edges[i] = dcel.NewEdge()
	edges[i].Origin = dcel.PointToVertex(geom.NewPoint(tr.left, tr.top[left], 0))
	edges[i].Origin.OutEdge = edges[i]
	p := geom.NewPoint(tr.right, tr.top[right], 0)
	if !p.Eq(edges[i].Origin) {
		i++
		edges = append(edges, dcel.NewEdge())
		edges[i].Origin = dcel.PointToVertex(p)
		edges[i].Origin.OutEdge = edges[i]
		edges[i-1].Next = edges[i]
		edges[i].Prev = edges[i-1]
	}
	p = geom.NewPoint(tr.right, tr.bot[right], 0)
	if !p.Eq(edges[i].Origin) {
		i++
		edges = append(edges, dcel.NewEdge())
		edges[i].Origin = dcel.PointToVertex(p)
		edges[i].Origin.OutEdge = edges[i]
		edges[i-1].Next = edges[i]
		edges[i].Prev = edges[i-1]
	}
	p = geom.NewPoint(tr.left, tr.bot[left], 0)
	if !p.Eq(edges[i].Origin) &&
		!p.Eq(edges[0].Origin) {
		i++
		edges = append(edges, dcel.NewEdge())
		edges[i].Origin = dcel.PointToVertex(p)
		edges[i].Origin.OutEdge = edges[i]
		edges[i-1].Next = edges[i]
		edges[i].Prev = edges[i-1]
	}
	// In the case of a trapezoid which is a point,
	// this will cause the edge to refer to itself by next
	// and prev, which is probably not expected by code
	// which iterates over edges.
	edges[0].Prev = edges[i]
	edges[i].Next = edges[0]
	return edges
}

// Rights is shorthand for setting both of
// tr's right neighbors to the same value
func (tr *Trapezoid) Rights(tr2 *Trapezoid) {
	tr.Neighbors[upright] = tr2
	tr.Neighbors[botright] = tr2
}

// Lefts is shorthand for setting both of
// tr's left neighbors to the same value.
func (tr *Trapezoid) Lefts(tr2 *Trapezoid) {
	tr.Neighbors[upleft] = tr2
	tr.Neighbors[botleft] = tr2
}

// Copy returns a trapezoid with identical edges
// and neighbors.
func (tr *Trapezoid) Copy() *Trapezoid {
	tr2 := new(Trapezoid)
	tr2.top = tr.top
	tr2.bot = tr.bot
	tr2.left = tr.left
	tr2.right = tr.right
	tr2.Neighbors = tr.Neighbors
	tr2.faces = tr.faces
	return tr2
}

// AsPoints converts a trapezoid's internal values
// into four points.
func (tr *Trapezoid) AsPoints() []geom.D2 {
	ds := make([]geom.D2, 4)
	ds[0] = geom.NewPoint(tr.left, tr.top[left], 0)
	ds[1] = geom.NewPoint(tr.right, tr.top[right], 0)
	ds[2] = geom.NewPoint(tr.right, tr.bot[right], 0)
	ds[3] = geom.NewPoint(tr.left, tr.bot[left], 0)
	return ds
}

// BotEdge returns a translation of tr's values to
// tr's bottom edge as a FullEdge
func (tr *Trapezoid) BotEdge() geom.FullEdge {
	return geom.FullEdge{
		geom.NewPoint(tr.right, tr.bot[right], 0),
		geom.NewPoint(tr.left, tr.bot[left], 0),
	}
}

// TopEdge acts as BotEdge for tr's top
func (tr *Trapezoid) TopEdge() geom.FullEdge {
	return geom.FullEdge{
		geom.NewPoint(tr.left, tr.top[left], 0),
		geom.NewPoint(tr.right, tr.top[right], 0),
	}
}

// HasDefinedPoint returns for a given Trapezoid
// whether or not any of the points on the Trapezoid's
// perimeter match the query point.
// We make an assumption here that there will be no
// edges who have vertices defined on other edges, aka
// that all intersections are represented through
// vertices.
func (tr *Trapezoid) HasDefinedPoint(p geom.D3) bool {
	for _, p2 := range tr.AsPoints() {
		if p2.X() == p.X() && p2.Y() == p.Y() {
			return true
		}
	}
	return false
}

func (tr *Trapezoid) String() string {
	s := ""
	for _, p := range tr.AsPoints() {
		s += "("
		s += printutil.Stringf64(p.X(), p.Y())
		s += ")"
	}
	return s
}

func newTrapezoid(sp geom.Span) *Trapezoid {
	t := new(Trapezoid)
	min := sp.At(geom.SPAN_MIN).(geom.Point)
	max := sp.At(geom.SPAN_MAX).(geom.Point)
	t.top[left] = max.Y()
	t.top[right] = max.Y()
	t.bot[left] = min.Y()
	t.bot[right] = min.Y()
	t.left = min.X()
	t.right = max.X()
	t.Neighbors = [4]*Trapezoid{nil, nil, nil, nil}
	return t
}

func (tr *Trapezoid) toPhysics() []physics.Vector {
	vs := make([]physics.Vector, 4)
	for i, p := range tr.AsPoints() {
		vs[i] = physics.NewVector(p.X(), p.Y())
	}
	return vs
}

func (tr *Trapezoid) visualize() {
	if tr == nil {
		fmt.Println("Nothing to visualize")
		return
	}
	visualize.HighlightColor = visualize.AddFaceColor
	visualize.DrawPoly(tr.toPhysics())
}

func (tr *Trapezoid) visualizeNeighbors() {
	if tr == nil {
		return
	}
	visualize.HighlightColor = color.RGBA{128, 0, 128, 128}
	for _, n := range tr.Neighbors {
		if n == nil {
			continue
		}
		visualize.DrawPoly(n.toPhysics())
	}
}

// Replace neighbors runs replace neighbor for all directions
// off of a trapezoid
func (tr *Trapezoid) replaceNeighbors(rep, new *Trapezoid) {
	if tr == nil {
		return
	}
	for i := range tr.Neighbors {
		tr.replaceNeighbor(i, rep, new)
	}
}

// Replace neighbor checks that the input is not nil,
// and if it is not, if it's neighbor in the given direction is the
// expected trapezoid to replace, replaces it with the given new trapezoid.
func (tr *Trapezoid) replaceNeighbor(dir int, rep, new *Trapezoid) {
	if tr == nil {
		return
	}
	if tr.Neighbors[dir] == rep {
		tr.Neighbors[dir] = new
	}
}

// Assign the neighbors of the trapezoid tr's upleft and upright
// neighbors (if they exist) dependant on tr being replaced by
// the two trapezoids u and b split at y value lpy.
//
//  ~ ~ ~ ~ ~ ~
//    ul |  u
//  ~ ~ ~ -lpy-----
//    bl |  b
//  ~ ~ ~ ~ ~ ~
func (tr *Trapezoid) replaceLeftPointers(u, b *Trapezoid, lpy float64) {
	replaceLeftPointers(tr, tr.Neighbors[upleft], tr.Neighbors[botleft], u, b, lpy)
}

// Given the trapezoid tr, being replaced by u and b where
// u is above b and lpy is the point at which u and be connect
// on tr's left edge, assign all pointers from ul and bl where ul
// is above bl that previously pointed to tr to the appropriate
// trapezoid of u and b.
func replaceLeftPointers(tr, ul, bl, u, b *Trapezoid, lpy float64) {
	if ul != nil && geom.F64eq(ul.bot[right], lpy) {
		fmt.Println("Case 0")
		// U matches exactly to ul,
		// B matches exactly to bl.
		//
		//  ~ ~ ~ ~ ~ ~
		//    ul |  u
		// -----lpy-----
		//    bl |  b
		//  ~ ~ ~ ~ ~ ~
		//
		ul.replaceNeighbors(tr, u)
		bl.replaceNeighbors(tr, b)
	} else if (ul != nil && geom.F64eq(ul.top[right], lpy)) ||
		(ul == nil && bl != nil && geom.F64eq(bl.top[right], lpy)) {
		fmt.Println("Case 1")
		// U does not border the left edge
		//
		// ~ ~ ~ lpy \
		//  (ul)  \   \ u
		// ~ ~ ~ ~ \ b \
		//  (bl)    \   \
		// ~ ~ ~ ~ ~ ~ ~
		ul.replaceNeighbors(tr, b)
		bl.replaceNeighbors(tr, b)
		if ul != nil {
			b.Neighbors[upleft] = ul
		} else {
			b.Neighbors[upleft] = bl
		}
		if bl != nil {
			b.Neighbors[botleft] = bl
		} else {
			b.Neighbors[botleft] = ul
		}
		u.Lefts(b)
	} else if (bl != nil && geom.F64eq(bl.bot[right], lpy)) ||
		(bl == nil && ul != nil && geom.F64eq(ul.bot[right], lpy)) {
		fmt.Println("Case 2")
		// D does not border the left edge
		//
		// ~ ~ ~ ~ ~ ~ ~ ~
		//  (ul)    /   /
		// ~ ~ ~ ~ / u /
		//  (bl)  /   / b
		// ~ ~ ~ lpy / ~ ~
		//
		ul.replaceNeighbors(tr, u)
		bl.replaceNeighbors(tr, u)
		if bl != nil {
			u.Neighbors[botleft] = bl
		} else {
			u.Neighbors[botleft] = ul
		}
		if ul != nil {
			u.Neighbors[upleft] = ul
		} else {
			u.Neighbors[upleft] = bl
		}
		b.Lefts(u)
	} else if ul != nil && ul.bot[right] < lpy {
		fmt.Println("Case 3")
		// UL expands past FE
		//
		// ~ ~ ~ ~ ~ ~ ~ ~ ~
		//   ul    |   u
		//        lpy ~ ~ ~
		// ~ ~ ~ ~ |   b
		//   bl    |
		// ~ ~ ~ ~ ~ ~ ~ ~ ~
		//
		if ul != bl {
			bl.replaceNeighbors(tr, b)
		}
		ul.replaceNeighbor(upright, tr, u)
		ul.replaceNeighbor(botright, tr, b)
		u.Lefts(ul)
		b.Neighbors[upleft] = ul
		b.Neighbors[botleft] = bl
	} else if bl != nil && bl.top[right] > lpy {
		fmt.Println("Case 4")
		// BL expands past FE
		//
		// ~ ~ ~ ~ ~ ~ ~ ~ ~
		//   ul    |   u
		// ~ ~ ~ ~ |
		//         lpy ~ ~ ~
		//   bl    |   b
		// ~ ~ ~ ~ ~ ~ ~ ~ ~
		//
		if ul != bl {
			ul.replaceNeighbors(tr, u)
		}
		bl.replaceNeighbor(upright, tr, u)
		bl.replaceNeighbor(botright, tr, b)
		b.Lefts(bl)
		u.Neighbors[upleft] = ul
		u.Neighbors[botleft] = bl
	} else {
		fmt.Println("No left pointer case satisfied.")
		fmt.Println("Both neighbors nil:", ul == nil && bl == nil)
		fmt.Println(ul, bl)
		fmt.Println(u, b)
		fmt.Println(lpy)
	}
}

//  ~ ~ ~ ~ ~ ~
//    u |  ur
//  --rpy ~ ~ ~
//    b |  br
//  ~ ~ ~ ~ ~ ~
func (tr *Trapezoid) replaceRightPointers(u, b *Trapezoid, rpy float64) {
	replaceRightPointers(tr, tr.Neighbors[upright], tr.Neighbors[botright], u, b, rpy)
}

func replaceRightPointers(tr, ur, br, u, b *Trapezoid, rpy float64) {
	if ur != nil && geom.F64eq(ur.bot[left], rpy) {
		fmt.Println("Case 0")
		// U matches exactly to ur,
		// B matches exactly to br.
		//
		//  ~ ~ ~ ~ ~ ~
		//    u  |  ur
		// -----rpy-----
		//    b  |  br
		//  ~ ~ ~ ~ ~ ~
		//
		ur.replaceNeighbors(tr, u)
		br.replaceNeighbors(tr, b)
	} else if (ur != nil && geom.F64eq(ur.top[left], rpy)) ||
		(ur == nil && br != nil && geom.F64eq(br.top[left], rpy)) {
		fmt.Println("Right case 1")
		// U does not border the right edge
		//
		//  ~ ~ rpy ~ ~ ~
		//   u /   / (ur)
		//    / b / ~ ~ ~
		//   /   /  (br)
		//  ~ ~ ~ ~ ~ ~
		//
		ur.replaceNeighbors(tr, b)
		br.replaceNeighbors(tr, b)
		u.Rights(b)
		if ur != nil {
			b.Neighbors[upright] = ur
		} else {
			b.Neighbors[upright] = br
		}
		if br != nil {
			b.Neighbors[botright] = br
		} else {
			b.Neighbors[botright] = ur
		}
	} else if (br != nil && geom.F64eq(br.bot[left], rpy)) ||
		(br == nil && ur != nil && geom.F64eq(ur.bot[left], rpy)) {
		fmt.Println("Right case 2")
		//
		//  ~ ~ rpy ~ ~ ~
		//  \   \    ur
		//   \ u \ ~ ~ ~
		//  b \   \ br
		//  ~ ~ rpy ~ ~ ~
		//
		ur.replaceNeighbors(tr, u)
		br.replaceNeighbors(tr, u)
		b.Rights(u)
		if ur != nil {
			u.Neighbors[upright] = ur
		} else {
			u.Neighbors[upright] = br
		}
		if br != nil {
			u.Neighbors[botright] = br
		} else {
			u.Neighbors[botright] = ur
		}
	} else if ur != nil && ur.bot[left] < rpy {
		fmt.Println("Right case 3")
		// UR expands past FE
		//
		// ~ ~ ~ ~ ~ ~ ~ ~ ~
		//   u     |   ur
		// ~ ~ ~ ~rpy
		//         | ~ ~ ~ ~
		//   b     |   br
		// ~ ~ ~ ~ ~ ~ ~ ~ ~
		//
		if ur != br {
			br.replaceNeighbors(tr, b)
		}
		ur.replaceNeighbor(upleft, tr, u)
		ur.replaceNeighbor(botleft, tr, b)
		u.Rights(ur)
		b.Neighbors[upright] = ur
		b.Neighbors[botright] = br
	} else if br != nil && br.top[left] > rpy {
		fmt.Println("Right case 4")
		// BR expands past FE
		//
		// ~ ~ ~ ~ ~ ~ ~ ~ ~
		//   u     |   ur
		//         | ~ ~ ~ ~
		// ~ ~ ~ ~rpy   br
		//   b     |
		// ~ ~ ~ ~ ~ ~ ~ ~ ~
		//
		if ur != br {
			ur.replaceNeighbors(tr, u)
		}
		br.replaceNeighbor(upleft, tr, u)
		br.replaceNeighbor(botleft, tr, b)
		b.Rights(br)
		u.Neighbors[upright] = ur
		u.Neighbors[botright] = br
	} else {
		fmt.Println("No right pointer case satisfied")
		fmt.Println("Both neighbors nil:", br == nil && ur == nil)
	}
}

func (tr *Trapezoid) twoRights(u, b *Trapezoid, lpy float64) {
	tr.Neighbors[upright] = u
	tr.Neighbors[botright] = b
	if geom.F64eq(tr.top[right], lpy) {
		tr.Neighbors[upright] = b
	} else if geom.F64eq(tr.bot[right], lpy) {
		tr.Neighbors[botright] = u
	}
	u.Lefts(tr)
	b.Lefts(tr)
}

func (tr *Trapezoid) twoLefts(u, b *Trapezoid, rpy float64) {
	tr.Neighbors[upleft] = u
	tr.Neighbors[botleft] = b
	if geom.F64eq(tr.top[left], rpy) {
		tr.Neighbors[upleft] = b
	} else if geom.F64eq(tr.bot[left], rpy) {
		tr.Neighbors[botleft] = u
	}
	u.Rights(tr)
	b.Rights(tr)
}

func splitExactly(u, d *Trapezoid, fe geom.FullEdge) {
	u.exactly(top, fe)
	d.exactly(bot, fe)
}

func (tr *Trapezoid) exactly(d int, fe geom.FullEdge) {
	lp := fe.Left()
	rp := fe.Right()
	tr.left = lp.X()
	tr.right = rp.X()
	if d == bot {
		tr.top[left] = lp.Y()
		tr.top[right] = rp.Y()
	} else if d == top {
		tr.bot[left] = lp.Y()
		tr.bot[right] = rp.Y()
	}
}

func annotatedVisualize(strs []string, trs []*Trapezoid) {
	for i, s := range strs {
		fmt.Println("Visualizing " + s)
		trs[i].visualize()
		//trs[i].visualizeNeighbors()
		fmt.Print(s + " visualized.")
	}
	fmt.Println("")
}

func (tr *Trapezoid) setBotleft(fe geom.FullEdge) {
	r := tr.right
	if r > fe.Right().X() {
		r = fe.Right().X()
	}
	l := tr.left
	if l < fe.Left().X() {
		l = fe.Left().X()
	}
	edge, _ := fe.SubEdge(0, l, r)
	tr.bot[left] = edge.Left().Y()
	tr.bot[right] = edge.Right().Y()
}

func (tr *Trapezoid) setTopleft(fe geom.FullEdge) {
	r := tr.right
	if r > fe.Right().X() {
		r = fe.Right().X()
	}
	l := tr.left
	if l < fe.Left().X() {
		l = fe.Left().X()
	}
	edge, _ := fe.SubEdge(0, l, r)
	tr.top[left] = edge.Left().Y()
	tr.top[right] = edge.Right().Y()
}
