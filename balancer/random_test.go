package balancer

import (
	"testing"
)

func TestRandomLoad(t *testing.T) {
	var balancer = NewRandomLoad()
	target, err := balancer.GetOne("name")
	if err == nil {
		t.Error("RandomLoad func have an error #1")
	}

	registryMap = nil
	NewTarget(RegistNode{
		Domain: "www.google.com",
		Items: []OriginItem{
			{"192.168.1.100", 80},
			{"192.168.1.101", 80},
			{"192.168.1.102", 80},
			{"192.168.1.103", 80},
			{"192.168.1.104", 80},
			{"192.168.1.105", 80},
			{"192.168.1.106", 80},
			{"192.168.1.107", 80},
			{"192.168.1.108", 80},
		},
	})

	var addrMap map[string]string = map[string]string{}
	for i := 0; i < 1000; i++ {
		target, err = balancer.GetOne("www.google.com")
		if err != nil {
			t.Error("RandomLoad func have an error #2")
		}
		addrMap[target.Addr] = "1"
	}

	if len(addrMap) < 5 {
		t.Error("RandomLoad func the random seed have an error #3")
	}
}

func TestAddAddrRandomLoad(t *testing.T) {
	var balancer = NewRandomLoad()
	domain := "www.google.com"
	registryMap = nil
	NewTarget(RegistNode{
		Domain: domain,
		Items: []OriginItem{
			{"192.168.1.100", 80},
		},
	})
	if len(registryMap[domain].Items) != 1 {
		t.Error("AddAddr func have an error #1")
	}

	balancer.AddAddr(domain, "192.168.1.101", 0)
	balancer.AddAddr(domain, "192.168.1.102", 0)

	if len(registryMap[domain].Items) != 3 {
		t.Error("AddAddr func have an error #2")
	}

}

func TestDelAddrRandomLoad(t *testing.T) {
	var balancer = NewRandomLoad()
	domain := "www.google.com"
	registryMap = nil
	NewTarget(RegistNode{
		Domain: domain,
		Items: []OriginItem{
			{"192.168.1.100", 80},
			{"192.168.1.101", 80},
			{"192.168.1.102", 80},
		},
	})

	balancer.DelAddr(domain, "192.168.1.101")

	if len(registryMap[domain].Items) != 2 {
		t.Error("DelAddr func have an error #1")
	}
}
