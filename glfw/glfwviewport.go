package glfw

import (
	"unicode"
	"unsafe"

	. "tryor/game2e"
	"tryor/game2e/log"

	"github.com/go-gl/glfw/v3.2/glfw"
	//	"github.com/google/gxui"
	"github.com/tryor/commons/event"
	. "github.com/tryor/eui"
	. "github.com/tryor/winapi"
)

func Run(appRoutine func(driver IDriver)) {
	var driverFuncs DriverFuncs
	driverFuncs.Init = glfw.Init
	driverFuncs.Terminate = glfw.Terminate
	driverFuncs.WaitEvents = glfw.WaitEvents
	driverFuncs.Wake = glfw.PostEmptyEvent
	driverFuncs.CreateViewport = CreateGLFWViewport
	StartDriver(&driverFuncs, appRoutine)
}

type NativeType int

const (
	NativeTypeWin32 NativeType = 0
	NativeTypeWGL   NativeType = 1
)

type GLFWViewport struct {
	glfwWindow *glfw.Window
	nativeType NativeType

	onResizeEvent  Event
	onCloseEvent   Event
	onMouseEvent   Event
	onKeyEvent     Event
	onKeyCharEvent Event
}

func CreateGLFWViewport(width, height int, title string) (IViewport, error) {

	glfw.WindowHint(glfw.Resizable, 0)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	window.MakeContextCurrent()

	viewport := &GLFWViewport{glfwWindow: window, nativeType: NativeTypeWin32}
	viewport.initEvents()
	return viewport, nil
}

func (this *GLFWViewport) SetNativeType(wintype NativeType) {
	this.nativeType = wintype
}

func (this *GLFWViewport) GetGlfwWindow() *glfw.Window {
	return this.glfwWindow
}

func (this *GLFWViewport) initEvents() {

	this.glfwWindow.SetSizeCallback(func(w *glfw.Window, width int, height int) {
		if this.onResizeEvent != nil {
			this.onResizeEvent.Fire()
		}
	})

	this.glfwWindow.SetCloseCallback(func(w *glfw.Window) {
		if this.onCloseEvent != nil {
			this.onCloseEvent.Fire()
		}

		this.glfwWindow.Destroy()
	})

	this.glfwWindow.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if this.onMouseEvent != nil {
			modifier := KeyboardModifier(mod)
			xpos, ypos := w.GetCursorPos()
			var typ event.Type

			switch action {
			case glfw.Press:
				typ = MOUSE_PRESS_EVENT_TYPE
			case glfw.Release:
				typ = MOUSE_RELEASE_EVENT_TYPE
			case glfw.Repeat:
				//typ = MOUSE_DOUBLE_CLICK_EVENT_TYPE
			}
			me := NewMouseEvent(typ, this, int(xpos), int(ypos), this.conversionMouseButton(button), modifier)
			this.onMouseEvent.Fire(me)
		}
	})

	this.glfwWindow.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		if this.onMouseEvent != nil {
			button := this.getMouseButton()
			modifier := this.getModifier()
			me := NewMouseEvent(MOUSE_MOVE_EVENT_TYPE, this, int(xpos), int(ypos), button, modifier)
			this.onMouseEvent.Fire(me)
		}

	})

	this.glfwWindow.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if this.onKeyEvent == nil {
			return
		}

		var typ event.Type
		switch action {
		case glfw.Press:
			typ = KEY_PRESS_EVENT_TYPE
		case glfw.Release:
			typ = KEY_RELEASE_EVENT_TYPE
		case glfw.Repeat:
			typ = KEY_REPEAT_EVENT_TYPE
		}

		//ke := NewKeyEvent(typ, this, KeyboardModifier(mods), translateKeyboardKey(key))
		kv := translateKeyboardKey(key)
		kvs := append(getKeyboardKeys(w, typ), kv)
		ke := NewKeyEvent(typ, this, KeyboardModifier(mods), kvs...)
		ke.Key = kv
		log.Debug("kvs:", kvs, kv, key)
		//		keys := []glfw.Key{}
		//		for _, idx := range kvs {
		//			keys = append(keys, glfwKeys[idx-1])
		//		}
		//		log.Debug("keys:", keys)
		this.onKeyEvent.Fire(ke)

	})

	this.glfwWindow.SetCharModsCallback(func(w *glfw.Window, char rune, mods glfw.ModifierKey) {
		if this.onKeyCharEvent == nil {
			return
		}

		if !unicode.IsControl(char) &&
			!unicode.IsGraphic(char) &&
			!unicode.IsLetter(char) &&
			!unicode.IsMark(char) &&
			!unicode.IsNumber(char) &&
			!unicode.IsPunct(char) &&
			!unicode.IsSpace(char) &&
			!unicode.IsSymbol(char) {
			return
		}
		ke := NewKeyCharEvent(KEY_CHAR_EVENT_TYPE, this, char, KeyboardModifier(mods))
		this.onKeyCharEvent.Fire(ke)
	})

	//this.glfwWindow.SetCharCallback()
}

