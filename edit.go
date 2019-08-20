package game2e

import (
	"bytes"
	"image"
	"image/color"
	"math"

	. "github.com/tryor/game2e/util"

	. "github.com/tryor/eui"
	. "github.com/tryor/winapi"
)

const key_LF = 10
const key_Return = 13
const key_Back = 8
const key_Delete = 46

type Edit struct {
	*Widget
	Label    *Label
	MaxLimit int
	Text     []rune

	editCursor        *Cursor //选择模式下移动元素时光标
	showStartPosition int     //可显示字符开始位置
	cursorPosition    int     //光标字符位置
	cx                int
	cy                int
	ch                int

	lastFontFillColor color.Color
	isSelectAll       bool

	selectAll  *Widget
	textCursor *Widget
	timer      *Timer

	MoveEditCursorEnable bool

	onCharPressEvent func(sender interface{}, key *uint16)
	onChangeEvent    func(sender interface{})
}

func (p *Edit) SetOnCharPress(fn func(sender interface{}, key *uint16)) {
	if p.onCharPressEvent == nil {
		p.onCharPressEvent = fn
	} else {
		onCharPressEvent := p.onCharPressEvent
		p.onCharPressEvent = func(sender interface{}, key *uint16) {
			fn(sender, key)
			if *key != 0 {
				onCharPressEvent(sender, key)
			}
		}
	}
}

func (p *Edit) SetOnChange(fn func(sender interface{})) {
	if p.onChangeEvent == nil {
		p.onChangeEvent = fn
	} else {
		onChangeEvent := p.onChangeEvent
		p.onChangeEvent = func(sender interface{}) {
			onChangeEvent(sender)
			fn(sender)
		}
	}

}

func NewEdit(ts *TaskSchedule, x, y, w, h int) *Edit {
	edit := &Edit{Widget: NewWidget(), MaxLimit: math.MaxInt32}
	//edit.Label = NewLabel2(x, y, w, h, "", "Courier New", 13, color.NRGBA{0, 0, 0, 255})
	//edit.Label = NewLabel2(x, y, w, h, "", "Courier New", 10, color.NRGBA{0, 0, 0, 255})
	edit.Label = NewLabel2(0, 0, w, h, "", "Arial", 13, color.NRGBA{0, 0, 0, 255})
	edit.Self = edit

	edit.SetCoordinate(x, y)
	edit.SetWidth(w)
	edit.SetHeight(h)
	edit.SetBorder(true)
	edit.BorderColor = color.NRGBA{0, 0, 0, 100}
	edit.Fill = true
	edit.FillColor = color.NRGBA{255, 255, 255, 0x00}

	//edit.SetAutoDrawChildren(false)
	edit.Label.SetAnchorPoint(0, 0)
	edit.Label.Multiline = false
	edit.MoveEditCursorEnable = true
	edit.Label.SetEventEnabled(false)
	//	edit.Multiline = gdiplus.StringFormatFlagsMeasureTrailingSpaces |
	//		gdiplus.StringFormatFlagsNoWrap
	//		gdiplus.StringFormatFlagsNoClip
	//gdiplus.StringFormatFlagsDisplayFormatControl

	edit.selectAll = NewWidget()
	edit.selectAll.SetAnchorPoint(0, 0)
	edit.selectAll.SetEventEnabled(false)
	edit.selectAll.SetVisible(false)
	edit.selectAll.SetFill(true)
	edit.selectAll.SetFillColor(color.RGBA{0, 112, 192, 0xff})
	edit.AddChild(edit.selectAll)

	edit.AddChild(edit.Label)

	edit.editCursor = NewCursor(InitCursor(IDC_IBEAM))
	edit.Text = make([]rune, 0)

	edit.cx = 2
	edit.cy = 0
	//	edit.ch = int(edit.Label.FontSize)
	edit.textCursor = NewWidget()
	edit.textCursor.SetWidth(1)
	edit.textCursor.SetHeight(edit.ch)
	//	edit.textCursor.SetCoordinate(x+edit.cx, y+edit.cy)
	edit.textCursor.Fill = true
	edit.textCursor.FillColor = color.NRGBA{0, 0, 0, 255}
	edit.textCursor.Border = true
	edit.textCursor.SetAnchorPoint(0, 0)
	edit.textCursor.SetEventEnabled(false)

	//	edit.AddSubWidget(edit.textCursor)

	//log.Println("edit.id:", edit.GetId())

	//光标闪烁
	edit.timer = NewTimer(ts, func(delayIndex int) {
		edit.switchTextCursor()
	}, 50)
	//	edit.timer.Stop()

	edit.OnMouseDoubleClick(func(me *MouseEvent) {
		edit.SelectAll()
	})

	edit.OnMouseDown(func(me *MouseEvent) {
		edit.cancelSelectAll()
		edit.SetFocus()
	})

	edit.OnMouseEnter(func(me *MouseEvent) {
		edit.onMouseEnterEvent(me)
	})

	edit.OnMouseExit(func(me *MouseEvent) {
		edit.onMouseLeaveEvent(me)
	})

	edit.OnFocus(func(fe *FocusEvent) {
		edit.onFocusEvent(fe)
	})

	edit.OnKeyChar(func(ke *KeyEvent) {
		key := uint16(ke.Char)
		if edit.onCharPressEvent != nil {
			edit.onCharPressEvent(edit, &key)
			if key == 0 {
				return
			}
		}
		ke.Char = rune(key)
		edit.onKeyCharEvent(ke)
	})

	edit.OnKeyPress(func(ke *KeyEvent) {
		edit.onKeyPressEvent(ke)
	})

	edit.OnKeyRepeat(func(ke *KeyEvent) {
		edit.onKeyPressEvent(ke)
	})

	return edit
}

