package generator

import (
	"path/filepath"
	"strings"
)

// reservedDirNames are content directories that are skipped in nav, generation,
// and sitemap. Keep this single source of truth in sync with the description
// in AGENTS.md ("Directories named media, assets, or scripts under content are
// skipped both in nav and generation.").
var reservedDirNames = []string{"media", "assets", "scripts", "templates"}

// hiddenSegmentsPrefixes are directory name prefixes treated as private/unpublished.
const publishableDirPrefixes = "._"

// skipContentDirName reports whether a top-level content dir name should be
// skipped everywhere (nav, generation, sitemap). Used for one directory-level
// segment; caller applies it to names as they are encountered.
func skipContentDirName(name string) bool {
	if name == "" {
		return false
	}
	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
		return true
	}
	return containsStr(reservedDirNames, name)
}

// shouldSkipContentRel reports whether a slash-separated content-relative path
// (file or dir) is private and must not be published. Handles both a leading
// "_prefix" and an embedded "/_segment". The "."-prefixed dirs are covered by
// skipContentDirName at the walk level; this function only guards "_" segments.
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

// skipContentDir reports whether a content directory (given its parent-relative
// slash path and its base name) should be pruned from nav, generation, and
// sitemap. Combines the segment-value and segment-prefix rules.
func skipContentDir(relSlash, name string) bool {
	if skipContentDirName(name) {
		return true
	}
	return shouldSkipContentRel(relSlash)
}

// skipSitemapDirName extends skipContentDirName with build-only exclusions that
// only matter for sitemap walking (e.g. dependency dirs) if the site content
// ever nests such a directory.
func skipSitemapDirName(name string) bool {
	if name == "node_modules" {
		return true
	}
	return skipContentDirName(name)
}

func containsStr(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
