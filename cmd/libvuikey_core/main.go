package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"os"
	"sync"
	"time"

	"github.com/mrtrinhvn/VuiKey/internal/config"
	"github.com/mrtrinhvn/VuiKey/internal/engine"
)

var (
	mu          sync.Mutex
	globalE     = engine.NewEngine(config.Load())
	lastModTime time.Time
)

func syncConfigLocked() {
	info, err := os.Stat(config.ConfigPath())
	if err == nil && !info.ModTime().Equal(lastModTime) {
		lastModTime = info.ModTime()
		cfg := config.Load()
		globalE.Cfg = cfg
	}
}

//export VuiCoreToggleEnabled
func VuiCoreToggleEnabled() C.int {
	mu.Lock()
	defer mu.Unlock()
	syncConfigLocked()
	globalE.Cfg.Enabled = !globalE.Cfg.Enabled
	globalE.Reset()
	_ = config.Save(globalE.Cfg)
	if info, err := os.Stat(config.ConfigPath()); err == nil {
		lastModTime = info.ModTime()
	}
	if globalE.Cfg.Enabled {
		return 1
	}
	return 0
}

//export VuiCoreReloadConfig
func VuiCoreReloadConfig() {
	mu.Lock()
	defer mu.Unlock()
	lastModTime = time.Time{}
	syncConfigLocked()
	globalE.Reset()
}

//export VuiCoreReset
func VuiCoreReset() *C.char {
	mu.Lock()
	defer mu.Unlock()
	syl := globalE.CurrentSyllable()
	globalE.Reset()
	return C.CString(syl)
}

//export VuiCoreCurrentSyllable
func VuiCoreCurrentSyllable() *C.char {
	mu.Lock()
	defer mu.Unlock()
	return C.CString(globalE.CurrentSyllable())
}

//export VuiCoreHandleBackspace
func VuiCoreHandleBackspace(outHandled *C.int) *C.char {
	mu.Lock()
	defer mu.Unlock()
	if len(globalE.CurrentSyllable()) == 0 {
		*outHandled = 0
		return C.CString("")
	}
	globalE.HandleBackspace()
	*outHandled = 1
	return C.CString(globalE.CurrentSyllable())
}

//export VuiCoreProcessKey
func VuiCoreProcessKey(ch C.uint, outCommit **C.char, outPreedit **C.char) C.int {
	mu.Lock()
	defer mu.Unlock()
	syncConfigLocked()
	commitPrefix, curPreedit, handled := globalE.ProcessKeyForIME(rune(ch))
	if !handled {
		*outCommit = nil
		*outPreedit = nil
		return 0
	}
	*outCommit = C.CString(commitPrefix)
	*outPreedit = C.CString(curPreedit)
	return 1
}

//export VuiCoreProcessKeyZeroPreedit
func VuiCoreProcessKeyZeroPreedit(ch C.uint, outBackspaces *C.int, outCommit **C.char) C.int {
	mu.Lock()
	defer mu.Unlock()
	syncConfigLocked()
	act := globalE.ProcessKey(rune(ch))
	if !act.Swallow {
		*outBackspaces = 0
		*outCommit = nil
		return 0
	}
	*outBackspaces = C.int(act.Backspaces)
	*outCommit = C.CString(act.CommitText)
	return 1
}

//export VuiCoreCheckMacroOnSpace
func VuiCoreCheckMacroOnSpace(outBackspaces *C.int, outCommit **C.char) C.int {
	mu.Lock()
	defer mu.Unlock()
	syncConfigLocked()
	act, ok := globalE.CheckMacroOnSpace()
	if !ok {
		*outBackspaces = 0
		*outCommit = nil
		return 0
	}
	*outBackspaces = C.int(act.Backspaces)
	*outCommit = C.CString(act.CommitText)
	return 1
}

//export VuiCoreTryModifySurrounding
func VuiCoreTryModifySurrounding(cText *C.char, cursorCharPos C.uint, ch C.uint, outBackspaces *C.int, outCommit **C.char) C.int {
	mu.Lock()
	defer mu.Unlock()
	syncConfigLocked()
	if cText == nil || len(globalE.CurrentSyllable()) > 0 {
		*outBackspaces = 0
		*outCommit = nil
		return 0
	}
	fullRunes := []rune(C.GoString(cText))
	pos := int(cursorCharPos)
	if pos <= 0 || pos > len(fullRunes) {
		*outBackspaces = 0
		*outCommit = nil
		return 0
	}
	textBefore := string(fullRunes[:pos])
	act, ok := globalE.TryModifySurroundingWord(textBefore, rune(ch))
	if !ok || !act.Swallow {
		*outBackspaces = 0
		*outCommit = nil
		return 0
	}
	*outBackspaces = C.int(act.Backspaces)
	*outCommit = C.CString(act.CommitText)
	return 1
}

func main() {}
