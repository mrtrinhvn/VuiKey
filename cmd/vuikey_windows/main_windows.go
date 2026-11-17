//go:build windows

package main

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/mrtrinhvn/VuiKey/internal/config"
	"github.com/mrtrinhvn/VuiKey/internal/engine"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	advapi32 = syscall.NewLazyDLL("advapi32.dll")
	dwmapi   = syscall.NewLazyDLL("dwmapi.dll")

	procGetModuleHandleW        = kernel32.NewProc("GetModuleHandleW")
	procCreateMutexW            = kernel32.NewProc("CreateMutexW")
	procGetLastError            = kernel32.NewProc("GetLastError")
	procRegisterClassExW        = user32.NewProc("RegisterClassExW")
	procCreateWindowExW         = user32.NewProc("CreateWindowExW")
	procDefWindowProcW          = user32.NewProc("DefWindowProcW")
	procShowWindow              = user32.NewProc("ShowWindow")
	procUpdateWindow            = user32.NewProc("UpdateWindow")
	procSetForegroundWindow     = user32.NewProc("SetForegroundWindow")
	procGetMessageW             = user32.NewProc("GetMessageW")
	procTranslateMessage        = user32.NewProc("TranslateMessage")
	procDispatchMessageW        = user32.NewProc("DispatchMessageW")
	procPostQuitMessage         = user32.NewProc("PostQuitMessage")
	procPostMessageW            = user32.NewProc("PostMessageW")
	procSendMessageW            = user32.NewProc("SendMessageW")
	procSetWindowTextW          = user32.NewProc("SetWindowTextW")
	procSetWindowsHookExW       = user32.NewProc("SetWindowsHookExW")
	procUnhookWindowsHookEx     = user32.NewProc("UnhookWindowsHookEx")
	procCallNextHookEx          = user32.NewProc("CallNextHookEx")
	procGetKeyState             = user32.NewProc("GetKeyState")
	procGetAsyncKeyState        = user32.NewProc("GetAsyncKeyState")
	procSendInput               = user32.NewProc("SendInput")
	procCreatePopupMenu         = user32.NewProc("CreatePopupMenu")
	procAppendMenuW             = user32.NewProc("AppendMenuW")
	procTrackPopupMenu          = user32.NewProc("TrackPopupMenu")
	procDestroyMenu             = user32.NewProc("DestroyMenu")
	procGetCursorPos            = user32.NewProc("GetCursorPos")
	procCreateIconIndirect      = user32.NewProc("CreateIconIndirect")
	procGetDC                   = user32.NewProc("GetDC")
	procReleaseDC               = user32.NewProc("ReleaseDC")
	procFillRect                = user32.NewProc("FillRect")
	procSetBkMode               = gdi32.NewProc("SetBkMode")
	procSetTextColor            = gdi32.NewProc("SetTextColor")
	procSetBkColor              = gdi32.NewProc("SetBkColor")
	procCreateSolidBrush        = gdi32.NewProc("CreateSolidBrush")
	procCreateFontW             = gdi32.NewProc("CreateFontW")
	procSelectObject            = gdi32.NewProc("SelectObject")
	procDeleteObject            = gdi32.NewProc("DeleteObject")
	procCreateCompatibleDC      = gdi32.NewProc("CreateCompatibleDC")
	procCreateBitmap            = gdi32.NewProc("CreateBitmap")
	procCreateDIBSection        = gdi32.NewProc("CreateDIBSection")
	procDeleteDC                = gdi32.NewProc("DeleteDC")
	procShellNotifyIconW        = shell32.NewProc("Shell_NotifyIconW")
	procShellExecuteW           = shell32.NewProc("ShellExecuteW")
	procRegOpenKeyExW           = advapi32.NewProc("RegOpenKeyExW")
	procRegSetValueExW          = advapi32.NewProc("RegSetValueExW")
	procRegDeleteValueW         = advapi32.NewProc("RegDeleteValueW")
	procRegQueryValueExW        = advapi32.NewProc("RegQueryValueExW")
	procRegCloseKey             = advapi32.NewProc("RegCloseKey")
	procDwmSetWindowAttribute   = dwmapi.NewProc("DwmSetWindowAttribute")
)

