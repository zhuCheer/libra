package balancer

import (
	"errors"
)

// RoundRobinLoad this is a round robin balancer
// without weight value
type RoundRobinLoad struct {
	activeIndex int
}

// NewRoundRobinLoad get a RoundRobin point
func NewRoundRobinLoad() Balancer {
	return &RoundRobinLoad{0}
}

// GetOne get an target by round robin
func (r *RoundRobinLoad) GetOne(domain string) (*ProxyTarget, error) {
	targetSrv, err := GetTarget(domain)
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
func (r *RoundRobinLoad) AddAddr(domain string, addr string, weight uint32) error {
	endpoint := OriginItem{
		Endpoint: addr,
		Weight:   weight,
	}
	return addEndpoint(domain, endpoint)
}

// DelAddr delete an endpoint
func (r *RoundRobinLoad) DelAddr(domain string, addr string) error {
	return delEndpoint(domain, addr)
}
