#!/bin/sh
# Format all sources for consistency. Tools used: sh + prettier.
#   ./formatter.sh          check mode: lists unformatted files, exit 1 if any
#   ./formatter.sh --write  format in place, then rebuild public/ via build.sh
#
# Scope: HTML sources only (content/**/body.html, templates/*.html).
# public/ is build output (see .prettierignore) and is never formatted.
# Shell files (build.sh, formatter.sh, page.meta) get a syntax check only.
set -e
cd "$(dirname "$0")"

HTML_FILES=$(find content templates -name '*.html' | sort)

check_shell() {
    sh -n build.sh
    sh -n formatter.sh
    python3 -c "import ast; ast.parse(open('smarten.py').read())"
    for meta in $(find content -name page.meta | sort); do
        sh -n "$meta" || exit 1
        ( . "./$meta" ) || exit 1
    done
}

case "${1:-}" in
    --write)
        check_shell
        prettier --write $HTML_FILES
        ./build.sh
        ;;
    *)
        check_shell
        prettier --check $HTML_FILES
        ;;
esac
