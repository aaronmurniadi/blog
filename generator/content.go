package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Page is data for default.html.
type Page struct {
	Title       string
	Path        string
	HTML        template.HTML
	Nav         []Link
	HideNav     bool
	Description string
	Canonical   string
	Image       string
	SiteName    string
	SiteImage   string
	Lang        string
	Type        string
}

// Link is a nav or directory-listing entry.
type Link struct {
	Path     string
	Title    string
	Date     string
	SortDate string
}

// DirIndexData is passed to dirindex.html.
type DirIndexData struct {
	DirTitle string
	Links    []Link
}

// FrontMatter holds parsed YAML front matter fields we care about.
type FrontMatter struct {
	Title       string `yaml:"title"`
	Date        string `yaml:"date"`
	NavBar      *bool  `yaml:"nav_bar"` // nil = show nav; explicit false hides <nav>
	Description string `yaml:"description"`
	Summary     string `yaml:"summary"`
	Image       string `yaml:"image"`
	Lang        string `yaml:"lang"`
}

// metaDescription returns the front-matter description if set, else summary,
// else an auto-extracted plain-text summary from the markdown body.
func (fm FrontMatter) metaDescription(body []byte) string {
	if fm.Description != "" {
		return fm.Description
	}
	if fm.Summary != "" {
		return fm.Summary
	}
	return autoDescription(body)
}

// parseFrontMatter splits YAML front matter from body; body is trimmed.
// If there is no well-formed `---` … `---` block, returns full src and empty fm.
func parseFrontMatter(src []byte) (body []byte, fm FrontMatter, err error) {
	trimmed := bytes.TrimSpace(src)
	if !bytes.HasPrefix(trimmed, []byte("---")) {
		return src, fm, nil
	}
	rest := bytes.TrimPrefix(trimmed, []byte("---"))
	rest = bytes.TrimPrefix(rest, []byte("\r\n"))
	rest = bytes.TrimPrefix(rest, []byte("\n"))

	var sep []byte
	var closeIdx int
	if i := bytes.Index(rest, []byte("\n---")); i >= 0 {
		sep = []byte("\n---")
		closeIdx = i
	} else if i := bytes.Index(rest, []byte("\r\n---")); i >= 0 {
		sep = []byte("\r\n---")
		closeIdx = i
	} else {
		return src, fm, nil
	}

	metaBlock := bytes.TrimSpace(rest[:closeIdx])
	body = bytes.TrimSpace(rest[closeIdx+len(sep):])
	if err := yaml.Unmarshal(metaBlock, &fm); err != nil {
		return nil, fm, fmt.Errorf("front matter: %w", err)
	}
	return body, fm, nil
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

func extractTitleFromPath(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, ".md")
}

// stripFencedCode removes triple-backtick or triple-tilde code fences (and
// their contents) so code is never used to derive a meta description.
func stripFencedCode(src string) string {
	lines := strings.Split(src, "\n")
	out := make([]string, 0, len(lines))
	inFence := false
	var fenceChar byte
	for _, ln := range lines {
		t := strings.TrimSpace(ln)
		if !inFence {
			if len(t) >= 3 && ((t[0] == '`' || t[0] == '~') && isAllSameChar(t)) {
				inFence = true
				fenceChar = t[0]
				continue
			}
			out = append(out, ln)
			continue
		}
		// inside a fence: only a closing run of the same marker ends it
		if len(t) >= 3 && t[0] == fenceChar && isAllSameChar(t) {
			inFence = false
		}
	}
	return strings.Join(out, "\n")
}

func isAllSameChar(s string) bool {
	if len(s) == 0 {
		return false
	}
	c := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != c {
			return false
		}
	}
	return true
}

// autoDescription produces a short plain-text excerpt from a markdown body by
// stripping markup, taking the first substantive paragraph, collapsing
// whitespace, and truncating to about maxLen runes at a word boundary.
func autoDescription(md []byte, maxLen ...int) string {
	n := 160
	if len(maxLen) > 0 && maxLen[0] > 0 {
		n = maxLen[0]
	}
	text := string(md)
	text = stripFencedCode(text)
	// Heading, image, and horizontal-rule lines carry no useful summary text.
	h := regexp.MustCompile("(?m)^#{1,6}\\s+.*$")
	text = h.ReplaceAllString(text, "")
	img := regexp.MustCompile("(?m)^\\s*!\\[[^]]*]\\([^)]*\\)\\s*$")
	text = img.ReplaceAllString(text, "")
	text = strings.ReplaceAll(text, "---", "\n")
	// Strip inline links/images and markdown emphasis/backticks.
	text = regexp.MustCompile(`!\[[^\]]*]\([^)]*\)`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`\[([^\]]+)]\([^)]*\)`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`[` + "`" + `*~]`).ReplaceAllString(text, "")
	// Blockquote and list markers.
	text = regexp.MustCompile("(?m)^\\s*>\\s?").ReplaceAllString(text, "")
	text = regexp.MustCompile("(?m)^\\s*[-+*]\\s+").ReplaceAllString(text, "")
	text = regexp.MustCompile("(?m)^\\s*\\d+\\.\\s+").ReplaceAllString(text, "")
	// Escape any stray HTML was preserved; strip tags.
	text = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(text, "")
	text = strings.NewReplacer("\u00a0", " ").Replace(text)

	// Take first block containing any text.
	var excerpt []rune
	for _, para := range strings.Split(text, "\n") {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}
		excerpt = []rune(para)
		break
	}
	if len(excerpt) == 0 {
		return ""
	}
	// Collapse internal whitespace.
	sw := regexp.MustCompile(`\s+`)
	s := sw.ReplaceAllString(string(excerpt), " ")
	if runes := []rune(s); len(runes) > n {
		s = string(runes[:n])
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	return s
}

