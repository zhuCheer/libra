# libra 一个动态的反向代理网关包

## 相关特性
- 动态进行多地址反向代理
- 动态管理源站地址
- 支持自定义头

## 使用方法
 
 ### 第一步 注册节点
 
 将一个域名和对应的ip:port进行注册
 ```
 domain :="www.qiproxy.cn"
	balancer.RegistTarget(domain,[]string{})
	balancer.AddAddrWithoutWeight(domain, []string{
		"192.168.1.101:8081",
		"192.168.1.101:8082",
	}...)
 ```
 
 ### 第二步 启动服务
 
 - 目前支持http代理，绑定代理端端口和管理端端口，管理端目前只包含一个错误页面，后续会增加相关管理接口;
- 当目标页面不可达或其他异常都会输出网关自定义的错误页面; 

 ```
 srv:=libra.ProxySrv{
		ProxyAddr:"127.0.0.1:5000",
		ManageAddr:"127.0.0.1:5001",
		Scheme:"http",
	}
	
	srv.Start()
 ```
 
 
 ### 第三步 启动
 
 直接通过`srv.Start()`就能将服务启动，实现一个高效的负载均衡代理器；