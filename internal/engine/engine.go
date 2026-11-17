package engine

import (
	"unicode"
)

// InputMode defines the Vietnamese typing method.
type InputMode int

const (
	ModeTelex InputMode = iota
	ModeVNI
	ModeTelexVNI
)

func (m InputMode) String() string {
	switch m {
	case ModeTelex:
		return "Telex"
	case ModeVNI:
		return "VNI"
	case ModeTelexVNI:
		return "Telex + VNI"
	default:
		return "Telex"
	}
}

// Tone represents the 6 Vietnamese tones (0 = ngang/none, 1 = sắc, 2 = huyền, 3 = hỏi, 4 = ngã, 5 = nặng).
type Tone int

const (
	ToneNone  Tone = 0
	ToneSac   Tone = 1 // s / 1
	ToneHuyen Tone = 2 // f / 2
	ToneHoi   Tone = 3 // r / 3
	ToneNga   Tone = 4 // x / 4
	ToneNang  Tone = 5 // j / 5
)

// vowelTable maps a base vowel (with hat/horn/breve) to its 6 toned forms:
// [ToneNone, ToneSac, ToneHuyen, ToneHoi, ToneNga, ToneNang]
var vowelTable = map[rune][6]rune{
	'a': {'a', 'á', 'à', 'ả', 'ã', 'ạ'},
	'ă': {'ă', 'ắ', 'ằ', 'ẳ', 'ẵ', 'ặ'},
	'â': {'â', 'ấ', 'ầ', 'ẩ', 'ẫ', 'ậ'},
	'e': {'e', 'é', 'è', 'ẻ', 'ẽ', 'ẹ'},
	'ê': {'ê', 'ế', 'ề', 'ể', 'ễ', 'ệ'},
	'i': {'i', 'í', 'ì', 'ỉ', 'ĩ', 'ị'},
	'o': {'o', 'ó', 'ò', 'ỏ', 'õ', 'ọ'},
	'ô': {'ô', 'ố', 'ồ', 'ổ', 'ỗ', 'ộ'},
	'ơ': {'ơ', 'ớ', 'ờ', 'ở', 'ỡ', 'ợ'},
	'u': {'u', 'ú', 'ù', 'ủ', 'ũ', 'ụ'},
	'ư': {'ư', 'ứ', 'ừ', 'ử', 'ữ', 'ự'},
	'y': {'y', 'ý', 'ỳ', 'ỷ', 'ỹ', 'ỵ'},

	'A': {'A', 'Á', 'À', 'Ả', 'Ã', 'Ạ'},
	'Ă': {'Ă', 'Ắ', 'Ằ', 'Ẳ', 'Ẵ', 'Ặ'},
	'Â': {'Â', 'Ấ', 'Ầ', 'Ẩ', 'Ẫ', 'Ậ'},
	'E': {'E', 'É', 'È', 'Ẻ', 'Ẽ', 'Ẹ'},
	'Ê': {'Ê', 'Ế', 'Ề', 'Ể', 'Ễ', 'Ệ'},
	'I': {'I', 'Í', 'Ì', 'Ỉ', 'Ĩ', 'Ị'},
	'O': {'O', 'Ó', 'Ò', 'Ỏ', 'Õ', 'Ọ'},
	'Ô': {'Ô', 'Ố', 'Ồ', 'Ổ', 'Ỗ', 'Ộ'},
	'Ơ': {'Ơ', 'Ớ', 'Ờ', 'Ở', 'Ỡ', 'Ợ'},
	'U': {'U', 'Ú', 'Ù', 'Ủ', 'Ũ', 'Ụ'},
	'Ư': {'Ư', 'Ứ', 'Ừ', 'Ử', 'Ữ', 'Ự'},
	'Y': {'Y', 'Ý', 'Ỳ', 'Ỷ', 'Ỹ', 'Ỵ'},
}

