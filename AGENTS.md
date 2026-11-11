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
  `DATE`, `SITEMAP`, `BODY_CLASS`). Must stay valid `sh` (`sh -n` is gated).
  `BODY_CLASS` becomes `<body class="...">`: `front` for `/`, `gallery` for
  `/photography`, `/typesettings`, `/fujipes`. Both lift the narrow reading
  measure (see CSS layout notes).
- `content/<path>/body.html` — inner `<main>` fragment (no `<main>` tags).
- `content/media/`, `content/*.{png,ico,svg,txt,webmanifest}` — copied verbatim.
- `templates/header.html`, `nav.html`, `footer.html` — shell around every page.
- Article images: full files in `content/media/images/`, 800px thumbs in
  `content/media/images/thumbnail/` (see `_thumbnail.sh`, magick `800x800>`),
  embedded as `figure > a[href=full] > img[src=thumb + alt] + figcaption > p`
  (cf. `content/photography/body.html`). Captions are styled globally (smaller
  type, tight under the image); `class="aside"` opts a figure into the small
  floated text-flow treatment (alternating sides by default; `aside-left` /
  `aside-right` force a side, headings clear). Only vendor
  PD/CC0, or CC BY(-SA) with author + license credited in the caption.
- `templates/font-switcher.html` — font picker partial, injected by `build.sh`
  as the first child of every page's `<main>` (top-right row).
- `templates/410.html`, `style.css`, `prism.css/js` — copied verbatim.
- `templates/fonts/*.woff2` — copied to `public/fonts/` (license `.txt` files
  in that dir are NOT copied; they are attribution carriers, keep them).
- `smarten.py` — stdin→stdout HTML filter, applied to `body.html` only.
- `public/` — build output. NEVER edit, never format, never commit.
- `mockups/` — local design mocks, git-ignored, never built or deployed. A
  tabbed harness (`mockups/index.html`) that loads four alternative layouts
  plus a copy of the built home page in an iframe, with viewport and theme
  controls. Run `mockups/serve.sh [port]` (default 8004); it symlinks
  `public/fonts`, `public/media`, and `public/style.css` at its web root.

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

- Dropdown lives in `templates/font-switcher.html` (injected into the
  `.utilities` row by `build.sh`, after `</main>`; + mirrored markup in
  `templates/410.html`).
- Scope: families from https://r2src.github.io/top10fonts/ (9 of 10 — Boisik is
  Metafont-sources-only upstream and cannot be vendored, so it is omitted),
  plus Baskervald X, Bembo, Palatino, Crimson. NO system fonts in the dropdown.
- Site default reading face is Crimson (`--reading`, `"Crimson", Georgia,
  serif`) — "Default" option = no `data-font` attribute. Headings always use
  `--display` (`"Baskervald X", Georgia, serif`), so the dropdown changes body
  and prose only.
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

## Theme switcher (dark/light)

- Dropdown lives in `templates/theme-switcher.html` (injected into the
  `.utilities` row by `build.sh`, right after the font switcher; + mirrored
  markup in `templates/410.html`).
- Writes/reads `localStorage` key `site-theme`; stored value is applied to the
  `data-theme` attribute on `<html>` plus `style.colorScheme` (so native
  controls, scrollbars, and form fields match). Both persisted in the `<head>`
  inline script and the trailing control script in `build.sh`/`410.html`.
- Every colour is a token on `:root` (`--paper`, `--ink`, `--muted`,
  `--rule`, `--accent`). The dark theme is one token swap at the bottom of
  `templates/style.css` (`html[data-theme="dark"]`), so no rule below the tokens
  may hardcode a colour. Keep the palette consistent (`#16130f` paper,
  `#ece5d9` ink, `#e08a5b` accent).
- The selected value is never allowlisted (only `light`/`dark` are offered), so
  the theme script needs no `ok={...}` stale-value guard (unlike fonts).

## CSS layout notes

- The sheet is `min(100% - 2.5rem, 60rem)` (`1.5rem` below `40rem`). `body` is
  a flex column: `.masthead`, `main`, `.utilities`, `footer`.
- Prose stops at `--measure` (`36rem`), centered inside the sheet: `main` is
  capped, not every child. `body.front` (`/`) and `body.gallery` lift the cap;
  gallery pages re-cap prose blocks to `--measure` and center them, so pictures
  grow and text does not.
