package engine

import (
	"strings"
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
	Cfg              Config
	syllable         []rune // Currently active syllable displayed on screen
	rawHistory       []rune // Raw keys typed for the current word
	wordSealed       bool   // True if the user undid a double-key or word is non-Vietnamese English
	lastUndoKey      rune   // Key that just triggered an undo
	delayedD         bool   // True if the immediately preceding keystroke applied a delayed 'd' -> 'đ'
	committedLine    []rune // Current line committed text buffer (strictly committed text only)
	committedX       []int  // Recorded caret X pixel coordinates for each cursor index 0..len(committedLine)
	committedCursor  int    // Cursor index within committedLine
	lineY            int    // Last observed caret Y pixel coordinate
	lineH            int    // Last observed caret height in pixels
	justResetByMouse bool   // True right after InputContextReset (mouse click / selection)
	reeditActive     bool   // True ONLY when user explicitly positioned cursor at a committed word via arrow keys, backspace, or mouse click
}

func NewEngine(cfg Config) *Engine {
	return &Engine{
		Cfg:           cfg,
		syllable:      make([]rune, 0, 16),
		rawHistory:    make([]rune, 0, 16),
		committedLine: make([]rune, 0, 128),
		committedX:    make([]int, 1, 129),
	}
}

func (e *Engine) ensureCommittedXLen() {
	want := len(e.committedLine) + 1
	for len(e.committedX) < want {
		e.committedX = append(e.committedX, 0)
	}
	if len(e.committedX) > want {
		e.committedX = e.committedX[:want]
	}
}

// ClearCommittedLine clears the committed line buffer (called on Enter, Tab, Esc, Up/Down, Ctrl+A, Shift+Select, FocusOut).
func (e *Engine) ClearCommittedLine() {
	e.committedLine = e.committedLine[:0]
	e.committedX = e.committedX[:0]
	e.committedX = append(e.committedX, 0)
	e.committedCursor = 0
	e.lineY = 0
	e.lineH = 0
	e.justResetByMouse = false
	e.reeditActive = false
}

// RecordCommittedText inserts committed text at the current cursor position in the committed line buffer.
func (e *Engine) RecordCommittedText(text string) {
	if text == "" {
		return
	}
	r := []rune(text)
	if e.committedCursor < 0 {
		e.committedCursor = 0
	}
	if e.committedCursor > len(e.committedLine) {
		e.committedCursor = len(e.committedLine)
	}
	e.ensureCommittedXLen()
	tail := append([]rune(nil), e.committedLine[e.committedCursor:]...)
	e.committedLine = append(append(e.committedLine[:e.committedCursor], r...), tail...)

	xZeros := make([]int, len(r))
	xTail := append([]int(nil), e.committedX[e.committedCursor+1:]...)
	e.committedX = append(append(e.committedX[:e.committedCursor+1], xZeros...), xTail...)

	e.committedCursor += len(r)
	if len(e.committedLine) > 256 {
		drop := len(e.committedLine) - 256
		e.committedLine = e.committedLine[drop:]
		if len(e.committedX) > drop {
			e.committedX = e.committedX[drop:]
		}
		e.committedCursor -= drop
		if e.committedCursor < 0 {
			e.committedCursor = 0
		}
	}
	e.ensureCommittedXLen()
}

// RecordCommittedBackspace deletes `count` runes before the cursor in the committed line buffer.
func (e *Engine) RecordCommittedBackspace(count int) {
	if e.justResetByMouse {
		// User clicked/dragged with the mouse (potentially selecting all text) and pressed BackSpace:
		// clear the committed line buffer so stale text is never kept!
		e.ClearCommittedLine()
		return
	}
	e.ensureCommittedXLen()
	for i := 0; i < count; i++ {
		if e.committedCursor > 0 && e.committedCursor <= len(e.committedLine) {
			e.committedLine = append(e.committedLine[:e.committedCursor-1], e.committedLine[e.committedCursor:]...)
			if e.committedCursor < len(e.committedX) {
				e.committedX = append(e.committedX[:e.committedCursor], e.committedX[e.committedCursor+1:]...)
			}
			e.committedCursor--
		}
	}
	e.ensureCommittedXLen()
	e.reeditActive = e.committedCursor > 0 && e.committedCursor <= len(e.committedLine) && unicode.IsLetter(e.committedLine[e.committedCursor-1])
}

