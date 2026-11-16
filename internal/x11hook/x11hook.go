package x11hook

/*
#cgo pkg-config: x11 xtst xi
#include <X11/Xlib.h>
#include <X11/XKBlib.h>
#include <X11/keysym.h>
#include <X11/extensions/XTest.h>
#include <X11/extensions/XInput2.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

static int probe_err = 0;
static int x11_error_handler(Display *dpy, XErrorEvent *err) {
    probe_err = 1;
    return 0;
}

static void install_error_handler() {
    XSetErrorHandler(x11_error_handler);
}

static unsigned int tracked_mods = 0;
static int synth_shift_depth = 0;

static int setup_xi2(Display *dpy, Window root) {
    int event, error;
    int xi_opcode = -1;
    if (!XQueryExtension(dpy, "XInputExtension", &xi_opcode, &event, &error)) {
        return -1;
    }
    int major = 2, minor = 0;
    if (XIQueryVersion(dpy, &major, &minor) != Success) {
        return -1;
    }
    XkbStateRec xkb_state;
    if (XkbGetState(dpy, XkbUseCoreKbd, &xkb_state) == Success) {
        tracked_mods = xkb_state.mods;
    }
    XIEventMask evmask;
    unsigned char mask[(XI_LASTEVENT + 7)/8];
    memset(mask, 0, sizeof(mask));
    XISetMask(mask, XI_RawKeyPress);
    XISetMask(mask, XI_RawKeyRelease);
    XISetMask(mask, XI_RawButtonPress);

    evmask.deviceid = XIAllMasterDevices;
    evmask.mask_len = sizeof(mask);
    evmask.mask = mask;
    XISelectEvents(dpy, root, &evmask, 1);
    XFlush(dpy);
    return xi_opcode;
}

static int is_excluded_hardware_kc(int kc) {
    // Standard physical keyboard range (9..119), except 93, 97, 103 which are unused on standard PC keyboards
    if (kc >= 9 && kc <= 119 && kc != 93 && kc != 97 && kc != 103) {
        return 1;
    }
    // Exclude Volume, Power, Modifiers, Media Play/Next/Prev, PrintScreen, Brightness, Touchpad
    switch (kc) {
    case 121: case 122: case 123: case 124: case 125: case 126: case 127: case 129:
    case 133: case 134: case 135:
    case 171: case 172: case 173: case 174:
    case 198: case 199: case 200: case 201:
    case 203: case 204: case 205: case 206: case 207:
    case 218:
    case 232: case 233: case 234: case 235: case 236: case 237: case 238:
        return 1;
    default:
        return 0;
    }
}

// Dynamically discovers safe keycodes across ANY Linux distro / Window Manager.
// Uses high non-physical keycodes (128..255) that have existing XKB entries so Xorg assigns TWO_LEVEL (Shift-aware) key types!
static int discover_safe_keycodes(Display *dpy, Window root, int *out_kcs, int max_need) {
    int count = 0;
    for (int kc = 128; kc <= 255 && count < max_need; kc++) {
        if (is_excluded_hardware_kc(kc)) continue;
        out_kcs[count++] = kc;
    }
    return count;
}

typedef struct {
    int kind;       // 1 = RawKeyPress, 2 = RawButtonPress, 3 = GrabbedKeyPress (Hotkeys only)
    int detail;     // keycode or button
    unsigned int state; // modifier mask
    unsigned long keysym;
} VuiXEvent;

static int next_vui_event(Display *dpy, int xi_opcode, VuiXEvent *out) {
    XEvent ev;
    XNextEvent(dpy, &ev);

    if (ev.type == KeyPress) {
        out->kind = 3;
        out->detail = ev.xkey.keycode;
        out->state = ev.xkey.state;
        int shift = ( (ev.xkey.state & ShiftMask) != 0 ) ^ ( (ev.xkey.state & LockMask) != 0 );
        out->keysym = XkbKeycodeToKeysym(dpy, ev.xkey.keycode, 0, shift ? 1 : 0);
        return 1;
    }

    if (ev.type == GenericEvent && ev.xcookie.extension == xi_opcode) {
        if (XGetEventData(dpy, &ev.xcookie)) {
            int evtype = ev.xcookie.evtype;
            XIRawEvent *raw = (XIRawEvent *)ev.xcookie.data;
            if (evtype == XI_RawKeyPress || evtype == XI_RawKeyRelease) {
                int kc = raw->detail;
                KeySym base_sym = XkbKeycodeToKeysym(dpy, kc, 0, 0);
                unsigned int mod_bit = 0;
                if (base_sym == XK_Shift_L || base_sym == XK_Shift_R) mod_bit = ShiftMask;
                else if (base_sym == XK_Control_L || base_sym == XK_Control_R) mod_bit = ControlMask;
                else if (base_sym == XK_Alt_L || base_sym == XK_Alt_R || base_sym == XK_Meta_L || base_sym == XK_Meta_R) mod_bit = Mod1Mask;
                else if (base_sym == XK_Super_L || base_sym == XK_Super_R) mod_bit = Mod4Mask;

                if (mod_bit != 0) {
                    if (evtype == XI_RawKeyPress) tracked_mods |= mod_bit;
                    else tracked_mods &= ~mod_bit;
                } else if (base_sym == XK_Caps_Lock && evtype == XI_RawKeyPress) {
                    tracked_mods ^= LockMask;
                }

                if (evtype == XI_RawKeyPress) {
                    out->kind = 1;
                    out->detail = kc;
                    out->state = tracked_mods;
                    int shift = ( (tracked_mods & ShiftMask) != 0 ) ^ ( (tracked_mods & LockMask) != 0 );
                    out->keysym = XkbKeycodeToKeysym(dpy, kc, 0, shift ? 1 : 0);
                    XFreeEventData(dpy, &ev.xcookie);
                    return 1;
                }
                XFreeEventData(dpy, &ev.xcookie);
                return 0;
            } else if (evtype == XI_RawButtonPress) {
                out->kind = 2;
                out->detail = raw->detail;
                out->state = 0;
                out->keysym = 0;
                XFreeEventData(dpy, &ev.xcookie);
                return 1;
            }
            XFreeEventData(dpy, &ev.xcookie);
        }
    }
    return 0;
}

static void grab_keycode(Display *dpy, Window root, int keycode, unsigned int base_mods) {
    unsigned int locks[4] = {0, LockMask, Mod2Mask, LockMask | Mod2Mask};
    for (int i = 0; i < 4; i++) {
        XGrabKey(dpy, keycode, base_mods | locks[i], root, False, GrabModeAsync, GrabModeAsync);
    }
}

static KeySym unicode_to_keysym(unsigned int cp) {
    switch (cp) {
    case 0x0111: return 0x01f0; // XK_dstroke (đ)
    case 0x0110: return 0x01d0; // XK_Dstroke (Đ)
    case 0x0103: return 0x01e3; // XK_abreve (ă)
    case 0x0102: return 0x01c3; // XK_Abreve (Ă)
    case 0x0129: return 0x03b5; // XK_itilde (ĩ)
    case 0x0128: return 0x03a5; // XK_Itilde (Ĩ)
    case 0x0169: return 0x03fd; // XK_utilde (ũ)
    case 0x0168: return 0x03dd; // XK_Utilde (Ũ)
    }
    return (cp < 0x100) ? (KeySym)cp : (KeySym)(0x01000000 | cp);
}

static void map_keycode_pair(Display *dpy, int kc, unsigned int lower_cp, unsigned int upper_cp) {
    KeySym lsym = unicode_to_keysym(lower_cp);
    KeySym usym = unicode_to_keysym(upper_cp);
    KeySym syms[4] = {lsym, usym, lsym, usym};
    XChangeKeyboardMapping(dpy, kc, 4, syms, 1);
}

static void release_active_grab(Display *dpy) {
    XUngrabKeyboard(dpy, CurrentTime);
}

static void tap_keycode(Display *dpy, int kc) {
    XTestFakeKeyEvent(dpy, kc, True, CurrentTime);
    XTestFakeKeyEvent(dpy, kc, False, CurrentTime);
}

static int is_shift_currently_active(Display *dpy) {
    XkbStateRec xkb_state;
    XkbGetState(dpy, XkbUseCoreKbd, &xkb_state);
    int shift = ((xkb_state.mods & ShiftMask) != 0) ^ ((xkb_state.mods & LockMask) != 0);
    return shift;
}

static void tap_with_shift_target(Display *dpy, int kc, int want_upper, int shift_l_kc, int shift_r_kc) {
    int phys_shift = (tracked_mods & ShiftMask) != 0;
    int caps_lock = (tracked_mods & LockMask) != 0;
    int effective_upper = phys_shift ^ caps_lock;

    if (want_upper == effective_upper) {
        tap_keycode(dpy, kc);
        return;
    }

    if (want_upper && !effective_upper) {
        XTestFakeKeyEvent(dpy, shift_l_kc, True, CurrentTime);
        tap_keycode(dpy, kc);
        XTestFakeKeyEvent(dpy, shift_l_kc, False, CurrentTime);
    } else {
        if (phys_shift) {
            XTestFakeKeyEvent(dpy, shift_l_kc, False, CurrentTime);
            XTestFakeKeyEvent(dpy, shift_r_kc, False, CurrentTime);
            tap_keycode(dpy, kc);
        } else {
            XTestFakeKeyEvent(dpy, shift_l_kc, True, CurrentTime);
            tap_keycode(dpy, kc);
            XTestFakeKeyEvent(dpy, shift_l_kc, False, CurrentTime);
        }
    }
}
*/
import "C"

