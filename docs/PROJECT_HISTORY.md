# Lịch Sử Phát Triển & Kiến Trúc Dự Án VuiKey (`PROJECT_HISTORY.md`)

Tài liệu này lưu trữ toàn bộ lịch sử thiết kế, các quyết định kiến trúc kỹ thuật và quá trình phát triển **VuiKey Studio (v1.0.0)** từ phiên làm việc gốc trên Antigravity (`conversation://b96e460b-37e5-4d60-bbd7-97fa242a1d42`).

---

## 1. Mục Tiêu Cốt Lõi Khi Khởi Tạo Dự Án
- Giải quyết triệt để lỗi nhân đôi chuỗi ký tự dài không dấu cách (`nhacdidangnhactiengdidangtiengchoitai` ➔ nhảy đúp khi ấn `Delete`/`Backspace`), lỗi rớt chữ khi gõ lướt nhanh (`mà thành khác nữa` ➔ `mà thành khá ữa`) và hiện tượng khựng giao diện (`40–80ms`) trên Linux Desktop khi dùng bàn phím cơ/Bluetooth.
- Xây dựng **VuiKey** thành bộ gõ Tiếng Việt mã nguồn mở thế hệ mới cho cả **Linux (`X11` & `Wayland`)** và **Windows 10/11 (`Portable .exe`)**, phát triển bởi **Công ty TNHH eChip (`https://echipco.com`)** và **VuiCaHat (`https://vuicahat.com`)**.

---

## 2. Các Đột Phá Kỹ Thuật Đã Hoàn Thiện Trong `v1.0.0`

### 2.1. Cửa Sổ Âm Tiết Trượt (`Sliding Syllable Window` — `internal/engine/phonology.go`)
- Tự động phân tích âm vị học Tiếng Việt theo thời gian thực để ngắt và chốt từng âm tiết ngay cả khi gõ chuỗi dài không dấu cách:
  - `nhacdidangnhactiengdidangtiengchoitai` ➔ `nhac | di | dang | nhac | tieng | di | dang | tieng | choi | tai`
  - `nhacjdiujdangfnhacjtieengschoitai` ➔ `nhạcdịudàngnhạctiếngchoitai`
- Khi nhấn `Delete`, `BackSpace`, `Return`, `Tab` hoặc phím mũi tên: xóa sạch `ClientPreedit` trước khi chốt chữ và tuần tự hóa phím điều hướng qua `forwardKey()`, triệt tiêu 100% lỗi nhân đôi chữ.

### 2.2. Kiểu Gõ Tự Do Toàn Diện (`Free-Style Telex / VNI`)
- **Chữ `đ` tự do (`delayedD`):** Cho phép gõ chữ `d` thứ hai ở bất kỳ vị trí nào trong âm tiết (`dodongj`, `doongdj`, `doongjd`, `dongodj`, `dongojd` ➔ **`động`**; `duocwjd`, `duowcjd` ➔ **`được`**), đồng thời tự động phân biệt khi nào chữ `d` bắt đầu một âm tiết mới trong chuỗi dính liền (`nhacjdiujdangf` ➔ **`nhạcdịudàng`**).
- **5 thứ tự gõ nguyên âm đôi `ươ`:** Hỗ trợ `dduowcj`, `dduocwj`, `dduwowcj`, `dduwocwj`, `dduwocj` ➔ **`được`**; và `phuwcs`, `phucws`, `phucsw`, `phwcs`, `phwucs` ➔ **`phức`**.

### 2.3. Mở Khóa Âm Tiết Tức Thì Khi Bấm `Backspace` (`internal/engine/engine.go`)
- Trong `HandleBackspace()`, cờ `e.wordSealed = false` được mở khóa ngay ở nhát `Backspace` đầu tiên, cho phép sửa lỗi gõ nhầm và bỏ dấu lại ngay (`phuwcf` + `Backspace` + `cs` ➔ **`phức`**) mà không phải xóa lùi về đầu từ.

### 2.4. Tự Động Bảo Vệ Từ Tiếng Anh (`Phonotactic Smart English Guard`)
- Kết hợp luật âm vị học `VCV` (`Nguyên âm + Phụ âm + Nguyên âm`) và từ điển kỹ thuật (`window`, `view`, `function`, `test`, `user`, `server`, `linux`, `docker`, `arrow`...) để tự động giữ nguyên từ Tiếng Anh ngay khi đang bật chế độ Tiếng Việt `[VN]`.

### 2.5. Kiến Trúc Đa Nền Tảng (Linux & Windows)
- **Linux (`fcitx5/vuikey_fcitx5.cpp` + `cmd/libvuikey_core/main.go`):**
  - Đón sự kiện phím ở tầng ưu tiên cao nhất `EventWatcherPhase::PreInputMethod` để mọi cửa sổ (`ibus` và `dbus` frontend) luôn đồng bộ 100% trạng thái `[VN]` / `[EN]` từ `~/.config/vuikey/config.json`.
  - Tắt các addon vẽ icon phụ (`notificationitem`, `classicui`, `kimpanel`) để trên khay hệ thống chỉ hiển thị duy nhất 1 biểu tượng **`[VN]`** xanh ngọc của **VuiKey Studio**.
- **Windows (`cmd/vuikey_windows/main_windows.go`):**
  - Thực thi độc lập (`VuiKey.exe` Portable `1.9 MB`, không phụ thuộc thư viện ngoài) sử dụng Win32 `WH_KEYBOARD_LL` + `SendInput(KEYEVENTF_UNICODE)` với nhãn `VUIKEY_EXTRA_INFO` chống lặp sự kiện.

---

## 3. Liên Kết Phát Hành Chính Thức
- **GitHub Repository:** https://github.com/mrtrinhvn/VuiKey
- **GitHub Release `v1.0.0`:** https://github.com/mrtrinhvn/VuiKey/releases/tag/v1.0.0
  - `vuikey_1.0.0_amd64.deb` (Ubuntu / Pop!_OS / Linux Mint / Debian)
  - `vuikey-1.0.0-linux-x86_64.tar.gz` (Fedora / Arch / openSUSE)
  - `VuiKey-1.0.0-windows-amd64.zip` (Windows 10 / 11 64-bit Portable)
- **Trang chủ sản phẩm:** https://echipco.com/vuikey/