// RecordCommittedDelete deletes 1 rune at the cursor in the committed line buffer.
func (e *Engine) RecordCommittedDelete() {
	if e.justResetByMouse {
		e.ClearCommittedLine()
		return
	}
	e.ensureCommittedXLen()
	if e.committedCursor >= 0 && e.committedCursor < len(e.committedLine) {
		e.committedLine = append(e.committedLine[:e.committedCursor], e.committedLine[e.committedCursor+1:]...)
		if e.committedCursor+1 < len(e.committedX) {
			e.committedX = append(e.committedX[:e.committedCursor+1], e.committedX[e.committedCursor+2:]...)
		}
	}
	e.ensureCommittedXLen()
	e.reeditActive = e.committedCursor > 0 && e.committedCursor <= len(e.committedLine) && unicode.IsLetter(e.committedLine[e.committedCursor-1])
}

// ApplySurroundingReplacement updates the committed line buffer when a surrounding word before the cursor
// has been modified in-place by TryModifySurroundingWord, keeping reeditActive armed for chained modifiers.
func (e *Engine) ApplySurroundingReplacement(backspaces int, commitText string) {
	e.justResetByMouse = false
	e.RecordCommittedBackspace(backspaces)
	e.RecordCommittedText(commitText)
	e.MoveCommittedCursor(0)
}

// MoveCommittedCursor moves the cursor within the committed line buffer (-1 = Left, +1 = Right, -999 = Home, +999 = End).
func (e *Engine) MoveCommittedCursor(delta int) {
	e.justResetByMouse = false
	switch {
	case delta <= -999:
		e.committedCursor = 0
	case delta >= 999:
		e.committedCursor = len(e.committedLine)
	default:
		e.committedCursor += delta
		if e.committedCursor < 0 {
			e.committedCursor = 0
		}
		if e.committedCursor > len(e.committedLine) {
			e.committedCursor = len(e.committedLine)
		}
	}
	e.reeditActive = e.committedCursor > 0 && e.committedCursor <= len(e.committedLine) && unicode.IsLetter(e.committedLine[e.committedCursor-1])
}

// OnClientReset handles client-initiated IME reset (e.g. mouse click to reposition cursor).
// Any active preedit syllable is recorded as committed at the current position, and justResetByMouse is set.
func (e *Engine) OnClientReset() string {
	syl := string(e.syllable)
	e.Reset()
	if syl != "" {
		e.RecordCommittedText(syl)
		e.reeditActive = true
	} else {
		e.reeditActive = false
	}
	e.justResetByMouse = true
	return syl
}