type vowelMeta struct {
	base  rune // 'a', 'ă', 'â', 'e', 'ê', 'i', 'o', 'ô', 'ơ', 'u', 'ư', 'y' (preserving case)
	ascii rune // 'a', 'e', 'i', 'o', 'u', 'y' (lowercase)
	tone  Tone
	upper bool
}

var runeToMeta = map[rune]vowelMeta{}

func init() {
	asciiMap := map[rune]rune{
		'a': 'a', 'ă': 'a', 'â': 'a',
		'e': 'e', 'ê': 'e',
		'i': 'i',
		'o': 'o', 'ô': 'o', 'ơ': 'o',
		'u': 'u', 'ư': 'u',
		'y': 'y',
	}
	for base, forms := range vowelTable {
		lowerBase := unicode.ToLower(base)
		isUp := unicode.IsUpper(base)
		asc := asciiMap[lowerBase]
		for t, r := range forms {
			runeToMeta[r] = vowelMeta{
				base:  base,
				ascii: asc,
				tone:  Tone(t),
				upper: isUp,
			}
		}
	}
}

func isVietnameseVowel(r rune) bool {
	_, ok := runeToMeta[r]
	return ok
}

func applyToneToBase(base rune, t Tone) rune {
	if forms, ok := vowelTable[base]; ok {
		return forms[t]
	}
	return base
}

func changeVowelShape(oldRune, newLowerBase rune) rune {
	meta, ok := runeToMeta[oldRune]
	if !ok {
		return oldRune
	}
	targetBase := newLowerBase
	if meta.upper {
		targetBase = unicode.ToUpper(newLowerBase)
	}
	return applyToneToBase(targetBase, meta.tone)
}

// Config holds the user preferences for VuiKey's linguistic engine.
type Config struct {
	Enabled         bool              `json:"enabled"`
	Mode            InputMode         `json:"mode"`
	SmartEnglish    bool              `json:"smart_english"`
	SyllableSplit   bool              `json:"syllable_split"`
	ModernToneStyle bool              `json:"modern_tone_style"` // oà, uý vs òa, úy
	EnableMacros    bool              `json:"enable_macros"`
	Macros          map[string]string `json:"macros"`
}

func DefaultConfig() Config {
	return Config{
		Enabled:         true,
		Mode:            ModeTelex,
		SmartEnglish:    true,
		SyllableSplit:   true,
		ModernToneStyle: false,
		EnableMacros:    false,
		Macros: map[string]string{
			"vn": "Việt Nam",
		},
	}
}

// Action describes what the input method should do in response to a key press.
type Action struct {
	Swallow    bool   // If true, intercept the keystroke so the target app does not receive the raw key
	Backspaces int    // Number of Backspaces to send before inserting CommitText
	CommitText string // UTF-8 string to inject after Backspaces
}

// Engine implements VuiKey's Zero-Preedit + Sliding Syllable Window state machine.
type Engine struct {
	Cfg         Config
	syllable    []rune // Currently active syllable displayed on screen
	rawHistory  []rune // Raw keys typed for the current word
	wordSealed  bool   // True if the user undid a double-key or word is non-Vietnamese English
	lastUndoKey rune   // Key that just triggered an undo
	delayedD    bool   // True if the immediately preceding keystroke applied a delayed 'd' -> 'đ'
}

func NewEngine(cfg Config) *Engine {
	return &Engine{
		Cfg:        cfg,
		syllable:   make([]rune, 0, 16),
		rawHistory: make([]rune, 0, 16),
	}
}

// Reset clears the syllable state (called on Space, Enter, Delete, Arrow keys, Mouse click, Focus change).
// Crucial property of Zero-Preedit: Reset NEVER commits or re-sends anything to the application!
func (e *Engine) Reset() {
	e.syllable = e.syllable[:0]
	e.rawHistory = e.rawHistory[:0]
	e.wordSealed = false
	e.lastUndoKey = 0
	e.delayedD = false
}

// CurrentSyllable returns the currently tracked syllable on screen.
func (e *Engine) CurrentSyllable() string {
	return string(e.syllable)
}

