package libra

import (
	"fmt"
	"github.com/zhuCheer/libra/balancer"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestProxyStart(t *testing.T) {
	proxy := NewHttpProxySrv("127.0.0.1:5001", "roundrobin", nil)

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
	proxy := NewHttpProxySrv("127.0.0.1:5000", "roundrobin", nil)
	if _, ok := proxy.balancer.(*balancer.RoundRobinLoad); ok == false {
		t.Error("NewHttpProxySrv loadType have an error #1")
	}
	proxy.ResetCustomHeader(map[string]string{"X-LIBRA": "the smart ReverseProxy"})

	header, ok := proxy.customHeader["X-LIBRA"]
	if ok == false || header != "the smart ReverseProxy" {
		t.Error("ResetCustomHeader func have an error #2")
	}

	proxy.ChangeLoadType("wroundrobin")
	if _, ok := proxy.balancer.(*balancer.WRoundRobinLoad); ok == false {
		t.Error("NewHttpProxySrv ChangeLoadType have an error #3")
	}

	proxy.ChangeLoadType("random")
	if _, ok := proxy.balancer.(*balancer.RandomLoad); ok == false {
		t.Error("NewHttpProxySrv ChangeLoadType have an error #4")
	}
}

func TestReverseProxySrv(t *testing.T) {
	targetHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "testing ReverseProxySrv")
	}))
	defer targetHttpServer.Close()

	proxy := NewHttpProxySrv("127.0.0.1:5000", "roundrobin", nil)
	reverseProxy := proxy.dynamicReverseProxy()
	proxy.ResetCustomHeader(map[string]string{"httptest": "01023"})
	ts := httptest.NewServer(proxy.httpMiddleware(reverseProxy))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 500 {
		t.Error("ReverseProxySrv have an error #1")
	}

	testHeader := res.Header.Get("httptest")
	if testHeader != "01023" {
		t.Error("ReverseProxySrv have an error #2")
	}

	tsUrl, _ := url.Parse(ts.URL)
	targetHttpUrl, _ := url.Parse(targetHttpServer.URL)
	proxy.balancer.AddAddr(tsUrl.Host, targetHttpUrl.Host, 0)
	res, err = http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Error("ReverseProxySrv have an error #3")
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if string(greeting) != "testing ReverseProxySrv" {
		t.Error("ReverseProxySrv have an error #4")
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
	if path != "/abc/" {
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

}
