package tree

import (
	"fmt"
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
	randomInputCt    = 50
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
		fmt.Println(tree.(*BST).root)
		assert.True(t, valid)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
	}
}

func TestRBRandomInput(t *testing.T) {
	tree := New(RedBlack)
	inserted := make(map[float64]bool)
	for i := 0; i < randomInputCt; i++ {
		n := testNode{
			float64(rand.Intn(randomInputRange)),
			float64(rand.Intn(randomInputRange)),
		}
		tree.Insert(n)
		inserted[n.key] = true
		findCycle(tree.(*BST))
		valid, err := RBValid(tree.(*BST))
		assert.True(t, valid)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
	}
	// These values might not be in the bst.
	for i := 0; i < randomInputCt; i++ {
		n := nilValNode{float64(rand.Intn(randomInputRange))}
		tree.Delete(n)
		findCycle(tree.(*BST))
		valid, err := RBValid(tree.(*BST))
		assert.True(t, valid)
		assert.Nil(t, err)
	}
}
