package balancer

import (
	"errors"
	"math/rand"
)

// RandomLoad Load Balancers By Random
// you will get an random origin address
type RandomLoad struct{}

// NewRandomLoad get a RandomLoad point
func NewRandomLoad() Balancer {

	return &RandomLoad{}
}

// GetOne get an target by random
func (r *RandomLoad) GetOne(domain string) (*ProxyTarget, error) {
	targetSrv, err := GetTarget(domain)
	if err != nil {
		return nil, err
	}
	if len(targetSrv.Items) == 0 {
		return nil, errors.New("not found endpoints")
	}
	randCode := rand.Intn(len(targetSrv.Items))

	return &ProxyTarget{targetSrv.Domain, targetSrv.Items[randCode].Endpoint}, nil
}

// AddAddr add an endpoint
func (r *RandomLoad) AddAddr(domain string, addr string, weight uint32) error {
	endpoint := OriginItem{
		Endpoint: addr,
		Weight:   weight,
	}
	return addEndpoint(domain, endpoint)
}

// DelAddr delete an endpoint
func (r *RandomLoad) DelAddr(domain string, addr string) error {
	return delEndpoint(domain, addr)
}
