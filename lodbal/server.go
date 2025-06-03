package loadbal

import "net/http"

type ServerConfig struct {
	// url of the server
	url string
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
}

func CreateServerConfig(url string, weightage uint8, health_check_api string, failure_threshold uint64) ServerConfig{
	if failure_threshold == 0 {
		failure_threshold = DEFAULT_FAILURE_THRESHOLD
	}
	return ServerConfig{
		url: url,
		weightage: weightage,
		is_healthy: true,
		health_check_api: health_check_api,
		failure_threshold: failure_threshold,
	}
}

func (sc* ServerConfig) CallHealthAPI() (bool, error){
	_, err := http.Get(sc.health_check_api)
	if err != nil {
		sc.IncrementFailureCount();
		return false, err
	}
	return true, nil
}

func (sc* ServerConfig) IncrementConnectionCount() (bool, error){
	sc.connection_count++
	return true, nil
}

func (sc* ServerConfig) SetHealthStatus(status bool) {
	sc.is_healthy = status
}

func (sc* ServerConfig) IncrementFailureCount() {
	if sc.current_failure_count < sc.failure_threshold {
		sc.current_failure_count++
	} else {
		sc.SetHealthStatus(false)
	}
}
