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
rm -rf "$PKG_DIR" "$TAR_STAGE"

ls -lh "$SCRIPT_DIR/dist/"
