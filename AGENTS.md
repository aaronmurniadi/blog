# AGENTS.md

Hand-rolled Go static site generator (no Hugo/Jekyll). Markdown in `content/` is rendered to HTML in `_site/` by `go run .`.

## Built-in workflow

- **Run from repo root**: the tool resolves all paths (`content`, `templates`, `media`, `fonts`, static assets) relative to the current working directory. Use `go run .` from the repo root.
- **Output**: writes `_site/` by default. `_site/` and the compiled `blog` binary are gitignored — never commit build output.
- **Flags**: `-out`, `-content`, `-sitemap-base` (or `BLOG_SITEMAP_BASE`), `-write-sitemap` (write `sitemap.xml` to cwd and exit), `-para-num-paths` (default `articles/,summaries/`).
- Two Go packages: the `main` entrypoint (`main.go`, root) imports `blog/generator`. All generator logic lives under `generator/`: the `generator.Site` core type + path-traversal safety (`safePathUnderContent`) in `site.go`, markdown→HTML + front matter in `content.go`/`render.go`, static-site flattening in `static.go`, the reserved-name skip policy in `policy.go`, sitemap in `sitemap.go`, WebP conversion in `mediawebp.go`, and the `{% postList %}` tag in `postlist.go`. Only exported symbols (`NewSite`, `SiteConfig`, `GenerateStaticSite`, `WriteSitemapFile`) are called from `main.go`.

## Content conventions

- Markdown files under `content/`. `index.md` at the site root serves `/`. Every other `.md` becomes `/<path>`, and `articles/foo.md` becomes `/articles/foo` (built as `_site/articles/foo/index.html`).
- **`_`-prefixed** files and dirs (and any dir whose segment starts with `_`) are skipped everywhere — dir listings, nav, sitemap, and generation. These are private/unpublished; files with `.md` under `_drafts/` won't be published.
- Directories named `media`, `assets`, or `scripts` under content are skipped both in nav and generation.
- A directory with no `index.md` gets an auto-generated directory-listing page (via `templates/dirindex.html`) listing child dirs and `.md` files, sorted newest-first by `date`. A sibling `X.md` ("shadow") suppresses the listing for `/X` rather than overwrite it.

### Front matter (YAML between `---` delimiters)

Supported fields: `title`, `date` (used for post-listing sort and sitemap `<lastmod>`; `time.Parse` formats `2006-01-02` or e.g. `January 2, 2006`), and `nav_bar: false` to hide the `<nav>`. `title`/`date` fall back to the first H1 and file mtime respectively. Older posts carry vestigial Jekyll keys (`layout`, `nav_order`) that the parser ignores — don't rely on them.

### Special rendering

- **Nav** is built automatically by walking top-level `content/` dirs that contain `.md` files (plus a Home link when `content/index.md` exists). You don't maintain nav by hand.
- **Paragraph numbers**: URLs under the `-para-num-paths` prefixes (default `articles/`, `summaries/`) get wrapped in `<div class="para-num">`, and the `default.html` template numbers each `p` after the first.
- **`{% postList collections.<name> %}`** in markdown body expands to a `<ul class="post-list">` of `content/<name>/`'s `.md` files, sorted by date desc.
- **WebP media**: during static generation, raster images under `media/` (`jpeg/png/gif`) are converted to WebP (quality 90, concurrency 4) into `_site/media/`; the rendered HTML has any `/media/...` raster URL rewritten to `.webp`. Output is idempotent (skips existing targets). The URL rewrite deliberately requires a leading quote so absolute off-site URLs aren't touched.
- Markdown rendering uses goldmark with figure, linked-images, footnote, table, and typographer extensions; `WithUnsafe()` is enabled, so raw HTML in markdown passes through.

## Security-sensitive areas (be careful when editing)

- `generator/site.go` defends against path traversal (`safePathUnderContent`). Don't weaken this check; it's the boundary between the renderer and the filesystem (used to render dir-listing pages).
- Reserved dir names (`media`, `assets`, `scripts`, `templates`) and the `_`-prefix rule are centralized in `generator/policy.go` (`skipContentDir`, `skipContentDirName`, `shouldSkipContentRel`); `content.go`, `static.go`, `sitemap.go` all use it. Keep policy changes there, not at call sites.
- Production is served by Caddy + Anubis in a two-port layout (see `Caddyfile`); the live site content root is registered at `/home/aaron/blog/_site`. Sitemap base defaults to `https://aaron.beago-cirius.ts.net/`.

## Tests and verification

Unit tests live next to the sources (`*_test.go`); run with `go test ./...`. After changing the generator, verify by running `go run .` from the repo root and inspecting the regenerated `_site/` for affected pages.
