package com.echipco.vuikey.engine

enum class InputMode(val code: Int, val displayName: String) {
    TELEX(0, "Telex"),
    VNI(1, "VNI"),
    TELEX_VNI(2, "Telex + VNI");

    companion object {
        fun fromCode(code: Int): InputMode = entries.find { it.code == code } ?: TELEX
    }
}

enum class Tone(val index: Int) {
    NONE(0),
    SAC(1),   // s / 1
    HUYEN(2), // f / 2
    HOI(3),   // r / 3
    NGA(4),   // x / 4
    NANG(5);  // j / 5

    companion object {
        fun fromIndex(idx: Int): Tone = entries.find { it.index == idx } ?: NONE
    }
}

data class VuiKeyConfig(
    var enabled: Boolean = true,
    var mode: InputMode = InputMode.TELEX,
    var smartEnglish: Boolean = true,
    var syllableSplit: Boolean = true,
    var modernToneStyle: Boolean = false,
    var enableMacros: Boolean = true,
    var hapticFeedback: Boolean = true,
    var macros: MutableMap<String, String> = mutableMapOf("vn" to "Việt Nam")
)

data class VuiKeyAction(
    val swallow: Boolean = false,
    val backspaces: Int = 0,
    val commitText: String = ""
)

private data class VowelMeta(
    val base: Char,
    val ascii: Char,
    val tone: Tone,
    val upper: Boolean
)

private data class DecomposedSyllable(
    val initial: List<Char> = emptyList(),
    val vowels: List<Char> = emptyList(),
    val finalConsonants: List<Char> = emptyList(),
    val vowelStart: Int = 0,
    val vowelEnd: Int = 0
)

/**
 * Pure-Kotlin Vietnamese Linguistic Engine (1:1 port of VuiKey `internal/engine`).
 * Zero external dependencies, < 1 MB memory footprint, 0.00ms latency for Android IME.
 */
class VuiKeyEngine(var config: VuiKeyConfig = VuiKeyConfig()) {

    private val syllable = ArrayList<Char>(16)
    private val rawHistory = ArrayList<Char>(16)
    private var wordSealed: Boolean = false
    private var lastUndoKey: Char = '\u0000'
    private var delayedD: Boolean = false

    fun reset() {
        syllable.clear()
        rawHistory.clear()
        wordSealed = false
        lastUndoKey = '\u0000'
        delayedD = false
    }

    val currentSyllable: String
        get() = syllable.joinToString("")

    fun handleBackspace() {
        if (syllable.isNotEmpty()) {
            syllable.removeAt(syllable.size - 1)
        }
        // Instant Backspace Unseal: never force the user to delete all the way to the first letter!
        wordSealed = false
        lastUndoKey = '\u0000'
        delayedD = false
        rawHistory.clear()
        rawHistory.addAll(syllable)
    }

    fun handleDelete() {
        reset()
    }

    fun checkMacroOnSpace(): VuiKeyAction? {
        if (!config.enabled || !config.enableMacros || syllable.isEmpty()) {
            reset()
            return null
        }
        val word = currentSyllable
        val replacement = config.macros[word]
        if (!replacement.isNullOrEmpty()) {
            val bs = syllable.size
            reset()
            return VuiKeyAction(swallow = false, backspaces = bs, commitText = replacement)
        }
        reset()
        return null
    }

