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
	page := Page{Title: headerTitle, Path: urlPath, HTML: pageHTML, Nav: nav, HideNav: hideNav}
	if err := s.templates.ExecuteTemplate(w, "default.html", page); err != nil {
		return fmt.Errorf("execute default.html: %w", err)
	}
	return nil
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
