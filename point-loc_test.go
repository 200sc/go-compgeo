package compgeo

import (
	"fmt"
	"testing"

	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/search/tree"
)

func TestPointLocSquare(t *testing.T) {
	dc := dcel.Rect(0, 0, 10, 10)
	sd, err := dc.SlabDecompose(tree.RedBlack)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sd.(*dcel.SlabPointLocator))
}