import (
	"sync"
	"unicode"

	"github.com/mrtrinhvn/VuiKey/internal/engine"
)

type runePair struct {
	lower rune
	upper rune
}

var allSyntheticPairs = []runePair{
	// 67 Vietnamese accented pairs
	{'đ', 'Đ'},
	{'á', 'Á'}, {'à', 'À'}, {'ả', 'Ả'}, {'ã', 'Ã'}, {'ạ', 'Ạ'},
	{'ă', 'Ă'}, {'ắ', 'Ắ'}, {'ằ', 'Ằ'}, {'ẳ', 'Ẳ'}, {'ẵ', 'Ẵ'}, {'ặ', 'Ặ'},
	{'â', 'Â'}, {'ấ', 'Ấ'}, {'ầ', 'Ầ'}, {'ẩ', 'Ẩ'}, {'ẫ', 'Ẫ'}, {'ậ', 'Ậ'},
	{'é', 'É'}, {'è', 'È'}, {'ẻ', 'Ẻ'}, {'ẽ', 'Ẽ'}, {'ẹ', 'Ẹ'},
	{'ê', 'Ê'}, {'ế', 'Ế'}, {'ề', 'Ề'}, {'ể', 'Ể'}, {'ễ', 'Ễ'}, {'ệ', 'Ệ'},
	{'í', 'Í'}, {'ì', 'Ì'}, {'ỉ', 'Ỉ'}, {'ĩ', 'Ĩ'}, {'ị', 'Ị'},
	{'ó', 'Ó'}, {'ò', 'Ò'}, {'ỏ', 'Ỏ'}, {'õ', 'Õ'}, {'ọ', 'Ọ'},
	{'ô', 'Ô'}, {'ố', 'Ố'}, {'ồ', 'Ồ'}, {'ổ', 'Ổ'}, {'ỗ', 'Ỗ'}, {'ộ', 'Ộ'},
	{'ơ', 'Ơ'}, {'ớ', 'Ớ'}, {'ờ', 'Ờ'}, {'ở', 'Ở'}, {'ỡ', 'Ỡ'}, {'ợ', 'Ợ'},
	{'ú', 'Ú'}, {'ù', 'Ù'}, {'ủ', 'Ủ'}, {'ũ', 'Ũ'}, {'ụ', 'Ụ'},
	{'ư', 'Ư'}, {'ứ', 'Ứ'}, {'ừ', 'Ừ'}, {'ử', 'Ử'}, {'ữ', 'Ữ'}, {'ự', 'Ự'},
	{'ý', 'Ý'}, {'ỳ', 'Ỳ'}, {'ỷ', 'Ỷ'}, {'ỹ', 'Ỹ'}, {'ỵ', 'Ỵ'},
	// 26 Synthetic ASCII letter pairs (isolated from physical keyboard keycodes 10..60!)
	{'a', 'A'}, {'b', 'B'}, {'c', 'C'}, {'d', 'D'}, {'e', 'E'}, {'f', 'F'},
	{'g', 'G'}, {'h', 'H'}, {'i', 'I'}, {'j', 'J'}, {'k', 'K'}, {'l', 'L'},
	{'m', 'M'}, {'n', 'N'}, {'o', 'O'}, {'p', 'P'}, {'q', 'Q'}, {'r', 'R'},
	{'s', 'S'}, {'t', 'T'}, {'u', 'U'}, {'v', 'V'}, {'w', 'W'}, {'x', 'X'},
	{'y', 'Y'}, {'z', 'Z'},
}

