package loadbal

type RoundRobinLoadBalancerAlgo struct {
}

func (rrlba RoundRobinLoadBalancerAlgo) LbAlgo() (*ServerConfig, error) {
	lodbal.LBmutex.Lock()
	defer lodbal.LBmutex.Unlock()
	lodbal.Current += 1
	lodbal.Current = (lodbal.Current) % uint(len(lodbal.Servers))
	if !lodbal.Servers[lodbal.Current].is_healthy {
		return rrlba.LbAlgo()
	}
	_, err := lodbal.Servers[lodbal.Current].CallHealthAPI()
	if err != nil {
		return rrlba.LbAlgo()
	}
	lodbal.Servers[lodbal.Current].IncrementConnectionCount()
	return lodbal.Servers[lodbal.Current], nil
}
