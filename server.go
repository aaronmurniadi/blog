package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ServerConfig holds paths and behavior flags for the site generator.
type ServerConfig struct {
	SiteRoot        string
	ContentRoot     string
	SitemapBase     string
	ParaNumPrefixes []string
}

// Server serves markdown pages and static assets.
type Server struct {
	cfg       ServerConfig
	templates *template.Template
	funcMap   template.FuncMap
}

var staticFileNames = []string{
	"style.css",
	"favicon.ico",
	"favicon-32x32.png",
	"favicon-16x16.png",
	"android-chrome-192x192.png",
	"android-chrome-512x512.png",
	"apple-touch-icon.png",
	"robots.txt",
	"site.webmanifest",
}

// NewServer validates roots, parses templates, and returns a Server.
func NewServer(cfg ServerConfig) (*Server, error) {
	siteRoot, err := filepath.Abs(cfg.SiteRoot)
	if err != nil {
		return nil, fmt.Errorf("site root: %w", err)
	}
	contentRoot, err := filepath.Abs(cfg.ContentRoot)
	if err != nil {
		return nil, fmt.Errorf("content root: %w", err)
	}
	cfg.SiteRoot = siteRoot
	cfg.ContentRoot = contentRoot

	if st, err := os.Stat(contentRoot); err != nil || !st.IsDir() {
		if err != nil {
			return nil, fmt.Errorf("content root %s: %w", contentRoot, err)
		}
		return nil, fmt.Errorf("content root is not a directory: %s", contentRoot)
	}

	s := &Server{cfg: cfg, funcMap: newTemplateFuncMap()}
	tmpl, err := s.parseTemplates()
	if err != nil {
		return nil, err
	}
	s.templates = tmpl
	return s, nil
}

func (s *Server) parseTemplates() (*template.Template, error) {
	tmplDir := filepath.Join(s.cfg.SiteRoot, "templates")
	templates := template.New("").Funcs(s.funcMap)
	return templates.ParseGlob(filepath.Join(tmplDir, "*.html"))
}

var errPathTraversal = errors.New("path escapes content root")

// contentRootResolved returns EvalSymlinks path for content root (cached per check in callers).
func (s *Server) contentRootResolved() (string, error) {
	r, err := filepath.EvalSymlinks(s.cfg.ContentRoot)
	if err != nil {
		return "", err
	}
	return filepath.Clean(r), nil
}

// safePathUnderContent maps a relative content path to an absolute path under content root.
// rel uses slash-separated segments, e.g. "articles/foo.md" or "index.md".
func (s *Server) safePathUnderContent(rel string) (string, error) {
	rel = strings.TrimPrefix(filepath.ToSlash(strings.TrimSpace(rel)), "/")
	if rel == "" {
		return "", errPathTraversal
	}
	for _, seg := range strings.Split(rel, "/") {
		if seg == ".." {
			return "", errPathTraversal
		}
	}
	root, err := s.contentRootResolved()
	if err != nil {
		return "", err
	}
	full := filepath.Join(root, filepath.FromSlash(rel))
	full = filepath.Clean(full)

	resolved, err := filepath.EvalSymlinks(full)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		resolved = full
	} else {
		resolved = filepath.Clean(resolved)
	}

	relOut, err := filepath.Rel(root, resolved)
	if err != nil {
		return "", errPathTraversal
	}
	if relOut == ".." || strings.HasPrefix(relOut, ".."+string(os.PathSeparator)) {
		return "", errPathTraversal
	}
	return full, nil
}

func pathHasUnderscoreSegment(p string) bool {
	p = strings.Trim(filepath.ToSlash(p), "/")
	if p == "" {
		return false
	}
	for _, seg := range strings.Split(p, "/") {
		if strings.HasPrefix(seg, "_") {
			return true
		}
	}
	return false
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("%s %s", r.Method, r.URL.Path)

	if pathHasUnderscoreSegment(r.URL.Path) {
		http.NotFound(w, r)
		return
	}

	path := r.URL.Path
	if path == "/" {
		path = "/index"
	}

	if strings.HasSuffix(path, "/") {
		indexPath := path + "index.md"
		rel := strings.TrimPrefix(indexPath, "/")
		abs, err := s.safePathUnderContent(rel)
		if err != nil {
			if errors.Is(err, errPathTraversal) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if _, err := os.Stat(abs); err == nil {
			path = indexPath
		} else {
			s.serveDirListing(w, r, path)
			return
		}
	} else if !strings.HasSuffix(path, ".md") && !strings.HasPrefix(path, "/media/") {
		path += ".md"
	}

	rel := strings.TrimPrefix(path, "/")
	abs, err := s.safePathUnderContent(rel)
	if err != nil {
		if errors.Is(err, errPathTraversal) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	info, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		indexAbs := filepath.Join(abs, "index.md")
		if st, err := os.Stat(indexAbs); err == nil && !st.IsDir() {
			abs = indexAbs
			rel = filepath.ToSlash(filepath.Join(rel, "index.md"))
		} else {
			s.serveDirListing(w, r, strings.TrimSuffix("/"+rel, "/")+"/")
			return
		}
	}

	if err := s.serveMarkdownPage(w, abs, rel, path); err != nil {
		log.Println(err)
	}
}

func (s *Server) handleSitemapFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	path := filepath.Join(s.cfg.SiteRoot, "sitemap.xml")
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "sitemap not generated: run with -write-sitemap", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	http.ServeFile(w, r, path)
}

func (s *Server) contentDir() string { return s.cfg.ContentRoot }

func (s *Server) paraNumPrefixes() []string { return s.cfg.ParaNumPrefixes }
