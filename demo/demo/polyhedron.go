package demo

import (
	"image"
	"image/color"
	"math"
	"sort"

	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
)

// Polyhedron is a type which extends render.Renderable,
// allowing it to be drawn to screen through oak's drawing
// functionality.
// Polyhedron specifically draws a DCEL, given some defined
// colors for the DCEL's faces and Edges. Vertices are given
// a generic color because they are barely visibile anyway
// with our drawing scheme.
// Polyhedrons are not drawn in a very sophisticated manner.
type Polyhedron struct {
	render.Sprite
	dcel.DCEL
	FaceColors []color.Color
	EdgeColors []color.Color
	Center     physics.Vector
}

var (
	// Default colors
	edgeColor   = color.RGBA{0, 0, 255, 255}
	faceColor   = color.RGBA{0, 150, 150, 255}
	ptColor     = color.RGBA{255, 255, 255, 255}
	zCorrection = 1.0
)

// NewPolyhedronFromDCEL creates a polyhedron from a dcel
// and an initial screen position
func NewPolyhedronFromDCEL(dc *dcel.DCEL, x, y float64) *Polyhedron {
	p := new(Polyhedron)
	p.SetPos(x, y)
	p.DCEL = *dc
	p.Update()
	p.Center = physics.NewVector(p.X+(1+p.MaxX())/2, p.Y+(1+p.MaxY())/2)
	return p
}

// Update keeps a polyhedron's drawn elements consistent
// with changes in the underlying DCEL.
func (p *Polyhedron) Update() {

	p.clearNegativePoints()

	// Reset p's rgba
	maxX := p.MaxX() + 1
	maxY := p.MaxY() + 1
	rect := image.Rect(0, 0, int(maxX), int(maxY))
	rgba := image.NewRGBA(rect)

	// We ignore Z in terms of where exactly we draw anything
	// on screen.
	// There isn't an alternative to this, aside from
	// recoloring things to account for different z values,
	// without having a camera system, which is a lot of work

	// Try to keep the center of this polyhedron to stay in
	// one place on screen. This is not exactly the expected
	// behavior from someone rotating a shape, but it is
	// close.
	if p.Center.X != 0 || p.Center.Y != 0 {
		cx := p.X + maxX/2
		cy := p.Y + maxY/2
		p.X -= (cx - p.Center.X)
		p.Y -= (cy - p.Center.Y)
	}

	// Eventually:
	// For all Faces, Edges, and Vertices, sort by z value
	// and draw them high-to-low
	zi := 0
	zOrder := make([]polyDraw, len(p.HalfEdges)/2+len(p.Faces)-1+len(p.Vertices))
	// I understand that this is not an accurate way of drawing things
	// in 3D space. It happens to be enough to usually draw things
	// in the right order, and as we don't have access to the graphics card,
	// we don't want to spend forever determining an exact draw order or if
	// certain things shouldn't be drawn because they are hidden.
	// Ultimately this is not a job for this specific renderable but for
	// the engine, which is right now only concerned with 2D ordering of
	// elements to draw.

	// Step 1: draw all edges
	// Given the edge twin mandate, we can just use
	// every other halfEdge.

	// If new edges have been added, make sure our
	// edge color slice is long enough to hold
	// colors for the new edges.
	if len(p.EdgeColors) < len(p.HalfEdges)/2 {
		diff := (len(p.HalfEdges) / 2) - len(p.EdgeColors)
		p.EdgeColors = append(p.EdgeColors, make([]color.Color, diff)...)
	}

	for i := 0; i < len(p.HalfEdges); i += 2 {
		points, err := p.FullEdge(i)
		if err != nil {
			continue
		}
		if i/2 >= len(p.EdgeColors) {
			return
		}
		if p.EdgeColors[i/2] == nil {
			p.EdgeColors[i/2] = edgeColor
		}
		zOrder[zi] = coloredEdge{points, p.EdgeColors[i/2]}
		zi++
	}

	// Step 2: draw all vertices
	for _, v := range p.Vertices {
		zOrder[zi] = drawPoint{v}
		zi++
	}

	if len(p.FaceColors) < len(p.Faces) {
		diff := len(p.Faces) - len(p.FaceColors)
		p.FaceColors = append(p.FaceColors, make([]color.Color, diff)...)
	}

	for i := 1; i < len(p.Faces); i++ {
		f := p.Faces[i]
		if p.FaceColors[i] == nil {
			p.FaceColors[i] = faceColor
		}
		verts := f.Vertices()
		maxZ := math.MaxFloat64 * -1
		physVerts := make([]physics.Vector, len(verts))
		for i, v := range verts {
			physVerts[i] = physics.NewVector(v.X(), v.Y())
			if v.Z() > maxZ {
				maxZ = v.Z()
			}
		}

		// We draw each individual face as a Polygon formed of
		// a list of vertices.
		poly, err := render.NewPolygon(physVerts)
		if err != nil {
			continue
		}
		fpoly := facePolygon{
			poly,
			maxZ,
			p.FaceColors[i],
		}

		zOrder[zi] = fpoly
		zi++
	}

	// Sort the elements of zOrder by their Z values.
	sort.Slice(zOrder, func(i, j int) bool {
		if zOrder[i] == nil || zOrder[j] == nil {
			return false
		}
		return zOrder[i].Z() < zOrder[j].Z()
	})

	for _, item := range zOrder {
		if item != nil {
			item.draw(rgba)
		}
	}

	p.SetRGBA(rgba)
}

