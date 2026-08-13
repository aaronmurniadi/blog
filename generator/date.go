package generator

import "time"

// dateCanonicalLayout is the normal form used for sitemap lastmod and link
// sorting. All accepted input layouts parse to this.
const dateCanonicalLayout = "2006-01-02"

// dateParseLayouts are the front-matter date formats accepted when reading posts.
// Add a new accepted format here; nothing downstream needs to change.
var dateParseLayouts = []string{
	"2006-01-02",
	"January 2, 2006",
	"Jan 2, 2006",
}

// canonicalDate parses an arbitrary accepted date string into the canonical
// layout form, returning "" if it is unparseable or empty.
func canonicalDate(date string) string {
	if date == "" {
		return ""
	}
	for _, layout := range dateParseLayouts {
		if t, err := time.ParseInLocation(layout, date, time.Local); err == nil {
			return t.Format(dateCanonicalLayout)
		}
	}
	return ""
}