const (
	WH_KEYBOARD_LL = 13
	WH_MOUSE_LL    = 14

	WM_DESTROY       = 0x0002
	WM_CLOSE         = 0x0010
	WM_COMMAND       = 0x0111
	WM_CTLCOLORSTATIC = 0x0138
	WM_CTLCOLOREDIT  = 0x0133
	WM_CTLCOLORBTN   = 0x0135
	WM_KEYDOWN       = 0x0100
	WM_KEYUP         = 0x0101
	WM_SYSKEYDOWN    = 0x0104
	WM_LBUTTONDOWN   = 0x0201
	WM_RBUTTONDOWN   = 0x0204
	WM_MBUTTONDOWN   = 0x0207
	WM_LBUTTONUP     = 0x0202
	WM_RBUTTONUP     = 0x0205
	WM_LBUTTONDBLCLK = 0x0203
	WM_USER          = 0x0400
	WM_TRAYICON      = WM_USER + 1
	WM_SYNC_UI       = WM_USER + 2

	WS_OVERLAPPED  = 0x00000000
	WS_CAPTION     = 0x00C00000
	WS_SYSMENU     = 0x00080000
	WS_MINIMIZEBOX = 0x00020000
	WS_VISIBLE     = 0x10000000
	WS_CHILD       = 0x40000000
	WS_TABSTOP     = 0x00010000
	WS_BORDER      = 0x00800000
	ES_AUTOHSCROLL = 0x0080

	BS_PUSHBUTTON   = 0x00000000
	BS_AUTOCHECKBOX = 0x00000003
	BM_GETCHECK     = 0x00F0
	BM_SETCHECK     = 0x00F1
	BST_CHECKED     = 1
	BST_UNCHECKED   = 0

	CBS_DROPDOWNLIST = 0x0003
	CB_ADDSTRING     = 0x0143
	CB_SETCURSEL     = 0x014E
	CB_GETCURSEL     = 0x0147
	CBN_SELCHANGE    = 1

	SW_HIDE    = 0
	SW_SHOW    = 5
	SW_RESTORE = 9

	NIM_ADD    = 0x00000000
	NIM_MODIFY = 0x00000001
	NIM_DELETE = 0x00000002
	NIF_MESSAGE = 0x00000001
	NIF_ICON    = 0x00000002
	NIF_TIP     = 0x00000004

	MF_STRING    = 0x00000000
	MF_SEPARATOR = 0x00000800
	MF_CHECKED   = 0x00000008
	TPM_BOTTOMALIGN = 0x0020
	TPM_LEFTALIGN   = 0x0000

	INPUT_KEYBOARD        = 1
	KEYEVENTF_KEYUP       = 0x0002
	KEYEVENTF_UNICODE     = 0x0004
	LLKHF_INJECTED        = 0x00000010
	VUIKEY_EXTRA_INFO     = 0x5655494B

	VK_BACK    = 0x08
	VK_TAB     = 0x09
	VK_RETURN  = 0x0D
	VK_SHIFT   = 0x10
	VK_CONTROL = 0x11
	VK_MENU    = 0x12
	VK_CAPITAL = 0x14
	VK_ESCAPE  = 0x1B
	VK_SPACE   = 0x20
	VK_PRIOR   = 0x21
	VK_NEXT    = 0x22
	VK_END     = 0x23
	VK_HOME    = 0x24
	VK_LEFT    = 0x25
	VK_UP      = 0x26
	VK_RIGHT   = 0x27
	VK_DOWN    = 0x28
	VK_DELETE  = 0x2E
	VK_LWIN    = 0x5B
	VK_RWIN    = 0x5C

	// Control IDs
	ID_BTN_TOGGLE     = 101
	ID_CMB_MODE       = 102
	ID_CHK_SPLIT      = 103
	ID_CHK_SMART_ENG  = 104
	ID_CHK_STARTUP    = 105
	ID_BTN_CLEAR_TEST = 106
	ID_BTN_NEXT_QUOTE = 107
	ID_BTN_VUICAHAT   = 108
	ID_BTN_ECHIP      = 109
	ID_BTN_HIDE       = 110

	// Tray Menu IDs
	ID_TRAY_SHOW      = 201
	ID_TRAY_VN_TELEX  = 202
	ID_TRAY_VN_VNI    = 203
	ID_TRAY_EN        = 204
	ID_TRAY_STARTUP   = 205
	ID_TRAY_VUICAHAT  = 206
	ID_TRAY_ECHIP     = 207
	ID_TRAY_EXIT      = 208

	HKEY_CURRENT_USER = 0x80000001
	KEY_READ          = 0x20019
	KEY_WRITE         = 0x20006
	REG_SZ            = 1
)

type WNDCLASSEXW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       uintptr
}

type POINT struct {
	X, Y int32
}

