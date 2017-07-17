package demo

import (
	"fmt"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/render"
)

// A Button is a UI element that has transient background,
// overlayed text, collision space, and binds to mouse events.
type Button struct {
	entities.Solid
	Text       *render.Text
	TxtX, TxtY float64
	Font       *render.Font
	Layer      int
}

// NewButton returns a new button rendering with a given font,
// bound to a given bindable function on being clicked.
func NewButton(bndb event.Bindable, f *render.Font) *Button {
	b := new(Button)
	b.Solid = entities.NewSolid(0, 0, 1, 1, render.EmptyRenderable(), 0)
	b.CID.Bind(bndb, "MouseReleaseOn")
	b.Font = f

	return b
}

// SetSpace overwrites entities.Solid,
// pointing this button to use the mouse collision Rtree
// instead of the entity collision space.
func (b *Button) SetSpace(sp *collision.Space) {
	if b.Space != nil {
		mouse.Remove(b.Space)
	}
	b.Space = sp
	mouse.Add(b.Space)
}

// SetPos acts as SetSpace does, overwriting entities.Solid.
func (b *Button) SetPos(x float64, y float64) {
	b.SetLogicPos(x, y)
	if b.R != nil {
		b.R.SetPos(x, y)
	}

	if b.Space != nil {
		mouse.UpdateSpace(b.X(), b.Y(), b.W, b.H, b.Space)
	}
}

// a stringer is just a string with a function to convert it to
// a string which lets it satisfy the fmt.Stringer interface.
type stringer string

func (s stringer) String() string {
	return string(s)
}

// SetString converts input strings into stringers.
func (b *Button) SetString(txt string) {
	b.SetText(stringer(txt))
}

// SetText changes the text on this button to be the input txt.
func (b *Button) SetText(txt fmt.Stringer) {
	if b.Text != nil {
		b.Text.UnDraw()
	}
	b.Text = b.Font.NewText(txt, b.X()+b.TxtX, b.Y()-b.TxtY+b.H)
	render.Draw(b.Text, b.Layer+1)
}
