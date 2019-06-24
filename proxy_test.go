package libra

import (
	"github.com/zhuCheer/libra/balancer"
	"testing"
)

func TestGetBalancerRemote(t *testing.T) {
	proxy := ProxySrv{
		Ip:        "127.0.0.1",
		ProxyPort: 5000,
		Scheme:    "http",
	}

	domain := "www.ks58.cc"
	target, err := proxy.getBalancerRemote(domain)
	if err == nil {
		t.Error("proxy.getBalancerRemote have an error #1")
	}

	balancer.RegistTarget(domain, []string{})
	balancer.AddAddrWithoutWeight(domain, "192.168.137.100")
	balancer.AddAddrWithoutWeight(domain, "192.168.137.101")
	balancer.AddAddrWithoutWeight(domain, "192.168.137.102")

	target, err = proxy.getBalancerRemote(domain)

	if err != nil || target == nil {
		t.Error("proxy.getBalancerRemote have an error #2")
	}

}
