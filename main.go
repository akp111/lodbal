package main

import (
	"github.com/akp111/lodbal/lodbal"
)

type DummyLoadBalAlgo struct {
}

func (lba DummyLoadBalAlgo) LbAlgo() (bool, error) {

	return true, nil
}

func main() {
	server := loadbal.CreateServerConfig("http://localhost:3000", 0, "http://localhost:3000/health", 10)
	var servers []loadbal.ServerConfig
	servers = append(servers, server)
	dummyLbAlgo := DummyLoadBalAlgo{}
	loadbal.CreateLodBal(servers, dummyLbAlgo, 100, 8000)
}
