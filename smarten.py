#!/usr/bin/env python3
"""Smarten straight single quotes in HTML prose at build time.

Reads an HTML fragment from stdin, writes it to stdout with U+0027 '
in prose text converted to U+2018 '/U+2019 ' (like the hand-set
&ldquo;/&rdquo; doubles). Byte-preserving otherwise: tags, attributes,
entities and comments pass through untouched, and text inside
<pre>/<code>/<script>/<style> is skipped so code samples keep straight
quotes. Apostrophes (didn't, Python's, O'Hagan) become '.
"""
import re
import sys

OPENERS = set('([{ "\u201c\u2018\xab\u2013\u2014-/')
SKIP_TAGS = {"pre", "code", "script", "style"}
LEFT = "\u2018"
RIGHT = "\u2019"

TOKEN = re.compile(
    r"<!--[^-]*-(?:[^-][^-]*-)*?>|"
    r"<\?.*?\?>|"
    r"<![^>]*>|"
    r"<(?:[^\x22\x27>]|\x22[^\x22]*\x22|\x27[^\x27]*\x27)*>|"
    r"[^<]+",
    re.DOTALL,
)
TAGNAME = re.compile(r"^</?\s*([A-Za-z][^\s/>]*)")
VOID = {"br", "hr", "img", "meta", "link", "input"}


def smarten_node(s, prev):
    out = []
    for i, ch in enumerate(s):
        if ch != "'":
            out.append(ch)
            continue
        p = out[-1] if out else prev
        n = s[i + 1] if i + 1 < len(s) else None
        if p is not None and p.isalnum():
            rep = RIGHT  # apostrophe mid-word / closer / possessive
        elif n is not None and n.isalnum():
            if n.isdigit() and (p is None or p.isspace() or p in OPENERS):
                rep = RIGHT  # decade abbreviation: '80s
            else:
                rep = LEFT  # opener before a word
        elif p is None or p.isspace() or p in OPENERS:
            # Quote right after an opening bracket with no word following
            # it closes a quoted bracket ('[', ']'). Otherwise it opens.
            rep = RIGHT if p in "([{" else LEFT
        else:
            rep = RIGHT  # closer
        out.append(rep)
    res = "".join(out)
    return res, (res[-1] if res else prev)


def smarten(src):
    stack = []
    prev = None
    parts = []
    for m in TOKEN.finditer(src):
        tok = m.group(0)
        if tok.startswith("<"):
            parts.append(tok)
            tm = TAGNAME.match(tok)
            if (
                tm
                and not tok.startswith("<!--")
                and not tok.startswith("<?")
                and not tok.startswith("<!")
            ):
                name = tm.group(1).lower()
                if tok.startswith("</"):
                    if name in stack[::-1]:
                        while stack and stack.pop() != name:
                            pass
                elif not tok.endswith("/>") and name not in VOID:
                    stack.append(name)
        elif any(t in SKIP_TAGS for t in stack) or "'" not in tok:
            parts.append(tok)
            if tok:
                prev = tok[-1]
        else:
            new, prev = smarten_node(tok, prev)
            parts.append(new)
    return "".join(parts)


def main():
    sys.stdout.write(smarten(sys.stdin.read()))


if __name__ == "__main__":
    main()
