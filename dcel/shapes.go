package dcel

import (
	"fmt"
	"math/rand"

	"github.com/200sc/go-compgeo/geom"
)

// Result:
//
// Vertices
//
// 0-1
// | |
// 3-2
//
// Edges
//
//   3>
// .---.
//^|v2<|
//1|0 4|5
// |>6^|v
// .---.
//  <7
//
// Faces
//
// .--.
// |  |
// |0 | 1
// .--.
//

// Rect is a wrapper around FourPoint to make
// a rectangle with top left position and dimensions.
func Rect(x, y, w, h float64) *DCEL {
	return FourPoint(
		geom.Point{x, y, 0},
		geom.Point{x + w, y, 0},
		geom.Point{x + w, y + h, 0},
		geom.Point{x, y + h, 0},
	)
}

// FourPoint creates a dcel from four points, connected
// in order around one face.
func FourPoint(p1, p2, p3, p4 geom.D3) *DCEL {
	dc := new(DCEL)
	dc.Vertices = make([]*Vertex, 4)
	dc.Vertices[0] = PointToVertex(p1)
	dc.Vertices[1] = PointToVertex(p2)
	dc.Vertices[2] = PointToVertex(p3)
	dc.Vertices[3] = PointToVertex(p4)
	dc.HalfEdges = make([]*Edge, 8)
	dc.Faces = make([]*Face, 2)
	dc.Faces[0] = NewFace()
	dc.Faces[1] = NewFace()
	dc.HalfEdges[0] = &Edge{
		Origin: dc.Vertices[0],
		Face:   dc.Faces[0],
	}
	dc.Vertices[0].OutEdge = dc.HalfEdges[0]
	dc.HalfEdges[1] = &Edge{
		Origin: dc.Vertices[3],
		Face:   dc.Faces[1],
	}
	dc.Vertices[3].OutEdge = dc.HalfEdges[1]
	dc.HalfEdges[2] = &Edge{
		Origin: dc.Vertices[1],
		Face:   dc.Faces[0],
	}
	dc.Vertices[1].OutEdge = dc.HalfEdges[2]
	dc.HalfEdges[3] = &Edge{
		Origin: dc.Vertices[0],
		Face:   dc.Faces[1],
	}
	dc.HalfEdges[4] = &Edge{
		Origin: dc.Vertices[2],
		Face:   dc.Faces[0],
	}
	dc.Vertices[2].OutEdge = dc.HalfEdges[4]
	dc.HalfEdges[5] = &Edge{
		Origin: dc.Vertices[1],
		Face:   dc.Faces[1],
	}
	dc.HalfEdges[6] = &Edge{
		Origin: dc.Vertices[3],
		Face:   dc.Faces[0],
	}
	dc.HalfEdges[7] = &Edge{
		Origin: dc.Vertices[2],
		Face:   dc.Faces[1],
	}
	//Twins
	for i := range dc.HalfEdges {
		dc.HalfEdges[i].Twin = dc.HalfEdges[EdgeTwin(i)]
	}
	dc.HalfEdges[0].Prev = dc.HalfEdges[2]
	dc.HalfEdges[0].Next = dc.HalfEdges[6]
	dc.HalfEdges[6].Prev = dc.HalfEdges[0]
	dc.HalfEdges[6].Next = dc.HalfEdges[4]
	dc.HalfEdges[4].Prev = dc.HalfEdges[6]
	dc.HalfEdges[4].Next = dc.HalfEdges[2]
	dc.HalfEdges[2].Prev = dc.HalfEdges[4]
	dc.HalfEdges[2].Next = dc.HalfEdges[0]

	dc.HalfEdges[1].Next = dc.HalfEdges[3]
	dc.HalfEdges[1].Prev = dc.HalfEdges[7]
	dc.HalfEdges[7].Next = dc.HalfEdges[1]
	dc.HalfEdges[7].Prev = dc.HalfEdges[5]
	dc.HalfEdges[5].Next = dc.HalfEdges[7]
	dc.HalfEdges[5].Prev = dc.HalfEdges[3]
	dc.HalfEdges[3].Next = dc.HalfEdges[5]
	dc.HalfEdges[3].Prev = dc.HalfEdges[1]

	dc.Faces[0].Outer = dc.HalfEdges[0]
	dc.Faces[1].Inner = dc.HalfEdges[1]

	// Correcting for faces[0] = the infinite exterior
	dc.Faces[0], dc.Faces[1] = dc.Faces[1], dc.Faces[0]

	dc.CorrectDirectionality(dc.Faces[0])
	dc.CorrectDirectionality(dc.Faces[1])

	return dc
}

