package tree

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nilValNode struct {
	key float64
}

func (n nilValNode) Key() float64 {
	return n.key
}

func (n nilValNode) Val() interface{} {
	return nil
}

type testNode struct {
	key, val float64
}

func (t testNode) Key() float64 {
	return t.key
}

func (t testNode) Val() interface{} {
	return t.val
}

var (
	test1Input = []testNode{
		{1, 10},
		{2, 9},
		{3, 8},
		{4, 7},
		{5, 6},
		{6, 5},
		{7, 4},
		{8, 3},
		{9, 2},
		{10, 1},
	}
	notInInput1      = 12
	randomInputCt    = 20000
	randomInputRange = 1000
)

func TestRBDefinedInput1(t *testing.T) {
	tree := New(RedBlack)
	for _, v := range test1Input {
		tree.Insert(v)
	}

	valid, err := RBValid(tree.(*BST))
	assert.True(t, valid)
	assert.Nil(t, err)

	// Should be in tree
	for _, v := range test1Input {
		b, found := tree.Search(v.key)
		assert.True(t, b)
		assert.Equal(t, found, v.val)
	}
	// Should not be in tree
	for i := notInInput1; i < notInInput1+10; i++ {
		b, found := tree.Search(float64(i))
		assert.False(t, b)
		assert.Nil(t, found)
	}

	for _, v := range test1Input {
		tree.Delete(v)
		b, found := tree.Search(v.key)
		assert.False(t, b)
		assert.Nil(t, found)
		valid, err := RBValid(tree.(*BST))
		t.Log(tree.(*BST).root)
		assert.True(t, valid)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
	}
}

func TestRBRandomInput(t *testing.T) {
	tree := New(RedBlack)
	t.Log("Uh")
	for i := 0; i < randomInputCt; i++ {
		n := testNode{
			float64(rand.Intn(randomInputRange)),
			float64(rand.Intn(randomInputRange)),
		}
		t.Log("Inserting", n)
		tree.Insert(n)
		valid, err := RBValid(tree.(*BST))
		assert.True(t, valid)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
	}
	t.Log("Insert Complete")
	// These values might not be in the bst.
	for i := 0; i < randomInputCt; i++ {
		n := nilValNode{float64(rand.Intn(randomInputRange))}

		t.Log("Deleting", n)
		tree.Delete(n)
		valid, err := RBValid(tree.(*BST))
		assert.True(t, valid)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
	}
}

func TestRBToStatic(t *testing.T) {
	tree := New(RedBlack)
	inserted := make(map[float64]bool)
	for i := 0; i < randomInputCt; i++ {
		n := testNode{
			float64(rand.Intn(randomInputRange)),
			float64(rand.Intn(randomInputRange)),
		}
		inserted[n.key] = true
		tree.Insert(n)
		// We don't check that the tree is valid, that's
		// another test's job.
	}
	t2 := tree.ToStatic()
	for i := 0; i < randomInputCt; i++ {
		key := float64(rand.Intn(randomInputRange))
		b, _ := t2.Search(key)
		assert.Equal(t, b, inserted[key])
	}
}

func BenchmarkRBDynamic(b *testing.B) {
	tree := New(RedBlack)
	for _, v := range test1Input {
		tree.Insert(v)
	}
	j := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b, _ := tree.Search(float64(rand.Intn(notInInput1)))
		// We do this to be fair to maps
		if b {
			j++
		}
	}
}

func BenchmarkRBStatic(b *testing.B) {
	tree := New(RedBlack)
	for _, v := range test1Input {
		tree.Insert(v)
	}
	t2 := tree.ToStatic()
	j := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b, _ := t2.Search(float64(rand.Intn(notInInput1)))
		// We do this to be fair to maps
		if b {
			j++
		}
	}
}

func BenchmarkMap(b *testing.B) {
	m := make(map[float64]map[float64]bool)
	for _, v := range test1Input {
		if _, ok := m[v.key]; !ok {
			m[v.key] = make(map[float64]bool)
		}
		m[v.key][v.val] = true
	}
	j := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := m[float64(rand.Intn(notInInput1))]
		// The Go compiler won't let m[...] exist
		// by itself, so we need to do something
		// with its output
		if k != nil {
			j++
		}
	}
}
