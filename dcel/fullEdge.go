package dcel

// FullEdge is a renaming of the output of
// DCEL.FullEdge(i) or Edge.FullEdge, so
// associated functions can be attached to the construct.
type FullEdge [2]*Vertex

func (fe FullEdge) Lesser(d int) *Vertex {
	if fe[0].Val(d) < fe[1].Val(d) {
		return fe[0]
	}
	return fe[1]
}

func (fe FullEdge) Greater(d int) *Vertex {
	if fe[0].Val(d) > fe[1].Val(d) {
		return fe[0]
	}
	return fe[1]
}

func (fe FullEdge) Left() *Vertex {
	return fe.Lesser(0)
}

func (fe FullEdge) Right() *Vertex {
	return fe.Greater(0)
}

func (fe FullEdge) Top() *Vertex {
	return fe.Greater(1)
}

func (fe FullEdge) Bottom() *Vertex {
	return fe.Lesser(1)
}
