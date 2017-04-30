package test

import (
	"math/rand"
	"testing"

	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/pointLoc"
	"github.com/200sc/go-compgeo/dcel/pointLoc/slab"
	"github.com/200sc/go-compgeo/dcel/pointLoc/trapezoid"
	"github.com/200sc/go-compgeo/geom"
	"github.com/200sc/go-compgeo/search/tree"
	"github.com/stretchr/testify/assert"
)

var (
	inputSize  = 5
	inputRange = 10000.0
	testCt     = 1000
)

func randomPt() geom.D3 {
	return geom.NewPoint(rand.Float64()*inputRange,
		rand.Float64()*inputRange, 0)
}

func testRandomDCEL(t *testing.T, pl pointLoc.LocatesPoints) {
	for i := 0; i < testCt; i++ {
		pt := randomPt()
		structIntersected, err := pl.PointLocate(pt.X(), pt.Y())
		assert.Nil(t, err)
		assert.True(t, structIntersected.Contains(pt))
	}
}

func TestRandomDCELSlab(t *testing.T) {
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	structure, err := slab.Decompose(dc, tree.RedBlack)
	assert.Nil(t, err)

	testRandomDCEL(t, structure)
}

func TestRandomDCELTrapezoid(t *testing.T) {
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	_, _, structure, err := trapezoid.TrapezoidalMap(dc)
	assert.Nil(t, err)

	testRandomDCEL(t, structure)
}
