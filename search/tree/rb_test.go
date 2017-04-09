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
	test2Input = []testNode{
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
		{11, 10},
		{20, 9},
		{30, 8},
		{40, 7},
		{50, 6},
		{60, 5},
		{70, 4},
		{80, 3},
		{90, 2},
		{100, 1},
		{211, 10},
		{12, 9},
		{13, 8},
		{14, 7},
		{15, 6},
		{16, 5},
		{17, 4},
		{18, 3},
		{19, 2},
		{111, 1},
		{110, 10},
		{120, 9},
		{130, 8},
		{140, 7},
		{150, 6},
		{160, 5},
		{170, 4},
		{180, 3},
		{190, 2},
		{1100, 1},
	}
	notInInput1      = 12
	notInInput2      = 2000
	randomInputCt    = 10000
	randomInputRange = 1000
)

func TestRBInOrder(t *testing.T) {
	tree := New(RedBlack)
	for _, v := range test1Input {
		tree.Insert(v)
	}
	inOrder := tree.InOrderTraverse()
	expected := [...]float64{
		10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	for i := range inOrder {
		//fmt.Println(inOrder[i].Key(), inOrder[i].Val())
		assert.Equal(t, expected[i], inOrder[i].Key())
	}
}

func TestRBDefinedInput1(t *testing.T) {
	tree := New(RedBlack)
	for _, v := range test1Input {
		tree.Insert(v)
	}

	valid, err := RBValid(tree.(*BST))
	assert.True(t, valid)
	assert.Nil(t, err)
	valid2 := tree.(*BST).isValid()
	assert.True(t, valid2)
	t.Log(tree.(*BST).root)

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
		//t.Log(tree.(*BST).root)
		assert.True(t, valid)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		valid2 := tree.(*BST).isValid()
		assert.True(t, valid2)
	}
}

func TestPredSucc(t *testing.T) {
	tree := New(RedBlack)
	for _, v := range test1Input {
		tree.Insert(v)
	}

	v := tree.SearchUp(9.5)
	assert.Equal(t, v, float64(1))
	v = tree.SearchDown(9.5)
	assert.Equal(t, v, float64(2))

	t2 := tree.ToStatic()

	v = t2.SearchUp(9.5)
	assert.Equal(t, v, float64(1))
	v = t2.SearchDown(9.5)
	assert.Equal(t, v, float64(2))

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
