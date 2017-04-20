package visualize

import (
	"fmt"
	"image/color"

	"github.com/200sc/go-compgeo/geom"

	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
)

var (
	HighlightColor = color.RGBA{255, 255, 255, 255}
	HighlightLayer = 10
)

type Visual struct {
	render.Renderable
	Layer int
}

func DrawLine(vc chan *Visual, p1, p2 geom.D2) {
	v := new(Visual)
	v.Renderable = render.NewThickLine(p1.X(), p1.Y(), p2.X(), p2.Y(), HighlightColor, 1)
	v.Layer = HighlightLayer
	vc <- v
}

func DrawVerticalLine(vc chan *Visual, p geom.D2) {
	v := new(Visual)
	y1 := p.Y() - 480
	y2 := p.Y() + 480
	v.Renderable = render.NewThickLine(p.X(), y1, p.X(), y2, HighlightColor, 1)
	v.Layer = HighlightLayer
	vc <- v
}

func DrawPoly(vc chan *Visual, ps []physics.Vector) {
	v := new(Visual)
	var err error
	v.Renderable, err = render.NewPolygon(ps)
	if err != nil {
		fmt.Println(err)
		return
	}
	v.Renderable.(*render.Polygon).Fill(HighlightColor)

	v.Layer = HighlightLayer
	vc <- v
}
