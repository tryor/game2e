package win

import (
	"bytes"
	"errors"
	"sync"
	"syscall"
	. "tryor/game2e"
	"tryor/game2e/log"
	"unsafe"

	"code.google.com/p/mahonia"

	"github.com/tryor/commons/event"
	. "github.com/tryor/eui"
	"github.com/tryor/winapi"
	. "github.com/tryor/winapi"
)

func wndProc(hwnd HWND, msg UINT, wparam WPARAM, lparam LPARAM) (rc uintptr) {
	//log.Info(hwnd, msg, wparam, lparam)
	var w *Win32Viewport
	var ok bool
	if w, ok = windows[hwnd]; !ok {
		return DefWindowProcW(hwnd, msg, wparam, lparam)
	}

	switch msg {
	//case WM_CREATE:
	//	log.Info("WM_CREATE ", hwnd)
	//case WM_COMMAND:
	//	log.Info("WM_COMMAND ", hwnd)
	case WM_PAINT:

	case WM_SHOWWINDOW:
		//log.Info("WM_SHOWWINDOW ", hwnd)
		w.procMsgShow()
	case WM_CLOSE:
		//log.Info("WM_CLOSE ", hwnd)
		w.Hide()
		w.procMsgClose()
		procMsgCloses(w)
		return 0

	case WM_MOUSEMOVE:
		//appn.TrackMouseMoveEvent(int(LOWORD(INT(lparam))), int(HIWORD(INT(lparam))), canvas.MButton(wparam))
		w.procMsgMouseMove(int(LOWORD(INT(lparam))), int(HIWORD(INT(lparam))), MButton(wparam))
	case WM_LBUTTONDOWN, WM_RBUTTONDOWN:
		w.procMsgMouseDown(int(LOWORD(INT(lparam))), int(HIWORD(INT(lparam))), MButton(wparam))
	case WM_LBUTTONUP, WM_RBUTTONUP:
		w.procMsgMouseUp(int(LOWORD(INT(lparam))), int(HIWORD(INT(lparam))), MButton(wparam))
	case WM_LBUTTONDBLCLK, WM_RBUTTONDBLCLK:

	case WM_KEYDOWN, WM_SYSKEYDOWN, WM_KEYUP, WM_SYSKEYUP:
		FireKeyEvent(w.onKeyEvent, w.keys[:], w.stickyKeys, wparam, lparam)
		/*

			key := translateKey(wparam, lparam)
			scancode := int((lparam >> 16) & 0x1ff)
			var action Action
			if ((lparam >> 31) & 1) != 0 {
				action = RELEASE
			} else {
				action = PRESS
			}
			if key == _KEY_INVALID {
				break
			}
			mods := getKeyMods()
			if action == RELEASE && wparam == VK_SHIFT {
				// Release both Shift keys on Shift up event, as only one event
				// is sent even if both keys are released
				w.inputKey(KEY_LEFT_SHIFT, scancode, action, mods)
				w.inputKey(KEY_RIGHT_SHIFT, scancode, action, mods)
			} else if wparam == VK_SNAPSHOT {
				// Key down is not reported for the Print Screen key
				w.inputKey(key, scancode, PRESS, mods)
				w.inputKey(key, scancode, RELEASE, mods)
			} else {
				w.inputKey(key, scancode, action, mods)
			}

		*/

	case WM_SYSCHAR:
		fallthrough
	case WM_UNICHAR:
		fallthrough
	case WM_CHAR:
		//plain := (msg != WM_SYSCHAR)
		if msg == WM_UNICHAR && wparam == UNICODE_NOCHAR {
			return 1
		}

		FireCharEvent(w.onKeyCharEvent, uint(wparam))
		//w.inputChar(uint(wparam), getModifier(), plain)

		return 0
	}
	//log.Info(w)
	_ = w
	return DefWindowProcW(hwnd, msg, wparam, lparam)
}

var windows = map[HWND]*Win32Viewport{}
var firstWindow *Win32Viewport
var moduleHandle HINSTANCE
var locker sync.Mutex

func init() {

}

