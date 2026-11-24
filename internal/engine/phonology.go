package engine

import (
	"strings"
	"unicode"
)

// decomposedSyllable splits a Vietnamese syllable into:
// - initial consonant cluster (phụ âm đầu)
// - vowel cluster (vần nguyên âm)
// - final consonant cluster (phụ âm cuối)
type decomposedSyllable struct {
	initial    []rune
	vowels     []rune
	final      []rune
	vowelStart int
	vowelEnd   int // exclusive
}

func decompose(syl []rune) decomposedSyllable {
	n := len(syl)
	if n == 0 {
		return decomposedSyllable{}
	}

	// Find first vowel
	vStart := -1
	for i, r := range syl {
		if isVietnameseVowel(r) {
			// Special case: "qu" and "gi" followed by another vowel -> 'u' or 'i' is part of initial consonant!
			if i == 1 && i+1 < n && isVietnameseVowel(syl[i+1]) {
				firstLower := unicode.ToLower(syl[0])
				curASCII := runeToMeta[r].ascii
				if (firstLower == 'q' && curASCII == 'u') || (firstLower == 'g' && curASCII == 'i') {
					continue
				}
			}
			vStart = i
			break
		}
	}

	if vStart == -1 {
		return decomposedSyllable{
			initial:    syl,
			vowelStart: n,
			vowelEnd:   n,
		}
	}

	vEnd := vStart
	for vEnd < n && isVietnameseVowel(syl[vEnd]) {
		vEnd++
	}

	return decomposedSyllable{
		initial:    syl[:vStart],
		vowels:     syl[vStart:vEnd],
		final:      syl[vEnd:],
		vowelStart: vStart,
		vowelEnd:   vEnd,
	}
}

// Valid Vietnamese initial consonants
var validInitials = map[string]bool{
	"": true, "b": true, "c": true, "ch": true, "d": true, "đ": true,
	"g": true, "gh": true, "gi": true, "h": true, "k": true, "kh": true,
	"l": true, "m": true, "n": true, "ng": true, "ngh": true, "nh": true,
	"p": true, "ph": true, "qu": true, "r": true, "s": true, "t": true,
	"th": true, "tr": true, "v": true, "x": true,
}

// Valid Vietnamese final consonants
var validFinals = map[string]bool{
	"": true, "c": true, "ch": true, "m": true, "n": true, "ng": true, "nh": true, "p": true, "t": true,
}

// Diphthongs/triphthongs that are ALREADY closed (end in semivowel i/y/u/o and cannot accept any final consonant)
var closedVowelNuclei = map[string]bool{
	"ai": true, "ao": true, "au": true, "ay": true, "âu": true, "ây": true,
	"eo": true, "êu": true, "iu": true, "oi": true, "ôi": true, "ơi": true,
	"ui": true, "ưi": true, "ưu": true, "uya": true, "uyu": true,
	"oai": true, "oao": true, "oay": true, "oeo": true, "uây": true, "uôi": true, "ươi": true, "ươu": true,
}

func lowerASCIIString(runes []rune) string {
	var b strings.Builder
	for _, r := range runes {
		if meta, ok := runeToMeta[r]; ok {
			b.WriteRune(meta.ascii)
		} else {
			lc := unicode.ToLower(r)
			if lc == 'đ' {
				b.WriteRune('d')
			} else {
				b.WriteRune(lc)
			}
		}
	}
	return b.String()
}

