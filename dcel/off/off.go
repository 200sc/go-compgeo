// Package off describes methods for interacting with OFF files and structures
// formatted as OFF files. The file loading code is modeled after Ryan Holmes'
// C++ code, http://www.holmes3d.net/graphics/offfiles/

package off

import "github.com/200sc/go-compgeo/geom"

func NewOFF() OFF {
	return OFF{
		Vertices: make([]Vertex, 0),
		Faces:    make([]Face, 0),
	}
}

// OFF represents the geometric values stored in the OFF format.
// Does not store color.
// Todo: if we could make a structure that satisfied io.Reader,
// we could have less duplicate code here.
type OFF struct {
	NumVertices, NumFaces, NumEdges int
	Vertices                        []Vertex
	Faces                           []Face
}

type Face []int

type Vertex [3]float64

func NewVertex(d geom.D3) Vertex {
	return Vertex{d.X(), d.Y(), d.Z()}

}
