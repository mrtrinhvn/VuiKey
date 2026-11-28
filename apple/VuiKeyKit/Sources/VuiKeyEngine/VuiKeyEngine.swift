import Foundation

/// InputMode defines the Vietnamese typing method.
public enum InputMode: Int, Codable, CaseIterable {
    case telex = 0
    case vni = 1
    case telexVNI = 2

    public var displayName: String {
        switch self {
        case .telex: return "Telex"
        case .vni: return "VNI"
        case .telexVNI: return "Telex + VNI"
        }
    }
}

/// Tone represents the 6 Vietnamese tones.
public enum Tone: Int {
    case none = 0
    case sac = 1   // s / 1
    case huyen = 2 // f / 2
    case hoi = 3   // r / 3
    case nga = 4   // x / 4
    case nang = 5  // j / 5
}

/// Config holds user preferences for VuiKey's linguistic engine on macOS & iOS.
public struct VuiKeyConfig: Codable {
    public var enabled: Bool
    public var mode: InputMode
    public var smartEnglish: Bool
    public var syllableSplit: Bool
    public var modernToneStyle: Bool
    public var enableMacros: Bool
    public var macros: [String: String]

    public init(
        enabled: Bool = true,
        mode: InputMode = .telex,
        smartEnglish: Bool = true,
        syllableSplit: Bool = true,
        modernToneStyle: Bool = false,
        enableMacros: Bool = false,
        macros: [String: String] = ["vn": "Việt Nam"]
    ) {
        self.enabled = enabled
        self.mode = mode
        self.smartEnglish = smartEnglish
        self.syllableSplit = syllableSplit
        self.modernToneStyle = modernToneStyle
        self.enableMacros = enableMacros
        self.macros = macros
    }
}

/// Action describes what the macOS (`CGEventTap`) or iOS (`UITextDocumentProxy`) layer should do.
public struct VuiKeyAction: Equatable {
    public let swallow: Bool
    public let backspaces: Int
    public let commitText: String

    public init(swallow: Bool = false, backspaces: Int = 0, commitText: String = "") {
        self.swallow = swallow
        self.backspaces = backspaces
        self.commitText = commitText
    }
}

private struct VowelMeta {
    let base: Character
    let ascii: Character
    let tone: Tone
    let upper: Bool
}

/// Pure-Swift Vietnamese Linguistic Engine (1:1 port of VuiKey `internal/engine`).
/// Ultra-lightweight (< 1 MB memory footprint), safe for iOS Keyboard Extensions & macOS CGEventTap.
public final class VuiKeyEngine {
    public var config: VuiKeyConfig

    private var syllable: [Character] = []
    private var rawHistory: [Character] = []
    private var wordSealed: Bool = false
    private var lastUndoKey: Character? = nil
    private var delayedD: Bool = false

    private static let vowelTable: [Character: [Character]] = [
        "a": ["a", "á", "à", "ả", "ã", "ạ"],
        "ă": ["ă", "ắ", "ằ", "ẳ", "ẵ", "ặ"],
        "â": ["â", "ấ", "ầ", "ẩ", "ẫ", "ậ"],
        "e": ["e", "é", "è", "ẻ", "ẽ", "ẹ"],
        "ê": ["ê", "ế", "ề", "ể", "ễ", "ệ"],
        "i": ["i", "í", "ì", "ỉ", "ĩ", "ị"],
        "o": ["o", "ó", "ò", "ỏ", "õ", "ọ"],
        "ô": ["ô", "ố", "ồ", "ổ", "ỗ", "ộ"],
        "ơ": ["ơ", "ớ", "ờ", "ở", "ỡ", "ợ"],
        "u": ["u", "ú", "ù", "ủ", "ũ", "ụ"],
        "ư": ["ư", "ứ", "ừ", "ử", "ữ", "ự"],
        "y": ["y", "ý", "ỳ", "ỷ", "ỹ", "ỵ"],

        "A": ["A", "Á", "À", "Ả", "Ã", "Ạ"],
        "Ă": ["Ă", "Ắ", "Ằ", "Ẳ", "Ẵ", "Ặ"],
        "Â": ["Â", "Ấ", "Ầ", "Ẩ", "Ẫ", "Ậ"],
        "E": ["E", "É", "È", "Ẻ", "Ẽ", "Ẹ"],
        "Ê": ["Ê", "Ế", "Ề", "Ể", "Ễ", "Ệ"],
        "I": ["I", "Í", "Ì", "Ỉ", "Ĩ", "Ị"],
        "O": ["O", "Ó", "Ò", "Ỏ", "Õ", "Ọ"],
        "Ô": ["Ô", "Ố", "Ồ", "Ổ", "Ỗ", "Ộ"],
        "Ơ": ["Ơ", "Ớ", "Ờ", "Ở", "Ỡ", "Ợ"],
        "U": ["U", "Ú", "Ù", "Ủ", "Ũ", "Ụ"],
        "Ư": ["Ư", "Ứ", "Ừ", "Ử", "Ữ", "Ự"],
        "Y": ["Y", "Ý", "Ỳ", "Ỷ", "Ỹ", "Ỵ"]
    ]

