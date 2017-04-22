package visualize

import (
	"fmt"
	"image/color"

	"github.com/200sc/go-compgeo/geom"

	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
)

var (
	// VisualCh is used to determine if
	// this process is currently visualizing anything.
	// if this is nil, the program will not attempt to
	// send visuals.
	VisualCh chan *Visual
	// HighlightColor is the color that (right now)
	// will be assigned to every Visual generated.
	HighlightColor = color.RGBA{255, 255, 255, 255}
	// HighlightLayer is the layer that (right now)
	// will be assigned to every Visual generated.
	HighlightLayer = 10
)

// Visual is a renderable with attached instructions
// to be given to a renderable at time of drawing.
type Visual struct {
	render.Renderable
	Layer int
}

// DrawLine sends a line instruction to the Visual Channel
func DrawLine(p1, p2 geom.D2) {
	v := new(Visual)
	v.Renderable = render.NewThickLine(p1.X(), p1.Y(), p2.X(), p2.Y(), HighlightColor, 1)
	v.Layer = HighlightLayer
	VisualCh <- v
}

// DrawVerticalLine sends a line extending through the screen
// vertically to the visual channel at a given point
func DrawVerticalLine(p geom.D2) {
	v := new(Visual)
	y1 := p.Y() - 480
	y2 := p.Y() + 480
	v.Renderable = render.NewThickLine(p.X(), y1, p.X(), y2, HighlightColor, 1)
	v.Layer = HighlightLayer
	VisualCh <- v
}

// DrawPoly sends a polygon made up of ps (assumed convex)
// to the visual channel
func DrawPoly(ps []physics.Vector) {
	v := new(Visual)
	var err error
	v.Renderable, err = render.NewPolygon(ps)
	if err != nil {
		fmt.Println(err)
		return
	}
	v.Renderable.(*render.Polygon).Fill(HighlightColor)

	v.Layer = HighlightLayer
	VisualCh <- v
}
