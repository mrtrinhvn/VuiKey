package com.echipco.vuikey.engine

import android.content.Context

object VuiKeyPreferences {
    private const val PREFS_NAME = "vuikey_prefs"
    private const val KEY_ENABLED = "enabled"
    private const val KEY_MODE = "mode"
    private const val KEY_SMART_ENGLISH = "smart_english"
    private const val KEY_SYLLABLE_SPLIT = "syllable_split"
    private const val KEY_MODERN_TONE = "modern_tone"
    private const val KEY_HAPTIC = "haptic"

    fun load(context: Context): VuiKeyConfig {
        val prefs = context.getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE)
        return VuiKeyConfig(
            enabled = prefs.getBoolean(KEY_ENABLED, true),
            mode = InputMode.fromCode(prefs.getInt(KEY_MODE, InputMode.TELEX.code)),
            smartEnglish = prefs.getBoolean(KEY_SMART_ENGLISH, true),
            syllableSplit = prefs.getBoolean(KEY_SYLLABLE_SPLIT, true),
            modernToneStyle = prefs.getBoolean(KEY_MODERN_TONE, false),
            hapticFeedback = prefs.getBoolean(KEY_HAPTIC, true)
        )
    }

    fun save(context: Context, cfg: VuiKeyConfig) {
        context.getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE)
            .edit()
            .putBoolean(KEY_ENABLED, cfg.enabled)
            .putInt(KEY_MODE, cfg.mode.code)
            .putBoolean(KEY_SMART_ENGLISH, cfg.smartEnglish)
            .putBoolean(KEY_SYLLABLE_SPLIT, cfg.syllableSplit)
            .putBoolean(KEY_MODERN_TONE, cfg.modernToneStyle)
            .putBoolean(KEY_HAPTIC, cfg.hapticFeedback)
            .apply()
    }
}
