package compgeo

type AVLNode struct {
	key int
	val interface{}
}

type AVLTree []AVLNode

func NewAVLTree() *AVLTree {
	t := new(AVLTree)
	// ...
	return t
}

func (avlt *AVLTree) Search(key int) (bool, interface{}) {
	t := *avlt
	i := 0
	for {
		n := t[i]
		if n.key == key {
			return true, n.val
		} else if n.key < key {
			i = 2 * i
		} else {
			i = (2 * i) + 1
		}
		if i >= len(t) {
			return false, nil
		}
	}
}