func (this *Edit) Destroy() {
	this.ClearChildren()
	this.Label.Destroy()
	if this.timer != nil {
		this.timer.Stop()
		this.timer = nil
	}
	this.textCursor.Destroy()
	this.selectAll.Destroy()
	this.Destroy()
}

func (this *Edit) fireOnChangeEvent() {
	if this.onChangeEvent != nil {
		this.onChangeEvent(this)
	}
}

//func (this *Edit) SetEditCursorEnable(b bool) {
//	if b {
//		this.timer.Start()
//	} else {
//		this.timer.Stop()
//		this.hideTextCursor()
//	}
//}

func (this *Edit) SetWidth(w int) {
	this.Widget.SetWidth(w)
	this.Label.SetWidth(w)
}
func (this *Edit) SetHeight(h int) {
	this.Widget.SetHeight(h)
	this.Label.SetHeight(h)
}

func (this *Edit) SelLength() int {
	if this.isSelectAll {
		return len(this.Text)
	}
	return 0
}

func (this *Edit) SelectAll() {
	this.isSelectAll = true
	this._selectAll()
}

func (this *Edit) _selectAll() {
	if this.isSelectAll && !this.selectAll.IsVisible() {
		boundingBox, _, _ := this.Label.MeasureString(this.Label.Caption)
		if boundingBox != nil {
			this.selectAll.SetCoordinate(boundingBox.Min.X, boundingBox.Min.Y)
			this.selectAll.SetWidth(boundingBox.Dx())
			this.selectAll.SetHeight(boundingBox.Dy())
			this.selectAll.SetVisible(true)

			this.lastFontFillColor = this.Label.FontFillColor

			//r, g, b, a := this.FontFillColor.RGBA()
			//this.FontFillColor = color.RGBA{uint8(^r), uint8(^g), uint8(^b), uint8(a)}
			this.Label.FontFillColor = color.RGBA{0xff, 0xff, 0xff, 0xff}
			this.Label.ClearBrush()
			this.Label.SetModified(true)

		} else {
			//log.Println("selectAll:", boundingBox, this.Caption)
		}
	}
}

