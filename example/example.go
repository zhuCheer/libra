package main

import (
	"github.com/zhuCheer/libra"
	"net/http"
)

func main() {
	go httpsrv01()
	go httpsrv02()
	go httpsrv03()
	var srv = libra.NewHttpProxySrv("127.0.0.1:5000", "roundrobin", nil)
	srv.GetBalancer().AddAddr("127.0.0.1:5000", "127.0.0.1:5001", 1)
	srv.GetBalancer().AddAddr("127.0.0.1:5000", "127.0.0.1:5002", 1)
	srv.GetBalancer().AddAddr("127.0.0.1:5000", "127.0.0.1:5003", 1)
	srv.Scheme = "http"
	srv.Start()
}

func httpsrv01() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.Write([]byte("view http server 01 <br> source code: <a href='https://github.com/zhuCheer/libra'>https://github.com/zhuCheer/libra</a>"))
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

func httpsrv03() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("view http server 03"))
	})
	http.ListenAndServe(":5003", mux)
}
