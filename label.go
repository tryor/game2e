package game2e

import (
	"image"
	"image/color"

	. "github.com/tryor/eui"
	"github.com/tryor/gdiplus"
)

type Label struct {
	*Widget
	TextAlignment Alignment //文本对齐方式
	//	StringFormatFlags gdiplus.StringFormatFlags
	Multiline bool

	Caption       string
	FontFillColor color.Color
	FontFamily    string
	FontStyle     gdiplus.FontStyle
	FontSize      float32
	FontUnit      IUnit
	TextHeight    int  //显示内容实际高度
	UsePath       bool //使用路径模式显示

	fontfamily   IFontFamily
	stringFormat IStringFormat
	font         IFont
	brush        IBrush
}

func NewLabel(x, y, w, h int, caption string, fontFamily string, fontSize float32) *Label {
	return NewLabel2(x, y, w, h, caption, fontFamily, fontSize, color.NRGBA{255, 255, 255, 255})
}

func NewLabel2(x, y, w, h int, caption string, fontFamily string, fontSize float32, fillColor color.Color) *Label {
	label := &Label{Widget: NewWidget()}
	label.Self = label
	label.SetWidth(w)
	label.SetHeight(h)
	label.SetCoordinate(x, y)

	label.Caption = caption
	if fillColor != nil {
		label.FontFillColor = fillColor
	} else {
		label.FontFillColor = label.FillColor
	}
	label.FontFamily = fontFamily
	label.FontSize = fontSize
	label.FontUnit = IUnit(gdiplus.UnitPixel)
	//	label.FontStyle = gdiplus.FontStyleBold

	return label
}

func (this *Label) Destroy() {
	this.Widget.Destroy()
	if this.font != nil {
		this.font.Release()
		this.font = nil
	}
	if this.fontfamily != nil {
		this.fontfamily.Release()
		this.fontfamily = nil
	}
	if this.stringFormat != nil {
		this.stringFormat.Release()
		this.stringFormat = nil
	}
	if this.brush != nil {
		this.brush.Release()
		this.brush = nil
	}
}

func (this *Label) GetFontFillColorA() uint8 {
	_, _, _, a := this.FontFillColor.RGBA()
	return uint8(a)
}

func (this *Label) SetFontFillColorA(a uint8) {
	r, g, b, _ := GetColorRGBA(this.FontFillColor)
	this.FontFillColor = color.NRGBA{r, g, b, a}
	this.ClearBrush()
	this.Self.SetModified(true)
}

func (this *Label) SetFontInfo(fontInfo *ResFontInfo) {
	this.Caption = fontInfo.Text
	this.FontFamily = fontInfo.Font
	this.FontSize = fontInfo.Size
	this.FontFillColor = fontInfo.Color
	this.FontStyle = fontInfo.Style
	this.Multiline = fontInfo.Multiline
	this.TextAlignment = fontInfo.Textalignment
	this.ClearBrush()
	this.ClearFontfamily()
	this.ClearStringFormat()
	this.ClearFont()
	this.Self.(IElement).SetModified(true)
}

func (this *Label) SetCaption(text string) {
	this.Caption = text
	this.SetModified(true)
}

func (this *Label) ClearBrush() {
	if this.brush != nil {
		this.brush.Release()
		this.brush = nil
	}
}

func (this *Label) ClearFontfamily() {
	if this.fontfamily != nil {
		this.fontfamily.Release()
		this.fontfamily = nil
	}
}

func (this *Label) ClearStringFormat() {
	if this.stringFormat != nil {
		this.stringFormat.Release()
		this.stringFormat = nil
	}
}

func (this *Label) ClearFont() {
	if this.font != nil {
		this.font.Release()
		this.font = nil
	}
}

func (this *Label) CreatePath() {
	this.Widget.CreatePath()
	ge := this.GetGraphicsEngine()
	if ge != nil && this.Caption != "" {
		this.Self.CreateBoundRect()
		var err error
		if this.fontfamily == nil {
			this.fontfamily, err = ge.NewFontFamily(this.FontFamily)
			if err != nil {
				this.LastError = err
			}
		}

		if this.stringFormat == nil {
			this.stringFormat, err = ge.NewStringFormat()
			if err != nil {
				this.LastError = err
			} else {
				this.stringFormat.SetHorizontalAlignment(this.TextAlignment.Horizontal)
				this.stringFormat.SetVerticalAlignment(this.TextAlignment.Vertical)
				this.stringFormat.SetMultiline(this.Multiline)
			}
		}

		if this.font == nil {
			this.font, err = ge.NewFont(this.fontfamily, this.FontSize,
				IFontStyle(this.FontStyle), this.FontUnit) //gdiplus.UnitPoint, gdiplus.UnitPixel
			if err != nil {
				this.LastError = err
			}
		}

		if this.brush == nil {
			this.brush, err = ge.NewBrush(this.FontFillColor)
			if err != nil {
				this.LastError = err
			}
		}

		if this.UsePath {
			this.InitPath(this.Layer.GetGraphicsEngine())
			rect := FormatRect(this.GetBoundRect())
			minwidth := int(this.BorderWidth*2) + 1
			if this.Caption != "" && rect.Dx() > minwidth && rect.Dy() > minwidth {
				this.Path.AddString(rect, this.Caption, this.FontSize, this.fontfamily, IFontStyle(this.FontStyle), this.stringFormat)
			}
		}

	}
}

func (this *Label) Draw(ge IGraphicsEngine) {
	this.Widget.Draw(ge)
	if this.Caption == "" {
		return
	}
	if this.UsePath {
		ge.SetRenderBrush(this.brush)
		this.RenderPath(ge, OnlyFill)
	} else {
		brect := this.Self.(IElement).GetBoundRect()
		minwidth := int(this.BorderWidth*2) + 1
		if this.Caption != "" && brect.Dx() > minwidth && brect.Dy() > minwidth {
			brect_ := *brect
			brect_.Min.Y += 4
			ge.DrawString(this.Caption, this.font, &brect_, this.stringFormat, this.brush)
		}
	}

}

func (this *Label) MeasureString(text ...string) (boundingBox *image.Rectangle, codepointsFitted, linesFilled int) {
	ge := this.GetGraphicsEngine()
	if ge != nil {
		var err error
		if this.font == nil {
			this.fontfamily, err = ge.NewFontFamily(this.FontFamily)
			if err == nil {
				this.font, err = ge.NewFont(this.fontfamily, this.FontSize,
					IFontStyle(this.FontStyle), this.FontUnit)
			}
			if err != nil {
				this.LastError = err
			}
		}

		if this.stringFormat == nil {
			this.stringFormat, err = ge.NewStringFormat()
			if err != nil {
				this.LastError = err
			} else {
				this.stringFormat.SetHorizontalAlignment(this.TextAlignment.Horizontal)
				this.stringFormat.SetVerticalAlignment(this.TextAlignment.Vertical)
				this.stringFormat.SetMultiline(this.Multiline)
			}
		}

		var text_ string
		if len(text) > 0 {
			text_ = text[0]
		} else {
			text_ = this.Caption
		}

		layoutRect := image.Rect(0, 0, this.Width(), this.Height())
		boundingBox, codepointsFitted, linesFilled = ge.MeasureString(text_, this.font, &layoutRect, this.stringFormat)
		this.TextHeight = boundingBox.Dy()
	}
	return
}
