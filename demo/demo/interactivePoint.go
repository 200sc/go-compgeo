package demo

import (
	"time"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/render"
	"github.com/200sc/go-compgeo/dcel"
)

// InteractivePoint is a struct to wrap around dcel points
// and extend them to be bindable and have collision space.
type InteractivePoint struct {
	*dcel.Vertex
	s            *collision.Space
	cID          event.CID
	index        int
	mousedOverCh chan bool
	showing      bool
}

// Init allows ip to satisfy the event.Entity interface,
// so it may be stored with other entities in a global
// list managed by oak.
func (ip *InteractivePoint) Init() event.CID {
	ip.cID = event.NextID(ip)
	return ip.cID
}

// NewInteractivePoint creates a new ip given a dcel point to base it off of.
func NewInteractivePoint(v *dcel.Vertex, i int) *InteractivePoint {
	ip := new(InteractivePoint)
	ip.Init()
	ip.s = collision.NewSpace(0, 0, 1, 1, ip.cID)
	mouse.Add(ip.s)
	ip.cID.Bind(vertexStartDrag, "MousePressOn")
	ip.cID.Bind(vertexShow, "MouseDragOn")
	ip.index = i
	ip.Vertex = v
	ip.mousedOverCh = make(chan bool)
	return ip
}

func vertexShow(cID int, nothing interface{}) int {
	ip := event.GetEntity(cID).(*InteractivePoint)
	if !ip.showing {
		ip.showing = true
		txt := font.NewText(ip.Point, ip.s.GetX(), ip.s.GetY())
		render.Draw(txt, 3)
		go func() {
			for {
				select {
				case <-time.After(250 * time.Millisecond):
					nme := mouse.LastMouseEvent
					if ip.s.Contains(nme.ToSpace()) {
						continue
					}
					txt.UnDraw()
					ip.showing = false
					return
				case v := <-ip.mousedOverCh:
					if !v {
						return
					}
				}
			}
		}()
	} else {
		ip.mousedOverCh <- true
	}
	return 0
}

func vertexStartDrag(cID int, nothing interface{}) int {
	if mode == MOVE_POINT {
		ip := event.GetEntity(cID).(*InteractivePoint)
		dragging = ip.index
	}
	return 0
}

func vertexStopDrag(no int, nothing interface{}) int {
	dragging = -1
	return 0
}
