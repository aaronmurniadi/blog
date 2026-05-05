package main

import (
	"fmt"
	"html"
	"path/filepath"
	"regexp"
	"strings"
)

var postListTagRE = regexp.MustCompile(`\{%\s*postList\s+collections\.([a-zA-Z0-9_-]+)\s*%\}`)

func (s *Server) expandPostListTags(md []byte) ([]byte, error) {
	var err error
	out := postListTagRE.ReplaceAllFunc(md, func(match []byte) []byte {
		if err != nil {
			return match
		}
		sub := postListTagRE.FindSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		name := string(sub[1])
		abs := filepath.Join(s.contentDir(), filepath.FromSlash(name))
		links, e := s.linksForMarkdownFilesInDir(abs)
		if e != nil {
			err = fmt.Errorf("postList %q: %w", name, e)
			return match
		}
		return []byte(renderPostListHTML(links))
	})
	if err != nil {
		return md, err
	}
	return out, nil
}

func renderPostListHTML(links []Link) string {
	if len(links) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("<ul class=\"post-list\">\n")
	for _, link := range links {
		b.WriteString("<li>")
		if link.Date != "" {
			b.WriteString(`<span class="post-list-date">`)
			b.WriteString(html.EscapeString(link.Date))
			b.WriteString(`</span> `)
		}
		b.WriteString(`<a href="`)
		b.WriteString(html.EscapeString(link.Path))
		b.WriteString(`">`)
		b.WriteString(html.EscapeString(link.Title))
		b.WriteString("</a></li>\n")
	}
	b.WriteString("</ul>\n")
	return b.String()
}
