package handlers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/hauke-cloud/bookmark-generator/pkg/bookmarks"
	"github.com/hauke-cloud/bookmark-generator/pkg/kubernetes"
)

// Handler handles HTTP requests
type Handler struct {
	k8sClient *kubernetes.Client
	generator *bookmarks.Generator
	tmpl      *template.Template
}

// NewHandler creates a new HTTP handler
func NewHandler(k8sClient *kubernetes.Client, tmpl *template.Template) *Handler {
	return &Handler{
		k8sClient: k8sClient,
		generator: bookmarks.NewGenerator(),
		tmpl:      tmpl,
	}
}

// ServeHome renders the home page
func (h *Handler) ServeHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	ingresses, err := h.k8sClient.GetIngresses(ctx)
	if err != nil {
		log.Printf("Error getting ingresses: %v", err)
		http.Error(w, "Failed to retrieve ingresses", http.StatusInternalServerError)
		return
	}

	data := struct {
		Ingresses []kubernetes.IngressInfo
		Count     int
	}{
		Ingresses: ingresses,
		Count:     len(ingresses),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}
}

// ServeFirefoxBookmarks generates and serves Firefox bookmarks
func (h *Handler) ServeFirefoxBookmarks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	ingressesByNamespace, err := h.k8sClient.GetIngressesByNamespace(ctx)
	if err != nil {
		log.Printf("Error getting ingresses: %v", err)
		http.Error(w, "Failed to retrieve ingresses", http.StatusInternalServerError)
		return
	}

	bookmarkData := h.generator.GenerateFirefoxGrouped(ingressesByNamespace)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=bookmarks.html")
	w.Write(bookmarkData)
}

// ServeChromeBookmarks generates and serves Chrome bookmarks
func (h *Handler) ServeChromeBookmarks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	ingressesByNamespace, err := h.k8sClient.GetIngressesByNamespace(ctx)
	if err != nil {
		log.Printf("Error getting ingresses: %v", err)
		http.Error(w, "Failed to retrieve ingresses", http.StatusInternalServerError)
		return
	}

	bookmarkData, err := h.generator.GenerateChromeGrouped(ingressesByNamespace)
	if err != nil {
		log.Printf("Error generating Chrome bookmarks: %v", err)
		http.Error(w, "Failed to generate bookmarks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=bookmarks.json")
	w.Write(bookmarkData)
}

// ServeHealth provides a health check endpoint
func (h *Handler) ServeHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Try to get ingresses to verify k8s connectivity
	_, err := h.k8sClient.GetIngresses(ctx)
	if err != nil {
		log.Printf("Health check failed: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Unhealthy: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

// ServeReadiness provides a readiness check endpoint
func (h *Handler) ServeReadiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Ready")
}
