# VuiKey Studio 🇻🇳⚡⌨️

[![Language: VI | EN](https://img.shields.io/badge/🌐_Language-Tiếng_Việt_%7C_English-38bdf8?style=flat-square)](#-giới-thiệu-vuikey)
[![Download macOS](https://img.shields.io/badge/🍎_macOS_(M1--M4_%26_Intel)-VuiKey.dmg_/_Universal_.zip-a855f7?style=flat-square)](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-macos-universal.dmg)
[![Download Android](https://img.shields.io/badge/🤖_Android-VuiKey.apk_%26_Play_Store_.aab-22c55e?style=flat-square)](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-android.apk)
[![Download Windows](https://img.shields.io/badge/🪟_Windows_10/11-VuiKey.exe_(Portable_.zip)-0ea5e9?style=flat-square)](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-windows-amd64.zip)
[![Download .deb](https://img.shields.io/badge/🐧_Linux_.deb-(Ubuntu_%7C_Mint_%7C_Pop!__OS)-10b981?style=flat-square)](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/vuikey_1.0.0_amd64.deb)
[![Platform](https://img.shields.io/badge/💻_Hệ_Sinh_Thái-macOS_%7C_iOS_%7C_Android_%7C_Windows_%7C_Linux-f97316?style=flat-square)](#-bảng-tương-thích-hệ-điều-hành-macos--ios--android--windows--linux)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)

<p align="center">
  <b><a href="#-giới-thiệu-vuikey">🇻🇳 Đọc bản Tiếng Việt</a></b> &nbsp;•&nbsp;
  <b><a href="#-english-documentation">🇺🇸 Read in English</a></b> &nbsp;•&nbsp;
  <b><a href="https://echipco.com/vuikey">🌐 Trang chủ VuiKey (echipco.com/vuikey)</a></b> &nbsp;•&nbsp;
  <b><a href="https://echipco.com">🏢 Công ty TNHH eChip (echipco.com)</a></b> &nbsp;•&nbsp;
  <b><a href="https://vuicahat.com">🎤 VuiCaHat.com</a></b>
</p>

<p align="center">
  <img src="docs/assets/vuikey-hero-overview.svg" alt="VuiKey Studio — Bộ Gõ Tiếng Việt Thế Hệ Mới Cho macOS, iOS, Windows & Linux" width="100%">
</p>

---

# 🇻🇳 Giới Thiệu VuiKey (macOS • iOS • Android • Windows • Linux)

> **Bộ gõ tiếng Việt ngày nay không còn là thứ gì đao to búa lớn — nhưng nó lập tức tạo ra khác biệt khi giải quyết dứt điểm những "nỗi đau" tồn tại suốt nhiều năm trên cả hệ sinh thái Apple (Mac, iPhone/iPad), Android, Windows lẫn Linux mà chưa bộ gõ nào xử lý trọn vẹn.**

### ❌ Những "nỗi đau" kinh điển tồn tại nhiều năm ở các bộ gõ khác:
- **Trên Mac & iOS — Gạch chân khó chịu, nuốt chữ đầu & lỗi thanh địa chỉ:** Bộ gõ mặc định trên macOS luôn hiện dấu gạch chân (`preedit`), khi gõ nhanh vào thanh tìm kiếm **Spotlight, Safari, Chrome, Arc** hoặc trong **VS Code, Notion** thường xuyên bị **nuốt mất chữ cái đầu**, nhảy đúp chữ khi có gợi ý URL tự động (`Autocomplete`), hoặc lộn thứ tự phím.
- **Đang gõ Tiếng Việt mà gõ Tiếng Anh là hỏng chữ:** Trên cả Mac, iPhone, Windows và Linux, khi đang bật Tiếng Việt mà gõ xen kẽ `window`, `view`, `test`, `user`, `server`, `linux`, `docker`... là lập tức bị biến thành `ưindoư`, `vieư`, `tét`, `úer`, `lĩnu` — buộc phải xóa đi và bấm đổi ngôn ngữ liên tục.
- **Sửa 1 ký tự phải xóa sạch cả từ:** Lỡ gõ nhầm 1 phím rồi bấm `Backspace` xóa đi, nhưng khi gõ tiếp dấu sắc/huyền/mũ/móc thì bộ gõ "đơ" không nhận nữa, bắt phải xóa lùi về tận chữ cái đầu tiên mới bỏ dấu lại được.
- **Nhân đôi cả cụm chữ & rớt dấu khi gõ nhanh:** Gõ chuỗi dài không dấu cách (tên file, URL, hashtag) rồi bấm `Delete`/`Backspace` bị nhảy đúp cả đoạn chữ; hoặc gõ chữ `đ` và nguyên âm đôi `ươ` lệch thứ tự một nhịp là rớt dấu.

### ⚡ VuiKey giải quyết thế nào? — Nhanh `0.00ms` & Thông minh thực sự trên mọi thiết bị:
1. **Không gạch chân & Triệt tiêu lỗi Spotlight / Safari / Chrome trên Mac (`Zero-Preedit + Autocomplete Fix`):** Can thiệp trực tiếp ở tầng `CGEventTap` (macOS), `WH_KEYBOARD_LL` (Windows), `PreInputMethod` (Linux) và `UITextDocumentProxy` (iOS) — không bao giờ có gạch chân, tự động xử lý vùng bôi đen gợi ý trên trình duyệt để không bao giờ nuốt chữ đầu hay nhảy đúp chữ.
2. **Tự hiểu Tiếng Anh mà không cần tắt bộ gõ (`Smart English Guard`):** Phân tích âm vị học thời gian thực để tự động giữ nguyên 100% từ tiếng Anh (`window`, `test`, `function`, `server`, `linux`...) ngay khi đang bật Tiếng Việt.
3. **Bấm `Backspace` 1 nhịp sửa dấu ăn ngay (`Instant Unseal`):** Xóa lùi 1 ký tự là mở khóa âm tiết ngay lập tức để bạn bỏ dấu hoặc sửa vần tiếp, không bao giờ phải xóa lại từ đầu.
4. **Gõ tự do mọi thứ tự & Tự tách âm tiết viết liền (`Free-Style Telex/VNI + Sliding Syllable Window`):** Cho phép gõ chữ `d` thứ hai ở bất kỳ đâu trong từ (`dodongj`, `doongjd` ➔ **`động`**), nhận trọn **5 cách gõ `ươ`** (`dduowcj`, `dduocwj`, `dduwowcj`, `dduwocwj`, `dduwocj` ➔ **`được`**) và tự động tách từng âm tiết kể cả khi viết dính liền không dấu cách (`nhacjdiujdangf` ➔ **`nhạcdịudàng`**).

*Phát triển bởi [**Công ty TNHH eChip (`echipco.com`)**](https://echipco.com) & [**VuiCaHat (`vuicahat.com`)**](https://vuicahat.com).*

---

## 🚀 Tải Về & Cài Đặt Nhanh (macOS • iOS • Android • Windows • Linux)

### 🍎 Cho macOS (Apple Silicon M1/M2/M3/M4 & Intel — Universal Binary) & iOS
- **Tải bản cài đặt `.dmg` (macOS 12+):** [**`VuiKey-1.0.0-macos-universal.dmg`**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-macos-universal.dmg)
- **Tải bản nén `.zip` (`VuiKey.app`):** [**`VuiKey-1.0.0-macos-universal.zip`**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-macos-universal.zip)
- Kéo **`VuiKey.app`** vào thư mục `Applications`, mở ứng dụng và cấp quyền **Accessibility** (`System Settings ➔ Privacy & Security ➔ Accessibility`). Biểu tượng **`🇻🇳 V` / `🇬🇧 E`** sẽ xuất hiện trên thanh Menu Bar (chuyển nhanh bằng `Ctrl + Space`).
- **Cho iOS (iPhone / iPad) & Xcode:** Mở gói Swift Package tại [`apple/VuiKeyKit`](apple/VuiKeyKit) và Keyboard Extension tại [`apple/iOS/VuiKeyKeyboardViewController.swift`](apple/iOS/VuiKeyKeyboardViewController.swift) (siêu nhẹ `< 2 MB RAM`, hỗ trợ Haptic Feedback & vuốt Spacebar di chuyển con trỏ).

### 🤖 Cho Android (Điện thoại, Máy tính bảng & Google Play Store)
- **Tải file cài đặt trực tiếp (`.apk`):** [**`VuiKey-1.0.0-android.apk`**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-android.apk)
- **Gói xuất bản Google Play Console (`.aab` API 35):** [**`VuiKey-1.0.0-android-playstore.aab`**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-android-playstore.aab) (xem hồ sơ đăng CH Play tại [`android/PLAY_STORE_LISTING.md`](android/PLAY_STORE_LISTING.md)).
- Mở ứng dụng **VuiKey** trên điện thoại ➔ Bấm **Bật VuiKey** ➔ **Chọn VuiKey làm bàn phím mặc định**. Không gạch chân (`Zero-Preedit Batch Edit`), tự hiểu Tiếng Anh, vuốt phím Spacebar di chuyển con trỏ, hỗ trợ cả bàn phím cứng Bluetooth / Samsung DeX và **100% Offline (không yêu cầu quyền Internet)**.

### 🪟 Cho Windows 10 / Windows 11 (64-bit Portable — Chạy ngay, không cần cài đặt)
- **Tải trực tiếp:** [**`VuiKey-1.0.0-windows-amd64.zip`**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-windows-amd64.zip)
- Giải nén và nhấp đúp vào **`VuiKey.exe`** (chỉ `1.9 MB`, độc lập 100%, không cần cài `.NET` hay thư viện phụ trợ, tích hợp sẵn khay hệ thống **`[VN]` / `[EN]`**, phím tắt `Ctrl+Space` và tùy chọn *Khởi động cùng Windows*).

### 🐧 Cách 1 (Linux): Cài từ gói `.deb` dựng sẵn (Ubuntu, Pop!_OS, Linux Mint, Debian, Zorin, Kali)
```bash
wget -O /tmp/vuikey_1.0.0_amd64.deb https://github.com/mrtrinhvn/VuiKey/releases/latest/download/vuikey_1.0.0_amd64.deb
sudo apt install -y /tmp/vuikey_1.0.0_amd64.deb
vuikey gui
```

### 🐧 Cách 2 (Linux): Cài đặt tự động 1 dòng lệnh cho mọi bản Linux (`apt`, `dnf`, `pacman`, `zypper`)
```bash
curl -fsSL https://raw.githubusercontent.com/mrtrinhvn/VuiKey/main/scripts/install-linux.sh | bash
```

### 🤖 Cách 3: Giao cho AI Agent (Google Antigravity / Claude Code / Codex / Cursor)
> *"Hãy clone `https://github.com/mrtrinhvn/VuiKey`, đọc `AGENTS.md`, tự động cài đặt VuiKey và bật bộ gõ cho tôi."*

---

## 🔥 5 Ưu Việt Đột Phá Của VuiKey (Giải Quyết Triệt Để Những Phiền Toái Khi Gõ Tiếng Việt)

### 1. 🛡️ Gõ Tiếng Anh Xen Kẽ Tự Nhiên Mà Không Cần Tắt Bộ Gõ (`Phonotactic Smart English Guard`)
- **Vấn đề thường gặp ở các bộ gõ khác:** Khi đang bật chế độ Tiếng Việt (trên máy tính hoặc điện thoại) mà gõ xen kẽ các từ Tiếng Anh như `window`, `view`, `function`, `test`, `user`, `server`, `linux`, `docker`, `arrow`..., chữ lập tức bị biến dạng thành `ưindoư`, `vieư`, `tét`, `úer`, `sẻver`, `lĩnu`... hoặc bị tự động sửa từ sai nghĩa khiến bạn phải xóa đi gõ lại rất mất thời gian.
- **Cách VuiKey xử lý:** Phân tích cấu trúc âm vị học Tiếng Việt theo thời gian thực (luật `VCV: Nguyên âm + Phụ âm + Nguyên âm` vốn không tồn tại trong đơn âm tiết Tiếng Việt) kết hợp từ điển kỹ thuật để **tự động nhận diện và giữ nguyên 100% từ Tiếng Anh** ngay trong lúc đang bật Tiếng Việt mà không cần bấm `Ctrl + Space`.

### 2. 🔓 Bấm `Backspace` Sửa Chữ Nhận Dấu Ngay Lập Tức (Không Bao Giờ Phải Xóa Lại Từ Đầu)
- **Vấn đề thường gặp ở các bộ gõ khác:** Khi bạn vô tình gõ nhầm 1 phím trong từ (làm từ bị "đóng" hoặc mất trạng thái âm tiết), dù bạn bấm `Backspace` xóa ký tự gõ nhầm đi rồi gõ lại dấu sắc/huyền/nặng/móc vẫn không ăn, buộc phải xóa lùi về tận chữ cái đầu tiên mới gõ lại được.
- **Cách VuiKey xử lý:** Ngay ở nhát bấm `Backspace` đầu tiên, **VuiKey** lập tức mở khóa trạng thái âm tiết và khôi phục ngữ cảnh vần trước đó (ví dụ gõ nhầm `phuwcf` ➔ bấm `Backspace` 1 lần ➔ gõ tiếp `cs` ra ngay **`phức`** mượt mà).

### 3. 🎯 Gõ Tự Do Toàn Diện (`Free-Style Telex / VNI`) — Hỗ Trợ `dodongj` ➔ `động` & Cả 5 Thứ Tự Gõ `ươ`
- **Vấn đề thường gặp ở các bộ gõ khác:** Bắt người dùng phải gõ 2 chữ `dd` liền nhau ngay từ đầu từ, hoặc chỉ nhận đúng 1–2 thứ tự bỏ dấu cho nguyên âm đôi `ươ`, khiến khi gõ nhanh ở tốc độ cao dễ bị rớt dấu mũ/móc hoặc mất chữ (`mà thành khác nữa` bị đứt thành `mà thành khá ữa`).
- **Cách VuiKey xử lý:**
  - Cho phép gõ chữ **`d` thứ hai ở bất kỳ vị trí nào trong âm tiết**: `dodongj`, `doongdj`, `doongjd`, `dongodj`, `dongojd` đều ra chuẩn **`động`**; `duocwjd`, `duowcjd` đều ra **`được`**.
  - Hỗ trợ trọn vẹn cả **5 thói quen gõ nguyên âm đôi `ươ`**: `dduowcj`, `dduocwj`, `dduwowcj`, `dduwocwj`, `dduwocj` ➔ **`được`**; `phuwcs`, `phucws`, `phucsw`, `phwcs`, `phwucs` ➔ **`phức`**.

### 4. ✂️ Cửa Sổ Âm Tiết Trượt Thông Minh (`Sliding Syllable Window`) & Chống Nhân Đôi Chữ 100%
- **Vấn đề thường gặp ở các bộ gõ khác:** Khi bạn gõ các chuỗi dài không dấu cách (như tên file, đường dẫn URL, mật khẩu, hashtag hoặc chuỗi thẻ `nhacdidangnhactiengdidangtiengchoitai`) rồi ấn phím `Delete`, `Backspace` hoặc phím mũi tên, bộ đệm gạch chân bị giữ quá dài sẽ gây xung đột sự kiện làm **nhân đôi cả đoạn chữ** (`nhacdidang...hacdidang...`) trên thanh địa chỉ trình duyệt và các trình soạn thảo hiện đại.
- **Cách VuiKey xử lý:** Tự động nhận biết ranh giới âm tiết Tiếng Việt để ngắt và chốt từng âm tiết ngay khi bạn đang gõ dính liền (`nhac | di | dang | nhac | tieng | di | dang | tieng | choi | tai`). Ngay cả khi gõ dính liền có dấu (`nhacjdiujdangfnhacjtieengschoitai`), VuiKey vẫn phân tách chuẩn xác thành **`nhạcdịudàngnhạctiếngchoitai`** và ấn `Delete` không bao giờ bị nhảy đúp chữ.

### 5. ⚡ Kiến Trúc Native `0.00ms` Trên Cả 4 Hệ Điều Hành (`macOS` • `iOS` • `Windows` • `Linux`)
- Trên **macOS (Apple Silicon & Intel)**, VuiKey bắt phím trực tiếp qua `CGEventTap` (`kCGHIDEventTap`) kết hợp cơ chế **Autocomplete Bar Fix** chuyên biệt cho Safari, Chrome, Arc, Brave, Edge và Spotlight — triệt tiêu hoàn toàn dấu gạch chân (`Zero-Preedit`) và lỗi nuốt chữ cái đầu.
- Trên **iOS (iPhone & iPad)**, gói `VuiKeyKit` viết bằng Pure Swift chiếm chưa tới `2 MB RAM` (an toàn tuyệt đối trong giới hạn bộ nhớ của `UIInputViewController`), tích hợp Haptic Feedback và vuốt phím Spacebar di chuyển con trỏ.
- Trên **Windows 10/11**, VuiKey sử dụng `WH_KEYBOARD_LL` kết hợp gói sự kiện `KEYEVENTF_UNICODE` nguyên tử ở mức nhân Win32 với độ trễ `0.00ms`.
- Trên **Linux**, VuiKey hoạt động ở tầng ưu tiên cao nhất (`PreInputMethod` trong `libvuikey-fcitx5.so`), truyền trực tiếp chuỗi UTF-8 nguyên tử vào ứng dụng mà không gây hiện tượng chuyển đổi thiết bị bàn phím ảo (`XTest Slave Device Switching`), triệt tiêu hoàn toàn độ trễ `40–80ms` trên GNOME Shell / Mutter khi dùng bàn phím rời hoặc Bluetooth.

---

## 📋 Bảng Tương Thích Hệ Điều Hành (macOS • iOS • Android • Windows • Linux)

| Nhóm Hệ Điều Hành | Các phiên bản / Thiết bị hỗ trợ | Kiến trúc lõi (`0.00ms`) | Gói cài đặt khuyên dùng |
| :--- | :--- | :--- | :--- |
| **1. macOS (Apple)** | • **macOS 12 Monterey ➔ macOS 15 Sequoia**<br>• Chạy gốc trên cả **Apple Silicon (`M1/M2/M3/M4`)** & **Intel (`x86_64`)** | ✅ **`CGEventTap` Zero-Preedit**<br>✅ **Safari/Spotlight Fix** | [**`VuiKey-1.0.0-macos-universal.dmg`**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-macos-universal.dmg) hoặc [`.zip`](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-macos-universal.zip) |
| **2. Android (CH Play / APK)** | • **Android 7.0 ➔ Android 15 (API 24–35)**<br>• Điện thoại, Tablet & bàn phím rời Samsung DeX | ✅ **Zero-Preedit Batch Edit**<br>✅ **100% Offline (0 Net Perm)** | [**`VuiKey-1.0.0-android.apk`**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-android.apk) & [**`.aab` (CH Play)**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-android-playstore.aab) |
| **3. iOS / iPadOS (Apple)** | • **iOS 15+ / iPadOS 15+**<br>• Custom Keyboard Extension (`UIInputViewController`) | ✅ **Pure Swift `VuiKeyKit`** (`< 2 MB RAM`) | Mã nguồn Xcode & SwiftPM tại [`apple/`](apple/) |
| **4. Windows 10 / 11** | • **Windows 10 & Windows 11 (64-bit)**<br>• Chạy ngay (Portable), không cần cài đặt | ✅ **Win32 Native (`0.00ms`)** | [**`VuiKey-1.0.0-windows-amd64.zip`**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-windows-amd64.zip) (`VuiKey.exe` Portable) |
| **5. Họ Debian / Ubuntu** | • **Ubuntu** *(20.04, 22.04, 24.04 LTS)*, **Pop!_OS**, **Linux Mint**, **Debian** *(11, 12)*, **Zorin**, **Kali** | ✅ **X11 (Xorg)**<br>✅ **Wayland** | [**`vuikey_1.0.0_amd64.deb`**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/vuikey_1.0.0_amd64.deb) hoặc `scripts/install-linux.sh` |
| **6. Họ Fedora / Arch / SUSE** | • **Fedora Workstation** *(38–41)*, **RHEL / Rocky / Alma**<br>• **Arch Linux**, **Manjaro**, **EndeavourOS**, **openSUSE** | ✅ **X11 & Wayland** *(KDE, GNOME, Hyprland, Sway, i3)* | `scripts/install-linux.sh` hoặc [**`vuikey-1.0.0-linux-x86_64.tar.gz`**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/vuikey-1.0.0-linux-x86_64.tar.gz) |

---

## 💡 Câu Chuyện Ra Đời Của VuiKey — Từ Công ty TNHH eChip (`echipco.com`) & `VuiCaHat.com`

Trong quá trình đội ngũ kỹ thuật [**Công ty TNHH eChip (`echipco.com`)**](https://echipco.com) phát triển hệ sinh thái sản phẩm công nghệ số và nền tảng phòng hát Karaoke trực tuyến & xử lý âm thanh Vocal DSP 48kHz [**VuiCaHat (`vuicahat.com`)**](https://vuicahat.com), việc nhập lời bài hát, gõ các chuỗi thẻ âm thanh viết liền không dấu cách và viết mã nguồn xen kẽ Tiếng Việt — Tiếng Anh trên cả máy **Mac, iPhone, Windows lẫn Linux** đòi hỏi một bộ gõ thật sự thông minh, phản hồi tức thì và không bao giờ phạm lỗi gạch chân nuốt chữ hay nhân đôi chữ.

Từ nhu cầu thực tiễn đó, chúng tôi đã nghiên cứu và xây dựng **VuiKey** với thuật toán âm vị học hoàn toàn mới, đồng thời mở mã nguồn miễn phí (MIT License) cho cả cộng đồng người dùng **macOS, iOS, Windows và Linux**.

- ⭐ Hãy tặng dự án 1 **Star** trên GitHub nếu **VuiKey** giúp công việc hàng ngày của bạn thuận tiện hơn!
- 🎤 Trải nghiệm phòng hát Karaoke & xử lý giọng hát trực tuyến tại [**vuicahat.com**](https://vuicahat.com).
- 🏢 Tìm hiểu thêm về hệ sinh thái sản phẩm công nghệ của **Công ty TNHH eChip** tại [**echipco.com**](https://echipco.com).

---
---

# 🇺🇸 English Documentation

**VuiKey** is a next-generation Vietnamese Input Method (`Telex`, `VNI`, `Telex + VNI`) for **macOS (`Apple Silicon & Intel`)**, **iOS (`iPhone & iPad`)**, **Windows 10 / 11 (`64-bit Portable`)**, and **Linux Desktop (`X11` & `Wayland`)**, engineered by [**eChip Co., Ltd. (`echipco.com`)**](https://echipco.com) and [**VuiCaHat (`vuicahat.com`)**](https://vuicahat.com).

### Key Technical Highlights
1. **Zero-Preedit & macOS Browser/Spotlight Fix:** Eliminates the intrusive preedit underline on macOS and prevents first-character swallowing or duplication in Safari, Chrome, Arc, and Spotlight autocomplete bars.
2. **Phonotactic Smart English Guard:** Automatically preserves English words (`window`, `view`, `function`, `test`, `user`, `server`, `linux`, `docker`, `arrow`...) while Vietnamese mode remains active.
3. **Instant Backspace Unsealing:** Pressing `Backspace` once immediately unseals the syllable state so you can re-apply tone or vowel marks without retyping the entire word.
4. **Full Free-Style Telex & VNI + Sliding Syllable Window:** Supports delayed `d ➔ đ` anywhere in the syllable (`dodongj`, `doongdj`, `doongjd` ➔ **`động`**), all 5 typing orders for `ươ`, and real-time segmentation of unspaced multi-syllable strings.
