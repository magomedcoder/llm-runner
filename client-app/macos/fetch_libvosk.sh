#!/bin/sh

set -e
ROOT="$(cd "$(dirname "$0")" && pwd)"
DEST="$ROOT/Frameworks"
MARKER="$DEST/libvosk.dylib"

if [ -f "$MARKER" ]; then
  exit 0
fi

if [ "${SKIP_VOSK_FETCH:-}" = "1" ]; then
  echo "fetch_libvosk.sh: SKIP_VOSK_FETCH=1, пропуск" >&2
  exit 0
fi

VER="${VOSK_OSX_PREBUILT_VERSION:-0.3.42}"
NAME="vosk-osx-${VER}"
URL="https://github.com/alphacep/vosk-api/releases/download/v${VER}/${NAME}.zip"

TMP="$(mktemp -d)"

trap 'rm -rf "$TMP"' EXIT

echo "Vosk macOS: загрузка ${URL}" >&2

curl -fsSL -o "$TMP/vosk.zip" "$URL"

unzip -q "$TMP/vosk.zip" -d "$TMP/ex"

SRC="$TMP/ex/${NAME}/libvosk.dylib"

if [ ! -f "$SRC" ]; then
  echo "fetch_libvosk.sh: в архиве нет ${NAME}/libvosk.dylib" >&2
  exit 1
fi

mkdir -p "$DEST"

cp -f "$SRC" "$MARKER"
