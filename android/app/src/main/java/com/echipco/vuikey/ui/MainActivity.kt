package com.echipco.vuikey.ui

import android.app.Activity
import android.content.Intent
import android.graphics.Color
import android.graphics.Typeface
import android.graphics.drawable.GradientDrawable
import android.net.Uri
import android.os.Bundle
import android.provider.Settings
import android.util.TypedValue
import android.view.Gravity
import android.view.inputmethod.InputMethodManager
import android.widget.Button
import android.widget.CheckBox
import android.widget.EditText
import android.widget.LinearLayout
import android.widget.ScrollView
import android.widget.TextView
import com.echipco.vuikey.engine.InputMode
import com.echipco.vuikey.engine.VuiKeyConfig
import com.echipco.vuikey.engine.VuiKeyPreferences

class MainActivity : Activity() {

    private lateinit var cfg: VuiKeyConfig
    private val modeButtons = mutableListOf<Button>()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        cfg = VuiKeyPreferences.load(this)
        window.statusBarColor = Color.parseColor("#090D16")
        window.navigationBarColor = Color.parseColor("#090D16")

        val scroll = ScrollView(this).apply {
            setBackgroundColor(Color.parseColor("#090D16"))
            isFillViewport = true
        }

        val container = LinearLayout(this).apply {
            orientation = LinearLayout.VERTICAL
            setPadding(dp(20), dp(24), dp(20), dp(36))
        }
        scroll.addView(container)

        // Header
        val badge = TextView(this).apply {
            text = "⚡ ZERO-PREEDIT • 100% OFFLINE & BẢO MẬT"
            setTextColor(Color.parseColor("#34D399"))
            setTextSize(TypedValue.COMPLEX_UNIT_SP, 11.5f)
            typeface = Typeface.DEFAULT_BOLD
        }
        container.addView(badge)

        val title = TextView(this).apply {
            text = "VuiKey Studio 🇻🇳"
            setTextColor(Color.parseColor("#F8FAFC"))
            setTextSize(TypedValue.COMPLEX_UNIT_SP, 26f)
            typeface = Typeface.DEFAULT_BOLD
            setPadding(0, dp(4), 0, dp(4))
        }
        container.addView(title)

        val subtitle = TextView(this).apply {
            text = "Giải quyết dứt điểm những nỗi đau gõ Tiếng Việt tồn tại nhiều năm: Tự hiểu Tiếng Anh không cần tắt bộ gõ, bấm Backspace 1 nhịp sửa dấu ăn ngay, gõ tự do 5 kiểu ươ & đ, chống nhảy đúp chữ 100%."
            setTextColor(Color.parseColor("#94A3B8"))
            setTextSize(TypedValue.COMPLEX_UNIT_SP, 13.5f)
            setPadding(0, 0, 0, dp(16))
        }
        container.addView(subtitle)

        // Card 1: 2-Step Quick Activation
        val setupCard = createCard("#131C2E", "#1E3A8A").apply {
            addView(sectionTitle("🚀 Kích Hoạt Nhanh Trong 2 Bước"))

            val step1Btn = Button(this@MainActivity).apply {
                text = "1. Bật VuiKey Trong Cài Đặt Hệ Thống ➔"
                isAllCaps = false
                setTextColor(Color.parseColor("#042F2E"))
                typeface = Typeface.DEFAULT_BOLD
                background = roundedRect("#10B981", dp(10))
                layoutParams = LinearLayout.LayoutParams(
                    LinearLayout.LayoutParams.MATCH_PARENT,
                    LinearLayout.LayoutParams.WRAP_CONTENT
                ).apply { topMargin = dp(10) }
                setOnClickListener {
                    startActivity(Intent(Settings.ACTION_INPUT_METHOD_SETTINGS))
                }
            }
            addView(step1Btn)

            val step2Btn = Button(this@MainActivity).apply {
                text = "2. Chọn VuiKey Làm Bàn Phím Mặc Định ⌨️"
                isAllCaps = false
                setTextColor(Color.WHITE)
                typeface = Typeface.DEFAULT_BOLD
                background = roundedRect("#0284C7", dp(10))
                layoutParams = LinearLayout.LayoutParams(
                    LinearLayout.LayoutParams.MATCH_PARENT,
                    LinearLayout.LayoutParams.WRAP_CONTENT
                ).apply { topMargin = dp(10) }
                setOnClickListener {
                    val imm = getSystemService(INPUT_METHOD_SERVICE) as? InputMethodManager
                    imm?.showInputMethodPicker()
                }
            }
            addView(step2Btn)
        }
        container.addView(setupCard)

        // Card 2: Typing Sandbox
        val sandboxCard = createCard("#172554", "#2563EB").apply {
            addView(sectionTitle("🎤 Khu Vực Gõ Thử (Typing Sandbox)"))

            val hint = TextView(this@MainActivity).apply {
                text = "Thử gõ: Vieetj Nam, dodongj (➔ động), dduwocj (➔ được), nhacjdiujdangf (➔ nhạcdịudàng), hoặc gõ window, test, server:"
                setTextColor(Color.parseColor("#CBD5E1"))
                setTextSize(TypedValue.COMPLEX_UNIT_SP, 12.5f)
                setPadding(0, dp(6), 0, dp(8))
            }
            addView(hint)

            val input = EditText(this@MainActivity).apply {
                setHint("Chạm vào đây để bật bàn phím VuiKey và gõ thử...")
                setHintTextColor(Color.parseColor("#64748B"))
                setTextColor(Color.parseColor("#38BDF8"))
                setTextSize(TypedValue.COMPLEX_UNIT_SP, 15.5f)
                background = roundedRect("#090D16", dp(8), "#3B82F6")
                setPadding(dp(12), dp(12), dp(12), dp(12))
                minLines = 2
            }
            addView(input)
        }
        container.addView(sandboxCard)

