package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	goldmarklinkedimages "github.com/aaronmurniadi/goldmark-linked-images"
	figure "github.com/mangoumbrella/goldmark-figure"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type Page struct {
	Title string
	Path  string
	HTML  template.HTML
	Nav   []Link
}

type Link struct {
	Path     string
	Title    string
	Date     string
	SortDate string
}

var (
	mdDir         string
	tmplDir       string
	port          int
	writeSitemapF bool
	funcMap = template.FuncMap{
		"trim":      strings.TrimSpace,
		"titlecase": func(s string) string { return strings.ToUpper(s[:1]) + strings.ToLower(s[1:]) },
	}
	templates *template.Template
)

func main() {
	flag.IntVar(&port, "port", 8080, "port to serve on")
	flag.IntVar(&port, "p", 8080, "port to serve on")
	flag.BoolVar(&writeSitemapF, "write-sitemap", false, "write sitemap.xml to current directory and exit")
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	mdDir = wd
	tmplDir = filepath.Join(wd, "templates")

	if writeSitemapF {
		out := filepath.Join(wd, "sitemap.xml")
		if err := writeSitemapFile(wd, out); err != nil {
			log.Fatal(err)
		}
		log.Printf("wrote %s", out)
		return
	}

	templates = template.New("").Funcs(funcMap)
	templates, err = templates.ParseGlob(filepath.Join(tmplDir, "*.html"))
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/sitemap.xml", handleSitemapFile)
	http.Handle("/style.css", http.FileServer(http.Dir(wd)))
	http.Handle("/favicon.ico", http.FileServer(http.Dir(wd)))
	http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir(filepath.Join(wd, "media")))))
	log.Printf("Serving at http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	path := r.URL.Path

	if path == "/" {
		path = "/index"
	}
	if strings.HasSuffix(path, "/") {
		indexPath := path + "index.md"
		absPath := filepath.Join(mdDir, filepath.Clean(indexPath))
		if _, err := os.Stat(absPath); err == nil {
			path = indexPath
		} else {
			serveDirIndex(filepath.Join(mdDir, filepath.Clean(path)), w, r)
			return
		}
	} else if !strings.HasSuffix(path, ".md") && !strings.HasPrefix(path, "/media/") {
		path = path + ".md"
	}

	cleanPath := filepath.Clean(path)
	absPath := filepath.Join(mdDir, cleanPath)

	if !strings.HasPrefix(absPath, mdDir) {
		http.Error(w, "Forbidden", 403)
		return
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Not Found", 404)
			return
		}
		http.Error(w, err.Error(), 500)
		return
	}

	if info.IsDir() {
		indexPath := filepath.Join(absPath, "index.md")
		if _, err := os.Stat(indexPath); err == nil {
			absPath = indexPath
		} else {
			serveDirIndex(absPath, w, r)
			return
		}
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	mdContent, pageTitle, date := stripFrontMatter(content)

	html, err := mdToHTML(mdContent)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	title := pageTitle
	if title == "" {
		title = extractTitle(absPath)
	}
	headerTitle := "Home"
	nav := buildNav()
	pageHTML := template.HTML(html)
	if date != "" {
		pageHTML = template.HTML("<p class=\"subtitle\">" + date + "</p>" + string(html))
	}
	if shouldParaNum(filepath.ToSlash(cleanPath)) {
		pageHTML = template.HTML(`<div class="para-num">` + string(pageHTML) + `</div>`)
	}
	page := Page{Title: headerTitle, Path: path, HTML: pageHTML, Nav: nav}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, "default.html", page); err != nil {
		log.Println(err)
	}
}

func serveDirIndex(dirPath string, w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var links []Link
	relDir, _ := filepath.Rel(mdDir, dirPath)
	for _, e := range entries {
		if e.IsDir() {
			ep := filepath.Join(relDir, e.Name())
			ep = strings.TrimPrefix(ep, string(filepath.Separator))
			ep = "/" + ep + "/"
			title := strings.ToUpper(e.Name()[:1]) + e.Name()[1:]
			links = append(links, Link{Path: ep, Title: title})
		} else if strings.HasSuffix(e.Name(), ".md") {
			mdPath := filepath.Join(dirPath, e.Name())
			content, err := os.ReadFile(mdPath)
			if err != nil {
				continue
			}
			_, title, date := stripFrontMatter(content)
			ep := filepath.Join(relDir, e.Name())
			ep = strings.TrimPrefix(ep, string(filepath.Separator))
			ep = "/" + strings.TrimSuffix(ep, ".md")
			sortDate := ""
			if date != "" {
				if t, err := time.Parse("2006-01-02", date); err == nil {
					sortDate = t.Format("2006-01-02")
				} else if t, err := time.Parse("January 2, 2006", date); err == nil {
					sortDate = t.Format("2006-01-02")
				} else if t, err := time.Parse("Jan 2, 2006", date); err == nil {
					sortDate = t.Format("2006-01-02")
				}
			}
			links = append(links, Link{Path: ep, Title: title, Date: date, SortDate: sortDate})
		}
	}

	sort.Slice(links, func(i, j int) bool {
		if links[i].SortDate == links[j].SortDate {
			return links[i].Title < links[j].Title
		}
		if links[i].SortDate == "" {
			return false
		}
		if links[j].SortDate == "" {
			return true
		}
		return links[i].SortDate > links[j].SortDate
	})

	title := strings.ToUpper(filepath.Base(dirPath)[:1]) + filepath.Base(dirPath)[1:]
	nav := buildNav()
	htmlStr := "<h1>" + title + "</h1><ul>" + linksToHTML(links) + "</ul>"
	page := Page{Title: title, Path: r.URL.Path, HTML: template.HTML(htmlStr), Nav: nav}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, "default.html", page); err != nil {
		log.Println(err)
	}
}

