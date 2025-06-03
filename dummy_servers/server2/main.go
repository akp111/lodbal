package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Root endpoint - this is what your load balancer forwards to
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Server2 - Received request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello from Server 2 (Port 3001)! Path: %s\n", r.URL.Path)
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Server2 - Health check - server is healthy")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Server2 OK")
	})

	// Call endpoint
	http.HandleFunc("/call", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Server2 - Call went successfully")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Call successful on Server 2")
	})

	log.Println("🚀 Backend Server 2 started on port 3001")
	log.Println("   Health endpoint: http://localhost:3001/health")
	log.Println("   Call endpoint: http://localhost:3001/call")
	log.Println("   Root endpoint: http://localhost:3001/")

	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal("Server2 crashed: ", err)
	}
}