        // Card 3: Engine Settings
        val settingsCard = createCard("#131C2E", "#334155").apply {
            addView(sectionTitle("⚙️ Cấu Hình Bộ Gõ Thông Minh"))

            val modeRow = LinearLayout(this@MainActivity).apply {
                orientation = LinearLayout.HORIZONTAL
                layoutParams = LinearLayout.LayoutParams(
                    LinearLayout.LayoutParams.MATCH_PARENT,
                    LinearLayout.LayoutParams.WRAP_CONTENT
                ).apply {
                    topMargin = dp(10)
                    bottomMargin = dp(8)
                }
            }

            modeButtons.clear()
            for (mode in InputMode.entries) {
                val btn = Button(this@MainActivity).apply {
                    text = mode.displayName
                    isAllCaps = false
                    setTextSize(TypedValue.COMPLEX_UNIT_SP, 13f)
                    typeface = Typeface.DEFAULT_BOLD
                    layoutParams = LinearLayout.LayoutParams(0, dp(40), 1f).apply {
                        marginEnd = dp(6)
                    }
                    setOnClickListener {
                        cfg.mode = mode
                        VuiKeyPreferences.save(this@MainActivity, cfg)
                        refreshModeButtons()
                    }
                }
                modeButtons.add(btn)
                modeRow.addView(btn)
            }
            addView(modeRow)
            refreshModeButtons()

            addView(createToggle("Tự hiểu Tiếng Anh (Smart English Guard)", cfg.smartEnglish) {
                cfg.smartEnglish = it
                VuiKeyPreferences.save(this@MainActivity, cfg)
            })
            addView(createToggle("Tự tách âm tiết viết liền (Sliding Syllable Window)", cfg.syllableSplit) {
                cfg.syllableSplit = it
                VuiKeyPreferences.save(this@MainActivity, cfg)
            })
            addView(createToggle("Đặt dấu kiểu mới (oà, uý thay vì òa, úy)", cfg.modernToneStyle) {
                cfg.modernToneStyle = it
                VuiKeyPreferences.save(this@MainActivity, cfg)
            })
            addView(createToggle("Rung phản hồi khi bấm phím (Haptic Feedback)", cfg.hapticFeedback) {
                cfg.hapticFeedback = it
                VuiKeyPreferences.save(this@MainActivity, cfg)
            })
        }
        container.addView(settingsCard)

        // Footer links
        val footer = TextView(this).apply {
            text = "Phát triển bởi Công ty TNHH eChip (echipco.com) & VuiCaHat.com"
            gravity = Gravity.CENTER
            setTextColor(Color.parseColor("#38BDF8"))
            setTextSize(TypedValue.COMPLEX_UNIT_SP, 12.5f)
            setPadding(0, dp(16), 0, dp(8))
            setOnClickListener {
                startActivity(Intent(Intent.ACTION_VIEW, Uri.parse("https://echipco.com/vuikey/")))
            }
        }
        container.addView(footer)

        setContentView(scroll)
    }

    private fun refreshModeButtons() {
        for ((idx, btn) in modeButtons.withIndex()) {
            val active = (cfg.mode.code == idx)
            btn.setTextColor(if (active) Color.parseColor("#042F2E") else Color.parseColor("#CBD5E1"))
            btn.background = roundedRect(if (active) "#10B981" else "#1E293B", dp(8))
        }
    }

    private fun createCard(bgHex: String, borderHex: String): LinearLayout {
        return LinearLayout(this).apply {
            orientation = LinearLayout.VERTICAL
            background = roundedRect(bgHex, dp(14), borderHex)
            setPadding(dp(16), dp(16), dp(16), dp(16))
            layoutParams = LinearLayout.LayoutParams(
                LinearLayout.LayoutParams.MATCH_PARENT,
                LinearLayout.LayoutParams.WRAP_CONTENT
            ).apply {
                bottomMargin = dp(14)
            }
        }
    }

    private fun sectionTitle(textStr: String): TextView {
        return TextView(this).apply {
            text = textStr
            setTextColor(Color.parseColor("#38BDF8"))
            setTextSize(TypedValue.COMPLEX_UNIT_SP, 15.5f)
            typeface = Typeface.DEFAULT_BOLD
        }
    }

    private fun createToggle(label: String, initial: Boolean, onChanged: (Boolean) -> Unit): CheckBox {
        return CheckBox(this).apply {
            text = label
            isChecked = initial
            setTextColor(Color.parseColor("#E2E8F0"))
            setTextSize(TypedValue.COMPLEX_UNIT_SP, 13.5f)
            setPadding(dp(6), dp(8), dp(6), dp(8))
            setOnCheckedChangeListener { _, checked -> onChanged(checked) }
        }
    }

    private fun roundedRect(fillHex: String, radiusPx: Int, strokeHex: String? = null): GradientDrawable {
        return GradientDrawable().apply {
            shape = GradientDrawable.RECTANGLE
            cornerRadius = radiusPx.toFloat()
            setColor(Color.parseColor(fillHex))
            if (strokeHex != null) {
                setStroke(dp(1), Color.parseColor(strokeHex))
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