    private static let runeToMeta: [Character: VowelMeta] = {
        let asciiMap: [Character: Character] = [
            "a": "a", "ă": "a", "â": "a",
            "e": "e", "ê": "e",
            "i": "i",
            "o": "o", "ô": "o", "ơ": "o",
            "u": "u", "ư": "u",
            "y": "y"
        ]
        var map: [Character: VowelMeta] = [:]
        for (base, forms) in vowelTable {
            let lowerBase = Character(base.lowercased())
            let isUp = base.isUppercase
            let asc = asciiMap[lowerBase] ?? lowerBase
            for (tIdx, r) in forms.enumerated() {
                if let tone = Tone(rawValue: tIdx) {
                    map[r] = VowelMeta(base: base, ascii: asc, tone: tone, upper: isUp)
                }
            }
        }
        return map
    }()

    private static let validInitials: Set<String> = [
        "", "b", "c", "ch", "d", "đ",
        "g", "gh", "gi", "h", "k", "kh",
        "l", "m", "n", "ng", "ngh", "nh",
        "p", "ph", "qu", "r", "s", "t",
        "th", "tr", "v", "x"
    ]

    private static let validFinals: Set<String> = [
        "", "c", "ch", "m", "n", "ng", "nh", "p", "t"
    ]

    private static let closedVowelNuclei: Set<String> = [
        "ai", "ao", "au", "ay", "âu", "ây",
        "eo", "êu", "iu", "oi", "ôi", "ơi",
        "ui", "ưi", "ưu", "uya", "uyu",
        "oai", "oao", "oay", "oeo", "uây", "uôi", "ươi", "ươu"
    ]

    private static let invalidEnglishStarts2: Set<String> = [
        "bl", "br", "cl", "cr", "dr",
        "fl", "fr", "gl", "gr", "kn",
        "pl", "pr", "sc", "sk", "sl",
        "sm", "sn", "sp", "sr", "st",
        "sw", "tw", "wr", "sh", "wh",
        "ck", "ab", "ad", "ag",
        "ak", "al", "ar", "av"
    ]

    private static let invalidEnglishStarts3: Set<String> = [
        "thr", "chr", "str", "spr", "scr", "shr", "spl"
    ]

    private static let commonEnglishTechWords: Set<String> = [
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
    ]

    public init(config: VuiKeyConfig = VuiKeyConfig()) {
        self.config = config
    }

    public func reset() {
        syllable.removeAll(keepingCapacity: true)
        rawHistory.removeAll(keepingCapacity: true)
        wordSealed = false
        lastUndoKey = nil
        delayedD = false
    }

    public var currentSyllable: String {
        String(syllable)
    }

    public func handleBackspace() {
        if !syllable.isEmpty {
            syllable.removeLast()
        }
        wordSealed = false
        lastUndoKey = nil
        delayedD = false
        rawHistory = syllable
    }

