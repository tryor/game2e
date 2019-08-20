package game2e

import (
	. "tryor/game2e/util"

	. "github.com/tryor/eui"
)

type Ellipse struct {
	*Widget
}

func NewEllipse(x, y, w, h int) *Ellipse {
	s := &Ellipse{Widget: NewWidget()}
	s.Self = s
	s.SetCoordinate(x, y)
	s.SetWidth(w)
	s.SetHeight(h)

	s.SetType(ElementTypeEllipse)
	s.Border = true

	return s
}

func (this *Ellipse) Intersects(x, y int) (ret bool) {
	if !this.Widget.Intersects(x, y) {
		return false
	}
	tx, ty := this.Self.(IElement).GetWorldCoordinate()
	a, b := this.Width()/2, this.Height()/2
	//	cx, cy := a, b
	//	x, y = x-tx, y-ty
	return PointInEllipse(x-tx, y-ty, a, b, a, b)
}
