package balancer

import (
	"errors"
)

// RoundRobinLoad this is a round robin balancer
// without weight value
type RoundRobinLoad struct {
	domain      string
	activeIndex int
}

// NewRoundRobinLoad get a RoundRobin point
func NewRoundRobinLoad(domain string) Balancer {
	return &RoundRobinLoad{domain, 0}
}

// GetOne get an target by round robin
func (r *RoundRobinLoad) GetOne() (*ProxyTarget, error) {
	targetSrv, err := getTarget(r.domain)
	if err != nil {
		return nil, err
	}
	if len(targetSrv.Items) == 0 {
		return nil, errors.New("not found endpoints")
	}

	target := &ProxyTarget{targetSrv.Domain, targetSrv.Items[r.activeIndex].Endpoint}
	r.activeIndex = (r.activeIndex + 1) % len(targetSrv.Items)

	return target, nil
}

// AddAddr add an endpoint
func (r *RoundRobinLoad) AddAddr(addr string, weight uint32) error {
	endpoint := OriginItem{
		Endpoint: addr,
		Weight:   weight,
	}
	return addEndpoint(r.domain, endpoint)
}

// DelAddr delete an endpoint
func (r *RoundRobinLoad) DelAddr(addr string) error {
	return delEndpoint(r.domain, addr)
}
