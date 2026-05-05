package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func outHTMLPath(outRoot, contentRel string) string {
	rel := filepath.ToSlash(contentRel)
	rel = strings.TrimSuffix(rel, ".md")
	if rel == "index" {
		return filepath.Join(outRoot, "index.html")
	}
	if strings.HasSuffix(rel, "/index") {
		parent := strings.TrimSuffix(rel, "/index")
		return filepath.Join(outRoot, filepath.FromSlash(parent), "index.html")
	}
	return filepath.Join(outRoot, filepath.FromSlash(rel), "index.html")
}

func outDirIndexPath(outRoot, urlPath string) string {
	trim := strings.TrimPrefix(filepath.ToSlash(strings.TrimSpace(urlPath)), "/")
	trim = strings.TrimSuffix(trim, "/")
	if trim == "" {
		return filepath.Join(outRoot, "index.html")
	}
	return filepath.Join(outRoot, filepath.FromSlash(trim), "index.html")
}

func shouldSkipContentRel(rel string) bool {
	rel = filepath.ToSlash(rel)
	if rel == "" {
		return false
	}
	if strings.HasPrefix(rel, "_") {
		return true
	}
	return strings.Contains(rel, "/_")
}

func skipContentDirName(name string) bool {
	if name == "" {
		return false
	}
	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
		return true
	}
	switch name {
	case "media", "assets", "scripts", "templates":
		return true
	default:
		return false
	}
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

func (s *Server) generateStaticSite(outRoot string) error {
	if err := os.RemoveAll(outRoot); err != nil {
		return fmt.Errorf("clear output: %w", err)
	}
	if err := os.MkdirAll(outRoot, 0755); err != nil {
		return fmt.Errorf("mkdir output: %w", err)
	}

	for _, name := range staticFileNames {
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
	}
	for _, extra := range []string{"404.html"} {
		src := filepath.Join(s.cfg.SiteRoot, extra)
		if _, err := os.Stat(src); err == nil {
			if err := copyFile(src, filepath.Join(outRoot, extra)); err != nil {
				return fmt.Errorf("copy %s: %w", extra, err)
			}
		}
	}

	mediaSrc := filepath.Join(s.cfg.SiteRoot, "media")
	if st, err := os.Stat(mediaSrc); err == nil && st.IsDir() {
		if err := copyTree(mediaSrc, filepath.Join(outRoot, "media")); err != nil {
			return fmt.Errorf("copy media: %w", err)
		}
	}

	contentRoot := s.cfg.ContentRoot
	err := filepath.WalkDir(contentRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if path == contentRoot {
				return nil
			}
			name := d.Name()
			if skipContentDirName(name) {
				return filepath.SkipDir
			}
			relDir, err := filepath.Rel(contentRoot, path)
			if err != nil {
				return err
			}
			if shouldSkipContentRel(relDir) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		rel, err := filepath.Rel(contentRoot, path)
		if err != nil {
			return err
		}
		relSlash := filepath.ToSlash(rel)
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
		if err := s.writeMarkdownPage(f, path, relSlash, urlPath); err != nil {
			f.Close()
			return fmt.Errorf("%s: %w", relSlash, err)
		}
		if err := f.Close(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = filepath.WalkDir(contentRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() || path == contentRoot {
			return nil
		}
		name := d.Name()
		if skipContentDirName(name) {
			return filepath.SkipDir
		}
		relDir, err := filepath.Rel(contentRoot, path)
		if err != nil {
			return err
		}
		relSlash := filepath.ToSlash(relDir)
		if shouldSkipContentRel(relSlash) {
			return filepath.SkipDir
		}
		indexAbs := filepath.Join(path, "index.md")
		if _, err := os.Stat(indexAbs); err == nil {
			return nil
		}
		shadowMD := filepath.Join(contentRoot, relSlash+".md")
		if _, err := os.Stat(shadowMD); err == nil {
			// e.g. articles.md serves /articles; do not overwrite with articles/ listing
			return nil
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		hasListed := false
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), "_") {
				continue
			}
			if e.IsDir() || strings.HasSuffix(e.Name(), ".md") {
				hasListed = true
				break
			}
		}
		if !hasListed {
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
		return f.Close()
	})
	if err != nil {
		return err
	}

	sitemapOut := filepath.Join(outRoot, "sitemap.xml")
	if err := writeSitemapFile(s.cfg.SiteRoot, contentRoot, s.cfg.SitemapBase, sitemapOut); err != nil {
		return fmt.Errorf("sitemap: %w", err)
	}
	return nil
}
