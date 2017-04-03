package compgeo

type staticNode struct {
	// eventually key should be a comparable interface
	// but that would probably poorly effect performance
	key float64
	val interface{}
}

// This implicitly says that
// a user cannot store nils in
// this tree. This is probably
// overly limiting.
func (sn staticNode) IsNil() bool {
	return sn.val == nil
}

func ancestor(i, tiersUp int) int {
	return i / (tiersUp * 2)
}

func parent(i int) int {
	return i / 2
}

func left(i int) int {
	return 2 * i
}

func right(i int) int {
	return (2 * i) + 1
}

func isLeftChild(i int) bool {
	return i%2 == 0
}

type staticSearchTree []staticNode

func (sst staticSearchTree) IsNil(i int) bool {
	if i < len(sst) {
		return sst[i].IsNil()
	}
	return true
}

// Search :
func (sst staticSearchTree) Search(key float64) (bool, interface{}) {
	i := 1
	for {
		n := sst[i]
		k := n.key
		if k == key {
			return true, n.val
		} else if k < key {
			i = left(i)
		} else {
			i = right(i)
		}
		if sst.IsNil(i) {
			return false, nil
		}
	}
}

func (sst staticSearchTree) minKey(i int) int {
	for sst.IsNil(left(i)) {
		i = left(i)
	}
	return i
}

// Traverse :
// There are multiple ways to traverse a tree.
// The most useful of these is the in-order traverse,
// and that's what we provide here.
// Other traversal methods can be added as needed.
func (sst staticSearchTree) InOrderTraverse() []staticNode {
	out := make([]staticNode, len(sst))
	i := 0
	sst.inOrderTraverse(out, 0, &i)
	return out
}

func (sst staticSearchTree) inOrderTraverse(out []staticNode, i int, nextIndex *int) {
	if !sst.IsNil(i) {
		sst.inOrderTraverse(out, left(i), nextIndex)
		out[*nextIndex] = sst[i]
		*nextIndex++
		sst.inOrderTraverse(out, right(i), nextIndex)
	}
}
