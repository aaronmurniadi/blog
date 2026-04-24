package main

import (
	"bytes"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
)

type Page struct {
	Title   string
	Path    string
	HTML    template.HTML
	Nav     []Link
}

type Link struct {
	Path  string
	Title string
}

var (
	mdDir      string
	tmplDir    string
	funcMap   = template.FuncMap{"trim": strings.TrimSpace}
	templates *template.Template
)

const defaultHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<style>
:root{--bg:#fafafa;--fg:#222;--a:#0066cc;--code:#f4f4f4}
@media(prefers-color-scheme:dark){:root{--bg:#1a1a1a;--f0:#ddd;--fg:#ddd;--a:#7ab3ff;--code:#2d2d2d}}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Oxygen,Ubuntu,sans-serif;max-width:800px;margin:0 auto;padding:2rem;background:var(--bg);color:var(--fg);line-height:1.6}
a{color:var(--a);text-decoration:none}
a:hover{text-decoration:underline}
nav{margin-bottom:2rem;padding-bottom:1rem;border-bottom:1px solid #ddd;display:flex;flex-wrap:wrap;gap:1rem}
code{background:var(--code);padding:0.2rem 0.4rem;border-radius:4px;font-size:0.9em}
pre{background:var(--code);padding:1rem;border-radius:4px;overflow-x:auto}
pre code{background:none;padding:0}
img{max-width:100%;height:auto}
h1,h2,h3{margin-top:1.5em}
hr{border:none;border-top:1px solid #ddd;margin:2rem 0}
blockquote{border-left:3px solid #ddd;margin:1rem 0;padding-left:1rem;color:#666}
table{border-collapse:collapse;width:100%;margin:1rem 0}
th,td{border:1px solid #ddd;padding:0.5rem;text-align:left}
</style>
</head>
<body>
<nav>{{range .Nav}}<a href="{{.Path}}">{{.Title}}</a>{{end}}</nav>
<main>{{.HTML}}</main>
</body>
</html>`

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	mdDir = wd
	tmplDir = filepath.Join(wd, "templates")

	if _, err := os.Stat(tmplDir); os.IsNotExist(err) {
		os.MkdirAll(tmplDir, 0755)
	}

	defaultTmpl := filepath.Join(tmplDir, "default.html")
	if _, err := os.Stat(defaultTmpl); os.IsNotExist(err) {
		os.WriteFile(defaultTmpl, []byte(defaultHTML), 0644)
	}

	templates = template.New("").Funcs(funcMap)
	templates, err = templates.ParseGlob(filepath.Join(tmplDir, "*.html"))
	if err != nil {
		log.Fatal(err)
	}

	nav := buildNav("")
	log.Println("Nav:", nav)

	http.HandleFunc("/", handleRequest)
	log.Println("Serving at http://localhost:8003")
	log.Fatal(http.ListenAndServe(":8003", nil))
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

	mdContent, pageTitle := stripFrontMatter(content)

	html, err := mdToHTML(mdContent)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	title := pageTitle
	if title == "" {
		title = extractTitle(absPath)
	}
	nav := buildNav(absPath)
	page := Page{Title: title, Path: path, HTML: template.HTML(html), Nav: nav}

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
			_, title := stripFrontMatter(content)
			ep := filepath.Join(relDir, e.Name())
			ep = strings.TrimPrefix(ep, string(filepath.Separator))
			ep = "/" + strings.TrimSuffix(ep, ".md")
			links = append(links, Link{Path: ep, Title: title})
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
		b.WriteString("</a></li>")
	}
	return b.String()
}

func stripFrontMatter(md []byte) ([]byte, string) {
	if !bytes.HasPrefix(bytes.TrimSpace(md), []byte("---")) {
		return md, extractH1(md)
	}

	lines := strings.Split(string(md), "\n")
	if len(lines) < 3 {
		return md, extractH1(md)
	}

	contentStart := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			contentStart = i + 1
			break
		}
	}

	if contentStart == -1 {
		return md, extractH1(md)
	}

	for _, line := range lines[1:contentStart-1] {
		if strings.HasPrefix(line, "title:") {
			title := strings.TrimSpace(strings.TrimPrefix(line, "title:"))
			title = strings.Trim(title, `"`)
			title = strings.Trim(title, `'`)
			content := []byte(strings.Join(lines[contentStart:], "\n"))
			return bytes.TrimSpace(content), title
		}
	}

	content := []byte(strings.Join(lines[contentStart:], "\n"))
	return bytes.TrimSpace(content), extractH1(content)
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
	markdown := goldmark.New()
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
	dirs := make(map[string]string)

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
		dirs[rel] = name
		return nil
	})

	var nav []Link

	if _, err := os.Stat(filepath.Join(mdDir, "index.md")); err == nil {
		nav = append(nav, Link{Path: "/", Title: "Home"})
	}

	for dir, title := range dirs {
		title = strings.ToUpper(title[:1]) + title[1:]
		nav = append(nav, Link{Path: "/" + dir + "/", Title: title})
	}

	return nav
}