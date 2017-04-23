package trapezoid

import (
	"fmt"
	"image/color"

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
	// Will bring this back in once the algorithm works
	// for i := range fullEdges {
	// 	fmt.Println(fullEdges[i])
	// 	j := i + rand.Intn(len(fullEdges)-i)
	// 	fullEdges[i], fullEdges[j] = fullEdges[j], fullEdges[i]
	// }
	if visualize.VisualCh != nil {
		visualize.HighlightColor = color.RGBA{0, 255, 0, 255}
	}
	for k, fe := range fullEdges {
		if visualize.VisualCh != nil {
			visualize.HighlightColor = color.RGBA{0, 255, 0, 255}
			visualize.DrawLine(fe.Left(), fe.Right())
		}
		// 1: Find the trapezoids intersected by fe
		trs := tree.Query(fe)
		// 2: Remove those and replace them with what they become
		//    due to the intersection of halfEdges[i]

		fmt.Println(tree)

		// Case A: A fe is contained in a single trapezoid tr
		// Then we make (up to) four trapezoids out of tr.
		if len(trs) == 0 {
			fmt.Println(fe, "intersected nothing?")
			continue
		}
		if len(trs) == 1 {
			fmt.Println(fe, "intersected one trapezoid", trs[0])
			if visualize.VisualCh != nil {
				visualize.HighlightColor = color.RGBA{0, 0, 128, 128}
				visualize.DrawPoly(trs[0].toPhysics())
			}
			mapSingleCase(trs[0], fe, faces[k])
		} else {
			fmt.Println(fe, "Intersected multiple zoids", trs)
			mapMultipleCase(trs, fe, faces[k])
		}
	}
	fmt.Println("Search:\n", tree)
	dc, m := tree.DCEL()
	return dc, m, tree, nil
}