// UpdateCursorRect updates or snaps the committedCursor based on the client's caret pixel coordinates (x, y, h).
func (e *Engine) UpdateCursorRect(x, y, h int) {
	if x <= 0 && y <= 0 {
		return
	}
	e.ensureCommittedXLen()
	if e.lineY != 0 && h > 0 {
		dy := y - e.lineY
		if dy < 0 {
			dy = -dy
		}
		if dy > h+6 {
			// Caret moved to a different visual line
			if e.justResetByMouse || len(e.syllable) == 0 {
				e.ClearCommittedLine()
			}
			e.lineY = y
			e.lineH = h
			return
		}
	}
	e.lineY = y
	if h > 0 {
		e.lineH = h
	}

	if !e.justResetByMouse {
		// Record the ground-truth X coordinate ONLY when no preedit syllable is active,
		// so committedX[committedCursor] always reflects the true committed DOM caret position.
		if len(e.syllable) == 0 && e.committedCursor >= 0 && e.committedCursor < len(e.committedX) {
			e.committedX[e.committedCursor] = x
		}
		return
	}

	// Mouse click repositioned the caret on the current line: snap committedCursor using recorded X coordinates!
	if len(e.committedLine) == 0 {
		e.committedCursor = 0
		e.reeditActive = false
		return
	}
	x0 := e.committedX[0]
	if x0 > 0 && x <= x0+3 {
		// Caret is at the very start of the input box (or input box was cleared)
		e.committedCursor = 0
		e.reeditActive = false
		return
	}

	// Interpolate any missing X coordinates between recorded anchor points, and extrapolate the trailing word
	interpX := make([]int, len(e.committedX))
	copy(interpX, e.committedX)
	firstKnownIdx := -1
	lastKnownIdx := -1
	for i := 0; i < len(interpX); i++ {
		if interpX[i] > 0 {
			if firstKnownIdx == -1 {
				firstKnownIdx = i
			}
			if lastKnownIdx >= 0 && i > lastKnownIdx+1 {
				startVal := interpX[lastKnownIdx]
				endVal := interpX[i]
				span := i - lastKnownIdx
				for k := lastKnownIdx + 1; k < i; k++ {
					interpX[k] = startVal + (endVal-startVal)*(k-lastKnownIdx)/span
				}
			}
			lastKnownIdx = i
		}
	}
	if lastKnownIdx < 0 {
		return
	}
	if lastKnownIdx == firstKnownIdx {
		// Only 1 anchor point known (e.g. x0 at start of line): if user clicked to the right of x0,
		// keep cursor at the end of the single committed word on this line.
		if x > interpX[firstKnownIdx]+3 {
			e.committedCursor = len(e.committedLine)
			e.reeditActive = e.committedCursor > 0 && unicode.IsLetter(e.committedLine[e.committedCursor-1])
		}
		return
	}

	// Extrapolate any trailing indices after lastKnownIdx using measured average character width on this line
	avgCharW := (interpX[lastKnownIdx] - interpX[firstKnownIdx]) / (lastKnownIdx - firstKnownIdx)
	if avgCharW < 4 {
		avgCharW = 8
	}
	for k := lastKnownIdx + 1; k < len(interpX); k++ {
		interpX[k] = interpX[lastKnownIdx] + (k-lastKnownIdx)*avgCharW
	}

	// Find the word-end index closest to the clicked X coordinate
	bestWordEnd := -1
	bestWordDist := 999999
	for i := 1; i <= len(e.committedLine); i++ {
		isWordEnd := unicode.IsLetter(e.committedLine[i-1]) && (i == len(e.committedLine) || !unicode.IsLetter(e.committedLine[i]))
		if !isWordEnd || interpX[i] <= 0 {
			continue
		}
		d := x - interpX[i]
		if d < 0 {
			d = -d
		}
		if d < bestWordDist {
			bestWordDist = d
			bestWordEnd = i
		}
	}
	if bestWordEnd > 0 && bestWordDist <= avgCharW*2+6 {
		e.committedCursor = bestWordEnd
		e.reeditActive = true
		return
	}

	// Fallback: nearest character boundary
	bestIdx := e.committedCursor
	bestDist := 999999
	for i := 0; i < len(interpX); i++ {
		if interpX[i] <= 0 {
			continue
		}
		d := x - interpX[i]
		if d < 0 {
			d = -d
		}
		if d < bestDist {
			bestDist = d
			bestIdx = i
		}
	}
	e.committedCursor = bestIdx
	e.reeditActive = e.committedCursor > 0 && e.committedCursor <= len(e.committedLine) && unicode.IsLetter(e.committedLine[e.committedCursor-1])
}

// CommittedTextBeforeCursor returns the committed text up to the cursor ONLY when reeditActive is true
// and the cursor is at the end of a word.
func (e *Engine) CommittedTextBeforeCursor() string {
	if !e.reeditActive || e.committedCursor <= 0 || e.committedCursor > len(e.committedLine) {
		return ""
	}
	// Only allow re-editing if the cursor is at the end of a word (not in the middle of letters)
	if e.committedCursor < len(e.committedLine) && unicode.IsLetter(e.committedLine[e.committedCursor]) {
		return ""
	}
	return string(e.committedLine[:e.committedCursor])
}

