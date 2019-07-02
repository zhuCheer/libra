package libra

import (
	"bytes"
	"context"
	"crypto/tls"
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

// ProxySrv Proxy server node struct
type ProxySrv struct {
	ProxyAddr    string
	Scheme       string
	customHeader map[string]string
	loadType     string
	balancer     balancer.Balancer
}

// Common variable.
var (
	Logger      = logger.NoopLogger{}
	errorHeader = "x-libra-err"
	version     = "v0.0.1"
	githubUrl   = "https://github.com/zhuCheer/libra"
)

// NewHttpProxySrv new http reverse proxy
func NewHttpProxySrv(addr string, loadType string, header map[string]string) *ProxySrv {
	if header == nil {
		header = map[string]string{}
	}

	return &ProxySrv{
		ProxyAddr:    addr,
		Scheme:       "http",
		customHeader: header,
		balancer:     getBalancerByLoadType(loadType),
		loadType:     loadType,
	}
}

// Start http proxy server
func (p *ProxySrv) Start() error {
	if p.Scheme == "" {
		p.Scheme = "http"
	}

	proxyHttpMux := http.NewServeMux()
	Logger.Printf("start proxy server bind " + p.ProxyAddr)
	proxyHttpMux.Handle("/", p.httpMiddleware(p.dynamicReverseProxy()))

	proxyServer := &http.Server{
		Addr:    p.ProxyAddr,
		Handler: proxyHttpMux,
	}
	err := proxyServer.ListenAndServe()
	panic(err)
}

// GetBalancer get an balancer point
func (p *ProxySrv) GetBalancer() balancer.Balancer {

	return p.balancer
}

// ChangeLoadType change balancer loadType
func (p *ProxySrv) ChangeLoadType(loadType string) {
	b := getBalancerByLoadType(loadType)
	p.balancer = b
}

// ResetCustomHeader reset custom header
func (p *ProxySrv) ResetCustomHeader(header map[string]string) {
	p.customHeader = map[string]string{}
	p.customHeader = header
}

// getBalancerByLoadType get balancer by lodeType
func getBalancerByLoadType(loadType string) balancer.Balancer {
	b := balancer.NewRandomLoad()
	switch loadType {
	case "random":
		b = balancer.NewRandomLoad()
	case "roundrobin":
		b = balancer.NewRoundRobinLoad()
	case "wroundrobin":
		b = balancer.NewWRoundRobinLoad()
	}

	return b
}

// httpMiddleware http middleware set some header
func (p *ProxySrv) httpMiddleware(handler *httputil.ReverseProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, value := range p.customHeader {
			w.Header().Add(key, value)
		}
		w.Header().Set("X-LIBRA-VERSION", version)
		w.Header().Set("X-LIBRA-CODE", githubUrl)
		handler.ServeHTTP(w, r)
	})
}

// dynamicDirector get ReverseProxy dynamic director func
// in this function proxy server knows where to forward to
// if the target is a error node, proxy will forward to a default error page in local address.
func (p *ProxySrv) dynamicDirector(req *http.Request) {
	proxyTarget, err := p.balancer.GetOne(req.Host)
	var target *url.URL
	if err == nil && proxyTarget != nil {
		target, err = url.Parse(p.Scheme + "://" + proxyTarget.Addr)
	}

	// if err not nil wirte an err in header
	if err != nil {
		req.Header.Set(errorHeader, err.Error())
	} else {
		targetQuery := target.RawQuery
		//req.Host = target.Host
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
	Logger.Printf("proxy to " + req.URL.String())
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
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
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

// RoundTrip http transport
func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	proxyErrHeader := req.Header.Get(errorHeader)
	if proxyErrHeader != "" {
		return getDefaultErrorPage(500, proxyErrHeader, req)
	}

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
