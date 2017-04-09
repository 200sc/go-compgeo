package compgeo

import (
	"fmt"
	"testing"

	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/search/tree"
)

func TestPointLocSquare(t *testing.T) {
	dc := dcel.FourPoint(
		dcel.Point{0, 0, 0},
		dcel.Point{10, 1, 0},
		dcel.Point{11, 11, 0},
		dcel.Point{1, 10, 0},
	)
	sd, err := dc.SlabDecompose(tree.RedBlack)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sd.(*dcel.SlabPointLocator))
}
