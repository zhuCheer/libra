package balancer

import (
	"sync"
	"testing"
)

func TestRoudRobinLoad(t *testing.T) {
	var balancer = NewRoundRobinLoad()
	target, err := balancer.GetOne("name")
	if err == nil {
		t.Error("NewRoundRobinLoad func have an error #1")
	}

	domain := "www.google.com"
	registryMap = nil
	NewTarget(RegistNode{
		Domain: domain,
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

	// concurrency test should be ok
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			target, _ = balancer.GetOne(domain)
			wg.Done()
		}()
	}
	wg.Wait()

	// loop times at 5
	if target.Addr != "192.168.1.104" {
		t.Error("RoundRobinLoad->getOne func have an error #2")
	}

	// loop times at 30
	for i := 0; i < 25; i++ {
		target, _ = balancer.GetOne(domain)
	}

	if target.Addr != "192.168.1.102" {
		t.Error("RoundRobinLoad->getOne func have an error #3")
	}

}
