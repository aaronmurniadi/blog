#!/bin/sh
# Beago Cirius tiny site builder.
# Tools used: sh, cat, printf, mkdir, cp, find, sort, date, python3 (only for
# smarten.py's quote filter). No other deps.
#
# Source of truth:
#   content/<path>/page.meta   shell vars: TITLE_FULL, PLAIN_TITLE, DESCRIPTION,
#                               CANONICAL, OG_TYPE, LANG, LOCALE, DATE, SITEMAP
#   content/<path>/body.html   inner <main> fragment (no <main> tags).
#                               Write plain ' for single quotes; the build
#                               smartens them to '/'. Double quotes are hand-set
#                               &ldquo;/&rdquo; in the source.
#   content/media/             site media, copied verbatim to public/media/
#   content/*.{png,ico,svg,txt,webmanifest}  site-root static files, copied to public/
#   templates/header.html, nav.html, footer.html, 410.html
#   smarten.py                 stdin->stdout HTML filter, applied to body.html only
#
# Output: public/<path>/index.html + public/sitemap.xml + public/410.html
#   + public/media/ + public/{favicons,robots.txt,...} (all copied from content/)
#
# Add a page: mkdir content/<path>, write page.meta + body.html, run ./build.sh
#   DATE='' means "use build date" (for section indexes).
#   SITEMAP=0 excludes the page from sitemap.xml (kept for /posts, /recaps,
#   /summaries indexes to match the existing sitemap).
set -e
cd "$(dirname "$0")"

command -v python3 >/dev/null || {
    echo "build.sh needs python3 (smarten.py quote filter)" >&2
    exit 1
}

T=templates
C=content
OUT=public

render_page() {
    meta=$1
    body=$2
    out=$3
    TITLE_FULL=
    PLAIN_TITLE=
    DESCRIPTION=
    CANONICAL=
    OG_TYPE=website
    LANG=en
    LOCALE=en
    # shellcheck disable=SC1090
    . "./$meta"
    {
        printf '<!DOCTYPE html>\n<html lang="%s">\n\n<head>\n' "$LANG"
        printf '    <meta charset="UTF-8">\n'
        printf '    <meta name="viewport" content="width=device-width, initial-scale=1.0">\n'
        printf '    <meta name="google-site-verification" content="WueWgxazZfzyRmReKOxkb36Ub42Dz-ZLndfjb3HkzVs" />\n'
        printf '    <title>%s</title>\n' "$TITLE_FULL"
        printf '    <meta name="description" content="%s">\n' "$DESCRIPTION"
        printf '    <link rel="canonical" href="%s">\n' "$CANONICAL"
        printf '    <meta name="author" content="Aaron P. Murniadi">\n'
        printf '    <meta name="robots" content="index, follow">\n'
        printf '    <meta name="generator" content="Beago Cirius static site generator">\n'
        printf '\n'
        printf '    <meta property="og:site_name" content="Beago Cirius">\n'
        printf '    <meta property="og:type" content="%s">\n' "$OG_TYPE"
        printf '    <meta property="og:title" content="%s">\n' "$PLAIN_TITLE"
        printf '    <meta property="og:description" content="%s">\n' "$DESCRIPTION"
        printf '    <meta property="og:url" content="%s">\n' "$CANONICAL"
        printf '    <meta property="og:image" content="https://aaron.beago-cirius.ts.net/avatar.jpg">\n'
        printf '    <meta property="og:locale" content="%s">\n' "$LOCALE"
        printf '\n'
        printf '    <meta name="twitter:card" content="summary_large_image">\n'
        printf '    <meta name="twitter:site" content="@aaronpmurniadi">\n'
        printf '    <meta name="twitter:title" content="%s">\n' "$PLAIN_TITLE"
        printf '    <meta name="twitter:description" content="%s">\n' "$DESCRIPTION"
        printf '    <meta name="twitter:image" content="https://aaron.beago-cirius.ts.net/avatar.jpg">\n'
        printf '\n'
        printf '    <link rel="icon" href="/favicon.ico" sizes="any">\n'
        printf '    <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">\n'
        printf '    <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">\n'
        printf '    <link rel="icon" type="image/png" sizes="192x192" href="/android-chrome-192x192.png">\n'
        printf '    <link rel="apple-touch-icon" href="/apple-touch-icon.png">\n'
        printf '    <link rel="manifest" href="/site.webmanifest">\n'
        printf '    <link rel="stylesheet" href="/style.css">\n'
        printf '    <link rel="stylesheet" href="/prism.css">\n'
        printf '    <script>try{var f=localStorage.getItem("site-font");if(f)document.documentElement.setAttribute("data-font",f)}catch(e){}</script>\n'
        printf '\n    \n\n    \n\n    \n</head>\n\n<body>\n    \n'
        printf '        <div class="sidebar">\n'
        cat "$T/header.html"
        cat "$T/nav.html"
        printf '        </div>\n'
        printf '        <main>\n'
        cat "$T/font-switcher.html"
        python3 ./smarten.py < "$body"
        printf '        </main>\n'
        cat "$T/footer.html"
        printf '    <script src="/prism.js" defer></script>\n'
        printf '    <script>(function(){var s=document.getElementById("font-switcher");if(!s)return;var ok={default:1,baskervaldx:1,computermodern:1,kpfonts:1,gfsdidot:1,utopia:1,venturis:1,libertine:1,gyrebonum:1,gyrepagella:1,gyreschola:1,gyretermes:1,antiqua:1,bera:1,bembo:1,palatino:1,crimson:1};var f="default";try{f=document.documentElement.getAttribute("data-font")||"default"}catch(e){}if(!ok[f]){f="default";try{document.documentElement.removeAttribute("data-font");localStorage.removeItem("site-font")}catch(e){}}s.value=f;s.addEventListener("change",function(){var v=s.value;try{if(v==="default"){document.documentElement.removeAttribute("data-font");localStorage.removeItem("site-font")}else{document.documentElement.setAttribute("data-font",v);localStorage.setItem("site-font",v)}}catch(e){}})})();</script>\n'
        printf '\n    \n    \n</body>\n\n</html>'
    } > "$out"
}

