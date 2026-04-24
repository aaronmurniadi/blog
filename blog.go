package main

import (
	"bytes"
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

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type Page struct {
	Title string
	Path  string
	HTML  template.HTML
	Nav   []Link
}

type Link struct {
	Path  string
	Title string
	Date  string
}

var (
	mdDir      string
	tmplDir    string
	port       int
	funcMap    = template.FuncMap{"trim": strings.TrimSpace}
	templates  *template.Template
)

func main() {
	flag.IntVar(&port, "port", 8080, "port to serve on")
	flag.IntVar(&port, "p", 8080, "port to serve on")
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	mdDir = wd
	tmplDir = filepath.Join(wd, "templates")

	templates = template.New("").Funcs(funcMap)
	templates, err = templates.ParseGlob(filepath.Join(tmplDir, "*.html"))
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handleRequest)
	http.Handle("/style.css", http.FileServer(http.Dir(wd)))
	log.Printf("Serving at http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
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
	} else if !strings.HasSuffix(path, ".md") {
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
	headerTitle := "Beago Cirius"
	nav := buildNav(absPath)
	pageHTML := template.HTML(html)
	if date != "" {
		pageHTML = template.HTML("<p class=\"subtitle\">" + date + "</p>" + string(html))
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
			links = append(links, Link{Path: ep, Title: title, Date: date})
		}
	}

	title := filepath.Base(dirPath)
	nav := buildNav("")
	html := "<ul>" + linksToHTML(links) + "</ul>"
	page := Page{Title: title, Path: r.URL.Path, HTML: template.HTML(html), Nav: nav}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, "default.html", page); err != nil {
		log.Println(err)
	}
}

func linksToHTML(links []Link) string {
	var b strings.Builder
	for _, l := range links {
		b.WriteString(`<li><a href="`)
		b.WriteString(l.Path)
		b.WriteString(`">`)
		b.WriteString(l.Title)
		b.WriteString("</a>")
		if l.Date != "" {
			b.WriteString(` <span style="color:#666;font-size:0.9em">`)
			b.WriteString(l.Date)
			b.WriteString(`</span>`)
		}
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
	md = bytes.ReplaceAll(md, []byte("---"), []byte("—"))
	markdown := goldmark.New(
		goldmark.WithExtensions(extension.Footnote, extension.Table),
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

func buildNav(currentPath string) []Link {
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