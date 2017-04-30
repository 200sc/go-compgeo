package test

import (
	"math/rand"
	"testing"

	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/pointLoc"
	"github.com/200sc/go-compgeo/dcel/pointLoc/bench/bruteForce"
	"github.com/200sc/go-compgeo/dcel/pointLoc/bench/slab"
	"github.com/200sc/go-compgeo/dcel/pointLoc/bench/trapezoid"
	"github.com/200sc/go-compgeo/dcel/pointLoc/rtree"
	"github.com/200sc/go-compgeo/geom"
	"github.com/200sc/go-compgeo/search/tree"
	"github.com/stretchr/testify/assert"
)

var (
	inputSize  = 100
	inputRange = 10000.0
	testCt     = 1000
)

func randomPt() geom.D3 {
	return geom.NewPoint(rand.Float64()*inputRange,
		rand.Float64()*inputRange, 0)
}

func testRandomPts(t *testing.T, pl pointLoc.LocatesPoints, limit int) {
	for i := 0; i < limit; i++ {
		pt := randomPt()
		structIntersected, err := pl.PointLocate(pt.X(), pt.Y())
		bruteForceContains := structIntersected.Contains(pt)
		assert.Nil(t, err)
		assert.True(t, bruteForceContains)
	}
}

func TestRandomDCELSlab(t *testing.T) {
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	structure, err := slab.Decompose(dc, tree.RedBlack)
	assert.Nil(t, err)

	testRandomPts(t, structure, testCt)
}

func TestRandomDCELTrapezoid(t *testing.T) {
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	_, _, structure, err := trapezoid.TrapezoidalMap(dc)
	assert.Nil(t, err)

	testRandomPts(t, structure, testCt)
}

func TestRandomDCELPlumbLine(t *testing.T) {
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	it := bruteForce.PlumbLine(dc)

	testRandomPts(t, it, testCt)
}

func TestRandomDCELRtree(t *testing.T) {
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	pl := rtree.DCELtoRtree(dc)

	testRandomPts(t, pl, testCt)
}

func BenchmarkRandomDCELSlab(b *testing.B) {
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	pl, _ := slab.Decompose(dc, tree.RedBlack)

	for i := 0; i < b.N; i++ {
		pt := randomPt()
		pl.PointLocate(pt.X(), pt.Y())
	}
}

func BenchmarkRandomDCELTrapezoid(b *testing.B) {
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	_, _, pl, _ := trapezoid.TrapezoidalMap(dc)

	for i := 0; i < b.N; i++ {
		pt := randomPt()
		pl.PointLocate(pt.X(), pt.Y())
	}
}

func BenchmarkRandomDCELRtree(b *testing.B) {
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	pl := rtree.DCELtoRtree(dc)

	for i := 0; i < b.N; i++ {
		pt := randomPt()
		pl.PointLocate(pt.X(), pt.Y())
	}
}

func BenchmarkRandomDCELPlumbLine(b *testing.B) {
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	pl := bruteForce.PlumbLine(dc)

	for i := 0; i < b.N; i++ {
		pt := randomPt()
		pl.PointLocate(pt.X(), pt.Y())
	}
}