type keyMapping struct {
	kc    int
	shift bool
}

type Daemon struct {
	mu             sync.Mutex
	Eng            *engine.Engine
	OnStateChanged func(enabled bool, mode engine.InputMode)
	OnStatsUpdated func(totalKeystrokes, vietnameseTransforms int64)

	dpy                   *C.Display
	root                  C.Window
	xiOpcode              C.int
	backspaceKC           int
	shiftLKC              int
	shiftRKC              int
	spaceKC               int
	tildeKC               int
	syntheticKCSet        map[int]bool
	runeToKey             map[rune]keyMapping
	pendingFakeBackspaces int

	TotalKeystrokes      int64
	VietnameseTransforms int64
}

func NewDaemon(eng *engine.Engine) *Daemon {
	return &Daemon{
		Eng:            eng,
		syntheticKCSet: make(map[int]bool),
		runeToKey:      make(map[rune]keyMapping),
	}
}

func (d *Daemon) initStaticKeymap() {
	for ch := rune(0x20); ch <= 0x7e; ch++ {
		if unicode.IsLetter(ch) {
			continue
		}
		kc := int(C.XKeysymToKeycode(d.dpy, C.KeySym(ch)))
		if kc != 0 {
			d.runeToKey[ch] = keyMapping{kc: kc, shift: false}
		}
	}

	need := len(allSyntheticPairs)
	cKCs := make([]C.int, need)
	found := int(C.discover_safe_keycodes(d.dpy, d.root, &cKCs[0], C.int(need)))

	for i := 0; i < found && i < need; i++ {
		pair := allSyntheticPairs[i]
		kc := int(cKCs[i])
		d.syntheticKCSet[kc] = true
		C.map_keycode_pair(d.dpy, C.int(kc), C.uint(pair.lower), C.uint(pair.upper))
		d.runeToKey[pair.lower] = keyMapping{kc: kc, shift: false}
		d.runeToKey[pair.upper] = keyMapping{kc: kc, shift: true}
	}
	C.XSync(d.dpy, C.False)
}