type MSG struct {
	HWnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

type KEYBDINPUT struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

type INPUT struct {
	Type uint32
	_    uint32 // 64-bit alignment padding before union
	Ki   KEYBDINPUT
	_    [8]byte // union padding to match sizeof(INPUT) = 40 bytes on x64
}

type NOTIFYICONDATAW struct {
	CbSize           uint32
	HWnd             uintptr
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	HIcon            uintptr
	SzTip            [128]uint16
	DwState          uint32
	DwStateMask      uint32
	SzInfo           [256]uint16
	UVersion         uint32
	SzInfoTitle      [64]uint16
	DwInfoFlags      uint32
	GuidItem         [16]byte
	HBalloonIcon     uintptr
}

type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

type ICONINFO struct {
	FIcon    int32
	XHotspot uint32
	YHotspot uint32
	HbmMask  uintptr
	HbmColor uintptr
}

var (
	hInst          uintptr
	hMainWnd       uintptr
	hBtnToggle     uintptr
	hCmbMode       uintptr
	hChkSplit      uintptr
	hChkSmartEng   uintptr
	hChkStartup    uintptr
	hLblQuote      uintptr
	hLblHint       uintptr
	hEditTest      uintptr
	hLblStats      uintptr
	hIconVN        uintptr
	hIconEN        uintptr
	hKeyHook       uintptr
	hMouseHook     uintptr
	hBgBrush       uintptr
	hCardBrush     uintptr
	hEditBrush     uintptr

	eng            *engine.Engine
	totalKeys      int64
	transformedCnt int64
	sampleIdx      int
	nid            NOTIFYICONDATAW
)

func u16(s string) *uint16 {
	p, _ := syscall.UTF16PtrFromString(s)
	return p
}

// createBadgeIcon generates a crisp 32x32 ARGB icon in memory ("VN" on emerald #10B981 or "EN" on slate #475569).
func createBadgeIcon(isVN bool) uintptr {
	const w, h = 32, 32
	var bih BITMAPINFOHEADER
	bih.BiSize = uint32(unsafe.Sizeof(bih))
	bih.BiWidth = w
	bih.BiHeight = -h // top-down DIB
	bih.BiPlanes = 1
	bih.BiBitCount = 32

	var bits unsafe.Pointer
	hdcScreen, _, _ := procGetDC.Call(0)
	hbmColor, _, _ := procCreateDIBSection.Call(hdcScreen, uintptr(unsafe.Pointer(&bih)), 0, uintptr(unsafe.Pointer(&bits)), 0, 0)
	procReleaseDC.Call(0, hdcScreen)
	if hbmColor == 0 || bits == nil {
		return 0
	}

	pixels := unsafe.Slice((*uint32)(bits), w*h)
	bgColor := uint32(0xFF475569) // Slate gray #475569
	borderColor := uint32(0xFF334155)
	if isVN {
		bgColor = uint32(0xFF10B981) // Emerald #10B981
		borderColor = uint32(0xFF059669)
	}

	// Draw rounded rectangle badge
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if y < 2 || y >= h-2 || x < 1 || x >= w-1 {
				continue
			}
			// Corner cutouts
			dx, dy := 0, 0
			if x < 5 {
				dx = 5 - x
			} else if x >= w-5 {
				dx = x - (w - 6)
			}
			if y < 6 {
				dy = 6 - y
			} else if y >= h-6 {
				dy = y - (h - 7)
			}
			if dx*dx+dy*dy > 16 {
				continue
			}
			if dx*dx+dy*dy >= 9 || x == 1 || x == w-2 || y == 2 || y == h-3 {
				pixels[y*w+x] = borderColor
			} else {
				pixels[y*w+x] = bgColor
			}
		}
	}

	// Draw crisp bitmap glyphs for "VN" or "EN" (each glyph 8x12, scaled cleanly)
	glyphV := [12]uint8{
		0b11000011,
		0b11000011,
		0b11000011,
		0b11000011,
		0b01100110,
		0b01100110,
		0b01100110,
		0b01100110,
		0b00111100,
		0b00111100,
		0b00011000,
		0b00011000,
	}
	glyphE := [12]uint8{
		0b11111111,
		0b11111111,
		0b11000000,
		0b11000000,
		0b11000000,
		0b11111110,
		0b11111110,
		0b11000000,
		0b11000000,
		0b11000000,
		0b11111111,
		0b11111111,
	}
	glyphN := [12]uint8{
		0b11000011,
		0b11100011,
		0b11110011,
		0b11110011,
		0b11011011,
		0b11011011,
		0b11001111,
		0b11001111,
		0b11000111,
		0b11000111,
		0b11000011,
		0b11000011,
	}

	leftGlyph := glyphE
	if isVN {
		leftGlyph = glyphV
	}

	drawGlyph := func(g [12]uint8, startX, startY int) {
		for row := 0; row < 12; row++ {
			bitsRow := g[row]
			for col := 0; col < 8; col++ {
				if (bitsRow & (1 << (7 - col))) != 0 {
					px := startX + col
					py := startY + row
					if px >= 0 && px < w && py >= 0 && py < h {
						pixels[py*w+px] = 0xFFFFFFFF
					}
				}
			}
		}
	}

	drawGlyph(leftGlyph, 6, 10)
	drawGlyph(glyphN, 17, 10)

	hbmMask, _, _ := procCreateBitmap.Call(w, h, 1, 1, 0)
	ii := ICONINFO{
		FIcon:    1,
		HbmMask:  hbmMask,
		HbmColor: hbmColor,
	}
	hIcon, _, _ := procCreateIconIndirect.Call(uintptr(unsafe.Pointer(&ii)))
	procDeleteObject.Call(hbmColor)
	procDeleteObject.Call(hbmMask)
	return hIcon
}

func isStartupEnabled() bool {
	var hKey uintptr
	subKey := u16(`Software\Microsoft\Windows\CurrentVersion\Run`)
	ret, _, _ := procRegOpenKeyExW.Call(HKEY_CURRENT_USER, uintptr(unsafe.Pointer(subKey)), 0, KEY_READ, uintptr(unsafe.Pointer(&hKey)))
	if ret != 0 {
		return false
	}
	defer procRegCloseKey.Call(hKey)
	valName := u16("VuiKey")
	ret, _, _ = procRegQueryValueExW.Call(hKey, uintptr(unsafe.Pointer(valName)), 0, 0, 0, 0)
	return ret == 0
}

func setStartupEnabled(enable bool) {
	var hKey uintptr
	subKey := u16(`Software\Microsoft\Windows\CurrentVersion\Run`)
	ret, _, _ := procRegOpenKeyExW.Call(HKEY_CURRENT_USER, uintptr(unsafe.Pointer(subKey)), 0, KEY_WRITE, uintptr(unsafe.Pointer(&hKey)))
	if ret != 0 {
		return
	}
	defer procRegCloseKey.Call(hKey)
	valName := u16("VuiKey")
	if enable {
		exePath, err := os.Executable()
		if err != nil {
			return
		}
		cmdStr := fmt.Sprintf(`"%s" --minimized`, exePath)
		u16Slice := utf16.Encode([]rune(cmdStr + "\x00"))
		procRegSetValueExW.Call(
			hKey,
			uintptr(unsafe.Pointer(valName)),
			0,
			REG_SZ,
			uintptr(unsafe.Pointer(&u16Slice[0])),
			uintptr(len(u16Slice)*2),
		)
	} else {
		procRegDeleteValueW.Call(hKey, uintptr(unsafe.Pointer(valName)))
	}
}

func openURL(url string) {
	procShellExecuteW.Call(0, uintptr(unsafe.Pointer(u16("open"))), uintptr(unsafe.Pointer(u16(url))), 0, 0, SW_SHOW)
}

