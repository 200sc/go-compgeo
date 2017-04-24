package demo

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"time"

	"bitbucket.org/oakmoundstudio/oak/collision"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"bitbucket.org/oakmoundstudio/oak/render"
)

// Slider is a little UI element with a movable slider
// representing a range. It has a specific purpose of
// setting visualization delay, but eventually it might
// be expanded to a more generalized structure.
type Slider struct {
	*Button
	min, val, max float64
	interval      time.Duration
	knub          render.Renderable
	knubLine      render.Renderable
}

// Init returns a CID for the Slider.
//
// Note on engine internals:
// All entities as defined by the engine
// need to have this function defined on them.
// this is because an entity is only a meaningful
// concept in terms of the engine for an entity's
// ability to have events bound to it and triggered
// on it, which this CID (caller ID) represents.
//
// Its literal meaning is, in our event bus, the value
// passed into NextID (which is the only way to get a
// legitimate CID), is stored at the array index of the
// returned CID.
func (sl *Slider) Init() event.CID {
	cID := event.NextID(sl)
	return cID
}

// NewSlider returns a slider with initialized values
// using the given font to render its text.
func NewSlider(layer int, f *render.Font) *Slider {
	sl := new(Slider)
	sl.Button = new(Button)
	sl.min = 2
	sl.Layer = layer
	sl.max = 101
	sl.val = 2
	sl.interval = time.Duration(0)
	sl.CID = sl.Init()
	sl.W = 1
	sl.H = 1
	sl.Space = collision.NewSpace(0, 0, 1, 1, sl.CID)

	sl.Font = f
	sl.CID.Bind(sliderDragStart, "MousePressOn")

	sl.knub = render.NewColorBox(5, 15, color.RGBA{240, 100, 100, 255})
	sl.knubLine = render.NewLine(0, 0, 100, 0, color.RGBA{255, 255, 255, 255})
	render.Draw(sl.knub, sl.Layer+2)
	render.Draw(sl.knubLine, sl.Layer+1)
	return sl
}

// SetPos is an overwrite of a lower-tiered function
// which sets this slider's position and the position
// of it's attached entities
func (sl *Slider) SetPos(x float64, y float64) {
	sl.SetLogicPos(x, y)
	if sl.R != nil {
		sl.R.SetPos(x, y)
	}
	// Todo: There's an obvious need to have
	// 'attached' renderables to any renderable
	// with some offsets. Not in the scope of this project.
	if sl.knub != nil {
		sl.knub.SetPos(x+7, y+17)
	}
	if sl.knubLine != nil {
		sl.knubLine.SetPos(x+7, y+25)
	}

	if sl.Space != nil {
		mouse.UpdateSpace(sl.X, sl.Y, sl.W, sl.H, sl.Space)
	}
}

// sliderDragStart tells this demo to ignore some mouse events
// until sliding = false, and binds to every frame sliderDrag.
func sliderDragStart(sl int, nothing interface{}) int {
	if sliding != true {
		sliding = true
		event.CID(sl).Bind(sliderDrag, "EnterFrame")
	}
	return 0
}

// sliderDrag updates the position and value of this slider's
// knub, within a defined range. Once the mouse is let go,
// it allows other mouse operations to resume and updates
// the visualizaton delay to the value it was left at.
func sliderDrag(sl int, nothing interface{}) int {
	slider := event.GetEntity(sl).(*Slider)
	me := mouse.LastMouseEvent
	if me.Event == "MouseRelease" || me.X < 515 {
		event.Trigger("Visualize", slider.interval)
		sliding = false
		return event.UNBIND_EVENT
	}
	x := float64(me.X) - (slider.X + 5)
	if x <= slider.min {
		if slider.val == slider.min {
			return 0
		}
		slider.val = slider.min
	} else if x >= slider.max {
		if slider.val == slider.max {
			return 0
		}
		slider.val = slider.max
	} else {
		slider.val = x
	}
	slider.knub.SetPos(slider.X+slider.val+5, slider.knub.GetY())
	slider.SetText(slider.valText())
	return 0
}

// valText updates slider's text given it's knub position.
func (sl *Slider) valText() fmt.Stringer {
	if sl.val == 2 {
		sl.interval = time.Duration(0)
		return stringer("No Visualization")
	}
	if sl.val == 101 {
		sl.interval = time.Duration(math.MaxInt64)
		return stringer("Stepwise")
	}

	scaled := math.Pow(sl.val, 1.25)
	sl.interval = time.Duration(scaled) * time.Millisecond
	return stringer(strconv.Itoa(int(scaled)) + " ms")
}
