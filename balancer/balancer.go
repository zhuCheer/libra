package balancer

import (
	"errors"

	"log"
	"sync"
)

// Balancer is an interface used to lookup the target host
// Interfaces can be implemented using different algorithms, loop/round robin or others
type Balancer interface {
	AddAddr(addr string, weight uint32) error //add target addr
	DelAddr(addr string) error                // del target addr
	GetOne() (*ProxyTarget, error)            // Return the endpoint by different algorithms
}

// Common errors.
var (
	ErrServiceNotFound = errors.New("the proxy srv not found")
	ErrServiceExisted  = errors.New("the proxy srv has existed")
	ErrEndpointExisted = errors.New("the endpoint has existed")
)

// Global lock for the default registry,
// edit map should use lock.
var lock sync.RWMutex

// global registry proxy data
var registryMap map[string]*RegistNode

// OriginItem struct addr and weight
type OriginItem struct {
	Endpoint string // ip:port
	Weight   uint32
}

// RegistNode register a proxy node struct
type RegistNode struct {
	Domain   string
	Items    []OriginItem
	Balancer Balancer
	Scheme   string
}

// ProxyTarget proxy target node struct
type ProxyTarget struct {
	Domain string
	Addr   string
}

// newTarget New Target server is register a node
func newTarget(node RegistNode) error {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := registryMap[node.Domain]; !ok {
		if registryMap == nil {
			registryMap = map[string]*RegistNode{}
		}

		registryMap[node.Domain] = &node
		return nil
	}
	return ErrServiceExisted
}

// RegistTargetNoAddr register a target server node target ip list is empty
func RegistTargetNoAddr(domain, loadType, scheme string) *RegistNode {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := registryMap[domain]; !ok {
		if registryMap == nil {
			registryMap = map[string]*RegistNode{}
		}

		registryMap[domain] = &RegistNode{
			Domain:   domain,
			Items:    []OriginItem{},
			Balancer: getBalancerByLoadType(domain, loadType),
			Scheme:   scheme,
		}
		return registryMap[domain]
	} else {
		return registryMap[domain]
	}
}

// getBalancerByLoadType get balancer by lodeType
func getBalancerByLoadType(domain, loadType string) Balancer {
	b := NewRandomLoad(domain)
	switch loadType {
	case "random":
		b = NewRandomLoad(domain)
	case "roundrobin":
		b = NewRoundRobinLoad(domain)
	case "wroundrobin":
		b = NewWRoundRobinLoad(domain)
	}
	return b
}

// GetSiteInfo get target for out package
func GetSiteInfo(domain string) (*RegistNode, error) {
	info, err := getTarget(domain)
	if err != nil {
		log.Printf("GetSiteInfo func have error %v", err)
		return nil, err
	}

	return info, nil
}

// getTarget get a Target server
func getTarget(domain string) (*RegistNode, error) {
	lock.RLock()
	node, ok := registryMap[domain]
	lock.RUnlock()

	if ok == false {
		return nil, ErrServiceNotFound
	}

	return node, nil
}

// FlushProxy flush an proxy server
func FlushProxy(domain string) {
	lock.Lock()
	defer lock.Unlock()
	delete(registryMap, domain)
}

// addEndpoint add an endpoint
func addEndpoint(domain string, endpoints ...OriginItem) error {
	lock.Lock()
	defer lock.Unlock()

	if registryMap == nil {
		registryMap = map[string]*RegistNode{}
	}

	service, ok := registryMap[domain]
	if ok == false {
		registryMap[domain] = &RegistNode{
			domain,
			endpoints,
			getBalancerByLoadType(domain, "random"),
			"http",
		}
	} else {
		for _, item := range endpoints {
			if stringInOriginItem(item.Endpoint, service.Items) {
				return ErrEndpointExisted
			}
			service.Items = append(service.Items, item)
		}
	}

	return nil
}

// delEndpoint remove an endpoint
func delEndpoint(domain string, addr string) error {
	lock.Lock()
	defer lock.Unlock()

	service, ok := registryMap[domain]
	if ok == false {
		return ErrServiceNotFound
	}
	for k, item := range service.Items {
		if item.Endpoint == addr {
			endpoints := append(service.Items[:k], service.Items[k+1:]...)
			service.Items = endpoints
			break
		}
	}
	return nil
}

// ChangeLoadType set site load type
func ChangeLoadType(domain string, loadType string) error {
	lock.Lock()
	defer lock.Unlock()

	service, ok := registryMap[domain]
	if ok == false {
		registryMap[domain] = &RegistNode{
			domain,
			[]OriginItem{},
			getBalancerByLoadType(domain, loadType),
			"http",
		}
	} else {
		service.Balancer = getBalancerByLoadType(domain, loadType)
	}
	return nil
}

// stringInOriginItem check endpoint is existed
func stringInOriginItem(needle string, haystack []OriginItem) bool {
	result := false
	for _, item := range haystack {
		if needle == item.Endpoint {
			result = true
			break
		}
	}
	return result
}