// Reset clears the syllable state (called on Space, Enter, Delete, Arrow keys, Mouse click, Focus change).
// Crucial property of Zero-Preedit: Reset NEVER commits or re-sends anything to the application!
func (e *Engine) Reset() {
	e.syllable = e.syllable[:0]
	e.rawHistory = e.rawHistory[:0]
	e.wordSealed = false
	e.lastUndoKey = 0
	e.delayedD = false
	e.reeditActive = false
}

// CurrentSyllable returns the currently tracked syllable on screen.
func (e *Engine) CurrentSyllable() string {
	return string(e.syllable)
}

// HandleBackspace updates the internal active syllable state when the user presses hardware BackSpace.
// Note: HandleBackspace ONLY modifies the uncommitted active syllable `e.syllable`.
// When `len(e.syllable) == 0`, Fcitx5 calls `RecordCommittedBackspace(1)` separately so committed text
// is never double-decremented when backspacing inside a preedit syllable!
func (e *Engine) HandleBackspace() {
	if len(e.syllable) > 0 {
		e.syllable = e.syllable[:len(e.syllable)-1]
	} else {
		e.RecordCommittedBackspace(1)
	}
	// Always unseal the word on BackSpace so the user never has to delete all the way to the first letter!
	e.wordSealed = false
	e.lastUndoKey = 0
	e.delayedD = false
	e.rawHistory = append([]rune(nil), e.syllable...)
}

// HandleDelete handles hardware Delete key: resets syllable state and deletes 1 char at cursor in committed line buffer.
func (e *Engine) HandleDelete() {
	e.Reset()
	e.RecordCommittedDelete()
}

// HandleNavigation updates the committed line cursor on Left(1), Right(2), Home(3), End(4), Delete(5), or Clear(0).
func (e *Engine) HandleNavigation(navType int) {
	if len(e.syllable) > 0 {
		e.RecordCommittedText(string(e.syllable))
	}
	e.Reset()
	switch navType {
	case 1:
		e.MoveCommittedCursor(-1)
	case 2:
		e.MoveCommittedCursor(1)
	case 3:
		e.MoveCommittedCursor(-999)
	case 4:
		e.MoveCommittedCursor(999)
	case 5:
		e.RecordCommittedDelete()
	default:
		e.ClearCommittedLine()
	}
}

