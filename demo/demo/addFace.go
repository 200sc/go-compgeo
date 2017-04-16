package demo

import (
	"fmt"
	"image/color"
	"time"

	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"bitbucket.org/oakmoundstudio/oak/timing"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/search/tree"
)

func addFace(cID int, ev interface{}) int {
	phd := event.GetEntity(cID).(*InteractivePolyhedron)
	me := ev.(mouse.MouseEvent)
	mx := float64(me.X) - phd.X
	my := float64(me.Y) - phd.Y
	if me.Button == "LeftMouse" {
		// Detect clicks
		// On first click, declare the first point and edge and face
		// First off, assume that this point is new
		if mode == ADD_DCEL {
			// Need a real copy function
			//undoPhd = append(undoPhd, *phd)

			faceVertices = make(map[*dcel.Point]bool)
			// Add Case A: The first point of the face
			// already exists.
			hits := mouse.Hits(me.ToSpace())
			if len(hits) > 0 {
				ip := event.GetEntity(int(hits[0].CID)).(*InteractivePoint)
				p := ip.Point
				faceVertices[p] = true
				mouseZ = p.Z()
				// Todo: consider if we should change how outEdges
				// are handled so this linear scan doesn't need to
				// happen
				firstAddedPoint = phd.ScanPoints(p)
			} else {
				firstAddedPoint = len(phd.Vertices)
				phd.Vertices = append(phd.Vertices,
					dcel.NewPoint(mx, my, mouseZ))
				faceVertices[phd.Vertices[len(phd.Vertices)-1]] = true
			}
			phd.HalfEdges = append(phd.HalfEdges, dcel.NewEdge())
			prev = phd.HalfEdges[len(phd.HalfEdges)-1]
			prev.Origin = phd.Vertices[firstAddedPoint]

			f := dcel.NewFace()
			f.Inner = prev
			phd.Faces = append(phd.Faces, f)
			addedFace = phd.Faces[len(phd.Faces)-1]

			prev.Face = addedFace
			if firstAddedPoint == len(phd.Vertices)-1 {
				phd.OutEdges = append(phd.OutEdges, prev)
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
			var vi int
			hits := mouse.Hits(me.ToSpace())
			if len(hits) > 0 {
				ip := event.GetEntity(int(hits[0].CID)).(*InteractivePoint)
				p := ip.Point
				// Add Case F: this point already exists in this
				// face. Reject it.
				if _, ok := faceVertices[p]; ok {
					return 0
				}
				faceVertices[p] = true
				mouseZ = p.Z()
				// Todo: consider if we should change how outEdges
				// are handled so this linear scan doesn't need to
				// happen
				vi = phd.ScanPoints(p)
			} else {
				phd.Vertices = append(phd.Vertices,
					dcel.NewPoint(mx, my, mouseZ))
				vi = len(phd.Vertices) - 1
				faceVertices[phd.Vertices[vi]] = true
			}
			// This twin points from the new point to the previous point.
			phd.HalfEdges = append(phd.HalfEdges, dcel.NewEdge())
			twin := phd.HalfEdges[len(phd.HalfEdges)-1]

			twin.Origin = phd.Vertices[vi]

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
			next.Origin = phd.Vertices[vi]
			next.Prev = prev
			next.Face = addedFace
			prev.Next = next

			prev = next

			if vi == len(phd.Vertices)-1 {
				phd.OutEdges = append(phd.OutEdges, prev)
			}

			phd.Update()
			phd.UpdateSpaces()
		} else if mode == POINT_LOCATE {
			sd, err := phd.SlabDecompose(tree.RedBlack)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(sd.(*dcel.SlabPointLocator))
			}
			f, _ := sd.PointLocate(mx, my)
			if f == phd.Faces[0] || f == nil {
				fmt.Println("Outer/No Face")
			} else {
				faceIndex := phd.ScanFaces(f)
				phd.FaceColors[faceIndex] = color.RGBA{255, 0, 0, 255}
				timing.DoAfter(50*time.Millisecond, func() {
					phd.Update()
				})
				timing.DoAfter(2500*time.Millisecond, func() {
					phd.FaceColors[faceIndex] = color.RGBA{0, 255, 255, 255}
					phd.Update()
				})
			}
		}
	} else if me.Button == "RightMouse" {
		if mode == ADDING_DCEL {
			first := addedFace.Inner
			prev.Next = first
			first.Prev = prev

			phd.HalfEdges = append(phd.HalfEdges, dcel.NewEdge())
			// The final twin starts at the first point of this face
			twin := phd.HalfEdges[len(phd.HalfEdges)-1]
			twin.Origin = phd.Vertices[firstAddedPoint]
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
			firstAddedPoint = -1

			mode = ADD_DCEL

			phd.Update()
			phd.UpdateSpaces()

			fmt.Println(phd.DCEL.HalfEdges)
		}
	}
	return 0
}
