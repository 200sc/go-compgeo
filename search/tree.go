package search

// Tree types are Modifiable types that can be
// converted to PersistentTrees. This may change as this package's
// understanding of what types are useful to be made Persistent.
type Tree interface {
	Modifiable
	ToPersistent() PersistentTree
}

// PersistentTree types have a concept of time instants elapsing
// as modifications are made on them, and are otherwise searchable.
// To search a given instant of a PersistentTree, either SetInstant()
// and then search, or search pst.AtInstant(...)
type PersistentTree interface {
	Modifiable
	AtInstant(float64) Modifiable
	MinInstant() float64
	MaxInstant() float64
	SetInstant(float64)
}
