package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ServerNode represents a backend instance in the load balancer pool
type ServerNode struct {
	URL          *url.URL
	Proxy        *httputil.ReverseProxy
	Alive        bool
	mu           sync.RWMutex
	ActiveConns  int64
	TotalServed  int64
	FailureCount int
}

func (s *ServerNode) SetAlive(alive bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Alive = alive
	if alive {
		s.FailureCount = 0
	} else {
		s.FailureCount++
	}
}

func (s *ServerNode) IsAlive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Alive
}

// ServerPool maintains a group of backend instances with load balancing logic
type ServerPool struct {
	Name    string
	Nodes   []*ServerNode
	current uint64
}

func (p *ServerPool) AddNode(nodeURL *url.URL, timeout time.Duration) {
	proxy := httputil.NewSingleHostReverseProxy(nodeURL)
	proxy.Transport = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 50,
		IdleConnTimeout:     90 * time.Second,
		ResponseHeaderTimeout: timeout,
	}

	node := &ServerNode{
		URL:   nodeURL,
		Proxy: proxy,
		Alive: true,
	}

	p.Nodes = append(p.Nodes, node)
}

// GetNextNode finds a live node using Least Connections algorithm (with Round Robin fallback)
func (p *ServerPool) GetNextNode() *ServerNode {
	var bestNode *ServerNode
	var minConns int64 = 1<<62 - 1

	for _, node := range p.Nodes {
		if node.IsAlive() {
			conns := atomic.LoadInt64(&node.ActiveConns)
			if conns < minConns {
				minConns = conns
				bestNode = node
			}
		}
	}

	if bestNode != nil {
		return bestNode
	}

	// Fallback to round-robin if all marked dead (graceful degradation)
	if len(p.Nodes) > 0 {
		idx := int(atomic.AddUint64(&p.current, 1) % uint64(len(p.Nodes)))
		return p.Nodes[idx]
	}

	return nil
}

// HealthCheck periodically verifies backend responsiveness
func (p *ServerPool) HealthCheck() {
	client := http.Client{Timeout: 2 * time.Second}
	for {
		for _, node := range p.Nodes {
			healthURL := fmt.Sprintf("%s/health", node.URL.String())
			resp, err := client.Get(healthURL)
			if err != nil || resp.StatusCode >= 500 {
				if node.IsAlive() {
					log.Printf("⚠️ [%s] Backend node %s is DOWN (err: %v)", p.Name, node.URL.String(), err)
				}
				node.SetAlive(false)
			} else {
				if !node.IsAlive() {
					log.Printf("✅ [%s] Backend node %s recovered to HEALTHY", p.Name, node.URL.String())
				}
				node.SetAlive(true)
			}
			if resp != nil {
				_ = resp.Body.Close()
			}
		}
		time.Sleep(5 * time.Second)
	}
}

// HeavyRoutePatterns are endpoints classified as resource-intensive
var HeavyRoutePatterns = []string{
	"/api/v1/reports",
	"/api/v1/hris/payrolls/run",
	"/api/v1/hris/payrolls/batch-approve",
	"/api/v1/inventory/opnames",
	"/api/v1/inventory/stock-movements",
	"/api/v1/finance/reconciliation/match",
}

func isHeavyRequest(r *http.Request) bool {
	path := r.URL.Path
	for _, pattern := range HeavyRoutePatterns {
		if strings.HasPrefix(path, pattern) {
			return true
		}
	}
	return false
}

func main() {
	lbPort := getEnv("LB_PORT", "8080")
	coreNodesStr := getEnv("CORE_BACKENDS", "http://127.0.0.1:8081,http://127.0.0.1:8082")
	heavyNodesStr := getEnv("HEAVY_BACKENDS", "http://127.0.0.1:8083")

	corePool := &ServerPool{Name: "CORE-REALTIME"}
	for _, raw := range strings.Split(coreNodesStr, ",") {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		u, err := url.Parse(raw)
		if err != nil {
			log.Fatalf("Invalid core backend URL '%s': %v", raw, err)
		}
		corePool.AddNode(u, 15*time.Second) // Core operations: fast 15s timeout
		log.Printf("Registered CORE node: %s", u.String())
	}

	heavyPool := &ServerPool{Name: "HEAVY-BATCH"}
	for _, raw := range strings.Split(heavyNodesStr, ",") {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		u, err := url.Parse(raw)
		if err != nil {
			log.Fatalf("Invalid heavy backend URL '%s': %v", raw, err)
		}
		heavyPool.AddNode(u, 180*time.Second) // Heavy operations: generous 3m timeout
		log.Printf("Registered HEAVY node: %s", u.String())
	}

	go corePool.HealthCheck()
	go heavyPool.HealthCheck()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Load Balancer status dashboard endpoint
		if r.URL.Path == "/lb-status" {
			serveLBStatus(w, corePool, heavyPool)
			return
		}

		var targetPool *ServerPool
		var targetType string

		if isHeavyRequest(r) {
			targetPool = heavyPool
			targetType = "HEAVY"
		} else {
			targetPool = corePool
			targetType = "CORE"
		}

		node := targetPool.GetNextNode()
		if node == nil {
			http.Error(w, fmt.Sprintf("503 Service Unavailable: No healthy %s backends available", targetType), http.StatusServiceUnavailable)
			return
		}

		// Track concurrency
		atomic.AddInt64(&node.ActiveConns, 1)
		atomic.AddInt64(&node.TotalServed, 1)
		defer atomic.AddInt64(&node.ActiveConns, -1)

		// Set forwarding headers
		r.Header.Set("X-Forwarded-Host", r.Host)
		if r.Header.Get("X-Forwarded-Proto") == "" {
			r.Header.Set("X-Forwarded-Proto", "http")
		}
		r.Header.Set("X-Load-Balancer-Route", targetType)

		node.Proxy.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:         ":" + lbPort,
		Handler:      handler,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 180 * time.Second,
	}

	log.Printf("═══════════════════════════════════════════════════════════════════════")
	log.Printf("🚀 CAFE ERP LAYER-7 PATH-BASED LOAD BALANCER ACTIVE ON PORT :%s", lbPort)
	log.Printf("🎯 Core Traffic Routing  -> %s", coreNodesStr)
	log.Printf("⚡ Heavy Traffic Routing -> %s", heavyNodesStr)
	log.Printf("📊 Live Metrics Dashboard -> http://localhost:%s/lb-status", lbPort)
	log.Printf("═══════════════════════════════════════════════════════════════════════")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Load balancer failed: %v", err)
	}
}

func serveLBStatus(w http.ResponseWriter, corePool, heavyPool *ServerPool) {
	status := map[string]interface{}{
		"system":    "Cafe ERP Smart Load Balancer",
		"timestamp": time.Now().Format(time.RFC3339),
		"pools": map[string]interface{}{
			"core":  buildPoolStats(corePool),
			"heavy": buildPoolStats(heavyPool),
		},
		"heavy_routed_endpoints": HeavyRoutePatterns,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

func buildPoolStats(p *ServerPool) []map[string]interface{} {
	stats := make([]map[string]interface{}, 0, len(p.Nodes))
	for _, n := range p.Nodes {
		stats = append(stats, map[string]interface{}{
			"url":          n.URL.String(),
			"healthy":      n.IsAlive(),
			"active_conns": atomic.LoadInt64(&n.ActiveConns),
			"total_served": atomic.LoadInt64(&n.TotalServed),
		})
	}
	return stats
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
