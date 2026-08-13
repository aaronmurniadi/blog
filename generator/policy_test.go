package generator

import "testing"

func TestSkipContentDirName(t *testing.T) {
	cases := []struct {
		name string
		want bool
	}{
		{"articles", false},
		{"media", true},
		{"assets", true},
		{"scripts", true},
		{"templates", true},
		{"_drafts", true},
		{".hidden", true},
		{"", false},
	}
	for _, c := range cases {
		if got := skipContentDirName(c.name); got != c.want {
			t.Errorf("skipContentDirName(%q) = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestShouldSkipContentRel(t *testing.T) {
	cases := []struct {
		rel  string
		want bool
	}{
		{"index.md", false},
		{"_drafts/foo.md", true},
		{"articles/_private/x.md", true},
		{"articles/private.md", false},
		{"", false},
	}
	for _, c := range cases {
		if got := shouldSkipContentRel(c.rel); got != c.want {
			t.Errorf("shouldSkipContentRel(%q) = %v, want %v", c.rel, got, c.want)
		}
	}
}

func TestSkipContentDir(t *testing.T) {
	cases := []struct {
		rel  string
		name string
		want bool
	}{
		{"articles", "articles", false},
		{"media", "media", true},     // reserved top-level dir pruned here
		{"articles/_x", "_x", true},  // underscore-prefixed segment
		{"_drafts", "_drafts", true}, // underscore-prefixed rel
	}
	for _, c := range cases {
		if got := skipContentDir(c.rel, c.name); got != c.want {
			t.Errorf("skipContentDir(%q, %q) = %v, want %v", c.rel, c.name, got, c.want)
		}
	}
}