func lowerBaseString(runes []rune) string {
	var b strings.Builder
	for _, r := range runes {
		if meta, ok := runeToMeta[r]; ok {
			b.WriteRune(unicode.ToLower(meta.base))
		} else {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

// isEnglishOnlyPrefix detects non-Vietnamese English prefixes (like "cl", "br", "st", "sh", "window", "view", "fix").
func isEnglishOnlyPrefix(syl []rune, nextCh rune) bool {
	combined := append(append([]rune(nil), syl...), nextCh)
	s := lowerBaseString(combined)
	if len(s) == 1 {
		// Starting a word with f, j, w, z in English (file, json, window, zero)
		return s == "f" || s == "j" || s == "w" || s == "z"
	}
	if len(s) == 2 {
		// Consonant clusters that never start a Vietnamese word
		invalidStarts := map[string]bool{
			"bl": true, "br": true, "cl": true, "cr": true, "dr": true,
			"fl": true, "fr": true, "gl": true, "gr": true, "kn": true,
			"pl": true, "pr": true, "sc": true, "sk": true, "sl": true,
			"sm": true, "sn": true, "sp": true, "sr": true, "st": true,
			"sw": true, "tw": true, "wr": true, "sh": true, "wh": true,
			"ck": true, "ab": true, "ad": true, "ag": true,
			"ak": true, "al": true, "ar": true, "av": true,
		}
		if invalidStarts[s] {
			return true
		}
	}
	if len(s) == 3 {
		if s == "thr" || s == "chr" || s == "str" || s == "spr" || s == "scr" || s == "shr" || s == "spl" {
			return true
		}
	}
	return false
}

// commonEnglishTechWords contains high-frequency English words whose keystrokes collide with Telex
// (e.g. "test" = t-e-s-t -> "tét", "list" = l-i-s-t -> "lít", "post" = p-o-s-t -> "pót").
var commonEnglishTechWords = map[string]bool{
	"test": true, "tests": true, "testing": true, "tested": true, "tester": true,
	"best": true, "rest": true, "list": true, "lists": true, "last": true,
	"post": true, "posts": true, "host": true, "hosts": true, "cost": true, "costs": true,
	"most": true, "lost": true, "must": true, "rust": true, "dust": true,
	"user": true, "users": true, "use": true, "used": true, "uses": true,
	"server": true, "servers": true, "service": true, "services": true,
	"system": true, "systems": true, "status": true, "static": true,
	"linux": true, "ubuntu": true, "debian": true, "fedora": true, "arch": true,
	"docker": true, "github": true, "gitlab": true, "google": true, "chrome": true,
	"python": true, "golang": true, "node": true, "nodes": true, "code": true, "codes": true,
	"mode": true, "modes": true, "model": true, "models": true, "module": true, "modules": true,
	"file": true, "files": true, "folder": true, "folders": true, "filter": true, "filters": true,
	"master": true, "worker": true, "router": true, "render": true, "buffer": true,
	"number": true, "member": true, "driver": true, "player": true, "super": true,
	"power": true, "tower": true, "lower": true, "upper": true, "inner": true, "outer": true,
	"other": true, "after": true, "under": true, "over": true, "never": true, "ever": true,
	"enter": true, "center": true, "header": true, "footer": true, "layer": true,
	"table": true, "enable": true, "disable": true, "value": true, "values": true,
	"issue": true, "issues": true, "true": true, "false": true, "style": true, "title": true,
	"time": true, "times": true, "name": true, "names": true, "game": true, "games": true,
	"same": true, "state": true, "date": true, "rate": true, "create": true, "update": true,
	"delete": true, "remote": true, "note": true, "notes": true, "size": true, "type": true,
	"types": true, "byte": true, "bytes": true, "site": true, "sites": true, "core": true,
	"more": true, "store": true, "score": true, "pure": true, "sure": true, "fire": true,
	"wire": true, "base": true, "case": true, "cases": true, "close": true, "mouse": true,
	"house": true, "route": true, "rule": true, "rules": true, "role": true, "roles": true,
	"scale": true, "sale": true, "mile": true, "profile": true, "mobile": true,
	"task": true, "tasks": true, "mask": true, "desk": true, "disk": true, "risk": true,
	"dark": true, "mark": true, "park": true, "card": true, "hard": true, "bird": true,
	"port": true, "ports": true, "sort": true, "part": true, "parts": true, "cart": true,
	"hurt": true, "next": true, "text": true, "texts": true, "left": true, "soft": true,
	"lift": true, "self": true, "half": true, "golf": true, "help": true, "hold": true,
	"cold": true, "gold": true, "bold": true, "told": true, "sold": true, "open": true,
	"item": true, "items": true, "icon": true, "icons": true, "edit": true, "init": true,
	"unit": true, "exit": true, "admin": true, "login": true, "token": true, "panel": true,
	"level": true, "hotel": true, "pixel": true, "local": true, "total": true, "final": true,
	"arrow": true, "arrows": true, "array": true, "arrays": true, "error": true, "errors": true,
	"allow": true, "asset": true, "assets": true, "async": true, "audio": true, "about": true,
}

// checkAutoRestoreEnglish checks if adding `nextCh` makes `rawHistory` an unmistakable English word
// (such as "test", "Linux", "left", "next", "port", "task", "server") and returns true if it should
// be restored to raw ASCII (`rawHistory`) and sealed as English.
func checkAutoRestoreEnglish(currentSyl []rune, rawHistory []rune) bool {
	if len(rawHistory) < 3 {
		return false
	}
	var b strings.Builder
	for _, r := range rawHistory {
		b.WriteRune(unicode.ToLower(r))
	}
	rawLower := b.String()

	if commonEnglishTechWords[rawLower] {
		return true
	}

	// Never auto-restore to English if the syllable already contains an explicit Vietnamese vowel/consonant
	// (ă, â, đ, ê, ô, ơ, ư)
	baseLower := lowerBaseString(currentSyl)
	if strings.ContainsAny(baseLower, "ăâđêôơư") {
		return false
	}

	dec := decompose(currentSyl)
	if len(dec.final) > 0 {
		// Phonotactic Law: Vowel + Consonant + Vowel (VCV, like "Linu" in "Linux", "bonu" in "bonus", "mode", "code")
		// never exists in a single Vietnamese syllable!
		for _, r := range dec.final {
			if isVietnameseVowel(r) {
				return true
			}
		}
	}

	return false
}

// isModifierForSyllable returns true if `ch` is a Telex/VNI modifier that applies to `dec`.
func isModifierForSyllable(dec decomposedSyllable, ch rune, mode InputMode) bool {
	lc := unicode.ToLower(ch)

	// 'd' is a modifier in Telex if the syllable starts with 'd'/'D' (allowing free-style delayed 'd' like "dodongj", "doongdj", "doongjd" -> "động"),
	// or if the syllable is literally "đ"/"Đ" (for immediate double-'d' undo "dd" -> "d").
	if lc == 'd' && (mode == ModeTelex || mode == ModeTelexVNI) {
		if len(dec.initial) == 1 {
			first := unicode.ToLower(dec.initial[0])
			if first == 'd' {
				return true
			}
			if first == 'đ' && len(dec.vowels) == 0 {
				return true
			}
		}
		return false
	}

	if len(dec.vowels) == 0 {
		return false
	}

	if mode == ModeTelex || mode == ModeTelexVNI {
		switch lc {
		case 's', 'f', 'r', 'x', 'j', 'z':
			return true
		case 'w':
			// 'w' modifies a, o, u
			vStr := lowerASCIIString(dec.vowels)
			return strings.ContainsAny(vStr, "aou")
		case 'a':
			// 'a' only modifies if there's already an 'a' in vowels (e.g. "nha" + 'a' -> "nhâ", "toan" + 'a' -> "toân" / wait: only when final is empty or valid with â)
			vBase := lowerBaseString(dec.vowels)
			return strings.Contains(vBase, "a") || strings.Contains(vBase, "â")
		case 'e':
			vBase := lowerBaseString(dec.vowels)
			return strings.Contains(vBase, "e") || strings.Contains(vBase, "ê")
		case 'o':
			vBase := lowerBaseString(dec.vowels)
			return strings.Contains(vBase, "o") || strings.Contains(vBase, "ô")
		}
	}

	if mode == ModeVNI || mode == ModeTelexVNI {
		switch lc {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return true
		}
	}

	return false
}

// shouldSplitNewSyllable determines if typing `ch` after `syl` begins a brand new Vietnamese syllable
// in a concatenated string like `nhacdidangnhactiengdidangtiengchoitai`.
func shouldSplitNewSyllable(syl []rune, ch rune, mode InputMode) bool {
	if len(syl) == 0 {
		return false
	}
	dec := decompose(syl)
	if len(dec.vowels) == 0 {
		return false
	}

	// If `ch` is a valid Telex/VNI modifier for the current syllable, do NOT split!
	if isModifierForSyllable(dec, ch, mode) {
		// Exception: what if `syl` already has a final consonant and `ch` is a vowel like `a`, `e`, `o`?
		// e.g. `nhac` + `a`: in Vietnamese, nobody types `nhaca` to mean `nhâc` when writing concatenated words,
		// OR wait: even if someone types `nhacdidang`, after `nhac`, the next char is `d` (not a modifier!).
		// What if someone types `tieng` + `a`? `tieng` has `e`, if they type `a`, `tienga` is a new syllable `a` unless they typed `tieng`!
		// Even better: if `len(dec.final) > 0` and `ch` is a vowel (`a`, `e`, `o`), only allow `a/e/o` hat modification if the vowel is not already hatted and forms a valid rhyme!
		return false
	}

	lc := unicode.ToLower(ch)

	// Case 1: The syllable ALREADY has a final consonant (e.g., `nhac`, `dang`, `tieng`).
	if len(dec.final) > 0 {
		newFinal := lowerASCIIString(dec.final) + string(lc)
		// Only "ch", "nh", "ng" can extend a 1-char final consonant ("c"->"ch", "n"->"nh", "n"->"ng")!
		if validFinals[newFinal] {
			return false
		}
		// If `syl` has NO Vietnamese tone/diacritic (pure ASCII like "Lin" in "Linux", "bon" in "bonus")
		// and `ch` is a non-modifier vowel ('u', 'i', 'e'...), do NOT split into a fake 1-vowel syllable:
		// let it stay in the same word so Smart English Guard seals "Linux", "bonus", "minus", "focus"!
		if isVietnameseVowel(ch) && string(syl) == lowerASCIIString(syl) || (len(syl) > 0 && strings.EqualFold(string(syl), lowerASCIIString(syl))) {
			if isVietnameseVowel(ch) {
				return false
			}
		}
		return true
	}

	// Case 2: The syllable has vowels and NO final consonant yet (e.g., `di`, `choi`, `tai`).
	if !isVietnameseVowel(ch) {
		// `ch` is a consonant after a vowel nucleus.
		// Check 2a: Can `ch` even be the start of a valid Vietnamese final consonant?
		// Only c, m, n, p, t can start a final consonant!
		// Consonants like b, d, g, h, k, l, q, r, s, v, x, f, j, w, z can NEVER start a final consonant after a vowel!
		switch lc {
		case 'c', 'm', 'n', 'p', 't':
			// Check 2b: Even if `lc` is c/m/n/p/t, is the current vowel nucleus ALREADY closed (like `choi`, `tai`, `nhau`, `kêu`)?
			vBase := lowerBaseString(dec.vowels)
			if closedVowelNuclei[vBase] {
				// Example: `choi` + `t` -> `oi` is a closed diphthong, cannot take `t` -> SPLITS into `choi` | `t`!
				return true
			}
			return false
		default:
			// Example: `di` + `d` -> `d` cannot be a final consonant -> SPLITS into `di` | `d`!
			return true
		}
	}

	// Case 3: `ch` is another vowel after an existing vowel nucleus (e.g., `nha` + `i` -> `nhai`, or `tai` + `a` -> `tai` | `a`).
	if len(dec.vowels) >= 3 {
		return true
	}
	vBase := lowerBaseString(dec.vowels)
	if closedVowelNuclei[vBase] {
		return true
	}

	return false
}

// findToneTargetIndex returns the index within `syl` (between `dec.vowelStart` and `dec.vowelEnd-1`)
// where the Vietnamese tone mark belongs according to standard orthography.
func findToneTargetIndex(dec decomposedSyllable, modernStyle bool) int {
	vCount := len(dec.vowels)
	if vCount == 0 {
		return -1
	}
	if vCount == 1 {
		return dec.vowelStart
	}

	// Rule 1: If any vowel already carries a diacritic (ă, â, ê, ô, ơ, ư), it takes the tone!
	// If there are two (like "ươ" in "trường", "uơ"), the second one ("ơ") takes the tone!
	lastSpecial := -1
	for i, r := range dec.vowels {
		meta := runeToMeta[r]
		lb := unicode.ToLower(meta.base)
		if lb == 'ă' || lb == 'â' || lb == 'ê' || lb == 'ô' || lb == 'ơ' || lb == 'ư' {
			lastSpecial = dec.vowelStart + i
		}
	}
	if lastSpecial != -1 {
		return lastSpecial
	}

	// Rule 2: 3 vowels without diacritics (e.g., "oai", "uay", "uya") -> middle vowel takes the tone
	if vCount == 3 {
		return dec.vowelStart + 1
	}

	// Rule 3: 2 vowels without diacritics (e.g., "ai", "ao", "oa", "oe", "uy", "ia", "ua")
	// If there is a final consonant (e.g., "hoan", "toan", "quyen"), the second vowel takes the tone!
	if len(dec.final) > 0 {
		return dec.vowelStart + 1
	}

	// Open syllable with 2 vowels:
	vStr := lowerASCIIString(dec.vowels)
	if modernStyle && (vStr == "oa" || vStr == "oe" || vStr == "uy") {
		return dec.vowelStart + 1
	}
	return dec.vowelStart
}

// rebalanceTone ensures the syllable's tone is on the correct vowel after adding vowels/final consonants,
// and also normalizes "ưo" + final/i/u -> "ươ" (e.g., dduwocj -> được, truwongf -> trường, nguwoif -> người).
func rebalanceTone(syl []rune, modernStyle bool) []rune {
	dec := decompose(syl)
	if len(dec.vowels) == 0 {
		return syl
	}

	// Auto-normalize "ưo" -> "ươ" when followed by a final consonant or 3rd vowel (i/u)
	if len(dec.vowels) >= 2 && (len(dec.final) > 0 || len(dec.vowels) >= 3) {
		for i := dec.vowelStart; i+1 < dec.vowelEnd; i++ {
			m1 := runeToMeta[syl[i]]
			m2 := runeToMeta[syl[i+1]]
			if unicode.ToLower(m1.base) == 'ư' && unicode.ToLower(m2.base) == 'o' {
				syl[i+1] = changeVowelShape(syl[i+1], 'ơ')
			}
		}
	}

	currentTone := ToneNone
	// Scan from index 0 (including "gi"/"qu" in case tone was placed on "gì" before 'a' was typed)
	for i := 0; i < dec.vowelEnd; i++ {
		if meta, ok := runeToMeta[syl[i]]; ok && meta.tone != ToneNone {
			currentTone = meta.tone
			syl[i] = applyToneToBase(meta.base, ToneNone)
		}
	}
	if currentTone != ToneNone {
		target := findToneTargetIndex(decompose(syl), modernStyle)
		if target >= 0 {
			meta := runeToMeta[syl[target]]
			syl[target] = applyToneToBase(meta.base, currentTone)
		}
	}
	return syl
}

// tryTransform attempts to apply a Telex or VNI modifier `ch` to `oldSyl`.
// Returns (newSyllable, transformed, undone).
func (e *Engine) tryTransform(oldSyl []rune, ch rune) ([]rune, bool, bool) {
	if len(oldSyl) == 0 {
		return nil, false, false
	}

	syl := append([]rune(nil), oldSyl...)
	dec := decompose(syl)
	lc := unicode.ToLower(ch)

	useTelex := e.Cfg.Mode == ModeTelex || e.Cfg.Mode == ModeTelexVNI
	useVNI := e.Cfg.Mode == ModeVNI || e.Cfg.Mode == ModeTelexVNI

	// 1. Handle 'd' -> 'đ' anywhere in the syllable (Telex free-style: "dd...", "dodongj", "doongdj", "doongjd" -> "động") or '9' -> 'đ' (VNI)
	if (useTelex && lc == 'd' && len(dec.initial) == 1) ||
		(useVNI && lc == '9' && len(dec.initial) >= 1) {
		first := syl[0]
		prevWasD := len(e.rawHistory) >= 2 && (unicode.ToLower(e.rawHistory[len(e.rawHistory)-2]) == 'd' || e.rawHistory[len(e.rawHistory)-2] == '9')
		switch first {
		case 'd':
			syl[0] = 'đ'
			e.delayedD = len(dec.vowels) > 0
			return syl, true, false
		case 'D':
			syl[0] = 'Đ'
			e.delayedD = len(dec.vowels) > 0
			return syl, true, false
		case 'đ':
			if len(dec.vowels) == 0 || prevWasD {
				// Double-key undo: đ + d -> dd (or dod + d -> dod)
				syl[0] = 'd'
				if len(dec.vowels) == 0 {
					syl = append(syl, ch)
				}
				return syl, true, true
			}
		case 'Đ':
			if len(dec.vowels) == 0 || prevWasD {
				syl[0] = 'D'
				if len(dec.vowels) == 0 {
					syl = append(syl, ch)
				}
				return syl, true, true
			}
		}
	}

	if len(dec.vowels) == 0 {
		// Support 'w' -> 'ư' after valid Vietnamese initial consonants (e.g. "phw" -> "phư", "nhw" -> "như", "thw" -> "thư")
		// This also fixes rollover typing when 'w' is pressed 1ms before 'u' in "phức" (p-h-w-u-c-s)!
		if useTelex && lc == 'w' && len(dec.initial) > 0 {
			initLower := lowerBaseString(dec.initial)
			if validInitials[initLower] {
				if unicode.IsUpper( ch ) {
					syl = append(syl, 'Ư')
				} else {
					syl = append(syl, 'ư')
				}
				return syl, true, false
			}
		}
		return nil, false, false
	}

	// If 'u' is typed right after 'ư' that was created by shorthand/rollover 'w' (e.g. p-h-w-u-c-s -> "phức"),
	// absorb the 'u' into 'ư' so it stays "phư" instead of becoming "phưu"!
	if useTelex && lc == 'u' && len(dec.vowels) == 1 && len(dec.final) == 0 && unicode.ToLower(runeToMeta[dec.vowels[0]].base) == 'ư' {
		if len(e.rawHistory) >= 2 && unicode.ToLower(e.rawHistory[len(e.rawHistory)-2]) == 'w' {
			return syl, true, false
		}
	}

	// 2. Tone keys: s/1 (Sắc), f/2 (Huyền), r/3 (Hỏi), x/4 (Ngã), j/5 (Nặng), z/0 (Xóa dấu)
	var wantedTone Tone = -1
	isToneKey := false
	if useTelex {
		switch lc {
		case 's':
			wantedTone, isToneKey = ToneSac, true
		case 'f':
			wantedTone, isToneKey = ToneHuyen, true
		case 'r':
			wantedTone, isToneKey = ToneHoi, true
		case 'x':
			wantedTone, isToneKey = ToneNga, true
		case 'j':
			wantedTone, isToneKey = ToneNang, true
		case 'z':
			wantedTone, isToneKey = ToneNone, true
		}
	}
	if !isToneKey && useVNI {
		switch lc {
		case '1':
			wantedTone, isToneKey = ToneSac, true
		case '2':
			wantedTone, isToneKey = ToneHuyen, true
		case '3':
			wantedTone, isToneKey = ToneHoi, true
		case '4':
			wantedTone, isToneKey = ToneNga, true
		case '5':
			wantedTone, isToneKey = ToneNang, true
		case '0':
			wantedTone, isToneKey = ToneNone, true
		}
	}

	if isToneKey {
		// Check current tone on syllable
		curTone := ToneNone
		curToneIdx := -1
		for i := dec.vowelStart; i < dec.vowelEnd; i++ {
			meta := runeToMeta[syl[i]]
			if meta.tone != ToneNone {
				curTone = meta.tone
				curToneIdx = i
				break
			}
		}

		if wantedTone == ToneNone {
			if curTone != ToneNone {
				meta := runeToMeta[syl[curToneIdx]]
				syl[curToneIdx] = applyToneToBase(meta.base, ToneNone)
				return syl, true, false
			}
			// Second-stage 'z'/'0': if syllable has no tone, strip Vietnamese vowel hats/horns/breve (ă, â, ê, ô, ơ, ư) or đ/Đ
			stripped := false
			for i, r := range syl {
				if r == 'đ' {
					syl[i] = 'd'
					stripped = true
				} else if r == 'Đ' {
					syl[i] = 'D'
					stripped = true
				} else if meta, ok := runeToMeta[r]; ok && unicode.ToLower(meta.base) != meta.ascii {
					syl[i] = changeVowelShape(r, meta.ascii)
					stripped = true
				}
			}
			if stripped {
				return syl, true, false
			}
			return nil, false, false
		}

		// Double-key undo: pressing the same tone key twice removes the tone and appends the raw key!
		if curTone == wantedTone {
			meta := runeToMeta[syl[curToneIdx]]
			syl[curToneIdx] = applyToneToBase(meta.base, ToneNone)
			syl = append(syl, ch)
			return syl, true, true
		}

		// Clear old tone if any
		if curToneIdx != -1 {
			meta := runeToMeta[syl[curToneIdx]]
			syl[curToneIdx] = applyToneToBase(meta.base, ToneNone)
		}

		targetIdx := findToneTargetIndex(decompose(syl), e.Cfg.ModernToneStyle)
		if targetIdx >= 0 {
			meta := runeToMeta[syl[targetIdx]]
			syl[targetIdx] = applyToneToBase(meta.base, wantedTone)
			return syl, true, false
		}
	}

	// 3. Vowel shape modifiers:
	// Telex 'a', 'e', 'o' or VNI '6' (circumflex ^ on a/e/o)
	if (useTelex && (lc == 'a' || lc == 'e' || lc == 'o')) || (useVNI && lc == '6') {
		for i := dec.vowelEnd - 1; i >= dec.vowelStart; i-- {
			meta := runeToMeta[syl[i]]
			if lc != '6' && meta.ascii != lc {
				continue
			}
			lb := unicode.ToLower(meta.base)
			switch lb {
			case 'a':
				syl[i] = changeVowelShape(syl[i], 'â')
				return rebalanceTone(syl, e.Cfg.ModernToneStyle), true, false
			case 'e':
				syl[i] = changeVowelShape(syl[i], 'ê')
				return rebalanceTone(syl, e.Cfg.ModernToneStyle), true, false
			case 'o':
				syl[i] = changeVowelShape(syl[i], 'ô')
				return rebalanceTone(syl, e.Cfg.ModernToneStyle), true, false
			case 'â', 'ê', 'ô':
				// Double-key undo (e.g. aa + a -> aa, ee + e -> ee, oo + o -> oo)
				syl[i] = changeVowelShape(syl[i], meta.ascii)
				syl = append(syl, ch)
				return syl, true, true
			}
		}
	}

	// 4. Telex 'w' or VNI '7' / '8' (horn ơ/ư and breve ă)
	if (useTelex && lc == 'w') || (useVNI && (lc == '7' || lc == '8')) {
		vASC := lowerASCIIString(dec.vowels)
		vBase := lowerBaseString(dec.vowels)

		// Count how many times 'w'/'7'/'8' was previously typed in rawHistory (excluding the current key `ch` which was just appended)
		prevWCount := 0
		prevWasW := false
		if len(e.rawHistory) >= 2 {
			prevKey := unicode.ToLower(e.rawHistory[len(e.rawHistory)-2])
			prevWasW = (prevKey == 'w' || prevKey == '7' || prevKey == '8')
			for _, rk := range e.rawHistory[:len(e.rawHistory)-1] {
				rl := unicode.ToLower(rk)
				if rl == 'w' || rl == '7' || rl == '8' {
					prevWCount++
				}
			}
		}

		// If vBase is ALREADY "ươ" because the user typed "uw" earlier and auto-normalization turned "ưo" -> "ươ",
		// and the user now presses a second 'w' (not immediately after the first 'w', e.g. dduwocwj or truwongwf),
		// treat this second 'w' as confirming "ươ" (swallow with no change) rather than undoing "ươ"!
		if vBase == "ươ" && prevWCount == 1 && !prevWasW {
			return syl, true, false
		}

		// Undo if already has horn/breve and user pressed 'w' again immediately (or after already finished)
		if vBase == "ă" || vBase == "ơ" || vBase == "ư" || vBase == "ươ" ||
			(prevWasW && (strings.ContainsRune(vBase, 'ă') || strings.ContainsRune(vBase, 'ơ') || strings.ContainsRune(vBase, 'ư'))) {
			for i := dec.vowelStart; i < dec.vowelEnd; i++ {
				meta := runeToMeta[syl[i]]
				syl[i] = changeVowelShape(syl[i], meta.ascii)
			}
			syl = append(syl, ch)
			return syl, true, true
		}

		// Special pair "uo" / "uô" / "ưo" -> "ươ" (e.g., "truong" + "w" -> "trương", "huou" + "w" -> "hươu")
		if strings.Contains(vASC, "uo") && lc != '8' {
			for i := dec.vowelStart; i+1 < dec.vowelEnd; i++ {
				m1 := runeToMeta[syl[i]]
				m2 := runeToMeta[syl[i+1]]
				if m1.ascii == 'u' && m2.ascii == 'o' {
					// Exception: "qu" where 'u' was not decomposed as initial (already handled by decompose)
					syl[i] = changeVowelShape(syl[i], 'ư')
					syl[i+1] = changeVowelShape(syl[i+1], 'ơ')
					return rebalanceTone(syl, e.Cfg.ModernToneStyle), true, false
				}
			}
		}

		// Single vowel modification (right-to-left search for u, o, a)
		for i := dec.vowelEnd - 1; i >= dec.vowelStart; i-- {
			meta := runeToMeta[syl[i]]
			if lc == '7' && meta.ascii == 'a' {
				continue
			}
			if lc == '8' && meta.ascii != 'a' {
				continue
			}
			switch meta.ascii {
			case 'a':
				// In "ua" (like "mua" + "w" -> "mưa"), 'u' takes the horn, not 'a' (unless "qua" + "w" -> "quă" where 'u' is in initial)
				if i > dec.vowelStart && runeToMeta[syl[i-1]].ascii == 'u' && len(dec.final) == 0 {
					syl[i-1] = changeVowelShape(syl[i-1], 'ư')
					return rebalanceTone(syl, e.Cfg.ModernToneStyle), true, false
				}
				syl[i] = changeVowelShape(syl[i], 'ă')
				return rebalanceTone(syl, e.Cfg.ModernToneStyle), true, false
			case 'o':
				syl[i] = changeVowelShape(syl[i], 'ơ')
				return rebalanceTone(syl, e.Cfg.ModernToneStyle), true, false
			case 'u':
				syl[i] = changeVowelShape(syl[i], 'ư')
				return rebalanceTone(syl, e.Cfg.ModernToneStyle), true, false
			}
		}
	}

	//Also rebalance tone if user just added a vowel or final consonant?
	//Wait: what if user typed "ho" + "f" -> "hò", and then typed "a" -> should become "hòa" or "hoà", and then typed "n" -> MUST become "hoàn"!
	if isVietnameseVowel(ch) || (len(dec.vowels) > 0 && validFinals[lowerASCIIString(dec.final)+string(lc)]) {
		candidate := append(append([]rune(nil), oldSyl...), ch)
		rebalanced := rebalanceTone(append([]rune(nil), candidate...), e.Cfg.ModernToneStyle)
		if string(rebalanced) != string(candidate) {
			return rebalanced, true, false
		}
	}

	return nil, false, false
}
