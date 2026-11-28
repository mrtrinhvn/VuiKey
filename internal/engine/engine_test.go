package engine

import (
	"testing"
)

// simulateEditor simulates a real application text box receiving Zero-Preedit events.
func simulateEditor(eng *Engine, keystrokes string) string {
	var screen []rune
	for _, ch := range keystrokes {
		switch ch {
		case '\b': // Hardware Backspace
			if len(screen) > 0 {
				screen = screen[:len(screen)-1]
			}
			eng.HandleBackspace()
		case '\x7f': // Hardware Delete
			eng.HandleDelete()
		case ' ':
			if act, ok := eng.CheckMacroOnSpace(); ok {
				if act.Backspaces <= len(screen) {
					screen = screen[:len(screen)-act.Backspaces]
				}
				screen = append(screen, []rune(act.CommitText)...)
			}
			screen = append(screen, ' ')
		default:
			act := eng.ProcessKey(ch)
			if act.Swallow {
				if act.Backspaces <= len(screen) {
					screen = screen[:len(screen)-act.Backspaces]
				}
				screen = append(screen, []rune(act.CommitText)...)
			} else {
				screen = append(screen, ch)
			}
		}
	}
	return string(screen)
}

func TestBugReproduction_NhacDiDang_ZeroDuplicateOnDelete(t *testing.T) {
	eng := NewEngine(DefaultConfig())
	input := "nhacdidangnhactiengdidangtiengchoitai\x7f"
	got := simulateEditor(eng, input)
	want := "nhacdidangnhactiengdidangtiengchoitai"
	if got != want {
		t.Fatalf("Expected %q, got %q (MUST NEVER DUPLICATE!)", want, got)
	}
	if eng.CurrentSyllable() != "" {
		t.Fatalf("Expected empty syllable after Delete, got %q", eng.CurrentSyllable())
	}
}

func TestSlidingSyllableWindow_ConcatenatedWithTones(t *testing.T) {
	eng := NewEngine(DefaultConfig())
	input := "nhacjdiujdangfnhacjtieengschoitai"
	got := simulateEditor(eng, input)
	want := "nhạcdịudàngnhạctiếngchoitai"
	if got != want {
		t.Fatalf("Concatenated syllables with tones: got %q, want %q", got, want)
	}
}

func TestComprehensiveVietnameseTypingOrders(t *testing.T) {
	eng := NewEngine(DefaultConfig())
	cases := []struct {
		in   string
		want string
	}{
		// User reported sentence & variations
		{"Tooi e rawngf conf nhieeuf looix vowis cacs chuwx khacs nuwax, maf thanhf khacs nuwax", "Tôi e rằng còn nhiều lỗi với các chữ khác nữa, mà thành khác nữa"},
		{"maf thanhf khasc nuawx", "mà thành khác nữa"},
		{"maf thanhf khacs nuwxa", "mà thành khác nữa"},
		{"maf thanhf khacs nuwax", "mà thành khác nữa"},
		{"phieenf phuwcs", "phiền phức"},
		{"phieenf phucws", "phiền phức"},
		{"phieenf phucsw", "phiền phức"},
		{"phieenf phwcs", "phiền phức"},
		{"phieenf phwucs", "phiền phức"},
		{"pl\bhuwcs", "phức"},
		{"phuwcf\bcs", "phức"},

		// All 5 ways of typing "được", "trường", "người", "nước", "mượn", "thương", and free-style delayed 'd' ("dodongj" -> "động")
		{"dodongj", "động"},
		{"doongdj", "động"},
		{"doongjd", "động"},
		{"dongodj", "động"},
		{"dongojd", "động"},
		{"duocwjd", "được"},
		{"duowcjd", "được"},
		{"dduowcj", "được"},
		{"dduocwj", "được"},
		{"dduwowcj", "được"},
		{"dduwocwj", "được"},
		{"dduwocj", "được"},
		{"truwowngf", "trường"},
		{"truowngf", "trường"},
		{"truongwf", "trường"},
		{"truwongf", "trường"},
		{"truwongwf", "trường"},
		{"nguwowif", "người"},
		{"nguowif", "người"},
		{"nguoiwf", "người"},
		{"nguwoif", "người"},
		{"nuwowcs", "nước"},
		{"nuowcs", "nước"},
		{"nuocws", "nước"},
		{"nuwocs", "nước"},

		// Tricky "gi" and "qu" combinations
		{"gias", "giá"},
		{"gisa", "giá"},
		{"giaf", "già"},
		{"gifa", "già"},
		{"gif", "gì"},
		{"ginf", "gìn"},
		{"gieengs", "giếng"},
		{"giuwax", "giữa"},
		{"giuawx", "giữa"},
		{"giuwowngf", "giường"},
		{"quaas", "quấ"},
		{"quaan", "quân"},
		{"quaans", "quấn"},
		{"quoocs", "quốc"},
		{"quocs", "quóc"},
		{"quys", "quý"},
		{"quyf", "quỳ"},
		{"quynhf", "quỳnh"},
		{"quyens", "quyén"},
		{"quyeens", "quyến"},

		// Tricky triphthongs & diphthongs
		{"nghiengs", "nghiéng"},
		{"nghieengs", "nghiêngs"}, // wait: nghiêng has tone ngang, nghiếng has sắc!
		{"nghieeng", "nghiêng"},
		{"khuyur", "khuỷu"},
		{"khuysu", "khuýu"},
		{"ngoafn", "ngoàn"},
		{"ngoawfn", "ngoằn"},
		{"ngoawnf", "ngoằn"},
		{"ngoanwf", "ngoằn"},
		{"hoawcj", "hoặc"},
		{"hoacwj", "hoặc"},
		{"xoawns", "xoắn"},
		{"xoanws", "xoắn"},
		{"chuaanr", "chuẩn"},
		{"chuanar", "chuẩn"},
		{"luajat", "luật"},
		{"luaatj", "luật"},
		{"luataj", "luật"},
		{"chuyeenj", "chuyện"},
		{"chuyenje", "chuyện"},
		{"tuyeets", "tuyết"},
		{"huynhf", "huỳnh"},
		{"huyfnh", "huỳnh"},
		{"tuooir", "tuổi"},
		{"tuoir", "tuỏi"},
		{"tuoior", "tuổi"},
		{"huwowu", "hươu"},
		{"huouw", "hươu"},
		{"thuowr", "thuở"},
		{"thuorw", "thuở"},
		{"thuwor", "thuở"},
		{"thuwowr", "thuở"},
		{"Thuowr xuwa", "Thuở xưa"},
		{"huow tay", "huơ tay"},

		// Capitalized & ALL CAPS Vietnamese
		{"Vieetj Nam", "Việt Nam"},
		{"VIEETJ NAM", "VIỆT NAM"},
		{"DDuwowngf", "Đường"},
		{"DDUWOWNGF", "ĐƯỜNG"},
		{"OWr ddaay", "Ở đây"},
		{"Aans tuwowngj", "Ấn tượng"},
	}

	for _, tc := range cases {
		if tc.in == "nghieengs" {
			tc.want = "nghiếng"
		}
		eng.Reset()
		got := simulateEditor(eng, tc.in)
		if got != tc.want {
			t.Errorf("Input %q:\n  got:  %q\n  want: %q", tc.in, got, tc.want)
		}
	}
}

