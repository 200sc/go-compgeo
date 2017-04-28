package tree

var (
	SplayFnSet = &FnSet{
		InsertFn: splay,
		DeleteFn: splayDelete,
		SearchFn: splay,
	}
)

func splay(n *node) *node {
	for n.parent != nil {
		if n.parent.parent == nil {
			if n.parent.left == n {
				n.parent.rightRotate()
			} else {
				n.parent.leftRotate()
			}
		} else {
			if n.parent.left == n {
				if n.parent.parent.left == n.parent {
					n.parent.parent.rightRotate()
					n.parent.rightRotate()
				} else {
					n.parent.rightRotate()
					n.parent.leftRotate()
				}
			} else {
				if n.parent.parent.left == n.parent {
					n.parent.leftRotate()
					n.parent.rightRotate()
				} else {
					n.parent.parent.leftRotate()
					n.parent.leftRotate()
				}
			}
		}
	}
	return n
}

func splayDelete(n *node) *node {

	// splay(n)
	// if n.left != nil {
	// 	n.right.parent = nil
	// 	return n.right
	// } else if n.right != nil {
	// 	n.left.parent = nil
	// 	return n.left
	// } else {

	// }
	return nil
}
