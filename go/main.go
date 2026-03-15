package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	mrand "math/rand"
	"net/http"
	"time"
)

const (
	streamDelay  = 1 * time.Second
	retryTimeout = 15000 // milliseconds
	maxEvents    = 60
	port         = 8000
)

// newUUID generates a version-4 (random) UUID using crypto/rand.
func newUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant bits
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// pollHandler streams SSE events with random numbers, one per second.
// It mirrors the Python FastAPI /poll endpoint exactly:
//   - same event name ("random_number")
//   - same retry value (15000 ms)
//   - same data shape ({"value": int, "timestamp": ISO8601})
//   - same max-event cap (60 events)
//   - immediate disconnect detection via context cancellation
func pollHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()

	for count := 0; count < maxEvents; count++ {
		// Non-blocking disconnect check before generating an event.
		select {
		case <-ctx.Done():
			log.Println("Client disconnected")
			return
		default:
		}

		value := mrand.Intn(100) + 1
		// Python's datetime.isoformat(sep="T", timespec="auto") gives microsecond
		// precision; we match that with a 6-decimal-place layout.
		timestamp := time.Now().Format("2006-01-02T15:04:05.000000")

		data, _ := json.Marshal(map[string]any{
			"value":     value,
			"timestamp": timestamp,
		})

		fmt.Fprintf(w, "id: %s\n", newUUID())
		fmt.Fprintf(w, "event: random_number\n")
		fmt.Fprintf(w, "retry: %d\n", retryTimeout)
		fmt.Fprintf(w, "data: %s\n", data)
		fmt.Fprintf(w, "\n")
		flusher.Flush()

		log.Printf("Generated random number: %d\n", value)

		// Sleep for streamDelay, but wake immediately on client disconnect.
		select {
		case <-ctx.Done():
			log.Println("Client disconnected")
			return
		case <-time.After(streamDelay):
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "SSE Server running. Visit /poll for the event stream.",
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/poll", pollHandler)
	mux.HandleFunc("/", rootHandler)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting SSE server on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