    public func handleDelete() {
        reset()
    }

    public func checkMacroOnSpace() -> VuiKeyAction? {
        guard config.enabled && config.enableMacros && !syllable.isEmpty else {
            reset()
            return nil
        }
        let word = String(syllable)
        if let replacement = config.macros[word], !replacement.isEmpty {
            let bs = syllable.count
            reset()
            return VuiKeyAction(swallow: false, backspaces: bs, commitText: replacement)
        }
        reset()
        return nil
    }

    public func processKey(_ ch: Character) -> VuiKeyAction {
        guard config.enabled else {
            return VuiKeyAction(swallow: false)
        }

        if !ch.isLetter && !ch.isNumber {
            reset()
            return VuiKeyAction(swallow: false)
        }

        if ch.isNumber && config.mode == .telex {
            reset()
            return VuiKeyAction(swallow: false)
        }

        rawHistory.append(ch)

        if wordSealed {
            syllable.append(ch)
            return VuiKeyAction(swallow: false)
        }

        if config.smartEnglish && Self.isEnglishOnlyPrefix(syllable, nextCh: ch) {
            wordSealed = true
            syllable.append(ch)
            return VuiKeyAction(swallow: false)
        }

        // Sliding Syllable Window: delayed 'd' followed by non-modifier vowel
        if config.syllableSplit && delayedD && Self.isVietnameseVowel(ch) {
            let dec = Self.decompose(syllable)
            let isUO = dec.vowels.count == 1 && Self.lowerASCIIString(dec.vowels) == "u" && ch.lowercased() == "o"
            if !Self.isModifierForSyllable(dec, ch: ch, mode: config.mode) && !isUO {
                let oldScreen = syllable
                if syllable[0] == "đ" {
                    syllable[0] = "d"
                } else if syllable[0] == "Đ" {
                    syllable[0] = "D"
                }
                var prevD: Character = "d"
                if rawHistory.count >= 2 && rawHistory[rawHistory.count - 2].isUppercase {
                    prevD = "D"
                }
                let targetScreen = syllable + [prevD, ch]
                syllable = [prevD, ch]
                rawHistory = [prevD, ch]
                delayedD = false
                lastUndoKey = nil
                let (bs, ins) = Self.computeMinimalDiff(oldSyl: oldScreen, newSyl: targetScreen)
                return VuiKeyAction(swallow: true, backspaces: bs, commitText: ins)
            }
        }
        delayedD = false

        if config.syllableSplit && Self.shouldSplitNewSyllable(syllable, ch: ch, mode: config.mode) {
            syllable = [ch]
            rawHistory = [ch]
            lastUndoKey = nil
            return VuiKeyAction(swallow: false)
        }

        let oldSyl = syllable
        let (newSyl, transformed, undone) = tryTransform(oldSyl, ch: ch)
        if transformed, let newSyl = newSyl {
            syllable = newSyl
            if undone {
                wordSealed = true
            } else if config.smartEnglish && Self.checkAutoRestoreEnglish(currentSyl: syllable, rawHistory: rawHistory) {
                syllable = rawHistory
                wordSealed = true
            }
            let (bs, ins) = Self.computeMinimalDiff(oldSyl: oldSyl, newSyl: syllable)
            return VuiKeyAction(swallow: true, backspaces: bs, commitText: ins)
        }

        syllable.append(ch)
        if config.smartEnglish && Self.checkAutoRestoreEnglish(currentSyl: syllable, rawHistory: rawHistory) {
            syllable = rawHistory
            wordSealed = true
            if String(oldSyl) + String(ch) != String(syllable) {
                let (bs, ins) = Self.computeMinimalDiff(oldSyl: oldSyl, newSyl: syllable)
                return VuiKeyAction(swallow: true, backspaces: bs, commitText: ins)
            }
        }
        return VuiKeyAction(swallow: false)
    }

