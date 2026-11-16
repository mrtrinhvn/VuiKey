package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"sync"

	"github.com/mrtrinhvn/VuiKey/internal/config"
	"github.com/mrtrinhvn/VuiKey/internal/engine"
)

var (
	mu      sync.Mutex
	globalE = engine.NewEngine(config.Load())
)

//export VuiCoreReloadConfig
func VuiCoreReloadConfig() {
	mu.Lock()
	defer mu.Unlock()
	cfg := config.Load()
	cfg.Enabled = true
	globalE.Cfg = cfg
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

func main() {}
