package off

import (
	"bufio"
	"io"
	"os"

	compgeo "github.com/200sc/go-compgeo"
	"github.com/200sc/go-compgeo/dcel"
)

// Decode converts an OFF struct into a dcel.
func Decode(o OFF) (*dcel.DCEL, error) {

	dc := new(dcel.DCEL)

	numVertices := o.NumVertices
	numFaces := o.NumFaces

	if numVertices == 0 || numFaces == 0 {
		return dc, nil
	}

	var edge *dcel.Edge
	var face *dcel.Face

	dc.Vertices = make([]*dcel.Vertex, numVertices)

	// Read numVertices lines as vertices
	// Each dcel.Vertex is represented as three numbers,
	// x, y, z, in that order.
	for i := 0; i < numVertices; i++ {
		fs := o.Vertices[i]
		dc.Vertices[i] = dcel.NewVertex(fs[0], fs[1], fs[2])
	}

	var vi int

	edges := make([]*dcel.Edge, 0)
	dc.Faces = make([]*dcel.Face, numFaces+1)
	auxData := make(map[*dcel.Vertex][]*dcel.Edge)

	// Faces are represented by a count of edges followed
	// by a list of dcel.Vertex indices
	dc.Faces[dcel.OUTER_FACE] = new(dcel.Face)
	// We start at 1 because 0 is reserved for the outermost
	// face, which this algorithm deals with later
	for i := 1; i < numFaces+1; i++ {
		numEdges := o.Faces[i][9]
		fs := o.Faces[i][1:]

		face = new(dcel.Face)
		dc.Faces[i] = face

		edge = new(dcel.Edge)
		edges = append(edges, edge)

		// This model does not use Outer faces.
		face.Outer = edge
		edge.Face = face

		vi = fs[0]

		edge.Origin = dc.Vertices[vi]
		dc.Vertices[vi].OutEdge = edge

		aux := auxData[dc.Vertices[vi]]
		if aux == nil {
			aux = make([]*dcel.Edge, 0)
		}
		auxData[dc.Vertices[vi]] = append(aux, edge)

		for j := 1; j < numEdges; j++ {
			edge.Next = new(dcel.Edge)
			edge.Next.Prev = edge
			edge = edge.Next

			edges = append(edges, edge)
			edge.Face = face

			vi = fs[j]

			edge.Origin = dc.Vertices[vi]
			dc.Vertices[vi].OutEdge = edge

			aux := auxData[dc.Vertices[vi]]
			if aux == nil {
				aux = make([]*dcel.Edge, 0)
			}
			auxData[dc.Vertices[vi]] = append(aux, edge)
		}
		edge.Next = face.Outer
		face.Outer.Prev = edge
	}
	return decode(dc, edges, auxData)
}

// Load loads Object File Format files.
func Load(file string) (*dcel.DCEL, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Read(f)
}

