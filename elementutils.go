package game2e

import (
	"errors"

	"github.com/tryor/game2e/log"
	. "github.com/tryor/game2e/util"

	. "github.com/tryor/eui"
)

func NewElementByResourceId(ts *TaskSchedule, id string, reselsmap ...ResElementInfoMap) (IWidget, error) {
	var resmap ResElementInfoMap
	if len(reselsmap) == 0 {
		resmap = ResElementInfos
	} else {
		resmap = reselsmap[0]
	}

	elinfo, ok := resmap[id]
	if !ok || elinfo == nil {
		err := errors.New("resource element is not exist, id is " + id)
		log.Error(err)
		return nil, err
	}
	el := NewElementByResource(ts, elinfo)
	if el != nil {
		return el, nil
	}
	return nil, nil
}

func NewElementByResource(ts *TaskSchedule, elinfo *ResElementInfo, reselsmap ...ResElementInfoMap) IWidget {
	if elinfo.Invalid {
		return nil
	}
	var widget IWidget
	switch elinfo.Type {
	case ElementTypeImage:
		img := NewImageByResource(elinfo.Feature.Bounds.X, elinfo.Feature.Bounds.Y, elinfo.RefId, elinfo.RefId2)
		if img != nil {
			widget = img
		}
	case ElementTypeAnimation:
		animn := NewAnimation(ts, elinfo.Feature.Bounds.X, elinfo.Feature.Bounds.Y, elinfo.RefId, nil, true)
		if animn != nil {
			widget = animn
		}
	case ElementTypeSpirit:
		spirit := NewSpirit(ts, elinfo.Feature.Bounds.X, elinfo.Feature.Bounds.Y, elinfo.RefId)
		if spirit != nil {
			spirit.Name.SetFontInfo(&elinfo.Text)
			spirit.Name.SetFeature(&elinfo.Text.Feature)
			widget = spirit
		}

	case ElementTypeEllipse:
		e := NewEllipse(elinfo.Feature.Bounds.X, elinfo.Feature.Bounds.Y, elinfo.Feature.Bounds.W, elinfo.Feature.Bounds.H)
		if e != nil {
			widget = e
		}

	case ElementTypeRectangle:
		shape := NewWidget()
		if shape != nil {
			widget = shape
		}

	case ElementTypePolygon:
		if len(elinfo.Points) > 2 {
			p := NewPolygon(elinfo.Feature.Bounds.X, elinfo.Feature.Bounds.Y, elinfo.Points...)
			if p != nil {
				widget = p
			}
		} else {
			log.Errorf("Polygon (%v) points error, points is %v", elinfo.Desc, elinfo.Points)
		}
	case ElementTypeText:
		text := NewLabel2(0, 0, 0, 0,
			elinfo.Text.Text, elinfo.Text.Font, elinfo.Text.Size, elinfo.Text.Color)
		//		text.FontStyle = elinfo.Text.Style
		//		text.Multiline = elinfo.Text.Multiline
		//		text.TextAlignment = elinfo.Text.Textalignment
		text.SetFontInfo(&elinfo.Text)
		widget = text

	case ElementTypeEdit:
		edit := NewEdit(ts, elinfo.Feature.Bounds.X, elinfo.Feature.Bounds.Y, elinfo.Feature.Bounds.W, elinfo.Feature.Bounds.H)
		//		log.Info("edit:", ts, elinfo.Feature.Bounds.X, elinfo.Feature.Bounds.Y, elinfo.Feature.Bounds.W, elinfo.Feature.Bounds.H)
		edit.Label.SetFontInfo(&elinfo.Text)
		edit.MaxLimit = elinfo.Text.MaxLimit
		edit.SetText(elinfo.Text.Text)
		widget = edit

	case ElementTypeButton:

		var err error
		var normal, hovering, pressdown IWidget
		var normalAX, normalAY, hoveringAX, hoveringAY, pressdownAX, pressdownAY REAL
		if elinfo.NormalRefid != "" {
			normal, err = NewElementByResourceId(ts, elinfo.NormalRefid, reselsmap...)
			if err != nil {
				log.Error(err)
			} else {
				normalAX, normalAY = normal.GetAnchorPoint()
			}
		}
		if elinfo.HoveringRefid != "" {
			hovering, err = NewElementByResourceId(ts, elinfo.HoveringRefid, reselsmap...)
			if err != nil {
				log.Error(err)
			} else {
				hoveringAX, hoveringAY = hovering.GetAnchorPoint()
			}
		}
		if elinfo.PressdownRefid != "" {
			pressdown, err = NewElementByResourceId(ts, elinfo.PressdownRefid, reselsmap...)
			if err != nil {
				log.Error(err)
			} else {
				pressdownAX, pressdownAY = pressdown.GetAnchorPoint()
			}
		}
		//log.Info(elinfo.Id, normal, hovering, pressdown)
		button := NewButton(elinfo.Feature.Bounds.X, elinfo.Feature.Bounds.Y, elinfo.Feature.Bounds.W, elinfo.Feature.Bounds.H,
			elinfo.Text.Text, elinfo.Text.Size, normal, hovering, pressdown)
		button.Text.SetFontInfo(&elinfo.Text)
		if normal != nil {
			normal.SetAnchorPoint(normalAX, normalAY)
		} else if elinfo.NormalColor != ZRGBA {
			button.Normal.SetFillColor(elinfo.NormalColor)
		}

		if hovering != nil {
			hovering.SetAnchorPoint(hoveringAX, hoveringAY)
		} else if elinfo.HoveringColor != ZRGBA {
			button.Hovering.SetFillColor(elinfo.HoveringColor)
		}

		if pressdown != nil {
			pressdown.SetAnchorPoint(pressdownAX, pressdownAY)
		} else if elinfo.PressdownColor != ZRGBA {
			button.Pressdown.SetFillColor(elinfo.PressdownColor)
		}

		widget = button

	case ElementTypeElement:
		var err error
		widget, err = NewElementByResourceId(ts, elinfo.RefId, reselsmap...)
		if err != nil {
			log.Error(err)
		}
	}

	if widget != nil {

		if elinfo.Id != "" {
			widget.SetId(elinfo.Id)
		}
		widget.SetFeature(&elinfo.Feature)

		widget.SetType(elinfo.Type)
		widget.SetTag(elinfo.Tag)
		//widget.SetDesc(elinfo.Desc)
		widget.SetEventEnabled(elinfo.EventEnabled)
		widget.SetVisible(!elinfo.Invisible)
		widget.SetObstacle(elinfo.Obstacle)

		for _, elinfo := range elinfo.Children {
			el := NewElementByResource(ts, elinfo)
			if el != nil {
				var err error
				err = widget.AddChild(el)
				if err != nil {
					log.Error(err)
				}
			}
		}

		return widget
	}
	return nil
}
