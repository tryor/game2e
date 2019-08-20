package game2e

import (
	"errors"
	"fmt"
	"image"

	"tryor/game2e/log"
	. "tryor/game2e/util"

	. "github.com/tryor/eui"
)

type Map struct {
	MapInfo *MapInfo
	Layers  []ILayer
}

//vx, vy 可视区域坐标
func NewMap(ts *TaskSchedule, vx, vy int, id string) (*Map, error) {
	mapinfo := MapInfos[id]
	if mapinfo == nil {
		err := errors.New(fmt.Sprint("map info not exist, map id is ", id))
		log.Error(err)
		return nil, err
	}
	m := &Map{MapInfo: mapinfo, Layers: make([]ILayer, 0)}

	for _, maplayer := range mapinfo.Layers {
		layer := newLayer(ts, vx, vy, maplayer)
		if layer != nil {
			m.Layers = append(m.Layers, layer)
		}
	}
	return m, nil
}

func newLayer(ts *TaskSchedule, vx, vy int, maplayer *MapLayer) ILayer {
	if maplayer.Width <= 0 || maplayer.Height <= 0 {
		log.Errorf("layer width or height is less than or equal 0, layer id is %v, desc is %v", maplayer.Id, maplayer.Desc)
		return nil
	}

	l := NewLayer(nil, nil)
	if maplayer.Id != "" {
		l.SetId(maplayer.Id)
	}
	l.SetTag(maplayer.Tag)
	l.SetLayerType(maplayer.Type)
	//l.SetDesc(maplayer.Desc)
	l.SetAlignment(maplayer.Alignment)
	l.SetAnchorPoint(maplayer.Anchorpoint.X, maplayer.Anchorpoint.Y)
	l.SetCoordinate(maplayer.X, maplayer.Y)
	l.SetWidth(maplayer.Width)
	l.SetHeight(maplayer.Height)
	l.SetBackground(maplayer.Background)
	l.SetVisibleRegion(image.Rect(vx, vy, vx+maplayer.VWidth, vy+maplayer.VHeight))
	l.SetVisible(!maplayer.Invisible)
	l.SetDrawMode(maplayer.DrawMode)
	l.SetScrollrateX(maplayer.ScrollrateX)
	l.SetScrollrateY(maplayer.ScrollrateY)
	l.SetScrollMode(maplayer.ScrollMode)

	for _, mapel := range maplayer.Children {
		el := NewElementByResource(ts, mapel)
		if el != nil {
			var err error
			//			if mapel.Feature.OrderZ > 0 {
			//				err = l.AddChild(el, mapel.Feature.OrderZ-1)
			//			} else {
			err = l.AddChild(el)
			//			}
			if err != nil {
				log.Error(err)
			}
		}
	}

	return l
}

//func newWidget(ts *TaskSchedule, mapel *ResElementInfo) IWidget {
//	if mapel.Invalid {
//		return nil
//	}
//	var widget IWidget
//	switch mapel.Type {
//	case ElementTypeImage:
//		img := NewCachedImage(mapel.Feature.Bounds.X, mapel.Feature.Bounds.Y, mapel.RefId, mapel.RefId2)
//		if img != nil {
//			widget = img
//		}
//	case ElementTypeAnimation:
//		animn := NewAnimation(ts, mapel.Feature.Bounds.X, mapel.Feature.Bounds.Y, mapel.RefId, true)
//		if animn != nil {
//			widget = animn
//		}
//	case ElementTypeDemon:
//		demon := NewDemon(ts, mapel.Feature.Bounds.X, mapel.Feature.Bounds.Y, mapel.RefId)
//		if demon != nil {
//			demon.Name.SetFontInfo(&mapel.Text)
//			demon.Name.SetFeature(&mapel.Text.Feature)
//			widget = demon
//		}

//	case ElementTypeEllipse:
//		e := NewEllipse(mapel.Feature.Bounds.X, mapel.Feature.Bounds.Y, mapel.Feature.Bounds.W, mapel.Feature.Bounds.H)
//		if e != nil {
//			widget = e
//		}

//	case ElementTypeRectangle:
//		shape := NewWidget()
//		if shape != nil {
//			widget = shape
//		}

