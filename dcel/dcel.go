package dcel

import (
	"fmt"
	"math"

	compgeo "github.com/200sc/go-compgeo"
	"github.com/200sc/go-compgeo/geom"
)

// A DCEL is a structure representin arbitrary plane
// divisions and 3d polytopes. Its values are relatively
// self-explanatory but constructing it is significantly
// harder.
type DCEL struct {
	Vertices  []*Vertex
	HalfEdges []*Edge
	// The first value in a face is the outside component
	// of the face, the second value is the inside component
	Faces []*Face
}

// New returns an empty DCEL with its inner
// fields initialized to empty slices, and a
// zeroth outside face.
func New() *DCEL {
	dc := new(DCEL)
	dc.Vertices = []*Vertex{}
	dc.HalfEdges = []*Edge{}
	dc.Faces = []*Face{NewFace()}
	return dc
}

// MaxX returns the Maximum of all X values
func (dc *DCEL) MaxX() float64 {
	return dc.Max(0)
}

// MaxY returns the Maximum of all Y values
func (dc *DCEL) MaxY() float64 {
	return dc.Max(1)
}

// MaxZ returns the Maximum of all Z values
func (dc *DCEL) MaxZ() float64 {
	return dc.Max(2)
}

// Max functions iterate through vertices
// to find the maximum value along a given axis
// in the DCEL
func (dc *DCEL) Max(i int) (x float64) {
	for _, p := range dc.Vertices {
		if p.Val(i) > x {
			x = p.Val(i)
		}
	}
	return x
}

// MinX returns the Minimum of all X values
func (dc *DCEL) MinX() float64 {
	return dc.Min(0)
}

// MinY returns the Minimum of all Y values
func (dc *DCEL) MinY() float64 {
	return dc.Min(1)
}

// MinZ returns the Minimum of all Z values
func (dc *DCEL) MinZ() float64 {
	return dc.Min(2)
}

// Min functions iterate through vertices
// to find the maximum value along a given axis
// in the DCEL
func (dc *DCEL) Min(i int) (x float64) {
	x = math.Inf(1)
	for _, p := range dc.Vertices {
		if p.Val(i) < x {
			x = p.Val(i)
		}
	}
	return x
}

// ScanFaces returns which index, if any, within dc matches f.
func (dc *DCEL) ScanFaces(f *Face) int {
	for i, f2 := range dc.Faces {
		if f2 == f {
			return i
		}
	}
	return -1
}

// FullEdge returns the ith edge in the form of its
// two vertices
func (dc *DCEL) FullEdge(i int) (geom.FullEdge, error) {
	if i >= len(dc.HalfEdges) {
		return geom.FullEdge{}, compgeo.BadEdgeError{}
	}
	return dc.HalfEdges[i].FullEdge()
}

// FullEdges returns the set of all FullEdges in DCEL.
func (dc *DCEL) FullEdges() ([]geom.FullEdge, [][2]*Face, error) {
	var err error
	fullEdges := make([]geom.FullEdge, len(dc.HalfEdges)/2)
	faces := make([][2]*Face, len(fullEdges))
	for i := 0; i < len(dc.HalfEdges); i += 2 {
		fullEdges[i/2], err = dc.HalfEdges[i].FullEdge()
		if err != nil {
			return nil, nil, err
		}
		faces[i/2] = [2]*Face{dc.HalfEdges[i].Face,
			dc.HalfEdges[i+1].Face}

		// Correct the edges so that the 0th indexed face
		// is below this edge, and vice versa.
		clkwz, _ := faces[i/2][0].Inner.IsClockwise()
		if clkwz && dc.HalfEdges[i].X() > dc.HalfEdges[i+1].X() {
			faces[i/2][0], faces[i/2][1] = faces[i/2][1], faces[i/2][0]
		}
	}
	return fullEdges, faces, nil
}

// CorrectDirectionality (rather innefficently)
// ensures that a face has the right clockwise/
// counter-clockwise orientation based on
// whether its chain is the inner or outer
// portion of a face.
func (dc *DCEL) CorrectDirectionality(f *Face) {
	// Inners need to be going CC
	// Outers need to be going Clockwise

	clock, err := f.Inner.IsClockwise()
	if err == nil && clock {
		f.Inner.Flip()
	} else {
		fmt.Println(err, clock)
	}
	clock, err = f.Outer.IsClockwise()
	if err == nil && !clock {
		f.Outer.Flip()

	}
}

// CorrectTwins modifies the ordering on twins
// inside the DCEL such that dc.HalfEdges[i] is
// the twin of dc.HalfEdges[i+1] for all even
// i values.
func (dc *DCEL) CorrectTwins() {
	newEdges := make([]*Edge, len(dc.HalfEdges))
	seen := make(map[*Edge]bool)
	for i := 0; i < len(dc.HalfEdges); i++ {
		if _, ok := seen[dc.HalfEdges[i]]; !ok {
			newEdges[i] = dc.HalfEdges[i]
			newEdges[i+1] = newEdges[i].Twin
			seen[newEdges[i]] = true
			seen[newEdges[i+1]] = true
		}
	}
	dc.HalfEdges = newEdges
}

// Copy duplicates a DCEL's internal values
// in a new DCEL.
func (dc *DCEL) Copy() *DCEL {
	dc2 := new(DCEL)
	dc2.Faces = make([]*Face, len(dc.Faces))
	dc2.HalfEdges = make([]*Edge, len(dc.HalfEdges))
	dc2.Vertices = make([]*Vertex, len(dc.Vertices))
	ePointerMap := make(map[*Edge]int)
	fPointerMap := make(map[*Face]int)
	vPointerMap := make(map[*Vertex]int)
	for i, e := range dc.HalfEdges {
		ePointerMap[e] = i
		dc2.HalfEdges[i] = NewEdge()
	}
	for i, f := range dc.Faces {
		fPointerMap[f] = i
		f2 := NewFace()
		dc2.Faces[i] = f2
		if f.Outer != nil {
			f2.Outer = dc2.HalfEdges[ePointerMap[f.Outer]]
			f2.Outer.Face = f2
		}
		if f.Inner != nil {
			f2.Inner = dc2.HalfEdges[ePointerMap[f.Inner]]
			f2.Inner.Face = f2
		}
	}
	for i, v := range dc.Vertices {
		vPointerMap[v] = i
		v2 := NewVertex(v.X(), v.Y(), v.Z())
		dc2.Vertices[i] = v2
		if v.OutEdge != nil {
			v2.OutEdge = dc2.HalfEdges[ePointerMap[v.OutEdge]]
		}
	}
	for i, e := range dc.HalfEdges {
		e2 := dc2.HalfEdges[i]
		if e.Prev != nil {
			e2.Prev = dc2.HalfEdges[ePointerMap[e.Prev]]
		}
		if e.Next != nil {
			e2.Next = dc2.HalfEdges[ePointerMap[e.Next]]
		}
		if e.Twin != nil {
			e2.Twin = dc2.HalfEdges[ePointerMap[e.Twin]]
		}
		if e.Origin != nil {
			e2.Origin = dc2.Vertices[vPointerMap[e.Origin]]
		}
	}

	return dc2
}
