package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "port to serve on")
	flag.IntVar(&port, "p", 8080, "port to serve on (alias)")

	writeSitemap := flag.Bool("write-sitemap", false, "write sitemap.xml to current directory and exit")
	contentDir := flag.String("content", "content", "directory containing markdown (relative to site root)")
	sitemapBase := flag.String("sitemap-base", os.Getenv("BLOG_SITEMAP_BASE"), "base URL for sitemap loc elements (or set BLOG_SITEMAP_BASE)")
	dev := flag.Bool("dev", false, "reload HTML templates on every request")
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
		Dev:             *dev,
		ParaNumPrefixes: prefixes,
	})
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.handleRequest)
	mux.HandleFunc("/sitemap.xml", srv.handleSitemapFile)

	for _, f := range staticFileNames {
		mux.Handle("/"+f, http.FileServer(http.Dir(wd)))
	}
	mux.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir(filepath.Join(wd, "media")))))

	log.Printf("Serving at http://localhost:%d (content %s)", port, contentRoot)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
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
