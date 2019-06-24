package libra

import (
	"bytes"
	"context"
	"errors"
	"github.com/zhuCheer/libra/balancer"
	"github.com/zhuCheer/libra/logger"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// Proxy server node struct
type ProxySrv struct {
	ProxyAddr  string
	ManageAddr string
	Scheme     string
	LoadType   string
}

// Common errors.
var (
	ErrNotRegistTarget = errors.New("you should regist target server first")
	Logger             = logger.NoopLogger{}
)

// start http proxy server
func (p *ProxySrv) Start() error {

	go func() {
		p.startManagerHttpSrv()
	}()

	proxyHttpMux := http.NewServeMux()
	Logger.Printf("start proxy server bind " + p.ProxyAddr)
	proxyHttpMux.Handle("/", p.dynamicReverseProxy())

	proxyServer := &http.Server{
		Addr:    p.ProxyAddr,
		Handler: proxyHttpMux,
	}
	err := proxyServer.ListenAndServe()
	panic(err)
}

// verify the ProxyAddr and ManageAddr
func (p *ProxySrv) verifyAddr() error {
	if p.ManageAddr == p.ProxyAddr {
		errors.New("ProxyAddr can not the same of ManageAddr")
	}
	return nil
}

// get a balancer
func (p *ProxySrv) getBalancerRemote(domain string) (*balancer.ProxyTarget, error) {
	b := balancer.NewRandomLoad()

	switch p.LoadType {
	case "random":
		b = balancer.NewRandomLoad()
	}

	return b.GetOne(domain)
}

// start a manager http server
func (p *ProxySrv) startManagerHttpSrv() {
	mux := http.NewServeMux()
	mux.HandleFunc("/notFoundPage", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not found remote page"))
	})
	Logger.Printf("start manage http server bind addr:" + p.ManageAddr)
	err := http.ListenAndServe(p.ManageAddr, mux)
	if err != nil {
		panic(err)
	}
}

// get ReverseProxy dynamic director func
// in this function proxy server knows where to forward to
// if the target is a error node, proxy will forward to a default error page in local address.
func (p *ProxySrv) dynamicDirector(req *http.Request) {
	proxyTarget, err := p.getBalancerRemote(req.Host)
	var target *url.URL
	if err == nil && proxyTarget != nil {
		target, err = url.Parse(p.Scheme + "://" + proxyTarget.Addr)
	}

	// if err not nil reverse proxy to a default error page
	if err != nil || target == nil {
		req.URL.Host = p.ManageAddr
		req.URL.Path = "/notFoundPage"
	} else {
		targetQuery := target.RawQuery
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	req.URL.Scheme = p.Scheme
}

// get ReverseProxy Http Handler
func (p *ProxySrv) dynamicReverseProxy() *httputil.ReverseProxy {
	roundTripper := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext(ctx, network, addr)
		},
		MaxIdleConns:          100,
		DisableKeepAlives:     false,
		IdleConnTimeout:       10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	transport := &transport{RoundTripper: roundTripper}

	httpProxy := &httputil.ReverseProxy{
		Director:  p.dynamicDirector,
		Transport: transport,
	}
	return httpProxy
}

// Implementing RoundTripper interface
type transport struct {
	http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return getDefaultErrorPage(502, err.Error(), req)
	}

	if resp.StatusCode > 400 {
		return getDefaultErrorPage(resp.StatusCode, "have an error", req)
	}

	remoteBody, _ := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewReader(remoteBody))
	return resp, nil
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