// HandleBackspace updates the internal syllable state when the user presses hardware BackSpace.
func (e *Engine) HandleBackspace() {
	if len(e.syllable) > 0 {
		e.syllable = e.syllable[:len(e.syllable)-1]
	}
	// Always unseal the word on BackSpace so the user never has to delete all the way to the first letter!
	e.wordSealed = false
	e.lastUndoKey = 0
	e.delayedD = false
	e.rawHistory = append([]rune(nil), e.syllable...)
}

// HandleDelete handles hardware Delete key: simply resets syllable state without emitting anything!
func (e *Engine) HandleDelete() {
	e.Reset()
}

// CheckMacroOnSpace checks if the word just completed before Space matches a configured macro.
func (e *Engine) CheckMacroOnSpace() (Action, bool) {
	if !e.Cfg.Enabled || !e.Cfg.EnableMacros || len(e.syllable) == 0 {
		e.Reset()
		return Action{}, false
	}
	word := string(e.syllable)
	if replacement, ok := e.Cfg.Macros[word]; ok && replacement != "" {
		bs := len(e.syllable)
		e.Reset()
		return Action{
			Swallow:    false, // Let space happen or include space
			Backspaces: bs,
			CommitText: replacement,
		}, true
	}
	e.Reset()
	return Action{}, false
}