// buildNav walks the content tree and returns top-level section links plus Home when index.md exists.
func buildNav(s *Site) ([]Link, error) {
	type dirInfo struct {
		path  string
		title string
	}
	var dirList []dirInfo

	contentRoot := s.cfg.ContentRoot
	err := filepath.WalkDir(contentRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if path == contentRoot {
			return nil
		}
		name := d.Name()
		if skipContentDirName(name) {
			return filepath.SkipDir
		}
		rel, err := filepath.Rel(contentRoot, path)
		if err != nil {
			return err
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		hasMD := false
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
	if err != nil {
		return nil, err
	}

	sort.Slice(dirList, func(i, j int) bool {
		return dirList[i].path < dirList[j].path
	})

	var nav []Link
	indexPath := filepath.Join(contentRoot, "index.md")
	if st, err := os.Stat(indexPath); err == nil && !st.IsDir() {
		nav = append(nav, Link{Path: "/", Title: "Home"})
	}

	for _, d := range dirList {
		rel := filepath.ToSlash(d.path)
		title := titleCaseDir(d.title)
		nav = append(nav, Link{Path: "/" + rel + "/", Title: title})
	}

	return nav, nil
}

func titleCaseDir(name string) string {
	if name == "" {
		return name
	}
	return strings.ToUpper(name[:1]) + strings.ToLower(name[1:])
}

func linkSortTime(date string) string {
	return canonicalDate(date)
}

func sortLinksByDateDesc(links []Link) {
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
}

// linksForMarkdownFilesInDir returns one Link per *.md file directly in absDir (non-recursive).
func (s *Site) linksForMarkdownFilesInDir(absDir string) ([]Link, error) {
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return nil, err
	}
	contentRoot := s.cfg.ContentRoot
	var links []Link
	for _, e := range entries {
		if e.IsDir() || strings.HasPrefix(e.Name(), "_") {
			continue
		}
		if !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		entryPath := filepath.Join(absDir, e.Name())
		relItem, err := filepath.Rel(contentRoot, entryPath)
		if err != nil {
			continue
		}
		relItem = filepath.ToSlash(relItem)
		content, err := os.ReadFile(entryPath)
		if err != nil {
			continue
		}
		body, meta, perr := parseFrontMatter(content)
		if perr != nil {
			continue
		}
		title := meta.Title
		if title == "" {
			title = extractH1(body)
		}
		if title == "" {
			title = titleCaseDir(strings.TrimSuffix(e.Name(), ".md"))
		}
		web := "/" + strings.TrimSuffix(relItem, ".md")
		date := meta.Date
		links = append(links, Link{
			Path:     web,
			Title:    title,
			Date:     date,
			SortDate: linkSortTime(date),
		})
	}
	sortLinksByDateDesc(links)
	return links, nil
}

// writeDirListing lists child directories and .md files for a path ending in /.
func (s *Site) writeDirListing(w io.Writer, urlPath string) error {
	rel := strings.TrimPrefix(filepath.ToSlash(strings.TrimSpace(urlPath)), "/")
	rel = strings.TrimSuffix(rel, "/")
	absDir, err := s.safePathUnderContent(rel)
	if err != nil {
		return err
	}
	st, err := os.Stat(absDir)
	if err != nil {
		return err
	}
	if !st.IsDir() {
		return fmt.Errorf("%s is not a directory", absDir)
	}

	entries, err := os.ReadDir(absDir)
	if err != nil {
		return err
	}

	var links []Link
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "_") {
			continue
		}
		entryPath := filepath.Join(absDir, e.Name())
		relItem, _ := filepath.Rel(s.cfg.ContentRoot, entryPath)
		relItem = filepath.ToSlash(relItem)

		if e.IsDir() {
			links = append(links, Link{
				Path:  "/" + relItem + "/",
				Title: titleCaseDir(e.Name()),
			})
			continue
		}
	}

	fileLinks, err := s.linksForMarkdownFilesInDir(absDir)
	if err != nil {
		return err
	}
	links = append(links, fileLinks...)

	sortLinksByDateDesc(links)

	dirTitle := titleCaseDir(filepath.Base(absDir))
	if rel == "" {
		dirTitle = "Home"
	}

	nav, err := buildNav(s)
	if err != nil {
		log.Println("buildNav:", err)
		nav = nil
	}
	var body bytes.Buffer
	data := DirIndexData{DirTitle: dirTitle, Links: links}
	if err := s.templates.ExecuteTemplate(&body, "dirindex.html", data); err != nil {
		return fmt.Errorf("dirindex template: %w", err)
	}

	page := Page{
		Title:       dirTitle,
		Path:        urlPath,
		HTML:        template.HTML(body.String()),
		Nav:         nav,
		Canonical:   s.canonicalURL(urlPath),
		Description: "Archive of pages in the " + dirTitle + " section of Beago Cirius.",
		Image:       "/avatar.jpg",
		SiteName:    "Beago Cirius",
		SiteImage:   "/beago-cirius-logo-white.svg",
		Lang:        "en",
		Type:        "website",
	}
	if err := s.templates.ExecuteTemplate(w, "default.html", page); err != nil {
		return fmt.Errorf("execute default.html: %w", err)
	}
	return nil
}

// shouldApplyParaNum returns true if the URL's content path should use the .para-num wrapper.
func (s *Site) shouldApplyParaNum(webRel string) bool {
	webRel = filepath.ToSlash(webRel)
	for _, p := range s.cfg.ParaNumPrefixes {
		if strings.HasPrefix(webRel, p) {
			return true
		}
	}
	return false
}
