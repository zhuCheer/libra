package libra

import (
	"fmt"
	"github.com/zhuCheer/libra/balancer"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestProxyStart(t *testing.T) {
	proxy := NewHttpProxySrv("127.0.0.1:5000", nil)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("catch error", r)
			}
		}()
		proxy.Start()
	}()
}

func TestProxySrvFun(t *testing.T) {
	domain := "www.google.com"
	proxy := NewHttpProxySrv("127.0.0.1:5001", nil)
	proxy.RegistSite(domain, "roundrobin", "http")
	siteInfo, err := proxy.GetSiteInfo(domain)
	if err != nil {
		t.Error("NewHttpProxySrv loadType have an error #0")
	}

	if _, ok := siteInfo.Balancer.(*balancer.RoundRobinLoad); ok == false {
		t.Error("NewHttpProxySrv loadType have an error #1")
	}

	proxy.ResetCustomHeader(map[string]string{"X-LIBRA": "the smart ReverseProxy"})

	header, ok := proxy.customHeader["X-LIBRA"]
	if ok == false || header != "the smart ReverseProxy" {
		t.Error("ResetCustomHeader func have an error #2")
	}

	proxy.ChangeLoadType(domain, "wroundrobin")
	if _, ok := siteInfo.Balancer.(*balancer.WRoundRobinLoad); ok == false {
		t.Error("NewHttpProxySrv ChangeLoadType have an error #3")
	}

	proxy.ChangeLoadType(domain, "random")
	if _, ok := siteInfo.Balancer.(*balancer.RandomLoad); ok == false {
		t.Error("NewHttpProxySrv ChangeLoadType have an error #4")
	}

}

func TestReverseProxySrv(t *testing.T) {
	targetHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "testing ReverseProxySrv")
	}))
	defer targetHttpServer.Close()

	notfoundHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprint(w, "testing ReverseProxySrv NotFound")
	}))
	defer notfoundHttpServer.Close()

	gateway := "127.0.0.1:5002"
	proxy := NewHttpProxySrv(gateway, nil)
	proxy.ResetCustomHeader(map[string]string{"httptest": "01023"})
	proxy.RegistSite(gateway, "random", "http")

}

func TestReverseProxySrvUnStart(t *testing.T) {
	targetHttpServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "testing ReverseProxySrvUnStart")
	}))
	defer targetHttpServer.Close()

	gateway := "127.0.0.1:5005"
	proxy := NewHttpProxySrv(gateway, nil)
	proxy.ResetCustomHeader(map[string]string{"httptest": "01023"})
	proxy.RegistSite(gateway, "random", "http")
	go proxy.Start()

	targetHttpUrl, _ := url.Parse(targetHttpServer.URL)

	proxy.AddAddr(gateway, targetHttpUrl.Host, 0)
	res, _ := http.Get("http://" + gateway)

	if res.StatusCode != 502 {
		t.Error("ReverseProxySrv have an error(UnStart) #1")
	}
	time.Sleep(1 * time.Second)
}

func TestAddDelAddr(t *testing.T) {
	gateway := "127.0.0.1:5005"
	domain := "www.google.cn"
	proxy := NewHttpProxySrv(gateway, nil)
	proxy.RegistSite(domain, "random", "http").
		AddAddr(domain, "192.168.1.100", 0).
		AddAddr(domain, "192.168.1.100:8080", 0).
		AddAddr(domain, "192.168.1.100", 0).
		AddAddr(domain, "192.168.1.101", 0)

	siteInfo, _ := proxy.GetSiteInfo(domain)
	if len(siteInfo.Items) != 3 {
		t.Error("proxy AddAddr have an error #1")
	}
	proxy.DelAddr(domain, "192.168.1.105")
	proxy.DelAddr(domain, "192.168.1.100")

	if len(siteInfo.Items) != 2 {
		t.Error("proxy DelAddr have an error #2")
	}
}

func TestGetSiteInfo(t *testing.T) {
	gateway := "127.0.0.1:5009"
	domain := "www.google.cn"
	proxy := NewHttpProxySrv(gateway, nil)
	proxy.RegistSite(domain, "random", "http")

	info1, _ := proxy.GetSiteInfo(domain)
	info2, _ := balancer.GetSiteInfo(domain)
	if info1.Domain != info2.Domain {
		t.Error("proxy GetSiteInfo have an error #1")
	}
	if _, ok := info1.Balancer.(*balancer.RandomLoad); !ok {
		t.Error("proxy.GetSiteInfo have an error #2")
	}
	if _, ok := info2.Balancer.(*balancer.RandomLoad); !ok {
		t.Error("balancer.GetSiteInfo have an error #3")
	}
	if info1.Balancer != info1.Balancer {
		t.Error("balancer.GetSiteInfo have an error #4")
	}

}

func TestSingleJoiningSlash(t *testing.T) {
	target, _ := url.Parse("http://192.168.1.100/")
	path := singleJoiningSlash(target.Path, "/")
	if path != "/" {
		t.Error("singleJoiningSlash func have an error #1")
	}

	target, _ = url.Parse("http://192.168.1.100/abc")
	path = singleJoiningSlash(target.Path, "/")
	if path != "/abc" {
		t.Error("singleJoiningSlash func have an error #2")
	}

	target, _ = url.Parse("http://192.168.1.100/abc")
	path = singleJoiningSlash(target.Path, "/efg")
	if path != "/abc/efg" {
		t.Error("singleJoiningSlash func have an error #3")
	}

	target, _ = url.Parse("http://192.168.1.100/abc/")
	path = singleJoiningSlash(target.Path, "/efg")
	if path != "/abc/efg" {
		t.Error("singleJoiningSlash func have an error #4")
	}

	target, _ = url.Parse("http://192.168.1.100/abc")
	path = singleJoiningSlash(target.Path, "efg")
	if path != "/abc/efg" {
		t.Error("singleJoiningSlash func have an error #5")
	}

}
