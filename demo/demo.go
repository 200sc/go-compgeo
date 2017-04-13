package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/200sc/go-compgeo/dcel"

	"time"

	"bitbucket.org/oakmoundstudio/oak"
	"bitbucket.org/oakmoundstudio/oak/collision"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"bitbucket.org/oakmoundstudio/oak/render"
)

type mouseMode int

// control mode constant
const (
	ROTATE mouseMode = iota
	MOVE_POINT
	POINT_LOCATE
	ADD_DCEL
	REM_DCEL
	LAST_MODE
	ADDING_DCEL
)

func (m mouseMode) String() string {
	switch m {
	case ROTATE:
		return "Rotate"
	case MOVE_POINT:
		return "Move Point"
	case POINT_LOCATE:
		return "Point Location"
	case ADD_DCEL:
		return "Define Face"
	case ADDING_DCEL:
		return "Defining Face..."
	case REM_DCEL:
		return "Define Inside Face"
	default:
		return "INVALID"
	}
}

const (
	zMoveSpeed    = 1
	shiftSpeed    = 3
	scaleSpeed    = .02
	vCollisionDim = 8
	defScale      = 20
	defRotZ       = math.Pi
	defRotY       = math.Pi
)

var (
	dragX           float32 = -1
	dragY           float32 = -1
	dragging                = -1
	offFile                 = filepath.Join("data", "A.off")
	mode                    = ROTATE
	loopDemo        bool
	firstAddedPoint int
	prev            *dcel.Edge
	addedFace       *dcel.Face
	mouseZ          = 0.0
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		offFile = args[0]
	}
	err := oak.LoadConf("oak.config")
	if err != nil {
		log.Fatal(err)
	}
	oak.AddCommand("load", func(strs []string) {
		if len(strs) > 1 {
			offFile = strs[1]
			loopDemo = false
		}
	})
	oak.AddScene("demo",
		func(prevScene string, data interface{}) {
			loopDemo = true
			//phd := render.NewCuboid(100, 100, 100, 100, 100, 100)
			dc, err := dcel.LoadOFF(offFile)
			if err != nil {
				log.Fatal(err)
			}
			phd := new(InteractivePolyhedron)
			phd.Polyhedron = render.NewPolyhedronFromDCEL(dc, 100, 100)
			phd.Polyhedron.Scale(defScale)
			phd.Polyhedron.RotZ(defRotZ)
			phd.Polyhedron.RotY(defRotY)
			phd.Init()
			render.Draw(phd, 2)

			modeStr := render.DefFont().NewText(mode.String(), 3, 40)
			render.Draw(modeStr, 1)

			mouseStr := render.DefFont().NewInterfaceText(
				dcel.Point{0, 0, 0}, 3, 465)

			render.Draw(mouseStr, 1)

			event.GlobalBind(vertexStopDrag, "MouseRelease")
			event.GlobalBind(func(no int, nothing interface{}) int {
				mode = (mode + 1) % LAST_MODE
				modeStr.SetText(mode.String())
				return 0
			}, "KeyDownQ")
			event.GlobalBind(func(no int, nothing interface{}) int {
				mode = 0
				modeStr.SetText(mode.String())
				return 0
			}, "KeyDown1")
			event.GlobalBind(func(no int, nothing interface{}) int {
				mode = 1
				modeStr.SetText(mode.String())
				return 0
			}, "KeyDown2")
			event.GlobalBind(func(no int, nothing interface{}) int {
				mode = 2
				modeStr.SetText(mode.String())
				return 0
			}, "KeyDown3")
			event.GlobalBind(func(no int, nothing interface{}) int {
				mode = 3
				modeStr.SetText(mode.String())
				return 0
			}, "KeyDown4")
			phd.cID.Bind(func(no int, nothing interface{}) int {
				shft := oak.IsDown("LeftShift")
				if oak.IsDown("LeftArrow") {
					phd.ShiftX(-shiftSpeed)
					phd.UpdateSpaces()
				} else if oak.IsDown("RightArrow") {
					phd.ShiftX(shiftSpeed)
					phd.UpdateSpaces()
				}
				if oak.IsDown("UpArrow") {
					if shft {
						phd.Scale(1 + scaleSpeed)
						phd.UpdateSpaces()
					} else {
						phd.ShiftY(-shiftSpeed)
						phd.UpdateSpaces()
					}
				} else if oak.IsDown("DownArrow") {
					if shft {
						phd.Scale(1 - scaleSpeed)
						phd.UpdateSpaces()
					} else {
						phd.ShiftY(shiftSpeed)
						phd.UpdateSpaces()
					}
				}
				nme := mouse.LastMouseEvent
				mouseStr.SetText(dcel.Point{float64(nme.X) - phd.X,
					float64(nme.Y) - phd.Y, mouseZ})
				if mode == ROTATE {
					if dragX != -1 {
						dx := float64(nme.X - dragX)
						dy := float64(nme.Y - dragY)
						if dx != 0 {
							if shft {
								phd.RotZ(.01 * dx)
								phd.UpdateSpaces()
							} else {
								phd.RotY(.01 * dx)
								phd.UpdateSpaces()
							}
						}
						if dy != 0 {
							phd.RotX(.01 * dy)
							phd.UpdateSpaces()
						}
					}
					if oak.IsDown("D") {
						mouseZ += zMoveSpeed
					} else if oak.IsDown("C") {
						mouseZ -= zMoveSpeed
					}
				} else if mode == MOVE_POINT && dragging != -1 {
					update := false
					mouseZ = phd.Vertices[dragging][2]
					if dragX != -1 {
						phd.Vertices[dragging][0] = float64(dragX) - phd.X
						update = true
					}
					if dragY != -1 {
						phd.Vertices[dragging][1] = float64(dragY) - phd.Y
						update = true
					}
					if oak.IsDown("D") {
						phd.Vertices[dragging][2] += zMoveSpeed
						update = true
					} else if oak.IsDown("C") {
						phd.Vertices[dragging][2] -= zMoveSpeed
						update = true
					}
					if update {
						phd.Update()
						phd.UpdateSpaces()
					}
				} else if mode == ADD_DCEL {
					if oak.IsDown("D") {
						mouseZ += zMoveSpeed
					} else if oak.IsDown("C") {
						mouseZ -= zMoveSpeed
					}
					// Detect clicks
					// On first click, declare the first point and edge and face
					// First off, assume that this point is new
				} else if mode == ADDING_DCEL {
					// On following clicks, prev.next = next, next.prev = prev
					//                      next.origin = origin
					//                      next.face = theface
					//                      prev.twin = make a twin at origin
					//                      twin.face = 0 I guess for now
					// Detect final click by clicking on first point or
					//                      by right clicking
				}
				if oak.IsDown("LeftMouse") {
					dragX = nme.X
					dragY = nme.Y
				} else {
					dragX = -1
					dragY = -1
				}
				return 0
			}, "EnterFrame")
			event.GlobalBind(func(no int, event interface{}) int {
				me := event.(mouse.MouseEvent)
				if me.Button == "LeftMouse" {
					if mode == ADD_DCEL {
						firstAddedPoint = len(phd.Vertices)
						phd.Vertices = append(phd.Vertices,
							dcel.NewPoint(float64(me.X)-phd.X, float64(me.Y)-phd.Y, mouseZ))

						phd.HalfEdges = append(phd.HalfEdges, dcel.NewEdge())
						prev = phd.HalfEdges[len(phd.HalfEdges)-1]
						prev.Origin = phd.Vertices[firstAddedPoint]

						f := dcel.NewFace()
						f.Inner = prev
						phd.Faces = append(phd.Faces, f)
						addedFace = phd.Faces[len(phd.Faces)-1]

						prev.Face = addedFace
						phd.OutEdges = append(phd.OutEdges, prev)

						mode = ADDING_DCEL

						phd.Update()
						phd.UpdateSpaces()

					} else if mode == ADDING_DCEL {
						phd.Vertices = append(phd.Vertices,
							dcel.NewPoint(float64(me.X)-phd.X, float64(me.Y)-phd.Y, mouseZ))

						phd.HalfEdges = append(phd.HalfEdges, dcel.NewEdge())
						twin := phd.HalfEdges[len(phd.HalfEdges)-1]
						twin.Origin = phd.Vertices[len(phd.Vertices)-1]
						twin.Face = phd.Faces[dcel.OUTER_FACE]
						prev.Twin = twin
						twin.Twin = prev
						// If prev.Prev is nil, we add the pointer
						// on right click.
						if prev.Prev != nil {
							twin.Next = prev.Prev.Twin
						}

						phd.HalfEdges = append(phd.HalfEdges, dcel.NewEdge())
						next := phd.HalfEdges[len(phd.HalfEdges)-1]
						next.Origin = phd.Vertices[len(phd.Vertices)-1]
						next.Prev = prev
						next.Face = addedFace
						prev.Next = next

						prev = next

						phd.Update()
						phd.UpdateSpaces()
					}
				} else if me.Button == "RightMouse" {
					if mode == ADDING_DCEL {
						prev.Next = phd.OutEdges[firstAddedPoint]
						phd.OutEdges[firstAddedPoint].Prev = prev

						phd.HalfEdges = append(phd.HalfEdges, dcel.NewEdge())
						// The final twin
						twin := phd.HalfEdges[len(phd.HalfEdges)-1]
						twin.Origin = phd.Vertices[firstAddedPoint]
						twin.Face = phd.Faces[dcel.OUTER_FACE]
						prev.Twin = twin
						twin.Twin = prev
						twin.Prev = prev.Next.Twin
						prev.Next.Twin.Prev = twin
						twin.Next = prev.Prev.Twin
						prev.Prev.Twin.Next = twin

						prev = nil
						addedFace = nil
						firstAddedPoint = -1

						mode = ADD_DCEL

						phd.Update()
						phd.UpdateSpaces()

					}
				}
				return 0
			}, "MouseRelease")
			// event.GlobalBind(func(no int, me interface{}) int {
			// 	event := me.(mouse.MouseEvent)
			// 	if event.Button == "LeftMouse" {
			// 		fmt.Println(event.X, event.Y, event.Button)
			// 		dragX = event.X
			// 		dragY = event.Y
			// 	}
			// 	return 0
			// }, "MousePress")
			// event.GlobalBind(func(no int, me interface{}) int {
			// 	event := me.(mouse.MouseEvent)
			// 	if event.Button == "LeftMouse" {
			// 		fmt.Println(event.X, event.Y, event.Button)
			// 		dragX = -1
			// 		dragY = -1
			// 	}
			// 	return 0
			// }, "MouseRelease")
		},
		func() bool {
			return loopDemo
		},
		func() (string, *oak.SceneResult) {
			return "demo", nil
		},
	)
	oak.Init("demo")
}

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