func linksToHTML(links []Link) string {
	var b strings.Builder
	for _, l := range links {
		b.WriteString(`<li>`)
		if l.Date != "" {
			b.WriteString(`<span style="color:#666;font-size:0.9em">`)
			b.WriteString(l.Date)
			b.WriteString(`</span> `)
		}
		b.WriteString(`<a href="`)
		b.WriteString(l.Path)
		b.WriteString(`">`)
		b.WriteString(l.Title)
		b.WriteString("</a>")
		b.WriteString(`</li>`)
	}
	return b.String()
}

func stripFrontMatter(md []byte) ([]byte, string, string) {
	trimmed := bytes.TrimSpace(md)
	if !bytes.HasPrefix(trimmed, []byte("---")) {
		return md, extractH1(md), ""
	}

	lines := strings.Split(string(trimmed), "\n")
	if len(lines) < 3 {
		return md, extractH1(md), ""
	}

	var metaLines []string
	contentStart := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			contentStart = i + 1
			break
		}
		metaLines = append(metaLines, lines[i])
	}

	if contentStart == -1 {
		return md, extractH1(md), ""
	}

	var date string
	var title string
	for _, line := range metaLines {
		if strings.HasPrefix(line, "title:") {
			title = strings.TrimSpace(strings.TrimPrefix(line, "title:"))
			title = strings.Trim(title, `"`)
			title = strings.Trim(title, "'")
		}
		if strings.HasPrefix(line, "date:") {
			date = strings.TrimSpace(strings.TrimPrefix(line, "date:"))
			date = strings.Trim(date, `"`)
			date = strings.Trim(date, "'")
		}
	}

	content := []byte(strings.Join(lines[contentStart:], "\n"))
	return bytes.TrimSpace(content), title, date
}

func extractH1(md []byte) string {
	for _, line := range strings.Split(string(md), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return ""
}

func mdToHTML(md []byte) (string, error) {
	md = bytes.ReplaceAll(md, []byte("\r\n"), []byte("\n"))
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.Footnote,
			extension.Table,
			extension.Typographer,
			figure.Figure,
			goldmarklinkedimages.LinkedImages,
		),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
	var buf bytes.Buffer
	if err := markdown.Convert(md, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func extractTitle(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, ".md")
}

// shouldParaNum: paragraph margin numbers only for long-form under articles/ and summaries/.
func shouldParaNum(cleanPath string) bool {
	return strings.Contains(cleanPath, "articles/") || strings.Contains(cleanPath, "summaries/")
}

func buildNav() []Link {
	type dirInfo struct {
		path  string
		title string
	}
	var dirList []dirInfo

	filepath.Walk(mdDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		if path == mdDir {
			return nil
		}
		rel, _ := filepath.Rel(mdDir, path)
		name := filepath.Base(path)
		if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") || name == "content" || name == "media" || name == "assets" || name == "scripts" || name == "templates" {
			return filepath.SkipDir
		}
		hasMD := false
		entries, _ := os.ReadDir(path)
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".md") {
				hasMD = true
				break
			}
		}
		if !hasMD {
			return filepath.SkipDir
		}
		dirList = append(dirList, dirInfo{path: rel, title: name})
		return nil
	})

	sort.Slice(dirList, func(i, j int) bool {
		return dirList[i].path < dirList[j].path
	})

	var nav []Link

	if _, err := os.Stat(filepath.Join(mdDir, "index.md")); err == nil {
		nav = append(nav, Link{Path: "/", Title: "Home"})
	}

	for _, d := range dirList {
		title := strings.ToUpper(d.title[:1]) + d.title[1:]
		nav = append(nav, Link{Path: "/" + d.path + "/", Title: title})
	}

	return nav
}

type URL struct {
	Loc      string `xml:"loc"`
	LastMod  string `xml:"lastmod,omitempty"`
	Priority string `xml:"priority,omitempty"`
}

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	XMLns   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

const sitemapBase = "https://aaron.beago-cirius.ts.net/"

func collectSitemapURLs(root string) ([]URL, error) {
	var urls []URL
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			name := info.Name()
			if name == "node_modules" || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}
		if strings.HasPrefix(rel, "content") || strings.HasPrefix(rel, "media") || strings.HasPrefix(rel, "assets") || strings.HasPrefix(rel, "templates") {
			return nil
		}
		rel = strings.TrimSuffix(rel, ".md")
		rel = strings.ReplaceAll(rel, string(filepath.Separator), "/")
		if strings.HasPrefix(rel, "_") || strings.Contains(rel, "/_") {
			return nil
		}
		loc := sitemapBase + rel
		if strings.HasSuffix(loc, "/index") {
			loc = strings.TrimSuffix(loc, "index")
		} else if rel == "index" {
			loc = sitemapBase
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		_, _, date := stripFrontMatter(content)
		u := URL{Loc: loc}
		if date != "" {
			u.LastMod = date
		} else {
			u.LastMod = info.ModTime().Format("2006-01-02")
		}
		urls = append(urls, u)
		return nil
	})
	return urls, err
}

func buildSitemap(root string) (Sitemap, error) {
	urls, err := collectSitemapURLs(root)
	if err != nil {
		return Sitemap{}, err
	}
	return Sitemap{
		XMLns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}, nil
}

func writeSitemapFile(root, outPath string) error {
	sitemap, err := buildSitemap(root)
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

func handleSitemapFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	path := filepath.Join(mdDir, "sitemap.xml")
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "sitemap not generated: run with -write-sitemap", 404)
			return
		}
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	http.ServeFile(w, r, path)
}
