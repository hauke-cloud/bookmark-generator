package bookmarks

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/hauke-cloud/bookmark-generator/pkg/kubernetes"
)

func TestGenerateFirefox(t *testing.T) {
	generator := NewGenerator()

	ingresses := []kubernetes.IngressInfo{
		{
			Name:      "test-ingress",
			Namespace: "default",
			Host:      "example.com",
			Path:      "/",
			URL:       "https://example.com/",
		},
		{
			Name:      "api-ingress",
			Namespace: "production",
			Host:      "api.example.com",
			Path:      "/v1",
			URL:       "https://api.example.com/v1",
		},
	}

	result := generator.GenerateFirefox(ingresses)
	html := string(result)

	// Check for DOCTYPE
	if !strings.Contains(html, "<!DOCTYPE NETSCAPE-Bookmark-file-1>") {
		t.Error("Firefox bookmark should contain DOCTYPE")
	}

	// Check for bookmarks
	if !strings.Contains(html, "https://example.com/") {
		t.Error("Firefox bookmark should contain first URL")
	}

	if !strings.Contains(html, "https://api.example.com/v1") {
		t.Error("Firefox bookmark should contain second URL")
	}

	// Check for metadata
	if !strings.Contains(html, "example.com (default)") {
		t.Error("Firefox bookmark should contain bookmark title")
	}
}

func TestGenerateChrome(t *testing.T) {
	generator := NewGenerator()

	ingresses := []kubernetes.IngressInfo{
		{
			Name:      "test-ingress",
			Namespace: "default",
			Host:      "example.com",
			Path:      "/",
			URL:       "https://example.com/",
		},
	}

	result, err := generator.GenerateChrome(ingresses)
	if err != nil {
		t.Fatalf("GenerateChrome failed: %v", err)
	}

	// Parse as JSON to verify structure
	var bookmarks map[string]interface{}
	if err := json.Unmarshal(result, &bookmarks); err != nil {
		t.Fatalf("Chrome bookmark should be valid JSON: %v", err)
	}

	// Check for roots
	roots, ok := bookmarks["roots"].(map[string]interface{})
	if !ok {
		t.Fatal("Chrome bookmark should have roots")
	}

	// Check for bookmark_bar
	bookmarkBar, ok := roots["bookmark_bar"].(map[string]interface{})
	if !ok {
		t.Fatal("Chrome bookmark should have bookmark_bar")
	}

	// Check for children
	children, ok := bookmarkBar["children"].([]interface{})
	if !ok {
		t.Fatal("bookmark_bar should have children")
	}

	if len(children) != 1 {
		t.Errorf("Expected 1 folder in bookmark_bar, got %d", len(children))
	}
}

func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"<script>", "&lt;script&gt;"},
		{"a&b", "a&amp;b"},
		{"'quote'", "&#39;quote&#39;"},
		{`"double"`, "&quot;double&quot;"},
	}

	for _, tt := range tests {
		result := escapeHTML(tt.input)
		if result != tt.expected {
			t.Errorf("escapeHTML(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