// CheckMacroOnSpace checks if the word just completed before Space matches a configured macro.
func (e *Engine) CheckMacroOnSpace() (Action, bool) {
	e.justResetByMouse = false
	if !e.Cfg.Enabled || !e.Cfg.EnableMacros || len(e.syllable) == 0 {
		if len(e.syllable) > 0 {
			e.RecordCommittedText(string(e.syllable))
		}
		e.Reset()
		e.RecordCommittedText(" ")
		return Action{}, false
	}
	word := string(e.syllable)
	if replacement, ok := e.Cfg.Macros[word]; ok && replacement != "" {
		bs := len(e.syllable)
		e.Reset()
		e.RecordCommittedText(replacement + " ")
		return Action{
			Swallow:    false, // Let space happen or include space
			Backspaces: bs,
			CommitText: replacement,
		}, true
	}
	e.RecordCommittedText(word + " ")
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

// ProcessKey processes a printable character keystroke (Zero-Preedit mode).
// If Action.Swallow is false, the caller lets the key pass straight through to the app (or commits `string(ch)`).
// If Action.Swallow is true, the caller intercepts the key, sends Action.Backspaces Backspaces, and injects Action.CommitText.
func (e *Engine) ProcessKey(ch rune) Action {
	e.justResetByMouse = false
	if !e.Cfg.Enabled {
		return Action{Swallow: false}
	}

	// Non-letter/digit punctuation commits any active syllable and resets the syllable state.
	if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
		if len(e.syllable) > 0 {
			e.RecordCommittedText(string(e.syllable))
		}
		e.Reset()
		e.RecordCommittedText(string(ch))
		return Action{Swallow: false}
	}

	// Digits only participate if VNI mode is active and syllable is non-empty.
	if unicode.IsDigit(ch) && e.Cfg.Mode == ModeTelex {
		if len(e.syllable) > 0 {
			e.RecordCommittedText(string(e.syllable))
		}
		e.Reset()
		e.RecordCommittedText(string(ch))
		return Action{Swallow: false}
	}

	// If active syllable is empty (e.g. cursor was moved via Left/Right/Home/End/Backspace to the end of a word),
	// check if `ch` is a tone/diacritic modifier that validly transforms the word before the cursor!
	if len(e.syllable) == 0 && e.reeditActive {
		if tb := e.CommittedTextBeforeCursor(); tb != "" {
			if surAct, ok := e.TryModifySurroundingWord(tb, ch); ok && surAct.Swallow {
				e.RecordCommittedBackspace(surAct.Backspaces)
				e.RecordCommittedText(surAct.CommitText)
				e.reeditActive = true
				return surAct
			}
		}
		e.reeditActive = false
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
			e.RecordCommittedText(string(e.syllable))
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
		e.RecordCommittedText(string(e.syllable))
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

// ExtractLastSyllableBeforeCursor extracts the contiguous letter sequence immediately
// preceding the cursor (at the end of `textBeforeCursor`).
func ExtractLastSyllableBeforeCursor(textBeforeCursor string) []rune {
	runes := []rune(textBeforeCursor)
	if len(runes) == 0 {
		return nil
	}
	end := len(runes)
	start := end
	for start > 0 && unicode.IsLetter(runes[start-1]) {
		start--
	}
	word := runes[start:end]
	if len(word) == 0 || len(word) > 7 {
		return nil
	}
	if start > 0 && unicode.IsDigit(runes[start-1]) {
		return nil
	}
	return word
}

// IsValidVietnameseSyllableForReedit checks whether a standalone word before the cursor
// has a valid Vietnamese syllable structure and is not a protected English word.
func (e *Engine) IsValidVietnameseSyllableForReedit(word []rune) bool {
	if len(word) == 0 || len(word) > 7 {
		return false
	}
	if e.Cfg.SmartEnglish {
		rawLower := strings.ToLower(string(word))
		if commonEnglishTechWords[rawLower] {
			return false
		}
		for i := 0; i < len(word) && i < 3; i++ {
			if isEnglishOnlyPrefix(word[:i], word[i]) {
				return false
			}
		}
	}
	dec := decompose(word)
	initLower := lowerBaseString(dec.initial)
	if !validInitials[initLower] {
		return false
	}
	if len(dec.vowels) == 0 {
		return len(dec.initial) == 1 && (initLower == "d" || initLower == "đ") && len(dec.final) == 0
	}
	if len(dec.vowels) > 3 {
		return false
	}
	finalLower := lowerASCIIString(dec.final)
	if !validFinals[finalLower] {
		return false
	}
	if len(dec.final) > 0 {
		vBase := lowerBaseString(dec.vowels)
		if closedVowelNuclei[vBase] {
			return false
		}
	}
	return true
}

// TryModifySurroundingWord inspects the word immediately before the cursor when the active syllable buffer
// is empty (e.g. after moving the cursor via arrow keys, mouse click, or backspace to the end of an existing word).
// If `word` is a valid Vietnamese syllable and `ch` is a tone/diacritic key (or vowel/final consonant that rebalances `word`)
// without triggering a double-key undo, it returns Action(Swallow=true, Backspaces, CommitText).
func (e *Engine) TryModifySurroundingWord(textBeforeCursor string, ch rune) (Action, bool) {
	if !e.Cfg.Enabled || len(e.syllable) > 0 {
		return Action{}, false
	}
	if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
		return Action{}, false
	}
	if unicode.IsDigit(ch) && e.Cfg.Mode == ModeTelex {
		return Action{}, false
	}
	word := ExtractLastSyllableBeforeCursor(textBeforeCursor)
	if !e.IsValidVietnameseSyllableForReedit(word) {
		return Action{}, false
	}
	dec := decompose(word)
	lc := unicode.ToLower(ch)
	canRebalanceContinuation := len(dec.vowels) > 0 && (isVietnameseVowel(ch) || validFinals[lowerASCIIString(dec.final)+string(lc)])
	if !isModifierForSyllable(dec, ch, e.Cfg.Mode) && !canRebalanceContinuation {
		return Action{}, false
	}

	savedRaw := e.rawHistory
	savedDelayedD := e.delayedD
	e.rawHistory = append(append([]rune(nil), word...), ch)
	newSyl, transformed, undone := e.tryTransform(word, ch)
	e.rawHistory = savedRaw
	e.delayedD = savedDelayedD

	if !transformed || len(newSyl) == 0 || undone {
		return Action{}, false
	}
	if e.Cfg.SmartEnglish && checkAutoRestoreEnglish(newSyl, append(append([]rune(nil), word...), ch)) {
		return Action{}, false
	}
	bs, ins := ComputeMinimalDiff(word, newSyl)
	if bs == 0 && ins == "" {
		return Action{}, false
	}
	return Action{
		Swallow:    true,
		Backspaces: bs,
		CommitText: ins,
	}, true
}

// ProcessKeyForIME processes a printable character for the Native Fcitx5 Engine.
// Returns:
// - commitPrefix: any completed syllable (from Sliding Syllable Window or word-ending punctuation/space) to commit immediately
// - currentPreedit: the active single syllable (at most 1 syllable, never an unspaced multi-syllable chain)
// - handled: true if the key was consumed by VuiKey
func (e *Engine) ProcessKeyForIME(ch rune) (commitPrefix string, currentPreedit string, handled bool) {
	e.justResetByMouse = false
	if !e.Cfg.Enabled {
		return "", "", false
	}

	// Non-letter/digit (such as Space, punctuation): commit current syllable + ch together atomically!
	if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
		if len(e.syllable) > 0 {
			out := string(e.syllable) + string(ch)
			e.Reset()
			e.RecordCommittedText(out)
			return out, "", true
		}
		e.Reset()
		e.RecordCommittedText(string(ch))
		return "", "", false
	}

	// Digits in Telex mode end the current syllable
	if unicode.IsDigit(ch) && e.Cfg.Mode == ModeTelex {
		if len(e.syllable) > 0 {
			out := string(e.syllable) + string(ch)
			e.Reset()
			e.RecordCommittedText(out)
			return out, "", true
		}
		e.Reset()
		e.RecordCommittedText(string(ch))
		return "", "", false
	}

	// If active syllable is empty and the cursor was explicitly placed at the end of an already-committed
	// valid Vietnamese syllable prefix (e.g. user backspaced "trương " to "trươ" and types 'n', 'g', 'f'),
	// commit continuation letters directly into `committedLine` so the word stays unified for subsequent tone keys!
	if len(e.syllable) == 0 && e.reeditActive && unicode.IsLetter(ch) {
		if tb := e.CommittedTextBeforeCursor(); tb != "" {
			if word := ExtractLastSyllableBeforeCursor(tb); len(word) > 0 && e.IsValidVietnameseSyllableForReedit(word) {
				if !shouldSplitNewSyllable(word, ch, e.Cfg.Mode) {
					combined := append(append([]rune(nil), word...), ch)
					if e.IsValidVietnameseSyllableForReedit(combined) {
						s := string(ch)
						e.RecordCommittedText(s)
						e.reeditActive = true
						return s, "", true
					}
				}
			}
		}
		e.reeditActive = false
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
			e.RecordCommittedText(finished)
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
		e.RecordCommittedText(finished)
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

