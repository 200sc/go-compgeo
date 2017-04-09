package dcel

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
func Rect(x, y, w, h float64) *DCEL {
	return FourPoint(
		Point{x, y, 0},
		Point{x + w, y, 0},
		Point{x + w, y + h, 0},
		Point{x, y + h, 0},
	)
}

func FourPoint(p1, p2, p3, p4 Point) *DCEL {
	dc := new(DCEL)
	dc.Vertices = make([]Point, 4)
	dc.Vertices[0] = p1
	dc.Vertices[1] = p2
	dc.Vertices[2] = p3
	dc.Vertices[3] = p4
	dc.OutEdges = make([]*Edge, 4)
	dc.HalfEdges = make([]Edge, 8)
	dc.Faces = make([]Face, 2)
	dc.Faces[0] = Face{}
	dc.Faces[1] = Face{}
	dc.HalfEdges[0] = Edge{
		Origin: &dc.Vertices[0],
		Face:   &dc.Faces[0],
	}
	dc.OutEdges[0] = &dc.HalfEdges[0]
	dc.HalfEdges[1] = Edge{
		Origin: &dc.Vertices[3],
		Face:   &dc.Faces[1],
	}
	dc.OutEdges[3] = &dc.HalfEdges[1]
	dc.HalfEdges[2] = Edge{
		Origin: &dc.Vertices[1],
		Face:   &dc.Faces[0],
	}
	dc.OutEdges[1] = &dc.HalfEdges[2]
	dc.HalfEdges[3] = Edge{
		Origin: &dc.Vertices[0],
		Face:   &dc.Faces[1],
	}
	dc.HalfEdges[4] = Edge{
		Origin: &dc.Vertices[2],
		Face:   &dc.Faces[0],
	}
	dc.OutEdges[2] = &dc.HalfEdges[4]
	dc.HalfEdges[5] = Edge{
		Origin: &dc.Vertices[1],
		Face:   &dc.Faces[1],
	}
	dc.HalfEdges[6] = Edge{
		Origin: &dc.Vertices[3],
		Face:   &dc.Faces[0],
	}
	dc.HalfEdges[7] = Edge{
		Origin: &dc.Vertices[2],
		Face:   &dc.Faces[1],
	}
	//Twins
	for i := range dc.HalfEdges {
		dc.HalfEdges[i].Twin = &dc.HalfEdges[EdgeTwin(i)]
	}
	dc.HalfEdges[0].Prev = &dc.HalfEdges[2]
	dc.HalfEdges[0].Next = &dc.HalfEdges[6]
	dc.HalfEdges[6].Prev = &dc.HalfEdges[0]
	dc.HalfEdges[6].Next = &dc.HalfEdges[4]
	dc.HalfEdges[4].Prev = &dc.HalfEdges[6]
	dc.HalfEdges[4].Next = &dc.HalfEdges[2]
	dc.HalfEdges[2].Prev = &dc.HalfEdges[4]
	dc.HalfEdges[2].Next = &dc.HalfEdges[0]

	dc.HalfEdges[1].Next = &dc.HalfEdges[3]
	dc.HalfEdges[1].Prev = &dc.HalfEdges[7]
	dc.HalfEdges[7].Next = &dc.HalfEdges[1]
	dc.HalfEdges[7].Prev = &dc.HalfEdges[5]
	dc.HalfEdges[5].Next = &dc.HalfEdges[7]
	dc.HalfEdges[5].Prev = &dc.HalfEdges[3]
	dc.HalfEdges[3].Next = &dc.HalfEdges[5]
	dc.HalfEdges[3].Prev = &dc.HalfEdges[1]

	dc.Faces[0].Outer = &dc.HalfEdges[0]
	dc.Faces[1].Inner = &dc.HalfEdges[1]

	return dc
}
