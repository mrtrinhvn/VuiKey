# VuiKey — AI Agent Instructions (`AGENTS.md`)

**VuiKey** is a standalone, Zero-Preedit Vietnamese Input Method (`Telex`, `VNI`, `Telex + VNI`) for Linux Desktop (`X11` / Ubuntu / Pop!_OS / Debian / Mint / Arch / Fedora), created by `mrtrinhvn` as part of the **Vui Ecosystem** ([VuiCaHat.com](https://vuicahat.com), [DanCaQuanHo.com](https://dancaquanho.com), [VuiRemote](https://github.com/mrtrinhvn/VuiRemote)).

## Autonomous Build & Deployment Workflow

When a user asks an AI Agent (Google Antigravity, Claude Code, OpenAI Codex, Cursor) to install or build **VuiKey**:

1. **Verify Prerequisites (Ubuntu / Pop!_OS / Debian)**:
   ```bash
   sudo apt-get update && sudo apt-get install -y gcc pkg-config libx11-dev libxtst-dev libxi-dev libgtk-3-dev libayatana-appindicator3-dev
   ```
2. **Build & Run Automated Verification (`5/5` Zero-Preedit Tests)**:
   ```bash
   go test -v ./internal/engine/...
   mkdir -p bin
   go build -o bin/vuikey ./cmd/vuikey
   ./bin/vuikey self-test
   ```
3. **Install System-Wide CLI, Desktop Entry & Autostart**:
   ```bash
   ./bin/vuikey install
   ```
4. **Launch VuiKey Studio / Daemon**:
   ```bash
   nohup ./bin/vuikey gui >/dev/null 2>&1 &
   ```
