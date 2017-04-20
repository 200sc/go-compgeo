package geom

type Vector interface {
	Lesser(dimension int) D3
	Greater(dimension int) D3
	// Renaming of Lesser and Greater
	// for the supplied X dimension
	Left() D3
	Right() D3
	// As above, for Y
	Top() D3
	Bottom() D3
	// As Above, for Z
	Inner() D3
	Outer() D3
}
