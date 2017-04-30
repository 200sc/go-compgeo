package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

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
	inputSize   = 3
	inputRange  = 10000.0
	testCt      = 10000
	slabErrors  = 0
	trapErrors  = 0
	rtreeErrors = 0
	plumbErrors = 0
	seed        int64
)

func randomPt() geom.D3 {
	return geom.NewPoint(rand.Float64()*inputRange,
		rand.Float64()*inputRange, 0)
}

func testRandomPts(t *testing.T, pl pointLoc.LocatesPoints, limit int, errs *int) {
	for i := 0; i < limit; i++ {
		pt := randomPt()
		structIntersected, err := pl.PointLocate(pt.X(), pt.Y())
		bruteForceContains := structIntersected.Contains(pt)
		assert.Nil(t, err)
		if !assert.True(t, bruteForceContains) {
			(*errs)++
		}
	}
}

func TestRandomDCELSlab(t *testing.T) {
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	structure, err := slab.Decompose(dc, tree.RedBlack)
	assert.Nil(t, err)

	testRandomPts(t, structure, testCt, &slabErrors)
	printErrors()
}

func TestRandomDCELTrapezoid(t *testing.T) {
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	_, _, structure, err := trapezoid.TrapezoidalMap(dc)
	assert.Nil(t, err)

	testRandomPts(t, structure, testCt, &trapErrors)
	printErrors()
}

func TestRandomDCELPlumbLine(t *testing.T) {
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	it := bruteForce.PlumbLine(dc)

	testRandomPts(t, it, testCt, &plumbErrors)
	printErrors()
}

func TestRandomDCELRtree(t *testing.T) {
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	pl := rtree.DCELtoRtree(dc)

	testRandomPts(t, pl, testCt, &rtreeErrors)
	printErrors()
}

func printErrors() {
	fmt.Println("Total errors")
	fmt.Println("Slab:", slabErrors)
	fmt.Println("Trapezoid:", trapErrors)
	fmt.Println("Plumb Line:", plumbErrors, "(Should be zero)")
	fmt.Println("Rtree:", rtreeErrors)
	fmt.Println()
}

func BenchmarkRandomDCELSlab(b *testing.B) {
	if seed == 0 {
		fmt.Println("Setting seed")
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	pl, _ := slab.Decompose(dc, tree.RedBlack)

	rand.Seed(seed)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt := randomPt()
		pl.PointLocate(pt.X(), pt.Y())
	}
}

func BenchmarkRandomDCELTrapezoid(b *testing.B) {
	// This seed pattern guarantees that
	// each benchmark is run with the same
	// input dcel and input points
	if seed == 0 {
		fmt.Println("Setting seed")
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	_, _, pl, _ := trapezoid.TrapezoidalMap(dc)

	rand.Seed(seed)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt := randomPt()
		pl.PointLocate(pt.X(), pt.Y())
	}
}

func BenchmarkRandomDCELRtree(b *testing.B) {
	if seed == 0 {
		fmt.Println("Setting seed")
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	pl := rtree.DCELtoRtree(dc)

	rand.Seed(seed)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt := randomPt()
		pl.PointLocate(pt.X(), pt.Y())
	}
}

func BenchmarkRandomDCELPlumbLine(b *testing.B) {
	if seed == 0 {
		fmt.Println("Setting seed")
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
	dc := dcel.Random2DDCEL(inputRange, inputSize)

	pl := bruteForce.PlumbLine(dc)

	rand.Seed(seed)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt := randomPt()
		pl.PointLocate(pt.X(), pt.Y())
	}
}

func BenchmarkRandomSetupSlab(b *testing.B) {
	if seed == 0 {
		fmt.Println("Setting seed")
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
	for i := 0; i < b.N; i++ {
		dc := dcel.Random2DDCEL(inputRange, inputSize)
		slab.Decompose(dc, tree.RedBlack)
	}
}

func BenchmarkRandomSetupTrapezoid(b *testing.B) {
	if seed == 0 {
		fmt.Println("Setting seed")
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
	for i := 0; i < b.N; i++ {
		dc := dcel.Random2DDCEL(inputRange, inputSize)
		trapezoid.TrapezoidalMap(dc)
	}
}

func BenchmarkRandomSetupRtree(b *testing.B) {
	if seed == 0 {
		fmt.Println("Setting seed")
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
	for i := 0; i < b.N; i++ {
		dc := dcel.Random2DDCEL(inputRange, inputSize)
		rtree.DCELtoRtree(dc)
	}
}

func BenchmarkRandomSetupPlumbLine(b *testing.B) {
	if seed == 0 {
		fmt.Println("Setting seed")
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
	for i := 0; i < b.N; i++ {
		dc := dcel.Random2DDCEL(inputRange, inputSize)
		bruteForce.PlumbLine(dc)
	}
}
