#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "=========================================================="
echo "🍎 VuiKey Apple Ecosystem Builder (macOS & iOS)"
echo "=========================================================="

# 1. Run Go Engine Unit Tests
echo "➔ [1/3] Running VuiKey Linguistic Engine Tests (Go)..."
go test -v ./internal/engine

# 2. Run Swift Package Unit Tests if Swift toolchain is installed
if command -v swift >/dev/null 2>&1; then
    echo "➔ [2/3] Running VuiKeyEngine Swift Package Tests..."
    (cd "$ROOT_DIR/apple/VuiKeyKit" && swift test)
else
    echo "➔ [2/3] Swift toolchain not detected on this host; skipping 'swift test'."
fi

# 3. Build macOS Universal App Bundle (.app, .zip, .dmg) if running on macOS (Darwin)
if [[ "$(uname -s)" == "Darwin" ]]; then
    echo "➔ [3/3] Building macOS Universal Binary (arm64 + amd64) -> dist/VuiKey.app..."
    APP_DIR="$ROOT_DIR/dist/VuiKey.app"
    rm -rf "$APP_DIR" "$ROOT_DIR/dist/VuiKey-1.0.0-macos-universal.zip" "$ROOT_DIR/dist/VuiKey-1.0.0-macos-universal.dmg"
    mkdir -p "$APP_DIR/Contents/MacOS" "$APP_DIR/Contents/Resources" "$ROOT_DIR/build"

    CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "$ROOT_DIR/build/vuikey_darwin_arm64" ./cmd/vuikey_darwin
    CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "$ROOT_DIR/build/vuikey_darwin_amd64" ./cmd/vuikey_darwin

    lipo -create -output "$APP_DIR/Contents/MacOS/VuiKey" \
        "$ROOT_DIR/build/vuikey_darwin_arm64" \
        "$ROOT_DIR/build/vuikey_darwin_amd64"
    chmod +x "$APP_DIR/Contents/MacOS/VuiKey"

    cat > "$APP_DIR/Contents/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>VuiKey</string>
    <key>CFBundleIdentifier</key>
    <string>com.echipco.vuikey</string>
    <key>CFBundleName</key>
    <string>VuiKey</string>
    <key>CFBundleDisplayName</key>
    <string>VuiKey Vietnamese Input</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>LSMinimumSystemVersion</key>
    <string>12.0</string>
    <key>LSUIElement</key>
    <true/>
    <key>NSHighResolutionCapable</key>
    <true/>
</dict>
</plist>
EOF

    # Ad-hoc code sign the app bundle
    codesign --force --deep --sign - "$APP_DIR"

    # Create .zip archive (preserving macOS bundle metadata)
    (cd "$ROOT_DIR/dist" && ditto -c -k --sequesterRsrc --keepParent "VuiKey.app" "VuiKey-1.0.0-macos-universal.zip")

    # Create .dmg installer image
    DMG_STAGE="$ROOT_DIR/build/dmg_stage"
    rm -rf "$DMG_STAGE"
    mkdir -p "$DMG_STAGE"
    cp -R "$APP_DIR" "$DMG_STAGE/"
    ln -s /Applications "$DMG_STAGE/Applications"
    hdiutil create -volname "VuiKey" -srcfolder "$DMG_STAGE" -ov -format UDZO "$ROOT_DIR/dist/VuiKey-1.0.0-macos-universal.dmg"

    echo "✅ Built macOS Universal App:"
    echo "   • dist/VuiKey.app"
    echo "   • dist/VuiKey-1.0.0-macos-universal.zip"
    echo "   • dist/VuiKey-1.0.0-macos-universal.dmg"
else
    echo "➔ [3/3] Host is $(uname -s) (macOS Cocoa/CGEventTap requires Xcode SDK on macOS to link)."
    echo "        Run ./scripts/build-apple.sh on a Mac or GitHub Actions macos-latest runner to produce dist/VuiKey.app."
fi