func Random2DDCEL(size float64, splits int) *DCEL {
	// Generate a bounding box as a DCEL with one face
	dc := FourPoint(
		geom.NewPoint(0, 0, 0),
		geom.NewPoint(0, size, 0),
		geom.NewPoint(size, size, 0),
		geom.NewPoint(size, 0, 0))

	fmt.Println(dc)

	for i := 0; i < splits; i++ {
		// choose a random face of the dcel
		fi := rand.Intn(len(dc.Faces)-1) + 1
		f := dc.Faces[fi]
		fmt.Println("Face", fi, f)
		// choose two random edges of that face
		edges := f.Outer.EdgeChain()
		fmt.Println("Edges in order of face:")
		e := edges[0]
		for {
			fmt.Println(e)
			e = e.Next
			if e == edges[0] {
				break
			}
		}
		r1 := rand.Intn(len(edges))
		r2 := rand.Intn(len(edges))
		if r2 == r1 {
			r2 = (r2 + 1) % len(edges)
		}
		e1 := edges[r1]
		e2 := edges[r2]
		fmt.Println("Edges chosen")
		fmt.Println("e1,e2", e1, e2)
		// On each edge choose a random point
		v1 := PointToVertex(e1.PointAlong(0, rand.Float64()))
		v2 := PointToVertex(e2.PointAlong(0, rand.Float64()))
		// Add new vertices to dc at p1 and p2,
		dc.Vertices = append(dc.Vertices, v1, v2)
		// Split e1 and e2 and their twins at v1 and v2
		//
		//  e1       e3
		// ------/------/
		//       v1
		// \------\------
		//  t1       t3
		//
		t1 := e1.Twin
		e3 := e1.Copy()
		t3 := t1.Copy()

		e3.SetPrev(e1)
		t1.SetPrev(t3)
		e3.SetTwin(t3)
		t3.Prev.SetNext(t3)
		e3.Next.SetPrev(e3)
		t3.Origin.OutEdge = t3
		t1.Origin = v1
		e3.Origin = v1
		v1.OutEdge = e3
		//
		t2 := e2.Twin
		e4 := e2.Copy()
		t4 := t2.Copy()

		e4.SetPrev(e2)
		t2.SetPrev(t4)
		e4.SetTwin(t4)
		t4.Prev.SetNext(t4)
		e4.Next.SetPrev(e4)
		t4.Origin.OutEdge = t4
		t2.Origin = v2
		e4.Origin = v2
		v2.OutEdge = e4

		fmt.Println("New Edges:")
		fmt.Println("e1, t1", e1, t1)
		fmt.Println("e2, t2", e2, t2)
		fmt.Println("e3, t3", e3, t3)
		fmt.Println("e4, t4", e4, t4)

		dc.HalfEdges = append(dc.HalfEdges, e3, t3, e4, t4)

		// Connect v1 and v2
		//
		//
		//  t4       t2
		// ------/------/
		//      v2
		// \------\------
		//  e4   ||   e2
		//      /||
		//       ||e6
		//     e5||
		//       ||/
		//  e1   ||   e3
		// ------/------/
		//      v1
		// \------\------
		//  t1       t3

		fmt.Println("Connecting", v1, v2)

		e5 := NewEdge()
		e5.Origin = v1

		e6 := NewEdge()
		e6.Origin = v2

		e5.SetTwin(e6)

		e1.SetNext(e5)
		e5.SetNext(e4)
		e2.SetNext(e6)
		e6.SetNext(e3)

		e5.Face = f

		dc.HalfEdges = append(dc.HalfEdges, e5, e6)

		f2 := NewFace()
		dc.Faces = append(dc.Faces, f2)

		e6.Face = f2
		fmt.Println("Walking", e6)
		for e7 := e6.Next; e7 != e6; e7 = e7.Next {
			//fmt.Println("Current Edge", e7)
			e7.Face = f2
		}
		f2.Outer = e6
		f.Outer = e5
	}
	fmt.Println("Random dcel end")
	fmt.Println(dc)

	return dc
}
