package tree

type PersistentBST struct {
	instant float64

	// Implicitly sorted
	instants []BSTInstant
}

type BSTInstant struct {
	BST
	instant float64
}