    public static func computeMinimalDiff(oldSyl: [Character], newSyl: [Character]) -> (backspaces: Int, insert: String) {
        var common = 0
        let minLen = min(oldSyl.count, newSyl.count)
        while common < minLen && oldSyl[common] == newSyl[common] {
            common += 1
        }
        let backspaces = oldSyl.count - common
        let insert = String(newSyl[common...])
        return (backspaces, insert)
    }

    // MARK: - Phonology Helpers

    private struct DecomposedSyllable {
        var initial: [Character] = []
        var vowels: [Character] = []
        var final: [Character] = []
        var vowelStart: Int = 0
        var vowelEnd: Int = 0
    }

    private static func isVietnameseVowel(_ r: Character) -> Bool {
        runeToMeta[r] != nil
    }

    private static func applyToneToBase(_ base: Character, tone: Tone) -> Character {
        if let forms = vowelTable[base] {
            return forms[tone.rawValue]
        }
        return base
    }

    private static func changeVowelShape(_ oldRune: Character, newLowerBase: Character) -> Character {
        guard let meta = runeToMeta[oldRune] else { return oldRune }
        let targetBase: Character = meta.upper ? Character(newLowerBase.uppercased()) : newLowerBase
        return applyToneToBase(targetBase, tone: meta.tone)
    }

    private static func decompose(_ syl: [Character]) -> DecomposedSyllable {
        let n = syl.count
        if n == 0 { return DecomposedSyllable() }

        var vStart = -1
        for (i, r) in syl.enumerated() {
            if isVietnameseVowel(r) {
                if i == 1 && i + 1 < n && isVietnameseVowel(syl[i + 1]) {
                    let firstLower = Character(syl[0].lowercased())
                    let curASCII = runeToMeta[r]?.ascii
                    if (firstLower == "q" && curASCII == "u") || (firstLower == "g" && curASCII == "i") {
                        continue
                    }
                }
                vStart = i
                break
            }
        }

        if vStart == -1 {
            return DecomposedSyllable(initial: syl, vowels: [], final: [], vowelStart: n, vowelEnd: n)
        }

        var vEnd = vStart
        while vEnd < n && isVietnameseVowel(syl[vEnd]) {
            vEnd += 1
        }

        return DecomposedSyllable(
            initial: Array(syl[0..<vStart]),
            vowels: Array(syl[vStart..<vEnd]),
            final: Array(syl[vEnd..<n]),
            vowelStart: vStart,
            vowelEnd: vEnd
        )
    }

    private static func lowerASCIIString(_ runes: [Character]) -> String {
        var out = ""
        out.reserveCapacity(runes.count)
        for r in runes {
            if let meta = runeToMeta[r] {
                out.append(meta.ascii)
            } else {
                let lc = Character(r.lowercased())
                out.append(lc == "đ" ? "d" : lc)
            }
        }
        return out
    }

    private static func lowerBaseString(_ runes: [Character]) -> String {
        var out = ""
        out.reserveCapacity(runes.count)
        for r in runes {
            if let meta = runeToMeta[r] {
                out.append(contentsOf: meta.base.lowercased())
            } else {
                out.append(contentsOf: r.lowercased())
            }
        }
        return out
    }

    private static func isEnglishOnlyPrefix(_ syl: [Character], nextCh: Character) -> Bool {
        let combined = syl + [nextCh]
        let s = lowerBaseString(combined)
        if combined.count == 1 {
            return s == "f" || s == "j" || s == "w" || s == "z"
        }
        if combined.count == 2 {
            return invalidEnglishStarts2.contains(s)
        }
        if combined.count == 3 {
            return invalidEnglishStarts3.contains(s)
        }
        return false
    }

    private static func checkAutoRestoreEnglish(currentSyl: [Character], rawHistory: [Character]) -> Bool {
        if rawHistory.count < 3 { return false }
        let rawLower = String(rawHistory).lowercased()
        if commonEnglishTechWords.contains(rawLower) {
            return true
        }
        let baseLower = lowerBaseString(currentSyl)
        if baseLower.contains(where: { "ăâđêôơư".contains($0) }) {
            return false
        }
        let dec = decompose(currentSyl)
        if !dec.final.isEmpty {
            for r in dec.final {
                if isVietnameseVowel(r) {
                    return true
                }
            }
        }
        return false
    }

