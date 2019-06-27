package balancer

import (
	"errors"
	"sync"
)

// Balancer is an interface used to lookup the target host
// Interfaces can be implemented using different algorithms, loop/round robin or others
type Balancer interface {
	GetOne(domain string) (*ProxyTarget, error) // Return the endpoint by different algorithms
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
var registryMap map[string]RegistNode

// origin list item
type OriginItem struct {
	Endpoint string // ip:port
	Weight   uint32
}

// register a proxy node
type RegistNode struct {
	//Name         string
	Domain string
	Items  []OriginItem
}

// proxy target node
type ProxyTarget struct {
	Domain string
	Addr   string
}

// New Target server is register a node
func NewTarget(node RegistNode) error {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := registryMap[node.Domain]; !ok {
		if registryMap == nil {
			registryMap = map[string]RegistNode{}
		}

		registryMap[node.Domain] = node
		return nil
	} else {
		return ErrServiceExisted
	}
}

// register a target server node target ip list is empty
func RegistTargetNoAddr(domain string) {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := registryMap[domain]; !ok {
		if registryMap == nil {
			registryMap = map[string]RegistNode{}
		}

		registryMap[domain] = RegistNode{
			Domain: domain,
			Items:  []OriginItem{},
		}
	}
}

// get a Target server
func GetTarget(domain string) (*RegistNode, error) {
	lock.RLock()
	node, ok := registryMap[domain]
	lock.RUnlock()

	if ok == false {
		return nil, ErrServiceNotFound
	}

	return &node, nil
}

// AddEndpointWithoutWeight is a add endpoint function without weight
// actually the weight used of default value one
func AddAddrWithoutWeight(domain string, addrs ...string) error {

	endpoints := make([]OriginItem, 0)
	for _, item := range addrs {

		endpoints = append(endpoints, OriginItem{
			Endpoint: item,
			Weight:   1,
		})
	}
	return addEndpoint(domain, endpoints...)
}

// AddAddrWithWeight is a add endpoint function with weight
// every addr should with a weight used for balanced choice algorithm.
func AddAddrWithWeight(domain string, addr string, weight uint32) error {
	endpoint := OriginItem{
		Endpoint: addr,
		Weight:   weight,
	}
	return addEndpoint(domain, endpoint)
}

// add an endpoint
func addEndpoint(domain string, endpoints ...OriginItem) error {
	lock.Lock()
	defer lock.Unlock()

	service, ok := registryMap[domain]
	if ok == false {
		return ErrServiceNotFound
	}
	for _, item := range endpoints {
		if stringInOriginItem(item.Endpoint, service.Items) {
			return ErrEndpointExisted
		}
		service.Items = append(service.Items, item)
	}

	registryMap[domain] = service
	return nil
}

// remove an endpoint
func DelEndpoint(domain string, addr string) error {
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
	registryMap[domain] = service
	return nil
}

// flush an proxy server
func FlushProxy(domain string) {
	lock.Lock()
	defer lock.Unlock()
	delete(registryMap, domain)
}

// check endpoint is existed
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
