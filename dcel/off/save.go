package off

import (
	"io/ioutil"
	"strconv"

	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/printutil"
)

// Save converts a DCEL into an OFF structure.
func Save(dc *dcel.DCEL) OFF {
	of := NewOFF()
	of.Vertices = make([]Vertex, len(dc.Vertices))
	vMap := make(map[*dcel.Vertex]int)
	for i, v := range dc.Vertices {
		vMap[v] = i
		of.Vertices[i] = NewVertex(v)
	}

	of.Faces = make([]Face, len(dc.Faces)-1)
	for i := 1; i < len(dc.Faces); i++ {
		vs := dc.Faces[i].Vertices()
		offFace := make(Face, len(vs))
		for j, v := range vs {
			offFace[j] = vMap[v]
		}
		of.Faces[i-1] = offFace
	}
	of.NumVertices = len(dc.Vertices)
	of.NumFaces = len(dc.Faces) - 1
	of.NumEdges = len(dc.HalfEdges)
	return of
}

// WriteFile takes an OFF structure and writes it to
// the given relative path.
func (of OFF) WriteFile(relPath string) error {
	// This could be made much faster using the bufio package
	bData := []byte("OFF\n")
	bData = append(bData, strconv.Itoa(of.NumVertices)...)
	bData = append(bData, ' ')
	bData = append(bData, strconv.Itoa(of.NumFaces)...)
	bData = append(bData, ' ')
	bData = append(bData, strconv.Itoa(of.NumEdges)...)
	bData = append(bData, '\n')
	for _, v := range of.Vertices {
		for i := 0; i < 3; i++ {
			bData = append(bData, printutil.Stringf64(v[i])...)
			if i != 2 {
				bData = append(bData, ' ')
			}
		}
		bData = append(bData, '\n')
	}
	for _, f := range of.Faces {
		bData = append(bData, strconv.Itoa(len(f))...)
		for i := 0; i < len(f); i++ {
			bData = append(bData, ' ')
			bData = append(bData, strconv.Itoa(f[i])...)
		}
		bData = append(bData, '\n')
	}
	return ioutil.WriteFile(relPath, bData, 0644)
}
