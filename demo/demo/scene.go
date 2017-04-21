package demo

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/200sc/go-compgeo/geom"
	"golang.org/x/sync/syncmap"

	"bitbucket.org/oakmoundstudio/oak"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/render"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/slab"
	"github.com/200sc/go-compgeo/dcel/visualize"
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
)

const (
	SLAB_DECOMPOSITION = iota
	TRAPEZOID_MAP
)

var (
	dragX             float64 = -1
	dragY             float64 = -1
	dragging                  = -1
	offFile                   = filepath.Join("data", "test.off")
	mode                      = ROTATE
	loopDemo          bool
	firstAddedPoint   *dcel.Vertex
	prev              *dcel.Edge
	addedFace         *dcel.Face
	mouseZ            = 0.0
	faceVertices      = &syncmap.Map{}
	err               error
	mouseStr          *render.IFText
	modeStr           *render.Text
	font              *render.Font
	phd               *InteractivePolyhedron
	undoPhd           []InteractivePolyhedron
	ticker            *time.Ticker
	stopTickerCh      = make(chan bool)
	sliding           bool
	locator           dcel.LocatesPoints
	pointLocationMode = SLAB_DECOMPOSITION
	modeBtn           *Button
)

// InitScene is called whenever the scene 'demo' starts.
// it creates the objects in our application.
func InitScene(prevScene string, data interface{}) {
	loopDemo = true
	//phd := render.NewCuboid(100, 100, 100, 100, 100, 100)
	var dc *dcel.DCEL
	if offFile == "none" {
		dc = dcel.New()
	} else {
		dc, err = dcel.LoadOFF(offFile)
		if err != nil {
			fmt.Println("Unable to load", offFile, ":", err)
			dc = dcel.New()
		}
	}
	pointLocationMode = SLAB_DECOMPOSITION
	phd = new(InteractivePolyhedron)
	phd.Polyhedron = NewPolyhedronFromDCEL(dc, 100, 100)
	phd.Polyhedron.Scale(defScale)
	phd.Polyhedron.RotZ(defRotZ)
	phd.Polyhedron.RotY(defRotY)
	phd.Init()
	render.Draw(phd, 2)

	fg := render.FontGenerator{File: "luxisr.ttf", Color: render.FontColor("white"), Size: 12}
	font = fg.Generate()

	modeStr = font.NewText(mode.String(), 3, 40)
	render.Draw(modeStr, 3)

	mouseStr = font.NewInterfaceText(
		geom.Point{0, 0, 0}, 3, 465)
	render.Draw(mouseStr, 3)

	clrBtn := NewButton(clear, font)
	clrBtn.SetLogicDim(70, 20)
	clrBtn.SetRenderable(render.NewColorBox(int(clrBtn.W), int(clrBtn.H), color.RGBA{50, 50, 100, 255}))
	clrBtn.SetPos(560, 10)
	clrBtn.TxtX = 10
	clrBtn.TxtY = 5
	clrBtn.SetString("Clear")

	modeBtn = NewButton(changeMode, font)
	modeBtn.SetLogicDim(115, 20)
	modeBtn.SetRenderable(render.NewColorBox(int(modeBtn.W), int(modeBtn.H), color.RGBA{50, 50, 100, 255}))
	modeBtn.SetPos(515, 410)
	modeBtn.TxtX = 5
	modeBtn.TxtY = 5
	modeBtn.SetString("Slab Decomposition")

	visSlider := NewSlider(font)
	visSlider.SetDim(115, 35)
	visSlider.SetRenderable(
		render.NewColorBox(int(visSlider.W), int(visSlider.H), color.RGBA{50, 50, 100, 255}))
	visSlider.SetPos(515, 440)
	visSlider.TxtX = 10
	visSlider.TxtY = 20
	visSlider.SetString("No Visualization")

	event.GlobalBind(clear, "Clear")
	event.GlobalBind(visuals, "Visualize")
	event.GlobalBind(vertexStopDrag, "MouseRelease")
	event.GlobalBind(func(no int, nothing interface{}) int {
		mode = (mode + 1) % LAST_MODE
		modeStr.SetText(mode.String())
		return 0
	}, "KeyDownQ")
	keys := []string{"1", "2", "3", "4"}
	for i, k := range keys {
		j := mouseMode(i)
		event.GlobalBind(func(no int, nothing interface{}) int {
			mode = j
			modeStr.SetText(mode.String())
			return 0
		}, "KeyDown"+k)
	}
	phd.cID.Bind(phdEnter, "EnterFrame")
	phd.cID.Bind(addFace, "MouseRelease")
	// phd.cID.Bind(func(cID int, nothing interface{}) int {
	// 	if oak.IsDown("LeftControl") && len(undoPhd) != 0 {
	// 		phd := event.GetEntity(cID).(*InteractivePolyhedron)
	// 		*phd = undoPhd[len(undoPhd)-1]
	// 		// Discarding right now,
	// 		// could offer redo later
	// 		undoPhd = undoPhd[len(undoPhd)-1:]
	// 	}
	// 	return 0
	// }, "KeyDownZ")
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
		if len(strs) > 1 {
			offFile = strs[1]
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
}

func clear(no int, nothing interface{}) int {
	offFile = "none"
	loopDemo = false
	return 0
}

func visuals(no int, rt interface{}) int {
	rate := rt.(time.Duration)
	if rate != 0 {
		if ticker != nil {
			close(slab.VisualCh)
			select {
			case stopTickerCh <- true:
			default:
			}
			ticker.Stop()
		}
		slab.VisualCh = make(chan *visualize.Visual)
		ticker = time.NewTicker(rate)
		go func() {
			var visual *visualize.Visual
			for {
				select {
				case <-stopTickerCh:
					return
				case <-ticker.C:
					if visual != nil {
						render.UndrawAfter(visual, 100*time.Millisecond)
					}
					visual = <-slab.VisualCh
					if visual == nil {
						fmt.Println("Nil visual recieved")
						return
					}
					fmt.Println("Drawing visual")
					visual.ShiftX(phd.X)
					visual.ShiftY(phd.Y)

					fmt.Println(visual.GetX(), visual.GetY())

					render.Draw(visual.Renderable, visual.Layer)
				}
			}
		}()
	} else {
		if ticker != nil {
			close(slab.VisualCh)
			select {
			case stopTickerCh <- true:
			default:
			}
			ticker.Stop()
			ticker = nil
		}
		slab.VisualCh = nil
	}

	return 0
}

func changeMode(no int, nothing interface{}) int {
	if pointLocationMode == SLAB_DECOMPOSITION {
		pointLocationMode = TRAPEZOID_MAP
		modeBtn.SetString("Trapezoidal Map")
	} else {
		pointLocationMode = SLAB_DECOMPOSITION
		modeBtn.SetString("Slab Decomposition")
	}
	locator = nil
	return 0
}
