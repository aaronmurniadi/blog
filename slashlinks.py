#!/usr/bin/env python3
"""Append a trailing slash to internal page links at build time.

Reads HTML from stdin, writes it to stdout. Every page is a directory served
at its trailing-slash form: Caddy answers /posts/ with 200 and 308s /posts to
it. Canonical URLs, the sitemap, and internal links must therefore all use the
slash form, otherwise Google files the page as "Page with redirect" or as an
alternate of a URL that redirects.

Only href="/..." values whose last path segment has no dot are rewritten, so
stylesheets, scripts, images, fonts, manifests, and downloads pass through
untouched. Query strings and fragments are kept after the slash, and the root
link "/" is already in slash form.
"""
import re
import sys

LINK = re.compile(r'href="([^"]*)"')


def add_slash(match):
    value = match.group(1)
    if not value.startswith("/") or value.startswith("//"):
        return match.group(0)
    cuts = [i for i in (value.find("#"), value.find("?")) if i != -1]
    cut = min(cuts) if cuts else len(value)
    path, rest = value[:cut], value[cut:]
    if path.endswith("/"):
        return match.group(0)
    last = path.rsplit("/", 1)[-1]
    if "." in last:
        return match.group(0)
    return 'href="%s/%s"' % (path, rest)


def main():
    sys.stdout.write(LINK.sub(add_slash, sys.stdin.read()))


if __name__ == "__main__":
    main()
