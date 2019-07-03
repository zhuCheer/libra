# libra is a dynamic reverse proxy / load balancer

[![Build Status](https://travis-ci.org/zhuCheer/libra.svg?branch=master)](https://travis-ci.org/zhuCheer/libra) [![Go Report Card](https://goreportcard.com/badge/github.com/zhuCheer/libra)](https://goreportcard.com/report/github.com/zhuCheer/libra) [![GoDoc](https://godoc.org/github.com/zhuCheer/libra?status.svg)](https://godoc.org/github.com/zhuCheer/libra)


[English document](https://github.com/zhuCheer/libra/blob/master/README.md)， [中文文档](https://github.com/zhuCheer/libra/blob/master/README_CN.md)

## Feature
- dynamic and multiple reverse proxy server
- dynamic change origin server addr
- dynamic change response header
- rigorous unit testing

You can use this package build a dynamic reverse proxy server faster, now it has three load balance algorithms, random, roundrobin, wroundrobin (round robin with weight)

## Getting Started

#### Installation

To install this package, you need to install Go and setup your Go workspace on your computer. The simplest way to install the library is to run:

`go get github.com/zhuCheer/libra`

#### run example
Change directory to libra package and run the example.go, you can start a reverse proxy.
```
> cd ../src/github.com/zhuCheer/libra/example
> go run example.go

```

Now,you can open browser to `http://127.0.0.1:5000`,you will see the reverse proxy, it running with round robin balance to `http://127.0.0.1:5001` and `http://127.0.0.1:5002` http server.


#### How to use

```
import "github.com/zhuCheer/libra"

    
// regist a reverse proxy，input three params bind ip:port, balancer algorithm, custom response header
// it has three balancer algorithm random，roundrobin，wroundrobin(round robin with weight)
srv := libra.NewHttpProxySrv("127.0.0.1:5000", "roundrobin", nil)


// add target domain and ip:port
srv.GetBalancer().AddAddr("www.yourappdomain.com", "127.0.0.1:5001", 1)
srv.GetBalancer().AddAddr("www.yourappdomain.com", "127.0.0.1:5002", 1)


// start reverse proxy server
srv.Start()
```


#### Principles

- The reverse proxy serves as a gateway between users and your application origin server. In so doing it handles all policy management and traffic routing;
- A reverse proxy operates by:
- 1. Receiving a user connection request
- 2. Completing a TCP three-way handshake, terminating the initial connection
- 3. Connecting with the origin server and forwarding the original request

![image](https://img.douyucdn.cn/data/yuba/weibo/2019/07/02/201907021730116899917826388.gif)


## Functions

```
import "github.com/zhuCheer/libra"
srv := libra.NewHttpProxySrv("127.0.0.1:5000", "roundrobin", nil)


// set response header
srv.ResetCustomHeader(map[string]string{"X-LIBRA": "the smart ReverseProxy"})

// change balance algorithm
srv.ChangeLoadType("random")


// add origin server addr, dynamic change without restarting
srv.GetBalancer().AddAddr("www.yourappdomain.com","192.168.1.100:8081", 1)

// delete origin server addr, dynamic change without restarting
srv.GetBalancer().DelAddr("www.yourappdomain.com","192.168.1.100:8081")

```

## Contributors
- [@Chase](https://www.facebook.com/profile.php?id=100017355485621)


## License

[MIT](./LICENSE)