    private static func isModifierForSyllable(_ dec: DecomposedSyllable, ch: Character, mode: InputMode) -> Bool {
        let lc = Character(ch.lowercased())

        if lc == "d" && (mode == .telex || mode == .telexVNI) {
            if dec.initial.count == 1 {
                let first = Character(dec.initial[0].lowercased())
                if first == "d" { return true }
                if first == "đ" && dec.vowels.isEmpty { return true }
            }
            return false
        }

        if dec.vowels.isEmpty { return false }

        if mode == .telex || mode == .telexVNI {
            switch lc {
            case "s", "f", "r", "x", "j", "z":
                return true
            case "w":
                let vStr = lowerASCIIString(dec.vowels)
                return vStr.contains(where: { "aou".contains($0) })
            case "a":
                let vBase = lowerBaseString(dec.vowels)
                return vBase.contains("a") || vBase.contains("â")
            case "e":
                let vBase = lowerBaseString(dec.vowels)
                return vBase.contains("e") || vBase.contains("ê")
            case "o":
                let vBase = lowerBaseString(dec.vowels)
                return vBase.contains("o") || vBase.contains("ô")
            default:
                break
            }
        }

        if mode == .vni || mode == .telexVNI {
            if ("0"..."9").contains(lc) {
                return true
            }
        }

        return false
    }

    private static func shouldSplitNewSyllable(_ syl: [Character], ch: Character, mode: InputMode) -> Bool {
        if syl.isEmpty { return false }
        let dec = decompose(syl)
        if dec.vowels.isEmpty { return false }

        if isModifierForSyllable(dec, ch: ch, mode: mode) {
            return false
        }

        let lc = Character(ch.lowercased())

        if !dec.final.isEmpty {
            let newFinal = lowerASCIIString(dec.final) + String(lc)
            if validFinals.contains(newFinal) {
                return false
            }
            if isVietnameseVowel(ch) && String(syl).lowercased() == lowerASCIIString(syl) {
                return false
            }
            return true
        }

        if !isVietnameseVowel(ch) {
            switch lc {
            case "c", "m", "n", "p", "t":
                let vBase = lowerBaseString(dec.vowels)
                if closedVowelNuclei.contains(vBase) {
                    return true
                }
                return false
            default:
                return true
            }
        }

        if dec.vowels.count >= 3 {
            return true
        }
        let vBase = lowerBaseString(dec.vowels)
        if closedVowelNuclei.contains(vBase) {
            return true
        }
        return false
    }

    private static func findToneTargetIndex(_ dec: DecomposedSyllable, modernStyle: Bool) -> Int {
        let vCount = dec.vowels.count
        if vCount == 0 { return -1 }
        if vCount == 1 { return dec.vowelStart }

        var lastSpecial = -1
        for (i, r) in dec.vowels.enumerated() {
            if let meta = runeToMeta[r] {
                let lb = Character(meta.base.lowercased())
                if "ăâêôơư".contains(lb) {
                    lastSpecial = dec.vowelStart + i
                }
            }
        }
        if lastSpecial != -1 {
            return lastSpecial
        }

        if vCount == 3 {
            return dec.vowelStart + 1
        }

        if !dec.final.isEmpty {
            return dec.vowelStart + 1
        }

        let vStr = lowerASCIIString(dec.vowels)
        if modernStyle && (vStr == "oa" || vStr == "oe" || vStr == "uy") {
            return dec.vowelStart + 1
        }
        return dec.vowelStart
    }

