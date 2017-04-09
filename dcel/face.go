package dcel

// A DCELFace points to the edges on its inner and
// outer portions. Any given face may have either
// of these values be nil, but never both.
type Face struct {
	Outer, Inner *Edge
}
