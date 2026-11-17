import XCTest
@testable import VuiKeyEngine

final class VuiKeyEngineTests: XCTestCase {

    private func simulateEditor(_ eng: VuiKeyEngine, _ keystrokes: String) -> String {
        var screen: [Character] = []
        for ch in keystrokes {
            if ch == "\u{0008}" { // Backspace
                if !screen.isEmpty { screen.removeLast() }
                eng.handleBackspace()
            } else if ch == "\u{007F}" { // Delete
                eng.handleDelete()
            } else if ch == " " {
                if let act = eng.checkMacroOnSpace() {
                    if act.backspaces <= screen.count {
                        screen.removeLast(act.backspaces)
                    }
                    screen.append(contentsOf: act.commitText)
                }
                screen.append(" ")
            } else {
                let act = eng.processKey(ch)
                if act.swallow {
                    if act.backspaces <= screen.count {
                        screen.removeLast(act.backspaces)
                    }
                    screen.append(contentsOf: act.commitText)
                } else {
                    screen.append(ch)
                }
            }
        }
        return String(screen)
    }

    func testNhacDiDangZeroDuplicateOnDelete() {
        let eng = VuiKeyEngine()
        let input = "nhacdidangnhactiengdidangtiengchoitai\u{007F}"
        let got = simulateEditor(eng, input)
        let want = "nhacdidangnhactiengdidangtiengchoitai"
        XCTAssertEqual(got, want)
        XCTAssertEqual(eng.currentSyllable, "")
    }

    func testSlidingSyllableWindowConcatenatedWithTones() {
        let eng = VuiKeyEngine()
        let input = "nhacjdiujdangfnhacjtieengschoitai"
        let got = simulateEditor(eng, input)
        let want = "nhạcdịudàngnhạctiếngchoitai"
        XCTAssertEqual(got, want)
    }

    func testComprehensiveVietnameseTypingOrders() {
        let eng = VuiKeyEngine()
        let cases: [(String, String)] = [
            ("Tooi e rawngf conf nhieeuf looix vowis cacs chuwx khacs nuwax, maf thanhf khacs nuwax",
             "Tôi e rằng còn nhiều lỗi với các chữ khác nữa, mà thành khác nữa"),
            ("maf thanhf khasc nuawx", "mà thành khác nữa"),
            ("maf thanhf khacs nuwxa", "mà thành khác nữa"),
            ("maf thanhf khacs nuwax", "mà thành khác nữa"),
            ("phieenf phuwcs", "phiền phức"),
            ("phieenf phucws", "phiền phức"),
            ("phieenf phucsw", "phiền phức"),
            ("phieenf phwcs", "phiền phức"),
            ("phieenf phwucs", "phiền phức"),
            ("pl\u{0008}huwcs", "phức"),
            ("phuwcf\u{0008}cs", "phức"),
            ("dodongj", "động"),
            ("doongdj", "động"),
            ("doongjd", "động"),
            ("dongodj", "động"),
            ("dongojd", "động"),
            ("duocwjd", "được"),
            ("duowcjd", "được"),
            ("dduowcj", "được"),
            ("dduocwj", "được"),
            ("dduwowcj", "được"),
            ("dduwocwj", "được"),
            ("dduwocj", "được"),
            ("truwowngf", "trường"),
            ("truowngf", "trường"),
            ("truongwf", "trường"),
            ("truwongf", "trường"),
            ("truwongwf", "trường"),
            ("nguwowif", "người"),
            ("nguowif", "người"),
            ("nguoiwf", "người"),
            ("nguwoif", "người"),
            ("nuwowcs", "nước"),
            ("nuowcs", "nước"),
            ("nuocws", "nước"),
            ("nuwocs", "nước"),
            ("Vieetj Nam", "Việt Nam"),
            ("VIEETJ NAM", "VIỆT NAM"),
            ("DDuwowngf", "Đường"),
            ("DDUWOWNGF", "ĐƯỜNG")
        ]

        for (input, want) in cases {
            eng.reset()
            let got = simulateEditor(eng, input)
            XCTAssertEqual(got, want, "Failed for input: \(input)")
        }
    }

    func testSmartEnglishAndDoubleKeyUndo() {
        let eng = VuiKeyEngine()
        let cases: [(String, String)] = [
            ("window", "window"),
            ("view", "view"),
            ("function", "function"),
            ("stream", "stream"),
            ("class", "class"),
            ("linux", "linux"),
            ("test", "test"),
            ("ddda", "dda"),
            ("tooi", "tôi"),
            ("toooi", "tooi")
        ]
        for (input, want) in cases {
            eng.reset()
            let got = simulateEditor(eng, input)
            XCTAssertEqual(got, want, "Failed for input: \(input)")
        }
    }
}
