#!/usr/bin/env python3
"""Render a page's Markdown body to an HTML fragment at build time.

Reads Markdown from stdin, writes an HTML fragment for <main> to stdout.

Only the subset the vendored READMEs use is supported: ATX headings (with
slug ids), paragraphs, unordered and ordered lists, tables, fenced code
blocks, blockquotes, horizontal rules, and raw HTML blocks (a line that
starts with "<" is copied through with its following lines up to a blank
line). Inline: code spans, images, links, strong and emphasis.

The output then goes through smarten.py like every other body, so straight
single quotes in prose still become the curly pair.

No third-party dependency: python3 stdlib only.
"""
import html
import re
import sys

# Prism ships bash, shell, python, typst, markup, css, and javascript in
# templates/prism.js. Anything else loses its class so Prism leaves it alone.
PRISM_LANGS = {
    "bash",
    "shell",
    "sh",
    "python",
    "py",
    "typst",
    "typc",
    "html",
    "markup",
    "css",
    "javascript",
    "js",
}
LANG_ALIASES = {"sh": "bash", "shell": "bash", "zsh": "bash", "py": "python"}

HEADING = re.compile(r"^(#{1,6})\s+(.*?)\s*#*\s*$")
FENCE = re.compile(r"^```(.*)$")
HR = re.compile(r"^\s*([-*_])(\s*\1){2,}\s*$")
LIST_ITEM = re.compile(r"^(\s*)([-*+]|\d+[.)])\s+(.*)$")
TABLE_ROW = re.compile(r"^\s*\|")
TABLE_SEP = re.compile(r"^\s*\|?[\s:|-]+\|[\s:|-]*$")
QUOTE = re.compile(r"^>\s?(.*)$")
ENTITY = re.compile(r"&(?:[A-Za-z][A-Za-z0-9]*|#[0-9]+|#[xX][0-9A-Fa-f]+);")
PLACEHOLDER = re.compile(r"\x00(\d+)\x00")
CODE_SPAN = re.compile(r"`([^`]+)`")
IMAGE = re.compile(r"!\[([^\]]*)\]\(([^)\s]+)\)")
LINK = re.compile(r"\[([^\]]+)\]\(([^)\s]+)\)")
STRONG = re.compile(r"\*\*(?!\s)(.+?)(?<!\s)\*\*")
EM = re.compile(r"(?<!\*)\*(?!\s)(.+?)(?<!\s)\*(?!\*)")
# A paragraph made only of images (optional link wrappers) is a badge row.
BADGE_ONLY = re.compile(
    r"^(\s*(?:<img\b[^>]*>|<a\b[^>]*>\s*<img\b[^>]*>\s*</a>))+\s*$"
)


def escape_text(s):
    """Escape bare &, <, > but keep entities and placeholders intact."""
    out = []
    for i, ch in enumerate(s):
        if ch == "&":
            if ENTITY.match(s, i):
                out.append(ch)
            else:
                out.append("&amp;")
            continue
        if ch == "<":
            out.append("&lt;")
            continue
        if ch == ">":
            out.append("&gt;")
            continue
        out.append(ch)
    return "".join(out)


def escape_attr(s):
    return html.escape(s, quote=True)


def inline(text, store=None):
    """Markdown inline markup -> HTML. `store` holds stashed fragments."""
    if store is None:
        store = []

    def stash(fragment):
        store.append(fragment)
        return "\x00%d\x00" % (len(store) - 1)

    # 1. Code spans keep their content literal.
    def code_span(m):
        return stash("<code>" + html.escape(m.group(1), quote=False) + "</code>")

    text = CODE_SPAN.sub(code_span, text)

    # 2. Images before links, so [![alt](img)](url) nests correctly.
    def image(m):
        return stash(
            '<img src="%s" alt="%s">' % (escape_attr(m.group(2)), escape_attr(m.group(1)))
        )

    text = IMAGE.sub(image, text)

    def link(m):
        return stash(
            '<a href="%s">%s</a>' % (escape_attr(m.group(2)), inline(m.group(1), store))
        )

    text = LINK.sub(link, text)

    # 3. Escape whatever is left, then apply emphasis.
    text = escape_text(text)
    text = STRONG.sub(r"<strong>\1</strong>", text)
    text = EM.sub(r"<em>\1</em>", text)

    # 4. Put the stashed fragments back (nested ones need several passes).
    for _ in range(len(store) + 1):
        if "\x00" not in text:
            break
        text = PLACEHOLDER.sub(lambda m: store[int(m.group(1))], text)
    return text


def slug(text):
    plain = re.sub(r"<[^>]+>", "", text)
    plain = re.sub(r"\x00\d+\x00", "", plain)
    plain = re.sub(r"`", "", plain)
    out = re.sub(r"[^a-z0-9]+", "-", plain.lower()).strip("-")
    if not out:
        return "section"
    return out


