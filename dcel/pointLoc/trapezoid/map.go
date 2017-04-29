package trapezoid

import (
	"math/rand"

	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/pointLoc/visualize"
	"github.com/200sc/go-compgeo/geom"
)

var (
	tree *Node
	err  error
)

// TrapezoidalMap converts a dcel into a version of itself split into
// trapezoids and a search structure to find a containing trapezoid in
// the map in response to a point location query.
func TrapezoidalMap(dc *dcel.DCEL) (*dcel.DCEL, map[*dcel.Face]*dcel.Face, *Node, error) {
	bounds := dc.Bounds()

	tree = NewRoot()
	tree.payload = dc.Faces[dcel.OUTER_FACE]
	tree.set(left, NewTrapNode(newTrapezoid(bounds)))

	fullEdges, faces, err := dc.FullEdges()
	if err != nil {
		return nil, nil, nil, err
	}
	// Get rid of bad (duplicate) edges
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
	// Scramble the edges
	for i := range fullEdges {
		j := i + rand.Intn(len(fullEdges)-i)
		fullEdges[i], fullEdges[j] = fullEdges[j], fullEdges[i]
	}
	for k, fe := range fullEdges {
		visualize.HighlightColor = visualize.AddColor
		visualize.DrawLine(fe.Left(), fe.Right())
		// 1: Find the trapezoids intersected by fe
		trs := tree.Query(fe)
		// 2: Remove those and replace them with what they become
		//    due to the intersection of halfEdges[i]
		// Case A: A fe is contained in a single trapezoid tr
		// Then we make (up to) four trapezoids out of tr.
		if len(trs) == 0 {
			continue
		}
		if len(trs) == 1 {
			visualize.HighlightColor = visualize.CheckFaceColor
			visualize.DrawPoly(trs[0].toPhysics())
			mapSingleCase(trs[0], fe, faces[k])
		} else {
			mapMultipleCase(trs, fe, faces[k])
		}
	}
	dc, m := tree.DCEL()
	return dc, m, tree, nil
}
