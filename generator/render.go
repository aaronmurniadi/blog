package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	goldmarklinkedimages "github.com/aaronmurniadi/goldmark-linked-images"
	figure "github.com/mangoumbrella/goldmark-figure"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func newTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"trim": strings.TrimSpace,
		"titlecase": func(s string) string {
			if s == "" {
				return s
			}
			return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
		},
		// absURL turns a site-relative path like "/avatar.jpg" into an absolute
		// URL using the origin from a canonical base URL (needed by OG scrapers).
		"absURL": func(canonical, path string) string {
			path = strings.TrimSpace(path)
			if strings.Contains(path, "://") || strings.HasPrefix(path, "//") {
				return path
			}
			base := canonical
			if base == "" || base == "/" {
				return path
			}
			segs := strings.Split(base, "/")
			const minSegs = 2 // protocol, host
			if len(segs) < minSegs+1 {
				return path
			}
			origin := strings.Join(segs[:minSegs+1], "/")
			rel := strings.TrimLeft(path, "/")
			if rel == "" {
				return origin
			}
			return origin + "/" + rel
		},
	}
}

func (s *Site) writeMarkdownPage(w io.Writer, abs, rel, urlPath string, rewriteMediaWebp bool) error {
	content, err := os.ReadFile(abs)
	if err != nil {
		return err
	}
	body, fm, err := parseFrontMatter(content)
	if err != nil {
		return fmt.Errorf("front matter: %w", err)
	}
	if fm.Title == "" {
		fm.Title = extractH1(body)
	}
	if fm.Title == "" {
		fm.Title = extractTitleFromPath(abs)
	}

	body, err = s.expandPostListTags(body)
	if err != nil {
		return err
	}

	htmlStr, err := mdToHTML(body)
	if err != nil {
		return err
	}
	if rewriteMediaWebp {
		htmlStr = rewriteMediaRasterToWebP(htmlStr)
	}

	headerTitle := "Home"
	hideNav := fm.NavBar != nil && !*fm.NavBar
	var nav []Link
	if !hideNav {
		var err error
		nav, err = buildNav(s)
		if err != nil {
			log.Println("buildNav:", err)
			nav = nil
		}
	}

	pageHTML := template.HTML(htmlStr)
	if fm.Date != "" {
		pageHTML = template.HTML("<p class=\"subtitle\">" + template.HTMLEscapeString(fm.Date) + "</p>" + string(pageHTML))
	}
	webRel := filepath.ToSlash(rel)
	if s.shouldApplyParaNum(webRel) {
		pageHTML = template.HTML(`<div class="para-num">` + string(pageHTML) + `</div>`)
	}
	if fm.Title != "" {
		headerTitle = fm.Title
	}
	page := Page{
		Title:       headerTitle,
		Path:        publicPath(urlPath),
		HTML:        pageHTML,
		Nav:         nav,
		HideNav:     hideNav,
		Canonical:   s.canonicalURL(urlPath),
		Description: fm.metaDescription(body),
		Image:       firstNonEmpty(fm.Image, "/avatar.jpg"),
		SiteName:    "Beago Cirius",
		SiteImage:   "/beago-cirius-logo-white.svg",
		Lang:        pageLang(fm, webRel),
		Type:        pageTypeFor(webRel),
	}
	if err := s.templates.ExecuteTemplate(w, "default.html", page); err != nil {
		return fmt.Errorf("execute default.html: %w", err)
	}
	return nil
}

// publicPath normalises a content URL like "/index.md" or "/typesettings/index.md"
// to its public form "/" or "/typesettings/".
func publicPath(urlPath string) string {
	path := strings.TrimSpace(urlPath)
	path = strings.TrimSuffix(path, ".md")
	nestedIndex := strings.HasSuffix(path, "/index")
	path = strings.TrimSuffix(path, "/index")
	if path == "" || path == "/" {
		return "/"
	}
	if nestedIndex {
		return strings.TrimRight(path, "/") + "/"
	}
	return path
}

// canonicalURL returns the absolute, canonical form of every page URL using the
// configured sitemap base (normalised to end with '/'). The input is a content
// URL like "/index.md" or "/articles/foo" and is normalised to its public form.
func (s *Site) canonicalURL(urlPath string) string {
	base := strings.TrimRight(s.cfg.SitemapBase, "/") + "/"
	path := strings.TrimSpace(urlPath)
	path = strings.TrimSuffix(path, ".md")
	nestedIndex := strings.HasSuffix(path, "/index")
	path = strings.TrimSuffix(path, "/index")
	if path == "" || path == "/" || path == "/index" {
		return base
	}
	rel := "/" + strings.Trim(path, "/")
	if nestedIndex {
		rel = rel + "/"
	}
	return base + strings.TrimLeft(rel, "/")
}

// pageTypeFor classifies a content-relative path into an Open Graph type.
func pageTypeFor(webRel string) string {
	for _, p := range []string{"posts/", "summaries/", "articles/"} {
		if strings.HasPrefix(webRel, p) {
			return "article"
		}
	}
	return "website"
}

// pageLang chooses an article language: explicit front-matter lang wins,
// otherwise Indonesian for paths under articles/, else English.
func pageLang(fm FrontMatter, webRel string) string {
	if fm.Lang != "" {
		return fm.Lang
	}
	if strings.HasPrefix(webRel, "articles/") {
		return "id"
	}
	return "en"
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func mdToHTML(md []byte) (string, error) {
	md = bytes.ReplaceAll(md, []byte("\r\n"), []byte("\n"))
	markdown := goldmark.New(
		goldmark.WithExtensions(
			figure.Figure,
			goldmarklinkedimages.LinkedImages,
			extension.Footnote,
			extension.Table,
			extension.Typographer,
		),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
	var buf bytes.Buffer
	if err := markdown.Convert(md, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
