package main

import (
	"log"
	"path/filepath"

	"github.com/200sc/go-compgeo/dcel"

	"bitbucket.org/oakmoundstudio/oak"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"bitbucket.org/oakmoundstudio/oak/render"
)

const (
	shiftSpeed = 3
	scaleSpeed = .02
)

var (
	dragX float32 = -1
	dragY float32 = -1
)

func main() {
	err := oak.LoadConf("oak.config")
	if err != nil {
		log.Fatal(err)
	}
	oak.AddScene("demo",
		func(prevScene string, data interface{}) {
			//phd := render.NewCuboid(100, 100, 100, 100, 100, 100)
			dc, err := dcel.LoadOFF(filepath.Join("data", "A.off"))
			if err != nil {
				log.Fatal(err)
			}
			phd := render.NewPolyhedronFromDCEL(dc, 100, 100)
			render.Draw(phd, 1)
			event.GlobalBind(func(no int, nothing interface{}) int {
				shft := oak.IsDown("LeftShift")
				if oak.IsDown("LeftArrow") {
					phd.ShiftX(-shiftSpeed)
				} else if oak.IsDown("RightArrow") {
					phd.ShiftX(shiftSpeed)
				}
				if oak.IsDown("UpArrow") {
					if shft {
						phd.Scale(1 + scaleSpeed)
					} else {
						phd.ShiftY(-shiftSpeed)
					}
				} else if oak.IsDown("DownArrow") {
					if shft {
						phd.Scale(1 - scaleSpeed)
					} else {
						phd.ShiftY(shiftSpeed)
					}
				}
				if oak.IsDown("LeftMouse") {
					nme := mouse.LastMouseEvent
					if dragX != -1 {
						dx := float64(nme.X - dragX)
						dy := float64(nme.Y - dragY)
						if dx != 0 {
							if shft {
								phd.RotZ(.01 * dx)
							} else {
								phd.RotY(.01 * dx)
							}
						}
						if dy != 0 {
							phd.RotX(.01 * dy)
						}
					}
					dragX = nme.X
					dragY = nme.Y
				} else {
					dragX = -1
					dragY = -1
				}
				//phd.RotZ(.01)
				//phd.RotX(.02)
				//phd.RotY(-.005)
				return 0
			}, "EnterFrame")
			// event.GlobalBind(func(no int, me interface{}) int {
			// 	event := me.(mouse.MouseEvent)
			// 	if event.Button == "LeftMouse" {
			// 		fmt.Println(event.X, event.Y, event.Button)
			// 		dragX = event.X
			// 		dragY = event.Y
			// 	}
			// 	return 0
			// }, "MousePress")
			// event.GlobalBind(func(no int, me interface{}) int {
			// 	event := me.(mouse.MouseEvent)
			// 	if event.Button == "LeftMouse" {
			// 		fmt.Println(event.X, event.Y, event.Button)
			// 		dragX = -1
			// 		dragY = -1
			// 	}
			// 	return 0
			// }, "MouseRelease")
		},
		func() bool {
			return true
		},
		func() (string, *oak.SceneResult) {
			return "demo", nil
		},
	)
	oak.Init("demo")
}
