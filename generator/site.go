package generator

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// SiteConfig holds paths and behavior flags for the site generator.
type SiteConfig struct {
	SiteRoot        string
	ContentRoot     string
	SitemapBase     string
	ParaNumPrefixes []string
}

// Site is the site generator core: it resolves content paths, parses
// templates, and renders the static site.
type Site struct {
	cfg       SiteConfig
	templates *template.Template
	funcMap   template.FuncMap
	// contentResolved is the symlink-resolved, cleaned content root, computed
	// once at construction instead of EvalSymlinks on every request.
	contentResolved string
}

var staticFileNames = []string{
	"beago-cirius-logo-white.svg",
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

// NewSite validates roots, parses templates, and returns a Site.
func NewSite(cfg SiteConfig) (*Site, error) {
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

	s := &Site{cfg: cfg, funcMap: newTemplateFuncMap()}
	s.contentResolved, err = filepath.EvalSymlinks(contentRoot)
	if err != nil {
		return nil, fmt.Errorf("content root %s: %w", contentRoot, err)
	}
	s.contentResolved = filepath.Clean(s.contentResolved)
	tmpl, err := s.parseTemplates()
	if err != nil {
		return nil, err
	}
	s.templates = tmpl
	log.Printf("blog: site core ready siteRoot=%s contentRoot=%s", siteRoot, contentRoot)
	return s, nil
}

func (s *Site) parseTemplates() (*template.Template, error) {
	tmplDir := filepath.Join(s.cfg.SiteRoot, "templates")
	templates := template.New("").Funcs(s.funcMap)
	return templates.ParseGlob(filepath.Join(tmplDir, "*.html"))
}

var errPathTraversal = errors.New("path escapes content root")

// safePathUnderContent maps a relative content path to an absolute path under content root.
// rel uses slash-separated segments, e.g. "articles/foo.md" or "index.md".
func (s *Site) safePathUnderContent(rel string) (string, error) {
	rel = strings.TrimPrefix(filepath.ToSlash(strings.TrimSpace(rel)), "/")
	if rel == "" {
		return "", errPathTraversal
	}
	for _, seg := range strings.Split(rel, "/") {
		if seg == ".." {
			return "", errPathTraversal
		}
	}
	root := s.contentResolved
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
