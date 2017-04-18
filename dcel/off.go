package dcel

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

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

	dc.Vertices = make([]*Vertex, numVertices)

	// Read numVertices lines as vertices
	// Each vertex is represented as three numbers,
	// x, y, z, in that order.
	for i := 0; i < numVertices; i++ {
		fs, err := readFloat64Line(scanner, 3)
		if err != nil {
			return nil, err
		}
		dc.Vertices[i] = NewVertex(fs[0], fs[1], fs[2])
	}

	var vi int

	edges := make([]*Edge, numEdges)
	edgeIndex := 0
	dc.Faces = make([]*Face, numFaces+1)
	auxData := make(map[*Vertex][]*Edge)

	// Faces are represented by a count of edges followed
	// by a list of vertex indices
	dc.Faces[OUTER_FACE] = new(Face)
	// We start at 1 because 0 is reserved for the outermost
	// face, which this algorithm deals with later
	for i := 1; i < numFaces+1; i++ {
		numEdges, fs, err := readIntsLineNoLength(scanner)
		if err != nil {
			return nil, err
		}

		face = new(Face)
		dc.Faces[i] = face

		edge = new(Edge)
		if edgeIndex >= len(edges) {
			// Some jerk gave us an incorrect definition of their
			// edge count, or we interpreted it wrong.
			edges = append(edges, edge)
		} else {
			edges[edgeIndex] = edge
		}
		edgeIndex++

		// This model does not use Outer faces.
		face.Inner = edge
		edge.Face = face

		vi = fs[0]

		edge.Origin = dc.Vertices[vi]
		dc.Vertices[vi].OutEdge = edge

		aux := auxData[dc.Vertices[vi]]
		if aux == nil {
			aux = make([]*Edge, 0)
		}
		auxData[dc.Vertices[vi]] = append(aux, edge)

		for j := 1; j < numEdges; j++ {
			edge.Next = new(Edge)
			edge.Next.Prev = edge
			edge = edge.Next

			if edgeIndex >= len(edges) {
				// Some jerk gave us an incorrect definition of their
				// edge count, or we interpreted it wrong.
				edges = append(edges, edge)
			} else {
				edges[edgeIndex] = edge
			}
			edgeIndex++
			edge.Face = face

			vi = fs[j]

			edge.Origin = dc.Vertices[vi]
			dc.Vertices[vi].OutEdge = edge

			aux := auxData[dc.Vertices[vi]]
			if aux == nil {
				aux = make([]*Edge, 0)
			}
			auxData[dc.Vertices[vi]] = append(aux, edge)
		}
		edge.Next = face.Inner
		face.Inner.Prev = edge
	}

	var numFound, foundIndex int
	var twin *Edge

	// Create twins
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
				edgeList = append(edgeList, twin)
				twin.Twin = edge
				edge.Twin = twin
				edge.Face = dc.Faces[OUTER_FACE]
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

	// The original algorithm here had auxData attached to each vertex,
	// and called this step "cleaning up" those pointers, when all it was
	// doing was making sure they were all empty.
	//
	// It seems like the algorithm will -always- have elements in some auxdata,
	// as it doesn't remove things from auxdata that already had twins in the
	// twin-attaching phase.
	// for _, v := range auxData {
	// 	if v == nil || len(v) == 0 {
	// 		continue
	// 	}
	// 	for _, v2 := range v {
	// 		if v2 != nil {
	// 			break
	// 		}
	// 	}
	// 	fmt.Println("AuxData was not empty because of course it wasn't")
	// 	isManifold = false
	// 	break
	// }

	if !isManifold {
		return nil, NotManifoldError{}
	}
	var prev *Edge
	for i := 0; i < len(edges); i++ {
		edge = edges[i]
		if edge.Face == dc.Faces[OUTER_FACE] {
			prev = edge.Twin.Next.Twin
			for prev.Next != nil { // Could infinite loop, apparently??
				prev = prev.Next.Twin
			}
			prev.Next = edge
		}
	}
	dc.HalfEdges = make([]*Edge, len(edges))
	ei := 0
	hei := 0
	marked := make(map[*Edge]bool)

	// Our internal DCEL format expects edges[i].Twin to be edges[i+1].
	for hei < len(dc.HalfEdges) {
		if _, ok := marked[edges[ei]]; !ok {
			dc.HalfEdges[hei] = edges[ei]
			dc.HalfEdges[hei+1] = edges[ei].Twin
			marked[edges[ei].Twin] = true
			hei += 2
		}
		ei++
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
		out[i], err = strconv.Atoi(ints[i+1])
		if err != nil {
			return 0, nil, TypeError{}
		}
	}

	return length, out, nil
}
