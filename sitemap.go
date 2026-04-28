package main

import (
	"bytes"
	"encoding/xml"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// SitemapURL is one loc entry for sitemap.xml.
type SitemapURL struct {
	Loc      string `xml:"loc"`
	LastMod  string `xml:"lastmod,omitempty"`
	Priority string `xml:"priority,omitempty"`
}

// SitemapRoot is the root element for sitemap.xml marshaling.
type SitemapRoot struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLns   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

func collectSitemapURLs(contentRoot, base string) ([]SitemapURL, error) {
	contentRoot, err := filepath.EvalSymlinks(filepath.Clean(contentRoot))
	if err != nil {
		return nil, err
	}
	base = strings.TrimRight(base, "/") + "/"

	var urls []SitemapURL
	err = filepath.WalkDir(contentRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if name == "node_modules" || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		rel, err := filepath.Rel(contentRoot, path)
		if err != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)
		if strings.HasPrefix(rel, "_") || strings.Contains(rel, "/_") {
			return nil
		}

		rel = strings.TrimSuffix(rel, ".md")
		loc := base + rel
		if strings.HasSuffix(loc, "/index") {
			loc = strings.TrimSuffix(loc, "index")
		} else if rel == "index" {
			loc = base
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		_, fm, err := parseFrontMatter(content)
		if err != nil {
			return nil
		}
		u := SitemapURL{Loc: loc}
		if fm.Date != "" {
			u.LastMod = fm.Date
		} else {
			u.LastMod = info.ModTime().Format("2006-01-02")
		}
		urls = append(urls, u)
		return nil
	})
	return urls, err
}

func buildSitemap(contentRoot, base string) (SitemapRoot, error) {
	urls, err := collectSitemapURLs(contentRoot, base)
	if err != nil {
		return SitemapRoot{}, err
	}
	return SitemapRoot{
		XMLns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}, nil
}

func writeSitemapFile(_, contentRoot, base, outPath string) error {
	sitemap, err := buildSitemap(contentRoot, base)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(sitemap); err != nil {
		return err
	}
	_ = buf.WriteByte('\n')
	return os.WriteFile(outPath, buf.Bytes(), 0644)
}