func TestSmartEnglishAndDoubleKeyUndo(t *testing.T) {
	eng := NewEngine(DefaultConfig())
	cases := []struct {
		in   string
		want string
	}{
		{"window", "window"},
		{"view", "view"},
		{"function", "function"},
		{"stream", "stream"},
		{"class", "class"},
		{"show", "show"},
		{"know", "know"},
		{"throw", "throw"},
		{"arrow", "arrow"},
		{"follow", "follow"},
		{"global", "global"},
		{"ddda", "dda"},
		{"tooi", "tôi"},
		{"toooi", "tooi"},
	}

	for _, tc := range cases {
		eng.Reset()
		got := simulateEditor(eng, tc.in)
		if got != tc.want {
			t.Errorf("Input %q:\n  got:  %q\n  want: %q", tc.in, got, tc.want)
		}
	}
}

func TestVNIModeAndCombinedMode(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Mode = ModeTelexVNI
	eng := NewEngine(cfg)

	cases := []struct {
		in   string
		want string
	}{
		{"Vie65t Nam", "Việt Nam"},
		{"Co65ng hoa2 xa4 ho65i chu3 nghi4a Vie65t Nam", "Cộng hòa xã hội chủ nghĩa Việt Nam"},
		{"d9u7o7ng2 tru7o7ng2", "đường trường"},
		{"a8n co7m", "ăn cơm"},
	}
	for _, tc := range cases {
		eng.Reset()
		got := simulateEditor(eng, tc.in)
		if got != tc.want {
			t.Errorf("VNI Input %q:\n  got:  %q\n  want: %q", tc.in, got, tc.want)
		}
	}
}

