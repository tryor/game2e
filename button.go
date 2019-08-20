package game2e

import (
	"image/color"

	. "github.com/tryor/eui"
)

type Button struct {
	*Widget
	Text           *Label
	current        IWidget
	Normal         IWidget //正常状态
	Hovering       IWidget //鼠标悬停
	hoveringExist  bool
	HoveringCover  bool
	Pressdown      IWidget //鼠标按下
	pressdownExist bool
}

func NewButton(x, y, w, h int, text string, fontSize float32, args ...IWidget) *Button {
	button := &Button{Widget: NewWidget()}
	button.Self = button
	button.SetCoordinate(x, y)
	button.Widget.SetWidth(w)
	button.Widget.SetHeight(h)
	button.Text = NewLabel(0, 0, w, h, text, "黑体", fontSize) //Verdana
	button.Text.TextAlignment.Horizontal = AlignmentCenter
	button.Text.TextAlignment.Vertical = AlignmentCenter
	//button.Text.SetEventEnabled(false)

	if len(args) > 0 && args[0] != nil {
		button.Normal = args[0]
	}

	if len(args) > 1 && args[1] != nil {
		button.Hovering = args[1]
		button.hoveringExist = true
	}

	if len(args) > 2 && args[2] != nil {
		button.Pressdown = args[2]
		button.pressdownExist = true
	}

	if button.Normal == nil {
		normal := NewWidget()
		normal.SetWidth(w)
		normal.SetHeight(h)
		normal.SetCoordinate(0, 0)
		normal.Fill = true
		normal.FillColor = color.NRGBA{0, 033, 99, 200}
		button.Normal = normal
		//button.Normal.SetAnchorPoint(0, 0)
	}
	button.current = button.Normal

	if button.Hovering == nil {
		hovering := NewWidget()
		hovering.SetWidth(w)
		hovering.SetHeight(h)
		hovering.SetCoordinate(0, 0)
		hovering.Fill = true
		hovering.FillColor = color.NRGBA{255, 255, 255, 50}
		button.Hovering = hovering
		//button.Hovering.SetAnchorPoint(0, 0)
	}

	if button.Pressdown == nil {
		pressdown := NewWidget()
		pressdown.SetWidth(w)
		pressdown.SetHeight(h)
		pressdown.SetCoordinate(0, 0)
		pressdown.Fill = true
		pressdown.FillColor = color.NRGBA{200, 255, 100, 50}
		button.Pressdown = pressdown
		//button.Pressdown.SetAnchorPoint(0, 0)
	}

	button.Normal.SetAnchorPoint(0, 0)
	button.Hovering.SetAnchorPoint(0, 0)
	button.Pressdown.SetAnchorPoint(0, 0)
	button.Text.SetAnchorPoint(0, 0)

	button.Text.SetEventEnabled(false)
	button.Normal.SetEventEnabled(false)
	button.Hovering.SetEventEnabled(false)
	button.Pressdown.SetEventEnabled(false)

	button.AddChild(button.current, 0)
	button.AddChild(button.Text, 1)

	button.OnMouseEnter(func(me *MouseEvent) {
		button.onMouseEnterEvent(me)
	})
	button.OnMouseExit(func(me *MouseEvent) {
		button.onMouseLeaveEvent(me)
	})
	button.OnMouseDown(func(me *MouseEvent) {
		button.onMouseDownEvent(me)
	})
	button.OnMouseUp(func(me *MouseEvent) {
		button.onMouseUpEvent(me)
	})

	return button
}

func (this *Button) Destroy() {
	//	this.ClearChildren()
	this.RemoveChild(this.Text)
	this.RemoveChild(this.Normal)
	this.RemoveChild(this.Hovering)
	this.RemoveChild(this.Pressdown)

	this.Text.Destroy()
	this.Normal.Destroy()
	this.Hovering.Destroy()
	this.Pressdown.Destroy()

	this.Widget.Destroy()
}

func (this *Button) SetLayer(l ILayer) {
	this.Widget.SetLayer(l)
	this.Text.SetLayer(l)
	this.Normal.SetLayer(l)
	this.Hovering.SetLayer(l)
}

func (this *Button) onMouseEnterEvent(e *MouseEvent) {
	if this.hoveringExist && !this.HoveringCover {
		this.current.SetVisible(false)
		//this.RemoveChild(this.current)
		this.current = this.Hovering
		//this.AddChild(this.current, 0)
		if !this.ExistChild(this.current.GetId()) {
			this.AddChild(this.current, 0)
		}
	} else {
		this.current = this.Hovering
		//this.AddChild(this.current, 1)
		if !this.ExistChild(this.current.GetId()) {
			this.AddChild(this.current, 1)
		} else {
			this.current.SetOrderZ(1)
		}
	}
	this.current.SetVisible(true)
	this.Self.(IElement).SetModified(true)
}

func (this *Button) onMouseLeaveEvent(e *MouseEvent) {
	if this.hoveringExist && !this.HoveringCover {
		this.current.SetVisible(false)
		//this.RemoveChild(this.current)
		this.current = this.Normal
		//this.AddChild(this.current, 0)
		if !this.ExistChild(this.current.GetId()) {
			this.AddChild(this.current, 0)
		}
	} else {
		this.current.SetVisible(false)
		//this.RemoveChild(this.current)
		this.current = this.Normal
	}
	this.current.SetVisible(true)
	this.Self.(IElement).SetModified(true)
}

func (this *Button) onMouseDownEvent(e *MouseEvent) {
	if this.pressdownExist {
		this.current.SetVisible(false)
		//this.RemoveChild(this.current)
		this.current = this.Pressdown
		//this.AddChild(this.current, 0)
		if !this.ExistChild(this.current.GetId()) {
			this.AddChild(this.current, 0)
		}
	} else {
		this.current = this.Pressdown
		//this.AddChild(this.current, 1)
		if !this.ExistChild(this.current.GetId()) {
			this.AddChild(this.current, 1)
		} else {
			this.current.SetOrderZ(1)
		}
	}
	this.current.SetVisible(true)
	this.Self.SetModified(true)
}

func (this *Button) onMouseUpEvent(e *MouseEvent) {
	//this.RemoveChild(this.current)
	this.current.SetVisible(false)
	this.onMouseEnterEvent(e)
}

func (this *Button) SetWidth(width int) {
	this.Widget.SetWidth(width)
	this.Text.SetWidth(width)
	this.Normal.SetWidth(width)
	this.Hovering.SetWidth(width)
	this.Pressdown.SetWidth(width)
}

func (this *Button) SetHeight(height int) {
	this.Widget.SetHeight(height)
	this.Text.SetHeight(height)
	this.Normal.SetHeight(height)
	this.Hovering.SetHeight(height)
	this.Pressdown.SetHeight(height)
}