    private static func rebalanceTone(_ inputSyl: [Character], modernStyle: Bool) -> [Character] {
        var syl = inputSyl
        let dec = decompose(syl)
        if dec.vowels.isEmpty { return syl }

        if dec.vowels.count >= 2 {
            if !dec.final.isEmpty || dec.vowels.count >= 3 {
                var i = dec.vowelStart
                while i + 1 < dec.vowelEnd {
                    if let m1 = runeToMeta[syl[i]], let m2 = runeToMeta[syl[i + 1]] {
                        let b1 = m1.base.lowercased()
                        let b2 = m2.base.lowercased()
                        if (b1 == "ư" && b2 == "o") || (b1 == "u" && b2 == "ơ") {
                            syl[i] = changeVowelShape(syl[i], newLowerBase: "ư")
                            syl[i + 1] = changeVowelShape(syl[i + 1], newLowerBase: "ơ")
                        }
                    }
                    i += 1
                }
            } else if dec.vowels.count == 2 && dec.final.isEmpty {
                let i = dec.vowelStart
                if let m1 = runeToMeta[syl[i]], let m2 = runeToMeta[syl[i + 1]] {
                    if m1.base.lowercased() == "ư" && m2.base.lowercased() == "o" {
                        syl[i] = changeVowelShape(syl[i], newLowerBase: "u")
                        syl[i + 1] = changeVowelShape(syl[i + 1], newLowerBase: "ơ")
                    }
                }
            }
        }

        var currentTone: Tone = .none
        for i in 0..<dec.vowelEnd {
            if let meta = runeToMeta[syl[i]], meta.tone != .none {
                currentTone = meta.tone
                syl[i] = applyToneToBase(meta.base, tone: .none)
            }
        }
        if currentTone != .none {
            let target = findToneTargetIndex(decompose(syl), modernStyle: modernStyle)
            if target >= 0, let meta = runeToMeta[syl[target]] {
                syl[target] = applyToneToBase(meta.base, tone: currentTone)
            }
        }
        return syl
    }

