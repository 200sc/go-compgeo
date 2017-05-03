package fullCopy

import (
	"math"

	"fmt"

	"github.com/200sc/go-compgeo/geom"
	"github.com/200sc/go-compgeo/printutil"
	"github.com/200sc/go-compgeo/search"
)

// FullPersistentBST is an implementation of a persistent
// binary search tree using full copies, with each
// instant represented by a separate BST.
type FullPersistentBST struct {
	instant float64
	index   int
	// Implicitly sorted
	instants []BSTInstant
}

func NewFullPersistentBST(dyn search.Dynamic) search.DynamicPersistent {
	pbst := new(FullPersistentBST)
	pbst.instant = math.MaxFloat64 * -1
	pbst.instants = []BSTInstant{{Dynamic: dyn, instant: pbst.instant}}
	return pbst
}

// BSTInstant is a single BST within a Persistent BST.
type BSTInstant struct {
	search.Dynamic
	instant float64
}

// ThisInstant returns the subtree at the most recent
// instant set
func (pbst *FullPersistentBST) ThisInstant() search.Dynamic {
	return pbst.instants[pbst.index]
}

// AtInstant returns the subtree of pbst at the given instant
func (pbst *FullPersistentBST) AtInstant(ins float64) search.Dynamic {
	// binary search
	bot := 0
	top := len(pbst.instants) - 1
	var mid int
	for {
		if top <= bot {
			// round down
			if pbst.instants[bot].instant > ins {
				bot--
			}
			return pbst.instants[bot]
		}
		mid = (bot + top) / 2
		v := pbst.instants[mid].instant
		if geom.F64eq(v, ins) {
			return pbst.instants[mid]
		} else if v < ins {
			bot = mid + 1
		} else {
			top = mid - 1
		}
	}
}

// ToStaticPersistent returns a static peristent version
// of the pbst
func (pbst *FullPersistentBST) ToStaticPersistent() search.StaticPersistent {
	// Todo
	return nil
}

// MinInstant returns the minimum instant ever set on pbst.
func (pbst *FullPersistentBST) MinInstant() float64 {
	return pbst.instants[0].instant
}

// MaxInstant returns the maximum instant ever set on pbst.
func (pbst *FullPersistentBST) MaxInstant() float64 {
	return pbst.instants[len(pbst.instants)-1].instant
}

// SetInstant increments the pbst to the given instant.
func (pbst *FullPersistentBST) SetInstant(ins float64) {
	if ins < pbst.instant {
		panic("Decreasing instants is not yet supported")
	} else if ins == pbst.instant {
		return
	}
	bsti := BSTInstant{}
	bsti.Dynamic = pbst.instants[len(pbst.instants)-1].Copy().(search.Dynamic)
	bsti.instant = ins
	pbst.instants = append(pbst.instants, bsti)
	pbst.instant = ins
	pbst.index++
}

// Insert peforms Insert on the current set instant's search tree.
func (pbst *FullPersistentBST) Insert(n search.Node) error {
	return pbst.AtInstant(pbst.instant).Insert(n)
}

// Delete performs Delete on the current set instant's search tree.
func (pbst *FullPersistentBST) Delete(n search.Node) error {
	return pbst.AtInstant(pbst.instant).Delete(n)
}

// ToStatic performs ToStatic on the current set instant's search tree.
func (pbst *FullPersistentBST) ToStatic() search.Static {
	return pbst.AtInstant(pbst.instant).ToStatic()
}

// Size performs Size on the current set instant's search tree.
func (pbst *FullPersistentBST) Size() int {
	return pbst.AtInstant(pbst.instant).Size()
}

// InOrderTraverse performs InOrderTraverse on the current
// set instant's search tree.
func (pbst *FullPersistentBST) InOrderTraverse() []search.Node {
	return pbst.AtInstant(pbst.instant).InOrderTraverse()
}

// Search performs Search on the current set instant's search tree.
func (pbst *FullPersistentBST) Search(f interface{}) (bool, interface{}) {
	return pbst.AtInstant(pbst.instant).Search(f)
}

// SearchDown performs SearchDown on the current set instant's search tree.
func (pbst *FullPersistentBST) SearchDown(f interface{}, d int) (search.Comparable, interface{}) {
	return pbst.AtInstant(pbst.instant).SearchDown(f, d)
}

// SearchUp performs SearchUp on the current set instant's search tree.
func (pbst *FullPersistentBST) SearchUp(f interface{}, u int) (search.Comparable, interface{}) {
	return pbst.AtInstant(pbst.instant).SearchUp(f, u)
}

// String returns a string representation of pbst.
func (pbst *FullPersistentBST) String() string {
	s := ""
	for _, ins := range pbst.instants {
		s += printutil.Stringf64(ins.instant) + ":\n"
		s += fmt.Sprintf("%v", ins.Dynamic)
	}
	return s
}

func (pbst *FullPersistentBST) Copy() interface{} {
	return nil
}
