package demo

import (
	"fmt"

	"golang.org/x/sync/syncmap"

	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/slab"
	"github.com/200sc/go-compgeo/dcel/triangulation"
	"github.com/200sc/go-compgeo/search/tree"
)

func addFace(cID int, ev interface{}) int {
	if sliding {
		return 0
	}
	phd := event.GetEntity(cID).(*InteractivePolyhedron)
	me := ev.(mouse.MouseEvent)
	if me.X < 0 || me.Y < 0 || (me.X > 515 && me.Y > 410) {
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
			f.Inner = prev
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
			if locator == nil {
				var err error
				if pointLocationMode == SLAB_DECOMPOSITION {
					locator, err = slab.Decompose(&phd.DCEL, tree.RedBlack)
				} else if pointLocationMode == TRAPEZOID_MAP {
					var dc *dcel.DCEL
					dc, _, locator, err = triangulation.TrapezoidalMap(&phd.DCEL)
					// for now
					if err == nil {
						phd.DCEL = *dc
						phd.Update()
					} else {
						fmt.Println("error", err)
					}
				}
				if err != nil {
					fmt.Println(err)
					return 0
				}
			}

			// f, _ := locator.PointLocate(mx, my)
			// if f == phd.Faces[0] || f == nil {
			// 	fmt.Println("Outer/No Face")
			// } else {
			// 	faceIndex := phd.ScanFaces(f)
			// 	phd.FaceColors[faceIndex] = color.RGBA{255, 0, 0, 255}

			// 	timing.DoAfter(50*time.Millisecond, func() {
			// 		phd.Update()
			// 	})
			// 	timing.DoAfter(2500*time.Millisecond, func() {
			// 		phd.FaceColors[faceIndex] = color.RGBA{0, 255, 255, 255}
			// 		phd.Update()
			// 	})
			// }
		}
	} else if me.Button == "RightMouse" {
		if mode == ADDING_DCEL {
			first := addedFace.Inner
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

			fmt.Println(phd.DCEL.HalfEdges)
			phd.CorrectDirectionality(addedFace)

			prev = nil
			addedFace = nil
			firstAddedPoint = nil

			mode = ADD_DCEL

			phd.Update()
			phd.UpdateSpaces()

			fmt.Println(phd.DCEL.HalfEdges)
		}
	}
	return 0
}
