#if os(iOS)
import UIKit
import VuiKeyEngine

/// Native iOS Custom Keyboard Extension (`UIInputViewController`) powered by `VuiKeyEngine`.
/// Memory footprint < 5 MB (well inside iOS's strict ~50 MB Keyboard Extension limit).
public final class VuiKeyKeyboardViewController: UIInputViewController {

    private let engine = VuiKeyEngine()
    private let haptic = UIImpactFeedbackGenerator(style: .light)

    private var isShifted: Bool = true
    private var isCapsLocked: Bool = false
    private var isNumberLayer: Bool = false
    private var lastShiftTapTime: TimeInterval = 0

    private let topBar = UIStackView()
    private let modeButton = UIButton(type: .system)
    private let langButton = UIButton(type: .system)
    private let syllablePreviewLabel = UILabel()
    private let rowsStack = UIStackView()

    private let letterRows: [[String]] = [
        ["q", "w", "e", "r", "t", "y", "u", "i", "o", "p"],
        ["a", "s", "d", "f", "g", "h", "j", "k", "l"],
        ["SHIFT", "z", "x", "c", "v", "b", "n", "m", "BACKSPACE"],
        ["123", "GLOBE", "SPACE", ".", "RETURN"]
    ]

    private let numberRows: [[String]] = [
        ["1", "2", "3", "4", "5", "6", "7", "8", "9", "0"],
        ["-", "/", ":", ";", "(", ")", "$", "&", "@", "\""],
        ["ABC", ".", ",", "?", "!", "'", "BACKSPACE"],
        ["ABC", "GLOBE", "SPACE", ".", "RETURN"]
    ]

    public override func viewDidLoad() {
        super.viewDidLoad()
        haptic.prepare()
        setupKeyboardUI()
    }

    public override func textWillChange(_ textInput: UITextInput?) {
        super.textWillChange(textInput)
        // Reset syllable when user taps to move cursor or switches text fields
        engine.reset()
        updateTopBar()
    }

