package tree

var (
	avlFnSet = &fnSet{
		insertFn: avlInsert,
		deleteFn: avlDelete,
		searchFn: nopNode,
	}
)

func (n *node) balance() int {
	if n == nil {
		return 0
	}
	return n.payload.(int)
}

func avlRotateRL(a, b *node) *node { return nil }
func avlRotateLR(a, b *node) *node { return nil }
func avlRotateR(a, b *node) *node  { return nil }
func avlRotateL(a, b *node) *node  { return nil }

func avlInsert(n *node) *node {
	var g, p, s *node
	for {
		p = n.parent
		if p == nil {
			break
		}

		if n == p.right {
			if p.balance() > 0 {
				g = p.parent
				if n.balance() < 0 {
					s = avlRotateRL(p, n)
				} else {
					s = avlRotateL(p, n)
				}
			} else {
				if p.balance() < 0 {
					p.payload = 0
					break
				}
				p.payload = 1
				n = p
				continue
			}
		} else {
			if p.balance() < 0 {
				g = p.parent
				if n.balance() > 0 {
					s = avlRotateLR(p, n)
				} else {
					s = avlRotateR(p, n)
				}
			} else {
				if p.balance() > 0 {
					p.payload = 0
					break
				}
				p.payload = -1
				n = p
				continue
			}
		}

		s.parent = g
		if g != nil {
			if p == g.left {
				g.left = s
			} else {
				g.right = s
			}
			break
		} else {
			return s
		}
	}
	return nil

}
func avlDelete(n *node) *node {
	//Todo
	return nil
}
