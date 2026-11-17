#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$SCRIPT_DIR"

VERSION="1.0.0"
ARCH="amd64"
PKG_DIR="$SCRIPT_DIR/dist/vuikey_${VERSION}_${ARCH}"
mkdir -p "$SCRIPT_DIR/dist"
rm -rf "$PKG_DIR"

mkdir -p "$PKG_DIR/DEBIAN"
mkdir -p "$PKG_DIR/usr/bin"
mkdir -p "$PKG_DIR/usr/lib/x86_64-linux-gnu/fcitx5"
mkdir -p "$PKG_DIR/usr/share/fcitx5/addon"
mkdir -p "$PKG_DIR/usr/share/fcitx5/inputmethod"
mkdir -p "$PKG_DIR/usr/share/applications"

cp bin/vuikey "$PKG_DIR/usr/bin/vuikey"
chmod 755 "$PKG_DIR/usr/bin/vuikey"

cp ~/.local/lib/fcitx5/libvuikey-fcitx5.so "$PKG_DIR/usr/lib/x86_64-linux-gnu/fcitx5/libvuikey-fcitx5.so"
chmod 644 "$PKG_DIR/usr/lib/x86_64-linux-gnu/fcitx5/libvuikey-fcitx5.so"

cp ~/.local/share/fcitx5/addon/vuikey.conf "$PKG_DIR/usr/share/fcitx5/addon/vuikey.conf"
cp ~/.local/share/fcitx5/inputmethod/vuikey.conf "$PKG_DIR/usr/share/fcitx5/inputmethod/vuikey.conf"

cat << 'EOF' > "$PKG_DIR/usr/share/applications/vuikey.desktop"
[Desktop Entry]
Type=Application
Name=VuiKey Studio (Bộ Gõ Tiếng Việt)
GenericName=Vietnamese Input Method
Comment=Bộ gõ Tiếng Việt thế hệ mới phát triển bởi Công ty TNHH eChip (echipco.com) & VuiCaHat.com
Exec=/usr/bin/vuikey gui
Icon=input-keyboard
Terminal=false
Categories=Utility;Settings;
Keywords=vietnamese;unikey;telex;vni;vuikey;vuicahat;echip;
EOF

cat << EOF > "$PKG_DIR/DEBIAN/control"
Package: vuikey
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Depends: fcitx5, fcitx5-frontend-gtk3, fcitx5-frontend-gtk4, fcitx5-frontend-qt5, libgtk-3-0, libayatana-appindicator3-1
Maintainer: Công ty TNHH eChip <contact@echipco.com>
Homepage: https://echipco.com
Description: VuiKey — Next-Generation Vietnamese Input Method for Linux
 Developed by Công ty TNHH eChip (https://echipco.com) & VuiCaHat (https://vuicahat.com).
 Features Sliding Syllable Window, Anti-Duplicate Delete Guard, and Phonotactic Smart English Guard.
EOF

dpkg-deb --build "$PKG_DIR" "$SCRIPT_DIR/dist/vuikey_${VERSION}_${ARCH}.deb"

# Also create universal .tar.gz bundle for Arch / Fedora / openSUSE
TAR_STAGE="$SCRIPT_DIR/dist/vuikey-linux-x86_64"
rm -rf "$TAR_STAGE"
mkdir -p "$TAR_STAGE"
cp bin/vuikey "$TAR_STAGE/"
cp ~/.local/lib/fcitx5/libvuikey-fcitx5.so "$TAR_STAGE/"
cp ~/.local/share/fcitx5/addon/vuikey.conf "$TAR_STAGE/addon-vuikey.conf"
cp ~/.local/share/fcitx5/inputmethod/vuikey.conf "$TAR_STAGE/im-vuikey.conf"
cp README.md LICENSE NOTICE "$TAR_STAGE/"
tar -czf "$SCRIPT_DIR/dist/vuikey-${VERSION}-linux-x86_64.tar.gz" -C "$SCRIPT_DIR/dist" "vuikey-linux-x86_64"

# Also build standalone Windows x64 Portable bundle (VuiKey.exe)
WIN_STAGE="$SCRIPT_DIR/dist/VuiKey-Windows-x64"
rm -rf "$WIN_STAGE"
mkdir -p "$WIN_STAGE"
GO_BIN="/home/tao/.local/go/bin/go"
if ! [ -x "$GO_BIN" ]; then GO_BIN="go"; fi
GOOS=windows GOARCH=amd64 "$GO_BIN" build -ldflags="-H=windowsgui -s -w" -o "$WIN_STAGE/VuiKey.exe" ./cmd/vuikey_windows
cp README.md LICENSE NOTICE "$WIN_STAGE/"
cat << 'WEOF' > "$WIN_STAGE/HUONG-DAN-SU-DUNG-WINDOWS.txt"
======================================================================
   VUIKEY STUDIO v1.0.0 — BỘ GÕ TIẾNG VIỆT THẾ HỆ MỚI CHO WINDOWS
   Phát triển bởi Công ty TNHH eChip (echipco.com) & VuiCaHat.com
======================================================================

1. CHẠY NGAY KHÔNG CẦN CÀI ĐẶT (PORTABLE):
   - Nhấp đúp vào file "VuiKey.exe" để mở bảng điều khiển VuiKey Studio
     và biểu tượng [VN] màu xanh ngọc dưới khay hệ thống (System Tray).

2. PHÍM TẮT & THAO TÁC NHANH:
   - Nhấn Ctrl + Space (hoặc Alt + Z), hoặc click chuột trái vào biểu
     tượng [VN] ở góc phải dưới màn hình để chuyển đổi Tiếng Việt [VN]
     <-> Tiếng Anh [EN].
   - Nhấp chuột phải vào biểu tượng [VN] để mở nhanh menu chọn Telex,
     VNI, hoặc bật "Khởi động cùng Windows".

3. TÍNH NĂNG VƯỢT TRỘI:
   - Gõ tự do (Free-style Telex): Hỗ trợ đầy đủ cả 5 cách gõ "ươ"
     (dduowcj, dduocwj, dduwowcj, dduwocwj, dduwocj -> được) và gõ chữ
     "đ" ở bất kỳ vị trí nào trong âm tiết (dodongj, doongdj, doongjd -> động).
   - Sliding Syllable Window: Tự động tách âm tiết thông minh ngay cả
     khi gõ liền không dấu cách, chống lỗi nhân đôi chữ khi nhấn Delete.
   - Smart English Guard: Tự động nhận diện và giữ nguyên các từ Tiếng
     Anh thông dụng (window, view, function, test, server, linux, docker...)
     mà không cần tắt bộ gõ Tiếng Việt.

Website chính thức: https://echipco.com/vuikey/
                    https://vuicahat.com
======================================================================
WEOF
python3 -c "
import shutil
shutil.make_archive('$SCRIPT_DIR/dist/VuiKey-${VERSION}-windows-amd64', 'zip', '$SCRIPT_DIR/dist', 'VuiKey-Windows-x64')
"
rm -rf "$PKG_DIR" "$TAR_STAGE" "$WIN_STAGE"

ls -lh "$SCRIPT_DIR/dist/"
