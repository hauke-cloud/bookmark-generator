package kubernetes

import (
	"context"
	"fmt"
	"sort"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Client wraps the Kubernetes client
type Client struct {
	clientset *kubernetes.Clientset
}

// IngressInfo contains information about an ingress
type IngressInfo struct {
	Name      string
	Namespace string
	Host      string
	Path      string
	URL       string
}

// NewClient creates a new Kubernetes client using in-cluster config
func NewClient() (*Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &Client{clientset: clientset}, nil
}

// GetIngresses retrieves all ingresses from all namespaces
func (c *Client) GetIngresses(ctx context.Context) ([]IngressInfo, error) {
	ingresses, err := c.clientset.NetworkingV1().Ingresses("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list ingresses: %w", err)
	}

	var result []IngressInfo

	for _, ingress := range ingresses.Items {
		result = append(result, c.extractIngressInfo(&ingress)...)
	}

	// Sort by URL for consistent output
	sort.Slice(result, func(i, j int) bool {
		return result[i].URL < result[j].URL
	})

	return result, nil
}

// GetIngressesByNamespace retrieves ingresses grouped by namespace, with unique hosts only
func (c *Client) GetIngressesByNamespace(ctx context.Context) (map[string][]IngressInfo, error) {
	ingresses, err := c.clientset.NetworkingV1().Ingresses("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list ingresses: %w", err)
	}

	grouped := make(map[string]map[string]IngressInfo) // namespace -> host -> IngressInfo

	for _, ingress := range ingresses.Items {
		infos := c.extractIngressInfo(&ingress)
		
		if _, exists := grouped[ingress.Namespace]; !exists {
			grouped[ingress.Namespace] = make(map[string]IngressInfo)
		}
		
		// Keep only one entry per unique host per namespace
		for _, info := range infos {
			if _, exists := grouped[ingress.Namespace][info.Host]; !exists {
				grouped[ingress.Namespace][info.Host] = info
			}
		}
	}

	// Convert to sorted slice per namespace
	result := make(map[string][]IngressInfo)
	for namespace, hosts := range grouped {
		var ingressList []IngressInfo
		for _, info := range hosts {
			ingressList = append(ingressList, info)
		}
		
		// Sort by host for consistent output
		sort.Slice(ingressList, func(i, j int) bool {
			return ingressList[i].Host < ingressList[j].Host
		})
		
		result[namespace] = ingressList
	}

	return result, nil
}

func (c *Client) extractIngressInfo(ingress *networkingv1.Ingress) []IngressInfo {
	var result []IngressInfo

	for _, rule := range ingress.Spec.Rules {
		host := rule.Host
		if host == "" {
			continue
		}

		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				pathStr := path.Path
				if pathStr == "" {
					pathStr = "/"
				}

				// Determine scheme (https if TLS is configured, http otherwise)
				scheme := "http"
				if len(ingress.Spec.TLS) > 0 {
					for _, tls := range ingress.Spec.TLS {
						for _, tlsHost := range tls.Hosts {
							if tlsHost == host {
								scheme = "https"
								break
							}
						}
						if scheme == "https" {
							break
						}
					}
				}

				url := fmt.Sprintf("%s://%s%s", scheme, host, pathStr)

				result = append(result, IngressInfo{
					Name:      ingress.Name,
					Namespace: ingress.Namespace,
					Host:      host,
					Path:      pathStr,
					URL:       url,
				})
			}
		} else {
			// No HTTP paths, just add the host
			scheme := "http"
			if len(ingress.Spec.TLS) > 0 {
				for _, tls := range ingress.Spec.TLS {
					for _, tlsHost := range tls.Hosts {
						if tlsHost == host {
							scheme = "https"
							break
						}
					}
					if scheme == "https" {
						break
					}
				}
			}

			url := fmt.Sprintf("%s://%s", scheme, host)

			result = append(result, IngressInfo{
				Name:      ingress.Name,
				Namespace: ingress.Namespace,
				Host:      host,
				Path:      "/",
				URL:       url,
			})
		}
	}

	return result
}
