#!/bin/sh
set -e
cd "$(dirname "$0")"
PORT="${1:-${PORT:-8000}}"
exec python3 -m http.server "$PORT" --directory public