class Renderer:
    def __init__(self):
        self.ids = {}

    def unique_id(self, text):
        base = slug(text)
        seen = self.ids.get(base, 0)
        self.ids[base] = seen + 1
        if seen:
            return "%s-%d" % (base, seen + 1)
        return base

    def render(self, src):
        lines = src.splitlines()
        out = []
        i = 0
        n = len(lines)
        while i < n:
            line = lines[i]
            if not line.strip():
                i += 1
                continue

            fence = FENCE.match(line)
            if fence:
                i = self.fence(lines, i, fence.group(1).strip(), out)
                continue

            heading = HEADING.match(line)
            if heading:
                level = len(heading.group(1))
                text = inline(heading.group(2))
                anchor = self.unique_id(heading.group(2))
                out.append(
                    '<h%d id="%s">%s</h%d>' % (level, anchor, text, level)
                )
                i += 1
                continue

            if HR.match(line):
                out.append("<hr>")
                i += 1
                continue

            if TABLE_ROW.match(line):
                i = self.table(lines, i, out)
                continue

            if QUOTE.match(line):
                i = self.quote(lines, i, out)
                continue

            if LIST_ITEM.match(line):
                i = self.list_block(lines, i, out)
                continue

            if line.lstrip().startswith("<"):
                i = self.raw_html(lines, i, out)
                continue

            i = self.paragraph(lines, i, out)
        return "\n".join(out) + "\n"

    def fence(self, lines, i, info, out):
        lang = info.split()[0] if info else ""
        lang = LANG_ALIASES.get(lang, lang)
        body = []
        i += 1
        while i < len(lines) and not lines[i].startswith("```"):
            body.append(lines[i])
            i += 1
        i += 1  # closing fence
        code = html.escape("\n".join(body), quote=False)
        if lang and lang in PRISM_LANGS:
            out.append('<pre><code class="language-%s">%s</code></pre>' % (lang, code))
        else:
            out.append("<pre><code>%s</code></pre>" % code)
        return i

    def raw_html(self, lines, i, out):
        block = []
        while i < len(lines) and lines[i].strip():
            block.append(lines[i])
            i += 1
        out.append("\n".join(block))
        return i

    def paragraph(self, lines, i, out):
        buf = []
        while i < len(lines) and lines[i].strip():
            line = lines[i]
            if HEADING.match(line) or FENCE.match(line) or HR.match(line):
                break
            if TABLE_ROW.match(line) or LIST_ITEM.match(line) or QUOTE.match(line):
                break
            if line.lstrip().startswith("<"):
                break
            buf.append(line.strip())
            i += 1
        if buf:
            html_text = inline(" ".join(buf))
            cls = ' class="badges"' if BADGE_ONLY.match(html_text) else ""
            out.append("<p%s>%s</p>" % (cls, html_text))
        return i

    def list_block(self, lines, i, out):
        ordered = LIST_ITEM.match(lines[i]).group(2)[0].isdigit()
        items = []
        while i < len(lines):
            line = lines[i]
            if not line.strip():
                break
            item = LIST_ITEM.match(line)
            if not item:
                # An indented line with no marker continues the last item.
                if items and line.startswith(("  ", "\t")):
                    items[-1].append(line.strip())
                    i += 1
                    continue
                break
            if item.group(2)[0].isdigit() != ordered:
                break
            items.append([item.group(3).strip()])
            i += 1
        tag = "ol" if ordered else "ul"
        out.append("<%s>" % tag)
        for parts in items:
            out.append("<li>%s</li>" % inline(" ".join(parts)))
        out.append("</%s>" % tag)
        return i

    def table(self, lines, i, out):
        rows = []
        while i < len(lines) and lines[i].strip():
            rows.append(lines[i])
            i += 1
        head = split_row(rows[0])
        body = rows[2:] if len(rows) > 1 and TABLE_SEP.match(rows[1]) else rows[1:]
        out.append("<table>")
        out.append("<thead>")
        out.append("<tr>")
        for cell in head:
            out.append('<th scope="col">%s</th>' % inline(cell))
        out.append("</tr>")
        out.append("</thead>")
        if body:
            out.append("<tbody>")
            for row in body:
                cells = split_row(row)
                out.append("<tr>")
                for cell in cells:
                    out.append("<td>%s</td>" % inline(cell))
                out.append("</tr>")
            out.append("</tbody>")
        out.append("</table>")
        return i

    def quote(self, lines, i, out):
        buf = []
        while i < len(lines):
            m = QUOTE.match(lines[i])
            if not m:
                break
            buf.append(m.group(1).strip())
            i += 1
        out.append("<blockquote><p>%s</p></blockquote>" % inline(" ".join(buf)))
        return i


def split_row(line):
    s = line.strip()
    if s.startswith("|"):
        s = s[1:]
    if s.endswith("|"):
        s = s[:-1]
    return [cell.strip() for cell in s.split("|")]


def main():
    src = sys.stdin.read()
    sys.stdout.write(Renderer().render(src))


if __name__ == "__main__":
    main()
