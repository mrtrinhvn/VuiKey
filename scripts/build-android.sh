#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR/android"

echo "=========================================================="
echo "🤖 VuiKey Android Builder (APK & Google Play Store AAB)"
echo "=========================================================="

mkdir -p "$ROOT_DIR/android/app/keystore" "$ROOT_DIR/dist"

# 1. Generate release signing keystore if not already present
KEYSTORE_FILE="$ROOT_DIR/android/app/keystore/vuikey-release.jks"
if [[ ! -f "$KEYSTORE_FILE" ]]; then
    echo "➔ [1/4] Generating release keystore for VuiKey Android..."
    keytool -genkeypair -v \
        -keystore "$KEYSTORE_FILE" \
        -alias "${VUIKEY_KEY_ALIAS:-vuikey}" \
        -keyalg RSA -keysize 2048 -validity 10000 \
        -storepass "${VUIKEY_KS_PASS:-vuikey2026}" \
        -keypass "${VUIKEY_KEY_PASS:-vuikey2026}" \
        -dname "CN=VuiKey, OU=Software, O=eChip Co Ltd, L=Ho Chi Minh, S=Ho Chi Minh, C=VN"
fi

# 2. Run Kotlin Linguistic Engine JUnit Tests
echo "➔ [2/4] Running VuiKeyEngine Kotlin Unit Tests..."
gradle testDebugUnitTest --no-daemon

# 3. Build Signed Release APK (for direct installation) & Signed Release AAB (for Google Play Console)
echo "➔ [3/4] Building Release APK & Google Play App Bundle (.aab)..."
gradle assembleRelease bundleRelease --no-daemon

cp -f "$ROOT_DIR/android/app/build/outputs/apk/release/app-release.apk" "$ROOT_DIR/dist/VuiKey-1.0.0-android.apk"
cp -f "$ROOT_DIR/android/app/build/outputs/bundle/release/app-release.aab" "$ROOT_DIR/dist/VuiKey-1.0.0-android-playstore.aab"

echo "✅ Built Android Release Artifacts:"
ls -lh "$ROOT_DIR/dist/VuiKey-1.0.0-android.apk" "$ROOT_DIR/dist/VuiKey-1.0.0-android-playstore.aab"
