// Entry point for your load balancer

package loadbal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	// "net/url"
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

		lb.forwardRequest(w, r, *server)

		log.Printf("✅ Request completed for: %s", server.GetURL())

	})
	err := http.ListenAndServe(":"+fmt.Sprintf("%d", lb.Port), nil)
	if err != nil {
		log.Fatalf("Error in starting load balancer: %v", err)
	}
}

func (lb *LodBal) forwardRequest(w http.ResponseWriter, r *http.Request, server ServerConfig) {
	// 1. construct the target url
	targetUrl := &url.URL{
		Scheme:   server.url.Scheme,
		Host:     server.url.Host,
		Path:     r.URL.Path,
		RawQuery: r.URL.RawQuery,
		Fragment: r.URL.Fragment,
	}
	if server.url.Path != "" && server.url.Path != "/" {
		targetUrl.Path = strings.TrimSuffix(server.url.Path, "/") + "/" + strings.TrimPrefix(r.URL.Path, "/")
	}

	fmt.Println("target url", targetUrl)

	// 2. New request url which has the context of call made to the proxy
	proxyRequest, err := http.NewRequestWithContext(r.Context(), r.Method, targetUrl.String(), r.Body)
	if err != nil {
		log.Fatal("Error while creating  new request")
	}
	fmt.Printf("proxy request obj %+v\n", proxyRequest)

	// 3. Copy headers
	fmt.Println("Request headers")
	for key, val := range r.Header {
		fmt.Println(key, val)

		// do not forward hopbyhop headers, ref: https://stackoverflow.com/questions/23560405/why-connection-header-should-be-removed-by-proxies
		if IsHopByHopHeader(key) {
			continue
		}

		for _, value := range val {
			proxyRequest.Header.Add(key, value)
		}
	}

	fmt.Printf("New proxy headers: \n %+v\n", proxyRequest.Header)
	// (additional) handle if user-agent is not passed
	if proxyRequest.Header.Get("User-Agent") == "" {
		proxyRequest.Header.Set("User-Agent", "")
	}
	// 4. attach proxy headers
	clientIp := GetClientIp(r)
	if clientIp != "" {
		// check if it was already forwarded from a proxy
		existing := proxyRequest.Header.Get("X-Forwarded-For")
		if existing != "" {
			proxyRequest.Header.Set("X-Forwarded-For", existing+","+clientIp)
		} else {
			proxyRequest.Header.Set("X-Forwarded-For", clientIp)
		}
	}
	proxyRequest.Header.Set("X-Forwarded-Host", r.Host)
	if r.TLS == nil {
		proxyRequest.Header.Set("X-Forwarded-Proto", "http")
	} else {
		proxyRequest.Header.Set("X-Forwarded-Proto", "https")
	}
	proxyRequest.Header.Set("X-Load-Balancer", "LodBal/"+Version)

	// 5. Make the request to the backend

	httpclient := http.Client{
		Timeout:   30 * time.Second,
		Transport: http.DefaultTransport,
	}

	resp, err := httpclient.Do(proxyRequest)
	if err != nil {
		log.Printf("❌ Backend request failed: %v", err)
		server.IncrementFailureCount()
		http.Error(w, "Backend server unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Response from the server:")
	fmt.Printf("%+v\n", resp)

	// 6. Copy response headers
	for key, values := range resp.Header {
		if IsHopByHopHeader(key) {
			continue
		}

		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.Header().Set("X-Served-By", "LodBal/"+Version)
	w.Header().Set("X-Response-Time", time.Now().Format(time.RFC3339))
	w.WriteHeader(resp.StatusCode)

	//7. Send the response back to the client
	stream_bytes, err := io.Copy(w, resp.Body)
	    if err != nil {
        log.Printf("❌ Error streaming response: %v", err)
        fmt.Printf("error streaming response: %v", err)
    }
	fmt.Printf("Streamed bytes: %d", stream_bytes)
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