func (this *Edit) clearSelectAll() {
	if this.isSelectAll {
		this.cancelSelectAll()
		this.setText("")
	}
}

func (this *Edit) cancelSelectAll() {
	if this.isSelectAll {
		this.isSelectAll = false
		this.selectAll.SetVisible(false)
		this.Label.FontFillColor = this.lastFontFillColor
		this.Label.ClearBrush()
		this.Label.SetModified(true)
	}
}

//func (this *Edit) onMouseDoubleClick() {

//}

func (this *Edit) SetText(text string) {
	this.clearSelectAll()
	this.setText(text)
}

func (this *Edit) setText(text string) {
	if string(this.Text) != text {
		this.Text = bytes.Runes([]byte(text))
		this.cursorPosition = len(this.Text)
		this.showStartPosition = 0
		this.calculateCursorCoordinate()
		this.Label.Caption = text
		this.Self.(IElement).SetModified(true)
		this.fireOnChangeEvent()
	}
}

func (this *Edit) SetTextAlignment(align Alignment) {
	if this.Label.TextAlignment != align {
		this.Label.TextAlignment = align
		this.Label.ClearStringFormat()
		this.Label.SetModified(true)
	}

	if !this.Label.Multiline && this.Label.TextAlignment.Vertical == AlignmentCenter {
		align := this.textCursor.GetAlignment()
		align.Vertical = this.Label.TextAlignment.Vertical
		this.textCursor.SetAlignment(align)
		this.textCursor.SetAnchorPoint(0, 0.5)
	}

}

func (this *Edit) showTextCursor() {
	this.AddChild(this.textCursor)
	this.calculateCursorCoordinate()

	this.Self.(IElement).SetModified(true)
}

func (this *Edit) hideTextCursor() {
	this.RemoveChild(this.textCursor)
	this.Self.(IElement).SetModified(true)
}

func (this *Edit) switchTextCursor() {
	if this.GetSelf().HasChild(this.textCursor) {
		this.hideTextCursor()
	} else {
		this.showTextCursor()
	}
}

func (this *Edit) flushTextCursor() {
	if this.ch == 0 {
		this.ch = int(this.Label.FontSize)
	}
	thisEL := this.Self.(IElement)

	ax, ay := this.Self.(IElement).GetAnchorPoint()
	aw := int(REAL(thisEL.Width()) * ax)
	ah := int(REAL(thisEL.Height()) * ay)

	x, y := this.cx-aw, this.cy-ah

	if this.Label.TextAlignment.Horizontal == AlignmentFar {
		if !this.MoveEditCursorEnable {
			x = 0 //@TODO ... 暂时不处理右对齐光标问题
		}
		x = this.Width() - x
	}

	//	log.Println("cursorPosition, showStartPosition:", this.cursorPosition, this.showStartPosition)
	//	log.Println("x, y:", x, y)
	//	log.Println("this.cx-aw, this.cy-ah:", this.cx-aw, this.cy-ah)
	//	log.Println("this.cx, this.cy:", this.cx, this.cy)
	//	log.Println("aw, ah:", aw, ah)

	this.textCursor.SetCoordinate(x, y)
	this.textCursor.SetHeight(this.ch)
}

//func (this *Edit) CreatePath() {
//	this.Label.CreatePath()
//	this._selectAll()
//}

//func (this *Edit) Draw(ge IGraphicsEngine) {
//	this.Widget.Draw(ge)
//	//this.DrawChildren(ge)
//	//this.Label.Draw(ge)
//}

func (this *Edit) MoveBy(dx, dy int, angle float32) {
	this.Label.MoveBy(dx, dy, angle)
	this.hideTextCursor()
}

func (this *Edit) onMouseEnterEvent(e *MouseEvent) {
	//	println("Edit.OnMouseEnterEvent.Focus:", this.IsFocus())
	//if this.IsFocus() {
	this.SetCursor(this.editCursor)
	//}
}

