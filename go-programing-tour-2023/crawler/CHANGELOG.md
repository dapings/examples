# CHANGELOG

## 2023/12/26
* 根据pod ip 生成分布式id
* 添加 docker-compose.yaml, k8s master, worker yaml files
* 添加 auth micro 中间件
* worker task add limiter
* master ratelimit(juju/ratelimit) and breaker(hystrix-go) for grpc server
* master请求转发到leader

## 2023/12/25
* Dockerfile, k8s worker deployment,service yaml file
* http pprof
* worker故障容错，任务分配到其他节点
* master GRPC
* 任务分配，查到最小负载
* master成为leader后加载资源
* master添加初始的种子任务
* master简单的任务分配

## 2023/12/22
* master elect and watch worker change
* proxy fuzz test
* sqldb cover test