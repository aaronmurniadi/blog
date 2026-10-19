#!/bin/bash

# cd to the directory containing this script
cd "$(dirname "$0")"

rebuild_all=false
for arg in "$@"; do
    case "$arg" in
        --all) rebuild_all=true ;;
        *)
            echo "Usage: $0 [--all]"
            echo "  --all  Recompile every .typ (default: missing/stale pdf or webp)"
            exit 1
            ;;
    esac
done

# Check for typst command
if ! command -v typst &> /dev/null; then
    echo "Error: typst is not installed or not in PATH."
    exit 1
fi

# Check for pdftocairo (Poppler)
if ! command -v pdftocairo &> /dev/null; then
    echo "Error: pdftocairo is not installed or not in PATH (install poppler)."
    exit 1
fi

# Check for cwebp (libwebp)
if ! command -v cwebp &> /dev/null; then
    echo "Error: cwebp is not installed or not in PATH (install webp)."
    exit 1
fi

list_typ_sources() {
    find . -name '*.typ' ! -name 'maid_of_orleans_style.typ' ! -path './.git/*' -print |
        sed 's|^\./||' | sort -u
}

# Rebuild when outputs are absent, the source is newer, or a co-located .typ
# import bundle changed (e.g. maid_of_orleans_style.typ). Does not use git status,
# so pull/clone with gitignored pdf/webp still compiles.
typ_needs_build() {
    local typfile=$1
    local base="${typfile%.typ}"
    local pdffile="${base}.pdf"
    local webpfile="${base}.webp"
    local dir
    local script="${BASH_SOURCE[0]}"

    [ ! -f "$pdffile" ] || [ ! -f "$webpfile" ] && return 0
    [ "$typfile" -nt "$pdffile" ] || [ "$typfile" -nt "$webpfile" ] && return 0
    [ "$script" -nt "$pdffile" ] || [ "$script" -nt "$webpfile" ] && return 0

    dir=$(dirname "$typfile")
    [ "$dir" = "." ] || dir="./$dir"

    # Rebuild if a local .typ imported/included by this document changed
    # (e.g. maid_of_orleans.typ imports maid_of_orleans_style.typ). Package
    # imports like "@preview/..." do not end in .typ and are ignored.
    while IFS= read -r dep; do
        [ -n "$dep" ] || continue
        case "$dep" in
            /*) ;;
            *) dep="$dir/$dep" ;;
        esac
        [ -f "$dep" ] || continue
        [ "$dep" -nt "$pdffile" ] || [ "$dep" -nt "$webpfile" ] && return 0
    done < <(grep -oE '#(import|include)[[:space:]]+"[^"]+\.typ"' "$typfile" \
        | sed -E 's/^#[a-z]+[[:space:]]+"//; s/"$//')

    return 1
}

typ_files=()
while IFS= read -r line; do
    [ -z "$line" ] && continue
    if [ "$rebuild_all" = true ] || typ_needs_build "$line"; then
        typ_files+=("$line")
    fi
done < <(list_typ_sources)

if [ "$rebuild_all" = true ]; then
    echo "Regenerating all first-page images (${#typ_files[@]} .typ files)"
elif [ ${#typ_files[@]} -eq 0 ]; then
    echo "No .typ files to process."
    exit 0
else
    echo "Stale or missing outputs for: ${typ_files[*]}"
fi

shopt -s nullglob
errors=0
for typfile in "${typ_files[@]}"; do
    base="${typfile%.typ}"
    pdffile="${base}.pdf"
    webpfile="${base}.webp"

    # Compile .typ to .pdf. --root .. (i.e. content/media/) widens the sandbox
    # so sources may reference sibling media (e.g. mini.typ's ../images/...).
    echo "Compiling $typfile to $pdffile"
    if ! typst compile --root .. "$typfile" "$pdffile"; then
        echo "Error compiling $typfile"
        errors=1
        continue
    fi

    # Convert first page of pdf to webp (pdftocairo has no webp; pipe via png + cwebp)
    echo "Converting first page of $pdffile to $webpfile"
    tmpdir=$(mktemp -d)
    tmpbase="$tmpdir/page"
    if ! pdftocairo -png -r 300 -singlefile "$pdffile" "$tmpbase"; then
        echo "Error rasterizing $pdffile to png"
        rm -rf "$tmpdir"
        errors=1
        continue
    fi
    if ! cwebp -q 90 -quiet "${tmpbase}.png" -o "$webpfile"; then
        echo "Error converting $pdffile to webp"
        rm -rf "$tmpdir"
        errors=1
        continue
    fi
    rm -rf "$tmpdir"
done

exit "$errors"
