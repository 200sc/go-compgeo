package dcel

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

type OFFdata struct {
}

type TypeError struct{}

func (fte TypeError) Error() string {
	return "The input was not of the expected type, or was malformed"
}

type EmptyError struct{}

func (ee EmptyError) Error() string {
	return "The input was empty"
}

type NotManifoldError struct{}

func (nme NotManifoldError) Error() string {
	return "The given shape was not manifold"
}

// LoadOFF loads Object File Format files. This function
// is modeled after Ryan Holmes' C++ code, http://www.holmes3d.net/graphics/offfiles/
func LoadOFF(file string) (*DCEL, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ReadOFF(f)
}

// ReadOFF peforms the underlying work to transform OFF data
// into a DCEL.
func ReadOFF(f io.Reader) (*DCEL, error) {
	scanner := bufio.NewScanner(f)

	isManifold := true

	if scanner.Scan() {
		if scanner.Text() != "OFF" {
			return nil, TypeError{}
		}
	} else {
		return nil, EmptyError{}
	}

	counts, err := readIntLine(scanner, 3)
	if err != nil {
		return nil, err
	}
	numVertices := counts[0]
	numFaces := counts[1]
	numEdges := counts[2]

	dc := new(DCEL)

	if numVertices == 0 || numFaces == 0 {
		return dc, nil
	}

	var edge *Edge
	var face *Face

	dc.Vertices = make([]Point, numVertices)

	// Read numVertices lines as vertices
	// Each vertex is represented as three numbers,
	// x, y, z, in that order.
	for i := 0; i < numVertices; i++ {
		fs, err := readFloat64Line(scanner, 3)
		if err != nil {
			return nil, err
		}
		dc.Vertices[i] = Point{fs[0], fs[1], fs[2]}
	}

	var vi int

	dc.OutEdges = make([]*Edge, numVertices)

	edges := make([]*Edge, numEdges)
	edgeIndex := 0
	faces := make([]*Face, numFaces)
	auxData := make(map[*Point][]*Edge)

	// Faces are represented by a count of edges followed
	// by a list of vertex indices
	faces[OUTER_FACE] = new(Face)
	// We start at 1 because 0 is reserved for the outermost
	// face, which this algorithm deals with later
	for i := 1; i < numFaces+1; i++ {
		numEdges, fs, err := readIntsLineNoLength(scanner)
		if err != nil {
			return nil, err
		}

		face = new(Face)
		faces[i] = face

		edge = new(Edge)
		edges[edgeIndex] = edge
		edgeIndex++

		// Inner or outer?
		// The model we are basing off of does not
		// have faces with both inner and outer edges.
		// How is the outer face represented?
		// What about holes??
		face.Inner = edge
		edge.Face = face

		vi = fs[0]

		edge.Origin = &dc.Vertices[vi]
		dc.OutEdges[vi] = edge

		aux := auxData[&dc.Vertices[vi]]
		if aux == nil {
			aux = make([]*Edge, 0)
		}
		auxData[&dc.Vertices[vi]] = append(aux, edge)

		for j := 1; j < numEdges; j++ {
			edge.Next = new(Edge)
			edge = edge.Next

			edges[edgeIndex] = edge
			edgeIndex++
			edge.Face = face

			vi = fs[j]

			edge.Origin = &dc.Vertices[vi]
			dc.OutEdges[vi] = edge

			aux := auxData[&dc.Vertices[vi]]
			if aux == nil {
				aux = make([]*Edge, 0)
			}
			auxData[&dc.Vertices[vi]] = append(aux, edge)
		}
		edge.Next = face.Inner
	}

	var numFound, foundIndex int
	var twin *Edge

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
				twin = new(Edge)
				edgeList[edgeIndex] = twin
				twin.Twin = edge
				edge.Twin = twin
				edge.Face = faces[OUTER_FACE]
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

	// Even if we've decided the mesh is non-manifold, we need to clean up the auxData pointers before we clear the mesh
	// (If the mesh is manifold, this is easy. If not, we almost definitely have auxData pointers to clear)
	for _, v := range auxData {
		if v == nil || len(v) == 0 {
			continue
		}
		isManifold = false
		break
	}

	if !isManifold {
		return nil, NotManifoldError{}
	}
	var prev *Edge
	for i := 0; i < len(edges); i++ {
		edge = edges[i]
		if edge.Face == faces[OUTER_FACE] {
			prev = edge.Twin.Next.Twin
			for prev.Next != nil { // Could infinite loop, apparently??
				prev = prev.Next.Twin
			}
			prev.Next = edge
		}
	}

	return dc, nil
}

func readIntLine(s *bufio.Scanner, l int) ([]int, error) {
	var err error
	out := make([]int, l)

	if !s.Scan() {
		return out, TypeError{}
	}

	ints := strings.Split(s.Text(), " ")
	if len(ints) < l {
		return nil, TypeError{}
	}

	for i := 0; i < l; i++ {
		out[i], err = strconv.Atoi(ints[i])
		if err != nil {
			return nil, TypeError{}
		}
	}

	return out, nil
}

func readFloat64Line(s *bufio.Scanner, l int) ([]float64, error) {
	var err error
	out := make([]float64, l)

	if !s.Scan() {
		return out, TypeError{}
	}

	ints := strings.Split(s.Text(), " ")
	if len(ints) < l {
		return nil, TypeError{}
	}

	for i := 0; i < l; i++ {
		out[i], err = strconv.ParseFloat(ints[i], 64)
		if err != nil {
			return nil, TypeError{}
		}
	}

	return out, nil
}

// The number of elements in this line is defined by the first value.
func readIntsLineNoLength(s *bufio.Scanner) (int, []int, error) {
	var err error
	if !s.Scan() {
		return 0, make([]int, 0), TypeError{}
	}

	ints := strings.Split(s.Text(), " ")

	length, err := strconv.Atoi(ints[0])
	if err != nil {
		return 0, nil, TypeError{}
	}

	out := make([]int, length)

	if len(ints) < (length + 1) {
		return 0, nil, TypeError{}
	}

	for i := 0; i < length; i++ {
		out[i], err = strconv.Atoi(ints[i])
		if err != nil {
			return 0, nil, TypeError{}
		}
	}

	return length, out, nil
}
