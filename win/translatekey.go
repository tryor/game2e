package win

import (
	//	. "tryor/game2e"

	. "github.com/tryor/eui"
	//	"github.com/tryor/winapi"
	. "github.com/tryor/winapi"
)

const (
	VK_CONTROL = 0x11
	VK_MENU    = 0x12
	VK_SHIFT   = 0x10
	VK_PAUSE   = 0x13
	VK_CAPITAL = 0x14

	VK_PROCESSKEY = 0xE5
	VK_SNAPSHOT   = 0x2C
)

const (
	PM_NOREMOVE = 0x0000
	PM_REMOVE   = 0x0001
	PM_NOYIELD  = 0x0002
)

const (
	WM_UNICHAR     = 0x0109
	UNICODE_NOCHAR = 0xFFFF
)

const (
	_KEY_INVALID = -2
)

type Action int8

const (
	RELEASE = 0
	PRESS   = 1
	REPEAT  = 2
	_STICK  = 3
)

const (
	MOD_SHIFT   = 0x0001
	MOD_CONTROL = 0x0002
	MOD_ALT     = 0x0004
	MOD_SUPER   = 0x0008
)

var publicKeys = map[INT]Key{}

func init() {
	publicKeys[0x00B] = KEY_0
	publicKeys[0x002] = KEY_1
	publicKeys[0x003] = KEY_2
	publicKeys[0x004] = KEY_3
	publicKeys[0x005] = KEY_4
	publicKeys[0x006] = KEY_5
	publicKeys[0x007] = KEY_6
	publicKeys[0x008] = KEY_7
	publicKeys[0x009] = KEY_8
	publicKeys[0x00A] = KEY_9
	publicKeys[0x01E] = KEY_A
	publicKeys[0x030] = KEY_B
	publicKeys[0x02E] = KEY_C
	publicKeys[0x020] = KEY_D
	publicKeys[0x012] = KEY_E
	publicKeys[0x021] = KEY_F
	publicKeys[0x022] = KEY_G
	publicKeys[0x023] = KEY_H
	publicKeys[0x017] = KEY_I
	publicKeys[0x024] = KEY_J
	publicKeys[0x025] = KEY_K
	publicKeys[0x026] = KEY_L
	publicKeys[0x032] = KEY_M
	publicKeys[0x031] = KEY_N
	publicKeys[0x018] = KEY_O
	publicKeys[0x019] = KEY_P
	publicKeys[0x010] = KEY_Q
	publicKeys[0x013] = KEY_R
	publicKeys[0x01F] = KEY_S
	publicKeys[0x014] = KEY_T
	publicKeys[0x016] = KEY_U
	publicKeys[0x02F] = KEY_V
	publicKeys[0x011] = KEY_W
	publicKeys[0x02D] = KEY_X
	publicKeys[0x015] = KEY_Y
	publicKeys[0x02C] = KEY_Z

	publicKeys[0x028] = KEY_APOSTROPHE
	publicKeys[0x02B] = KEY_BACKSLASH
	publicKeys[0x033] = KEY_COMMA
	publicKeys[0x00D] = KEY_EQUAL
	publicKeys[0x029] = KEY_GRAVE_ACCENT
	publicKeys[0x01A] = KEY_LEFT_BRACKET
	publicKeys[0x00C] = KEY_MINUS
	publicKeys[0x034] = KEY_PERIOD
	publicKeys[0x01B] = KEY_RIGHT_BRACKET
	publicKeys[0x027] = KEY_SEMICOLON
	publicKeys[0x035] = KEY_SLASH
	publicKeys[0x056] = KEY_WORLD_2

	publicKeys[0x00E] = KEY_BACKSPACE
	publicKeys[0x153] = KEY_DELETE
	publicKeys[0x14F] = KEY_END
	publicKeys[0x01C] = KEY_ENTER
	publicKeys[0x001] = KEY_ESCAPE
	publicKeys[0x147] = KEY_HOME
	publicKeys[0x152] = KEY_INSERT
	publicKeys[0x15D] = KEY_MENU
	publicKeys[0x151] = KEY_PAGE_DOWN
	publicKeys[0x149] = KEY_PAGE_UP
	publicKeys[0x045] = KEY_PAUSE
	publicKeys[0x146] = KEY_PAUSE
	publicKeys[0x039] = KEY_SPACE
	publicKeys[0x00F] = KEY_TAB
	publicKeys[0x03A] = KEY_CAPS_LOCK
	publicKeys[0x145] = KEY_NUM_LOCK
	publicKeys[0x046] = KEY_SCROLL_LOCK
	publicKeys[0x03B] = KEY_F1
	publicKeys[0x03C] = KEY_F2
	publicKeys[0x03D] = KEY_F3
	publicKeys[0x03E] = KEY_F4
	publicKeys[0x03F] = KEY_F5
	publicKeys[0x040] = KEY_F6
	publicKeys[0x041] = KEY_F7
	publicKeys[0x042] = KEY_F8
	publicKeys[0x043] = KEY_F9
	publicKeys[0x044] = KEY_F10
	publicKeys[0x057] = KEY_F11
	publicKeys[0x058] = KEY_F12
	publicKeys[0x064] = KEY_F13
	publicKeys[0x065] = KEY_F14
	publicKeys[0x066] = KEY_F15
	publicKeys[0x067] = KEY_F16
	publicKeys[0x068] = KEY_F17
	publicKeys[0x069] = KEY_F18
	publicKeys[0x06A] = KEY_F19
	publicKeys[0x06B] = KEY_F20
	publicKeys[0x06C] = KEY_F21
	publicKeys[0x06D] = KEY_F22
	publicKeys[0x06E] = KEY_F23
	publicKeys[0x076] = KEY_F24
	publicKeys[0x038] = KEY_LEFT_ALT
	publicKeys[0x01D] = KEY_LEFT_CONTROL
	publicKeys[0x02A] = KEY_LEFT_SHIFT
	publicKeys[0x15B] = KEY_LEFT_SUPER
	publicKeys[0x137] = KEY_PRINT_SCREEN
	publicKeys[0x138] = KEY_RIGHT_ALT
	publicKeys[0x11D] = KEY_RIGHT_CONTROL
	publicKeys[0x036] = KEY_RIGHT_SHIFT
	publicKeys[0x15C] = KEY_RIGHT_SUPER
	publicKeys[0x150] = KEY_DOWN
	publicKeys[0x14B] = KEY_LEFT
	publicKeys[0x14D] = KEY_RIGHT
	publicKeys[0x148] = KEY_UP

	publicKeys[0x052] = KEY_KP_0
	publicKeys[0x04F] = KEY_KP_1
	publicKeys[0x050] = KEY_KP_2
	publicKeys[0x051] = KEY_KP_3
	publicKeys[0x04B] = KEY_KP_4
	publicKeys[0x04C] = KEY_KP_5
	publicKeys[0x04D] = KEY_KP_6
	publicKeys[0x047] = KEY_KP_7
	publicKeys[0x048] = KEY_KP_8
	publicKeys[0x049] = KEY_KP_9
	publicKeys[0x04E] = KEY_KP_ADD
	publicKeys[0x053] = KEY_KP_DECIMAL
	publicKeys[0x135] = KEY_KP_DIVIDE
	publicKeys[0x11C] = KEY_KP_ENTER
	publicKeys[0x037] = KEY_KP_MULTIPLY
	publicKeys[0x04A] = KEY_KP_SUBTRACT
}