func (this *GLFWViewport) getMouseButton() MButton {
	var mb MButton
	if this.glfwWindow.GetMouseButton(glfw.MouseButtonLeft) == 1 {
		mb |= MButton_Left
	}
	if this.glfwWindow.GetMouseButton(glfw.MouseButtonRight) == 1 {
		mb |= MButton_Right
	}
	if this.glfwWindow.GetMouseButton(glfw.MouseButtonMiddle) == 1 {
		mb |= MButton_Mid
	}
	return mb
}

func (this *GLFWViewport) getModifier() KeyboardModifier {
	var mod KeyboardModifier
	if this.glfwWindow.GetKey(glfw.KeyLeftShift) == glfw.Press || this.glfwWindow.GetKey(glfw.KeyRightShift) == glfw.Press {
		mod |= ModShift
	}
	if this.glfwWindow.GetKey(glfw.KeyLeftControl) == glfw.Press || this.glfwWindow.GetKey(glfw.KeyRightControl) == glfw.Press {
		mod |= ModControl
	}
	if this.glfwWindow.GetKey(glfw.KeyLeftAlt) == glfw.Press || this.glfwWindow.GetKey(glfw.KeyRightAlt) == glfw.Press {
		mod |= ModAlt
	}
	if this.glfwWindow.GetKey(glfw.KeyLeftSuper) == glfw.Press || this.glfwWindow.GetKey(glfw.KeyRightSuper) == glfw.Press {
		mod |= ModSuper
	}
	return mod
}

func (this *GLFWViewport) conversionMouseButton(button glfw.MouseButton) MButton {
	switch button {
	case glfw.MouseButtonLeft:
		return MButton_Left
	case glfw.MouseButtonRight:
		return MButton_Right
	case glfw.MouseButtonMiddle:
		return MButton_Mid
	}
	return MButton_No
}

func (this *GLFWViewport) GetHandle() HANDLE {
	switch this.nativeType {
	case NativeTypeWin32:
		return HANDLE(unsafe.Pointer(this.glfwWindow.GetWin32Window()))
	case NativeTypeWGL:
		return HANDLE(unsafe.Pointer(this.glfwWindow.GetWGLContext()))
	}
	return HANDLE(0)
}

func (this *GLFWViewport) OnKeyChar(f func(*KeyEvent)) EventSubscription {
	if this.onKeyCharEvent == nil {
		this.onKeyCharEvent = CreateEvent(func(*KeyEvent) {})
	}
	return this.onKeyCharEvent.Listen(f)
}

func (this *GLFWViewport) OnKey(f func(*KeyEvent)) EventSubscription {
	if this.onKeyEvent == nil {
		this.onKeyEvent = CreateEvent(func(*KeyEvent) {})
	}
	return this.onKeyEvent.Listen(f)
}

func (this *GLFWViewport) OnMouse(f func(IMouseEvent)) EventSubscription {
	if this.onMouseEvent == nil {
		this.onMouseEvent = CreateEvent(func(IMouseEvent) {})
	}
	return this.onMouseEvent.Listen(f)
}

