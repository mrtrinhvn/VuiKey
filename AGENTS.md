# VuiKey — AI Agent Instructions (`AGENTS.md`)

**VuiKey** is a next-generation Vietnamese Input Method (`Telex`, `VNI`, `Telex + VNI`) for **Linux Desktop** (`X11` & `Wayland` across Ubuntu, Pop!_OS, Linux Mint, Debian, Fedora, Arch, openSUSE) and **Windows 10/11 (64-bit Portable)**, developed by [**Công ty TNHH eChip (`echipco.com`)**](https://echipco.com) and [**VuiCaHat (`vuicahat.com`)**](https://vuicahat.com).

- **Full Project History & Architecture Reference:** See [docs/PROJECT_HISTORY.md](docs/PROJECT_HISTORY.md) (originating from Antigravity conversation `b96e460b-37e5-4d60-bbd7-97fa242a1d42`).

## Key Repository Layout
- `internal/engine/engine.go` & `internal/engine/phonology.go`: Pure-Go Vietnamese linguistic engine (Sliding Syllable Window, Free-style delayed `d ➔ đ`, 5-order `ươ`, Instant Backspace Unseal, Phonotactic Smart English Guard).
- `fcitx5/vuikey_fcitx5.cpp` & `cmd/libvuikey_core/main.go`: Native C++/Go Fcitx5 plugin (`libvuikey-fcitx5.so`) using `EventWatcherPhase::PreInputMethod` for `0.00ms` latency across X11 & Wayland.
- `internal/ui/ui.go` & `cmd/vuikey/main.go`: Linux GTK3 + Ayatana AppIndicator3 GUI (`VuiKey Studio`) and CLI.
- `cmd/vuikey_windows/main_windows.go`: Native Win32 GUI + System Tray + `WH_KEYBOARD_LL` / `KEYEVENTF_UNICODE` portable executable (`VuiKey.exe`) for Windows 10/11.
- `scripts/install-linux.sh`: Universal 1-line installer for Debian/Ubuntu (`apt`), Fedora/RHEL (`dnf`), Arch (`pacman`), and openSUSE (`zypper`).
- `scripts/build-release.sh`: Builds `dist/vuikey_1.0.0_amd64.deb`, `dist/vuikey-1.0.0-linux-x86_64.tar.gz`, and `dist/VuiKey-1.0.0-windows-amd64.zip`.

## Autonomous Build & Verification Workflow
1. **Run Linguistic Engine Unit Tests:**
   ```bash
   go test -v ./internal/engine
   ```
2. **Build All Linux & Windows Release Artifacts:**
   ```bash
   ./scripts/install-linux.sh
   ./scripts/build-release.sh
   ```
