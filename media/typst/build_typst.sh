#!/bin/bash

# cd to the directory containing this script
cd "$(dirname "$0")"

rebuild_all=false
for arg in "$@"; do
    case "$arg" in
        --all) rebuild_all=true ;;
        *)
            echo "Usage: $0 [--all]"
            echo "  --all  Regenerate first-page webp for every .typ (default: only git-changed)"
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

if [ "$rebuild_all" = true ]; then
    # All standalone .typ under this dir (skip style-only import bundle)
    typ_files=()
    while IFS= read -r line; do
        [ -n "$line" ] && typ_files+=("$line")
    done < <(find . -name '*.typ' ! -name 'maid_of_orleans_style.typ' ! -path './.git/*' -print | sed 's|^\./||' | sort -u)
    echo "Regenerating all first-page images (${#typ_files[@]} .typ files)"
else
    # Changed .typ only (paths relative to this dir, same as git media/typst/)
    typ_files=()
    while IFS= read -r line; do
        [ -n "$line" ] && typ_files+=("$line")
    done < <(git status --porcelain | grep '^...media/typst/.*\.typ$' | sed 's/^...//' | sed 's|media/typst/||')
fi

if [ ${#typ_files[@]} -eq 0 ]; then
    echo "No .typ files to process."
    exit 0
fi

if [ "$rebuild_all" != true ]; then
    echo "Found changed .typ files: ${typ_files[*]}"
fi

shopt -s nullglob
for typfile in "${typ_files[@]}"; do
    base="${typfile%.typ}"
    pdffile="${base}.pdf"
    webpfile="${base}.webp"

    # Compile .typ to .pdf
    echo "Compiling $typfile to $pdffile"
    typst compile "$typfile" "$pdffile"
    if [ $? -ne 0 ]; then
        echo "Error compiling $typfile"
        continue
    fi

    # Convert first page of pdf to webp (pdftocairo has no webp; pipe via png + cwebp)
    echo "Converting first page of $pdffile to $webpfile"
    tmpdir=$(mktemp -d)
    tmpbase="$tmpdir/page"
    pdftocairo -png -r 300 -singlefile "$pdffile" "$tmpbase"
    if [ $? -ne 0 ]; then
        echo "Error rasterizing $pdffile to png"
        rm -rf "$tmpdir"
        continue
    fi
    cwebp -q 90 -quiet "${tmpbase}.png" -o "$webpfile"
    status=$?
    rm -rf "$tmpdir"
    if [ $status -ne 0 ]; then
        echo "Error converting $pdffile to webp"
        continue
    fi
done
