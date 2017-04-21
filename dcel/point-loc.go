package dcel

// LocatesPoints is an interface to represent point location
// queries.
type LocatesPoints interface {
	PointLocate(vs ...float64) (*Face, error)
}
