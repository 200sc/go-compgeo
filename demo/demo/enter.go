package demo

import (
	"bitbucket.org/oakmoundstudio/oak"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"github.com/200sc/go-compgeo/geom"
)

func phdEnter(cID int, nothing interface{}) int {
	if mode == LOCATING {
		return 0
	}

	// Move the polyhedron with the arrow keys
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
	mX := float64(nme.X)
	mY := float64(nme.Y)
	mouseStr.SetText(geom.Point{mX - phd.X, mY - phd.Y, mouseZ})
	if mX < 0 || mY < 0 || mX > 515 {
		dragX = -1
		dragY = -1
		return 0
	}

	// Certain modes allow for rotaing the dcel
	if mode == ROTATE || ((mode == ADD_DCEL || mode == ADDING_DCEL) && shft) {
		if dragX != -1 {
			dx := mX - dragX
			dy := mY - dragY
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

		// The Move point mode allows for changing the position of points.
	} else if mode == MOVE_POINT && dragging != -1 {
		mouseZ = phd.Vertices[dragging].Z()
		update := false
		if dragX != -1 {
			phd.Vertices[dragging].Point =
				phd.Vertices[dragging].Set(0, float64(dragX)-phd.X).(geom.Point)
			update = true
		}
		if dragY != -1 {
			phd.Vertices[dragging].Point =
				phd.Vertices[dragging].Set(1, float64(dragY)-phd.Y).(geom.Point)
			update = true
		}
		if oak.IsDown("D") {
			phd.Vertices[dragging].Add(2, zMoveSpeed)
			update = true
		} else if oak.IsDown("C") {
			phd.Vertices[dragging].Add(2, zMoveSpeed)
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
		dragX = mX
		dragY = mY
	} else {
		dragX = -1
		dragY = -1
	}
	return 0
}
