package demo

import (
	"fmt"
	"time"

	"bitbucket.org/oakmoundstudio/oak/collision"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"bitbucket.org/oakmoundstudio/oak/render"
	"github.com/200sc/go-compgeo/dcel"
)

type InteractivePoint struct {
	*dcel.Point
	s            *collision.Space
	cID          event.CID
	index        int
	mousedOverCh chan bool
	showing      bool
}

func NewInteractivePoint(v *dcel.Point, i int) *InteractivePoint {
	ip := new(InteractivePoint)
	ip.Init()
	ip.s = collision.NewSpace(0, 0, 1, 1, ip.cID)
	mouse.Add(ip.s)
	ip.cID.Bind(vertexStartDrag, "MousePressOn")
	ip.cID.Bind(vertexShow, "MouseDragOn")
	ip.index = i
	ip.Point = v
	ip.mousedOverCh = make(chan bool)
	return ip
}

func vertexShow(cID int, nothing interface{}) int {
	ip := event.GetEntity(cID).(*InteractivePoint)
	fmt.Println("Mouse drag on triggered")
	if !ip.showing {
		ip.showing = true
		txt := render.DefFont().NewInterfaceText(ip.Point, ip.s.GetX(), ip.s.GetY())
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
		fmt.Println("Start drag")
		ip := event.GetEntity(cID).(*InteractivePoint)
		dragging = ip.index
	}
	return 0
}

func vertexStopDrag(no int, nothing interface{}) int {
	dragging = -1
	return 0
}

func (ip *InteractivePoint) Init() event.CID {
	ip.cID = event.NextID(ip)
	return ip.cID
}
