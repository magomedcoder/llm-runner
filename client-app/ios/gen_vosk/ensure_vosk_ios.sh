#!/bin/sh

set -e
ROOT="$(cd "$(dirname "$0")" && pwd)"
XCFW="$ROOT/Vosk.xcframework"

if [ -d "$XCFW" ]; then
  echo "gen_vosk: Vosk.xcframework уже есть."
  exit 0
fi

if [ "${SKIP_VOSK_FETCH:-}" = "1" ] && [ -z "${VOSK_IOS_XCFRAMEWORK_ZIP_URL:-}" ]; then
  echo "gen_vosk: SKIP_VOSK_FETCH=1 и нет URL — положите Vosk.xcframework в ios/gen_vosk/ вручную." >&2
  exit 1
fi

if [ -n "${VOSK_IOS_XCFRAMEWORK_ZIP_URL:-}" ]; then
  echo "gen_vosk: загрузка $VOSK_IOS_XCFRAMEWORK_ZIP_URL" >&2
  
  TMP="$(mktemp -d)"
  
  trap 'rm -rf "$TMP"' EXIT
  
  curl -fsSL -o "$TMP/vosk_ios.zip" "$VOSK_IOS_XCFRAMEWORK_ZIP_URL"
  
  unzip -q "$TMP/vosk_ios.zip" -d "$TMP/ex"
  
  FOUND="$(find "$TMP/ex" -maxdepth 4 -type d -name 'Vosk.xcframework' | head -1)"
  
  if [ -z "$FOUND" ]; then
    echo "gen_vosk: в архиве нет Vosk.xcframework" >&2
    exit 1
  fi

  rm -rf "$XCFW"

  cp -R "$FOUND" "$XCFW"

  exit 0
fi

if [ -n "${VOSK_IOS_DEVICE_LIB:-}" ] && [ -n "${VOSK_IOS_SIM_LIB:-}" ]; then
  echo "gen_vosk: xcodebuild -create-xcframework из двух libvosk.a" >&2
  
  TMP="$(mktemp -d)"
  
  trap 'rm -rf "$TMP"' EXIT
  
  mkdir -p "$TMP/device" "$TMP/sim"
  
  cp "$VOSK_IOS_DEVICE_LIB" "$TMP/device/libvosk.a"
  
  cp "$VOSK_IOS_SIM_LIB" "$TMP/sim/libvosk.a"
  
  xcrun xcodebuild -create-xcframework \
    -library "$TMP/device/libvosk.a" \
    -library "$TMP/sim/libvosk.a" \
    -output "$XCFW"

  exit 0
fi