func (this *GLFWViewport) OnResize(f func()) EventSubscription {
	if this.onResizeEvent == nil {
		this.onResizeEvent = CreateEvent(func() {})
	}
	return this.onResizeEvent.Listen(f)
}

func (this *GLFWViewport) OnClose(f func()) EventSubscription {
	if this.onCloseEvent == nil {
		this.onCloseEvent = CreateEvent(func() {})
	}
	return this.onCloseEvent.Listen(f)
}

func (this *GLFWViewport) Size() Size {
	w, h := this.glfwWindow.GetSize()
	return Size{w, h}
}

func (this *GLFWViewport) Close() {
	this.glfwWindow.SetShouldClose(true)
	if this.onCloseEvent != nil {
		this.onCloseEvent.Fire()
	}
	this.glfwWindow.Destroy()
}

var glfwKeys []glfw.Key = []glfw.Key{
	//glfw.KeyUnknown,
	glfw.KeySpace,
	glfw.KeyApostrophe,
	glfw.KeyComma,
	glfw.KeyMinus,
	glfw.KeyPeriod,
	glfw.KeySlash,
	glfw.Key0,
	glfw.Key1,
	glfw.Key2,
	glfw.Key3,
	glfw.Key4,
	glfw.Key5,
	glfw.Key6,
	glfw.Key7,
	glfw.Key8,
	glfw.Key9,
	glfw.KeySemicolon,
	glfw.KeyEqual,
	glfw.KeyA,
	glfw.KeyB,
	glfw.KeyC,
	glfw.KeyD,
	glfw.KeyE,
	glfw.KeyF,
	glfw.KeyG,
	glfw.KeyH,
	glfw.KeyI,
	glfw.KeyJ,
	glfw.KeyK,
	glfw.KeyL,
	glfw.KeyM,
	glfw.KeyN,
	glfw.KeyO,
	glfw.KeyP,
	glfw.KeyQ,
	glfw.KeyR,
	glfw.KeyS,
	glfw.KeyT,
	glfw.KeyU,
	glfw.KeyV,
	glfw.KeyW,
	glfw.KeyX,
	glfw.KeyY,
	glfw.KeyZ,
	glfw.KeyLeftBracket,
	glfw.KeyBackslash,
	glfw.KeyRightBracket,
	glfw.KeyGraveAccent,
	glfw.KeyWorld1,
	glfw.KeyWorld2,
	glfw.KeyEscape,
	glfw.KeyEnter,
	glfw.KeyTab,
	glfw.KeyBackspace,
	glfw.KeyInsert,
	glfw.KeyDelete,
	glfw.KeyRight,
	glfw.KeyLeft,
	glfw.KeyDown,
	glfw.KeyUp,
	glfw.KeyPageUp,
	glfw.KeyPageDown,
	glfw.KeyHome,
	glfw.KeyEnd,
	glfw.KeyCapsLock,
	glfw.KeyScrollLock,
	glfw.KeyNumLock,
	glfw.KeyPrintScreen,
	glfw.KeyPause,
	glfw.KeyF1,
	glfw.KeyF2,
	glfw.KeyF3,
	glfw.KeyF4,
	glfw.KeyF5,
	glfw.KeyF6,
	glfw.KeyF7,
	glfw.KeyF8,
	glfw.KeyF9,
	glfw.KeyF10,
	glfw.KeyF11,
	glfw.KeyF12,
	glfw.KeyF13,
	glfw.KeyF14,
	glfw.KeyF15,
	glfw.KeyF16,
	glfw.KeyF17,
	glfw.KeyF18,
	glfw.KeyF19,
	glfw.KeyF20,
	glfw.KeyF21,
	glfw.KeyF22,
	glfw.KeyF23,
	glfw.KeyF24,
	glfw.KeyF25,
	glfw.KeyKP0,
	glfw.KeyKP1,
	glfw.KeyKP2,
	glfw.KeyKP3,
	glfw.KeyKP4,
	glfw.KeyKP5,
	glfw.KeyKP6,
	glfw.KeyKP7,
	glfw.KeyKP8,
	glfw.KeyKP9,
	glfw.KeyKPDecimal,
	glfw.KeyKPDivide,
	glfw.KeyKPMultiply,
	glfw.KeyKPSubtract,
	glfw.KeyKPAdd,
	glfw.KeyKPEnter,
	glfw.KeyKPEqual,
	glfw.KeyLeftShift,
	glfw.KeyLeftControl,
	glfw.KeyLeftAlt,
	glfw.KeyLeftSuper,
	glfw.KeyRightShift,
	glfw.KeyRightControl,
	glfw.KeyRightAlt,
	glfw.KeyRightSuper,
	glfw.KeyMenu,
	glfw.KeyLast}