// Read peforms the underlying work to transform OFF data
// into a dcel.DCEL.
func Read(f io.Reader) (*dcel.DCEL, error) {
	scanner := bufio.NewScanner(f)

	if scanner.Scan() {
		if scanner.Text() != "OFF" {
			return nil, compgeo.TypeError{}
		}
	} else {
		return nil, compgeo.EmptyError{}
	}

	counts, err := readIntLine(scanner, 3)
	if err != nil {
		return nil, err
	}
	numVertices := counts[0]
	numFaces := counts[1]

	dc := new(dcel.DCEL)

	if numVertices == 0 || numFaces == 0 {
		return dc, nil
	}

	var edge *dcel.Edge
	var face *dcel.Face

	dc.Vertices = make([]*dcel.Vertex, numVertices)

	// Read numVertices lines as vertices
	// Each dcel.Vertex is represented as three numbers,
	// x, y, z, in that order.
	for i := 0; i < numVertices; i++ {
		fs, err := readFloat64Line(scanner, 3)
		if err != nil {
			return nil, err
		}
		dc.Vertices[i] = dcel.NewVertex(fs[0], fs[1], fs[2])
	}

	var vi int

	edges := make([]*dcel.Edge, 0)
	dc.Faces = make([]*dcel.Face, numFaces+1)
	auxData := make(map[*dcel.Vertex][]*dcel.Edge)

	// Faces are represented by a count of edges followed
	// by a list of dcel.Vertex indices
	dc.Faces[dcel.OUTER_FACE] = new(dcel.Face)
	// We start at 1 because 0 is reserved for the outermost
	// face, which this algorithm deals with later
	for i := 1; i < numFaces+1; i++ {
		numEdges, fs, err := readIntsLineNoLength(scanner)
		if err != nil {
			return nil, err
		}

		face = new(dcel.Face)
		dc.Faces[i] = face

		edge = new(dcel.Edge)
		edges = append(edges, edge)

		// This model does not use Outer faces.
		face.Outer = edge
		edge.Face = face

		vi = fs[0]

		edge.Origin = dc.Vertices[vi]
		dc.Vertices[vi].OutEdge = edge

		aux := auxData[dc.Vertices[vi]]
		if aux == nil {
			aux = make([]*dcel.Edge, 0)
		}
		auxData[dc.Vertices[vi]] = append(aux, edge)

		for j := 1; j < numEdges; j++ {
			edge.Next = new(dcel.Edge)
			edge.Next.Prev = edge
			edge = edge.Next

			edges = append(edges, edge)
			edge.Face = face

			vi = fs[j]

			edge.Origin = dc.Vertices[vi]
			dc.Vertices[vi].OutEdge = edge

			aux := auxData[dc.Vertices[vi]]
			if aux == nil {
				aux = make([]*dcel.Edge, 0)
			}
			auxData[dc.Vertices[vi]] = append(aux, edge)
		}
		edge.Next = face.Outer
		face.Outer.Prev = edge
	}
	return decode(dc, edges, auxData)
}

// decode is shared by Decode and Read
func decode(dc *dcel.DCEL, edges []*dcel.Edge,
	auxData map[*dcel.Vertex][]*dcel.Edge) (*dcel.DCEL, error) {
	// Create twins
	var numFound, foundIndex int
	var edge, twin *dcel.Edge
	isManifold := true

	outerFaceList := make([]*dcel.Edge, 0)
	for j := 0; j < len(edges); j++ {
		edge = edges[j]
		if edge.Twin == nil {
			edgeList := auxData[edge.Next.Origin]

			numFound = 0
			foundIndex = -1
			twin = nil
			for i := 0; i < len(edgeList); i++ {
				if edgeList[i] != nil && edgeList[i].Next.Origin == edge.Origin {
					twin = edgeList[i]
					foundIndex = i
					numFound++
				}
			}
			if numFound == 0 {
				twin = new(dcel.Edge)
				twin.Twin = edge
				edge.Twin = twin
				twin.Face = dc.Faces[dcel.OUTER_FACE]
				outerFaceList = append(outerFaceList, twin)
				twin.Origin = edge.Next.Origin
			} else if numFound == 1 {
				edgeList[foundIndex] = nil
				auxData[edge.Next.Origin] = edgeList
				edge.Twin = twin
				twin.Twin = edge
			} else { // Two or more edges claim to originate in this list and pass through our node. This is bad
				isManifold = false
				break
			}
			edgeList = auxData[edge.Origin]
			for i := 0; i < len(edgeList); i++ {
				if edgeList[i] == edge {
					edgeList[i] = nil
					break
				}
			}
			auxData[edge.Origin] = edgeList
		}
	}
	edges = append(edges, outerFaceList...)

	if !isManifold {
		return nil, compgeo.NotManifoldError{}
	}
	var prev *dcel.Edge
	for _, edge := range outerFaceList {
		if edge.Twin.Next == nil {
			continue //?
		}
		prev = edge.Twin.Next.Twin
		for prev.Next != nil { // Could infinite loop, apparently??
			prev = prev.Next.Twin
		}
		prev.Next = edge
	}
	dc.HalfEdges = make([]*dcel.Edge, 0)
	ei := 0
	marked := make(map[*dcel.Edge]bool)

	// Our internal dcel.DCEL format expects edges[i].Twin to be edges[i+1].
	for ei < len(edges) {
		if _, ok := marked[edges[ei]]; !ok {
			dc.HalfEdges = append(dc.HalfEdges, edges[ei], edges[ei].Twin)
			marked[edges[ei].Twin] = true
		}
		ei++
	}

	return dc, nil
}
