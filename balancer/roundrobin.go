package balancer

import (
	"errors"
)

// this is a round robin balancer
// without weight value
type RoundRobinLoad struct {
	activeIndex int
}

// get a RoundRobin point
func NewRoundRobinLoad() Balancer {
	return &RoundRobinLoad{0}
}

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