    fun processKey(ch: Char): VuiKeyAction {
        if (!config.enabled) {
            return VuiKeyAction(swallow = false)
        }

        if (!ch.isLetter() && !ch.isDigit()) {
            reset()
            return VuiKeyAction(swallow = false)
        }

        if (ch.isDigit() && config.mode == InputMode.TELEX) {
            reset()
            return VuiKeyAction(swallow = false)
        }

        rawHistory.add(ch)

        if (wordSealed) {
            syllable.add(ch)
            return VuiKeyAction(swallow = false)
        }

        if (config.smartEnglish && isEnglishOnlyPrefix(syllable, ch)) {
            wordSealed = true
            syllable.add(ch)
            return VuiKeyAction(swallow = false)
        }

        // Sliding Syllable Window: delayed 'd' followed by non-modifier vowel
        if (config.syllableSplit && delayedD && isVietnameseVowel(ch)) {
            val dec = decompose(syllable)
            val isUO = dec.vowels.size == 1 && lowerASCIIString(dec.vowels) == "u" && ch.lowercaseChar() == 'o'
            if (!isModifierForSyllable(dec, ch, config.mode) && !isUO) {
                val oldScreen = ArrayList(syllable)
                if (syllable[0] == 'đ') {
                    syllable[0] = 'd'
                } else if (syllable[0] == 'Đ') {
                    syllable[0] = 'D'
                }
                var prevD = 'd'
                if (rawHistory.size >= 2 && rawHistory[rawHistory.size - 2].isUpperCase()) {
                    prevD = 'D'
                }
                val targetScreen = ArrayList(syllable).apply {
                    add(prevD)
                    add(ch)
                }
                syllable.clear()
                syllable.add(prevD)
                syllable.add(ch)
                rawHistory.clear()
                rawHistory.add(prevD)
                rawHistory.add(ch)
                delayedD = false
                lastUndoKey = '\u0000'
                val (bs, ins) = computeMinimalDiff(oldScreen, targetScreen)
                return VuiKeyAction(swallow = true, backspaces = bs, commitText = ins)
            }
        }
        delayedD = false

        // Sliding Syllable Window: unspaced multi-syllable boundary detection
        if (config.syllableSplit && shouldSplitNewSyllable(syllable, ch, config.mode)) {
            syllable.clear()
            syllable.add(ch)
            rawHistory.clear()
            rawHistory.add(ch)
            lastUndoKey = '\u0000'
            return VuiKeyAction(swallow = false)
        }

        val oldSyl = ArrayList(syllable)
        val (newSyl, transformed, undone) = tryTransform(oldSyl, ch)
        if (transformed && newSyl != null) {
            syllable.clear()
            syllable.addAll(newSyl)
            if (undone) {
                wordSealed = true
            } else if (config.smartEnglish && checkAutoRestoreEnglish(syllable, rawHistory)) {
                syllable.clear()
                syllable.addAll(rawHistory)
                wordSealed = true
            }
            val (bs, ins) = computeMinimalDiff(oldSyl, syllable)
            return VuiKeyAction(swallow = true, backspaces = bs, commitText = ins)
        }

        syllable.add(ch)
        if (config.smartEnglish && checkAutoRestoreEnglish(syllable, rawHistory)) {
            syllable.clear()
            syllable.addAll(rawHistory)
            wordSealed = true
            val oldPlusCh = oldSyl.joinToString("") + ch
            if (oldPlusCh != currentSyllable) {
                val (bs, ins) = computeMinimalDiff(oldSyl, syllable)
                return VuiKeyAction(swallow = true, backspaces = bs, commitText = ins)
            }
        }
        return VuiKeyAction(swallow = false)
    }

    fun isValidVietnameseSyllableForReedit(word: List<Char>): Boolean {
        if (word.isEmpty() || word.size > 7) return false
        if (config.smartEnglish) {
            val rawLower = word.joinToString("") { it.lowercaseChar().toString() }
            if (commonEnglishTechWords.contains(rawLower)) return false
            for (i in 0 until minOf(word.size, 3)) {
                if (isEnglishOnlyPrefix(word.subList(0, i), word[i])) return false
            }
        }
        val dec = decompose(word)
        val initLower = lowerBaseString(dec.initial)
        if (!validInitials.contains(initLower)) return false
        if (dec.vowels.isEmpty()) {
            return dec.initial.size == 1 && (initLower == "d" || initLower == "đ") && dec.finalConsonants.isEmpty()
        }
        if (dec.vowels.size > 3) return false
        val finalLower = lowerASCIIString(dec.finalConsonants)
        if (!validFinals.contains(finalLower)) return false
        if (dec.finalConsonants.isNotEmpty()) {
            val vBase = lowerBaseString(dec.vowels)
            if (closedVowelNuclei.contains(vBase)) return false
        }
        return true
    }

    fun tryModifySurroundingWord(textBeforeCursor: String, ch: Char): VuiKeyAction? {
        if (!config.enabled || syllable.isNotEmpty()) return null
        if (!ch.isLetter() && !ch.isDigit()) return null
        if (ch.isDigit() && config.mode == InputMode.TELEX) return null
        val word = extractLastSyllableBeforeCursor(textBeforeCursor) ?: return null
        if (!isValidVietnameseSyllableForReedit(word)) return null
        val dec = decompose(word)
        if (!isModifierForSyllable(dec, ch, config.mode)) return null

        val savedRaw = ArrayList(rawHistory)
        val savedDelayedD = delayedD
        rawHistory.clear()
        rawHistory.addAll(word)
        rawHistory.add(ch)
        val (newSyl, transformed, undone) = tryTransform(word, ch)
        rawHistory.clear()
        rawHistory.addAll(savedRaw)
        delayedD = savedDelayedD

        if (!transformed || newSyl.isNullOrEmpty() || undone) return null
        val combinedHistory = ArrayList(word).apply { add(ch) }
        if (config.smartEnglish && checkAutoRestoreEnglish(newSyl, combinedHistory)) return null

        val (bs, ins) = computeMinimalDiff(word, newSyl)
        if (bs == 0 && ins.isEmpty()) return null
        return VuiKeyAction(swallow = true, backspaces = bs, commitText = ins)
    }

