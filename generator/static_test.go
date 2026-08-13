package generator

import (
	"path/filepath"
	"testing"
)

func TestOutputIndexPath(t *testing.T) {
	tests := []struct {
		outRoot string
		rel     string
		want    string
	}{
		{"out", "", filepath.Join("out", "index.html")},
		{"out", "index", filepath.Join("out", "index.html")},
		{"out", "articles", filepath.Join("out", "articles", "index.html")},
		{"out", "articles/foo", filepath.Join("out", "articles", "foo", "index.html")},
	}
	for _, tt := range tests {
		if got := outputIndexPath(tt.outRoot, tt.rel); got != tt.want {
			t.Errorf("outputIndexPath(%q, %q) = %q, want %q", tt.outRoot, tt.rel, got, tt.want)
		}
	}
}

func TestOutHTMLPathMatchesOutDirIndexPath(t *testing.T) {
	// articles/foo.md builds to the same file a dir listing for /articles/foo/
	// (modulo the "index" component) would: articles/foo/index.html.
	if got := outHTMLPath("out", "articles/foo.md"); got != filepath.Join("out", "articles", "foo", "index.html") {
		t.Errorf("outHTMLPath = %q", got)
	}
}

func TestOutHTMLPathRootIndex(t *testing.T) {
	if got := outHTMLPath("out", "index.md"); got != filepath.Join("out", "index.html") {
		t.Errorf("outHTMLPath(index.md) = %q", got)
	}
}
