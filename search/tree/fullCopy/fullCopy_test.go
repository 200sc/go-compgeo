package fullCopy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	instantInputs1 = [][]testNode{
		{{1, 1}},
		{{2, 1}},
		{{3, 1}},
		{{4, 1}},
		{{5, 1}},
		{{6, 1}},
		{{7, 1}},
		{{8, 1}},
		{{9, 1}},
		{{10, 1}},
	}
)

func TestPBSTDefinedInput1(t *testing.T) {
	tree := New(RedBlack).ToPersistent()
	for i, ls := range instantInputs1 {
		tree.SetInstant(float64(i))
		for _, v := range ls {
			err := tree.Insert(v)
			assert.Nil(t, err)
		}
	}
	for i := range instantInputs1 {
		t2 := tree.AtInstant(float64(i))
		for j, ls2 := range instantInputs1 {
			if j > i {
				break
			}
			for _, v := range ls2 {
				found, _ := t2.Search(v.key)
				assert.True(t, found)
			}
		}
	}
	for i, ls := range instantInputs1 {
		tree.SetInstant(float64(len(instantInputs1) + i))
		for _, v := range ls {
			err := tree.Delete(v)
			assert.Nil(t, err)
		}
	}
	for i := range instantInputs1 {
		t2 := tree.AtInstant(float64(len(instantInputs1) + i))
		for j, ls2 := range instantInputs1 {
			if j > i {
				break
			}
			for _, v := range ls2 {
				found, _ := t2.Search(v.key)
				assert.False(t, found)
			}
		}
	}
}
