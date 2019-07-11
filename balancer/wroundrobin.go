package balancer

import (
	"errors"
)

// WRoundRobinLoad this is a round robin by weight balancer
// weighted round robin struct
type WRoundRobinLoad struct {
	domain       string
	activeIndex  int
	activeWeight uint32
	activeItems  []OriginItem
}

// NewWRoundRobinLoad get a WRoundRobin point
func NewWRoundRobinLoad(domain string) Balancer {
	return &WRoundRobinLoad{domain, 0, 0, []OriginItem{}}
}

// GetOne get an target by round robin with weight
func (r *WRoundRobinLoad) GetOne() (*ProxyTarget, error) {
	targetSrv, err := getTarget(r.domain)
	if err != nil {
		return nil, err
	}
	if len(targetSrv.Items) == 0 {
		return nil, errors.New("not found endpoints")
	}

	initAllZero := 1
	for _, item := range targetSrv.Items {
		if item.Weight > 0 {
			initAllZero = 0
		}
	}
	if initAllZero == 1 {
		return nil, errors.New("not found available endpoints")
	}

	isAllZero := 1
	for _, item := range r.activeItems {
		if item.Weight > 0 {
			isAllZero = 0
		}
	}
	if isAllZero == 1 {
		r.reloadActiveItems()
	}

	var target *ProxyTarget

	for i := 0; i < len(r.activeItems); i++ {
		if r.activeItems[r.activeIndex].Weight > 0 {
			r.activeItems[r.activeIndex].Weight--
			target = &ProxyTarget{targetSrv.Domain, r.activeItems[r.activeIndex].Endpoint}
			break
		} else {
			r.activeIndex = (r.activeIndex + 1) % len(targetSrv.Items)
		}
	}
	r.activeIndex = (r.activeIndex + 1) % len(targetSrv.Items)

	return target, nil
}

// AddAddr add an endpoint
func (r *WRoundRobinLoad) AddAddr(addr string, weight uint32) error {
	endpoint := OriginItem{
		Endpoint: addr,
		Weight:   weight,
	}
	err := addEndpoint(r.domain, endpoint)
	if err != nil {
		return err
	}

	err = r.reloadActiveItems()
	if err != nil {
		return err
	}

	return err
}

// DelAddr delete an endpoint
func (r *WRoundRobinLoad) DelAddr(addr string) error {
	err := delEndpoint(r.domain, addr)
	if err != nil {
		return err
	}

	err = r.reloadActiveItems()
	if err != nil {
		return err
	}
	return err
}

// reloadActiveItems reload active weight
// when edit items,should run this func
func (r *WRoundRobinLoad) reloadActiveItems() error {
	target, err := getTarget(r.domain)
	if err != nil {
		return err
	}

	originItems := make([]OriginItem, 0)
	for _, item := range target.Items {
		originItems = append(originItems, item)
	}

	gcdWeight, err := getGCDWeight(originItems)
	if err != nil {
		return err
	}

	for k, item := range target.Items {
		if originItems[k].Weight == 0 {
			originItems[k].Weight = 0
		} else {
			originItems[k].Weight = item.Weight / gcdWeight
		}
	}

	r.activeItems = originItems
	r.activeIndex = 0
	return nil
}

// getMaxWeight get max weight origin item
func getMaxWeight(items []OriginItem) (uint32, error) {
	activeItem := OriginItem{}
	if len(items) == 0 {
		return 0, errors.New("origin items is empty")
	}

	activeItem = items[0]
	for _, v := range items {
		if v.Weight > activeItem.Weight {
			activeItem = v
		}
	}
	return activeItem.Weight, nil
}

// getGCDWeight Greatest Common Divisor By weight item
func getGCDWeight(items []OriginItem) (uint32, error) {
	if len(items) == 0 {
		return 0, errors.New("origin items is empty")
	}

	var gcdWeight uint32
	for i := 1; i < len(items); i++ {
		if i == 1 {
			gcdWeight = gcd(items[0].Weight, items[i].Weight)
		} else {
			gcdWeight = gcd(gcdWeight, items[i].Weight)
		}
	}

	return gcdWeight, nil
}

// gcd Greatest Common Divisor  of two Numbers
// this is called Euclidean algorithm
func gcd(m, n uint32) uint32 {
	if m < n {
		m, n = n, m
	}

	for {
		if n == 0 {
			break
		}
		m, n = n, m%n
	}
	return m
}