func (this *Edit) onMouseLeaveEvent(e *MouseEvent) {
	//	println("Edit.OnMouseLeaveEvent.Focus:", this.IsFocus())
	this.ResetCursor()
	//	this.Label.OnMouseLeaveEvent(e)
}

func (this *Edit) onFocusEvent(e *FocusEvent) {
	//log.Println("Edit.onFocusEvent.Focus:", e.Focus)
	if e.Focus {
		this.SetCursor(this.editCursor)
		this.timer.Start()
	} else {
		this.timer.Stop()
		this.hideTextCursor()
	}
	//	this.Label.OnFocusEvent(e)
}

//在指定位置插入字符
func (this *Edit) insertChar(pos int, c rune) {
	text := make([]rune, 0, len(this.Text)+1)

	//	if this.TextAlignment.Horizontal == AlignmentNear {
	text = append(text, this.Text[0:pos]...)
	text = append(text, c)
	text = append(text, this.Text[pos:]...)
	this.Text = text
	this.cursorPosition++
	this.fireOnChangeEvent()
	//if this.MoveEditCursorEnable {

	//}

	//	} else { //if this.TextAlignment.Horizontal == AlignmentFar {
	//		text = append(text, this.Text[pos:]...)
	//		text = append(text, c)
	//		text = append(text, this.Text[0:pos]...)
	//		this.Text = text
	//	}

}

func (this *Edit) removeChar(pos int) {
	text := make([]rune, 0, len(this.Text)-1)
	text = append(text, this.Text[0:pos-1]...)
	text = append(text, this.Text[pos:]...)
	this.Text = text
	this.cursorPosition--
	if this.showStartPosition > 1 {
		this.showStartPosition -= 2
	} else if this.showStartPosition > 0 {
		this.showStartPosition--
	}
	this.fireOnChangeEvent()
}

func (this *Edit) onKeyCharEvent(e *KeyEvent) {

	if e.Char < 32 || (e.Char > 126 && e.Char < 160) { //不可见字符
		return
	}
	this.clearSelectAll()
	if !noShowKeys[e.Char] && len(this.Text) < this.MaxLimit {
		c := e.Char
		this.insertChar(this.cursorPosition, c)
		this.calculateCursorCoordinate(true)
		this.Label.Caption = string(this.Text[this.showStartPosition:])
		this.Self.(IElement).SetModified(true)
	}
}

func (this *Edit) onKeyPressEvent(e *KeyEvent) {
	switch e.Key {
	case KeyEnter:
		if this.Label.Multiline {
			this.clearSelectAll()
			this.insertChar(this.cursorPosition, key_LF)
			this.calculateCursorCoordinate(true)
			this.Label.Caption = string(this.Text[this.showStartPosition:])
			this.Self.(IElement).SetModified(true)
		}

	case KeyBackspace: //Key_Back:
		this.clearSelectAll()
		if this.cursorPosition > 0 {
			this.removeChar(this.cursorPosition)
			if this.cursorPosition > 0 && this.Text[this.cursorPosition-1] == key_LF {
				this.removeChar(this.cursorPosition)
			}
			this.Label.Caption = string(this.Text[this.showStartPosition:])
			this.Self.(IElement).SetModified(true)
			//			this.Self.(IElement).CreatePath()
			this.calculateCursorCoordinate()
		}
	case KeyLeft, KeyRight, KeyHome, KeyEnd: //Key_Left, Key_Right, Key_Home, Key_End:
		this.cancelSelectAll()
		if !this.MoveEditCursorEnable {
			return
		}

		switch e.Key {
		case KeyLeft:
			//			if this.TextAlignment.Horizontal == AlignmentNear {
			if this.cursorPosition > 0 {
				this.cursorPosition--
				if this.cursorPosition < this.showStartPosition && this.showStartPosition > 0 {
					this.showStartPosition--
				}
			}
			//			} else { //if this.TextAlignment.Horizontal == AlignmentFar {
			//				if this.cursorPosition < len(this.Text) {
			//					this.cursorPosition++
			//				}
			//			}
		case KeyRight:
			//			if this.TextAlignment.Horizontal == AlignmentNear {
			if this.cursorPosition < len(this.Text) {
				this.cursorPosition++
			}
			//			} else { //if this.TextAlignment.Horizontal == AlignmentFar {
			//				if this.cursorPosition > 0 {
			//					this.cursorPosition--
			//					if this.cursorPosition < this.showStartPosition && this.showStartPosition > 0 {
			//						this.showStartPosition--
			//					}
			//				}
			//			}
		case KeyHome:
			this.showStartPosition = 0
			this.cursorPosition = 0
		case KeyEnd:
			//			this.cursorPosition = len(this.Text)
			//			this.showStartPosition
			for this.cursorPosition < len(this.Text) {
				this.cursorPosition++
				this.calculateCursorCoordinate()
			}
		}
		this.calculateCursorCoordinate()
		this.Label.Caption = string(this.Text[this.showStartPosition:])
		this.Self.(IElement).SetModified(true)
		//		this.Self.(IElement).CreatePath()
	}

	//	this.Label.OnKeyPressEvent(e)
}

