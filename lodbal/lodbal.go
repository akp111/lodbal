// Entry point for your load balancer

package loadbal

import (
	"fmt"
	"log"
	"net/http"
)

const (
	Version = "0.0.1"
)

// load balancer algorithm contract
type LodBalAlgo interface{

	LbAlgo() (bool, error)
}

type LodBal struct {
	// array of servers that the lb needs to distribute load. Optional weigtage for weight based lb algo
	servers []ServerConfig
	// lb algorithm that needs to be used
	lbAlgorithm LodBalAlgo
	// intervals in which health check api needs to be called
	healthcheckFrequency uint8
	// port to which the lb will accept request
	port uint
}

func (lb* LodBal) startServer(){
	log.Printf("Starting server at port %d", lb.port);
	// handler to accept any type of requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("It has hit the rout successfully")
		log.Println(r.Method, r.Host, r.URL , r.Body)
		// get the server to which the request needs to forwarded
	})
	err := http.ListenAndServe(":"+ fmt.Sprintf("%d", lb.port), nil)
	if err != nil {
		log.Fatalf("Error in starting load balancer: %v", err)
	}
}

func CreateLodBal(servers []ServerConfig, lbAlgorithm LodBalAlgo, healthcheckFrequency uint8, port uint) {
	fmt.Println("Initializing the LodBal!!")
	lodbal := LodBal {
		servers: servers,
		lbAlgorithm: lbAlgorithm,
		healthcheckFrequency: healthcheckFrequency,
		port: port,
	}
	lodbal.startServer()
}