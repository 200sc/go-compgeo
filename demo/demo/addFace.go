package demo

import (
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"github.com/200sc/go-compgeo/dcel"
)

func addFace(cID int, ev interface{}) int {
	phd := event.GetEntity(cID).(*InteractivePolyhedron)
	me := ev.(mouse.MouseEvent)
	if me.Button == "LeftMouse" {
		// Detect clicks
		// On first click, declare the first point and edge and face
		// First off, assume that this point is new
		if mode == ADD_DCEL {
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
					dcel.NewPoint(float64(me.X)-phd.X, float64(me.Y)-phd.Y, mouseZ))
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
					dcel.NewPoint(float64(me.X)-phd.X, float64(me.Y)-phd.Y, mouseZ))
				vi = len(phd.Vertices) - 1
				faceVertices[phd.Vertices[vi]] = true
			}
			phd.HalfEdges = append(phd.HalfEdges, dcel.NewEdge())
			twin := phd.HalfEdges[len(phd.HalfEdges)-1]
			twin.Origin = phd.Vertices[vi]
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
		}
	} else if me.Button == "RightMouse" {
		if mode == ADDING_DCEL {
			prev.Next = addedFace.Inner
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
}
