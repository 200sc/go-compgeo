package demo

import (
	"fmt"
	"image/color"
	"time"

	"golang.org/x/sync/syncmap"

	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"bitbucket.org/oakmoundstudio/oak/render"
	"bitbucket.org/oakmoundstudio/oak/timing"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/pointLoc/bruteForce"
	"github.com/200sc/go-compgeo/dcel/pointLoc/kirkpatrick"
	"github.com/200sc/go-compgeo/dcel/pointLoc/slab"
	"github.com/200sc/go-compgeo/dcel/pointLoc/trapezoid"
	"github.com/200sc/go-compgeo/search/tree"
)

func addFace(cID int, ev interface{}) int {
	phd := event.GetEntity(cID).(*InteractivePolyhedron)
	me := ev.(mouse.MouseEvent)
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
				faceVertices.Store(firstAddedPoint, true)
				mouseZ = firstAddedPoint.Z()
			} else {
				firstAddedPoint = dcel.NewVertex(mx, my, mouseZ)
				phd.Vertices = append(phd.Vertices, firstAddedPoint)
				faceVertices.Store(phd.Vertices[len(phd.Vertices)-1], true)
			}
			phd.HalfEdges = append(phd.HalfEdges, dcel.NewEdge())
			prev = phd.HalfEdges[len(phd.HalfEdges)-1]
			prev.Origin = firstAddedPoint

			f := dcel.NewFace()
			f.Outer = prev
			phd.Faces = append(phd.Faces, f)
			addedFace = phd.Faces[len(phd.Faces)-1]

			prev.Face = addedFace
			if firstAddedPoint.OutEdge == nil {
				firstAddedPoint.OutEdge = prev
			}

			mode = ADDING_DCEL

			phd.Update()
			phd.UpdateSpaces()
			// On following clicks, prev.next = next, next.prev = prev
			//                      next.origin = origin
			//                      next.face = theface
			//                      prev.twin = make a twin at origin
			//                      twin.face = 0 I guess for now
			// Detect final click by clicking on first point or
			//                      by right clicking
		} else if mode == ADDING_DCEL {
			// Add Case D:
			// Some node other than the first in the face already
			// exists
			var p *dcel.Vertex
			hits := mouse.Hits(me.ToSpace())
			if len(hits) > 0 {
				ip := event.GetEntity(int(hits[0].CID)).(*InteractivePoint)
				p = ip.Vertex
				// Add Case F: this point already exists in this
				// face. Reject it.
				_, ok := faceVertices.Load(p)
				if ok {
					return 0
				}
				faceVertices.Store(p, true)
				mouseZ = p.Z()
			} else {
				p = dcel.NewVertex(mx, my, mouseZ)
				phd.Vertices = append(phd.Vertices, p)
				faceVertices.Store(p, true)
			}
			// This twin points from the new point to the previous point.
			phd.HalfEdges = append(phd.HalfEdges, dcel.NewEdge())
			twin := phd.HalfEdges[len(phd.HalfEdges)-1]

			twin.Origin = p

			// We make an assumption here that all twins have the outer face
			// as their face.
			twin.Face = phd.Faces[dcel.OUTER_FACE]

			// This should be the twin of the previous edge,
			// and vice versa.
			prev.Twin = twin
			twin.Twin = prev

			// If the previous edge has a previous edge,
			if prev.Prev != nil {
				// This twin's next should be the previous twin,
				// and the previous twin's previous should be this.
				twin.Next = prev.Prev.Twin
				prev.Prev.Twin.Prev = twin
			}

			phd.HalfEdges = append(phd.HalfEdges, dcel.NewEdge())
			next := phd.HalfEdges[len(phd.HalfEdges)-1]
			next.Origin = p
			next.Prev = prev
			next.Face = addedFace
			prev.Next = next

			prev = next

			if p.OutEdge == nil {
				p.OutEdge = prev
			}

			phd.Update()
			phd.UpdateSpaces()
		} else if mode == POINT_LOCATE {
			mode = LOCATING
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
						return
					}
				}
				modeBtn.SetRenderable(render.NewColorBox(int(modeBtn.W),
					int(modeBtn.H), color.RGBA{50, 100, 50, 255}))
				modeBtn.SetPos(515, 410)
				modeBtn.R.SetLayer(4)

				f, _ := locator.PointLocate(mx, my)
				if f == phd.Faces[0] || f == nil {
					fmt.Println("Outer/No Face")
				} else {
					faceIndex := phd.ScanFaces(f)
					if faceIndex < 0 {
						mode = POINT_LOCATE
						return
					}
					poly := PolygonFromFace(f)
					poly.Fill(color.RGBA{125, 0, 0, 125})
					poly.ShiftX(phd.X)
					poly.ShiftY(phd.Y)
					render.Draw(poly, 10)
					render.UndrawAfter(poly, 1500*time.Millisecond)
					phd.FaceColors[faceIndex] = color.RGBA{255, 0, 0, 255}

					go timing.DoAfter(50*time.Millisecond, func() {
						phd.Update()
					})
					go timing.DoAfter(1500*time.Millisecond, func() {
						phd.FaceColors[faceIndex] = faceColor
						phd.Update()
					})
				}
				mode = POINT_LOCATE
			}()
		}
	} else if me.Button == "RightMouse" {
		if mode == ADDING_DCEL {
			first := addedFace.Outer
			prev.Next = first
			first.Prev = prev

			phd.HalfEdges = append(phd.HalfEdges, dcel.NewEdge())
			// The final twin starts at the first point of this face
			twin := phd.HalfEdges[len(phd.HalfEdges)-1]
			twin.Origin = firstAddedPoint
			twin.Face = phd.Faces[dcel.OUTER_FACE]

			prev.Twin = twin
			twin.Twin = prev

			// This twin's next should be the previous twin,
			// and the previous twin's previous should be this.
			twin.Next = prev.Prev.Twin
			prev.Prev.Twin.Prev = twin
			// This twin's previous should be the first edge
			// we added's twin. and vice versa
			twin.Prev = first.Twin
			first.Twin.Next = twin

			phd.CorrectDirectionality(addedFace)

			prev = nil
			addedFace = nil
			firstAddedPoint = nil

			mode = ADD_DCEL

			phd.Update()
			phd.UpdateSpaces()
		}
	}
	return 0
}
