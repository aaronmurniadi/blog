# blog

Beago Cirius tiny static site.

Build only depends on: `sh`, `cat`, `printf`, `mkdir`, `cp`, `find`, `sort`, `date`.

## Layout

| Path                                                                        | What it is                                                                                                                   |
| --------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| `content/<path>/page.meta`                                                  | Per-page shell vars: `TITLE_FULL`, `PLAIN_TITLE`, `DESCRIPTION`, `CANONICAL`, `OG_TYPE`, `LANG`, `LOCALE`, `DATE`, `SITEMAP` |
| `content/<path>/body.html`                                                  | Inner `<main>` fragment (no `<main>` tags)                                                                                   |
| `content/media/`                                                            | Site media (`images/`, `typst/`, …), copied verbatim to `public/media/`                                                      |
| `content/*.{png,ico,svg,txt,webmanifest}`                                   | Site-root static files (favicons, `robots.txt`, `site.webmanifest`, logo), copied to `public/`                               |
| `templates/header.html`, `nav.html`, `footer.html`                          | Shell wrapped around every page                                                                                              |
| `templates/410.html`, `style.css`, `prism.css`, `prism.js`, `fonts/*.woff2` | Copied verbatim to `public/`                                                                                                 |
| `public/`                                                                   | Build output, git-ignored (`.gitignore`)                                                                                     |
| `build.sh`                                                                  | The site builder                                                                                                             |
| `formatter.sh`                                                              | `prettier --check` on HTML sources + `sh -n` on shell/`page.meta`; `--write` formats and rebuilds                            |
| `serve-local.sh`                                                            | Local preview: `python3 -m http.server --directory public`                                                                   |
| `Caddyfile`                                                                 | Local/prod static server for `public/` with 410 handling                                                                     |

## Build

```bash
./build.sh
```

1. Pages: `content/<stem>/` → `public/<stem>/index.html` (`content/index/` → `public/index.html`).
2. Copies: `templates/410.html`, `style.css`, `prism.css/js`, `fonts/*.woff2`; site-root files + `media/` from `content/`.
3. Sitemap: `public/sitemap.xml` from `page.meta` (`CANONICAL` + `DATE`; `DATE=''` means build date, `SITEMAP=0` excludes).

## Add a page

```bash
mkdir content/<path>
# write content/<path>/page.meta + content/<path>/body.html
./build.sh
```

## Format / check

```bash
./formatter.sh          # prettier --check on content/**/body.html + templates/*.html, sh -n on shell files
./formatter.sh --write  # prettier --write, then rebuild via build.sh
```

`public/` is never formatted.

## Preview

```bash
./serve-local.sh        # PORT env or arg, default 8000
./serve-local.sh 8080
```

Or via Caddy (serves `public/` on `127.0.0.1:8002`, missing paths → `410.html` with 410 status, long cache on static assets):

```bash
caddy run --config Caddyfile
```

## Deploy layout (Caddy + Anubis)

`Caddyfile` in repo is the internal static backend:

```caddy
:8002 {
	bind 127.0.0.1
	root * public
	encode zstd gzip
	# …long cache on *.ico *.gif *.jpg *.jpeg *.png *.svg *.woff *.woff2,
	# missing paths raise 410 and serve /410.html with X-Robots-Tag: noindex
}
```

Regenerate `public/` after editing content, then reload Caddy if needed.
