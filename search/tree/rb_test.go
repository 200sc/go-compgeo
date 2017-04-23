package tree

import (
	"math/rand"
	"testing"

	"github.com/200sc/go-compgeo/printutil"
	"github.com/200sc/go-compgeo/search"
	"github.com/stretchr/testify/assert"
)

type nilValNode struct {
	key compFloat
}

func (n nilValNode) Key() search.Comparable {
	return n.key
}

func (n nilValNode) Val() search.Equalable {
	return search.Nil{}
}

type compFloat float64

func (f compFloat) Compare(i interface{}) search.CompareResult {
	var f3 compFloat
	switch f2 := i.(type) {
	case float64:
		f3 = compFloat(f2)
	case compFloat:
		f3 = f2
	default:
		return search.Invalid
	}
	if f == f3 {
		return search.Equal
	} else if f < f3 {
		return search.Less
	}
	return search.Greater
}

func (f compFloat) String() string {
	return printutil.Stringf64(float64(f))
}

func (f compFloat) Equals(e search.Equalable) bool {
	switch f2 := e.(type) {
	case compFloat:
		return f == f2
	}
	return false
}

type testNode struct {
	key compFloat
	val compFloat
}

func (t testNode) Key() search.Comparable {
	return t.key
}

func (t testNode) Val() search.Equalable {
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
	randomInputCt    = 20000
	randomInputRange = 5000
)

func TestRBInOrder(t *testing.T) {
	tree := New(RedBlack)
	for _, v := range test1Input {
		tree.Insert(v)
	}
	inOrder := tree.InOrderTraverse()
	expected := [...]compFloat{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := range inOrder {
		assert.Equal(t, expected[i], inOrder[i].Key())
	}
}

func TestRBDefinedInput1(t *testing.T) {
	tree := New(RedBlack)
	for _, v := range test1Input {
		tree.Insert(v)
		assert.Equal(t, tree.Size(), tree.(*BST).calcSize())
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
		err = tree.Delete(v)
		assert.Nil(t, err)
		assert.Equal(t, tree.Size(), tree.(*BST).calcSize())
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

	_, v := tree.SearchUp(9.5)
	assert.Equal(t, v, compFloat(1))
	_, v = tree.SearchDown(9.5)
	assert.Equal(t, v, compFloat(2))

	t2 := tree.ToStatic()

	_, v = t2.SearchUp(9.5)
	assert.Equal(t, v, compFloat(1))
	_, v = t2.SearchDown(9.5)
	assert.Equal(t, v, compFloat(2))

}

func TestRBRandomInput(t *testing.T) {
	tree := New(RedBlack)
	for i := 0; i < randomInputCt; i++ {
		n := testNode{
			compFloat(float64(rand.Intn(randomInputRange))),
			compFloat(float64(rand.Intn(randomInputRange))),
		}
		t.Log("Inserting", n)
		tree.Insert(n)
		valid, err := RBValid(tree.(*BST))
		assert.True(t, valid)
		assert.Equal(t, tree.Size(), tree.(*BST).calcSize())
		if !assert.Nil(t, err) {
			t.FailNow()
		}
	}
	t.Log("Insert Complete")
	totalSize := tree.Size()
	// These values might not be in the bst.
	for i := 0; i < randomInputCt; i++ {
		n := nilValNode{compFloat(float64(rand.Intn(randomInputRange)))}

		t.Log("Deleting", n)
		err := tree.Delete(n)
		if err == nil {
			totalSize--
		}
		valid, err := RBValid(tree.(*BST))
		assert.Equal(t, totalSize, tree.Size())
		assert.Equal(t, tree.Size(), tree.(*BST).calcSize())
		assert.True(t, valid)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
	}
}

func TestRBToStatic(t *testing.T) {
	tree := New(RedBlack)
	inserted := make(map[compFloat]bool)
	for i := 0; i < randomInputCt; i++ {
		n := testNode{
			compFloat(float64(rand.Intn(randomInputRange))),
			compFloat(float64(rand.Intn(randomInputRange))),
		}
		inserted[n.key] = true
		tree.Insert(n)
		// We don't check that the tree is valid, that's
		// another test's job.
	}
	t2 := tree.ToStatic()
	for i := 0; i < randomInputCt; i++ {
		key := compFloat(float64(rand.Intn(randomInputRange)))
		b, _ := t2.Search(key)
		assert.Equal(t, b, inserted[key])
	}
}
