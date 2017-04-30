// Demo is a Point Location visualization demo.

package demo

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/sync/syncmap"

	"bitbucket.org/oakmoundstudio/oak"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/render"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/off"
	"github.com/200sc/go-compgeo/dcel/pointLoc"
	"github.com/200sc/go-compgeo/dcel/pointLoc/visualize"
	"github.com/200sc/go-compgeo/geom"
)

const (
	zMoveSpeed    = 1
	shiftSpeed    = 3
	scaleSpeed    = .02
	rotSpeed      = .01
	vCollisionDim = 8
	defScale      = 20
	defRotZ       = math.Pi
	defRotY       = math.Pi
	defShiftX     = 200
	defShiftY     = 200
)

// Point Location Mode Const
const (
	SLAB_DECOMPOSITION = iota
	TRAPEZOID_MAP
	PLUMB_LINE
	LAST_PL_MODE
	KIRKPATRICK_MONOTONE
	KIRKPATRICK_TRAPEZOID
)

var (
	dragX             float64 = -1
	dragY             float64 = -1
	dragging                  = -1
	offFile                   = filepath.Join("data", "test3.off")
	mode                      = ROTATE
	loopDemo          bool
	mouseZ            = 0.0
	faceVertices      = &syncmap.Map{}
	err               error
	mouseStr          *render.IFText
	font              *render.Font
	phd               *InteractivePolyhedron
	undoPhd           []InteractivePolyhedron
	ticker            *DynamicTicker
	stopTickerCh      = make(chan bool)
	sliding           bool
	locator           pointLoc.LocatesPoints
	pointLocationMode = SLAB_DECOMPOSITION
	modeBtn           *Button
	mouseModeBtn      *Button
	locating          bool
	btnColor          = color.RGBA{50, 50, 140, 255}
	createdColor      = color.RGBA{50, 140, 50, 255}
	visSlider         *Slider

	randomize           = true
	randomSplits        = 1
	defaultRandomSplits = 5
)

// InitScene is called whenever the scene 'demo' starts.
// it creates the objects in our application.
func InitScene(prevScene string, data interface{}) {
	if visualize.VisualCh != nil {
		close(visualize.VisualCh)
		select {
		case stopTickerCh <- true:
		default:
		}
		visualize.VisualCh = nil
	}
	ticker = NewDynamicTicker()
	loopDemo = true
	//phd := render.NewCuboid(100, 100, 100, 100, 100, 100)
	var dc *dcel.DCEL
	if randomize {
		dc = dcel.Random2DDCEL(100, randomSplits)
	} else if offFile == "none" {
		dc = dcel.New()
	} else {
		dc, err = off.Load(offFile)
		if err != nil {
			fmt.Println("Unable to load", offFile, ":", err)
			dc = dcel.New()
		}
	}
	mode = ROTATE
	pointLocationMode = SLAB_DECOMPOSITION
	phd = new(InteractivePolyhedron)
	phd.Polyhedron = NewPolyhedronFromDCEL(dc, defShiftX, defShiftY)
	if !randomize {
		phd.Polyhedron.Scale(defScale)
		phd.Polyhedron.RotZ(defRotZ)
		phd.Polyhedron.RotY(defRotY)
	} else {
		randomize = false
	}
	phd.Init()
	render.Draw(phd, 0)

	fg := render.FontGenerator{File: "luxisr.ttf", Color: render.FontColor("white"), Size: 12}
	font = fg.Generate()

	mouseStr = font.NewInterfaceText(
		geom.Point{0, 0, 0}, 3, 465)
	render.Draw(mouseStr, 3)

	bkgrnd := render.NewColorBox(140, 480, color.RGBA{50, 50, 80, 255})
	bkgrnd.SetPos(514, 0)
	render.Draw(bkgrnd, 0)

	clrBtn := NewButton(clear, font)
	clrBtn.SetLogicDim(50, 20)
	clrBtn.SetRenderable(render.NewColorBox(int(clrBtn.W), int(clrBtn.H), btnColor))
	clrBtn.SetPos(560, 10)
	clrBtn.TxtX = 10
	clrBtn.TxtY = 5
	clrBtn.Layer = 4
	clrBtn.R.SetLayer(4)
	clrBtn.SetString("Clear")

	stepBtn := NewButton(step, font)
	stepBtn.SetLogicDim(70, 20)
	stepBtn.SetRenderable(render.NewColorBox(int(stepBtn.W), int(stepBtn.H), btnColor))
	stepBtn.SetPos(560, 350)
	stepBtn.TxtX = 5
	stepBtn.TxtY = 5
	stepBtn.Layer = 4
	stepBtn.R.SetLayer(4)
	stepBtn.SetString("Step")

	modeBtn = NewButton(changeMode, font)
	modeBtn.SetLogicDim(115, 20)
	modeBtn.SetRenderable(render.NewColorBox(int(modeBtn.W), int(modeBtn.H), btnColor))
	modeBtn.SetPos(515, 410)
	modeBtn.TxtX = 5
	modeBtn.TxtY = 5
	modeBtn.Layer = 4
	modeBtn.R.SetLayer(4)
	modeBtn.SetString("Slab Decomposition")

	visSlider = NewSlider(4, font)
	visSlider.SetDim(115, 35)
	visSlider.SetRenderable(
		render.NewColorBox(int(visSlider.W), int(visSlider.H), btnColor))
	visSlider.SetPos(515, 440)
	visSlider.TxtX = 10
	visSlider.TxtY = 20
	visSlider.R.SetLayer(4)
	visSlider.SetString("No Visualization")

	mouseModeBtn = NewButton(changeMouseMode, font)
	mouseModeBtn.SetLogicDim(90, 20)
	mouseModeBtn.SetRenderable(render.NewColorBox(int(mouseModeBtn.W), int(mouseModeBtn.H), btnColor))
	mouseModeBtn.SetPos(540, 380)
	mouseModeBtn.TxtX = 5
	mouseModeBtn.TxtY = 5
	mouseModeBtn.Layer = 4
	mouseModeBtn.R.SetLayer(4)
	mouseModeBtn.SetString("Rotate")

	event.GlobalBind(clear, "Clear")
	event.GlobalBind(visuals, "Visualize")
	event.GlobalBind(vertexStopDrag, "MouseRelease")

	// Bind mode setting buttons
	keys := []string{"1", "2", "3", "4"}
	for i, k := range keys {
		j := mouseMode(i)
		event.GlobalBind(func(no int, nothing interface{}) int {
			if mode == LOCATING || mode == ADDING_DCEL {
				return 0
			}
			mode = j
			mouseModeBtn.SetString(mode.String())
			return 0
		}, "KeyDown"+k)
	}

	phd.cID.Bind(phdEnter, "EnterFrame")
	phd.cID.Bind(addFace, "MouseRelease")
}