func createWin32Viewport(width, height int, title string) (IViewport, error) {
	locker.Lock()
	defer locker.Unlock()

	wcname := "Game2e Window Class"

	var wh HWND
	var err error
	if firstWindow == nil {
		moduleHandle, err = GetModuleHandle("")
		if err != nil {
			log.Error(err)
			return nil, err
		}

		myicon, err := LoadIcon(0, IDI_APPLICATION)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		mycursor, err := LoadCursor(0, IDC_ARROW)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		// RegisterClassEx
		var wc Wndclassex
		wc.Size = uint32(unsafe.Sizeof(wc))
		wc.Style = CS_HREDRAW | CS_VREDRAW
		wc.WndProc = syscall.NewCallback(wndProc)
		wc.Instance = HINSTANCE(moduleHandle)
		wc.Icon = myicon
		wc.Cursor = mycursor
		wc.Background = COLOR_BTNFACE + 1
		wc.MenuName = nil
		wc.ClassName = syscall.StringToUTF16Ptr(wcname)
		wc.IconSm = myicon
		if _, err := RegisterClassExW(&wc); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	wh, err = CreateWindowExW(
		0, //WS_EX_TOOLWINDOW,
		wcname,
		title,
		WS_OVERLAPPEDWINDOW, //|WS_VISIBLE,
		CW_USEDEFAULT, CW_USEDEFAULT,
		int32(width)+2, int32(height)+32,
		0, 0, HINSTANCE(moduleHandle), 0)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	viewport := &Win32Viewport{winHwnd: wh, width: width, height: height}
	viewport.initEvents()
	windows[wh] = viewport

	if err := UpdateWindow(wh); err != nil {
		log.Error(err)
		return nil, err
	}

	ShowWindow(wh, SW_SHOW)

	if firstWindow == nil {
		firstWindow = viewport
	}

	return viewport, nil
}

type Win32Viewport struct {
	IViewport
	stickyKeys    bool
	keys          [KEY_LAST + 1]Action
	winHwnd       HWND
	width, height int

	onResizeEvent  Event
	onCloseEvent   Event
	onMouseEvent   Event
	onKeyEvent     Event
	onKeyCharEvent Event
}

func (this *Win32Viewport) initEvents() {

}

func (v *Win32Viewport) GetHandle() HANDLE {
	return HANDLE(v.winHwnd)
}

func (v *Win32Viewport) OnKeyChar(f func(*KeyEvent)) EventSubscription {
	if v.onKeyCharEvent == nil {
		v.onKeyCharEvent = CreateEvent(func(*KeyEvent) {})
	}
	return v.onKeyCharEvent.Listen(f)
}

func (v *Win32Viewport) OnKey(f func(*KeyEvent)) EventSubscription {
	if v.onKeyEvent == nil {
		v.onKeyEvent = CreateEvent(func(*KeyEvent) {})
	}
	return v.onKeyEvent.Listen(f)
}
func (v *Win32Viewport) OnMouse(f func(IMouseEvent)) EventSubscription {
	if v.onMouseEvent == nil {
		v.onMouseEvent = CreateEvent(func(IMouseEvent) {})
	}
	return v.onMouseEvent.Listen(f)
}
func (v *Win32Viewport) OnResize(f func()) EventSubscription {
	if v.onResizeEvent == nil {
		v.onResizeEvent = CreateEvent(func() {})
	}
	return v.onResizeEvent.Listen(f)
}
func (v *Win32Viewport) OnClose(f func()) EventSubscription {
	if v.onCloseEvent == nil {
		v.onCloseEvent = CreateEvent(func() {})
	}
	return v.onCloseEvent.Listen(f)
}

func (v *Win32Viewport) Size() Size {
	rect, err := winapi.GetWindowRect(v.winHwnd)
	if err != nil {
		log.Error(err)
		return Size{v.width, v.height}
	} else {
		return Size{int(rect.Right - rect.Left), int(rect.Bottom - rect.Top)}
	}

}
func (v *Win32Viewport) Close() {
	winapi.PostMessage(v.winHwnd, WM_CLOSE, 0, 0)
}

func (v *Win32Viewport) SetVisible(b bool) {
	if b {
		winapi.ShowWindow(v.winHwnd, winapi.SW_SHOW)
	} else {
		winapi.ShowWindow(v.winHwnd, winapi.SW_HIDE)
	}
}

func (v *Win32Viewport) Hide() {
	v.SetVisible(false)
}

func (v *Win32Viewport) IsHide() bool {
	long, err := winapi.GetWindowLongPtr(v.winHwnd, GWL_STYLE)
	if err != nil {
		log.Warn(err)
		return true
	}
	return long&WS_VISIBLE == 0
}

//-----------------------------------------------

func (w *Win32Viewport) procMsgMouseMove(x, y int, buttonsStatus MButton) {
	if w.onMouseEvent != nil {
		me := NewMouseEvent(MOUSE_MOVE_EVENT_TYPE, w, x, y, buttonsStatus, getModifier())
		//me.KeySequence = NewKeySequence(w.getKeyboardKeys()...)
		w.onMouseEvent.Fire(me)
	}
}

func (w *Win32Viewport) procMsgMouseDown(x, y int, buttonsStatus MButton) {
	if w.onMouseEvent != nil {
		me := NewMouseEvent(MOUSE_PRESS_EVENT_TYPE, w, x, y, buttonsStatus, getModifier())
		//me.KeySequence = NewKeySequence(w.getKeyboardKeys()...)
		w.onMouseEvent.Fire(me)
	}
}

func (w *Win32Viewport) procMsgMouseUp(x, y int, buttonsStatus MButton) {
	if w.onMouseEvent != nil {
		me := NewMouseEvent(MOUSE_RELEASE_EVENT_TYPE, w, x, y, buttonsStatus, getModifier())
		//me.KeySequence = NewKeySequence(w.getKeyboardKeys()...)
		w.onMouseEvent.Fire(me)
	}
}

//func (w *Win32Viewport) procMsgKeydown(key Key) {
//	w.procMsgKey(KEY_PRESS_EVENT_TYPE, key)
//}

//func (w *Win32Viewport) procMsgKeyup(key Key) {
//	w.procMsgKey(KEY_RELEASE_EVENT_TYPE, key)
//}

//func (w *Win32Viewport) procMsgKeyrepeat(key Key) {
//	w.procMsgKey(KEY_REPEAT_EVENT_TYPE, key)
//}

//func (w *Win32Viewport) procMsgKey(typ event.Type, key Key) {

//	if w.onKeyEvent == nil {
//		return
//	}
//	log.Debug(key)
//	kv := TranslateKeyboardKey(key)
//	kvs := append(w.getKeyboardKeys(), kv)
//	ke := NewKeyEvent(typ, w, getModifier(), kvs...)
//	ke.Key = kv

//	w.onKeyEvent.Fire(ke)

//}

func (w *Win32Viewport) inputKey(key Key, scancode int, action Action, mods KeyboardModifier) {

	if key >= 0 && key <= KEY_LAST {
		repeated := false
		if action == RELEASE && w.keys[key] == RELEASE {
			return
		}
		if action == PRESS && w.keys[key] == PRESS {
			repeated = true
		}
		if action == RELEASE && w.stickyKeys {
			w.keys[key] = _STICK
		} else {
			w.keys[key] = action
		}
		if repeated {
			action = REPEAT
		}
	}

	if w.onKeyEvent == nil {
		return
	}

	var typ event.Type
	switch action {
	case PRESS:
		typ = KEY_PRESS_EVENT_TYPE
	case RELEASE:
		typ = KEY_RELEASE_EVENT_TYPE
	case REPEAT:
		typ = KEY_REPEAT_EVENT_TYPE
	}

	kv := TranslateKeyboardKey(key)
	kvs := getKeyboardKeys(typ, key, w.keys[:]) // append(, kv)
	ke := NewKeyEvent(typ, w, mods, kvs...)
	ke.Key = kv
	w.onKeyEvent.Fire(ke)
	//log.Debug("kvs:", kvs, kv, key)
}

func (w *Win32Viewport) inputChar(codepoint uint, mods KeyboardModifier, plain bool) {
	if codepoint < 32 || (codepoint > 126 && codepoint < 160) {
		return
	}

	if w.onKeyCharEvent == nil {
		return
	}
	s := string(codepoint)
	enc := mahonia.NewDecoder("utf-8")
	v := enc.ConvertString(s)
	var c rune
	if v != "" {
		runes := bytes.Runes([]byte(v))
		c = runes[0]
	}
	kce := NewKeyCharEvent(KEY_CHAR_EVENT_TYPE, w, c, getModifier())
	w.onKeyCharEvent.Fire(kce)
}

/*
func (w *Win32Viewport) procMsgChar(char int) {
	if w.onKeyCharEvent == nil {
		return
	}
	s := string(char)
	enc := mahonia.NewDecoder("utf-8")
	v := enc.ConvertString(s)
	var c rune
	if v != "" {
		runes := bytes.Runes([]byte(v))
		c = runes[0]
	}
	//kce := NewKeyCharEvent(KEY_CHAR_EVENT_TYPE, nil, c, getModifier(), w.getKeyboardKeys()...)
	kce := NewKeyCharEvent(KEY_CHAR_EVENT_TYPE, w, c, getModifier())
	w.onKeyCharEvent.Fire(kce)
}
*/

func (w *Win32Viewport) procMsgShow() {

}

func (w *Win32Viewport) procMsgClose() {
	if w.onCloseEvent != nil {
		w.onCloseEvent.Fire()
	}
}

func getKeyboardKeys(typ event.Type, k Key, keysStatus []Action) []KeyboardKey {
	kvs := make([]KeyboardKey, 0)
	var has bool
	for _, key := range Keys {
		if int(key) >= 0 && int(key) < len(keysStatus) {
			action := keysStatus[key]
			switch action {
			case PRESS:
				if typ == KEY_PRESS_EVENT_TYPE {
					kvs = append(kvs, TranslateKeyboardKey(key))
					if k == key {
						has = true
					}
				}
			case RELEASE:
			case REPEAT:
				if typ == KEY_REPEAT_EVENT_TYPE {
					kvs = append(kvs, TranslateKeyboardKey(key))
					if k == key {
						has = true
					}
				}
			}
		}
	}
	if !has {
		kvs = append(kvs, TranslateKeyboardKey(k))
	}
	return kvs
}

func procMsgCloses(cw *Win32Viewport) {
	for _, w := range windows {
		if !w.IsHide() {
			return
		}
	}
	PostQuitMessage(0)
}

func getModifier() KeyboardModifier {
	return ModNone
}

func FireCharEvent(charEvent Event, codepoint uint) {
	//	if codepoint < 32 || (codepoint > 126 && codepoint < 160) {
	//		return
	//	}

	if charEvent == nil {
		return
	}

	//mods := getKeyMods()

	s := string(codepoint)
	enc := mahonia.NewDecoder("utf-8")
	v := enc.ConvertString(s)
	var c rune
	if v != "" {
		runes := bytes.Runes([]byte(v))
		c = runes[0]
	}
	kce := NewKeyCharEvent(KEY_CHAR_EVENT_TYPE, nil, c, getModifier())
	charEvent.Fire(kce)
}

func FireKeyEvent(keyEvent Event, keysStatus []Action, stickyKeys bool, wparam WPARAM, lparam LPARAM) error {
	key := translateKey(wparam, lparam)
	scancode := int((lparam >> 16) & 0x1ff)
	var action Action
	if ((lparam >> 31) & 1) != 0 {
		action = RELEASE
	} else {
		action = PRESS
	}
	if key == _KEY_INVALID {
		return errors.New("Key invalid")
	}
	mods := getModifier()
	if action == RELEASE && wparam == VK_SHIFT {
		// Release both Shift keys on Shift up event, as only one event
		// is sent even if both keys are released
		fireKeyEvent(keyEvent, keysStatus, stickyKeys, KEY_LEFT_SHIFT, scancode, action, mods)
		fireKeyEvent(keyEvent, keysStatus, stickyKeys, KEY_RIGHT_SHIFT, scancode, action, mods)
	} else if wparam == VK_SNAPSHOT {
		// Key down is not reported for the Print Screen key
		fireKeyEvent(keyEvent, keysStatus, stickyKeys, key, scancode, PRESS, mods)
		fireKeyEvent(keyEvent, keysStatus, stickyKeys, key, scancode, RELEASE, mods)
	} else {
		fireKeyEvent(keyEvent, keysStatus, stickyKeys, key, scancode, action, mods)
	}
	return nil
}

func fireKeyEvent(keyEvent Event, keysStatus []Action, stickyKeys bool, key Key, scancode int, action Action, mods KeyboardModifier) {

	if key >= 0 && key <= KEY_LAST {
		repeated := false
		if action == RELEASE && keysStatus[key] == RELEASE {
			return
		}
		if action == PRESS && keysStatus[key] == PRESS {
			repeated = true
		}
		if action == RELEASE && stickyKeys {
			keysStatus[key] = _STICK
		} else {
			keysStatus[key] = action
		}
		if repeated {
			action = REPEAT
		}
	}

	if keyEvent == nil {
		return
	}

	var typ event.Type
	switch action {
	case PRESS:
		typ = KEY_PRESS_EVENT_TYPE
	case RELEASE:
		typ = KEY_RELEASE_EVENT_TYPE
	case REPEAT:
		typ = KEY_REPEAT_EVENT_TYPE
	}

	kv := TranslateKeyboardKey(key)
	kvs := getKeyboardKeys(typ, key, keysStatus)
	ke := NewKeyEvent(typ, nil, mods, kvs...)
	ke.Key = kv
	keyEvent.Fire(ke)
}
