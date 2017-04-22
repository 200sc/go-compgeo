package triangulation

import "github.com/200sc/go-compgeo/dcel"

// YMonotoneSplit converts a dcel into another dcel of its
// triangles, along with a mapping of faces in the new set
// to faces in the input set.
func YMonotoneSplit(dc *dcel.DCEL) (*dcel.DCEL, map[*dcel.Face]*dcel.Face) {
	ypts := dc.VerticesSorted(1)
	for _, i := range ypts {
		v := dc.Vertices[i]
		switch MonotoneVertexType(v, dc) {
		case MONO_START:
			// Begin a monotone
		case MONO_END:
			// Terminate a monotone
		case MONO_REGULAR:
			// Continue a monotone
		case MONO_MERGE:
			// Somehow label we need to draw a diagonal
			// starting here to some other split
		case MONO_SPLIT:
			// Somehow find what merge we already saw
			// before we should draw a diagonal to.
		}

	}
	return nil, nil
}
