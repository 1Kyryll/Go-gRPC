package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/1kyryll/go-grpc/internal/middleware"
)

type SSEServer struct {
	mu      sync.Mutex
	clients map[chan []byte]struct{}
}

func NewSSEServer() *SSEServer {
	return &SSEServer{
		clients: make(map[chan []byte]struct{}),
	}
}

func (s *SSEServer) subscribe() chan []byte {
	ch := make(chan []byte, 16)
	s.mu.Lock()
	s.clients[ch] = struct{}{}
	s.mu.Unlock()
	return ch
}

func (s *SSEServer) unsubscribe(ch chan []byte) {
	s.mu.Lock()
	delete(s.clients, ch)
	s.mu.Unlock()
	close(ch)
}

// Broadcast sends data to all connected SSE clients.
func (s *SSEServer) Broadcast(data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal SSE data: %v", err)
		return
	}

	s.mu.Lock()
	clients := make([]chan []byte, 0, len(s.clients))
	for ch := range s.clients {
		clients = append(clients, ch)
	}
	s.mu.Unlock()

	for _, ch := range clients {
		select {
		case ch <- bytes:
		default:
		}
	}
}

// ServeHTTP handles GET /stream — the SSE endpoint. Kitchen staff only.
func (s *SSEServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := middleware.RequireRole(r.Context(), middleware.RoleKitchenStaff); err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ch := s.subscribe()
	defer s.unsubscribe(ch)

	fmt.Fprintf(w, "data: {\"connected\":true}\n\n")
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case data := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
