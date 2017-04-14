package demo

import (
	"bitbucket.org/oakmoundstudio/oak"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"github.com/200sc/go-compgeo/dcel"
)

func phdEnter(cID int, nothing interface{}) int {
	phd := event.GetEntity(cID).(*InteractivePolyhedron)
	shft := oak.IsDown("LeftShift")
	if oak.IsDown("LeftArrow") {
		phd.ShiftX(-shiftSpeed)
		phd.UpdateSpaces()
	} else if oak.IsDown("RightArrow") {
		phd.ShiftX(shiftSpeed)
		phd.UpdateSpaces()
	}
	if oak.IsDown("UpArrow") {
		if shft {
			phd.Scale(1 + scaleSpeed)
			phd.UpdateSpaces()
		} else {
			phd.ShiftY(-shiftSpeed)
			phd.UpdateSpaces()
		}
	} else if oak.IsDown("DownArrow") {
		if shft {
			phd.Scale(1 - scaleSpeed)
			phd.UpdateSpaces()
		} else {
			phd.ShiftY(shiftSpeed)
			phd.UpdateSpaces()
		}
	}
	nme := mouse.LastMouseEvent
	mouseStr.SetText(dcel.Point{float64(nme.X) - phd.X,
		float64(nme.Y) - phd.Y, mouseZ})
	if mode == ROTATE || ((mode == ADD_DCEL || mode == ADDING_DCEL) && shft) {
		if dragX != -1 {
			dx := float64(nme.X - dragX)
			dy := float64(nme.Y - dragY)
			if dx != 0 {
				if shft {
					phd.RotZ(rotSpeed * dx)
					phd.UpdateSpaces()
				} else {
					phd.RotY(rotSpeed * dx)
					phd.UpdateSpaces()
				}
			}
			if dy != 0 {
				phd.RotX(rotSpeed * dy)
				phd.UpdateSpaces()
			}
		}
	} else if mode == MOVE_POINT && dragging != -1 {
		update := false
		mouseZ = phd.Vertices[dragging][2]
		if dragX != -1 {
			phd.Vertices[dragging][0] = float64(dragX) - phd.X
			update = true
		}
		if dragY != -1 {
			phd.Vertices[dragging][1] = float64(dragY) - phd.Y
			update = true
		}
		if oak.IsDown("D") {
			phd.Vertices[dragging][2] += zMoveSpeed
			update = true
		} else if oak.IsDown("C") {
			phd.Vertices[dragging][2] -= zMoveSpeed
			update = true
		}
		if update {
			phd.Update()
			phd.UpdateSpaces()
		}
	}
	if mode != MOVE_POINT {
		if oak.IsDown("D") {
			mouseZ += zMoveSpeed
		} else if oak.IsDown("C") {
			mouseZ -= zMoveSpeed
		}
	}
	if oak.IsDown("LeftMouse") {
		dragX = nme.X
		dragY = nme.Y
	} else {
		dragX = -1
		dragY = -1
	}
	return 0
}