// WantedGrabKeys returns the exact lowercase ASCII characters that would trigger a Vietnamese
// transformation on the current syllable right now.
// This allows the X11 hook to dynamically grab ONLY the active modifier keys (0.0ms pass-through for all other keys).
func (e *Engine) WantedGrabKeys() []rune {
	if !e.Cfg.Enabled || e.wordSealed || len(e.syllable) == 0 {
		return nil
	}
	var keys []rune
	candidates := []rune{'s', 'f', 'r', 'x', 'j', 'z', 'a', 'e', 'o', 'w', 'd', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	for _, k := range candidates {
		if e.wouldTransform(k) {
			keys = append(keys, k)
		}
	}
	return keys
}

func (e *Engine) wouldTransform(ch rune) bool {
	clone := &Engine{
		Cfg:         e.Cfg,
		syllable:    append([]rune(nil), e.syllable...),
		rawHistory:  append([]rune(nil), e.rawHistory...),
		wordSealed:  e.wordSealed,
		lastUndoKey: e.lastUndoKey,
	}
	act := clone.ProcessKey(ch)
	return act.Swallow
}

// ProcessKey processes a printable character keystroke.
// If Action.Swallow is false, the caller lets the key pass straight through to the app (or commits `string(ch)`).
// If Action.Swallow is true, the caller intercepts the key, sends Action.Backspaces Backspaces, and injects Action.CommitText.
func (e *Engine) ProcessKey(ch rune) Action {
	if !e.Cfg.Enabled {
		return Action{Swallow: false}
	}

	// Non-letter/digit punctuation resets the syllable state.
	if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
		e.Reset()
		return Action{Swallow: false}
	}

	// Digits only participate if VNI mode is active and syllable is non-empty.
	if unicode.IsDigit(ch) && e.Cfg.Mode == ModeTelex {
		e.Reset()
		return Action{Swallow: false}
	}

	e.rawHistory = append(e.rawHistory, ch)

	if e.wordSealed {
		e.syllable = append(e.syllable, ch)
		return Action{Swallow: false}
	}

	// Check Smart English Guard on raw prefix
	if e.Cfg.SmartEnglish && isEnglishOnlyPrefix(e.syllable, ch) {
		e.wordSealed = true
		e.syllable = append(e.syllable, ch)
		return Action{Swallow: false}
	}

	// Sliding Syllable Window: if the previous keystroke was a delayed 'd' (e.g. "di" + 'd' -> "đi" or "dịu" + 'd' -> "địu")
	// and the user now types a non-modifier vowel (like 'a' in "didang" or "dịudàng"), that 'd' actually started a new syllable "da..."!
	if e.Cfg.SyllableSplit && e.delayedD && isVietnameseVowel(ch) {
		dec := decompose(e.syllable)
		isUO := len(dec.vowels) == 1 && lowerASCIIString(dec.vowels) == "u" && unicode.ToLower(ch) == 'o'
		if !isModifierForSyllable(dec, ch, e.Cfg.Mode) && !isUO {
			oldScreen := append([]rune(nil), e.syllable...)
			if e.syllable[0] == 'đ' {
				e.syllable[0] = 'd'
			} else if e.syllable[0] == 'Đ' {
				e.syllable[0] = 'D'
			}
			prevD := 'd'
			if len(e.rawHistory) >= 2 && unicode.IsUpper(e.rawHistory[len(e.rawHistory)-2]) {
				prevD = 'D'
			}
			targetScreen := append(append([]rune(nil), e.syllable...), prevD, ch)
			e.syllable = []rune{prevD, ch}
			e.rawHistory = []rune{prevD, ch}
			e.delayedD = false
			e.lastUndoKey = 0
			bs, ins := ComputeMinimalDiff(oldScreen, targetScreen)
			return Action{
				Swallow:    true,
				Backspaces: bs,
				CommitText: ins,
			}
		}
	}
	e.delayedD = false

	// Sliding Syllable Window: check if `ch` starts a brand new syllable in a concatenated string
	// (e.g., `nhac` + `d` -> `nhac` | `d`, `di` + `d` -> `di` | `d`, `dang` + `n` -> `dang` | `n`).
	if e.Cfg.SyllableSplit && shouldSplitNewSyllable(e.syllable, ch, e.Cfg.Mode) {
		e.syllable = []rune{ch}
		e.rawHistory = []rune{ch}
		e.lastUndoKey = 0
		return Action{Swallow: false}
	}

	// Try applying Vietnamese transformation (dd -> đ, vowel hats/horns, tones, or double-key undo)
	oldSyl := append([]rune(nil), e.syllable...)
	if newSyl, transformed, undone := e.tryTransform(oldSyl, ch); transformed {
		e.syllable = newSyl
		if undone {
			e.wordSealed = true
		} else if e.Cfg.SmartEnglish && checkAutoRestoreEnglish(e.syllable, e.rawHistory) {
			e.syllable = append([]rune(nil), e.rawHistory...)
			e.wordSealed = true
		}
		bs, ins := ComputeMinimalDiff(oldSyl, e.syllable)
		return Action{
			Swallow:    true,
			Backspaces: bs,
			CommitText: ins,
		}
	}

	// Normal character append
	e.syllable = append(e.syllable, ch)
	if e.Cfg.SmartEnglish && checkAutoRestoreEnglish(e.syllable, e.rawHistory) {
		e.syllable = append([]rune(nil), e.rawHistory...)
		e.wordSealed = true
		if string(oldSyl)+string(ch) != string(e.syllable) {
			bs, ins := ComputeMinimalDiff(oldSyl, e.syllable)
			return Action{
				Swallow:    true,
				Backspaces: bs,
				CommitText: ins,
			}
		}
	}
	return Action{Swallow: false}
}

// ComputeMinimalDiff computes the minimal number of Backspaces and replacement suffix
// needed to transform `oldSyl` into `newSyl`.
func ComputeMinimalDiff(oldSyl, newSyl []rune) (backspaces int, insert string) {
	common := 0
	minLen := len(oldSyl)
	if len(newSyl) < minLen {
		minLen = len(newSyl)
	}
	for common < minLen && oldSyl[common] == newSyl[common] {
		common++
	}
	backspaces = len(oldSyl) - common
	insert = string(newSyl[common:])
	return backspaces, insert
}

