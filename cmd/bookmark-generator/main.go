package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hauke-cloud/bookmark-generator/internal/handlers"
	"github.com/hauke-cloud/bookmark-generator/pkg/kubernetes"
)

//go:embed templates/*.html
var templatesFS embed.FS

func main() {
	log.Println("Starting Kubernetes Bookmark Generator...")

	// Get configuration from environment
	port := getEnv("PORT", "8080")

	// Initialize Kubernetes client
	k8sClient, err := kubernetes.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}
	log.Println("Successfully connected to Kubernetes API")

	// Parse templates
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	// Create handler
	handler := handlers.NewHandler(k8sClient, tmpl)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.ServeHome)
	mux.HandleFunc("/firefox/bookmarks.html", handler.ServeFirefoxBookmarks)
	mux.HandleFunc("/chrome/bookmarks.json", handler.ServeChromeBookmarks)
	mux.HandleFunc("/health", handler.ServeHealth)
	mux.HandleFunc("/readiness", handler.ServeReadiness)

	// Create server with timeouts
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
