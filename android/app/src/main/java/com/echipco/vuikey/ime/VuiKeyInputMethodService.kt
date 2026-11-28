package com.echipco.vuikey.ime

import android.content.Intent
import android.graphics.Color
import android.graphics.Typeface
import android.graphics.drawable.GradientDrawable
import android.inputmethodservice.InputMethodService
import android.os.Handler
import android.os.Looper
import android.util.TypedValue
import android.view.Gravity
import android.view.HapticFeedbackConstants
import android.view.KeyEvent
import android.view.MotionEvent
import android.view.View
import android.view.inputmethod.EditorInfo
import android.view.inputmethod.InputConnection
import android.view.inputmethod.InputMethodManager
import android.widget.Button
import android.widget.LinearLayout
import android.widget.TextView
import com.echipco.vuikey.engine.InputMode
import com.echipco.vuikey.engine.VuiKeyEngine
import com.echipco.vuikey.engine.VuiKeyPreferences
import com.echipco.vuikey.ui.MainActivity
import kotlin.math.abs

class VuiKeyInputMethodService : InputMethodService() {

    private val engine = VuiKeyEngine()
    private val handler = Handler(Looper.getMainLooper())

    private var isShifted = true
    private var isCapsLocked = false
    private var isSymbolLayer = false
    private var lastShiftTapMs = 0L

    private var rootLayout: LinearLayout? = null
    private var rowsContainer: LinearLayout? = null
    private var langButton: Button? = null
    private var modeButton: Button? = null
    private var statusPreview: TextView? = null

    private val letterRows = listOf(
        listOf("q", "w", "e", "r", "t", "y", "u", "i", "o", "p"),
        listOf("a", "s", "d", "f", "g", "h", "j", "k", "l"),
        listOf("SHIFT", "z", "x", "c", "v", "b", "n", "m", "BACKSPACE"),
        listOf("123", "GLOBE", ",", "SPACE", ".", "ENTER")
    )

    private val numberRows = listOf(
        listOf("1", "2", "3", "4", "5", "6", "7", "8", "9", "0"),
        listOf("@", "#", "$", "%", "&", "-", "+", "(", ")", "/"),
        listOf("ABC", "*", "\"", "'", ":", ";", "!", "?", "BACKSPACE"),
        listOf("ABC", "GLOBE", ",", "SPACE", ".", "ENTER")
    )

    override fun onCreate() {
        super.onCreate()
        engine.config = VuiKeyPreferences.load(this)
    }

    override fun onCreateInputView(): View {
        engine.config = VuiKeyPreferences.load(this)

        val main = LinearLayout(this).apply {
            orientation = LinearLayout.VERTICAL
            setBackgroundColor(Color.parseColor("#0F172A"))
            setPadding(dp(4), dp(6), dp(4), dp(8))
        }
        rootLayout = main

        // Top Smart Bar: [🇻🇳 Tiếng Việt / 🇬🇧 English] | [Current Syllable] | [⚙︎ Telex] | [🛠]
        val topBar = LinearLayout(this).apply {
            orientation = LinearLayout.HORIZONTAL
            gravity = Gravity.CENTER_VERTICAL
            setPadding(dp(6), dp(2), dp(6), dp(6))
        }

        val langBtn = Button(this).apply {
            isAllCaps = false
            setTextSize(TypedValue.COMPLEX_UNIT_SP, 12.5f)
            setTextColor(Color.WHITE)
            typeface = Typeface.DEFAULT_BOLD
            background = roundedRect("#10B981", dp(16))
            setPadding(dp(12), dp(2), dp(12), dp(2))
            minimumHeight = dp(30)
            minHeight = dp(30)
            setOnClickListener {
                vibrateKey(it)
                engine.config.enabled = !engine.config.enabled
                engine.reset()
                VuiKeyPreferences.save(this@VuiKeyInputMethodService, engine.config)
                updateTopBar()
                rebuildKeyboardRows()
            }
        }
        langButton = langBtn

        val preview = TextView(this).apply {
            gravity = Gravity.CENTER
            setTextSize(TypedValue.COMPLEX_UNIT_SP, 13.5f)
            setTextColor(Color.parseColor("#94A3B8"))
            typeface = Typeface.DEFAULT_BOLD
            layoutParams = LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1f)
        }
        statusPreview = preview

