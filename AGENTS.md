# AGENTS.md — Beago Cirius tiny static site

Personal blog. Static HTML built by shell + Python, no framework, no bundler.

## Commands

```bash
./build.sh            # regenerate public/ (requires: sh + python3)
./formatter.sh        # gate: prettier --check on HTML sources, sh -n on shell/page.meta, AST parse of smarten.py
./formatter.sh --write  # prettier --write, then rebuild
./serve-local.sh [port] # preview public/ (default 8000)
```

Always run `./formatter.sh` and `./build.sh` after touching content, templates,
`build.sh`, or `smarten.py`. Both must exit 0.

## Layout (source of truth)

- `content/<path>/page.meta` — shell vars sourced by `build.sh` (`TITLE_FULL`,
  `PLAIN_TITLE`, `DESCRIPTION`, `CANONICAL`, `OG_TYPE`, `LANG`, `LOCALE`,
  `DATE`, `SITEMAP`). Must stay valid `sh` (`sh -n` is gated).
- `content/<path>/body.html` — inner `<main>` fragment (no `<main>` tags).
- `content/media/`, `content/*.{png,ico,svg,txt,webmanifest}` — copied verbatim.
- `templates/header.html`, `nav.html`, `footer.html` — shell around every page.
- `templates/font-switcher.html` — font picker partial, injected by `build.sh`
  as the first child of every page's `<main>` (top-right row).
- `templates/410.html`, `style.css`, `prism.css/js` — copied verbatim.
- `templates/fonts/*.woff2` — copied to `public/fonts/` (license `.txt` files
  in that dir are NOT copied; they are attribution carriers, keep them).
- `smarten.py` — stdin→stdout HTML filter, applied to `body.html` only.
- `public/` — build output. NEVER edit, never format, never commit.

## Typography pipeline (quotes)

- In `body.html` prose, write plain straight `'` for single quotes/apostrophes.
  The build smartens them to `‘`/`’` via `smarten.py`.
- `smarten.py` transforms text nodes only: tags, attributes, entities, comments
  pass through byte-identical; text inside `<pre>/<code>/<script>/<style>` is
  skipped so code samples keep straight quotes.
- Double quotes are HAND-SET as `&ldquo;`/`&rdquo;` in sources (no automation).
- Do NOT hand-place `‘`/`’` in `body.html` prose (exception: authentic glyphs
  inside `<code>`, e.g. PDF extracts — leave those alone).
- `page.meta` values are copied verbatim (no smartening) AND cannot contain
  straight `'` (single-quoted shell syntax). Use entities or curly chars there.

## Font system

- Dropdown lives in `templates/font-switcher.html` (top-right row of `<main>`,
  injected by `build.sh`; + mirrored markup in `templates/410.html`).
- Scope: families from https://r2src.github.io/top10fonts/ (9 of 10 — Boisik is
  Metafont-sources-only upstream and cannot be vendored, so it is omitted),
  plus Baskervald X, Bembo, Palatino, Crimson. NO system fonts in the dropdown.
- Site default (base `body`) is the system stack
  `"Times New Roman", Times, Georgia, ui-serif, serif` — "Default" option =
  no `data-font` attribute.
- Adding a font requires ALL of these (they must stay in sync):
  1. `templates/fonts/<name>-400.woff2`, `-400-italic`, `-700`, `-700-italic`
     (omit a style only if upstream never made it — see gaps below) + a license
     file (`OFL-*.txt`, `LICENCE-*`, `COPYRIGHT-*`, …).
  2. `@font-face` blocks + `html[data-font="<value>"] body` rule in
     `templates/style.css` (document source/license in the switcher comment).
  3. `<option value="<value>">` in BOTH `font-switcher.html` and `410.html`.
  4. The value in the `ok={...}` allowlist in BOTH inline scripts
     (`build.sh` printf line + `410.html`) — stale stored values fall back
     to default.
- Known gaps (documented in the CSS comment, do not "fix" by synthesis):
  URW Antiqua ships regular-only; Bera Serif never had italics
  (browser-synthesized); Utopia entry uses Erewhon outlines (Utopia-derived).
- Conversions: OTF/TTF→woff2 via `fontTools` (`TTFont.flavor = "woff2"`,
  needs `brotli`); Type1 PFB→OTF via AFDKO `makeotf` (`tx -dump` to inspect;
  verify glyph/cmap counts — subset TeX fonts like pxfonts `rpx*` are traps).
  Sources: `https://mirrors.ctan.org/fonts/<pkg>.zip`,
  metadata: `https://ctan.org/json/2.0/pkg/<pkg>`.

## CSS layout notes

- Narrow: single centered column `min(100% - 2rem, 42rem)`.
- Wide (`@media (min-width: 60rem)`): grid `55rem` total, `13rem` sidebar +
  `1fr` main, gap `clamp(1.5rem, 4vw, 3rem)`. Keep main ≈39rem when retuning.
- Flush tops: `text-box-trim: trim-start` on `header h1` + `main > :first-child`,
  with matching `0.1em` top gap on both. The gap is load-bearing: without it,
  cap overshoot hits the sidebar's `overflow-y` clip edge. Do not use a
  negative margin to compensate (it puts ink back on the clip edge).

## Verification checklist (after font/layout/build changes)

- `./formatter.sh` and `./build.sh` exit 0.
- If inline JS changed: `node --check` the extracted `<script>` from both
  `public/index.html` and `public/410.html`.
- If fonts changed: every dropdown `value` (minus `default`) has a `data-font`
  rule; every `url("/fonts/...")` exists in `public/fonts/`; no references to
  retired families remain in CSS/HTML.

## Notes

- `README.md` dependency line predates `smarten.py` (build needs `python3`).
- Prettier formats HTML sources only; `public/` is excluded by scope, not by
  `.prettierignore` (there is none).

## Maintaining this file

- Keep `AGENTS.md` current in the same change: when you find important repo
  facts a future agent needs (moved files, renamed templates, new or changed
  constraints, anything above gone stale), update the relevant section here.
  Docs and code must stay in sync — a correct change with a stale doc is
  incomplete.
