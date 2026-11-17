# VuiKey Studio 🇻🇳⚡⌨️

[![Language: VI | EN](https://img.shields.io/badge/🌐_Language-Tiếng_Việt_%7C_English-38bdf8?style=flat-square)](#-tài-liệu-tiếng-việt-vietnamese-documentation)
[![Download .deb](https://img.shields.io/badge/📦_Linux_.deb-(Ubuntu_%7C_Mint_%7C_Pop!__OS)-10b981?style=flat-square)](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/vuikey_1.0.0_amd64.deb)
[![Download Windows](https://img.shields.io/badge/🪟_Windows_10/11-VuiKey.exe_(Portable_.zip)-0ea5e9?style=flat-square)](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-windows-amd64.zip)
[![AI Agent Native](https://img.shields.io/badge/🤖_AI_Agent_Native-Antigravity_%7C_Claude_Code_%7C_Codex-8b5cf6?style=flat-square)](AGENTS.md)
[![Linux & Windows](https://img.shields.io/badge/💻_OS-Linux_(X11_%26_Wayland)_%7C_Windows_10/11_x64-f97316?style=flat-square)](#-bảng-tương-thích-hệ-điều-hành-linux)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)

<p align="center">
  <b><a href="#-tài-liệu-tiếng-việt-vietnamese-documentation">🇻🇳 Đọc bản Tiếng Việt</a></b> &nbsp;•&nbsp;
  <b><a href="#-english-documentation">🇺🇸 Read in English</a></b> &nbsp;•&nbsp;
  <b><a href="https://echipco.com/vuikey">🌐 Trang chủ VuiKey (echipco.com/vuikey)</a></b> &nbsp;•&nbsp;
  <b><a href="https://echipco.com">🏢 Công ty TNHH eChip (echipco.com)</a></b> &nbsp;•&nbsp;
  <b><a href="https://vuicahat.com">🎤 VuiCaHat.com</a></b>
</p>

<p align="center">
  <img src="docs/assets/vuikey-hero-overview.svg" alt="VuiKey Studio — Bộ Gõ Tiếng Việt Thế Hệ Mới Cho Linux" width="100%">
</p>

---

# 🇻🇳 Tài Liệu Tiếng Việt (Vietnamese Documentation)

**VuiKey** là bộ gõ Tiếng Việt thế hệ mới dành cho Linux Desktop (`Ubuntu`, `Pop!_OS`, `Linux Mint`, `Debian`, `Fedora`, `Arch`, `Manjaro`, `openSUSE` — hỗ trợ cả **X11** và **Wayland**), được phát triển bởi đội ngũ kỹ thuật [**Công ty TNHH eChip (`echipco.com`)**](https://echipco.com) và nền tảng âm nhạc trực tuyến [**VuiCaHat (`vuicahat.com`)**](https://vuicahat.com).

> ⭐ **Nếu bạn từng ức chế vì lỗi nhảy đúp chữ khi ấn `Delete`/`Backspace`, mất chữ khi gõ nhanh hoặc giật lag của `ibus-bamboo` trên Linux, hãy tặng VuiKey 1 ngôi sao (Star) ở góc trên bên phải GitHub để ủng hộ dự án và giúp cộng đồng Linux Việt Nam biết đến nhé!**

---

## 🚀 Cài Đặt Nhanh (Linux & Windows)

### 🪟 Cho Windows 10 / Windows 11 (64-bit Portable — Chạy ngay như UniKey, không cần cài đặt)
- **Tải trực tiếp:** [**`VuiKey-1.0.0-windows-amd64.zip`**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-windows-amd64.zip)
- Giải nén và nhấp đúp vào **`VuiKey.exe`** (chỉ `1.9 MB`, không phụ thuộc thư viện ngoài, tích hợp sẵn khay hệ thống **`[VN]` / `[EN]`**, phím tắt `Ctrl+Space` và tùy chọn *Khởi động cùng Windows*).

### 🐧 Cách 1 (Linux): Cài từ gói `.deb` dựng sẵn (Ubuntu, Pop!_OS, Linux Mint, Debian, Zorin, Kali)
```bash
wget -O /tmp/vuikey_1.0.0_amd64.deb https://github.com/mrtrinhvn/VuiKey/releases/latest/download/vuikey_1.0.0_amd64.deb
sudo apt install -y /tmp/vuikey_1.0.0_amd64.deb
vuikey gui
```

### Cách 2: Cài đặt tự động 1 dòng lệnh cho mọi bản Linux (`apt`, `dnf`, `pacman`, `zypper`)
```bash
curl -fsSL https://raw.githubusercontent.com/mrtrinhvn/VuiKey/main/scripts/install-linux.sh | bash
```

### Cách 3: Giao cho AI Agent (Google Antigravity / Claude Code / Codex / Cursor)
> *"Hãy clone `https://github.com/mrtrinhvn/VuiKey`, đọc `AGENTS.md`, tự động cài đặt VuiKey và bật giao diện `VuiKey Studio` cho tôi."*

---

## 🔥 4 Đột Phá Kỹ Thuật Của VuiKey So Với `IBus-Bamboo` & `Unikey` Cũ

### 1. ⚡ Động Cơ C++/Go Native (`libvuikey-fcitx5.so`) — Độ Trễ `0.00ms`, Không Giật Focus
- Loại bỏ hoàn toàn hiện tượng xung đột chuyển đổi thiết bị bàn phím (`XTest Slave Device Switching`) gây khựng `40–80ms` trên GNOME Shell / Mutter khi dùng bàn phím rời hoặc Bluetooth.
- Tuần tự hóa tuyệt đối các phím điều hướng (`Delete`, `BackSpace`, `Return`, `Tab`, phím mũi tên) trên cùng một luồng FIFO: xóa sạch `ClientPreedit` trước khi chốt chữ, **triệt tiêu 100% căn bệnh nhân đôi chữ** trên Chrome, Brave, Edge, VS Code và Antigravity.

### 2. ✂️ Cửa Sổ Âm Tiết Trượt Thông Minh (`Sliding Syllable Window`)
- Tự động phân tích âm vị học Tiếng Việt theo thời gian thực để ngắt và chốt từng âm tiết ngay cả khi bạn gõ chuỗi dài không dấu cách:
  `nhacdidangnhactiengdidangtiengchoitai` ➔ tự tách thành `nhac | di | dang | nhac | tieng | di | dang | tieng | choi | tai`.
- Gõ dính liền có dấu (`nhacjdiujdangfnhacjtieengschoitai`) vẫn ra chuẩn **`nhạcdịudàngnhạctiếngchoitai`** và ấn `Delete` không bao giờ bị nhảy đúp chữ.

### 3. 🎯 Hỗ Trợ Linh Hoạt Mọi Thứ Tự Gõ Telex / VNI & Mở Khóa Tức Thì Khi Ấn `Backspace`
- Hỗ trợ cả 5 thói quen gõ nguyên âm đôi `ươ` (`dduowcj`, `dduocwj`, `dduwowcj`, `dduwocwj`, `dduwocj` ➔ **`được`**; `phuwcs`, `phucws`, `phucsw`, `phwcs`, `phwucs` ➔ **`phức`**).
- Ngay khi bấm `Backspace` xóa 1 ký tự gõ nhầm, trạng thái âm tiết được mở khóa tức thì để bạn sửa chữ ngay mà không bao giờ phải xóa lùi về tận chữ cái đầu tiên.

### 4. 🛡️ Tự Động Nhận Diện & Giữ Nguyên Từ Tiếng Anh (`Phonotactic Smart English Guard`)
- Áp dụng luật âm vị học `VCV` (`Nguyên âm + Phụ âm + Nguyên âm`) kết hợp bộ từ vựng kỹ thuật để tự động giữ nguyên các từ tiếng Anh phổ biến (`test`, `Linux`, `window`, `view`, `function`, `stream`, `server`, `docker`, `user`, `port`, `left`, `next`...) ngay trong lúc đang bật chế độ Tiếng Việt.

---

## 📋 Bảng Tương Thích Hệ Điều Hành Linux

| Nhóm Hệ Điều Hành | Các bản phân phối (Distros) hỗ trợ | Môi trường hiển thị | Gói cài đặt khuyên dùng |
| :--- | :--- | :--- | :--- |
| **1. Họ Debian / Ubuntu** | • **Ubuntu** *(20.04, 22.04, 24.04 LTS)*<br>• **Pop!_OS** *(22.04, 24.04)*<br>• **Linux Mint** *(Cinnamon, MATE, XFCE)*<br>• **Debian** *(11, 12)*, **Zorin OS**, **Kali Linux**, **Elementary OS** | ✅ **X11 (Xorg)**<br>✅ **Wayland** | [**`vuikey_1.0.0_amd64.deb`**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/vuikey_1.0.0_amd64.deb) hoặc `scripts/install-linux.sh` |
| **2. Họ Red Hat / Fedora** | • **Fedora Workstation** *(38, 39, 40, 41)*<br>• **RHEL / Rocky Linux / AlmaLinux / CentOS Stream** | ✅ **X11 (Xorg)**<br>✅ **Wayland** | `scripts/install-linux.sh` hoặc [**`vuikey-1.0.0-linux-x86_64.tar.gz`**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/vuikey-1.0.0-linux-x86_64.tar.gz) |
| **3. Họ Arch Linux** | • **Arch Linux**<br>• **Manjaro**, **EndeavourOS**, **Garuda Linux** | ✅ **X11 & Wayland** *(KDE, GNOME, Hyprland, Sway, i3)* | `scripts/install-linux.sh` hoặc [**`vuikey-1.0.0-linux-x86_64.tar.gz`**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/vuikey-1.0.0-linux-x86_64.tar.gz) |
| **4. Họ SUSE** | • **openSUSE Leap & Tumbleweed** | ✅ **X11 & Wayland** | `scripts/install-linux.sh` |
| **5. Windows 10 / 11** | • **Windows 10 & Windows 11 (64-bit)** | ✅ **Win32 Native (`WH_KEYBOARD_LL` + `KEYEVENTF_UNICODE`)** | [**`VuiKey-1.0.0-windows-amd64.zip`**](https://github.com/mrtrinhvn/VuiKey/releases/latest/download/VuiKey-1.0.0-windows-amd64.zip) (`VuiKey.exe` Portable) |

---

## 💡 Câu Chuyện Ra Đời Của VuiKey — Từ Công ty TNHH eChip (`echipco.com`) & `VuiCaHat.com`

Trong quá trình đội ngũ kỹ thuật [**Công ty TNHH eChip (`echipco.com`)**](https://echipco.com) phát triển hệ sinh thái sản phẩm công nghệ số và nền tảng phòng hát Karaoke trực tuyến & xử lý âm thanh Vocal DSP 48kHz [**VuiCaHat (`vuicahat.com`)**](https://vuicahat.com) trên các máy trạm Linux, việc nhập lời bài hát, gõ các chuỗi tag âm thanh viết liền (`nhacdidangnhactiengdidangtiengchoitai`) và làm việc với các AI Agent thường xuyên gặp phải lỗi nhân đôi chữ và giật lag của `ibus-bamboo`.

Để giải quyết triệt để bài toán này, chúng tôi đã xây dựng **VuiKey** và quyết định mở hoàn toàn mã nguồn (MIT License) để đóng góp cho cộng đồng người dùng Linux Việt Nam.

- ⭐ Hãy tặng dự án 1 **Star** trên GitHub nếu **VuiKey** giúp công việc hàng ngày của bạn thuận tiện hơn!
- 🎤 Trải nghiệm phòng hát Karaoke & xử lý giọng hát trực tuyến tại [**vuicahat.com**](https://vuicahat.com).
- 🏢 Tìm hiểu thêm về hệ sinh thái sản phẩm công nghệ của **Công ty TNHH eChip** tại [**echipco.com**](https://echipco.com).

---
---

# 🇺🇸 English Documentation

**VuiKey** is a next-generation Vietnamese Input Method (`Telex`, `VNI`, `Telex + VNI`) for Linux Desktop (`X11` and `Wayland` across Ubuntu, Pop!_OS, Linux Mint, Debian, Fedora, Arch, openSUSE), engineered in **Go + Native C/C++ (`libvuikey-fcitx5.so`, `GTK3`, `Ayatana AppIndicator3`)** by [**eChip Co., Ltd. (`echipco.com`)**](https://echipco.com) and [**VuiCaHat (`vuicahat.com`)**](https://vuicahat.com).
