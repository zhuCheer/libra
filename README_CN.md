# libra 一个动态的反向代理均衡器

[English document](https://github.com/zhuCheer/libra/blob/master/README.md)， [中文文档](https://github.com/zhuCheer/libra/blob/master/README_CN.md)

## 相关特性
- 动态进行多地址反向代理
- 动态管理源站地址
- 支持动态修改定响应头
- 代码都有严格的单元测试，可靠度高

你可以通过此包来快速构建一个动态的负载均衡服务器，目前有三种负载均衡算法，分别是随机，轮询，带权轮询；此包皆由原生代码构建，不依赖其他第三方包；如果你想快速的获得一个负载均衡服务器，快来使用它吧。


## 快速开始


#### 下载安装

在控制台输入
`go get github.com/zhuCheer/libra`

#### 运行示例

进入到包中 example 目录，然后运行 example.go 即可开启一个反向代理服务器
```
> cd ../src/github.com/zhuCheer/libra/example
> go run example.go

```

此时可以通过浏览器访问 `http://127.0.0.1:5000` 即可看到代理效果，会以轮询的方式访问`http://127.0.0.1:5001`, `http://127.0.0.1:5002` 这两个 http 服务。


#### 详细步骤说明
```
import "github.com/zhuCheer/libra"

    
// 注册一个反向代理服务器，反代服务器访问的 ip 和端口，负载均衡类型，和自定义响应头 三个参数
// 负载均衡类型目前有三种可选 random:随机，roundrobin:轮询，wroundrobin:带权轮询
srv := libra.NewHttpProxySrv("127.0.0.1:5000", "roundrobin", nil)


// 添加代理目标的域名和 ip 地址端口
srv.GetBalancer().AddAddr("www.yourappdomain.com", "127.0.0.1:5001", 1)
srv.GetBalancer().AddAddr("www.yourappdomain.com", "127.0.0.1:5002", 1)


// 启动反向代理服务
srv.Start()
```

#### 原理详解

- 将我们的各类应用得域名都直接解析到此代理服务器，我们可称此服务为一个网关；
- 反向代理服务启动后，我们通过不同的域名即可访问到我们添加的指定 ip 上，大多数情况我们会有多个服务器，访问请求将按照指定的负载均衡算法进行调度；
- 访问过程如下图所示：

![image](https://img.douyucdn.cn/data/yuba/weibo/2019/07/02/201907021730116899917826388.gif)

## 相关操作方法介绍

```
import "github.com/zhuCheer/libra"
srv := libra.NewHttpProxySrv("127.0.0.1:5000", "roundrobin", nil)


// 设置响应头
srv.ResetCustomHeader(map[string]string{"X-LIBRA": "the smart ReverseProxy"})

// 切换负载均衡类型
srv.ChangeLoadType("random")


// 添加目标服务器节点信息，添加后即生效，无需重启
srv.GetBalancer().AddAddr("www.yourappdomain.com","192.168.1.100:8081", 1)

// 删除目标服务器节点信息，添加后即生效，无需重启
srv.GetBalancer().DelAddr("www.yourappdomain.com","192.168.1.100:8081")

```