package loadbal

import (
	"log"
	"net/http"
	"net/http/httputil"
	urlutil "net/url"
	"sync"
)

type ServerConfig struct {
	// url of the server
	url *urlutil.URL
	// optional weightage for weightage lb algo
	weightage uint8
	// for determining if a server is healthy or not based on health check api call
	is_healthy bool
	// health check api
	health_check_api string
	// failure threshold to trigger circuit breaker
	failure_threshold uint64
	// current failure count
	current_failure_count uint64
	// connections made to the server
	connection_count uint
	// to lock the server when in use
	Server_mutex sync.Mutex
}

func CreateServerConfig(url string, weightage uint8, health_check_api string, failure_threshold uint64) ServerConfig {
	if failure_threshold == 0 {
		failure_threshold = DEFAULT_FAILURE_THRESHOLD
	}
	parsedUrl, err := urlutil.Parse(url)
	if err != nil {
		log.Fatal("Error parsing " + url)
	}
	return ServerConfig{
		url:               parsedUrl,
		weightage:         weightage,
		is_healthy:        true,
		health_check_api:  health_check_api,
		failure_threshold: failure_threshold,
	}
}

func (sc *ServerConfig) CallHealthAPI() (bool, error) {
	_, err := http.Get(sc.health_check_api)
	if err != nil {
		sc.IncrementFailureCount()
		return false, err
	}
	return true, nil
}

func (sc *ServerConfig) IncrementConnectionCount() (bool, error) {
	sc.Server_mutex.Lock()
	defer sc.Server_mutex.Unlock()
	sc.connection_count++
	return true, nil
}

func (sc *ServerConfig) SetHealthStatus(status bool) {
	sc.Server_mutex.Lock()
	defer sc.Server_mutex.Unlock()
	sc.is_healthy = status
}

func (sc *ServerConfig) IncrementFailureCount() {
	sc.Server_mutex.Lock()
	defer sc.Server_mutex.Unlock()
	if sc.current_failure_count < sc.failure_threshold {
		sc.current_failure_count++
	} else {
		sc.is_healthy = false
	}
}

func (sc *ServerConfig) IsHealthy() bool {
	sc.Server_mutex.Lock()
	defer sc.Server_mutex.Unlock()
	return sc.is_healthy
}

func (sc *ServerConfig) GetURL() string {
	return sc.url.String()
}

func (sc *ServerConfig) GetWeight() uint8 {
	return sc.weightage
}

func (sc *ServerConfig) GetConnectionCount() uint {
	sc.Server_mutex.Lock()
	defer sc.Server_mutex.Unlock()
	return sc.connection_count
}

func (sc *ServerConfig) GetFailureCount() uint64 {
	sc.Server_mutex.Lock()
	defer sc.Server_mutex.Unlock()
	return sc.current_failure_count
}

func (sc *ServerConfig) ProxyCall() *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(sc.url)
}