func updateTrayAndUI() {
	if eng.Cfg.Enabled {
		nid.HIcon = hIconVN
		copyTip(&nid.SzTip, "VuiKey [VN] — Đang bật Tiếng Việt ("+eng.Cfg.Mode.String()+") | Ctrl+Space để đổi")
		if hBtnToggle != 0 {
			procSetWindowTextW.Call(hBtnToggle, uintptr(unsafe.Pointer(u16("🟢 ĐANG BẬT TIẾNG VIỆT [VN]  —  Nhấn vào đây hoặc Ctrl+Space để tạm tắt"))))
		}
	} else {
		nid.HIcon = hIconEN
		copyTip(&nid.SzTip, "VuiKey [EN] — Chế độ Tiếng Anh | Ctrl+Space để bật Tiếng Việt")
		if hBtnToggle != 0 {
			procSetWindowTextW.Call(hBtnToggle, uintptr(unsafe.Pointer(u16("⚪ ĐANG TẠM TẮT TIẾNG VIỆT [EN]  —  Nhấn vào đây hoặc Ctrl+Space để bật"))))
		}
	}
	procShellNotifyIconW.Call(NIM_MODIFY, uintptr(unsafe.Pointer(&nid)))

	if hCmbMode != 0 {
		procSendMessageW.Call(hCmbMode, CB_SETCURSEL, uintptr(eng.Cfg.Mode), 0)
	}
	if hChkSplit != 0 {
		chk := uintptr(BST_UNCHECKED)
		if eng.Cfg.SyllableSplit {
			chk = BST_CHECKED
		}
		procSendMessageW.Call(hChkSplit, BM_SETCHECK, chk, 0)
	}
	if hChkSmartEng != 0 {
		chk := uintptr(BST_UNCHECKED)
		if eng.Cfg.SmartEnglish {
			chk = BST_CHECKED
		}
		procSendMessageW.Call(hChkSmartEng, BM_SETCHECK, chk, 0)
	}
	if hChkStartup != 0 {
		chk := uintptr(BST_UNCHECKED)
		if isStartupEnabled() {
			chk = BST_CHECKED
		}
		procSendMessageW.Call(hChkStartup, BM_SETCHECK, chk, 0)
	}
	updateStatsLabel()
}

func updateStatsLabel() {
	if hLblStats != 0 {
		txt := fmt.Sprintf("📊 Phiên hiện tại: %d phím đã gõ  •  %d dấu Tiếng Việt đã biến đổi mượt mà (Zero-Preedit 0.0ms)", totalKeys, transformedCnt)
		procSetWindowTextW.Call(hLblStats, uintptr(unsafe.Pointer(u16(txt))))
	}
}

func updateQuoteUI() {
	s := config.LyricSamples[sampleIdx%len(config.LyricSamples)]
	if hLblQuote != 0 {
		procSetWindowTextW.Call(hLblQuote, uintptr(unsafe.Pointer(u16(fmt.Sprintf("%s: “%s”", s.Category, s.Quote)))))
	}
	if hLblHint != 0 {
		procSetWindowTextW.Call(hLblHint, uintptr(unsafe.Pointer(u16(s.TelexHint))))
	}
}

func copyTip(dst *[128]uint16, s string) {
	for i := range dst {
		dst[i] = 0
	}
	enc := utf16.Encode([]rune(s))
	if len(enc) > 127 {
		enc = enc[:127]
	}
	copy(dst[:], enc)
}

func sendUnicodeReplacement(backspaces int, text string) {
	runes := []rune(text)
	u16Units := utf16.Encode(runes)
	totalEvents := backspaces*2 + len(u16Units)*2
	if totalEvents == 0 {
		return
	}

	inputs := make([]INPUT, 0, totalEvents)
	for i := 0; i < backspaces; i++ {
		inputs = append(inputs,
			INPUT{
				Type: INPUT_KEYBOARD,
				Ki: KEYBDINPUT{
					WVk:         VK_BACK,
					DwFlags:     0,
					DwExtraInfo: VUIKEY_EXTRA_INFO,
				},
			},
			INPUT{
				Type: INPUT_KEYBOARD,
				Ki: KEYBDINPUT{
					WVk:         VK_BACK,
					DwFlags:     KEYEVENTF_KEYUP,
					DwExtraInfo: VUIKEY_EXTRA_INFO,
				},
			},
		)
	}

	for _, codeUnit := range u16Units {
		inputs = append(inputs,
			INPUT{
				Type: INPUT_KEYBOARD,
				Ki: KEYBDINPUT{
					WVk:         0,
					WScan:       codeUnit,
					DwFlags:     KEYEVENTF_UNICODE,
					DwExtraInfo: VUIKEY_EXTRA_INFO,
				},
			},
			INPUT{
				Type: INPUT_KEYBOARD,
				Ki: KEYBDINPUT{
					WVk:         0,
					WScan:       codeUnit,
					DwFlags:     KEYEVENTF_UNICODE | KEYEVENTF_KEYUP,
					DwExtraInfo: VUIKEY_EXTRA_INFO,
				},
			},
		)
	}

	procSendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(INPUT{}),
	)
}

func isKeyDown(vk int) bool {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(vk))
	return (uint16(ret) & 0x8000) != 0
}

func isCapsLockOn() bool {
	ret, _, _ := procGetKeyState.Call(VK_CAPITAL)
	return (uint16(ret) & 0x0001) != 0
}