// LoopScene is a basic scene-loop function,
// returning the value of some boolean defined
// in this oak project package.
// When loopDemo is false, the scene will stop
// (and then immediately reset, as it is defined
// to be followed by itself).
func LoopScene() bool {
	return loopDemo
}

// AddCommands opens up some command line functions
// to the application.
func AddCommands() {
	args := os.Args[1:]
	if len(args) > 0 {
		offFile = args[0]
	}
	oak.AddCommand("load", func(strs []string) {
		if mode != LOCATING {
			if len(strs) > 1 {
				offFile = strs[1]
				loopDemo = false
			}
		}
	})
	oak.AddCommand("random", func(strs []string) {
		if mode != LOCATING {
			loopDemo = false
			randomize = true
			if len(strs) > 1 {
				randomSplits, err = strconv.Atoi(strs[1])
				if err != nil {
					randomSplits = defaultRandomSplits
					fmt.Println(err)
				}
			} else {
				randomSplits = defaultRandomSplits
			}
		}
	})
	oak.AddCommand("reset", func(strs []string) {
		if mode != LOCATING {
			loopDemo = false
		}
	})
	oak.AddCommand("clear", func(strs []string) {
		event.Trigger("Clear", nil)
	})
	oak.AddCommand("visualize", func(strs []string) {
		if len(strs) > 1 {
			rate, _ := time.ParseDuration(strs[1])
			event.Trigger("Visualize", rate)
		}
	})
	oak.AddCommand("print", func(strs []string) {
		fmt.Println(phd.DCEL.String())
	})
	oak.AddCommand("save", func(strs []string) {
		if len(strs) < 2 {
			fmt.Println("usage: c save <filepath>")
			return
		}
		err := off.Save(&phd.DCEL).WriteFile(strs[1])
		if err != nil {
			fmt.Println("Error in write: ", err)
		}
	})
}

func clear(no int, nothing interface{}) int {
	if mode != LOCATING {
		offFile = "none"
		loopDemo = false
	}
	return 0
}

func visuals(no int, rt interface{}) int {
	rate := rt.(time.Duration)
	if rate != 0 {
		if visualize.VisualCh == nil {
			visualize.VisualCh = make(chan *visualize.Visual)
		}
		select {
		case stopTickerCh <- true:
		default:
		}
		ticker.SetTick(rate)
		go func() {
			var visual *visualize.Visual
			for {
				select {
				case <-stopTickerCh:
					return
				case <-ticker.ch:
					if visual != nil {
						render.UndrawAfter(visual, 100*time.Millisecond)
					}
					visual = <-visualize.VisualCh
					if visual == nil {
						return
					}
					visual.ShiftX(phd.X)
					visual.ShiftY(phd.Y)

					render.Draw(visual.Renderable, visual.Layer)
					render.UndrawAfter(visual, 2000*time.Millisecond)
				}
			}
		}()
	} else {
		if visualize.VisualCh != nil {
			close(visualize.VisualCh)
			select {
			case stopTickerCh <- true:
			default:
			}
			visualize.VisualCh = nil
		}
	}
	return 0
}

func changeMode(no int, nothing interface{}) int {
	pointLocationMode = (pointLocationMode + 1) % LAST_PL_MODE
	switch pointLocationMode {
	case TRAPEZOID_MAP:
		modeBtn.SetString("Trapezoidal Map")
	case SLAB_DECOMPOSITION:
		modeBtn.SetString("Slab Decomposition")
	case KIRKPATRICK_MONOTONE:
		modeBtn.SetString("Kirkpatrick (mono)")
	case KIRKPATRICK_TRAPEZOID:
		modeBtn.SetString("Kirkpatrick (trap)")
	case PLUMB_LINE:
		modeBtn.SetString("Plumb Line")
	}
	if locator != nil {
		locator = nil
		modeBtn.SetRenderable(render.NewColorBox(int(modeBtn.W), int(modeBtn.H), btnColor))
		modeBtn.SetPos(515, 410)
		modeBtn.R.SetLayer(4)
	}
	return 0
}

func changeMouseMode(no int, nothing interface{}) int {
	if mode == LOCATING || mode == ADDING_DCEL {
		return 0
	}
	mode = (mode + 1) % LAST_MODE
	mouseModeBtn.SetString(mode.String())
	return 0
}

func step(no int, nothing interface{}) int {
	if mode == LOCATING {
		ticker.Step()
	}
	return 0
}
