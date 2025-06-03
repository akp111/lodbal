package main

import (
	"github.com/akp111/lodbal/lodbal"
)

type DummyLoadBalAlgo struct {
}

func (lba DummyLoadBalAlgo) LbAlgo(servers []loadbal.ServerConfig) (loadbal.ServerConfig, error) {
	return servers[0], nil
}

func main() {
	server1 := loadbal.CreateServerConfig("http://localhost:3000", 0, "http://localhost:3000/health", 10)
	server2 := loadbal.CreateServerConfig("http://localhost:3001", 0, "http://localhost:3001/health", 10)
	var servers []*loadbal.ServerConfig;
	servers = append(servers, &server1)
	servers = append(servers, &server2)
	// Initialize server map - safe in single-threaded main function
	// dummyLbAlgo := DummyLoadBalAlgo{}
	loadbal.CreateLodBal(servers, loadbal.RoundRobinLoadBalancerAlgo{}, 100, 8000)
}
