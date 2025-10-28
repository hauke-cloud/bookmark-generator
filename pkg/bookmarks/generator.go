package bookmarks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hauke-cloud/bookmark-generator/pkg/kubernetes"
)

// Generator handles bookmark generation
type Generator struct{}

// NewGenerator creates a new bookmark generator
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateFirefox creates a Firefox-compatible HTML bookmark file
func (g *Generator) GenerateFirefox(ingresses []kubernetes.IngressInfo) []byte {
	timestamp := time.Now().Unix()

	html := fmt.Sprintf(`<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks Menu</H1>

<DL><p>
    <DT><H3 ADD_DATE="%d" LAST_MODIFIED="%d">Kubernetes Ingresses</H3>
    <DL><p>
`, timestamp, timestamp)

	for _, ingress := range ingresses {
		title := fmt.Sprintf("%s (%s)", ingress.Host, ingress.Namespace)
		html += fmt.Sprintf(`        <DT><A HREF="%s" ADD_DATE="%d">%s</A>
`, ingress.URL, timestamp, escapeHTML(title))
	}

	html += `    </DL><p>
</DL><p>
`

	return []byte(html)
}

// GenerateChrome creates a Chrome-compatible JSON bookmark file
func (g *Generator) GenerateChrome(ingresses []kubernetes.IngressInfo) ([]byte, error) {
	timestamp := time.Now().Unix()

	children := make([]map[string]interface{}, 0, len(ingresses))
	for _, ingress := range ingresses {
		title := fmt.Sprintf("%s (%s)", ingress.Host, ingress.Namespace)
		children = append(children, map[string]interface{}{
			"date_added": fmt.Sprintf("%d000000", timestamp), // Chrome uses microseconds
			"guid":       generateGUID(ingress.URL),
			"id":         fmt.Sprintf("%d", len(children)+1),
			"name":       title,
			"type":       "url",
			"url":        ingress.URL,
		})
	}

	bookmarks := map[string]interface{}{
		"checksum": "computed_checksum",
		"roots": map[string]interface{}{
			"bookmark_bar": map[string]interface{}{
				"children": []map[string]interface{}{
					{
						"date_added":    fmt.Sprintf("%d000000", timestamp),
						"date_modified": fmt.Sprintf("%d000000", timestamp),
						"guid":          "00000000-0000-4000-a000-000000000001",
						"id":            "1",
						"name":          "Kubernetes Ingresses",
						"type":          "folder",
						"children":      children,
					},
				},
				"date_added":    fmt.Sprintf("%d000000", timestamp),
				"date_modified": fmt.Sprintf("%d000000", timestamp),
				"guid":          "00000000-0000-4000-a000-000000000000",
				"id":            "0",
				"name":          "Bookmarks bar",
				"type":          "folder",
			},
			"other": map[string]interface{}{
				"children":      []map[string]interface{}{},
				"date_added":    fmt.Sprintf("%d000000", timestamp),
				"date_modified": fmt.Sprintf("%d000000", timestamp),
				"guid":          "00000000-0000-4000-a000-000000000002",
				"id":            "2",
				"name":          "Other bookmarks",
				"type":          "folder",
			},
			"synced": map[string]interface{}{
				"children":      []map[string]interface{}{},
				"date_added":    fmt.Sprintf("%d000000", timestamp),
				"date_modified": fmt.Sprintf("%d000000", timestamp),
				"guid":          "00000000-0000-4000-a000-000000000003",
				"id":            "3",
				"name":          "Mobile bookmarks",
				"type":          "folder",
			},
		},
		"version": 1,
	}

	return json.MarshalIndent(bookmarks, "", "  ")
}

func escapeHTML(s string) string {
	// Basic HTML escaping
	replacements := map[rune]string{
		'&':  "&amp;",
		'<':  "&lt;",
		'>':  "&gt;",
		'"':  "&quot;",
		'\'': "&#39;",
	}

	result := ""
	for _, c := range s {
		if replacement, ok := replacements[c]; ok {
			result += replacement
		} else {
			result += string(c)
		}
	}
	return result
}

func generateGUID(url string) string {
	// Simple GUID generation based on URL hash
	// In production, you might want to use a proper UUID library
	hash := 0
	for _, c := range url {
		hash = (hash * 31) + int(c)
	}
	return fmt.Sprintf("%08x-0000-4000-a000-%012x", hash&0xffffffff, hash&0xffffffffffff)
}