func (d *Daemon) Start() error {
	C.install_error_handler()
	d.dpy = C.XOpenDisplay(nil)
	if d.dpy == nil {
		return ErrNoX11Display
	}
	d.root = C.XDefaultRootWindow(d.dpy)
	d.xiOpcode = C.setup_xi2(d.dpy, d.root)
	d.backspaceKC = int(C.XKeysymToKeycode(d.dpy, C.XK_BackSpace))
	d.shiftLKC = int(C.XKeysymToKeycode(d.dpy, C.XK_Shift_L))
	d.shiftRKC = int(C.XKeysymToKeycode(d.dpy, C.XK_Shift_R))
	d.spaceKC = int(C.XKeysymToKeycode(d.dpy, C.XK_space))
	d.tildeKC = int(C.XKeysymToKeycode(d.dpy, C.XK_grave))

	d.initStaticKeymap()

	if d.spaceKC != 0 {
		C.grab_keycode(d.dpy, d.root, C.int(d.spaceKC), C.ControlMask)
	}
	if d.tildeKC != 0 {
		C.grab_keycode(d.dpy, d.root, C.int(d.tildeKC), C.ShiftMask)
	}
	C.XFlush(d.dpy)

	go d.eventLoop()
	return nil
}

var ErrNoX11Display = &x11Error{"cannot open X11 display"}

type x11Error struct{ msg string }

func (e *x11Error) Error() string { return e.msg }

func (d *Daemon) SetEnabled(enabled bool) {
	d.mu.Lock()
	d.Eng.Cfg.Enabled = enabled
	d.Eng.Reset()
	cb := d.OnStateChanged
	mode := d.Eng.Cfg.Mode
	d.mu.Unlock()
	if cb != nil {
		cb(enabled, mode)
	}
}

func (d *Daemon) SetMode(mode engine.InputMode) {
	d.mu.Lock()
	d.Eng.Cfg.Mode = mode
	d.Eng.Reset()
	cb := d.OnStateChanged
	enabled := d.Eng.Cfg.Enabled
	d.mu.Unlock()
	if cb != nil {
		cb(enabled, mode)
	}
}

