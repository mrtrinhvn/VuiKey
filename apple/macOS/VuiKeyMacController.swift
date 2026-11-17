#if os(macOS)
import Cocoa
import ApplicationServices
import Carbon
import VuiKeyEngine

/// Native Swift macOS EventTap & Menu Bar Controller for VuiKey.
public final class VuiKeyMacController: NSObject {
    public static let shared = VuiKeyMacController()

    public let engine = VuiKeyEngine()
    private var eventTap: CFMachPort?
    private var statusItem: NSStatusItem?

    private static let eventMarker: Int64 = 0x5655494B // "VUIK"

    private static let autocompleteBundleIDs: Set<String> = [
        "com.apple.Safari",
        "com.google.Chrome",
        "company.thebrowser.Browser", // Arc
        "com.brave.Browser",
        "com.microsoft.edgemac",
        "org.mozilla.firefox",
        "com.apple.Spotlight"
    ]

    public func start() {
        requestAccessibilityPermission()
        setupStatusBar()
        setupEventTap()
    }

    public func requestAccessibilityPermission() {
        let options = [kAXTrustedCheckOptionPrompt.takeUnretainedValue() as String: true] as CFDictionary
        _ = AXIsProcessTrustedWithOptions(options)
    }

    private func setupStatusBar() {
        let item = NSStatusBar.system.statusItem(withLength: NSStatusItem.variableLength)
        self.statusItem = item
        updateStatusBarTitle()

        let menu = NSMenu(title: "VuiKey")
        let toggleItem = NSMenuItem(title: "Bật/Tắt Tiếng Việt (Ctrl+Space)", action: #selector(toggleEnabled), keyEquivalent: "")
        toggleItem.target = self
        menu.addItem(toggleItem)

        menu.addItem(.separator())

        for mode in InputMode.allCases {
            let modeItem = NSMenuItem(title: "Kiểu gõ: \(mode.displayName)", action: #selector(selectMode(_:)), keyEquivalent: "")
            modeItem.tag = mode.rawValue
            modeItem.target = self
            menu.addItem(modeItem)
        }

        menu.addItem(.separator())

        let quitItem = NSMenuItem(title: "Thoát VuiKey", action: #selector(NSApplication.terminate(_:)), keyEquivalent: "q")
        menu.addItem(quitItem)

        item.menu = menu
    }

    private func updateStatusBarTitle() {
        statusItem?.button?.title = engine.config.enabled ? "🇻🇳 V" : "🇬🇧 E"
    }

    @objc private func toggleEnabled() {
        engine.config.enabled.toggle()
        engine.reset()
        updateStatusBarTitle()
    }

    @objc private func selectMode(_ sender: NSMenuItem) {
        if let mode = InputMode(rawValue: sender.tag) {
            engine.config.mode = mode
            engine.reset()
        }
    }

    private func setupEventTap() {
        let mask: CGEventMask =
            (1 << CGEventType.keyDown.rawValue) |
            (1 << CGEventType.leftMouseDown.rawValue) |
            (1 << CGEventType.rightMouseDown.rawValue) |
            (1 << CGEventType.otherMouseDown.rawValue)

        let callback: CGEventTapCallBack = { proxy, type, event, refcon in
            guard let refcon = refcon else { return Unmanaged.passUnretained(event) }
            let controller = Unmanaged<VuiKeyMacController>.fromOpaque(refcon).takeUnretainedValue()
            return controller.handleEvent(proxy: proxy, type: type, event: event)
        }

        let selfPtr = Unmanaged.passUnretained(self).toOpaque()
        guard let tap = CGEvent.tapCreate(
            tap: .cghidEventTap,
            place: .headInsertEventTap,
            options: .defaultTap,
            eventsOfInterest: mask,
            callback: callback,
            userInfo: selfPtr
        ) else {
            NSLog("[VuiKey] Cannot create CGEventTap — please enable Accessibility permission.")
            return
        }

        self.eventTap = tap
        let runLoopSource = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, tap, 0)
        CFRunLoopAddSource(CFRunLoopGetCurrent(), runLoopSource, .commonModes)
        CGEvent.tapEnable(tap: tap, enable: true)
    }

    private func handleEvent(proxy: CGEventTapProxy, type: CGEventType, event: CGEvent) -> Unmanaged<CGEvent>? {
        if type == .tapDisabledByTimeout || type == .tapDisabledByUserInput {
            if let tap = eventTap {
                CGEvent.tapEnable(tap: tap, enable: true)
            }
            return Unmanaged.passUnretained(event)
        }

        if event.getIntegerValueField(.eventSourceUserData) == Self.eventMarker {
            return Unmanaged.passUnretained(event)
        }

        if type == .leftMouseDown || type == .rightMouseDown || type == .otherMouseDown {
            engine.reset()
            return Unmanaged.passUnretained(event)
        }

        guard type == .keyDown else {
            return Unmanaged.passUnretained(event)
        }

        let keyCode = Int(event.getIntegerValueField(.keyboardEventKeycode))
        let flags = event.flags

        // Ctrl + Space toggles V/E
        if keyCode == kVK_Space && flags.contains(.maskControl) && !flags.contains(.maskCommand) {
            toggleEnabled()
            return nil
        }

        if flags.contains(.maskCommand) || flags.contains(.maskControl) || flags.contains(.maskAlternate) {
            engine.reset()
            return Unmanaged.passUnretained(event)
        }

        switch keyCode {
        case kVK_Delete:
            engine.handleBackspace()
            return Unmanaged.passUnretained(event)
        case kVK_ForwardDelete:
            engine.handleDelete()
            return Unmanaged.passUnretained(event)
        case kVK_Space:
            if let act = engine.checkMacroOnSpace() {
                sendBackspacesAndText(backspaces: act.backspaces, text: act.commitText + " ")
                return nil
            }
            return Unmanaged.passUnretained(event)
        case kVK_Return, kVK_Tab, kVK_Escape, kVK_LeftArrow, kVK_RightArrow, kVK_UpArrow, kVK_DownArrow:
            engine.reset()
            return Unmanaged.passUnretained(event)
        default:
            break
        }

        var length = 0
        var chars = [UniChar](repeating: 0, count: 4)
        event.keyboardGetUnicodeString(maxStringLength: 4, actualStringLength: &length, unicodeString: &chars)
        guard length > 0, let scalar = UnicodeScalar(chars[0]) else {
            return Unmanaged.passUnretained(event)
        }

        let act = engine.processKey(Character(scalar))
        guard act.swallow else {
            return Unmanaged.passUnretained(event)
        }

        sendBackspacesAndText(backspaces: act.backspaces, text: act.commitText)
        return nil
    }

    private func isFrontmostAutocompleteApp() -> Bool {
        guard let bid = NSWorkspace.shared.frontmostApplication?.bundleIdentifier else {
            return false
        }
        return Self.autocompleteBundleIDs.contains(bid)
    }

    private func sendBackspacesAndText(backspaces: Int, text: String) {
        let src = CGEventSource(stateID: .privateState)
        var totalBackspaces = backspaces

        if totalBackspaces > 0 && isFrontmostAutocompleteApp() {
            // Narrow No-Break Space trick for Safari/Chrome/Spotlight URL bar autocomplete
            var fixChar: [UniChar] = [0x202F]
            if let down = CGEvent(keyboardEventSource: src, virtualKey: 0, keyDown: true),
               let up = CGEvent(keyboardEventSource: src, virtualKey: 0, keyDown: false) {
                down.keyboardSetUnicodeString(stringLength: 1, unicodeString: &fixChar)
                up.keyboardSetUnicodeString(stringLength: 1, unicodeString: &fixChar)
                down.setIntegerValueField(.eventSourceUserData, value: Self.eventMarker)
                up.setIntegerValueField(.eventSourceUserData, value: Self.eventMarker)
                down.post(tap: .cghidEventTap)
                up.post(tap: .cghidEventTap)
            }
            totalBackspaces += 1
        }

        for _ in 0..<totalBackspaces {
            if let bsDown = CGEvent(keyboardEventSource: src, virtualKey: CGKeyCode(kVK_Delete), keyDown: true),
               let bsUp = CGEvent(keyboardEventSource: src, virtualKey: CGKeyCode(kVK_Delete), keyDown: false) {
                bsDown.setIntegerValueField(.eventSourceUserData, value: Self.eventMarker)
                bsUp.setIntegerValueField(.eventSourceUserData, value: Self.eventMarker)
                bsDown.post(tap: .cghidEventTap)
                bsUp.post(tap: .cghidEventTap)
            }
        }

        if !text.isEmpty {
            var utf16 = Array(text.utf16)
            if let txtDown = CGEvent(keyboardEventSource: src, virtualKey: 0, keyDown: true),
               let txtUp = CGEvent(keyboardEventSource: src, virtualKey: 0, keyDown: false) {
                txtDown.keyboardSetUnicodeString(stringLength: utf16.count, unicodeString: &utf16)
                txtUp.keyboardSetUnicodeString(stringLength: utf16.count, unicodeString: &utf16)
                txtDown.setIntegerValueField(.eventSourceUserData, value: Self.eventMarker)
                txtUp.setIntegerValueField(.eventSourceUserData, value: Self.eventMarker)
                txtDown.post(tap: .cghidEventTap)
                txtUp.post(tap: .cghidEventTap)
            }
        }
    }
}
#endif
