#!/usr/bin/env bash
set -e

echo "╔══════════════════════════════════════════════════════════════════════╗"
echo "║      🚀 VUIKEY UNIVERSAL LINUX INSTALLER (eChip Co., Ltd.)          ║"
echo "╚══════════════════════════════════════════════════════════════════════╝"

# 1. Detect Linux Distribution & Install Native Dependencies
if command -v apt-get >/dev/null 2>&1; then
    echo "📦 Detected Debian / Ubuntu / Pop!_OS / Linux Mint / Zorin / Kali (apt)..."
    sudo apt-get update
    sudo apt-get install -y gcc g++ pkg-config fcitx5 fcitx5-frontend-gtk3 fcitx5-frontend-gtk4 fcitx5-frontend-qt5 libx11-dev libxtst-dev libxi-dev libgtk-3-dev libayatana-appindicator3-dev libfcitx5core-dev libfcitx5config-dev libfcitx5utils-dev
elif command -v dnf >/dev/null 2>&1; then
    echo "📦 Detected Fedora / RHEL / Rocky / AlmaLinux (dnf)..."
    sudo dnf install -y gcc gcc-c++ pkgconf-pkg-config fcitx5 fcitx5-gtk fcitx5-qt libX11-devel libXtst-devel libXi-devel gtk3-devel libayatana-appindicator-gtk3-devel fcitx5-devel
elif command -v pacman >/dev/null 2>&1; then
    echo "📦 Detected Arch Linux / Manjaro / EndeavourOS (pacman)..."
    sudo pacman -Sy --needed --noconfirm base-devel pkgconf fcitx5 fcitx5-gtk fcitx5-qt libx11 libxtst libxi gtk3 libayatana-appindicator
elif command -v zypper >/dev/null 2>&1; then
    echo "📦 Detected openSUSE (zypper)..."
    sudo zypper install -y gcc gcc-c++ pkg-config fcitx5 libX11-devel libXtst-devel libXi-devel gtk3-devel libayatana-appindicator3-devel fcitx5-devel
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$SCRIPT_DIR"

if command -v go >/dev/null 2>&1; then
    GO_BIN="go"
elif [ -x "$HOME/.local/go/bin/go" ]; then
    GO_BIN="$HOME/.local/go/bin/go"
fi

mkdir -p build bin ~/.local/lib/fcitx5 ~/.local/share/fcitx5/addon ~/.local/share/fcitx5/inputmethod

if [ -n "$GO_BIN" ]; then
    echo "🔨 Building VuiKey Core & Studio UI..."
    "$GO_BIN" build -buildmode=c-archive -o build/libvuikey_core.a ./cmd/libvuikey_core
    "$GO_BIN" build -o bin/vuikey ./cmd/vuikey
fi

INC_FLAGS="-I./build"
if [ -d "./third_party/fcitx5/usr/include/Fcitx5" ]; then
    INC_FLAGS="$INC_FLAGS -I./third_party/fcitx5/usr/include/Fcitx5/Core -I./third_party/fcitx5/usr/include/Fcitx5/Config -I./third_party/fcitx5/usr/include/Fcitx5/Utils"
else
    INC_FLAGS="$INC_FLAGS -I/usr/include/Fcitx5/Core -I/usr/include/Fcitx5/Config -I/usr/include/Fcitx5/Utils"
fi

echo "🔨 Building Native Fcitx5 Engine (libvuikey-fcitx5.so)..."
g++ -shared -fPIC -O2 -std=c++17 \
    $INC_FLAGS \
    fcitx5/vuikey_fcitx5.cpp \
    build/libvuikey_core.a \
    $(pkg-config --libs Fcitx5Core Fcitx5Config Fcitx5Utils 2>/dev/null || echo "/usr/lib/x86_64-linux-gnu/libFcitx5Core.so.7 /usr/lib/x86_64-linux-gnu/libFcitx5Config.so.6 /usr/lib/x86_64-linux-gnu/libFcitx5Utils.so.2") \
    -lpthread -ldl \
    -o ~/.local/lib/fcitx5/libvuikey-fcitx5.so

cat << 'EOF' > ~/.local/share/fcitx5/addon/vuikey.conf
[Addon]
Name[vi]=VuiKey (Bộ gõ Tiếng Việt thế hệ mới - eChip Co., Ltd.)
Name=VuiKey Vietnamese Engine
Category=InputMethod
Version=1.0.0
Enabled=True
Library=libvuikey-fcitx5
Type=SharedLibrary
OnDemand=False
Configurable=False
EOF

cat << 'EOF' > ~/.local/share/fcitx5/inputmethod/vuikey.conf
[InputMethod]
Name[vi]=VuiKey (Tiếng Việt)
Name=VuiKey (Tiếng Việt)
Icon=fcitx-unikey
LangCode=vi
Addon=vuikey
Configurable=False
EOF

./bin/vuikey self-test
./bin/vuikey install

echo "🎉 Hoàn tất cài đặt VuiKey! Khởi chạy bằng lệnh: vuikey gui"