func vkToASCII(vk uint32, shift, caps bool) (rune, bool) {
	if vk >= 'A' && vk <= 'Z' {
		upper := shift != caps
		if upper {
			return rune(vk), true
		}
		return rune(vk + 32), true
	}
	if vk >= '0' && vk <= '9' {
		if !shift {
			return rune(vk), true
		}
		shiftedDigits := ")!@#$%^&*("
		return rune(shiftedDigits[vk-'0']), true
	}
	switch vk {
	case 0xBA: // OEM_1 (;:)
		if shift {
			return ':', true
		}
		return ';', true
	case 0xBB: // OEM_PLUS (=+)
		if shift {
			return '+', true
		}
		return '=', true
	case 0xBC: // OEM_COMMA (,<)
		if shift {
			return '<', true
		}
		return ',', true
	case 0xBD: // OEM_MINUS (-_)
		if shift {
			return '_', true
		}
		return '-', true
	case 0xBE: // OEM_PERIOD (.>)
		if shift {
			return '>', true
		}
		return '.', true
	case 0xBF: // OEM_2 (/?)
		if shift {
			return '?', true
		}
		return '/', true
	case 0xC0: // OEM_3 (`~)
		if shift {
			return '~', true
		}
		return '`', true
	case 0xDB: // OEM_4 ([{)
		if shift {
			return '{', true
		}
		return '[', true
	case 0xDC: // OEM_5 (\|)
		if shift {
			return '|', true
		}
		return '\\', true
	case 0xDD: // OEM_6 (]})
		if shift {
			return '}', true
		}
		return ']', true
	case 0xDE: // OEM_7 ('")
		if shift {
			return '"', true
		}
		return '\'', true
	}
	return 0, false
}

func lowLevelKeyboardProc(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 && (wParam == WM_KEYDOWN || wParam == WM_SYSKEYDOWN) {
		kb := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))

		// Ignore our own synthetic SendInput events so they never feed back into the state machine!
		if (kb.Flags&LLKHF_INJECTED) != 0 || kb.DwExtraInfo == VUIKEY_EXTRA_INFO {
			ret, _, _ := procCallNextHookEx.Call(hKeyHook, uintptr(nCode), wParam, lParam)
			return ret
		}

		ctrl := isKeyDown(VK_CONTROL)
		alt := isKeyDown(VK_MENU)
		shift := isKeyDown(VK_SHIFT)
		win := isKeyDown(VK_LWIN) || isKeyDown(VK_RWIN)

		// Toggle [VN] / [EN] on Ctrl+Space or Alt+Z
		if (ctrl && !alt && !shift && kb.VkCode == VK_SPACE) || (alt && !ctrl && !shift && kb.VkCode == 'Z') {
			eng.Cfg.Enabled = !eng.Cfg.Enabled
			eng.Reset()
			_ = config.Save(eng.Cfg)
			procPostMessageW.Call(hMainWnd, WM_SYNC_UI, 0, 0)
			return 1
		}

		// Any Ctrl/Alt/Win shortcut resets the active syllable and passes straight through
		if ctrl || alt || win {
			eng.Reset()
			ret, _, _ := procCallNextHookEx.Call(hKeyHook, uintptr(nCode), wParam, lParam)
			return ret
		}

		switch kb.VkCode {
		case VK_BACK:
			eng.HandleBackspace()
			ret, _, _ := procCallNextHookEx.Call(hKeyHook, uintptr(nCode), wParam, lParam)
			return ret

		case VK_DELETE, VK_RETURN, VK_TAB, VK_ESCAPE,
			VK_LEFT, VK_RIGHT, VK_UP, VK_DOWN, VK_HOME, VK_END, VK_PRIOR, VK_NEXT:
			eng.HandleDelete()
			ret, _, _ := procCallNextHookEx.Call(hKeyHook, uintptr(nCode), wParam, lParam)
			return ret

		case VK_SPACE:
			if act, ok := eng.CheckMacroOnSpace(); ok {
				sendUnicodeReplacement(act.Backspaces, act.CommitText+" ")
				return 1
			}
			eng.Reset()
			ret, _, _ := procCallNextHookEx.Call(hKeyHook, uintptr(nCode), wParam, lParam)
			return ret
		}

		if ch, ok := vkToASCII(kb.VkCode, shift, isCapsLockOn()); ok {
			totalKeys++
			act := eng.ProcessKey(ch)
			if act.Swallow {
				transformedCnt++
				sendUnicodeReplacement(act.Backspaces, act.CommitText)
				procPostMessageW.Call(hMainWnd, WM_SYNC_UI, 1, 0)
				return 1
			}
		}
	}
	ret, _, _ := procCallNextHookEx.Call(hKeyHook, uintptr(nCode), wParam, lParam)
	return ret
}

func lowLevelMouseProc(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		if wParam == WM_LBUTTONDOWN || wParam == WM_RBUTTONDOWN || wParam == WM_MBUTTONDOWN {
			eng.Reset()
		}
	}
	ret, _, _ := procCallNextHookEx.Call(hMouseHook, uintptr(nCode), wParam, lParam)
	return ret
}

