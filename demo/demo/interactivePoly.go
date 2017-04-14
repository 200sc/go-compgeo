package demo

import (
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"bitbucket.org/oakmoundstudio/oak/render"
)

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

func (ip *InteractivePolyhedron) Init() event.CID {
	ip.cID = event.NextID(ip)
	ip.vs = make([]*InteractivePoint, len(ip.Vertices))
	for i, v := range ip.Vertices {
		ip.vs[i] = NewInteractivePoint(v, i)
	}
	return ip.cID
}

func (ip *InteractivePolyhedron) UpdateSpaces() {
	if len(ip.vs) < len(ip.Vertices) {
		diff := len(ip.Vertices) - len(ip.vs)
		ip.vs = append(ip.vs, make([]*InteractivePoint, diff)...)
	}
	for i, v := range ip.Vertices {
		if ip.vs[i] == nil {
			ip.vs[i] = NewInteractivePoint(v, i)
		}
		ip.vs[i].Point = v
		mouse.UpdateSpace(ip.X+(v[0]-vCollisionDim/2),
			ip.Y+(v[1]-vCollisionDim/2),
			vCollisionDim, vCollisionDim, ip.vs[i].s)
	}
}
