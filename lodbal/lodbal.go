// Entry point for your load balancer

package loadbal

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

const (
	Version = "0.0.1"
)

// load balancer algorithm contract
type LodBalAlgo interface {
	LbAlgo(servers []*ServerConfig) (*ServerConfig, error)
}

var lodbal LodBal

type LodBal struct {
	// array of servers that the lb needs to distribute load. Optional weigtage for weight based lb algo
	Servers []*ServerConfig
	// current server in use
	Current uint
	// lb algorithm that needs to be used
	lbAlgorithm LodBalAlgo
	// intervals in which health check api needs to be called
	HealthcheckFrequency uint8
	// port to which the lb will accept request
	Port uint
	// lock to protect shared server config
	LBmutex sync.Mutex
}

func (lb *LodBal) startServer() {
	log.Printf("Starting server at port %d", lb.Port)
	// handler to accept any type of requests
	// Thanks claude for beautiful logs
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🔄 Request received: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		server, err := lodbal.lbAlgorithm.LbAlgo(lodbal.Servers)
		if err != nil {
			log.Printf("❌ Error getting next server: %v", err)
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}

		// Log detailed server information
		log.Printf("🎯 Selected server: URL=%s, Healthy=%t, Connections=%d",
			server.GetURL(), server.IsHealthy(), server.GetConnectionCount())
		log.Printf("📊 Server details: Weight=%d, FailureCount=%d",
			server.GetWeight(), server.GetFailureCount())
		log.Printf("🔢 Current server index: %d", lb.Current)
		w.Header().Add("X-Forwarded-Server", server.GetURL())

		log.Printf("⏩ Forwarding request to: %s", server.GetURL())
		server.ProxyCall().ServeHTTP(w, r)

		log.Printf("✅ Request completed for: %s", server.GetURL())

	})
	err := http.ListenAndServe(":"+fmt.Sprintf("%d", lb.Port), nil)
	if err != nil {
		log.Fatalf("Error in starting load balancer: %v", err)
	}
}

func CreateLodBal(servers []*ServerConfig, lbAlgorithm LodBalAlgo, healthcheckFrequency uint8, port uint) {
	fmt.Println("Initializing the LodBal!!")
	lodbal = LodBal{
		Servers:              servers,
		lbAlgorithm:          lbAlgorithm,
		HealthcheckFrequency: healthcheckFrequency,
		Port:                 port,
	}
	lodbal.startServer()
}
