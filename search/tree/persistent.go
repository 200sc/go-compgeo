package tree

import (
	"fmt"

	"github.com/200sc/go-compgeo/printutil"
	"github.com/200sc/go-compgeo/search"
)

type PersistentBST struct {
	instant float64
	index   int
	// Implicitly sorted
	instants []BSTInstant
}

type BSTInstant struct {
	*BST
	instant float64
}

func (pbst *PersistentBST) AtInstant(ins float64) search.Dynamic {
	// binary search
	bot := 0
	top := len(pbst.instants) - 1
	var mid int
	fmt.Println("bot,top", bot, top)
	for {
		if top <= bot {
			fmt.Println("Returning index", bot)
			if pbst.instants[bot].instant > ins {
				return pbst.instants[bot-1]
			}
			return pbst.instants[bot]
		}
		mid = (bot + top) / 2
		v := pbst.instants[mid].instant
		if v == ins {
			return pbst.instants[mid]
		} else if v < ins {
			bot = mid + 1
		} else {
			top = mid - 1
		}
	}
}

func (pbst *PersistentBST) ToStaticPersitent() search.StaticPersistent {
	return nil
}
func (pbst *PersistentBST) MinInstant() float64 {
	return pbst.instants[0].instant
}
func (pbst *PersistentBST) MaxInstant() float64 {
	return pbst.instants[len(pbst.instants)-1].instant
}
func (pbst *PersistentBST) SetInstant(ins float64) {
	if ins < pbst.instant {
		panic("Decreasing instants is not yet supported")
	} else if ins == pbst.instant {
		return
	}
	bsti := BSTInstant{}
	bsti.BST = pbst.instants[len(pbst.instants)-1].copy()
	bsti.instant = ins
	pbst.instants = append(pbst.instants, bsti)
	pbst.instant = ins
	pbst.index++
}

func (pbst *PersistentBST) Insert(n search.Node) error {
	return pbst.AtInstant(pbst.instant).Insert(n)
}
func (pbst *PersistentBST) Delete(n search.Node) error {
	return pbst.AtInstant(pbst.instant).Delete(n)
}
func (pbst *PersistentBST) ToStatic() search.Static {
	return pbst.AtInstant(pbst.instant).ToStatic()
}
func (pbst *PersistentBST) Size() int {
	return pbst.AtInstant(pbst.instant).Size()
}
func (pbst *PersistentBST) InOrderTraverse() []search.Node {
	return pbst.AtInstant(pbst.instant).InOrderTraverse()
}
func (pbst *PersistentBST) Search(f interface{}) (bool, interface{}) {
	return pbst.AtInstant(pbst.instant).Search(f)
}
func (pbst *PersistentBST) SearchDown(f interface{}) interface{} {
	return pbst.AtInstant(pbst.instant).SearchDown(f)
}
func (pbst *PersistentBST) SearchUp(f interface{}) interface{} {
	return pbst.AtInstant(pbst.instant).SearchUp(f)
}
func (pbst *PersistentBST) String() string {
	s := ""
	for _, ins := range pbst.instants {
		s += printutil.Stringf64(ins.instant) + ":\n"
		s += ins.BST.String()
	}
	return s
}