func showTrayPopupMenu(hwnd uintptr) {
	hMenu, _, _ := procCreatePopupMenu.Call()
	if hMenu == 0 {
		return
	}
	defer procDestroyMenu.Call(hMenu)

	procAppendMenuW.Call(hMenu, MF_STRING, ID_TRAY_SHOW, uintptr(unsafe.Pointer(u16("⚙️ Mở Bảng Điều Khiển VuiKey Studio..."))))
	procAppendMenuW.Call(hMenu, MF_SEPARATOR, 0, 0)

	telexFlags := uintptr(MF_STRING)
	vniFlags := uintptr(MF_STRING)
	enFlags := uintptr(MF_STRING)
	if !eng.Cfg.Enabled {
		enFlags |= MF_CHECKED
	} else if eng.Cfg.Mode == engine.ModeTelex || eng.Cfg.Mode == engine.ModeTelexVNI {
		telexFlags |= MF_CHECKED
	} else {
		vniFlags |= MF_CHECKED
	}

	procAppendMenuW.Call(hMenu, telexFlags, ID_TRAY_VN_TELEX, uintptr(unsafe.Pointer(u16("🇻🇳 Bật Tiếng Việt — Kiểu gõ Telex (Zero-Preedit)"))))
	procAppendMenuW.Call(hMenu, vniFlags, ID_TRAY_VN_VNI, uintptr(unsafe.Pointer(u16("🇻🇳 Bật Tiếng Việt — Kiểu gõ VNI (Zero-Preedit)"))))
	procAppendMenuW.Call(hMenu, enFlags, ID_TRAY_EN, uintptr(unsafe.Pointer(u16("🇺🇸 Tạm Tắt Tiếng Việt — Chế độ Tiếng Anh (Ctrl+Space)"))))
	procAppendMenuW.Call(hMenu, MF_SEPARATOR, 0, 0)

	startupFlags := uintptr(MF_STRING)
	if isStartupEnabled() {
		startupFlags |= MF_CHECKED
	}
	procAppendMenuW.Call(hMenu, startupFlags, ID_TRAY_STARTUP, uintptr(unsafe.Pointer(u16("🚀 Khởi động cùng Windows"))))
	procAppendMenuW.Call(hMenu, MF_SEPARATOR, 0, 0)
	procAppendMenuW.Call(hMenu, MF_STRING, ID_TRAY_VUICAHAT, uintptr(unsafe.Pointer(u16("🎤 Phòng hát & Thu âm trực tuyến — VuiCaHat.com"))))
	procAppendMenuW.Call(hMenu, MF_STRING, ID_TRAY_ECHIP, uintptr(unsafe.Pointer(u16("🏢 Công ty TNHH eChip — echipco.com"))))
	procAppendMenuW.Call(hMenu, MF_SEPARATOR, 0, 0)
	procAppendMenuW.Call(hMenu, MF_STRING, ID_TRAY_EXIT, uintptr(unsafe.Pointer(u16("❌ Thoát VuiKey"))))

	var pt POINT
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	procSetForegroundWindow.Call(hwnd)
	procTrackPopupMenu.Call(hMenu, TPM_BOTTOMALIGN|TPM_LEFTALIGN, uintptr(pt.X), uintptr(pt.Y), 0, hwnd, 0)
}

func wndProc(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_CTLCOLORSTATIC:
		hdc := wParam
		procSetBkMode.Call(hdc, 1) // TRANSPARENT
		procSetTextColor.Call(hdc, 0x00F8FAFC)
		return hBgBrush

	case WM_CTLCOLOREDIT:
		hdc := wParam
		procSetBkColor.Call(hdc, 0x002B1E09)
		procSetTextColor.Call(hdc, 0x00F8BD38)
		return hEditBrush

	case WM_SYNC_UI:
		if wParam == 1 {
			updateStatsLabel()
		} else {
			updateTrayAndUI()
		}
		return 0

	case WM_TRAYICON:
		switch uint32(lParam) {
		case WM_LBUTTONUP:
			// Single left-click on tray icon toggles [VN] / [EN] immediately (just like UniKey!)
			eng.Cfg.Enabled = !eng.Cfg.Enabled
			eng.Reset()
			_ = config.Save(eng.Cfg)
			updateTrayAndUI()
		case WM_LBUTTONDBLCLK:
			procShowWindow.Call(hwnd, SW_RESTORE)
			procSetForegroundWindow.Call(hwnd)
		case WM_RBUTTONUP:
			showTrayPopupMenu(hwnd)
		}
		return 0

	case WM_COMMAND:
		id := int(wParam & 0xFFFF)
		notif := int((wParam >> 16) & 0xFFFF)

		switch id {
		case ID_BTN_TOGGLE:
			eng.Cfg.Enabled = !eng.Cfg.Enabled
			eng.Reset()
			_ = config.Save(eng.Cfg)
			updateTrayAndUI()

		case ID_CMB_MODE:
			if notif == CBN_SELCHANGE {
				idx, _, _ := procSendMessageW.Call(hCmbMode, CB_GETCURSEL, 0, 0)
				eng.Cfg.Mode = engine.InputMode(idx)
				eng.Reset()
				_ = config.Save(eng.Cfg)
				updateTrayAndUI()
			}

		case ID_CHK_SPLIT:
			st, _, _ := procSendMessageW.Call(hChkSplit, BM_GETCHECK, 0, 0)
			eng.Cfg.SyllableSplit = (st == BST_CHECKED)
			_ = config.Save(eng.Cfg)

		case ID_CHK_SMART_ENG:
			st, _, _ := procSendMessageW.Call(hChkSmartEng, BM_GETCHECK, 0, 0)
			eng.Cfg.SmartEnglish = (st == BST_CHECKED)
			_ = config.Save(eng.Cfg)

		case ID_CHK_STARTUP:
			st, _, _ := procSendMessageW.Call(hChkStartup, BM_GETCHECK, 0, 0)
			setStartupEnabled(st == BST_CHECKED)

		case ID_BTN_CLEAR_TEST:
			procSetWindowTextW.Call(hEditTest, uintptr(unsafe.Pointer(u16(""))))

		case ID_BTN_NEXT_QUOTE:
			sampleIdx = (sampleIdx + 1) % len(config.LyricSamples)
			updateQuoteUI()

		case ID_BTN_VUICAHAT, ID_TRAY_VUICAHAT:
			openURL("https://vuicahat.com")

		case ID_BTN_ECHIP, ID_TRAY_ECHIP:
			openURL("https://echipco.com")

		case ID_BTN_HIDE:
			procShowWindow.Call(hwnd, SW_HIDE)

		case ID_TRAY_SHOW:
			procShowWindow.Call(hwnd, SW_RESTORE)
			procSetForegroundWindow.Call(hwnd)

		case ID_TRAY_VN_TELEX:
			eng.Cfg.Enabled = true
			eng.Cfg.Mode = engine.ModeTelex
			eng.Reset()
			_ = config.Save(eng.Cfg)
			updateTrayAndUI()

		case ID_TRAY_VN_VNI:
			eng.Cfg.Enabled = true
			eng.Cfg.Mode = engine.ModeVNI
			eng.Reset()
			_ = config.Save(eng.Cfg)
			updateTrayAndUI()

		case ID_TRAY_EN:
			eng.Cfg.Enabled = false
			eng.Reset()
			_ = config.Save(eng.Cfg)
			updateTrayAndUI()

		case ID_TRAY_STARTUP:
			setStartupEnabled(!isStartupEnabled())
			updateTrayAndUI()

		case ID_TRAY_EXIT:
			procShellNotifyIconW.Call(NIM_DELETE, uintptr(unsafe.Pointer(&nid)))
			procPostQuitMessage.Call(0)
		}
		return 0

	case WM_CLOSE:
		// Minimize to system tray instead of quitting when clicking X
		procShowWindow.Call(hwnd, SW_HIDE)
		return 0

	case WM_DESTROY:
		procShellNotifyIconW.Call(NIM_DELETE, uintptr(unsafe.Pointer(&nid)))
		procPostQuitMessage.Call(0)
		return 0
	}

	ret, _, _ := procDefWindowProcW.Call(hwnd, uintptr(msg), wParam, lParam)
	return ret
}