        val modeBtn = Button(this).apply {
            isAllCaps = false
            setTextSize(TypedValue.COMPLEX_UNIT_SP, 12f)
            setTextColor(Color.parseColor("#38BDF8"))
            background = roundedRect("#1E293B", dp(14), strokeColor = "#334155")
            setPadding(dp(10), dp(2), dp(10), dp(2))
            minimumHeight = dp(30)
            minHeight = dp(30)
            setOnClickListener {
                vibrateKey(it)
                val nextCode = (engine.config.mode.code + 1) % InputMode.entries.size
                engine.config.mode = InputMode.fromCode(nextCode)
                engine.reset()
                VuiKeyPreferences.save(this@VuiKeyInputMethodService, engine.config)
                updateTopBar()
            }
        }
        modeButton = modeBtn

        val settingsBtn = Button(this).apply {
            text = "⚙"
            isAllCaps = false
            setTextSize(TypedValue.COMPLEX_UNIT_SP, 14f)
            setTextColor(Color.parseColor("#CBD5E1"))
            background = roundedRect("#1E293B", dp(14))
            val lp = LinearLayout.LayoutParams(dp(34), dp(30)).apply {
                marginStart = dp(6)
            }
            layoutParams = lp
            setPadding(0, 0, 0, 0)
            setOnClickListener {
                vibrateKey(it)
                val intent = Intent(this@VuiKeyInputMethodService, MainActivity::class.java).apply {
                    addFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
                }
                startActivity(intent)
            }
        }

        topBar.addView(langBtn)
        topBar.addView(preview)
        topBar.addView(modeBtn)
        topBar.addView(settingsBtn)
        main.addView(topBar)

        val rows = LinearLayout(this).apply {
            orientation = LinearLayout.VERTICAL
        }
        rowsContainer = rows
        main.addView(rows)

