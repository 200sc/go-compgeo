package demo

import (
	"log"
	"math"
	"os"
	"path/filepath"

	"bitbucket.org/oakmoundstudio/oak"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/render"
	"github.com/200sc/go-compgeo/dcel"
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

var (
	dragX           float32 = -1
	dragY           float32 = -1
	dragging                = -1
	offFile                 = filepath.Join("data", "A.off")
	mode                    = ROTATE
	loopDemo        bool
	firstAddedPoint int
	prev            *dcel.Edge
	addedFace       *dcel.Face
	mouseZ          = 0.0
	faceVertices    = make(map[*dcel.Point]bool)
	err             error
	mouseStr        *render.IFText
	modeStr         *render.Text
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
			log.Fatal(err)
		}
	}
	phd := new(InteractivePolyhedron)
	phd.Polyhedron = render.NewPolyhedronFromDCEL(dc, 100, 100)
	phd.Polyhedron.Scale(defScale)
	phd.Polyhedron.RotZ(defRotZ)
	phd.Polyhedron.RotY(defRotY)
	phd.Init()
	render.Draw(phd, 2)

	modeStr = render.DefFont().NewText(mode.String(), 3, 40)
	render.Draw(modeStr, 3)

	mouseStr = render.DefFont().NewInterfaceText(
		dcel.Point{0, 0, 0}, 3, 465)
	render.Draw(mouseStr, 3)

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
		offFile = "none"
		loopDemo = false
	})
}