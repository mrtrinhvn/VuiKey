package com.echipco.vuikey.engine

import org.junit.Assert.assertEquals
import org.junit.Test

class VuiKeyEngineTest {

    private fun simulateEditor(eng: VuiKeyEngine, keystrokes: String): String {
        val screen = StringBuilder()
        for (ch in keystrokes) {
            when (ch) {
                '\b' -> {
                    if (screen.isNotEmpty()) {
                        screen.deleteCharAt(screen.length - 1)
                    }
                    eng.handleBackspace()
                }
                '\u007F' -> {
                    eng.handleDelete()
                }
                ' ' -> {
                    val act = eng.checkMacroOnSpace()
                    if (act != null) {
                        if (act.backspaces <= screen.length) {
                            screen.delete(screen.length - act.backspaces, screen.length)
                        }
                        screen.append(act.commitText)
                    }
                    screen.append(' ')
                }
                else -> {
                    val act = eng.processKey(ch)
                    if (act.swallow) {
                        if (act.backspaces <= screen.length) {
                            screen.delete(screen.length - act.backspaces, screen.length)
                        }
                        screen.append(act.commitText)
                    } else {
                        screen.append(ch)
                    }
                }
            }
        }
        return screen.toString()
    }

    @Test
    fun testNhacDiDangZeroDuplicateOnDelete() {
        val eng = VuiKeyEngine()
        val input = "nhacdidangnhactiengdidangtiengchoitai\u007F"
        val got = simulateEditor(eng, input)
        assertEquals("nhacdidangnhactiengdidangtiengchoitai", got)
        assertEquals("", eng.currentSyllable)
    }

    @Test
    fun testSlidingSyllableWindowConcatenatedWithTones() {
        val eng = VuiKeyEngine()
        val input = "nhacjdiujdangfnhacjtieengschoitai"
        val got = simulateEditor(eng, input)
        assertEquals("nhạcdịudàngnhạctiếngchoitai", got)
    }

    @Test
    fun testComprehensiveVietnameseTypingOrders() {
        val eng = VuiKeyEngine()
        val cases = listOf(
            "Tooi e rawngf conf nhieeuf looix vowis cacs chuwx khacs nuwax, maf thanhf khacs nuwax" to
                "Tôi e rằng còn nhiều lỗi với các chữ khác nữa, mà thành khác nữa",
            "maf thanhf khasc nuawx" to "mà thành khác nữa",
            "maf thanhf khacs nuwxa" to "mà thành khác nữa",
            "maf thanhf khacs nuwax" to "mà thành khác nữa",
            "phieenf phuwcs" to "phiền phức",
            "phieenf phucws" to "phiền phức",
            "phieenf phucsw" to "phiền phức",
            "phieenf phwcs" to "phiền phức",
            "phieenf phwucs" to "phiền phức",
            "pl\bhuwcs" to "phức",
            "phuwcf\bcs" to "phức",
            "dodongj" to "động",
            "doongdj" to "động",
            "doongjd" to "động",
            "dongodj" to "động",
            "dongojd" to "động",
            "duocwjd" to "được",
            "duowcjd" to "được",
            "dduowcj" to "được",
            "dduocwj" to "được",
            "dduwowcj" to "được",
            "dduwocwj" to "được",
            "dduwocj" to "được",
            "truwowngf" to "trường",
            "truowngf" to "trường",
            "truongwf" to "trường",
            "truwongf" to "trường",
            "truwongwf" to "trường",
            "nguwowif" to "người",
            "nguowif" to "người",
            "nguoiwf" to "người",
            "nguwoif" to "người",
            "nuwowcs" to "nước",
            "nuowcs" to "nước",
            "nuocws" to "nước",
            "nuwocs" to "nước",
            "Vieetj Nam" to "Việt Nam",
            "VIEETJ NAM" to "VIỆT NAM",
            "DDuwowngf" to "Đường",
            "DDUWOWNGF" to "ĐƯỜNG"
        )

        for ((input, want) in cases) {
            eng.reset()
            val got = simulateEditor(eng, input)
            assertEquals("Failed for input: $input", want, got)
        }
    }

    @Test
    fun testSmartEnglishAndDoubleKeyUndo() {
        val eng = VuiKeyEngine()
        val cases = listOf(
            "window" to "window",
            "view" to "view",
            "function" to "function",
            "stream" to "stream",
            "class" to "class",
            "linux" to "linux",
            "test" to "test",
            "ddda" to "dda",
            "tooi" to "tôi",
            "toooi" to "tooi"
        )
        for ((input, want) in cases) {
            eng.reset()
            val got = simulateEditor(eng, input)
            assertEquals("Failed for input: $input", want, got)
        }
    }
}
