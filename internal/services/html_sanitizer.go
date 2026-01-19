package services

import (
	"github.com/microcosm-cc/bluemonday"
)

type HTMLSanitizer struct {
	policy *bluemonday.Policy
}

func NewHTMLSanitizer() *HTMLSanitizer {
	// Create policy: allow safe HTML tags for rich content
	policy := bluemonday.UGCPolicy()

	// Headings
	policy.AllowElements("h1", "h2", "h3", "h4", "h5", "h6")

	// Text formatting
	policy.AllowElements("p", "div", "span", "br", "hr")
	policy.AllowElements("strong", "em", "u", "s", "del", "ins")
	policy.AllowElements("blockquote", "code", "pre")
	policy.AllowElements("sub", "sup")

	// Lists
	policy.AllowElements("ul", "ol", "li")

	// Tables
	policy.AllowElements("table", "thead", "tbody", "tfoot", "tr", "td", "th", "caption")

	// Links and images
	policy.AllowElements("a", "img")

	// Allow style attributes for formatting
	policy.AllowAttrs("style").OnElements(
		"p", "div", "span", "h1", "h2", "h3", "h4", "h5", "h6",
		"td", "th", "table", "ul", "ol", "li",
	)

	// Allow class attributes
	policy.AllowAttrs("class").OnElements(
		"p", "div", "span", "table", "ul", "ol", "li",
		"h1", "h2", "h3", "h4", "h5", "h6",
	)

	// Allow href on links (with safe protocols)
	policy.AllowAttrs("href", "target", "rel").OnElements("a")
	policy.RequireNoFollowOnLinks(false) // Allow follow links

	// Allow src on images
	policy.AllowAttrs("src", "alt", "title", "width", "height").OnElements("img")

	// Allow table attributes
	policy.AllowAttrs("colspan", "rowspan").OnElements("td", "th")
	policy.AllowAttrs("align", "valign").OnElements("td", "th", "tr")

	// Allow list attributes
	policy.AllowAttrs("type", "start").OnElements("ol")
	policy.AllowAttrs("type").OnElements("ul")

	return &HTMLSanitizer{policy: policy}
}

// Sanitize cleans HTML content and removes potentially dangerous elements
func (s *HTMLSanitizer) Sanitize(html string) string {
	return s.policy.Sanitize(html)
}

// SanitizeBytes is the same as Sanitize but works with byte slices
func (s *HTMLSanitizer) SanitizeBytes(html []byte) []byte {
	return s.policy.SanitizeBytes(html)
}