- Masthead: site name + byline left, sections right, one rule under both. Below
  `52rem` it stacks into two rows. The active link gets `aria-current="page"`
  from a small script at the end of `<body>` (prefix match, so a post under
  `/posts/` marks Posts).
- Front page only: `.standfirst` is a two-column grid (title left, intro copy
  right with a rule between), then a `.section-head` label and a `.spread`
  two-column index. Both collapse to one column below `52rem`.
- Date-index listings (`<li><time>date</time> <a>title</a></li>`, e.g.
  `/posts`, `/articles`, `/summaries`) use the `main li:has(> time)` flex rules:
  the date is a fixed no-wrap tabular-nums column so wrapped titles stay in
  their own column instead of sliding under the date. Below `28rem` the date
  stacks above the title.
- Tables use the `main table` rules (`border-collapse`, thin `--rule` borders,
  top-aligned padded cells, `0.92rem`).
- Figure images get a `1px solid var(--rule)` frame, `max-height: 90dvh` with
  `width: auto` / `max-width: 100%` so aspect ratio is preserved. `main figure`
  is `display: table` with `figcaption` as `table-caption` so caption width
  matches the image, not the column. `.aside*` figures float from `44rem` up;
  headings `clear: both`.
- `content/fujipes/body.html` carries its own `<style>` block that reads
  `--border`, `--text`, `--header-bg`, `--heading-font-family`; `:root` keeps
  those four as aliases of the new tokens. Do not remove them without editing
  that block.

## Article rewrite workflow (repeatable)

First applied to `content/articles/2019-06-11-pertanyaan-pertanyaan-metafisika/`.
Reuse for other long texts.

- Prose: natural flowing Indonesian, one idea developed per paragraph (never
  choppy one-sentence-per-line style). NO em dashes anywhere, literal `—` or
  `&mdash;` — use commas, periods, colons, semicolons, or parentheses instead.
  The count must be 0 afterwards.
- Structure: one article = one chapter. `h1` = chapter title (when retitling,
  update `TITLE_FULL`, `PLAIN_TITLE`, `DESCRIPTION` in `page.meta` too, but
  keep `CANONICAL` stable so URLs never break). Sections are unnumbered
  topical `h2`s with `id`s plus `h3` subsections, fronted by a `Daftar Isi`
  `h2` + `ol` of anchor links (Daftar Pustaka included). Exam-style Q&A
  headings are rewritten as descriptive headings; question headings with no
  body are deleted (zero content loss).
- Quotations are frozen: never edit or translate `<blockquote>` blocks or
  original-language quotations/terms. Verify afterwards with a
  whitespace-normalized `<blockquote>` old-vs-new comparison (must report
  identical).
- Facts: verify risky claims (page numbers, ordinals, Greek terms, doctrinal
  attributions) via web search, then cite ORIGINAL BOOKS, never websites:
  append an `h2#daftar-pustaka` "Daftar Pustaka" — a `ul` of book entries
  (author, title, editor/translator, publisher, year), each annotated with the
  cited sections and which part of the article they support. Fix what
  verification disproves (unattested terms are removed, not footnoted); keep
  hedges the original already had (e.g. "selalu dianggap").
- Tables: `thead`+`tbody`, `scope="col"`/`"row"`, no presentational attributes
  (`cellspacing`, `cellpadding`, `width`, `valign`), no spacer rows, straight
  `"` in cells → `&ldquo;`/`&rdquo;`. Cells crammed with numbered sub-items
  are split into clean two/three-column tables.
- After any rewrite: `./formatter.sh` (+ `--write` if needed) and `./build.sh`
  must exit 0; re-run the em-dash count and the blockquote-identity check.

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
- `mockups/` holds throwaway UI exploration. It is not part of the build and
  `formatter.sh` does not check it. Delete the directory when a direction is
  chosen; port the winning rules into `templates/style.css` by hand. The
  `Current` tab is a copy of `public/index.html`, so run `./build.sh` before
  comparing.

## Maintaining this file

- Keep `AGENTS.md` current in the same change: when you find important repo
  facts a future agent needs (moved files, renamed templates, new or changed
  constraints, anything above gone stale), update the relevant section here.
  Docs and code must stay in sync — a correct change with a stale doc is
  incomplete.
