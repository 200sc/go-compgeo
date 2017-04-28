package dcel

import (
	"fmt"
	"math"
	"sort"
	"strconv"

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

func (dc *DCEL) String() string {
	s := "DCEL\n"
	s += "-----\n"
	s += "Faces\n"
	s += "-----\n"
	for i, f := range dc.Faces {
		s += "f" + strconv.Itoa(i)
		s += " Inner: "
		s += f.Inner.String()
		s += " Outer: "
		s += f.Outer.String()
		s += "\n"
	}
	s += "-----\n"
	s += "Vertices\n"
	s += "--------\n"
	for _, v := range dc.Vertices {
		s += v.String()
		s += " OutEdge: "
		s += v.OutEdge.String()
		s += "\n"
	}
	s += "-----\n"
	s += "HalfEdges\n"
	s += "-----\n"
	for _, e := range dc.HalfEdges {
		s += "Origin: "
		s += e.Origin.String()
		s += " Next: "
		s += e.Next.String()
		s += " Prev: "
		s += e.Prev.String()
		s += " Twin: "
		s += e.Twin.String()
		s += " Face: "
		faceIndex := 0
		for i, f := range dc.Faces {
			if f == e.Face {
				faceIndex = i
				break
			}
		}
		s += "f" + strconv.Itoa(faceIndex)
		s += "\n"
	}
	s += "-----\n"
	return s
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

	clock, err := f.Outer.IsClockwise()
	if err == nil && clock {
		f.Outer.Flip()
	} else {
		fmt.Println(err, clock)
	}
	clock, err = f.Inner.IsClockwise()
	if err == nil && !clock {
		f.Inner.Flip()

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
		if f.Inner != nil {
			f2.Inner = dc2.HalfEdges[ePointerMap[f.Inner]]
			f2.Inner.Face = f2
		}
		if f.Outer != nil {
			f2.Outer = dc2.HalfEdges[ePointerMap[f.Outer]]
			f2.Outer.Face = f2
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

// VerticesSorted returns a list indicating the sorted order
// of this dcel's vertices in dimensions ds.
// Example: to get points sorted by x, use with (0)
//          to get points sorted by y, breaking ties
//             on lesser x, use with (1,0).
func (dc *DCEL) VerticesSorted(ds ...int) []int {
	pts := make([]int, len(dc.Vertices))
	for i := range dc.Vertices {
		pts[i] = i
	}
	sort.Slice(pts, func(i, j int) bool {
		p1 := dc.Vertices[pts[i]]
		p2 := dc.Vertices[pts[j]]
		for _, d := range ds {
			v1 := p1.Val(d)
			v2 := p2.Val(d)
			if v1 != v2 {
				return v1 < v2
			}
		}
		return false
	})
	return pts
}

// ConnectVerts takes two vertices and adds edges
// to the dcel containing them to connect the two vertices by a full edge.
// the added edges will be at dc.HalfEdges[len-1] and len-2.
// no face is added, and connectVerts assumes the provided face
// is the face in which the diagonal will land.
//
// The job of creating new faces is delayed because if a series of
// ConnectVerts is called on the same face, calling code won't be
// able to easily tell which face the sequential diagonals land in.
func (dc *DCEL) ConnectVerts(a, b *Vertex, f *Face) {
	// If a and b's outEdges and twins do not
	// share a face, this connection would
	// cross a face and is not allowed. Hypothetically
	// it may be allowed in the future, recursively
	// adding vertices and edges until the connection
	// is complete.
	edgesA := a.AllEdges()
	edgesB := b.AllEdges()
	// It should be impossible for the same face
	// to exist on two edges off of the one vertex
	// in a well-formed DCEL, and if that happens we
	// may change the wrong edges.
	var e1, e2 *Edge
	for _, e := range edgesA {
		if e.Face == f {
			e1 = e
			break
		}
	}
	for _, e := range edgesB {
		if e.Face == f {
			e2 = e
			break
		}
	}
	// If the two vertices share more than one face,
	// there's already an edge here and we can't add
	// another one, unless one of those faces encloses
	// the other, in which case we use the enclosed face.
	//
	// ^^ this is not correct. Consider two vertices on
	// a pseudo-vertical line with both the left and right faces
	// defined. What we need is a good algorithm to determine
	// which of the two faces contains the diagonal
	//
	// The way this algorithm solves this is by taking the
	// face being split in as a hint.

	new1 := NewEdge()
	new1.Origin = b

	new2 := NewEdge()
	new2.Origin = a

	new1.Twin = new2
	new2.Twin = new1

	e1.Prev.Next = new1
	e1.Prev = new2

	e2.Prev.Next = new2
	e2.Prev = new1

	new2.Prev = e2.Prev
	new2.Next = e1

	new1.Prev = e1.Prev
	new1.Next = e2

	new1.Face = e1.Face
	new2.Face = e1.Face

	dc.HalfEdges = append(dc.HalfEdges, new1, new2)
}
