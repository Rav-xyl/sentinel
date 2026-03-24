package proxy

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Rav-xyl/sentinel/store"
	"golang.org/x/time/rate"
)

// IPLimiter manages rate limiters for individual IP addresses
type IPLimiter struct {
	limiters sync.Map
	rate     rate.Limit
	burst    int
}

// NewIPLimiter creates a new per-IP rate limiter
func NewIPLimiter(r rate.Limit, b int) *IPLimiter {
	return &IPLimiter{
		rate:  r,
		burst: b,
	}
}

// GetLimiter returns the rate.Limiter for the provided IP address
func (i *IPLimiter) GetLimiter(ip string) *rate.Limiter {
	limiter, _ := i.limiters.LoadOrStore(ip, rate.NewLimiter(i.rate, i.burst))
	return limiter.(*rate.Limiter)
}

// Router handles incoming HTTP requests and forwards them to the appropriate target based on the database rules.
type Router struct {
	db      *store.DB
	limiter *IPLimiter
}

// NewRouter creates a new router backed by the given database store and a global per-IP rate limiter.
func NewRouter(db *store.DB) *Router {
	return &Router{
		db:      db,
		limiter: NewIPLimiter(rate.Limit(5), 10), // Default: 5 req/sec with a burst of 10
	}
}

// ServeHTTP implements the http.Handler interface for the Router.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 1. Rate Limiting
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		ip = req.RemoteAddr // Fallback
	}

	limiter := r.limiter.GetLimiter(ip)
	if !limiter.Allow() {
		log.Printf("[Router] 429 - Rate limit exceeded for IP: %s\n", ip)
		http.Error(w, "Sentinel: Too Many Requests", http.StatusTooManyRequests)
		return
	}

	// 2. Domain Routing
	host := strings.ToLower(strings.Split(req.Host, ":")[0]) // Drop port if exists

	targetStr, err := r.db.GetTarget(host)
	if err != nil {
		log.Printf("[Router] Error retrieving route for %s: %v\n", host, err)
		http.Error(w, "Sentinel: Internal Server Error", http.StatusInternalServerError)
		return
	}

	if targetStr == "" {
		log.Printf("[Router] 404 - Unknown Host: %s\n", host)
		http.Error(w, "Sentinel: Host Not Found", http.StatusNotFound)
		return
	}

	targetURL, err := url.Parse(targetStr)
	if err != nil {
		log.Printf("[Router] Invalid target URL in DB for %s: %v\n", host, err)
		http.Error(w, "Sentinel: Invalid Configuration", http.StatusInternalServerError)
		return
	}

	log.Printf("[Router] Proxying %s -> %s\n", host, targetURL.String())

	// 3. Dynamic Proxying
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Optimize the transport layer for aggressive connection pooling
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	
	proxy.Director = func(outReq *http.Request) {
		outReq.URL.Scheme = targetURL.Scheme
		outReq.URL.Host = targetURL.Host
		outReq.URL.Path = targetURL.Path + outReq.URL.Path
		outReq.Header.Set("X-Forwarded-Host", req.Host)
		outReq.Header.Set("X-Forwarded-For", ip)
		outReq.Header.Set("X-Origin-Host", targetURL.Host)

		// Day 18: WebSocket and Edge Case Hardening
		if strings.ToLower(req.Header.Get("Connection")) == "upgrade" && strings.ToLower(req.Header.Get("Upgrade")) == "websocket" {
			outReq.Header.Set("Connection", "Upgrade")
			outReq.Header.Set("Upgrade", "websocket")
		}
	}

	proxy.ServeHTTP(w, req)
}
