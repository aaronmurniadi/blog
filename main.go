package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	writeSitemap := flag.Bool("write-sitemap", false, "write sitemap.xml to current directory and exit")
	outDir := flag.String("out", "_site", "output directory for generated static site")
	contentDir := flag.String("content", "content", "directory containing markdown (relative to site root)")
	sitemapBase := flag.String("sitemap-base", os.Getenv("BLOG_SITEMAP_BASE"), "base URL for sitemap loc elements (or set BLOG_SITEMAP_BASE)")
	paraNumPaths := flag.String("para-num-paths", "articles/,summaries/", "comma-separated URL path prefixes that get .para-num wrapper")
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	wd, err = filepath.Abs(wd)
	if err != nil {
		log.Fatal(err)
	}

	contentRoot := filepath.Join(wd, *contentDir)
	if *sitemapBase == "" {
		*sitemapBase = "https://aaron.beago-cirius.ts.net/"
	}
	if !strings.HasSuffix(*sitemapBase, "/") {
		*sitemapBase += "/"
	}

	prefixes := parseCommaPrefixes(*paraNumPaths)

	if *writeSitemap {
		out := filepath.Join(wd, "sitemap.xml")
		if err := writeSitemapFile(wd, contentRoot, *sitemapBase, out); err != nil {
			log.Fatal(err)
		}
		log.Printf("wrote %s", out)
		return
	}

	srv, err := NewServer(ServerConfig{
		SiteRoot:        wd,
		ContentRoot:     contentRoot,
		SitemapBase:     *sitemapBase,
		ParaNumPrefixes: prefixes,
	})
	if err != nil {
		log.Fatal(err)
	}

	outAbs := filepath.Join(wd, *outDir)
	if err := srv.generateStaticSite(outAbs); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote site to %s", outAbs)
}

func parseCommaPrefixes(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if !strings.HasSuffix(p, "/") {
			p += "/"
		}
		out = append(out, p)
	}
	return out
}
