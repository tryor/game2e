package game2e

import (
	//	"github.com/google/gxui"
	. "github.com/tryor/eui"
	. "github.com/tryor/winapi"
)

type IViewport interface {
	GetHandle() HANDLE

	OnKeyChar(f func(*KeyEvent)) EventSubscription
	OnKey(f func(*KeyEvent)) EventSubscription
	OnMouse(f func(IMouseEvent)) EventSubscription
	OnResize(f func()) EventSubscription
	OnClose(f func()) EventSubscription

	Size() Size
	//	Loop(func(viewport IViewport))
	Close()
}
