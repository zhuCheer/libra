package balancer

import (
	"strconv"
	"sync"
	"testing"
)

type testBalancer struct{}

func (b *testBalancer) GetOne(host string) (target *ProxyTarget, err error) {

	return &ProxyTarget{"www.google.com", "192.168.1.100:8080"}, nil
}

func TestBalancerInterface(t *testing.T) {
	var _ Balancer = (*testBalancer)(nil)

	var testB interface{} = &testBalancer{}
	_, ok := testB.(Balancer)

	if ok == false {
		t.Error("BalancerInterface implemention have an error #1")
	}
}

func TestNewTarget(t *testing.T) {
	registryMap = nil
	err := NewTarget(RegistNode{
		Domain: "www.google.com",
		Items:  []OriginItem{},
	})
	if err != nil {
		t.Error("NewTarget func has an error #1")
	}
	err = NewTarget(RegistNode{
		Domain: "www.google.com",
		Items:  []OriginItem{},
	})

	if err != ErrServiceExisted {
		t.Error("NewTarget func has an error #2")
	}
}

func TestGetTarget(t *testing.T) {
	registryMap = nil

	node, err := GetTarget("www.google.com")
	if node != nil || err == nil {
		t.Error("GetTarget func has an error #1")
	}

	NewTarget(RegistNode{
		Domain: "www.google.com",
		Items:  []OriginItem{},
	})
	node, err = GetTarget("www.google.com")

	if err != nil || node == nil {
		t.Error("GetTarget func has an error #2")
	}
}

func TestRegistTarget(t *testing.T) {

	RegistTarget("www.google.com", []string{})

	RegistTarget("www.google.com", []string{})

	if registryMap == nil || len(registryMap) > 1 {
		t.Error("RegistTarget func have an error")
	}
}

func TestAddEndpoint(t *testing.T) {
	registryMap = nil
	RegistTarget("www.google.com", []string{})
	err := addEndpoint("www.facebook.com", OriginItem{"192.168.1.1:80", 10})
	if err == nil {
		t.Error("AddEndpoint func have an error #1")
	}
	err = addEndpoint("www.google.com", OriginItem{"192.168.1.100:80", 10})
	err = addEndpoint("www.google.com", OriginItem{"192.168.1.100:80", 10})
	if err == nil || err.Error() != "the endpoint has existed" {
		t.Error("AddEndpoint func have an error #2")
	}

	addEndpoint("www.google.com", []OriginItem{{"192.168.1.101:80", 10}, {"192.168.1.102:80", 10}}...)
	if len(registryMap["www.google.com"].Items) != 3 {
		t.Error("AddEndpoint func have an error #3")
	}

	addEndpoint("www.google.com", []OriginItem{
		{"192.168.1.101:8080", 10},
		{"192.168.1.102:8080", 10},
		{"192.168.1.102:8080", 10},
	}...)

	if len(registryMap["www.google.com"].Items) != 3 {
		t.Error("AddEndpoint func have an error #4")
	}

}

func TestAddAddrWithoutWeight(t *testing.T) {
	registryMap = nil
	domain := "www.google.com"
	RegistTarget(domain, []string{})

	err := AddAddrWithoutWeight("www.facebook.com", "192.168.1.1:80")
	if err == nil {
		t.Error("AddAddrWithoutWeight func have an error #1")
	}

	AddAddrWithoutWeight(domain, "192.168.137.100")
	if len(registryMap[domain].Items) != 1 {
		t.Error("AddAddrWithoutWeight func have an error #2")
	}

	err = AddAddrWithoutWeight(domain, "192.168.137.100")
	if err == nil || err.Error() != "the endpoint has existed" {
		t.Error("AddAddrWithoutWeight func have an error #3")
	}

	AddAddrWithoutWeight(domain, "192.168.137.101", "192.168.137.102")
	if len(registryMap[domain].Items) != 3 {
		t.Error("AddAddrWithoutWeight func have an error #41")
	}
}

func TestAddAddrWithWeight(t *testing.T) {
	registryMap = nil
	domain := "www.google.com"
	RegistTarget(domain, []string{})
	err := AddAddrWithWeight(domain, "192.168.1.100:80", 80)
	if err != nil {
		t.Error("AddAddrWithWeight func have an error #1")
	}

	if registryMap[domain].Items[0].Weight != 80 {
		t.Error("AddAddrWithWeight func have an error #2")
	}
}