        rebuildKeyboardRows()
        updateTopBar()
        return main
    }

    private fun isPasswordInputField(): Boolean {
        val info = currentInputEditorInfo ?: return false
        val inputType = info.inputType
        val klass = inputType and android.text.InputType.TYPE_MASK_CLASS
        val variation = inputType and android.text.InputType.TYPE_MASK_VARIATION
        if (klass == android.text.InputType.TYPE_CLASS_TEXT) {
            return variation == android.text.InputType.TYPE_TEXT_VARIATION_PASSWORD ||
                variation == android.text.InputType.TYPE_TEXT_VARIATION_VISIBLE_PASSWORD ||
                variation == android.text.InputType.TYPE_TEXT_VARIATION_WEB_PASSWORD
        }
        if (klass == android.text.InputType.TYPE_CLASS_NUMBER) {
            return variation == android.text.InputType.TYPE_NUMBER_VARIATION_PASSWORD
        }
        return false
    }

    override fun onStartInputView(info: EditorInfo?, restarting: Boolean) {
        super.onStartInputView(info, restarting)
        engine.config = VuiKeyPreferences.load(this)
        engine.reset()
        updateTopBar()
        rebuildKeyboardRows()
    }

    override fun onUpdateSelection(
        oldSelStart: Int,
        oldSelEnd: Int,
        newSelStart: Int,
        newSelEnd: Int,
        candidatesStart: Int,
        candidatesEnd: Int
    ) {
        super.onUpdateSelection(oldSelStart, oldSelEnd, newSelStart, newSelEnd, candidatesStart, candidatesEnd)
        // If the user tapped elsewhere in the text box to move the cursor, reset the active syllable
        if (oldSelStart != newSelStart && abs(newSelStart - oldSelStart) > 2 && engine.currentSyllable.isNotEmpty()) {
            engine.reset()
            updateTopBar()
        }
    }

    // Full Hardware Keyboard support (Bluetooth / USB keyboards on Android phones, tablets & DeX)
    override fun onKeyDown(keyCode: Int, event: KeyEvent?): Boolean {
        if (event == null) return super.onKeyDown(keyCode, event)
        val ic = currentInputConnection ?: return super.onKeyDown(keyCode, event)

        // Ctrl + Space toggles Vietnamese / English on physical keyboards
        if (keyCode == KeyEvent.KEYCODE_SPACE && event.isCtrlPressed) {
            engine.config.enabled = !engine.config.enabled
            engine.reset()
            VuiKeyPreferences.save(this, engine.config)
            updateTopBar()
            rebuildKeyboardRows()
            return true
        }

        if (event.isCtrlPressed || event.isAltPressed || event.isMetaPressed) {
            engine.reset()
            return super.onKeyDown(keyCode, event)
        }

        if (keyCode == KeyEvent.KEYCODE_DEL) {
            engine.handleBackspace()
            updateTopBar()
            return super.onKeyDown(keyCode, event)
        }

        if (keyCode == KeyEvent.KEYCODE_FORWARD_DEL) {
            engine.handleDelete()
            updateTopBar()
            return super.onKeyDown(keyCode, event)
        }

        if (keyCode == KeyEvent.KEYCODE_SPACE) {
            val macroAct = engine.checkMacroOnSpace()
            if (macroAct != null) {
                ic.beginBatchEdit()
                if (macroAct.backspaces > 0) {
                    ic.deleteSurroundingText(macroAct.backspaces, 0)
                }
                ic.commitText(macroAct.commitText + " ", 1)
                ic.endBatchEdit()
                updateTopBar()
                return true
            }
            updateTopBar()
            return super.onKeyDown(keyCode, event)
        }

        val uni = event.unicodeChar
        if (uni > 0) {
            val ch = uni.toChar()
            if (isPasswordInputField()) {
                engine.reset()
                return super.onKeyDown(keyCode, event)
            }
            val act = engine.processKey(ch)
            if (act.swallow) {
                ic.beginBatchEdit()
                if (act.backspaces > 0) {
                    ic.deleteSurroundingText(act.backspaces, 0)
                }
                ic.commitText(act.commitText, 1)
                ic.endBatchEdit()
                updateTopBar()
                return true
            }
            updateTopBar()
        }
        return super.onKeyDown(keyCode, event)
    }

    private fun updateTopBar() {
        val enabled = engine.config.enabled
        langButton?.text = if (enabled) "🇻🇳 Tiếng Việt" else "🇬🇧 English"
        langButton?.background = roundedRect(if (enabled) "#10B981" else "#475569", dp(16))
        modeButton?.text = "Kiểu: ${engine.config.mode.displayName}"
        val syl = engine.currentSyllable
        statusPreview?.text = if (syl.isEmpty()) "VuiKey • Zero-Preedit" else syl
        statusPreview?.setTextColor(
            if (syl.isEmpty()) Color.parseColor("#64748B") else Color.parseColor("#38BDF8")
        )
    }

    private fun rebuildKeyboardRows() {
        val container = rowsContainer ?: return
        container.removeAllViews()

        val activeRows = if (isSymbolLayer) numberRows else letterRows
        for ((rowIdx, row) in activeRows.withIndex()) {
            val rowLayout = LinearLayout(this).apply {
                orientation = LinearLayout.HORIZONTAL
                gravity = Gravity.CENTER
                val lp = LinearLayout.LayoutParams(
                    LinearLayout.LayoutParams.MATCH_PARENT,
                    dp(48)
                ).apply {
                    topMargin = dp(6)
                }
                layoutParams = lp
                if (rowIdx == 1 && !isSymbolLayer) {
                    setPadding(dp(16), 0, dp(16), 0)
                }
            }

            for (key in row) {
                rowLayout.addView(createKeyView(key))
            }
            container.addView(rowLayout)
        }
    }

    private fun createKeyView(key: String): View {
        val weight = when (key) {
            "SPACE" -> 4.2f
            "SHIFT", "BACKSPACE", "123", "ABC", "ENTER" -> 1.45f
            "GLOBE" -> 1.1f
            else -> 1.0f
        }

        val isSpecial = key in setOf("SHIFT", "BACKSPACE", "123", "ABC", "GLOBE", "ENTER")
        val label = when (key) {
            "SHIFT" -> if (isCapsLocked) "⇪" else "⇧"
            "BACKSPACE" -> "⌫"
            "GLOBE" -> "🌐"
            "SPACE" -> if (engine.config.enabled) "Tiếng Việt • VuiKey" else "English"
            "ENTER" -> "↵"
            "123", "ABC" -> key
            else -> if (!isSymbolLayer && (isShifted || isCapsLocked)) key.uppercase() else key
        }

        val bgColor = when {
            key == "ENTER" -> "#0284C7"
            key == "SHIFT" && (isShifted || isCapsLocked) -> "#3B82F6"
            isSpecial -> "#1E293B"
            else -> "#334155"
        }

        val btn = Button(this).apply {
            text = label
            isAllCaps = false
            setTextColor(Color.parseColor("#F8FAFC"))
            setTextSize(TypedValue.COMPLEX_UNIT_SP, if (isSpecial || key == "SPACE") 14f else 19f)
            typeface = if (isSpecial) Typeface.DEFAULT_BOLD else Typeface.DEFAULT
            background = roundedRect(bgColor, dp(8))
            setPadding(0, 0, 0, 0)
            layoutParams = LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.MATCH_PARENT, weight).apply {
                marginStart = dp(3)
                marginEnd = dp(3)
            }
        }

        when (key) {
            "BACKSPACE" -> setupBackspaceRepeat(btn)
            "SPACE" -> setupSpacebarWithTrackpad(btn)
            "GLOBE" -> btn.setOnClickListener {
                vibrateKey(it)
                val imm = getSystemService(INPUT_METHOD_SERVICE) as? InputMethodManager
                imm?.showInputMethodPicker()
            }
            else -> btn.setOnClickListener {
                vibrateKey(it)
                handleVirtualKeyTap(key)
            }
        }

        return btn
    }

    private fun setupBackspaceRepeat(btn: Button) {
        val repeatRunnable = object : Runnable {
            override fun run() {
                performBackspace()
                handler.postDelayed(this, 65L)
            }
        }
        btn.setOnTouchListener { v, event ->
            when (event.actionMasked) {
                MotionEvent.ACTION_DOWN -> {
                    vibrateKey(v)
                    performBackspace()
                    handler.postDelayed(repeatRunnable, 380L)
                    true
                }
                MotionEvent.ACTION_UP, MotionEvent.ACTION_CANCEL -> {
                    handler.removeCallbacks(repeatRunnable)
                    true
                }
                else -> false
            }
        }
    }

    // Spacebar doubles as a smooth cursor trackpad when swiping left/right
    private fun setupSpacebarWithTrackpad(btn: Button) {
        var startX = 0f
        var lastSteppedX = 0f
        var dragged = false

        btn.setOnTouchListener { v, event ->
            when (event.actionMasked) {
                MotionEvent.ACTION_DOWN -> {
                    startX = event.rawX
                    lastSteppedX = startX
                    dragged = false
                    vibrateKey(v)
                    true
                }
                MotionEvent.ACTION_MOVE -> {
                    val stepPx = dp(11).toFloat()
                    val diff = event.rawX - lastSteppedX
                    if (abs(event.rawX - startX) > stepPx) {
                        dragged = true
                        engine.reset()
                    }
                    if (dragged && abs(diff) >= stepPx) {
                        val steps = (diff / stepPx).toInt()
                        val ic = currentInputConnection
                        if (ic != null && steps != 0) {
                            val keyCode = if (steps > 0) KeyEvent.KEYCODE_DPAD_RIGHT else KeyEvent.KEYCODE_DPAD_LEFT
                            repeat(abs(steps)) {
                                ic.sendKeyEvent(KeyEvent(KeyEvent.ACTION_DOWN, keyCode))
                                ic.sendKeyEvent(KeyEvent(KeyEvent.ACTION_UP, keyCode))
                            }
                            vibrateKey(v)
                            lastSteppedX += steps * stepPx
                        }
                    }
                    true
                }
                MotionEvent.ACTION_UP -> {
                    if (!dragged) {
                        handleVirtualKeyTap("SPACE")
                    }
                    updateTopBar()
                    true
                }
                else -> false
            }
        }
    }

    private fun performBackspace() {
        val ic = currentInputConnection ?: return
        val selected = ic.getSelectedText(0)
        if (!selected.isNullOrEmpty()) {
            ic.commitText("", 1)
            engine.reset()
        } else {
            ic.deleteSurroundingText(1, 0)
            engine.handleBackspace()
        }
        updateTopBar()
    }

    private fun handleVirtualKeyTap(key: String) {
        val ic: InputConnection = currentInputConnection ?: return

        when (key) {
            "SHIFT" -> {
                val now = System.currentTimeMillis()
                if (now - lastShiftTapMs < 320L) {
                    isCapsLocked = true
                    isShifted = true
                } else {
                    isCapsLocked = false
                    isShifted = !isShifted
                }
                lastShiftTapMs = now
                rebuildKeyboardRows()
                return
            }
            "123" -> {
                isSymbolLayer = true
                rebuildKeyboardRows()
                return
            }
            "ABC" -> {
                isSymbolLayer = false
                rebuildKeyboardRows()
                return
            }
            "SPACE" -> {
                val macroAct = engine.checkMacroOnSpace()
                if (macroAct != null) {
                    ic.beginBatchEdit()
                    if (macroAct.backspaces > 0) {
                        ic.deleteSurroundingText(macroAct.backspaces, 0)
                    }
                    ic.commitText(macroAct.commitText + " ", 1)
                    ic.endBatchEdit()
                } else {
                    ic.commitText(" ", 1)
                }
                updateTopBar()
                return
            }
            "ENTER" -> {
                engine.reset()
                val imeOptions = currentInputEditorInfo?.imeOptions ?: 0
                val actionId = imeOptions and EditorInfo.IME_MASK_ACTION
                if (actionId != EditorInfo.IME_ACTION_NONE &&
                    actionId != EditorInfo.IME_ACTION_UNSPECIFIED &&
                    (imeOptions and EditorInfo.IME_FLAG_NO_ENTER_ACTION) == 0
                ) {
                    ic.performEditorAction(actionId)
                } else {
                    ic.commitText("\n", 1)
                }
                isShifted = true
                rebuildKeyboardRows()
                updateTopBar()
                return
            }
        }

        val chStr = if (!isSymbolLayer && (isShifted || isCapsLocked)) key.uppercase() else key.lowercase()
        val ch = chStr.firstOrNull() ?: return

        if (isPasswordInputField()) {
            engine.reset()
            ic.commitText(ch.toString(), 1)
            if (isShifted && !isCapsLocked) {
                isShifted = false
                rebuildKeyboardRows()
            }
            updateTopBar()
            return
        }

        val act = engine.processKey(ch)
        if (act.swallow) {
            // Zero-Preedit Atomic Batch Edit: never uses setComposingText, so zero underline & zero duplicate bugs!
            ic.beginBatchEdit()
            if (act.backspaces > 0) {
                ic.deleteSurroundingText(act.backspaces, 0)
            }
            ic.commitText(act.commitText, 1)
            ic.endBatchEdit()
        } else {
            ic.commitText(ch.toString(), 1)
        }

        if (isShifted && !isCapsLocked) {
            isShifted = false
            rebuildKeyboardRows()
        }
        updateTopBar()
    }

    private fun vibrateKey(view: View) {
        if (engine.config.hapticFeedback) {
            view.performHapticFeedback(
                HapticFeedbackConstants.KEYBOARD_TAP,
                HapticFeedbackConstants.FLAG_IGNORE_GLOBAL_SETTING
            )
        }
    }

    private fun roundedRect(fillHex: String, radiusPx: Int, strokeColor: String? = null): GradientDrawable {
        return GradientDrawable().apply {
            shape = GradientDrawable.RECTANGLE
            cornerRadius = radiusPx.toFloat()
            setColor(Color.parseColor(fillHex))
            if (strokeColor != null) {
                setStroke(dp(1), Color.parseColor(strokeColor))
            }
        }
    }

    private fun dp(value: Int): Int {
        return TypedValue.applyDimension(
            TypedValue.COMPLEX_UNIT_DIP,
            value.toFloat(),
            resources.displayMetrics
        ).toInt()
    }
}
