package pointLoc

import (
	"math/rand"
	"testing"

	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/pointLoc/rtree"
	"github.com/200sc/go-compgeo/dcel/pointLoc/slab"
	"github.com/200sc/go-compgeo/dcel/pointLoc/voronoi"
	"github.com/200sc/go-compgeo/geom"
	"github.com/stretchr/testify/assert"
)

var (
	inputSize  = 1000
	inputRange = 10000.0
	testCt     = 500
)

func randomPt() geom.D3 {
	return geom.NewPoint(rand.Float64()*inputRange,
		rand.Float64()*inputRange, 0)
}

func randomDCEL() *dcel.DCEL {
	// generate n random points
	pts := make([]geom.D2, inputSize)
	for i := 0; i < inputSize; i++ {
		pts[i] = randomPt()
	}
	// from those points generate a set of non-overlapping faces
	// ^ ^ ^ the hard bit
	// we'll use a voronoi diagram, see fortune.go
	return voronoi.Fortune(pts)
}

func TestRandomDCEL(t *testing.T) {
	dc := randomDCEL()

	// We assume an rtree will be correct, and test against it.

	tree := rtree.DCELtoRtree(dc)
	structure := slab.Decompose(dc)

	for i := 0; i < testCt; i++ {
		pt := randomPt()
		treeIntersected := rtree.SearchIntersect(tree, pt)
		structIntersected := structure.PointLocate(pt[0], pt[1])
		assert.Contains(t, treeIntersected, structIntersected)
	}
}