func TestDelEndpoint(t *testing.T) {
	registryMap = nil
	RegistTarget("www.google.com", []string{})
	addEndpoint("www.google.com", []OriginItem{{"192.168.1.101:80", 10}, {"192.168.1.102:80", 10}}...)
	DelEndpoint("www.google.com", "192.168.1.101:80")
	DelEndpoint("www.google.com", "192.168.1.101:80")
	DelEndpoint("www.google.com", "192.168.1.105:80")
	if registryMap["www.google.com"].Items[0].Endpoint != "192.168.1.102:80" ||
		len(registryMap["www.google.com"].Items) != 1 {
		t.Error("DelEndpoint func have an error #1")
	}

	DelEndpoint("www.google.com", "192.168.1.102:80")
	if len(registryMap["www.google.com"].Items) != 0 {
		t.Error("DelEndpoint func have an error #2")
	}

	addEndpoint("www.google.com", []OriginItem{
		{"192.168.1.101:80", 10},
		{"192.168.1.102:80", 10},
		{"192.168.1.103:80", 10},
		{"192.168.1.104:80", 10},
	}...)

	DelEndpoint("www.google.com", "192.168.1.102:80")
	if len(registryMap["www.google.com"].Items) != 3 {
		t.Error("DelEndpoint func have an error #3")
	}
}

func TestFlushProxy(t *testing.T) {
	registryMap = nil
	RegistTarget("www.google.com", []string{})
	RegistTarget("www.google1.com", []string{})
	RegistTarget("www.google2.com", []string{})

	addEndpoint("www.google.com", []OriginItem{{"192.168.1.101:80", 10}, {"192.168.1.102:80", 10}}...)
	addEndpoint("www.google2.com", []OriginItem{{"192.168.1.101:80", 10}, {"192.168.1.102:80", 10}}...)

	FlushProxy("www.google4.com")
	if len(registryMap) != 3 {
		t.Error("FlushProxy func have an error #1")
	}

	FlushProxy("www.google1.com")
	if len(registryMap) != 2 || len(registryMap["www.google.com"].Items) != 2 {
		t.Error("FlushProxy func have an error #2")
	}

	FlushProxy("www.google2.com")
	if len(registryMap) != 1 || len(registryMap["www.google.com"].Items) != 2 {
		t.Error("FlushProxy func have an error #3")
	}
}

func TestRegistryMapLog(t *testing.T) {
	registryMap = nil
	RegistTarget("www.google.com", []string{})

	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			endpoint := "192.168.1.1" + strconv.Itoa(i) + ":8080"
			addEndpoint("www.google.com", OriginItem{endpoint, 10})
			wg.Done()
		}(i)
	}
	wg.Wait()

	if len(registryMap["www.google.com"].Items) != 100 {
		t.Error("RegistryMap Concurrency have an error #1")
	}
}

func TestStringInOriginItem(t *testing.T) {

	testItem := []OriginItem{
		{"192.168.137.100:80", 1},
		{"192.168.137.101:80", 1},
		{"192.168.137.102:80", 1},
	}

	t1 := stringInOriginItem("192.168.137.100", testItem)
	t2 := stringInOriginItem("192.168.137.101:80", testItem)
	t3 := stringInOriginItem("192.168.137.102:8080", testItem)
	t4 := stringInOriginItem("102:8080", testItem)
	t5 := stringInOriginItem("192.168.137.102:8080 ", testItem)
	t6 := stringInOriginItem(" 192.168.137.102:8080 ", testItem)
	t7 := stringInOriginItem(" 192.168.137. 102:8080 ", testItem)

	if t1 != false {
		t.Error("stringInOriginItem have an error #1")
	}

	if t2 != true {
		t.Error("stringInOriginItem have an error #2")
	}
	if t3 != false {
		t.Error("stringInOriginItem have an error #3")
	}
	if t4 != false {
		t.Error("stringInOriginItem have an error #4")
	}
	if t5 != false {
		t.Error("stringInOriginItem have an error #5")
	}
	if t6 != false {
		t.Error("stringInOriginItem have an error #6")
	}
	if t7 != false {
		t.Error("stringInOriginItem have an error #7")
	}

}
