package balancer

import (
	"errors"
)

// this is a round robin by weight balancer
// weighted round robin struct
type WRoundRobinLoad struct {
	activeIndex  int
	activeWeight uint32
	activeItems  []OriginItem
}

// get a WRoundRobin point
func NewWRoundRobinLoad() Balancer {
	return &WRoundRobinLoad{0, 0, []OriginItem{}}
}

func (r *WRoundRobinLoad) GetOne(domain string) (*ProxyTarget, error) {
	targetSrv, err := GetTarget(domain)
	if err != nil {
		return nil, err
	}
	if len(targetSrv.Items) == 0 {
		return nil, errors.New("not found endpoints")
	}

	isAllZero := 1
	for _, item := range r.activeItems {
		if item.Weight > 0 {
			isAllZero = 0
		}
	}
	if isAllZero == 1 {
		r.reloadActiveItems(domain)
	}

	var target *ProxyTarget
	for {
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

func (r *WRoundRobinLoad) AddAddr(domain string, addr string, weight uint32) error {
	endpoint := OriginItem{
		Endpoint: addr,
		Weight:   weight,
	}
	err := addEndpoint(domain, endpoint)
	if err != nil {
		return err
	}

	err = r.reloadActiveItems(domain)
	if err != nil {
		return err
	}

	return err
}

func (r *WRoundRobinLoad) DelAddr(domain string, addr string) error {
	err := delEndpoint(domain, addr)
	if err != nil {
		return err
	}

	err = r.reloadActiveItems(domain)
	if err != nil {
		return err
	}
	return err
}

// reload active weight
// when edit items,should run this func
func (r *WRoundRobinLoad) reloadActiveItems(domain string) error {
	target, err := GetTarget(domain)
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
		originItems[k].Weight = item.Weight / gcdWeight
	}

	r.activeItems = originItems
	r.activeIndex = 0
	return nil
}

// get max weight origin item
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

// Greatest Common Divisor By weight item
func getGCDWeight(items []OriginItem) (uint32, error) {
	if len(items) == 0 {
		return 0, errors.New("origin items is empty")
	}

	var gcdWeight uint32 = 0
	for i := 1; i < len(items); i++ {
		if i == 1 {
			gcdWeight = gcd(items[0].Weight, items[i].Weight)
		} else {
			gcdWeight = gcd(gcdWeight, items[i].Weight)
		}
	}

	return gcdWeight, nil
}

// Greatest Common Divisor  of two Numbers
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
