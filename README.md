# blog

Small Go static site: markdown under `content/`, HTML under `_site/`.

## Run

From repo root (paths are relative to cwd):

```bash
go run .
```

Writes `_site/` by default. Useful flags:

| Flag | Default | Notes |
|------|---------|--------|
| `-content` | `content` | Markdown tree |
| `-out` | `_site` | Output directory |
| `-sitemap-base` | `BLOG_SITEMAP_BASE` or built-in URL | Trailing `/` added if missing |
| `-write-sitemap` | off | Writes `sitemap.xml` in cwd and exits |
| `-para-num-paths` | `articles/,summaries/` | Comma-separated path prefixes for paragraph numbering |

Example:

```bash
BLOG_SITEMAP_BASE=https://example.com/ go run .
```

## Deploy layout (Caddy + Anubis)

Repo includes `Caddyfile` for a two-port setup:

```caddy
# --- Public Entrance ---
:8002 {
    # Send everything to Anubis (assuming Anubis runs on 3000)
    reverse_proxy localhost:3000 {
        # Pass the real IP so Anubis can track the client
        header_up X-Real-Ip {remote_host}
    }
}

# --- Internal "Clean" Backend ---
:8003 {
    # Bind to localhost so this port isn't exposed to the internet
    bind 127.0.0.1
    
    root * /home/aaron/blog/_site
    
    # Speed optimizations from before
    encode zstd gzip
    file_server
    
    # Optional: Aggressive Caching
    @static path *.ico *.css *.js *.gif *.jpg *.jpeg *.png *.svg *.woff *.woff2
    header @static Cache-Control "public, max-age=31536000"
}
```

1. **`:8002` (public)** — reverse proxy to Anubis on `localhost:3000`, passes `X-Real-Ip`.
2. **`:8003` (localhost only)** — serves `_site` with `file_server`, zstd/gzip, long cache on static assets.

Regenerate `_site` after editing content, then reload Caddy if needed.

```bash
caddy start
```

Anubis should sit between public Caddy and the inner static server: set its upstream to the internal listener (e.g. `http://127.0.0.1:8003`). Exact variable names depend on your Anubis version; keep them in a dedicated env file.

## Anubis env

Production env for this stack lives at:

`/etc/anubis/blog.env`

```bash
# The port Anubis will listen on (for Caddy port 8002 to point to)
BIND=127.0.0.1:3000

# The port where Caddy is "secretly" serving your files
TARGET=http://127.0.0.1:8003

# The URL users see in their browser
PUBLIC_URL=http://localhost:8002

# Challenge difficulty (4 is standard)
DIFFICULTY=4

# Automatically block common AI scrapers
SERVE_ROBOTS_TXT=true
```

Start anubis with:

```bash
sudo systemctl enable --now anubis@blog.service
systemctl status anubis@blog.service
```

Point your Anubis systemd unit (or container) `EnvironmentFile=` at that path so secrets and `TARGET` (or equivalent upstream URL) stay out of the repo.