func TestTryModifySurroundingWordAndNoGhostText(t *testing.T) {
	eng := NewEngine(DefaultConfig())

	applyAtCursor := func(before string, keys string) string {
		eng.Reset()
		screen := []rune(before)
		for _, ch := range keys {
			if act, ok := eng.TryModifySurroundingWord(string(screen), ch); ok && act.Swallow {
				screen = screen[:len(screen)-act.Backspaces]
				screen = append(screen, []rune(act.CommitText)...)
			} else {
				screen = append(screen, ch)
			}
		}
		return string(screen)
	}

	cases := []struct {
		before string
		keys   string
		want   string
	}{
		{"Xin chào tôi", "s", "Xin chào tối"},
		{"Xin chào tối", "f", "Xin chào tồi"},
		{"Xin chào tối", "z", "Xin chào tôi"},
		{"Xin chào tối", "zz", "Xin chào toi"},
		{"công tuong", "ws", "công tướng"},
		{"làm được", "z", "làm đươc"},
		{"làm được", "zz", "làm duoc"},
		{"chủ dong", "doj", "chủ động"},
		{"gí", "a", "giá"},
		{"hò", "a", "hòa"},
		{"run test", "s", "run tests"},
		{"active user", "s", "active users"},
		{"main server", "s", "main servers"},
		{"nghe nhạc", "k", "nghe nhạck"},
	}

	for _, tc := range cases {
		got := applyAtCursor(tc.before, tc.keys)
		if got != tc.want {
			t.Errorf("Cursor after %q + keys %q: got %q, want %q", tc.before, tc.keys, got, tc.want)
		}
	}

	// Regression test: Backspacing an uncommitted preedit syllable in IME mode must NEVER eat committedLine
	eng.Reset()
	eng.ClearCommittedLine()
	for _, ch := range "Xin chao " {
		eng.ProcessKeyForIME(ch)
	}
	if got := eng.CommittedTextBeforeCursor(); got != "" {
		// reeditActive is false during normal typing after space
		t.Errorf("Expected empty CommittedTextBeforeCursor after trailing space, got %q", got)
	}
	// Type "ban" in preedit (no space), then backspace 3 times via HandleBackspace()
	for _, ch := range "ban" {
		eng.ProcessKeyForIME(ch)
	}
	eng.HandleBackspace()
	eng.HandleBackspace()
	eng.HandleBackspace()
	if eng.CurrentSyllable() != "" {
		t.Fatalf("Expected empty syllable after 3 backspaces, got %q", eng.CurrentSyllable())
	}
	// committedLine must STILL be "Xin chao " (with the trailing space intact!)
	if string(eng.committedLine) != "Xin chao " || eng.committedCursor != 9 {
		t.Fatalf("committedLine corrupted by preedit backspace: got %q (cursor=%d), want %q (cursor=9)",
			string(eng.committedLine), eng.committedCursor, "Xin chao ")
	}
	// And because the cursor is after a trailing space, reeditActive must be false!
	if eng.CommittedTextBeforeCursor() != "" {
		t.Fatalf("Expected empty CommittedTextBeforeCursor after preedit backspace, got %q", eng.CommittedTextBeforeCursor())
	}
	// Now press Left arrow once (or Backspace once) so cursor touches the end of "chao": reeditActive must become true!
	eng.HandleNavigation(1) // Left arrow across the trailing space
	if got := eng.CommittedTextBeforeCursor(); got != "Xin chao" {
		t.Fatalf("Expected CommittedTextBeforeCursor='Xin chao' after Left arrow, got %q", got)
	}

	// Test Mouse Click cursor snapping via OnClientReset + UpdateCursorRect:
	// Type "tooi ddi hocj" with simulated X coordinates (10px per character starting at x=100, y=500, h=20)
	eng.Reset()
	eng.ClearCommittedLine()
	eng.UpdateCursorRect(100, 500, 20) // x0 = 100 at index 0
	for _, ch := range "tooi " {
		eng.ProcessKeyForIME(ch)
	}
	eng.UpdateCursorRect(140, 500, 20) // x = 140 at index 4 (after "tôi ")
	for _, ch := range "ddi " {
		eng.ProcessKeyForIME(ch)
	}
	eng.UpdateCursorRect(170, 500, 20) // x = 170 at index 7 (after "tôi đi ")
	for _, ch := range "hocj" {
		eng.ProcessKeyForIME(ch)
	}
	// Simulate Mouse Click at the end of "tôi" (index 3, x ≈ 130):
	eng.OnClientReset()
	eng.UpdateCursorRect(131, 500, 20)
	if got := eng.CommittedTextBeforeCursor(); got != "tôi" {
		t.Fatalf("Expected mouse click at x=131 to snap to 'tôi', got %q", got)
	}
	// Simulate Mouse Click at the end of "học" (index 10, x ≈ 200):
	eng.OnClientReset()
	eng.UpdateCursorRect(199, 500, 20)
	if got := eng.CommittedTextBeforeCursor(); got != "tôi đi học" {
		t.Fatalf("Expected mouse click at x=199 to snap to 'tôi đi học', got %q", got)
	}
	// Simulate Mouse Click/Clear at start of input box (x = 100): must NOT return any word before cursor!
	eng.OnClientReset()
	eng.UpdateCursorRect(100, 500, 20)
	if got := eng.CommittedTextBeforeCursor(); got != "" {
		t.Fatalf("Expected mouse click at x=100 (start of line) to return empty, got %q", got)
	}
}

