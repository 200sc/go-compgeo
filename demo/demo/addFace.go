package demo

import (
	"fmt"
	"image/color"
	"time"

	"golang.org/x/sync/syncmap"

	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"bitbucket.org/oakmoundstudio/oak/render"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/pointLoc/bruteForce"
	"github.com/200sc/go-compgeo/dcel/pointLoc/kirkpatrick"
	"github.com/200sc/go-compgeo/dcel/pointLoc/slab"
	"github.com/200sc/go-compgeo/dcel/pointLoc/trapezoid"
	"github.com/200sc/go-compgeo/search/tree"
)

var (
	prevVertExisted   bool
	prevVert          *dcel.Vertex
	firstAddedPoint   *dcel.Vertex
	firstAddedExisted bool
	prevEdge          *dcel.Edge
	addedFace         *dcel.Face
)

func addFace(cID int, ev interface{}) int {
	phd := event.GetEntity(cID).(*InteractivePolyhedron)
	me := ev.(mouse.Event)
	if me.X < 0 || me.Y < 0 || me.X > 515 {
		return 0
	}
	mx := float64(me.X) - phd.X
	my := float64(me.Y) - phd.Y
	if me.Button == "LeftMouse" {
		// Detect clicks
		// On first click, declare the first point and edge and face
		// First off, assume that this point is new
		if mode == ADD_DCEL {
			// Need a real copy function
			//undoPhd = append(undoPhd, *phd)

			faceVertices = &syncmap.Map{}
			// Add Case A: The first point of the face
			// already exists.
			hits := mouse.Hits(me.ToSpace())
			if len(hits) > 0 {
				ip := event.GetEntity(int(hits[0].CID)).(*InteractivePoint)
				firstAddedPoint = ip.Vertex
				mouseZ = firstAddedPoint.Z()
				prevVertExisted = true
				firstAddedExisted = true
			} else {
				firstAddedPoint = dcel.NewVertex(mx, my, mouseZ)
				phd.Vertices = append(phd.Vertices, firstAddedPoint)
				prevVertExisted = false
				firstAddedExisted = false
				phd.Update()
				phd.UpdateSpaces()
			}
			prevVert = firstAddedPoint
			faceVertices.Store(prevVert, true)

			addedFace = dcel.NewFace()
			phd.Faces = append(phd.Faces, addedFace)

			mode = ADDING_DCEL
			// On following clicks, prev.next = next, next.prev = prev
			//                      next.origin = origin
			//                      next.face = theface
			//                      prev.twin = make a twin at origin
			//                      twin.face = 0 I guess for now
			// Detect final click by right clicking
		} else if mode == ADDING_DCEL {
			hits := mouse.Hits(me.ToSpace())
			// Old Point <- ...
			if len(hits) > 0 {
				ip := event.GetEntity(int(hits[0].CID)).(*InteractivePoint)
				p := ip.Vertex
				mouseZ = p.Z()
				// This vertex already exists in this face.
				if _, ok := faceVertices.Load(p); ok {
					return 0
				}
				// Old Point -> Old Point
				if prevVertExisted {
					// This vertex must be connected to the previous
					// vertex. If it is not then we need to split a
					// face which we can't do with the information we've
					// been given, in the middle of making a new face.
					// Specifically, we don't know which face this new
					// face is supposed to be consuming in the split.
					// Todo: consider the assumption that addedFace is
					// consuming OUTER_FACE. Is this then reasonable?
					consumedEdge := p.EdgeToward(prevVert)
					if consumedEdge == nil {
						panic("Splitting face with add face")
					}
					check := consumedEdge.Next
					if consumedEdge.Face != phd.Faces[dcel.OUTER_FACE] {
						consumedEdge = consumedEdge.Twin
						check = consumedEdge.Prev
					}
					consumedEdge.Face = addedFace

					// If there was a 'T' before this
					if prevEdge != check && prevEdge != nil {
						check.SetNext(prevEdge.Twin)
						// This causes an issue because it assumes the
						// consumed edge follows our directionality. If it
						// doesn't, because we're defining the points clockwise,
						// we'll loop forever down the line.
						consumedEdge.SetPrev(prevEdge)
					}

					prevEdge = consumedEdge

					// New Point -> Old Point
				} else {
					// As NewPoint -> NewPoint, but
					// we make no new point
					e := dcel.NewEdge()
					tw := dcel.NewEdge()
					phd.HalfEdges = append(phd.HalfEdges, e, tw)
					e.SetTwin(tw)
					e.Face = addedFace
					tw.Face = phd.Faces[dcel.OUTER_FACE]
					tw.Origin = p
					e.Origin = prevVert
					if prevVert == firstAddedPoint {
						prevVert.OutEdge = e
					}
					e.SetPrev(prevEdge)
					if prevEdge != nil {
						prevEdge.Twin.SetPrev(tw)
					}
					prevEdge = e
				}
				prevVert = p
				prevVertExisted = true
				// New Point <- ...
			} else {
				p := dcel.NewVertex(mx, my, mouseZ)
				phd.Vertices = append(phd.Vertices, p)
				e := dcel.NewEdge()
				tw := dcel.NewEdge()
				phd.HalfEdges = append(phd.HalfEdges, e, tw)
				e.SetTwin(tw)
				e.Face = addedFace
				tw.Face = phd.Faces[dcel.OUTER_FACE]
				tw.Origin = p
				e.Origin = prevVert
				if prevVert == firstAddedPoint {
					prevVert.OutEdge = e
				}
				p.OutEdge = tw
				// Old Point -> New Point
				if prevVertExisted {
					if prevEdge == nil {
						// This is illegal, we
						// have no idea what edge off
						// of the previous vertex this
						// face should wrap around
						//
						// This actually prevents zero-width
						// spaces at vertices, where two
						// faces do not touch but at
						// a vertex point.
						panic("Illegal vertex added")
					}
					// If the previous point was old, then
					// this new twin wraps around and connects
					// to the previous's next edge, like this:
					//
					// \ -prevEdge.next
					//  \     /
					//   \   / -prevEdge
					//    \ /
					// (prevVert)
					//     |
					// tw- | -e
					//     |
					//     |
					//    (p)
					tw.SetNext(prevEdge.Next)
					// New Point -> New Point
				} else {
					if prevEdge != nil {
						prevEdge.Twin.SetPrev(tw)
					}
				}
				e.SetPrev(prevEdge)
				prevEdge = e
				prevVert = p
				prevVertExisted = false
			}
			if addedFace.Outer == nil {
				addedFace.Outer = firstAddedPoint.OutEdge
			}
			faceVertices.Store(prevVert, true)
			phd.Update()
			phd.UpdateSpaces()
		} else if mode == POINT_LOCATE {
			mode = LOCATING
			visSlider.min++
			mouseModeBtn.SetString(mode.String())
			go func() {
				if locator == nil {
					var err error
					switch pointLocationMode {
					case SLAB_DECOMPOSITION:
						locator, err = slab.Decompose(&phd.DCEL, tree.RedBlack)
					case TRAPEZOID_MAP:
						_, _, locator, err = trapezoid.TrapezoidalMap(&phd.DCEL)
					case KIRKPATRICK_MONOTONE:
						locator, err = kirkpatrick.TriangleTree(&phd.DCEL, kirkpatrick.MONOTONE)
					case KIRKPATRICK_TRAPEZOID:
						locator, err = kirkpatrick.TriangleTree(&phd.DCEL, kirkpatrick.TRAPEZOID)
					case PLUMB_LINE:
						locator = bruteForce.PlumbLine(&phd.DCEL)
					}
					if err != nil {
						fmt.Println(err)
						mode = POINT_LOCATE
						mouseModeBtn.SetString(mode.String())
						visSlider.min--
						return
					}
				}
				modeBtn.SetRenderable(render.NewColorBox(int(modeBtn.W),
					int(modeBtn.H), createdColor))
				modeBtn.SetPos(515, 410)
				modeBtn.R.SetLayer(4)

				f, _ := locator.PointLocate(mx, my)
				if f == phd.Faces[0] || f == nil {
					fmt.Println("Outer/No Face")
				} else {
					faceIndex := phd.ScanFaces(f)
					if faceIndex < 0 {
						mode = POINT_LOCATE
						mouseModeBtn.SetString(mode.String())
						visSlider.min--
						return
					}
					poly := PolygonFromFace(f)
					poly.Fill(color.RGBA{125, 0, 0, 125})
					poly.ShiftX(phd.X)
					poly.ShiftY(phd.Y)
					render.Draw(poly, 10)
					render.UndrawAfter(poly, 1500*time.Millisecond)
				}
				mode = POINT_LOCATE
				mouseModeBtn.SetString(mode.String())
				visSlider.min--
			}()
		}
	} else if me.Button == "RightMouse" {
		if mode == ADDING_DCEL {
			// Special case
			if firstAddedExisted && prevVertExisted {
				// If firstAddedPoint and prevVert are
				// connected, we set first.prev.face to addedface
				// and stop.
				if firstAddedPoint.EdgeToward(prevVert) != nil {
					addedFace.Outer.Prev.Face = addedFace
				} else {
					// Otherwise we try to split faces on the two verts
					// But this won't correct faces properly
					fmt.Println("Connecting verts")
					phd.ConnectVerts(firstAddedPoint, prevVert, addedFace)
				}
			} else {
				firstEdge := addedFace.Outer
				e := dcel.NewEdge()
				tw := dcel.NewEdge()
				phd.HalfEdges = append(phd.HalfEdges, e, tw)
				e.SetTwin(tw)
				e.SetNext(firstEdge)
				tw.SetPrev(firstEdge.Twin)
				e.Face = addedFace
				tw.Face = phd.Faces[dcel.OUTER_FACE]
				tw.Origin = firstAddedPoint
				e.Origin = prevVert
				if prevVertExisted {
					// !firstAddedExisted
					// else covered above
					tw.SetNext(prevEdge.Next)
				} else {
					tw.SetNext(prevEdge.Twin)
				}
				e.SetPrev(prevEdge)
			}

			phd.CorrectDirectionality(addedFace)

			prevEdge = nil
			prevVert = nil
			addedFace = nil
			firstAddedPoint = nil

			mode = ADD_DCEL

			phd.Update()
			phd.UpdateSpaces()
		}
	}
	return 0
}
