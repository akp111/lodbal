package loadbal

type RoundRobinLoadBalancerAlgo struct {
}

func (rrlba RoundRobinLoadBalancerAlgo) LbAlgo(servers []*ServerConfig) (*ServerConfig, error) {
	lodbal.LBmutex.Lock()
	defer lodbal.LBmutex.Unlock()
	lodbal.Current += 1
	lodbal.Current = (lodbal.Current) % uint(len(servers))
	_, err := servers[lodbal.Current].CallHealthAPI()
	if err != nil {
		return rrlba.LbAlgo(servers)
	}
	lodbal.Servers[lodbal.Current].IncrementConnectionCount()
	return servers[lodbal.Current], nil
}