// ProcessKeyForIME processes a printable character for the Native Fcitx5 Engine.
// Returns:
// - commitPrefix: any completed syllable (from Sliding Syllable Window or word-ending punctuation/space) to commit immediately
// - currentPreedit: the active single syllable (at most 1 syllable, never an unspaced multi-syllable chain)
// - handled: true if the key was consumed by VuiKey
func (e *Engine) ProcessKeyForIME(ch rune) (commitPrefix string, currentPreedit string, handled bool) {
	if !e.Cfg.Enabled {
		return "", "", false
	}

	// Non-letter/digit (such as Space, punctuation): commit current syllable + ch together atomically!
	if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
		if len(e.syllable) > 0 {
			out := string(e.syllable) + string(ch)
			e.Reset()
			return out, "", true
		}
		e.Reset()
		return "", "", false
	}

	// Digits in Telex mode end the current syllable
	if unicode.IsDigit(ch) && e.Cfg.Mode == ModeTelex {
		if len(e.syllable) > 0 {
			out := string(e.syllable) + string(ch)
			e.Reset()
			return out, "", true
		}
		e.Reset()
		return "", "", false
	}

	e.rawHistory = append(e.rawHistory, ch)

	if e.wordSealed {
		e.syllable = append(e.syllable, ch)
		return "", string(e.syllable), true
	}

	if e.Cfg.SmartEnglish && isEnglishOnlyPrefix(e.syllable, ch) {
		e.wordSealed = true
		e.syllable = append(e.syllable, ch)
		return "", string(e.syllable), true
	}

	// Sliding Syllable Window: if the previous keystroke was a delayed 'd' (e.g. "di" + 'd' -> "đi" or "dịu" + 'd' -> "địu")
	// and the user now types a non-modifier vowel (like 'a' in "didang" or "dịudàng"), that 'd' actually started a new syllable "da..."!
	if e.Cfg.SyllableSplit && e.delayedD && isVietnameseVowel(ch) {
		dec := decompose(e.syllable)
		isUO := len(dec.vowels) == 1 && lowerASCIIString(dec.vowels) == "u" && unicode.ToLower(ch) == 'o'
		if !isModifierForSyllable(dec, ch, e.Cfg.Mode) && !isUO {
			if e.syllable[0] == 'đ' {
				e.syllable[0] = 'd'
			} else if e.syllable[0] == 'Đ' {
				e.syllable[0] = 'D'
			}
			prevD := 'd'
			if len(e.rawHistory) >= 2 && unicode.IsUpper(e.rawHistory[len(e.rawHistory)-2]) {
				prevD = 'D'
			}
			finished := string(e.syllable)
			e.syllable = []rune{prevD, ch}
			e.rawHistory = []rune{prevD, ch}
			e.delayedD = false
			e.lastUndoKey = 0
			return finished, string(e.syllable), true
		}
	}
	e.delayedD = false

	// Sliding Syllable Window: if `ch` starts a new syllable in an unspaced compound string
	// (e.g. nhacdidangnhactieng...), immediately commit the finished syllable and start a fresh 1-char syllable!
	if e.Cfg.SyllableSplit && shouldSplitNewSyllable(e.syllable, ch, e.Cfg.Mode) {
		finished := string(e.syllable)
		e.syllable = []rune{ch}
		e.rawHistory = []rune{ch}
		e.lastUndoKey = 0
		return finished, string(e.syllable), true
	}

	oldSyl := append([]rune(nil), e.syllable...)
	if newSyl, transformed, undone := e.tryTransform(oldSyl, ch); transformed {
		e.syllable = newSyl
		if undone {
			e.wordSealed = true
		} else if e.Cfg.SmartEnglish && checkAutoRestoreEnglish(e.syllable, e.rawHistory) {
			e.syllable = append([]rune(nil), e.rawHistory...)
			e.wordSealed = true
		}
		return "", string(e.syllable), true
	}

	e.syllable = append(e.syllable, ch)
	if e.Cfg.SmartEnglish && checkAutoRestoreEnglish(e.syllable, e.rawHistory) {
		e.syllable = append([]rune(nil), e.rawHistory...)
		e.wordSealed = true
	}
	return "", string(e.syllable), true
}

