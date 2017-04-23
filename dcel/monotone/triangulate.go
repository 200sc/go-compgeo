package monotone

import (
	"errors"

	"github.com/200sc/go-compgeo/dcel"
)

// Triangulate uses Monotonization to convert a dcel into
// a dcel made up of just triangles
func Triangulate(inDc *dcel.DCEL) (*dcel.DCEL, map[*dcel.Face]*dcel.Face, error) {
	monotonized, faceMap, err := Split(inDc)
	if err != nil {
		return monotonized, faceMap, err
	}
	return TriangulateSplit(monotonized, faceMap)
}

// TriangulateSplit takes in a dcel whose faces are already
// monotone. If there is no existing faceMap, it will
// create its own.
func TriangulateSplit(monotonized *dcel.DCEL, faceMap map[*dcel.Face]*dcel.Face) (*dcel.DCEL, map[*dcel.Face]*dcel.Face, error) {
	edgeLen := len(monotonized.HalfEdges)
	if faceMap == nil {
		faceMap = make(map[*dcel.Face]*dcel.Face)
	}
	// Triangulate each monotone polygon
	for _, f := range monotonized.Faces {
		ypts := f.VerticesSorted(1, 0)
		chainMap, err := Chains(monotonized, f, ypts)
		if err != nil {
			return monotonized, faceMap, err
		}
		plTree, err := NewDoubleIntervalTree(f, monotonized)
		if err != nil {
			return monotonized, faceMap, err
		}
		stack := VertexStack{}
		stack.Push(ypts[0], ypts[1])
		var v *dcel.Vertex
		for i := 2; i < len(ypts)-1; i++ {
			if chainMap[ypts[i]] != chainMap[stack.first.Vertex] {
				v = stack.Pop()
				for !stack.IsEmpty() {
					monotonized.ConnectVerts(v, ypts[i], f)
					v = stack.Pop()
				}
				stack.Push(ypts[i-1], ypts[i])
			} else {
				stack.Pop()
				for {
					v = stack.Pop()
					if DiagonalWithinFace(plTree, v, ypts[i]) {
						monotonized.ConnectVerts(v, ypts[i], f)
					} else {
						break
					}
				}
				stack.Push(v, ypts[i])
			}
		}
		// Walk each new edge to create new faces.
		for i := edgeLen; i < len(monotonized.HalfEdges); i++ {
			e := monotonized.HalfEdges[i]
			if e.Face == f {
				newFace := dcel.NewFace()
				newFace.Outer = e
				e.Face = newFace
				for e = e.Next; e != newFace.Outer; e = e.Next {
					e.Face = newFace
				}
				faceMap[newFace] = f
			} // else we've already walked this edge
		}
		edgeLen = len(monotonized.HalfEdges)
	}
	return monotonized, faceMap, nil
}

type chain bool

const (
	aChain chain = true
	bChain chain = false
)

func Chains(dc *dcel.DCEL, f *dcel.Face, pts []*dcel.Vertex) (map[*dcel.Vertex]chain, error) {
	// find the single start and end vertices
	var start, end *dcel.Vertex
	m := make(map[*dcel.Vertex]chain)
	for _, p := range pts {
		typ := VertexType(p, dc)
		if typ == START {
			if start != nil {
				return m, errors.New("A face on the input DCEL was not monotone")
			}
			start = p
		} else if typ == END {
			if end != nil {
				return m, errors.New("A face on the input DCEL was not monotone")
			}
			end = p
		}
	}

	e := f.Outer
	for {
		if e.Origin == start {
			break
		}
		e = e.Next
	}

	// We'll say that start falls into the chain a.
	for e.Origin != end {
		m[e.Origin] = aChain
		// We don't care about the directionaliy
		// of the face, we just need to distinguish
		// the left and right chains. Whichever of
		// left and right ends up as chain A is
		// insignificant.
		e = e.Next
	}

	return m, nil
}

func DiagonalWithinFace(tree dcel.LocatesPoints, a, b *dcel.Vertex) bool {
	mid := a.Mid2D(b)
	f, _ := tree.PointLocate(mid.X(), mid.Y())
	return f != nil
}

type VertexStackItem struct {
	*dcel.Vertex
	next, prev *VertexStackItem
}

type VertexStack struct {
	first, last *VertexStackItem
}

func (vst *VertexStack) IsEmpty() bool {
	return vst.first == nil
}

func (vst *VertexStack) Push(vs ...*dcel.Vertex) {
	for _, v := range vs {
		item := &VertexStackItem{Vertex: v}
		if vst.last == nil {
			vst.first = item
			vst.last = item
			return
		}
		vst.last.next = item
		item.prev = vst.last
		vst.last = item
	}
}

func (vst *VertexStack) Pop() *dcel.Vertex {
	if vst.first == nil {
		return nil
	}
	v := vst.first.Vertex
	vst.first = vst.first.next
	return v
}
