package game2e

import (
	"image"

	. "github.com/tryor/game2e/util"

	. "github.com/tryor/eui"
)

type Polygon struct {
	*Widget

	points     []Point
	pointCount int
}

func NewPolygon(x, y int, points ...Point) *Polygon {

	if len(points) < 3 {
		panic("polygon vertices cannot be less than 3")
	}

	p := &Polygon{Widget: NewWidget()}
	p.Self = p
	p.SetCoordinate(x, y)
	p.SetType(ElementTypePolygon)
	p.Border = true
	//		p.Fill = true

	//	p.BorderColor = color.NRGBA{255, 0, 0, 200} //color.NRGBA{255, 0, 0, 200}
	//	p.FillColor = color.NRGBA{0, 177, 199, 200} //color.NRGBA{0, 177, 199, 200}

	p.points = points
	p.pointCount = len(points)

	p.CreateBoundRect()

	//	fmt.Println("NewPolygon", points)

	return p
}

func (this *Polygon) CreatePath() {
	this.Widget.CreatePath()
	x, y := this.Self.(IElement).GetWorldCoordinate()
	if this.pointCount > 2 {
		ps := make([]image.Point, this.pointCount)
		for i, p := range this.points {
			ps[i].X = p.X + x
			ps[i].Y = p.Y + y
		}
		this.BorderPath.AddPolygon(ps)
	}
}

func (this *Polygon) CreateBoundRect() *image.Rectangle {
	p0, p1 := this.points[0], this.points[1]

	x1, y1, x2, y2 := p0.X, p0.Y, p1.X, p1.Y
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	for _, p := range this.points[2:] {
		if x1 > p.X {
			x1 = p.X
		} else if x2 < p.X {
			x2 = p.X
		}

		if y1 > p.Y {
			y1 = p.Y
		} else if y2 < p.Y {
			y2 = p.Y
		}
	}

	//	this.Self.(IElement).SetCoordinate(x1, y1)
	this.Self.(IElement).SetWidth(x2 - x1)
	this.Self.(IElement).SetHeight(y2 - y1)
	return this.Widget.CreateBoundRect()
}

func (this *Polygon) Intersects(x, y int) (ret bool) {
	if !this.Widget.Intersects(x, y) {
		return false
	}

	tx, ty := this.Self.(IElement).GetWorldCoordinate()
	return PointInPolygon(x-tx, y-ty, this.points)
}
