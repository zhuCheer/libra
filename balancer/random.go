package balancer

import (
	"errors"
	"math/rand"
)

// Load Balancers By Random
// you will get an random origin address
type RandomLoad struct{}

// get a RandomLoad point
func NewRandomLoad() Balancer {

	return &RandomLoad{}
}

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

func (r *RandomLoad) AddAddr(domain string, addr string, weight uint32) error {
	endpoint := OriginItem{
		Endpoint: addr,
		Weight:   weight,
	}
	return addEndpoint(domain, endpoint)
}

func (r *RandomLoad) DelAddr(domain string, addr string) error {
	return delEndpoint(domain, addr)
}