    private fun tryTransform(oldSyl: List<Char>, ch: Char): Triple<List<Char>?, Boolean, Boolean> {
        if (oldSyl.isEmpty()) return Triple(null, false, false)

        val syl = ArrayList(oldSyl)
        val dec = decompose(syl)
        val lc = ch.lowercaseChar()

        val useTelex = config.mode == InputMode.TELEX || config.mode == InputMode.TELEX_VNI
        val useVNI = config.mode == InputMode.VNI || config.mode == InputMode.TELEX_VNI

        // 1. Handle 'd' -> 'đ' anywhere in syllable or '9' -> 'đ'
        if ((useTelex && lc == 'd' && dec.initial.size == 1) ||
            (useVNI && lc == '9' && dec.initial.isNotEmpty())
        ) {
            val first = syl[0]
            val prevWasD = rawHistory.size >= 2 &&
                (rawHistory[rawHistory.size - 2].lowercaseChar() == 'd' || rawHistory[rawHistory.size - 2] == '9')
            when (first) {
                'd' -> {
                    syl[0] = 'đ'
                    delayedD = dec.vowels.isNotEmpty()
                    return Triple(syl, true, false)
                }
                'D' -> {
                    syl[0] = 'Đ'
                    delayedD = dec.vowels.isNotEmpty()
                    return Triple(syl, true, false)
                }
                'đ' -> if (dec.vowels.isEmpty() || prevWasD) {
                    syl[0] = 'd'
                    if (dec.vowels.isEmpty()) syl.add(ch)
                    return Triple(syl, true, true)
                }
                'Đ' -> if (dec.vowels.isEmpty() || prevWasD) {
                    syl[0] = 'D'
                    if (dec.vowels.isEmpty()) syl.add(ch)
                    return Triple(syl, true, true)
                }
            }
        }

        if (dec.vowels.isEmpty()) {
            if (useTelex && lc == 'w' && dec.initial.isNotEmpty()) {
                val initLower = lowerBaseString(dec.initial)
                if (validInitials.contains(initLower)) {
                    syl.add(if (ch.isUpperCase()) 'Ư' else 'ư')
                    return Triple(syl, true, false)
                }
            }
            return Triple(null, false, false)
        }

        // Absorb 'u' typed right after shorthand/rollover 'w' (e.g. p-h-w-u-c-s -> "phức")
        if (useTelex && lc == 'u' && dec.vowels.size == 1 && dec.finalConsonants.isEmpty()) {
            val vMeta = runeToMeta[dec.vowels[0]]
            if (vMeta != null && vMeta.base.lowercaseChar() == 'ư') {
                if (rawHistory.size >= 2 && rawHistory[rawHistory.size - 2].lowercaseChar() == 'w') {
                    return Triple(syl, true, false)
                }
            }
        }

        // 2. Tone keys
        var wantedTone: Tone? = null
        if (useTelex) {
            wantedTone = when (lc) {
                's' -> Tone.SAC
                'f' -> Tone.HUYEN
                'r' -> Tone.HOI
                'x' -> Tone.NGA
                'j' -> Tone.NANG
                'z' -> Tone.NONE
                else -> null
            }
        }
        if (wantedTone == null && useVNI) {
            wantedTone = when (lc) {
                '1' -> Tone.SAC
                '2' -> Tone.HUYEN
                '3' -> Tone.HOI
                '4' -> Tone.NGA
                '5' -> Tone.NANG
                '0' -> Tone.NONE
                else -> null
            }
        }

        if (wantedTone != null) {
            var curTone = Tone.NONE
            var curToneIdx = -1
            for (i in dec.vowelStart until dec.vowelEnd) {
                val meta = runeToMeta[syl[i]] ?: continue
                if (meta.tone != Tone.NONE) {
                    curTone = meta.tone
                    curToneIdx = i
                    break
                }
            }

            if (wantedTone == Tone.NONE) {
                if (curTone != Tone.NONE) {
                    val meta = runeToMeta[syl[curToneIdx]]!!
                    syl[curToneIdx] = applyToneToBase(meta.base, Tone.NONE)
                    return Triple(syl, true, false)
                }
                var stripped = false
                for (i in syl.indices) {
                    val r = syl[i]
                    if (r == 'đ') {
                        syl[i] = 'd'
                        stripped = true
                    } else if (r == 'Đ') {
                        syl[i] = 'D'
                        stripped = true
                    } else {
                        val meta = runeToMeta[r]
                        if (meta != null && meta.base.lowercaseChar() != meta.ascii) {
                            syl[i] = changeVowelShape(r, meta.ascii)
                            stripped = true
                        }
                    }
                }
                return if (stripped) Triple(syl, true, false) else Triple(null, false, false)
            }

            if (curTone == wantedTone) {
                val meta = runeToMeta[syl[curToneIdx]]!!
                syl[curToneIdx] = applyToneToBase(meta.base, Tone.NONE)
                syl.add(ch)
                return Triple(syl, true, true)
            }

            if (curToneIdx != -1) {
                val meta = runeToMeta[syl[curToneIdx]]!!
                syl[curToneIdx] = applyToneToBase(meta.base, Tone.NONE)
            }

            val targetIdx = findToneTargetIndex(decompose(syl), config.modernToneStyle)
            if (targetIdx >= 0) {
                val meta = runeToMeta[syl[targetIdx]]!!
                syl[targetIdx] = applyToneToBase(meta.base, wantedTone)
                return Triple(syl, true, false)
            }
        }

        // 3. Vowel shape modifiers: a, e, o or VNI 6
        if ((useTelex && (lc == 'a' || lc == 'e' || lc == 'o')) || (useVNI && lc == '6')) {
            for (i in (dec.vowelEnd - 1) downTo dec.vowelStart) {
                val meta = runeToMeta[syl[i]] ?: continue
                if (lc != '6' && meta.ascii != lc) continue
                when (meta.base.lowercaseChar()) {
                    'a' -> {
                        syl[i] = changeVowelShape(syl[i], 'â')
                        return Triple(rebalanceTone(syl, config.modernToneStyle), true, false)
                    }
                    'e' -> {
                        syl[i] = changeVowelShape(syl[i], 'ê')
                        return Triple(rebalanceTone(syl, config.modernToneStyle), true, false)
                    }
                    'o' -> {
                        syl[i] = changeVowelShape(syl[i], 'ô')
                        return Triple(rebalanceTone(syl, config.modernToneStyle), true, false)
                    }
                    'â', 'ê', 'ô' -> {
                        syl[i] = changeVowelShape(syl[i], meta.ascii)
                        syl.add(ch)
                        return Triple(syl, true, true)
                    }
                }
            }
        }

        // 4. Telex 'w' or VNI '7'/'8'
        if ((useTelex && lc == 'w') || (useVNI && (lc == '7' || lc == '8'))) {
            val vASC = lowerASCIIString(dec.vowels)
            val vBase = lowerBaseString(dec.vowels)

            var prevWCount = 0
            var prevWasW = false
            if (rawHistory.size >= 2) {
                val prevKey = rawHistory[rawHistory.size - 2].lowercaseChar()
                prevWasW = (prevKey == 'w' || prevKey == '7' || prevKey == '8')
                for (idx in 0 until (rawHistory.size - 1)) {
                    val rl = rawHistory[idx].lowercaseChar()
                    if (rl == 'w' || rl == '7' || rl == '8') {
                        prevWCount++
                    }
                }
            }

            if (vBase == "ươ" && prevWCount == 1 && !prevWasW) {
                return Triple(syl, true, false)
            }

            if (vBase == "ă" || vBase == "ơ" || vBase == "ư" || vBase == "ươ" ||
                (prevWasW && (vBase.contains('ă') || vBase.contains('ơ') || vBase.contains('ư')))
            ) {
                for (i in dec.vowelStart until dec.vowelEnd) {
                    val meta = runeToMeta[syl[i]] ?: continue
                    syl[i] = changeVowelShape(syl[i], meta.ascii)
                }
                syl.add(ch)
                return Triple(syl, true, true)
            }

            if (vASC.contains("uo") && lc != '8') {
                for (i in dec.vowelStart until (dec.vowelEnd - 1)) {
                    val m1 = runeToMeta[syl[i]] ?: continue
                    val m2 = runeToMeta[syl[i + 1]] ?: continue
                    if (m1.ascii == 'u' && m2.ascii == 'o') {
                        syl[i] = changeVowelShape(syl[i], 'ư')
                        syl[i + 1] = changeVowelShape(syl[i + 1], 'ơ')
                        return Triple(rebalanceTone(syl, config.modernToneStyle), true, false)
                    }
                }
            }

            for (i in (dec.vowelEnd - 1) downTo dec.vowelStart) {
                val meta = runeToMeta[syl[i]] ?: continue
                if (lc == '7' && meta.ascii == 'a') continue
                if (lc == '8' && meta.ascii != 'a') continue
                when (meta.ascii) {
                    'a' -> {
                        if (i > dec.vowelStart && runeToMeta[syl[i - 1]]?.ascii == 'u' && dec.finalConsonants.isEmpty()) {
                            syl[i - 1] = changeVowelShape(syl[i - 1], 'ư')
                            return Triple(rebalanceTone(syl, config.modernToneStyle), true, false)
                        }
                        syl[i] = changeVowelShape(syl[i], 'ă')
                        return Triple(rebalanceTone(syl, config.modernToneStyle), true, false)
                    }
                    'o' -> {
                        syl[i] = changeVowelShape(syl[i], 'ơ')
                        return Triple(rebalanceTone(syl, config.modernToneStyle), true, false)
                    }
                    'u' -> {
                        syl[i] = changeVowelShape(syl[i], 'ư')
                        return Triple(rebalanceTone(syl, config.modernToneStyle), true, false)
                    }
                }
            }
        }

        if (isVietnameseVowel(ch) || (dec.vowels.isNotEmpty() && validFinals.contains(lowerASCIIString(dec.finalConsonants) + lc))) {
            val candidate = ArrayList(oldSyl).apply { add(ch) }
            val rebalanced = rebalanceTone(ArrayList(candidate), config.modernToneStyle)
            if (rebalanced != candidate) {
                return Triple(rebalanced, true, false)
            }
        }

        return Triple(null, false, false)
    }