type polyDraw interface {
	Z() float64
	draw(*image.RGBA)
}

type drawPoint struct {
	*dcel.Vertex
}

func (dp drawPoint) Z() float64 {
	return dp.Vertex.Z() + zCorrection*2
}

func (dp drawPoint) draw(rgba *image.RGBA) {
	rgba.Set(int(dp.Val(0)), int(dp.Val(1)), ptColor)
}

type coloredEdge struct {
	ps geom.FullEdge
	c  color.Color
}

func (ce coloredEdge) Z() float64 {
	return ce.ps.High(2).Val(2) + zCorrection
}

func (ce coloredEdge) draw(rgba *image.RGBA) {
	render.DrawLineOnto(rgba, int(ce.ps[0].X()), int(ce.ps[0].Y()),
		int(ce.ps[1].X()), int(ce.ps[1].Y()), ce.c)
}

type facePolygon struct {
	*render.Polygon
	z float64
	c color.Color
}

func (fp facePolygon) Z() float64 {
	return fp.z
}

func (fp facePolygon) draw(rgba *image.RGBA) {
	for x := fp.Rect.MinX; x < fp.Rect.MaxX; x++ {
		for y := fp.Rect.MinY; y < fp.Rect.MaxY; y++ {
			if fp.Contains(x, y) {
				rgba.Set(int(x), int(y), fp.c)
			}
		}
	}
}

// PolygonFromFace converts a dcelFace into a polygon
// to be rendered
func PolygonFromFace(f *dcel.Face) *render.Polygon {
	verts := f.Vertices()
	maxZ := math.MaxFloat64 * -1
	physVerts := make([]physics.Vector, len(verts))
	for i, v := range verts {
		physVerts[i] = physics.NewVector(v.X(), v.Y())
		if v.Z() > maxZ {
			maxZ = v.Z()
		}
	}
	poly, _ := render.NewPolygon(physVerts)
	return poly
}

// RotZ rotates the polyhedron around the Z axis
func (p *Polyhedron) RotZ(theta float64) {
	st := math.Sin(theta)
	ct := math.Cos(theta)

	for _, v := range p.Vertices {
		v0 := v.X()*ct - v.Y()*st
		v.Point[1] = v.Y()*ct + v.X()*st
		v.Point[0] = v0
	}
	p.Update()
}

// RotX rotates the polyhedron around the X axis
func (p *Polyhedron) RotX(theta float64) {
	st := math.Sin(theta)
	ct := math.Cos(theta)

	for _, v := range p.Vertices {
		v1 := v.Y()*ct - v.Z()*st
		v.Point[2] = v.Z()*ct + v.Y()*st
		v.Point[1] = v1
	}
	p.Update()
}

// RotY rotates the polyhedron around the Y axis
func (p *Polyhedron) RotY(theta float64) {
	st := math.Sin(theta)
	ct := math.Cos(theta)

	for _, v := range p.Vertices {
		v0 := v.X()*ct - v.Z()*st
		v.Point[2] = v.Z()*ct + v.X()*st
		v.Point[0] = v0
	}
	p.Update()
}

// Scale scales up or down the given polyhedron
func (p *Polyhedron) Scale(factor float64) {
	for _, v := range p.Vertices {
		v.Point = geom.Point{
			v.X() * factor,
			v.Y() * factor,
			v.Z() * factor,
		}
	}
	p.Update()
}

func (p *Polyhedron) clearNegativePoints() {
	// Anything with an x,y less than 0 needs to be increased,
	// this is a limitation so we stay in the bounds of a given rgba
	// rectangle on screen, so we increase everything by minX,minY
	x := p.MinX()
	y := p.MinY()
	for _, v := range p.Vertices {
		v.Point[0] = v.X() - x
	}
	for _, v := range p.Vertices {
		v.Point[1] = v.Y() - y
	}
}

// Utilities
func (p *Polyhedron) String() string {
	return "Polyhedron"
}

// ShiftX moves a polyhedron and its center along the x axis
func (p *Polyhedron) ShiftX(x float64) {
	p.Center.X += x
	p.Sprite.ShiftX(x)
}

// ShiftY moves a polyhedron and its center along the y axis
func (p *Polyhedron) ShiftY(y float64) {
	p.Center.Y += y
	p.Sprite.ShiftY(y)
}
