package printutil

import (
	"fmt"

	"github.com/200sc/go-compgeo/search"
)

// String converts any search.Comparable to
// an appropriate string representation,
// if possible.
func String(f search.Comparable) string {
	switch f2 := f.(type) {
	case fmt.Stringer:
		return f2.String()
	}
	return "?"
}
