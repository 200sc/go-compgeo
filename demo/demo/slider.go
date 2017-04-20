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

type Slider struct {
	*Button
	min, val, max float64
	interval      time.Duration
	knub          render.Renderable
	knubLine      render.Renderable
}

func (sl *Slider) Init() event.CID {
	cID := event.NextID(sl)
	return cID
}

func NewSlider(f *render.Font) *Slider {
	sl := new(Slider)
	sl.Button = new(Button)
	sl.min = 2
	sl.max = 100
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

func sliderDragStart(sl int, nothing interface{}) int {
	sliding = true
	event.CID(sl).Bind(sliderDrag, "EnterFrame")
	return 0
}

func sliderDrag(sl int, nothing interface{}) int {
	slider := event.GetEntity(sl).(*Slider)
	me := mouse.LastMouseEvent
	if me.Event == "MouseRelease" {
		event.Trigger("Visualize", slider.interval)
		sliding = false
		fmt.Println("WHat")
		return event.UNBIND_SINGLE
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

func (sl *Slider) valText() fmt.Stringer {
	if sl.val == 2 {
		sl.interval = time.Duration(0)
		return stringer("No Visualization")
	}
	// We'd like--
	// 1 == 5 ms
	// 100 == 1 second
	// y = x^1.5 maps to this well

	scaled := math.Pow(sl.val, 1.5)
	sl.interval = time.Duration(scaled) * time.Millisecond
	return stringer(strconv.Itoa(int(scaled)) + " ms")
}