//返回当前行首到当前位置的字符串
func (this *Edit) getCurrentLineText() (line int, str []rune, size int) {
	lineFirstPos := this.showStartPosition

	for i := this.showStartPosition; i < this.cursorPosition; i++ {
		if this.Text[i] == key_LF {
			lineFirstPos = i
			line++
		}
	}
	//	println("this.showStartPosition:", this.showStartPosition, lineFirstPos, this.cursorPosition, len(this.Text))
	if lineFirstPos < this.cursorPosition {
		size = this.cursorPosition - lineFirstPos
		str = this.Text[lineFirstPos:this.cursorPosition]
	}
	return
}

func (this *Edit) calculateCursorCoordinate(insert ...bool) {

	if this.cursorPosition > 0 {
		cline, cltext, clsize := this.getCurrentLineText()
		if clsize > 0 {

			boundingBox, codepointsFitted, linesFilled := this.Label.MeasureString(string(cltext))
			//println("Edit.getCurrentLineText:", string(cltext), cline, clsize, len(this.Text), linesFilled)
			if boundingBox != nil && *boundingBox != image.ZR {
				h := boundingBox.Dy()
				if linesFilled == 1 && h > 0 && this.ch != h {
					this.ch = h
				}

				if !this.Label.Multiline && codepointsFitted < clsize && this.showStartPosition < len(this.Text) {
					this.showStartPosition++
				}

				endpost := (codepointsFitted < clsize || linesFilled > cline+1) && len(insert) > 0 && insert[0]
				//log.Info("", codepointsFitted, clsize, linesFilled, cline, len(insert), insert)
				//log.Info("this.Multiline, endpost:", this.Multiline, endpost)
				if this.Label.Multiline && endpost {
					this.insertChar(this.cursorPosition-1, key_LF)
					this.calculateCursorCoordinate()
				} else {

					cy := cline * this.ch
					//					println("cy:", this.cy, cy+this.ch, this.Label.H)
					//					if cy+this.ch > this.Label.H {
					//						this.ch = this.Label.H - cy
					//					}
					if cy <= this.Label.Height() {
						this.cy = cy
						if clsize == 1 && cltext[0] == key_LF {
							this.cx = 1
						} else {
							this.cx = boundingBox.Dx()
						}
					}
				}
			}
		} else {
			this.cx = 1
		}
	} else {
		this.cx = 1
		this.cy = 0
	}
	this.flushTextCursor()
}

var noShowKeys = map[rune]bool{
	//	key_Back:   true,
	//	key_Delete: true
}