func translateKey(wparam WPARAM, lparam LPARAM) Key {

	if wparam == VK_CONTROL {

		if (lparam & 0x01000000) > 0 {
			return KEY_RIGHT_CONTROL
		}

		var next Msg
		var time DWORD

		time = GetMessageTime()

		if PeekMessage(&next, 0, 0, 0, PM_NOREMOVE) {
			if next.Message == WM_KEYDOWN ||
				next.Message == WM_SYSKEYDOWN ||
				next.Message == WM_KEYUP ||
				next.Message == WM_SYSKEYUP {

				if next.Wparam == VK_MENU &&
					(next.Lparam&0x01000000) > 0 &&
					next.Time == uint32(time) {
					return _KEY_INVALID
				}
			}
		}
		return KEY_LEFT_CONTROL
	}

	if wparam == VK_PROCESSKEY {
		return _KEY_INVALID
	}

	return publicKeys[HIWORD(INT(lparam))&0x1FF]
}

//func getKeyMods() KeyboardModifier {
//	out := ModNone

//	return out
//}

//func translateKeyboardModifier(in glfw.ModifierKey) KeyboardModifier {
//	out := ModNone
//	if in&glfw.ModShift != 0 {
//		out |= ModShift
//	}
//	if in&glfw.ModControl != 0 {
//		out |= ModControl
//	}
//	if in&glfw.ModAlt != 0 {
//		out |= ModAlt
//	}
//	if in&glfw.ModSuper != 0 {
//		out |= ModSuper
//	}
//	return out
//}

/*
void _glfwInputKey(_GLFWwindow* window, int key, int scancode, int action, int mods)
{
    if (key >= 0 && key <= KEY_LAST)
    {
        GLFWbool repeated = FALSE;

        if (action == RELEASE && window->keys[key] == RELEASE)
            return;

        if (action == PRESS && window->keys[key] == PRESS)
            repeated = TRUE;

        if (action == RELEASE && window->stickyKeys)
            window->keys[key] = _STICK;
        else
            window->keys[key] = (char) action;

        if (repeated)
            action = REPEAT;
    }

    if (window->callbacks.key)
        window->callbacks.key((GLFWwindow*) window, key, scancode, action, mods);
}



static int getKeyMods(void)
{
    int mods = 0;

    if (GetKeyState(VK_SHIFT) & (1 << 31))
        mods |= MOD_SHIFT;
    if (GetKeyState(VK_CONTROL) & (1 << 31))
        mods |= MOD_CONTROL;
    if (GetKeyState(VK_MENU) & (1 << 31))
        mods |= MOD_ALT;
    if ((GetKeyState(VK_LWIN) | GetKeyState(VK_RWIN)) & (1 << 31))
        mods |= MOD_SUPER;

    return mods;
}

*/
