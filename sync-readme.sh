#!/bin/sh
# Vendor the OfxHost README.md into content/ofxhost/body.md.
#
# The blog page /ofxhost/ is a copy of the project README, so the two cannot
# drift by hand. Run this after the README changes, then commit body.md (and
# the logo, if it changed). build.sh converts body.md to HTML with md.py, so
# the cron build stays offline and deterministic and never touches GitHub.
#
# Rewrites applied to the downloaded Markdown:
#   - the logo  docs/logo.png  becomes the site copy /media/images/ofxhost-logo.png
#   - repo-relative links such as [LICENSE](LICENSE) become absolute GitHub URLs
# An HTML comment at the top marks the file as generated.
set -e
cd "$(dirname "$0")"

REPO=aaronmurniadi/ofxhost
BRANCH=main
RAW="https://raw.githubusercontent.com/$REPO/$BRANCH"
BLOB="https://github.com/$REPO/blob/$BRANCH"
DEST=content/ofxhost/body.md
LOGO=content/media/images/ofxhost-logo.png

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT INT TERM

mkdir -p content/ofxhost content/media/images

curl -fsSL "$RAW/README.md" -o "$tmp/README.md"
curl -fsSL "$RAW/docs/logo.png" -o "$tmp/logo.png"

{
    printf '<!-- Vendored copy of %s/README.md.\n' "$BLOB"
    printf '     Regenerate with ./sync-readme.sh. Do not edit by hand. -->\n'
    sed -e "s|src=\"docs/logo.png\"|src=\"/media/images/ofxhost-logo.png\"|" \
        -e "s|](LICENSE)|]($BLOB/LICENSE)|" "$tmp/README.md"
} > "$tmp/body.md"

if cmp -s "$tmp/body.md" "$DEST" && cmp -s "$tmp/logo.png" "$LOGO"; then
    echo "sync-readme: content is unchanged"
    exit 0
fi

mv "$tmp/body.md" "$DEST"
mv "$tmp/logo.png" "$LOGO"
echo "sync-readme: updated $DEST and $LOGO"
./build.sh
