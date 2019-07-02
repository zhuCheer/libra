package main

import (
	"github.com/zhuCheer/libra"
	"net/http"
)

func main() {
	go httpsrv01()
	go httpsrv02()
	var srv = libra.NewHttpProxySrv("127.0.0.1:5000", "roundrobin", nil)
	srv.GetBalancer().AddAddr("127.0.0.1:5000", "127.0.0.1:5001", 1)
	srv.GetBalancer().AddAddr("127.0.0.1:5000", "127.0.0.1:5002", 1)
	srv.Scheme = "http"
	srv.Start()
}

func httpsrv01() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("view http server 01"))
	})
	http.ListenAndServe(":5001", mux)
}

func httpsrv02() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("view http server 02"))
	})
	http.ListenAndServe(":5002", mux)
}