//	case ElementTypePolygon:
//		if len(mapel.Points) > 2 {
//			p := NewPolygon(mapel.Feature.Bounds.X, mapel.Feature.Bounds.Y, mapel.Points...)
//			if p != nil {
//				widget = p
//			}
//		} else {
//			log.Errorf("Polygon (%v) points error, points is %v", mapel.Desc, mapel.Points)
//		}
//	case ElementTypeText:
//		text := NewLabel2(0, 0, 0, 0,
//			mapel.Text.Text, mapel.Text.Font, mapel.Text.Size, mapel.Text.Color)
//		text.FontStyle = mapel.Text.Style
//		text.StringFormatFlags = mapel.Text.Format
//		text.TextAlignment = mapel.Text.Textalignment
//		//		text.SetFeature(&mapel.Feature)
//		widget = text
//	}

//	if widget != nil {
//		if mapel.Id != "" {
//			widget.SetId(mapel.Id)
//		}
//		widget.SetFeature(&mapel.Feature)

//		widget.SetType(mapel.Type)
//		widget.SetTag(mapel.Tag)
//		widget.SetDesc(mapel.Desc)
//		widget.SetEventEnabled(mapel.EventEnabled)
//		widget.SetVisible(!mapel.Invisible)
//		widget.SetObstacle(mapel.Obstacle)

//		for _, mapel := range mapel.Elements {
//			el := newWidget(ts, mapel)
//			if el != nil {
//				err := widget.AddChild(el)
//				if err != nil {
//					log.Error(err)
//				}
//			}
//		}

//		return widget
//	}
//	return nil
//}

//func newWidget_old(ts *TaskSchedule, mapel *MapElement) IWidget {
//	if mapel.Invalid {
//		return nil
//	}

//	var widget IWidget
//	switch mapel.Type {
//	case ElementTypeImage:
//		img := NewCachedImage(mapel.X, mapel.Y, mapel.RefId, mapel.RefId2)
//		if img != nil {
//			widget = img
//		}
//	case ElementTypeAnimation:
//		animn := NewAnimation(ts, mapel.X, mapel.Y, mapel.RefId, true)
//		if animn != nil {
//			widget = animn
//		}
//	case ElementTypeDemon:
//		demon := NewDemon(ts, mapel.X, mapel.Y, mapel.RefId)
//		if demon != nil {
//			widget = demon
//		}

//	case ElementTypeEllipse:
//		if len(mapel.Points) > 1 {
//			p1, p2 := mapel.Points[0], mapel.Points[1]
//			e := NewEllipse(p1.X+mapel.X, p1.Y+mapel.Y, p2.X-p1.X, p2.Y-p1.Y)
//			if e != nil {
//				widget = e
//			}
//		} else {
//			log.Errorf("ellipse (%v) points error, points is %v", mapel.Desc, mapel.Points)
//		}

//	case ElementTypeRectangle:
//		if len(mapel.Points) > 1 {
//			shape := NewWidget()
//			if shape != nil {
//				widget = shape
//				p1, p2 := mapel.Points[0], mapel.Points[1]
//				shape.SetCoordinate(p1.X+mapel.X, p1.Y+mapel.Y)
//				shape.SetWidth(p2.X - p1.X)
//				shape.SetHeight(p2.Y - p1.Y)
//			}
//		} else {
//			log.Errorf("rectangle (%v) points error, points is %v", mapel.Desc, mapel.Points)
//		}

//	case ElementTypePolygon:
//		if len(mapel.Points) > 2 {
//			p := NewPolygon(mapel.X, mapel.Y, mapel.Points...)
//			if p != nil {
//				widget = p
//			}
//		} else {
//			log.Errorf("Polygon (%v) points error, points is %v", mapel.Desc, mapel.Points)
//		}

//	}

//	if widget != nil {
//		if mapel.Id != "" {
//			widget.SetId(mapel.Id)
//		}
//		//		fmt.Printf("%f,%f\n", mapel.Anchorpoint.X, mapel.Anchorpoint.Y)
//		widget.SetAlignment(mapel.Alignment)
//		widget.SetAnchorPoint(mapel.Anchorpoint.X, mapel.Anchorpoint.Y)
//		widget.SetType(mapel.Type)
//		widget.SetTag(mapel.Tag)
//		widget.SetDesc(mapel.Desc)
//		widget.SetEventEnabled(mapel.EventEnabled)
//		widget.SetBorder(mapel.Border)
//		if mapel.BorderColor != ZRGBA {
//			widget.SetBorderColor(mapel.BorderColor)
//		}
//		if mapel.BorderWidth > 0 {
//			widget.SetBorderWidth(mapel.BorderWidth)
//		}
//		widget.SetFill(mapel.Fill)
//		if mapel.FillColor != ZRGBA {
//			widget.SetFillColor(mapel.FillColor)
//		}
//		widget.SetVisible(!mapel.Invisible)
//		widget.SetObstacle(mapel.Obstacle)
//		return widget
//	}
//	return nil
//}
