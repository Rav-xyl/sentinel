package plugins

import "net/http"

// Middleware defines a standard HTTP middleware signature for Sentinel plugins
type Middleware func(http.Handler) http.Handler

// Chain applies a series of middleware to a final handler
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// ExamplePlugin is a sample plugin that adds a custom header to responses
func ExamplePlugin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Sentinel-Plugin", "active")
		next.ServeHTTP(w, r)
	})
}
