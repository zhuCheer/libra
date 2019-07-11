package libra

import (
	"fmt"
	"github.com/zhuCheer/libra/balancer"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

	time.Sleep(2 * time.Second)
}

func TestProxySrvFun(t *testing.T) {
	domain := "www.google.com"
	proxy := NewHttpProxySrv("127.0.0.1:5001", nil)
	proxy.RegistSite(domain, "roundrobin", "http").GetBalancer(domain)
	balancerDo, err := proxy.GetBalancer(domain)
	if err != nil {
		t.Error("NewHttpProxySrv loadType have an error #0")
	}

	if _, ok := balancerDo.(*balancer.RoundRobinLoad); ok == false {
		t.Error("NewHttpProxySrv loadType have an error #1")
	}

	proxy.ResetCustomHeader(map[string]string{"X-LIBRA": "the smart ReverseProxy"})

	header, ok := proxy.customHeader["X-LIBRA"]
	if ok == false || header != "the smart ReverseProxy" {
		t.Error("ResetCustomHeader func have an error #2")
	}

	proxy.ChangeLoadType(domain, "wroundrobin")
	balancerDo, _ = proxy.GetBalancer(domain)
	if _, ok := balancerDo.(*balancer.WRoundRobinLoad); ok == false {
		t.Error("NewHttpProxySrv ChangeLoadType have an error #3")
	}

	proxy.ChangeLoadType(domain, "random")
	balancerDo, _ = proxy.GetBalancer(domain)
	if _, ok := balancerDo.(*balancer.RandomLoad); ok == false {
		t.Error("NewHttpProxySrv ChangeLoadType have an error #4")
	}

}

func TestReverseProxySrv(t *testing.T) {
	targetHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "testing ReverseProxySrv")
	}))
	defer targetHttpServer.Close()

	gateway := "127.0.0.1:5002"
	proxy := NewHttpProxySrv(gateway, nil)
	proxy.ResetCustomHeader(map[string]string{"httptest": "01023"})
	proxy.RegistSite(gateway, "random", "http")
	go proxy.Start()

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+gateway, strings.NewReader(""))
	if err != nil {
		t.Error(err)
	}
	res, err := client.Do(req)
	defer res.Body.Close()

	if res.StatusCode != 500 {
		t.Error("ReverseProxySrv have an error #1")
	}

	testHeader := res.Header.Get("httptest")
	if testHeader != "01023" {
		t.Error("ReverseProxySrv have an error #2")
	}

	targetHttpUrl, _ := url.Parse(targetHttpServer.URL)

	balancerDo, _ := proxy.GetBalancer(gateway)
	balancerDo.AddAddr(targetHttpUrl.Host, 0)

	res, err = http.Get("http://" + gateway + "?abc=123")

	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Error("ReverseProxySrv have an error #3")
	}

	if res.Request.URL.RawQuery != "abc=123" {
		t.Error("ReverseProxySrv have an error #4")
	}

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Error(err)
	}
	if string(greeting) != "testing ReverseProxySrv" {
		t.Error("ReverseProxySrv have an error #5")
	}

	time.Sleep(2 * time.Second)
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

	balancerDo, _ := proxy.GetBalancer(gateway)
	balancerDo.AddAddr(targetHttpUrl.Host, 0)
	res, _ := http.Get("http://" + gateway)

	if res.StatusCode != 502 {
		t.Error("ReverseProxySrv have an error(UnStart) #1")
	}
	time.Sleep(1 * time.Second)
}

func TestReverseProxySrvNotFound(t *testing.T) {
	targetHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprint(w, "testing ReverseProxySrv NotFound")
	}))
	defer targetHttpServer.Close()

	gateway := "127.0.0.1:5006"
	proxy := NewHttpProxySrv(gateway, nil)
	proxy.RegistSite(gateway, "random", "http")
	go proxy.Start()

	targetHttpUrl, _ := url.Parse(targetHttpServer.URL)
	balancerDo, _ := proxy.GetBalancer(gateway)
	balancerDo.AddAddr(targetHttpUrl.Host, 0)
	res, _ := http.Get("http://" + gateway)

	if res.StatusCode != 404 {
		t.Error("ReverseProxySrv have an error(NotFound) #1")
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