func createControl(className, text string, style uint32, x, y, w, h int, id uintptr, font uintptr) uintptr {
	hwnd, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(u16(className))),
		uintptr(unsafe.Pointer(u16(text))),
		uintptr(style),
		uintptr(x), uintptr(y), uintptr(w), uintptr(h),
		hMainWnd,
		id,
		hInst,
		0,
	)
	if hwnd != 0 && font != 0 {
		procSendMessageW.Call(hwnd, 0x0030, font, 1) // WM_SETFONT
	}
	return hwnd
}

func main() {
	runtime.LockOSThread()

	// Single-instance check via named Win32 mutex
	mutexName := u16("Local\\VuiKey_SingleInstance_Mutex")
	procCreateMutexW.Call(0, 1, uintptr(unsafe.Pointer(mutexName)))
	if errCode, _, _ := procGetLastError.Call(); errCode == 183 { // ERROR_ALREADY_EXISTS
		return
	}

	cfg := config.Load()
	cfg.Enabled = true
	_ = config.Save(cfg)
	eng = engine.NewEngine(cfg)

	hInst, _, _ = procGetModuleHandleW.Call(0)
	hIconVN = createBadgeIcon(true)
	hIconEN = createBadgeIcon(false)
	hBgBrush, _, _ = procCreateSolidBrush.Call(0x002A170F)   // #0F172A Slate dark
	hCardBrush, _, _ = procCreateSolidBrush.Call(0x003B291E) // #1E293B
	hEditBrush, _, _ = procCreateSolidBrush.Call(0x002B1E09) // #091E2B

	className := u16("VuiKeyStudioWndClass")
	wc := WNDCLASSEXW{
		CbSize:        uint32(unsafe.Sizeof(WNDCLASSEXW{})),
		LpfnWndProc:   syscall.NewCallback(wndProc),
		HInstance:     hInst,
		HIcon:         hIconVN,
		HbrBackground: hBgBrush,
		LpszClassName: className,
		HIconSm:       hIconVN,
	}
	procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

	hMainWnd, _, _ = procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(u16("VuiKey Studio v1.0.0 — Bộ Gõ Tiếng Việt Thế Hệ Mới (Windows & Linux)"))),
		WS_OVERLAPPED|WS_CAPTION|WS_SYSMENU|WS_MINIMIZEBOX,
		220, 140, 660, 520,
		0, 0, hInst, 0,
	)

	// Enable Windows 10/11 Dark Titlebar (DWMWA_USE_IMMERSIVE_DARK_MODE = 20)
	darkVal := int32(1)
	procDwmSetWindowAttribute.Call(hMainWnd, 20, uintptr(unsafe.Pointer(&darkVal)), 4)

	hFontTitle, _, _ := procCreateFontW.Call(22, 0, 0, 0, 700, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(u16("Segoe UI"))))
	hFontBold, _, _ := procCreateFontW.Call(16, 0, 0, 0, 700, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(u16("Segoe UI"))))
	hFontNorm, _, _ := procCreateFontW.Call(15, 0, 0, 0, 400, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(u16("Segoe UI"))))

	createControl("STATIC", "⚡ VuiKey Studio — Bộ Gõ Tiếng Việt Zero-Preedit (0.0ms)", WS_VISIBLE|WS_CHILD, 24, 18, 600, 28, 0, hFontTitle)
	createControl("STATIC", "Kiến trúc Sliding Syllable Window • Gõ tự do chữ đ/ươ • Phát triển bởi Công ty TNHH eChip & VuiCaHat.com", WS_VISIBLE|WS_CHILD, 24, 46, 600, 22, 0, hFontNorm)

	hBtnToggle = createControl("BUTTON", "", WS_VISIBLE|WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON, 24, 78, 595, 42, ID_BTN_TOGGLE, hFontBold)

	createControl("STATIC", "⌨️ Kiểu gõ Tiếng Việt:", WS_VISIBLE|WS_CHILD, 24, 136, 175, 24, 0, hFontBold)
	hCmbMode = createControl("COMBOBOX", "", WS_VISIBLE|WS_CHILD|WS_TABSTOP|CBS_DROPDOWNLIST, 205, 132, 414, 120, ID_CMB_MODE, hFontNorm)
	procSendMessageW.Call(hCmbMode, CB_ADDSTRING, 0, uintptr(unsafe.Pointer(u16("Telex (Mặc định — Hỗ trợ gõ tự do dodongj, dduowcj, phuwcs...)"))))
	procSendMessageW.Call(hCmbMode, CB_ADDSTRING, 0, uintptr(unsafe.Pointer(u16("VNI (Gõ dấu bằng phím số 1..9)"))))
	procSendMessageW.Call(hCmbMode, CB_ADDSTRING, 0, uintptr(unsafe.Pointer(u16("Tích hợp Telex + VNI song song"))))

	hChkSplit = createControl("BUTTON", "✂️ Tự động tách âm tiết thông minh (Sliding Syllable Window — Chống nhân đôi chữ)", WS_VISIBLE|WS_CHILD|WS_TABSTOP|BS_AUTOCHECKBOX, 24, 172, 595, 26, ID_CHK_SPLIT, hFontNorm)
	hChkSmartEng = createControl("BUTTON", "🛡️ Tự nhận diện từ Tiếng Anh (Smart English Guard — Giữ nguyên window, view, test, server...)", WS_VISIBLE|WS_CHILD|WS_TABSTOP|BS_AUTOCHECKBOX, 24, 202, 595, 26, ID_CHK_SMART_ENG, hFontNorm)
	hChkStartup = createControl("BUTTON", "🚀 Tự khởi động VuiKey cùng Windows (Chạy ngầm trên khay hệ thống)", WS_VISIBLE|WS_CHILD|WS_TABSTOP|BS_AUTOCHECKBOX, 24, 232, 595, 26, ID_CHK_STARTUP, hFontNorm)

	createControl("STATIC", "🎵 Góc Thử Phím & Cảm Hứng Âm Nhạc (VuiCaHat.com):", WS_VISIBLE|WS_CHILD, 24, 272, 450, 22, 0, hFontBold)
	hLblQuote = createControl("STATIC", "", WS_VISIBLE|WS_CHILD, 24, 296, 595, 22, 0, hFontBold)
	hLblHint = createControl("STATIC", "", WS_VISIBLE|WS_CHILD, 24, 318, 595, 22, 0, hFontNorm)

	hEditTest = createControl("EDIT", "", WS_VISIBLE|WS_CHILD|WS_BORDER|WS_TABSTOP|ES_AUTOHSCROLL, 24, 346, 595, 32, 0, hFontBold)

	createControl("BUTTON", "🧹 Xóa ô thử", WS_VISIBLE|WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON, 24, 386, 120, 30, ID_BTN_CLEAR_TEST, hFontNorm)
	createControl("BUTTON", "🎲 Câu hát mẫu khác", WS_VISIBLE|WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON, 154, 386, 165, 30, ID_BTN_NEXT_QUOTE, hFontNorm)
	createControl("BUTTON", "🎤 Mở VuiCaHat.com", WS_VISIBLE|WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON, 329, 386, 145, 30, ID_BTN_VUICAHAT, hFontNorm)
	createControl("BUTTON", "🏢 Mở echipco.com", WS_VISIBLE|WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON, 484, 386, 135, 30, ID_BTN_ECHIP, hFontNorm)

	hLblStats = createControl("STATIC", "", WS_VISIBLE|WS_CHILD, 24, 432, 470, 24, 0, hFontNorm)
	createControl("BUTTON", "🔽 Thu xuống Khay", WS_VISIBLE|WS_CHILD|WS_TABSTOP|BS_PUSHBUTTON, 500, 426, 119, 30, ID_BTN_HIDE, hFontNorm)

	// Register System Tray Icon
	nid = NOTIFYICONDATAW{
		CbSize:           uint32(unsafe.Sizeof(NOTIFYICONDATAW{})),
		HWnd:             hMainWnd,
		UID:              1,
		UFlags:           NIF_MESSAGE | NIF_ICON | NIF_TIP,
		UCallbackMessage: WM_TRAYICON,
		HIcon:            hIconVN,
	}
	copyTip(&nid.SzTip, "VuiKey [VN] — Bộ Gõ Tiếng Việt Thế Hệ Mới")
	procShellNotifyIconW.Call(NIM_ADD, uintptr(unsafe.Pointer(&nid)))

	updateQuoteUI()
	updateTrayAndUI()

	// Install Win32 Low-Level Keyboard & Mouse Hooks
	hKeyHook, _, _ = procSetWindowsHookExW.Call(WH_KEYBOARD_LL, syscall.NewCallback(lowLevelKeyboardProc), hInst, 0)
	hMouseHook, _, _ = procSetWindowsHookExW.Call(WH_MOUSE_LL, syscall.NewCallback(lowLevelMouseProc), hInst, 0)
	defer func() {
		if hKeyHook != 0 {
			procUnhookWindowsHookEx.Call(hKeyHook)
		}
		if hMouseHook != 0 {
			procUnhookWindowsHookEx.Call(hMouseHook)
		}
	}()

	minimized := false
	for _, arg := range os.Args[1:] {
		if arg == "--minimized" || arg == "-m" || arg == "daemon" {
			minimized = true
		}
	}
	if !minimized {
		procShowWindow.Call(hMainWnd, SW_SHOW)
		procUpdateWindow.Call(hMainWnd)
	}

	var m MSG
	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(ret) <= 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}
