package balancer

import (
	"testing"
)

func TestWRoudRobinLoadZeroWeight(t *testing.T) {
	domain := "www.facebook.com"
	var balancer = NewWRoundRobinLoad(domain)

	registryMap = nil
	newTarget(RegistNode{
		Domain: domain,
		Items: []OriginItem{
			{"192.168.1.100", 0},
			{"192.168.1.101", 0},
		},
	})
	balancer.AddAddr("192.168.1.102", 0)
	_, err := balancer.GetOne()
	if err == nil {
		t.Error("WRoudRobinLoadZeroWeight GetOne have an error")
	}
}

func TestWRoudRobinLoad(t *testing.T) {
	domain := "www.google.com"
	var balancer = NewWRoundRobinLoad(domain)

	registryMap = nil
	newTarget(RegistNode{
		Domain: domain,
		Items: []OriginItem{
			{"192.168.1.100", 80},
			{"192.168.1.101", 40},
		},
	})
	balancer.AddAddr("192.168.1.102", 40)

	v1 := 0
	v2 := 0
	v3 := 0
	// loop times at 30
	for i := 0; i < 60; i++ {
		target, _ := balancer.GetOne()

		if target.Addr == "192.168.1.100" {
			v1++
		}
		if target.Addr == "192.168.1.101" {
			v2++
		}
		if target.Addr == "192.168.1.102" {
			v3++
		}
	}
	if v1 != 30 || v2 != 15 || v3 != 15 {
		t.Error("WRoundRobin GetOne have an error")
	}
}

func TestAddAddrWRoundRobin(t *testing.T) {
	domain := "www.google.com"
	var balancer = NewWRoundRobinLoad(domain)

	registryMap = nil
	newTarget(RegistNode{
		Domain: domain,
		Items: []OriginItem{
			{"192.168.1.100", 80},
		},
	})
	if len(registryMap[domain].Items) != 1 {
		t.Error("AddAddr func have an error #1")
	}

	balancer.AddAddr("192.168.1.101", 40)
	balancer.AddAddr("192.168.1.102", 40)

	if len(registryMap[domain].Items) != 3 {
		t.Error("AddAddr func have an error #2")
	}
}

func TestDelAddrWRoundRobin(t *testing.T) {
	domain := "www.google.com"
	var balancer = NewWRoundRobinLoad(domain)
	registryMap = nil
	newTarget(RegistNode{
		Domain: domain,
		Items: []OriginItem{
			{"192.168.1.100", 80},
			{"192.168.1.101", 80},
			{"192.168.1.102", 80},
		},
	})

	balancer.DelAddr("192.168.1.101")

	if len(registryMap[domain].Items) != 2 {
		t.Error("DelAddr func have an error #1")
	}
}

func TestGetMaxWeightIndex(t *testing.T) {
	items := []OriginItem{}
	_, err := getMaxWeight(items)
	if err == nil {
		t.Error("getMaxWeightIndex func have an error #1")
	}

	items = []OriginItem{
		{"192.168.137.100", 80},
		{"192.168.137.100", 130},
		{"192.168.137.100", 40},
		{"192.168.137.100", 20},
	}
	max2, err := getMaxWeight(items)

	if err != nil {
		t.Error("getMaxWeightIndex func have an error #2")
	}

	if max2 != 130 {
		t.Error("getMaxWeightIndex func have an error #3")
	}
}

func TestGetGCDWeight(t *testing.T) {
	items := []OriginItem{}
	_, err := getMaxWeight(items)
	if err == nil {
		t.Error("getGCDWeight func have an error #1")
	}
	items = []OriginItem{
		{"127.0.0.1:8081", 0},
	}
	gcdWeight, _ := getGCDWeight(items)
	if gcdWeight != 1 {
		t.Error("getGCDWeight func have an error #1.1")
	}

	items = []OriginItem{
		{"192.168.137.100", 80},
		{"192.168.137.100", 130},
		{"192.168.137.100", 40},
		{"192.168.137.100", 20},
	}
	gcdWeight, _ = getGCDWeight(items)
	if gcdWeight != 10 {
		t.Error("getGCDWeight func have an error #2")
	}

	items = []OriginItem{
		{"192.168.137.100", 244200},
		{"192.168.137.100", 111},
		{"192.168.137.100", 888},
	}

	gcdWeight, _ = getGCDWeight(items)
	if gcdWeight != 111 {
		t.Error("getGCDWeight func have an error #3")
	}

}

func TestGcd(t *testing.T) {
	r1 := gcd(0, 0)
	if r1 != 0 {
		t.Error("the gcd func have an error #1")
	}

	r2 := gcd(244200, 888)

	if r2 != 888 {
		t.Error("the gcd func have an error #2")
	}

	r3 := gcd(11, 244200)
	if r3 != 11 {
		t.Error("the gcd func have an error #3")
	}

	r4 := gcd(1800, 90)
	if r4 != 90 {
		t.Error("the gcd func have an error #4")
	}

}
