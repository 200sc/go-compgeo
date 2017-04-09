package tree

import (
	"math/rand"
	"testing"
)

var (
	j int
)

func BenchmarkRBDynamic1(b *testing.B) {
	benchmarkRBDynamic(b, test1Input, notInInput1)
}
func BenchmarkRBStatic1(b *testing.B) {
	benchmarkRBStatic(b, test1Input, notInInput1)
}
func BenchmarkMap1(b *testing.B) {
	benchmarkMap(b, test1Input, notInInput1)
}
func BenchmarkRBDynamic2(b *testing.B) {
	benchmarkRBDynamic(b, test2Input, notInInput2)
}
func BenchmarkRBStatic2(b *testing.B) {
	benchmarkRBStatic(b, test2Input, notInInput2)
}
func BenchmarkMap2(b *testing.B) {
	benchmarkMap(b, test2Input, notInInput2)
}
func BenchmarkRBDynamic3(b *testing.B) {
	randomInput := randomInput()
	benchmarkRBDynamic(b, randomInput, randomInputRange+1)
}
func BenchmarkRBStatic3(b *testing.B) {
	randomInput := randomInput()
	benchmarkRBStatic(b, randomInput, randomInputRange+1)
}
func BenchmarkMap3(b *testing.B) {
	randomInput := randomInput()
	benchmarkMap(b, randomInput, randomInputRange+1)
}
func BenchmarkRBDynamic4(b *testing.B) {
	randomInput := randomInputNoDupes()
	benchmarkRBDynamic(b, randomInput, randomInputRange+1)
}
func BenchmarkRBStatic4(b *testing.B) {
	randomInput := randomInputNoDupes()
	benchmarkRBStatic(b, randomInput, randomInputRange+1)
}
func BenchmarkMap4(b *testing.B) {
	randomInput := randomInputNoDupes()
	benchmarkMap(b, randomInput, randomInputRange+1)
}

func randomInput() []testNode {
	randomInput := make([]testNode, randomInputCt)
	for i := range randomInput {
		randomInput[i] = testNode{
			float64(rand.Intn(randomInputRange)),
			float64(rand.Intn(randomInputRange)),
		}
	}
	return randomInput
}

func randomInputNoDupes() []testNode {
	randomInput := make([]testNode, randomInputCt)
	for i := range randomInput {
		randomInput[i] = testNode{
			float64(i),
			float64(i),
		}
	}
	return randomInput
}

func benchmarkRBDynamic(b *testing.B, input []testNode, inputLimit int) {
	tree := New(RedBlack)
	for _, v := range input {
		tree.Insert(v)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b, _ := tree.Search(float64(rand.Intn(inputLimit)))
		// We do this to be fair to maps
		if b {
			j++
		}
	}
}

func benchmarkRBStatic(b *testing.B, input []testNode, inputLimit int) {
	tree := New(RedBlack)
	for _, v := range input {
		tree.Insert(v)
	}
	t2 := tree.ToStatic()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b, _ := t2.Search(float64(rand.Intn(inputLimit)))
		// We do this to be fair to maps
		if b {
			j++
		}
	}
}

func benchmarkMap(b *testing.B, input []testNode, inputLimit int) {
	m := make(map[float64]map[float64]bool)
	for _, v := range input {
		if _, ok := m[v.key]; !ok {
			m[v.key] = make(map[float64]bool)
		}
		m[v.key][v.val] = true
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := m[float64(rand.Intn(inputLimit))]
		// The Go compiler won't let m[...] exist
		// by itself, so we need to do something
		// with its output
		if k != nil {
			j++
		}
	}
}