func (d *Daemon) ToggleEnabled() bool {
	d.mu.Lock()
	next := !d.Eng.Cfg.Enabled
	d.mu.Unlock()
	d.SetEnabled(next)
	return next
}

func (d *Daemon) injectReplacementLocked(backspaces int, text string) {
	if backspaces > 0 && d.backspaceKC != 0 {
		d.pendingFakeBackspaces += backspaces
		for i := 0; i < backspaces; i++ {
			C.tap_keycode(d.dpy, C.int(d.backspaceKC))
		}
	}

	for _, r := range text {
		km, ok := d.runeToKey[r]
		if !ok || km.kc == 0 {
			continue
		}
		wantUpper := 0
		if km.shift {
			wantUpper = 1
		}
		C.tap_with_shift_target(d.dpy, C.int(km.kc), C.int(wantUpper), C.int(d.shiftLKC), C.int(d.shiftRKC))
	}
	C.XFlush(d.dpy)
}

func (d *Daemon) eventLoop() {
	var ev C.VuiXEvent
	for {
		if C.next_vui_event(d.dpy, d.xiOpcode, &ev) == 0 {
			continue
		}

		d.mu.Lock()
		kc := int(ev.detail)
		state := uint(ev.state)
		sym := uint32(ev.keysym)

		// Ignore all 93 synthetic keycodes
		if d.syntheticKCSet[kc] {
			d.mu.Unlock()
			continue
		}

		// Mouse click resets syllable buffer
		if ev.kind == 2 {
			d.Eng.Reset()
			d.mu.Unlock()
			continue
		}

		// Check global toggle hotkeys (GrabbedKeyPress on Ctrl+Space or Shift+grave)
		if ev.kind == 3 {
			if (kc == d.spaceKC && (state&C.ControlMask) != 0) ||
				(kc == d.tildeKC && (state&C.ShiftMask) != 0) {
				C.release_active_grab(d.dpy)
				next := !d.Eng.Cfg.Enabled
				d.Eng.Cfg.Enabled = next
				d.Eng.Reset()
				cb := d.OnStateChanged
				mode := d.Eng.Cfg.Mode
				d.mu.Unlock()
				if cb != nil {
					cb(next, mode)
				}
				continue
			}
			d.mu.Unlock()
			continue
		}

		if !d.Eng.Cfg.Enabled {
			d.mu.Unlock()
			continue
		}

		// Ignore synthetic Backspaces that we injected ourselves
		if kc == d.backspaceKC && d.pendingFakeBackspaces > 0 {
			d.pendingFakeBackspaces--
			d.mu.Unlock()
			continue
		}

		// If Ctrl, Alt (Mod1), or Super (Mod4) is held, reset syllable state
		if (state & (C.ControlMask | C.Mod1Mask | C.Mod4Mask)) != 0 {
			d.Eng.Reset()
			d.mu.Unlock()
			continue
		}

		switch sym {
		case C.XK_BackSpace:
			d.Eng.HandleBackspace()
			d.mu.Unlock()
			continue
		case C.XK_Delete:
			d.Eng.HandleDelete()
			d.mu.Unlock()
			continue
		case C.XK_space:
			d.Eng.Reset()
			d.mu.Unlock()
			continue
		case C.XK_Return, C.XK_KP_Enter, C.XK_Tab, C.XK_Escape,
			C.XK_Left, C.XK_Right, C.XK_Up, C.XK_Down, C.XK_Home, C.XK_End, C.XK_Page_Up, C.XK_Page_Down:
			d.Eng.Reset()
			d.mu.Unlock()
			continue
		case C.XK_Shift_L, C.XK_Shift_R, C.XK_Caps_Lock, C.XK_Num_Lock:
			d.mu.Unlock()
			continue
		}

		if sym >= 0x20 && sym <= 0x7e {
			ch := rune(sym)
			d.TotalKeystrokes++
			act := d.Eng.ProcessKey(ch)
			if act.Swallow {
				d.VietnameseTransforms++
				d.injectReplacementLocked(act.Backspaces+1, act.CommitText)
			}
			d.mu.Unlock()
			continue
		}

		d.mu.Unlock()
	}
}