func getKeyboardKeys(w *glfw.Window, typ event.Type) []KeyboardKey {
	kvs := make([]KeyboardKey, 0)
	for _, key := range glfwKeys {
		action := w.GetKey(key)

		switch action {
		case glfw.Press:
			//log.Info("Press:", action)
			if typ == KEY_PRESS_EVENT_TYPE {
				kvs = append(kvs, translateKeyboardKey(key))
			}
		case glfw.Release:
			//log.Info("Release:", action)
		case glfw.Repeat:
			//log.Info("Repeat:", action)
			if typ == KEY_REPEAT_EVENT_TYPE {
				kvs = append(kvs, translateKeyboardKey(key))
			}
		}
	}

	return kvs
}

func translateKeyboardKey(in glfw.Key) KeyboardKey {
	switch in {
	case glfw.KeySpace:
		return KeySpace
	case glfw.KeyApostrophe:
		return KeyApostrophe
	case glfw.KeyComma:
		return KeyComma
	case glfw.KeyMinus:
		return KeyMinus
	case glfw.KeyPeriod:
		return KeyPeriod
	case glfw.KeySlash:
		return KeySlash
	case glfw.Key0:
		return Key0
	case glfw.Key1:
		return Key1
	case glfw.Key2:
		return Key2
	case glfw.Key3:
		return Key3
	case glfw.Key4:
		return Key4
	case glfw.Key5:
		return Key5
	case glfw.Key6:
		return Key6
	case glfw.Key7:
		return Key7
	case glfw.Key8:
		return Key8
	case glfw.Key9:
		return Key9
	case glfw.KeySemicolon:
		return KeySemicolon
	case glfw.KeyEqual:
		return KeyEqual
	case glfw.KeyA:
		return KeyA
	case glfw.KeyB:
		return KeyB
	case glfw.KeyC:
		return KeyC
	case glfw.KeyD:
		return KeyD
	case glfw.KeyE:
		return KeyE
	case glfw.KeyF:
		return KeyF
	case glfw.KeyG:
		return KeyG
	case glfw.KeyH:
		return KeyH
	case glfw.KeyI:
		return KeyI
	case glfw.KeyJ:
		return KeyJ
	case glfw.KeyK:
		return KeyK
	case glfw.KeyL:
		return KeyL
	case glfw.KeyM:
		return KeyM
	case glfw.KeyN:
		return KeyN
	case glfw.KeyO:
		return KeyO
	case glfw.KeyP:
		return KeyP
	case glfw.KeyQ:
		return KeyQ
	case glfw.KeyR:
		return KeyR
	case glfw.KeyS:
		return KeyS
	case glfw.KeyT:
		return KeyT
	case glfw.KeyU:
		return KeyU
	case glfw.KeyV:
		return KeyV
	case glfw.KeyW:
		return KeyW
	case glfw.KeyX:
		return KeyX
	case glfw.KeyY:
		return KeyY
	case glfw.KeyZ:
		return KeyZ
	case glfw.KeyLeftBracket:
		return KeyLeftBracket
	case glfw.KeyBackslash:
		return KeyBackslash
	case glfw.KeyRightBracket:
		return KeyRightBracket
	case glfw.KeyGraveAccent:
		return KeyGraveAccent
	case glfw.KeyWorld1:
		return KeyWorld1
	case glfw.KeyWorld2:
		return KeyWorld2
	case glfw.KeyEscape:
		return KeyEscape
	case glfw.KeyEnter:
		return KeyEnter
	case glfw.KeyTab:
		return KeyTab
	case glfw.KeyBackspace:
		return KeyBackspace
	case glfw.KeyInsert:
		return KeyInsert
	case glfw.KeyDelete:
		return KeyDelete
	case glfw.KeyRight:
		return KeyRight
	case glfw.KeyLeft:
		return KeyLeft
	case glfw.KeyDown:
		return KeyDown
	case glfw.KeyUp:
		return KeyUp
	case glfw.KeyPageUp:
		return KeyPageUp
	case glfw.KeyPageDown:
		return KeyPageDown
	case glfw.KeyHome:
		return KeyHome
	case glfw.KeyEnd:
		return KeyEnd
	case glfw.KeyCapsLock:
		return KeyCapsLock
	case glfw.KeyScrollLock:
		return KeyScrollLock
	case glfw.KeyNumLock:
		return KeyNumLock
	case glfw.KeyPrintScreen:
		return KeyPrintScreen
	case glfw.KeyPause:
		return KeyPause
	case glfw.KeyF1:
		return KeyF1
	case glfw.KeyF2:
		return KeyF2
	case glfw.KeyF3:
		return KeyF3
	case glfw.KeyF4:
		return KeyF4
	case glfw.KeyF5:
		return KeyF5
	case glfw.KeyF6:
		return KeyF6
	case glfw.KeyF7:
		return KeyF7
	case glfw.KeyF8:
		return KeyF8
	case glfw.KeyF9:
		return KeyF9
	case glfw.KeyF10:
		return KeyF10
	case glfw.KeyF11:
		return KeyF11
	case glfw.KeyF12:
		return KeyF12
	case glfw.KeyKP0:
		return KeyKp0
	case glfw.KeyKP1:
		return KeyKp1
	case glfw.KeyKP2:
		return KeyKp2
	case glfw.KeyKP3:
		return KeyKp3
	case glfw.KeyKP4:
		return KeyKp4
	case glfw.KeyKP5:
		return KeyKp5
	case glfw.KeyKP6:
		return KeyKp6
	case glfw.KeyKP7:
		return KeyKp7
	case glfw.KeyKP8:
		return KeyKp8
	case glfw.KeyKP9:
		return KeyKp9
	case glfw.KeyKPDecimal:
		return KeyKpDecimal
	case glfw.KeyKPDivide:
		return KeyKpDivide
	case glfw.KeyKPMultiply:
		return KeyKpMultiply
	case glfw.KeyKPSubtract:
		return KeyKpSubtract
	case glfw.KeyKPAdd:
		return KeyKpAdd
	case glfw.KeyKPEnter:
		return KeyKpEnter
	case glfw.KeyKPEqual:
		return KeyKpEqual
	case glfw.KeyLeftShift:
		return KeyLeftShift
	case glfw.KeyLeftControl:
		return KeyLeftControl
	case glfw.KeyLeftAlt:
		return KeyLeftAlt
	case glfw.KeyLeftSuper:
		return KeyLeftSuper
	case glfw.KeyRightShift:
		return KeyRightShift
	case glfw.KeyRightControl:
		return KeyRightControl
	case glfw.KeyRightAlt:
		return KeyRightAlt
	case glfw.KeyRightSuper:
		return KeyRightSuper
	case glfw.KeyMenu:
		return KeyMenu
	default:
		return KeyUnknown
	}
}

func translateKeyboardModifier(in glfw.ModifierKey) KeyboardModifier {
	out := ModNone
	if in&glfw.ModShift != 0 {
		out |= ModShift
	}
	if in&glfw.ModControl != 0 {
		out |= ModControl
	}
	if in&glfw.ModAlt != 0 {
		out |= ModAlt
	}
	if in&glfw.ModSuper != 0 {
		out |= ModSuper
	}
	return out
}