    private func setupKeyboardUI() {
        view.backgroundColor = UIColor.secondarySystemBackground

        let mainStack = UIStackView()
        mainStack.axis = .vertical
        mainStack.spacing = 6
        mainStack.translatesAutoresizingMaskIntoConstraints = false
        view.addSubview(mainStack)

        NSLayoutConstraint.activate([
            mainStack.leadingAnchor.constraint(equalTo: view.leadingAnchor, constant: 4),
            mainStack.trailingAnchor.constraint(equalTo: view.trailingAnchor, constant: -4),
            mainStack.topAnchor.constraint(equalTo: view.topAnchor, constant: 6),
            mainStack.bottomAnchor.constraint(equalTo: view.bottomAnchor, constant: -4),
            view.heightAnchor.constraint(equalToConstant: 256)
        ])

        // Top toolbar (VN/EN toggle, Current Syllable indicator, Telex/VNI selector)
        topBar.axis = .horizontal
        topBar.distribution = .equalSpacing
        topBar.alignment = .center
        topBar.heightAnchor.constraint(equalToConstant: 32).isActive = true

        langButton.titleLabel?.font = .systemFont(ofSize: 13, weight: .semibold)
        langButton.addTarget(self, action: #selector(onToggleLanguage), for: .touchUpInside)

        syllablePreviewLabel.font = .systemFont(ofSize: 14, weight: .medium)
        syllablePreviewLabel.textColor = .secondaryLabel
        syllablePreviewLabel.textAlignment = .center

        modeButton.titleLabel?.font = .systemFont(ofSize: 13, weight: .medium)
        modeButton.addTarget(self, action: #selector(onCycleMode), for: .touchUpInside)

        topBar.addArrangedSubview(langButton)
        topBar.addArrangedSubview(syllablePreviewLabel)
        topBar.addArrangedSubview(modeButton)
        mainStack.addArrangedSubview(topBar)

        rowsStack.axis = .vertical
        rowsStack.spacing = 8
        rowsStack.distribution = .fillEqually
        mainStack.addArrangedSubview(rowsStack)

        rebuildKeys()
        updateTopBar()
    }

    private func updateTopBar() {
        langButton.setTitle(engine.config.enabled ? "🇻🇳 Tiếng Việt" : "🇬🇧 English", for: .normal)
        modeButton.setTitle("⚙︎ \(engine.config.mode.displayName)", for: .normal)
        let cur = engine.currentSyllable
        syllablePreviewLabel.text = cur.isEmpty ? "VuiKey" : cur
    }

    private func rebuildKeys() {
        rowsStack.arrangedSubviews.forEach { $0.removeFromSuperview() }
        let activeRows = isNumberLayer ? numberRows : letterRows

        for row in activeRows {
            let rowStack = UIStackView()
            rowStack.axis = .horizontal
            rowStack.spacing = 5
            rowStack.distribution = .fillProportionally

            for key in row {
                let btn = makeKeyButton(for: key)
                rowStack.addArrangedSubview(btn)
            }
            rowsStack.addArrangedSubview(rowStack)
        }
    }

    private func makeKeyButton(for key: String) -> UIButton {
        let btn = UIButton(type: .system)
        btn.layer.cornerRadius = 6
        btn.layer.shadowColor = UIColor.black.cgColor
        btn.layer.shadowOpacity = 0.15
        btn.layer.shadowOffset = CGSize(width: 0, height: 1)
        btn.layer.shadowRadius = 0.5

        var title = key
        var isSpecial = false

        switch key {
        case "SHIFT":
            title = isCapsLocked ? "⇪" : (isShifted ? "⇧" : "⇧")
            isSpecial = true
        case "BACKSPACE":
            title = "⌫"
            isSpecial = true
        case "123", "ABC":
            isSpecial = true
        case "GLOBE":
            title = "🌐"
            isSpecial = true
        case "SPACE":
            title = engine.config.enabled ? "Tiếng Việt • VuiKey" : "English"
            // Spacebar cursor trackpad gesture
            let pan = UIPanGestureRecognizer(target: self, action: #selector(onSpacebarPan(_:)))
            btn.addGestureRecognizer(pan)
        case "RETURN":
            title = "Nhập"
            isSpecial = true
        default:
            if !isNumberLayer {
                title = (isShifted || isCapsLocked) ? key.uppercased() : key.lowercased()
            }
        }

        btn.setTitle(title, for: .normal)
        btn.titleLabel?.font = .systemFont(ofSize: isSpecial ? 15 : 20, weight: isSpecial ? .medium : .regular)
        btn.setTitleColor(.label, for: .normal)
        btn.backgroundColor = isSpecial ? UIColor.systemGray4 : UIColor.systemBackground
        btn.accessibilityIdentifier = key

        if key == "GLOBE" {
            btn.addTarget(self, action: #selector(handleInputModeList(from:with:)), for: .allTouchEvents)
        } else {
            btn.addTarget(self, action: #selector(onKeyTapped(_:)), for: .touchUpInside)
        }

        return btn
    }

    @objc private func onToggleLanguage() {
        haptic.impactOccurred()
        engine.config.enabled.toggle()
        engine.reset()
        updateTopBar()
        rebuildKeys()
    }

    @objc private func onCycleMode() {
        haptic.impactOccurred()
        let nextRaw = (engine.config.mode.rawValue + 1) % InputMode.allCases.count
        engine.config.mode = InputMode(rawValue: nextRaw) ?? .telex
        engine.reset()
        updateTopBar()
    }

    private var panAccumulator: CGFloat = 0
    @objc private func onSpacebarPan(_ gesture: UIPanGestureRecognizer) {
        let translation = gesture.translation(in: view)
        if gesture.state == .began {
            panAccumulator = 0
            engine.reset()
        } else if gesture.state == .changed {
            let step: CGFloat = 10.0
            let delta = Int((translation.x - panAccumulator) / step)
            if delta != 0 {
                textDocumentProxy.adjustTextPosition(byCharacterOffset: delta)
                panAccumulator += CGFloat(delta) * step
                haptic.impactOccurred()
            }
        }
    }

    @objc private func onKeyTapped(_ sender: UIButton) {
        guard let key = sender.accessibilityIdentifier else { return }
        haptic.impactOccurred()

        switch key {
        case "SHIFT":
            let now = CACurrentMediaTime()
            if now - lastShiftTapTime < 0.3 {
                isCapsLocked = true
                isShifted = true
            } else {
                isCapsLocked = false
                isShifted.toggle()
            }
            lastShiftTapTime = now
            rebuildKeys()
            return

        case "123":
            isNumberLayer = true
            rebuildKeys()
            return

        case "ABC":
            isNumberLayer = false
            rebuildKeys()
            return

        case "BACKSPACE":
            textDocumentProxy.deleteBackward()
            engine.handleBackspace()
            updateTopBar()
            return

        case "SPACE":
            if let act = engine.checkMacroOnSpace() {
                for _ in 0..<act.backspaces {
                    textDocumentProxy.deleteBackward()
                }
                textDocumentProxy.insertText(act.commitText + " ")
            } else {
                textDocumentProxy.insertText(" ")
            }
            updateTopBar()
            return

        case "RETURN":
            engine.reset()
            textDocumentProxy.insertText("\n")
            isShifted = true
            rebuildKeys()
            updateTopBar()
            return

        default:
            break
        }

        let charStr = (!isNumberLayer && (isShifted || isCapsLocked)) ? key.uppercased() : key.lowercased()
        guard let ch = charStr.first else { return }

        if engine.currentSyllable.isEmpty,
           (textDocumentProxy.selectedText ?? "").isEmpty,
           !((textDocumentProxy.documentContextAfterInput ?? "").first?.isLetter ?? false),
           let before = textDocumentProxy.documentContextBeforeInput,
           let surAct = engine.tryModifySurroundingWord(before, ch: ch),
           surAct.swallow {
            for _ in 0..<surAct.backspaces {
                textDocumentProxy.deleteBackward()
            }
            textDocumentProxy.insertText(surAct.commitText)
            if isShifted && !isCapsLocked {
                isShifted = false
                rebuildKeys()
            }
            updateTopBar()
            return
        }

        let act = engine.processKey(ch)
        if act.swallow {
            for _ in 0..<act.backspaces {
                textDocumentProxy.deleteBackward()
            }
            textDocumentProxy.insertText(act.commitText)
        } else {
            textDocumentProxy.insertText(String(ch))
        }

        if isShifted && !isCapsLocked {
            isShifted = false
            rebuildKeys()
        }
        updateTopBar()
    }
}
#endif
