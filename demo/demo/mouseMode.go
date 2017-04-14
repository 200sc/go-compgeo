package demo

type mouseMode int

// control mode constant
const (
	ROTATE mouseMode = iota
	MOVE_POINT
	POINT_LOCATE
	ADD_DCEL
	REM_DCEL
	LAST_MODE
	ADDING_DCEL
)

func (m mouseMode) String() string {
	switch m {
	case ROTATE:
		return "Rotate"
	case MOVE_POINT:
		return "Move Point"
	case POINT_LOCATE:
		return "Point Location"
	case ADD_DCEL:
		return "Define Face"
	case ADDING_DCEL:
		return "Defining Face..."
	case REM_DCEL:
		return "Define Inside Face"
	default:
		return "INVALID"
	}
}
