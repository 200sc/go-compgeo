package demo

import (
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"bitbucket.org/oakmoundstudio/oak/render"
)

// An InteractivePolyhedron is a wrapper around
// an oak polyhedron which defines mouse collision
// areas to interact with parts of the underlying DCEL
// structure.
type InteractivePolyhedron struct {
	*render.Polyhedron
	vs []*InteractivePoint
	// This is more than a little impractical
	// until collision spaces can contain internal
	// polygons
	// eSpaces []*collision.Space
	//overSpace *collision.Space
	cID event.CID
}

// Init establishes ip's non-polyhedron variables.
func (ip *InteractivePolyhedron) Init() event.CID {
	ip.cID = event.NextID(ip)
	ip.vs = make([]*InteractivePoint, len(ip.Vertices))
	for i, v := range ip.Vertices {
		ip.vs[i] = NewInteractivePoint(v, i)
	}
	return ip.cID
}

// UpdateSpaces is a helper function to polyhedron's Update
// which similarly resets parts of the ip as it is moved
// through space. In this case the large job here is making sure
// all of the vertex collision areas stay in the right spots.
func (ip *InteractivePolyhedron) UpdateSpaces() {
	if len(ip.vs) < len(ip.Vertices) {
		diff := len(ip.Vertices) - len(ip.vs)
		ip.vs = append(ip.vs, make([]*InteractivePoint, diff)...)
	}
	for i, v := range ip.Vertices {
		if ip.vs[i] == nil {
			ip.vs[i] = NewInteractivePoint(v, i)
		}
		ip.vs[i].Vertex = v
		mouse.UpdateSpace(ip.X+(v.X()-vCollisionDim/2),
			ip.Y+(v.Y()-vCollisionDim/2),
			vCollisionDim, vCollisionDim, ip.vs[i].s)
	}
}