    companion object {
        private val vowelTable: Map<Char, CharArray> = mapOf(
            'a' to charArrayOf('a', 'á', 'à', 'ả', 'ã', 'ạ'),
            'ă' to charArrayOf('ă', 'ắ', 'ằ', 'ẳ', 'ẵ', 'ặ'),
            'â' to charArrayOf('â', 'ấ', 'ầ', 'ẩ', 'ẫ', 'ậ'),
            'e' to charArrayOf('e', 'é', 'è', 'ẻ', 'ẽ', 'ẹ'),
            'ê' to charArrayOf('ê', 'ế', 'ề', 'ể', 'ễ', 'ệ'),
            'i' to charArrayOf('i', 'í', 'ì', 'ỉ', 'ĩ', 'ị'),
            'o' to charArrayOf('o', 'ó', 'ò', 'ỏ', 'õ', 'ọ'),
            'ô' to charArrayOf('ô', 'ố', 'ồ', 'ổ', 'ỗ', 'ộ'),
            'ơ' to charArrayOf('ơ', 'ớ', 'ờ', 'ở', 'ỡ', 'ợ'),
            'u' to charArrayOf('u', 'ú', 'ù', 'ủ', 'ũ', 'ụ'),
            'ư' to charArrayOf('ư', 'ứ', 'ừ', 'ử', 'ữ', 'ự'),
            'y' to charArrayOf('y', 'ý', 'ỳ', 'ỷ', 'ỹ', 'ỵ'),

            'A' to charArrayOf('A', 'Á', 'À', 'Ả', 'Ã', 'Ạ'),
            'Ă' to charArrayOf('Ă', 'Ắ', 'Ằ', 'Ẳ', 'Ẵ', 'Ặ'),
            'Â' to charArrayOf('Â', 'Ấ', 'Ầ', 'Ẩ', 'Ẫ', 'Ậ'),
            'E' to charArrayOf('E', 'É', 'È', 'Ẻ', 'Ẽ', 'Ẹ'),
            'Ê' to charArrayOf('Ê', 'Ế', 'Ề', 'Ể', 'Ễ', 'Ệ'),
            'I' to charArrayOf('I', 'Í', 'Ì', 'Ỉ', 'Ĩ', 'Ị'),
            'O' to charArrayOf('O', 'Ó', 'Ò', 'Ỏ', 'Õ', 'Ọ'),
            'Ô' to charArrayOf('Ô', 'Ố', 'Ồ', 'Ổ', 'Ỗ', 'Ộ'),
            'Ơ' to charArrayOf('Ơ', 'Ớ', 'Ờ', 'Ở', 'Ỡ', 'Ợ'),
            'U' to charArrayOf('U', 'Ú', 'Ù', 'Ủ', 'Ũ', 'Ụ'),
            'Ư' to charArrayOf('Ư', 'Ứ', 'Ừ', 'Ử', 'Ữ', 'Ự'),
            'Y' to charArrayOf('Y', 'Ý', 'Ỳ', 'Ỷ', 'Ỹ', 'Ỵ')
        )

        private val runeToMeta: Map<Char, VowelMeta> = buildMap {
            val asciiMap = mapOf(
                'a' to 'a', 'ă' to 'a', 'â' to 'a',
                'e' to 'e', 'ê' to 'e',
                'i' to 'i',
                'o' to 'o', 'ô' to 'o', 'ơ' to 'o',
                'u' to 'u', 'ư' to 'u',
                'y' to 'y'
            )
            for ((base, forms) in vowelTable) {
                val lowerBase = base.lowercaseChar()
                val isUp = base.isUpperCase()
                val asc = asciiMap[lowerBase] ?: lowerBase
                for (tIdx in forms.indices) {
                    put(forms[tIdx], VowelMeta(base, asc, Tone.fromIndex(tIdx), isUp))
                }
            }
        }

        private val validInitials = setOf(
            "", "b", "c", "ch", "d", "đ",
            "g", "gh", "gi", "h", "k", "kh",
            "l", "m", "n", "ng", "ngh", "nh",
            "p", "ph", "qu", "r", "s", "t",
            "th", "tr", "v", "x"
        )

        private val validFinals = setOf(
            "", "c", "ch", "m", "n", "ng", "nh", "p", "t"
        )

        private val closedVowelNuclei = setOf(
            "ai", "ao", "au", "ay", "âu", "ây",
            "eo", "êu", "iu", "oi", "ôi", "ơi",
            "ui", "ưi", "ưu", "uya", "uyu",
            "oai", "oao", "oay", "oeo", "uây", "uôi", "ươi", "ươu"
        )

        private val invalidEnglishStarts2 = setOf(
            "bl", "br", "cl", "cr", "dr",
            "fl", "fr", "gl", "gr", "kn",
            "pl", "pr", "sc", "sk", "sl",
            "sm", "sn", "sp", "sr", "st",
            "sw", "tw", "wr", "sh", "wh",
            "ck", "ab", "ad", "ag",
            "ak", "al", "ar", "av"
        )

        private val invalidEnglishStarts3 = setOf(
            "thr", "chr", "str", "spr", "scr", "shr", "spl"
        )

        private val commonEnglishTechWords = setOf(
            "test", "tests", "testing", "tested", "tester",
            "best", "rest", "list", "lists", "last",
            "post", "posts", "host", "hosts", "cost", "costs",
            "most", "lost", "must", "rust", "dust",
            "user", "users", "use", "used", "uses",
            "server", "servers", "service", "services",
            "system", "systems", "status", "static",
            "linux", "ubuntu", "debian", "fedora", "arch",
            "docker", "github", "gitlab", "google", "chrome",
            "python", "golang", "node", "nodes", "code", "codes",
            "mode", "modes", "model", "models", "module", "modules",
            "file", "files", "folder", "folders", "filter", "filters",
            "master", "worker", "router", "render", "buffer",
            "number", "member", "driver", "player", "super",
            "power", "tower", "lower", "upper", "inner", "outer",
            "other", "after", "under", "over", "never", "ever",
            "enter", "center", "header", "footer", "layer",
            "table", "enable", "disable", "value", "values",
            "issue", "issues", "true", "false", "style", "title",
            "time", "times", "name", "names", "game", "games",
            "same", "state", "date", "rate", "create", "update",
            "delete", "remote", "note", "notes", "size", "type",
            "types", "byte", "bytes", "site", "sites", "core",
            "more", "store", "score", "pure", "sure", "fire",
            "wire", "base", "case", "cases", "close", "mouse",
            "house", "route", "rule", "rules", "role", "roles",
            "scale", "sale", "mile", "profile", "mobile",
            "task", "tasks", "mask", "desk", "disk", "risk",
            "dark", "mark", "park", "card", "hard", "bird",
            "port", "ports", "sort", "part", "parts", "cart",
            "hurt", "next", "text", "texts", "left", "soft",
            "lift", "self", "half", "golf", "help", "hold",
            "cold", "gold", "bold", "told", "sold", "open",
            "item", "items", "icon", "icons", "edit", "init",
            "unit", "exit", "admin", "login", "token", "panel",
            "level", "hotel", "pixel", "local", "total", "final",
            "arrow", "arrows", "array", "arrays", "error", "errors",
            "allow", "asset", "assets", "async", "audio", "about"
        )

        fun computeMinimalDiff(oldSyl: List<Char>, newSyl: List<Char>): Pair<Int, String> {
            var common = 0
            val minLen = minOf(oldSyl.size, newSyl.size)
            while (common < minLen && oldSyl[common] == newSyl[common]) {
                common++
            }
            val backspaces = oldSyl.size - common
            val insert = newSyl.subList(common, newSyl.size).joinToString("")
            return Pair(backspaces, insert)
        }

        fun extractLastSyllableBeforeCursor(textBeforeCursor: String): List<Char>? {
            if (textBeforeCursor.isEmpty()) return null
            val end = textBeforeCursor.length
            var start = end
            while (start > 0 && textBeforeCursor[start - 1].isLetter()) {
                start--
            }
            val len = end - start
            if (len == 0 || len > 7) return null
            if (start > 0 && textBeforeCursor[start - 1].isDigit()) return null
            return textBeforeCursor.substring(start, end).toList()
        }

        private fun isVietnameseVowel(r: Char): Boolean = runeToMeta.containsKey(r)

        private fun applyToneToBase(base: Char, tone: Tone): Char {
            val forms = vowelTable[base] ?: return base
            return forms[tone.index]
        }

        private fun changeVowelShape(oldRune: Char, newLowerBase: Char): Char {
            val meta = runeToMeta[oldRune] ?: return oldRune
            val targetBase = if (meta.upper) newLowerBase.uppercaseChar() else newLowerBase
            return applyToneToBase(targetBase, meta.tone)
        }

        private fun decompose(syl: List<Char>): DecomposedSyllable {
            val n = syl.size
            if (n == 0) return DecomposedSyllable()

            var vStart = -1
            for (i in 0 until n) {
                val r = syl[i]
                if (isVietnameseVowel(r)) {
                    if (i == 1 && i + 1 < n && isVietnameseVowel(syl[i + 1])) {
                        val firstLower = syl[0].lowercaseChar()
                        val curASCII = runeToMeta[r]?.ascii
                        if ((firstLower == 'q' && curASCII == 'u') || (firstLower == 'g' && curASCII == 'i')) {
                            continue
                        }
                    }
                    vStart = i
                    break
                }
            }

            if (vStart == -1) {
                return DecomposedSyllable(initial = syl, vowelStart = n, vowelEnd = n)
            }

            var vEnd = vStart
            while (vEnd < n && isVietnameseVowel(syl[vEnd])) {
                vEnd++
            }

            return DecomposedSyllable(
                initial = syl.subList(0, vStart),
                vowels = syl.subList(vStart, vEnd),
                finalConsonants = syl.subList(vEnd, n),
                vowelStart = vStart,
                vowelEnd = vEnd
            )
        }

        private fun lowerASCIIString(runes: List<Char>): String {
            val sb = StringBuilder(runes.size)
            for (r in runes) {
                val meta = runeToMeta[r]
                if (meta != null) {
                    sb.append(meta.ascii)
                } else {
                    val lc = r.lowercaseChar()
                    sb.append(if (lc == 'đ') 'd' else lc)
                }
            }
            return sb.toString()
        }

        private fun lowerBaseString(runes: List<Char>): String {
            val sb = StringBuilder(runes.size)
            for (r in runes) {
                val meta = runeToMeta[r]
                if (meta != null) {
                    sb.append(meta.base.lowercaseChar())
                } else {
                    sb.append(r.lowercaseChar())
                }
            }
            return sb.toString()
        }

        private fun isEnglishOnlyPrefix(syl: List<Char>, nextCh: Char): Boolean {
            val combined = ArrayList(syl).apply { add(nextCh) }
            val s = lowerBaseString(combined)
            return when (s.length) {
                1 -> s == "f" || s == "j" || s == "w" || s == "z"
                2 -> invalidEnglishStarts2.contains(s)
                3 -> invalidEnglishStarts3.contains(s)
                else -> false
            }
        }

        private fun checkAutoRestoreEnglish(currentSyl: List<Char>, rawHistory: List<Char>): Boolean {
            if (rawHistory.size < 3) return false
            val rawLower = rawHistory.joinToString("") { it.lowercaseChar().toString() }
            if (commonEnglishTechWords.contains(rawLower)) return true

            val baseLower = lowerBaseString(currentSyl)
            if (baseLower.any { "ăâđêôơư".indexOf(it) >= 0 }) return false

            val dec = decompose(currentSyl)
            if (dec.finalConsonants.isNotEmpty()) {
                for (r in dec.finalConsonants) {
                    if (isVietnameseVowel(r)) return true
                }
            }
            return false
        }

        private fun isModifierForSyllable(dec: DecomposedSyllable, ch: Char, mode: InputMode): Boolean {
            val lc = ch.lowercaseChar()

            if (lc == 'd' && (mode == InputMode.TELEX || mode == InputMode.TELEX_VNI)) {
                if (dec.initial.size == 1) {
                    val first = dec.initial[0].lowercaseChar()
                    if (first == 'd') return true
                    if (first == 'đ' && dec.vowels.isEmpty()) return true
                }
                return false
            }

            if (dec.vowels.isEmpty()) return false

            if (mode == InputMode.TELEX || mode == InputMode.TELEX_VNI) {
                when (lc) {
                    's', 'f', 'r', 'x', 'j', 'z' -> return true
                    'w' -> {
                        val vStr = lowerASCIIString(dec.vowels)
                        return vStr.any { it == 'a' || it == 'o' || it == 'u' }
                    }
                    'a' -> {
                        val vBase = lowerBaseString(dec.vowels)
                        return vBase.contains('a') || vBase.contains('â')
                    }
                    'e' -> {
                        val vBase = lowerBaseString(dec.vowels)
                        return vBase.contains('e') || vBase.contains('ê')
                    }
                    'o' -> {
                        val vBase = lowerBaseString(dec.vowels)
                        return vBase.contains('o') || vBase.contains('ô')
                    }
                }
            }

            if (mode == InputMode.VNI || mode == InputMode.TELEX_VNI) {
                if (lc in '0'..'9') return true
            }
            return false
        }

        private fun shouldSplitNewSyllable(syl: List<Char>, ch: Char, mode: InputMode): Boolean {
            if (syl.isEmpty()) return false
            val dec = decompose(syl)
            if (dec.vowels.isEmpty()) return false

            if (isModifierForSyllable(dec, ch, mode)) return false

            val lc = ch.lowercaseChar()

            if (dec.finalConsonants.isNotEmpty()) {
                val newFinal = lowerASCIIString(dec.finalConsonants) + lc
                if (validFinals.contains(newFinal)) return false
                if (isVietnameseVowel(ch) && syl.joinToString("").equals(lowerASCIIString(syl), ignoreCase = true)) {
                    return false
                }
                return true
            }

            if (!isVietnameseVowel(ch)) {
                return when (lc) {
                    'c', 'm', 'n', 'p', 't' -> {
                        val vBase = lowerBaseString(dec.vowels)
                        closedVowelNuclei.contains(vBase)
                    }
                    else -> true
                }
            }

            if (dec.vowels.size >= 3) return true
            val vBase = lowerBaseString(dec.vowels)
            return closedVowelNuclei.contains(vBase)
        }

        private fun findToneTargetIndex(dec: DecomposedSyllable, modernStyle: Boolean): Int {
            val vCount = dec.vowels.size
            if (vCount == 0) return -1
            if (vCount == 1) return dec.vowelStart

            var lastSpecial = -1
            for (i in dec.vowels.indices) {
                val meta = runeToMeta[dec.vowels[i]] ?: continue
                val lb = meta.base.lowercaseChar()
                if (lb == 'ă' || lb == 'â' || lb == 'ê' || lb == 'ô' || lb == 'ơ' || lb == 'ư') {
                    lastSpecial = dec.vowelStart + i
                }
            }
            if (lastSpecial != -1) return lastSpecial
            if (vCount == 3) return dec.vowelStart + 1
            if (dec.finalConsonants.isNotEmpty()) return dec.vowelStart + 1

            val vStr = lowerASCIIString(dec.vowels)
            if (modernStyle && (vStr == "oa" || vStr == "oe" || vStr == "uy")) {
                return dec.vowelStart + 1
            }
            return dec.vowelStart
        }

        private fun rebalanceTone(syl: MutableList<Char>, modernStyle: Boolean): List<Char> {
            val dec = decompose(syl)
            if (dec.vowels.isEmpty()) return syl

            if (dec.vowels.size >= 2 && (dec.finalConsonants.isNotEmpty() || dec.vowels.size >= 3)) {
                for (i in dec.vowelStart until (dec.vowelEnd - 1)) {
                    val m1 = runeToMeta[syl[i]] ?: continue
                    val m2 = runeToMeta[syl[i + 1]] ?: continue
                    if (m1.base.lowercaseChar() == 'ư' && m2.base.lowercaseChar() == 'o') {
                        syl[i + 1] = changeVowelShape(syl[i + 1], 'ơ')
                    }
                }
            }

            var currentTone = Tone.NONE
            for (i in 0 until dec.vowelEnd) {
                val meta = runeToMeta[syl[i]] ?: continue
                if (meta.tone != Tone.NONE) {
                    currentTone = meta.tone
                    syl[i] = applyToneToBase(meta.base, Tone.NONE)
                }
            }
            if (currentTone != Tone.NONE) {
                val target = findToneTargetIndex(decompose(syl), modernStyle)
                if (target >= 0) {
                    val meta = runeToMeta[syl[target]]!!
                    syl[target] = applyToneToBase(meta.base, currentTone)
                }
            }
            return syl
        }
    }
}
