package generator

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseFrontMatter(t *testing.T) {
	tests := []struct {
		name        string
		src         string
		wantTitle   string
		wantDate    string
		wantBodySub string
	}{
		{
			name:        "no front matter",
			src:         "# Hello\n\njust body\n",
			wantTitle:   "",
			wantBodySub: "# Hello",
		},
		{
			name:        "full front matter",
			src:         "---\ntitle: My Post\ndate: 2026-01-15\n---\n\nBody content\n",
			wantTitle:   "My Post",
			wantDate:    "2026-01-15",
			wantBodySub: "Body content",
		},
		{
			name:        "front matter without closing delimit lines",
			src:         "---\ntitle: Broken\nno closing\n",
			wantTitle:   "",
			wantBodySub: "title: Broken", // still passed through as body
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, fm, err := parseFrontMatter([]byte(tt.src))
			if err != nil {
				t.Fatalf("parseFrontMatter error: %v", err)
			}
			if fm.Title != tt.wantTitle {
				t.Errorf("title = %q, want %q", fm.Title, tt.wantTitle)
			}
			if fm.Date != tt.wantDate {
				t.Errorf("date = %q, want %q", fm.Date, tt.wantDate)
			}
			if tt.wantBodySub != "" && !bytes.Contains(body, []byte(tt.wantBodySub)) {
				t.Errorf("body %q does not contain %q", string(body), tt.wantBodySub)
			}
		})
	}
}

func TestCanonicalDate(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"2026-01-15", "2026-01-15"},
		{"January 2, 2006", "2006-01-02"},
		{"Jan 2, 2006", "2006-01-02"},
		{"not a date", ""},
	}
	for _, tt := range tests {
		if got := canonicalDate(tt.in); got != tt.want {
			t.Errorf("canonicalDate(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestExtractH1(t *testing.T) {
	src := "intro\ntext\n# The Title\nmore\n"
	if got := extractH1([]byte(src)); got != "The Title" {
		t.Errorf("extractH1 = %q, want %q", got, "The Title")
	}
	if strings.TrimSpace(extractH1([]byte("title: not here"))) != "" {
		t.Errorf("expected no H1 in front-matter-only text")
	}
}
