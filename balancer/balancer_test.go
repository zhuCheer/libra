package balancer

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

type testBalancer struct{}

func (b *testBalancer) GetOne() (target *ProxyTarget, err error) {
	fmt.Println("GetOne func")
	return &ProxyTarget{"www.google.com", "192.168.1.100:8080"}, nil
}

func (b *testBalancer) AddAddr(addr string, weight uint32) error {

	return nil
}

func (b *testBalancer) DelAddr(addr string) error {

	return nil
}

func (b *testBalancer) ChangeLoadType(loadType string) error {

	return nil
}

func (b *testBalancer) GetSiteInfo(domain string) (*RegistNode, error) {

	return nil, nil
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
	err := newTarget(RegistNode{
		Domain: "www.google.com",
		Items:  []OriginItem{},
	})
	if err != nil {
		t.Error("NewTarget func has an error #1")
	}
	err = newTarget(RegistNode{
		Domain: "www.google.com",
		Items:  []OriginItem{},
	})

	if err != ErrServiceExisted {
		t.Error("NewTarget func has an error #2")
	}
}

func TestGetTarget(t *testing.T) {
	registryMap = nil

	node, err := getTarget("www.google.com")
	if node != nil || err == nil {
		t.Error("GetTarget func has an error #1")
	}

	newTarget(RegistNode{
		Domain: "www.google.com",
		Items:  []OriginItem{},
	})
	node, err = getTarget("www.google.com")

	if err != nil || node == nil {
		t.Error("GetTarget func has an error #2")
	}
}

func TestRegistTarget(t *testing.T) {

	RegistTargetNoAddr("www.google.com", "random", "http")

	RegistTargetNoAddr("www.google.com", "random", "http")

	if registryMap == nil || len(registryMap) > 1 {
		t.Error("RegistTarget func have an error")
	}
}

func TestAddEndpoint(t *testing.T) {
	registryMap = nil
	RegistTargetNoAddr("www.google.com", "random", "http")
	err := addEndpoint("www.facebook.com", OriginItem{"192.168.1.1:80", 10})

	if _, ok := registryMap["www.facebook.com"]; ok == false {
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

	if len(registryMap["www.google.com"].Items) != 5 {
		t.Error("AddEndpoint func have an error #4")
	}

}

func TestDelEndpoint(t *testing.T) {
	registryMap = nil
	RegistTargetNoAddr("www.google.com", "random", "http")
	addEndpoint("www.google.com", []OriginItem{{"192.168.1.101:80", 10}, {"192.168.1.102:80", 10}}...)
	delEndpoint("www.google.com", "192.168.1.101:80")
	delEndpoint("www.google.com", "192.168.1.101:80")
	delEndpoint("www.google.com", "192.168.1.105:80")

	if registryMap["www.google.com"].Items[0].Endpoint != "192.168.1.102:80" ||
		len(registryMap["www.google.com"].Items) != 1 {
		t.Error("DelEndpoint func have an error #1")
	}

	delEndpoint("www.google.com", "192.168.1.102:80")
	if len(registryMap["www.google.com"].Items) != 0 {
		t.Error("DelEndpoint func have an error #2")
	}

	addEndpoint("www.google.com", []OriginItem{
		{"192.168.1.101:80", 10},
		{"192.168.1.102:80", 10},
		{"192.168.1.103:80", 10},
		{"192.168.1.104:80", 10},
	}...)

	delEndpoint("www.google.com", "192.168.1.102:80")
	if len(registryMap["www.google.com"].Items) != 3 {
		t.Error("DelEndpoint func have an error #3")
	}
}

func TestFlushProxy(t *testing.T) {
	registryMap = nil
	RegistTargetNoAddr("www.google.com", "random", "http")
	RegistTargetNoAddr("www.google1.com", "random", "http")
	RegistTargetNoAddr("www.google2.com", "random", "http")

	addEndpoint("www.google.com", []OriginItem{{"192.168.1.101:80", 10}, {"192.168.1.102:80", 10}}...)
	addEndpoint("www.google2.com", []OriginItem{{"192.168.1.101:80", 10}, {"192.168.1.102:80", 10}}...)
	addEndpoint("www.google2.com", OriginItem{"192.168.1.1013:80", 10})

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
	RegistTargetNoAddr("www.google.com", "random", "http")

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

func TestChangeLoadType(t *testing.T) {
	registryMap = nil
	RegistTargetNoAddr("www.google.com", "random", "http")

	siteInfo, _ := getTarget("www.google.com")
	_, ok := (siteInfo.Balancer).(*RandomLoad)
	if ok != true {
		t.Error("changeLoadType have an error #1")
	}
	ChangeLoadType("www.google.com", "roundrobin")

	_, ok = (siteInfo.Balancer).(*RoundRobinLoad)
	if ok != true {
		t.Error("changeLoadType have an error #2")
	}
	ChangeLoadType("www.xxx.com", "xxx")

	if len(registryMap) != 2 || len(registryMap["www.google.com"].Items) != 0 {
		t.Error("changeLoadType func have an error #3")
	}
}

func TestGetSiteInfo(t *testing.T) {
	registryMap = nil

	siteInfo, _ := GetSiteInfo("www.google.com")
	if siteInfo != nil {
		t.Error("siteInfo have an error #1")
	}
	RegistTargetNoAddr("www.google.com", "random", "http")

	siteInfo, _ = GetSiteInfo("www.google.com")
	if siteInfo.Domain != "www.google.com" {
		t.Error("siteInfo have an error #2")
	}
}
