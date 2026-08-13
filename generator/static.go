package generator

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func outHTMLPath(outRoot, contentRel string) string {
	rel := filepath.ToSlash(contentRel)
	rel = strings.TrimSuffix(rel, ".md")
	return outputIndexPath(outRoot, rel)
}

func outDirIndexPath(outRoot, urlPath string) string {
	rel := strings.Trim(filepath.ToSlash(strings.TrimSpace(urlPath)), "/")
	return outputIndexPath(outRoot, rel)
}

// outputIndexPath maps a URL-relative path (slash-separated, no leading/trailing
// slash) to the generated index.html on disk. "" and "index" both map to the
// site root index.html; "articles" maps to articles/index.html.
func outputIndexPath(outRoot, rel string) string {
	rel = strings.TrimSuffix(rel, "/index")
	if rel == "" || rel == "index" {
		return filepath.Join(outRoot, "index.html")
	}
	return filepath.Join(outRoot, filepath.FromSlash(rel), "index.html")
}

func markdownURLPath(contentRel string) string {
	return "/" + filepath.ToSlash(contentRel)
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func copyTree(srcDir, dstDir string) error {
	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dstDir, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFile(path, target)
	})
}

// GenerateStaticSite writes the full static site to outRoot.
func (s *Site) GenerateStaticSite(outRoot string) error {
	if err := os.MkdirAll(outRoot, 0755); err != nil {
		return fmt.Errorf("mkdir output: %w", err)
	}

	passes := []struct {
		name string
		fn   func(outRoot string) error
	}{
		{"static assets", s.copySiteAssets},
		{"media", s.copyMedia},
		{"fonts", s.copyFonts},
		{"markdown", s.renderMarkdownPass},
		{"directory listings", s.renderDirListingPass},
		{"sitemap", s.writeSitemapPass},
	}
	for _, p := range passes {
		log.Printf("site: pass %s start", p.name)
		if err := p.fn(outRoot); err != nil {
			return fmt.Errorf("%s: %w", p.name, err)
		}
		log.Printf("site: pass %s done", p.name)
	}
	log.Printf("site: generate complete")
	return nil
}

func (s *Site) copySiteAssets(outRoot string) error {
	for _, name := range append(append([]string{}, staticFileNames...), "404.html") {
		src := filepath.Join(s.cfg.SiteRoot, name)
		if _, err := os.Stat(src); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		dst := filepath.Join(outRoot, name)
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("copy %s: %w", name, err)
		}
		log.Printf("site: static %s", name)
	}
	return nil
}

func (s *Site) copyMedia(outRoot string) error {
	mediaSrc := filepath.Join(s.cfg.SiteRoot, "media")
	st, err := os.Stat(mediaSrc)
	if err != nil || !st.IsDir() {
		log.Printf("site: skip media (missing or not a directory)")
		return nil
	}
	return copyMediaTreeAsWebP(mediaSrc, filepath.Join(outRoot, "media"))
}

func (s *Site) copyFonts(outRoot string) error {
	fontsSrc := filepath.Join(s.cfg.SiteRoot, "fonts")
	st, err := os.Stat(fontsSrc)
	if err != nil || !st.IsDir() {
		log.Printf("site: skip fonts (missing or not a directory)")
		return nil
	}
	return copyTree(fontsSrc, filepath.Join(outRoot, "fonts"))
}

func (s *Site) renderMarkdownPass(outRoot string) error {
	contentRoot := s.cfg.ContentRoot
	return filepath.WalkDir(contentRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if path == contentRoot {
				return nil
			}
			if skipContentDir(filepath.ToSlash(mustRel(contentRoot, path)), d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		relSlash := filepath.ToSlash(mustRel(contentRoot, path))
		if shouldSkipContentRel(relSlash) {
			return nil
		}
		outPath := outHTMLPath(outRoot, relSlash)
		urlPath := markdownURLPath(relSlash)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}
		f, err := os.Create(outPath)
		if err != nil {
			return err
		}
		if err := s.writeMarkdownPage(f, path, relSlash, urlPath, true); err != nil {
			f.Close()
			return fmt.Errorf("%s: %w", relSlash, err)
		}
		if err := f.Close(); err != nil {
			return err
		}
		log.Printf("site: page %s -> %s", relSlash, outPath)
		return nil
	})
}

func (s *Site) renderDirListingPass(outRoot string) error {
	contentRoot := s.cfg.ContentRoot
	return filepath.WalkDir(contentRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() || path == contentRoot {
			return nil
		}
		relSlash := filepath.ToSlash(mustRel(contentRoot, path))
		if skipContentDir(relSlash, d.Name()) {
			return filepath.SkipDir
		}
		if _, err := os.Stat(filepath.Join(path, "index.md")); err == nil {
			return nil
		}
		if _, err := os.Stat(filepath.Join(contentRoot, relSlash+".md")); err == nil {
			// e.g. articles.md serves /articles; do not overwrite with articles/ listing
			return nil
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		if !dirHasListableEntries(entries) {
			return nil
		}

		urlPath := "/" + relSlash + "/"
		outPath := outDirIndexPath(outRoot, urlPath)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}
		f, err := os.Create(outPath)
		if err != nil {
			return err
		}
		if err := s.writeDirListing(f, urlPath); err != nil {
			f.Close()
			return fmt.Errorf("dir index %s: %w", relSlash, err)
		}
		if err := f.Close(); err != nil {
			return err
		}
		log.Printf("site: dir index %s -> %s", urlPath, outPath)
		return nil
	})
}

func (s *Site) writeSitemapPass(outRoot string) error {
	sitemapOut := filepath.Join(outRoot, "sitemap.xml")
	return WriteSitemapFile(s.cfg.ContentRoot, s.cfg.SitemapBase, sitemapOut)
}

func dirHasListableEntries(entries []fs.DirEntry) bool {
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "_") {
			continue
		}
		if e.IsDir() || strings.HasSuffix(e.Name(), ".md") {
			return true
		}
	}
	return false
}

func mustRel(base, target string) string {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return target
	}
	return rel
}