# 1. Pages: content/<stem>/ -> public/<stem>/index.html (stem "index" -> public/)
for meta in $(find "$C" -name page.meta | sort); do
    dir=${meta%/*}
    stem=${dir#"$C/"}
    body="$dir/body.html"
    if [ "$stem" = "index" ]; then
        out="$OUT/index.html"
    else
        out="$OUT/$stem/index.html"
    fi
    mkdir -p "${out%/*}"
    render_page "$meta" "$body" "$out"
done

# 2. Error page: verbatim copy, never templated
cp "$T/410.html" "$OUT/410.html"
cp "$T/style.css" "$OUT/style.css"
cp "$T/prism.css" "$OUT/prism.css"
cp "$T/prism.js" "$OUT/prism.js"
mkdir -p "$OUT/fonts"
cp "$T/fonts/"*.woff2 "$OUT/fonts/"

# 2b. Static assets from content/: site-root files + media/
# Any regular file directly under content/ (favicons, robots.txt,
# site.webmanifest, svg) is copied to public/. content/media/ -> public/media/.
mkdir -p "$OUT"
for f in "$C"/*; do
    [ -f "$f" ] && cp "$f" "$OUT/"
done
if [ -d "$C/media" ]; then
    rm -rf "$OUT/media"
    cp -R "$C/media" "$OUT/media"
fi

# 3. Sitemap: loc + lastmod come straight from page.meta, so existing links
#    (incl. trailing-slash quirks and the 2025-12-01 lastmod) are preserved.
{
    printf '<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n'
    for meta in $(find "$C" -name page.meta | sort); do
        CANONICAL=
        DATE=
        SITEMAP=1
        # shellcheck disable=SC1090
        . "./$meta"
        [ "$SITEMAP" = "0" ] && continue
        [ -z "$CANONICAL" ] && continue
        [ -z "$DATE" ] && DATE=$(date +%F)
        printf '  <url>\n    <loc>%s</loc>\n    <lastmod>%s</lastmod>\n  </url>\n' "$CANONICAL" "$DATE"
    done
    printf '</urlset>\n'
} > "$OUT/sitemap.xml"