    private func tryTransform(_ oldSyl: [Character], ch: Character) -> ([Character]?, Bool, Bool) {
        if oldSyl.isEmpty { return (nil, false, false) }

        var syl = oldSyl
        let dec = Self.decompose(syl)
        let lc = Character(ch.lowercased())

        let useTelex = config.mode == .telex || config.mode == .telexVNI
        let useVNI = config.mode == .vni || config.mode == .telexVNI

        // 1. 'd' -> 'đ'
        if (useTelex && lc == "d" && dec.initial.count == 1) ||
           (useVNI && lc == "9" && !dec.initial.isEmpty) {
            let first = syl[0]
            let prevWasD: Bool = {
                guard rawHistory.count >= 2 else { return false }
                let prev = Character(rawHistory[rawHistory.count - 2].lowercased())
                return prev == "d" || prev == "9"
            }()
            switch first {
            case "d":
                syl[0] = "đ"
                delayedD = !dec.vowels.isEmpty
                return (syl, true, false)
            case "D":
                syl[0] = "Đ"
                delayedD = !dec.vowels.isEmpty
                return (syl, true, false)
            case "đ":
                if dec.vowels.isEmpty || prevWasD {
                    syl[0] = "d"
                    if dec.vowels.isEmpty { syl.append(ch) }
                    return (syl, true, true)
                }
            case "Đ":
                if dec.vowels.isEmpty || prevWasD {
                    syl[0] = "D"
                    if dec.vowels.isEmpty { syl.append(ch) }
                    return (syl, true, true)
                }
            default:
                break
            }
        }

        if dec.vowels.isEmpty {
            if useTelex && lc == "w" && !dec.initial.isEmpty {
                let initLower = Self.lowerBaseString(dec.initial)
                if Self.validInitials.contains(initLower) {
                    syl.append(ch.isUppercase ? "Ư" : "ư")
                    return (syl, true, false)
                }
            }
            return (nil, false, false)
        }

        if useTelex && lc == "u" && dec.vowels.count == 1 && dec.final.isEmpty,
           let meta = Self.runeToMeta[dec.vowels[0]], meta.base.lowercased() == "ư" {
            if rawHistory.count >= 2 && rawHistory[rawHistory.count - 2].lowercased() == "w" {
                return (syl, true, false)
            }
        }

        // 2. Tone keys
        var wantedTone: Tone? = nil
        if useTelex {
            switch lc {
            case "s": wantedTone = .sac
            case "f": wantedTone = .huyen
            case "r": wantedTone = .hoi
            case "x": wantedTone = .nga
            case "j": wantedTone = .nang
            case "z": wantedTone = Tone.none
            default: break
            }
        }
        if wantedTone == nil && useVNI {
            switch lc {
            case "1": wantedTone = .sac
            case "2": wantedTone = .huyen
            case "3": wantedTone = .hoi
            case "4": wantedTone = .nga
            case "5": wantedTone = .nang
            case "0": wantedTone = Tone.none
            default: break
            }
        }

        if let wantedTone = wantedTone {
            var curTone: Tone = .none
            var curToneIdx = -1
            for i in 0..<dec.vowelEnd {
                if let meta = Self.runeToMeta[syl[i]], meta.tone != .none {
                    curTone = meta.tone
                    curToneIdx = i
                    break
                }
            }

            if wantedTone == .none {
                if curTone != .none {
                    if let meta = Self.runeToMeta[syl[curToneIdx]] {
                        syl[curToneIdx] = Self.applyToneToBase(meta.base, tone: .none)
                    }
                    return (syl, true, false)
                }
                var stripped = false
                for (i, r) in syl.enumerated() {
                    if r == "đ" { syl[i] = "d"; stripped = true }
                    else if r == "Đ" { syl[i] = "D"; stripped = true }
                    else if let meta = Self.runeToMeta[r], Character(meta.base.lowercased()) != meta.ascii {
                        syl[i] = Self.changeVowelShape(r, newLowerBase: meta.ascii)
                        stripped = true
                    }
                }
                return stripped ? (syl, true, false) : (nil, false, false)
            }

            if curTone == wantedTone {
                if let meta = Self.runeToMeta[syl[curToneIdx]] {
                    syl[curToneIdx] = Self.applyToneToBase(meta.base, tone: .none)
                }
                syl.append(ch)
                return (syl, true, true)
            }

            if curToneIdx != -1, let meta = Self.runeToMeta[syl[curToneIdx]] {
                syl[curToneIdx] = Self.applyToneToBase(meta.base, tone: .none)
            }

            let targetIdx = Self.findToneTargetIndex(Self.decompose(syl), modernStyle: config.modernToneStyle)
            if targetIdx >= 0, let meta = Self.runeToMeta[syl[targetIdx]] {
                syl[targetIdx] = Self.applyToneToBase(meta.base, tone: wantedTone)
                return (syl, true, false)
            }
        }

        // 3. Vowel shape modifiers: a, e, o or VNI 6
        if (useTelex && (lc == "a" || lc == "e" || lc == "o")) || (useVNI && lc == "6") {
            var i = dec.vowelEnd - 1
            while i >= dec.vowelStart {
                if let meta = Self.runeToMeta[syl[i]] {
                    if lc != "6" && meta.ascii != lc {
                        i -= 1
                        continue
                    }
                    let lb = Character(meta.base.lowercased())
                    switch lb {
                    case "a":
                        syl[i] = Self.changeVowelShape(syl[i], newLowerBase: "â")
                        return (Self.rebalanceTone(syl, modernStyle: config.modernToneStyle), true, false)
                    case "e":
                        syl[i] = Self.changeVowelShape(syl[i], newLowerBase: "ê")
                        return (Self.rebalanceTone(syl, modernStyle: config.modernToneStyle), true, false)
                    case "o":
                        syl[i] = Self.changeVowelShape(syl[i], newLowerBase: "ô")
                        return (Self.rebalanceTone(syl, modernStyle: config.modernToneStyle), true, false)
                    case "â", "ê", "ô":
                        syl[i] = Self.changeVowelShape(syl[i], newLowerBase: meta.ascii)
                        syl.append(ch)
                        return (syl, true, true)
                    default:
                        break
                    }
                }
                i -= 1
            }
        }

        // 4. Telex 'w' or VNI '7'/'8'
        if (useTelex && lc == "w") || (useVNI && (lc == "7" || lc == "8")) {
            let vASC = Self.lowerASCIIString(dec.vowels)
            let vBase = Self.lowerBaseString(dec.vowels)

            var prevWCount = 0
            var prevWasW = false
            if rawHistory.count >= 2 {
                let prevKey = Character(rawHistory[rawHistory.count - 2].lowercased())
                prevWasW = (prevKey == "w" || prevKey == "7" || prevKey == "8")
                for rk in rawHistory.dropLast() {
                    let rl = Character(rk.lowercased())
                    if rl == "w" || rl == "7" || rl == "8" {
                        prevWCount += 1
                    }
                }
            }

            if (vBase == "ươ" || vBase == "uơ") && prevWCount == 1 && !prevWasW {
                return (syl, true, false)
            }

            if vBase == "ă" || vBase == "ơ" || vBase == "ư" || vBase == "ươ" || vBase == "uơ" ||
               (prevWasW && vBase.contains(where: { "ăơư".contains($0) })) {
                for i in dec.vowelStart..<dec.vowelEnd {
                    if let meta = Self.runeToMeta[syl[i]] {
                        syl[i] = Self.changeVowelShape(syl[i], newLowerBase: meta.ascii)
                    }
                }
                syl.append(ch)
                return (syl, true, true)
            }

            if vASC.contains("uo") && lc != "8" {
                var i = dec.vowelStart
                while i + 1 < dec.vowelEnd {
                    if let m1 = Self.runeToMeta[syl[i]], let m2 = Self.runeToMeta[syl[i + 1]],
                       m1.ascii == "u" && m2.ascii == "o" {
                        if !dec.final.isEmpty || dec.vowels.count >= 3 {
                            syl[i] = Self.changeVowelShape(syl[i], newLowerBase: "ư")
                            syl[i + 1] = Self.changeVowelShape(syl[i + 1], newLowerBase: "ơ")
                        } else {
                            syl[i] = Self.changeVowelShape(syl[i], newLowerBase: "u")
                            syl[i + 1] = Self.changeVowelShape(syl[i + 1], newLowerBase: "ơ")
                        }
                        return (Self.rebalanceTone(syl, modernStyle: config.modernToneStyle), true, false)
                    }
                    i += 1
                }
            }

            var i = dec.vowelEnd - 1
            while i >= dec.vowelStart {
                if let meta = Self.runeToMeta[syl[i]] {
                    if lc == "7" && meta.ascii == "a" { i -= 1; continue }
                    if lc == "8" && meta.ascii != "a" { i -= 1; continue }
                    switch meta.ascii {
                    case "a":
                        if i > dec.vowelStart, let prevMeta = Self.runeToMeta[syl[i - 1]], prevMeta.ascii == "u", dec.final.isEmpty {
                            syl[i - 1] = Self.changeVowelShape(syl[i - 1], newLowerBase: "ư")
                            return (Self.rebalanceTone(syl, modernStyle: config.modernToneStyle), true, false)
                        }
                        syl[i] = Self.changeVowelShape(syl[i], newLowerBase: "ă")
                        return (Self.rebalanceTone(syl, modernStyle: config.modernToneStyle), true, false)
                    case "o":
                        syl[i] = Self.changeVowelShape(syl[i], newLowerBase: "ơ")
                        return (Self.rebalanceTone(syl, modernStyle: config.modernToneStyle), true, false)
                    case "u":
                        syl[i] = Self.changeVowelShape(syl[i], newLowerBase: "ư")
                        return (Self.rebalanceTone(syl, modernStyle: config.modernToneStyle), true, false)
                    default:
                        break
                    }
                }
                i -= 1
            }
        }

        if Self.isVietnameseVowel(ch) || (!dec.vowels.isEmpty && Self.validFinals.contains(Self.lowerASCIIString(dec.final) + String(lc))) {
            let candidate = oldSyl + [ch]
            let rebalanced = Self.rebalanceTone(candidate, modernStyle: config.modernToneStyle)
            if rebalanced != candidate {
                return (rebalanced, true, false)
            }
        }

        return (nil, false, false)
    }
}
