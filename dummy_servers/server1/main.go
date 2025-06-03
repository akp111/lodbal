package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Root endpoint - this is what your load balancer forwards to
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Server1 - Received request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello from Server 1 (Port 3000)! Path: %s\n", r.URL.Path)
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Server1 - Health check - server is healthy")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Server1 OK")
	})

	// Call endpoint
	http.HandleFunc("/call", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Server1 - Call went successfully")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Call successful on Server 1")
	})

	log.Println("🚀 Backend Server 1 started on port 3000")
	log.Println("   Health endpoint: http://localhost:3000/health")
	log.Println("   Call endpoint: http://localhost:3000/call")
	log.Println("   Root endpoint: http://localhost:3000/")

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("Server1 crashed: ", err)
	}
}
