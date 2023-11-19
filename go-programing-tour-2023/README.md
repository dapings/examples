# Book: Go 语言编程之旅：一起用 Go 做项目

[tour-book source code](https://github.com/go-programming-tour-book)
[how-can-i-track-tool-dependencies-for-a-module](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module)

[protocol buffers documentation](https://protobuf.dev/)
[protocol buffers go tutorial](https://protobuf.dev/getting-started/gotutorial/)
[protocol buffers - protobuf](https://github.com/protocolbuffers/protobuf)

[grpc.io](https://grpc.io/)
[protoc-installation](https://grpc.io/docs/protoc-installation/)
[grpc-go-quick-start](https://grpc.io/docs/languages/go/quickstart/)
[golang-protobuf](https://github.com/golang/protobuf/)
[grpc/grpc](https://github.com/grpc/grpc)
[grpc/grpc-go](https://github.com/grpc/grpc-go)
[grpc-ecosystem grpc-gateway](https://grpc-ecosystem.github.io/grpc-gateway/)
[grpcurl](https://github.com/fullstorydev/grpcurl)

[go-bindata 数据文件转换为Go代码，摆脱静态资源文件](https://github.com/go-bindata/go-bindata)

1. tour 命令行程序

    功能清单：
    - 单词格式转换
    - 便利的时间工具
    - SQL语句到结构体的转换
    - JSON到结构体的转换

2. blog-service 博客程序

   HTTP应用：写一个完整的博客后端
   - 博客之旅
   - 项目设计
   - 公共组件
   - 接口文档
   - 接口校验
   - 模块开发：标签管理
   - 模块开发：文章管理
   - 上传图片和文件服务
   - API访问控制
      考虑做纵深防御，对API接口进行访问控制，比较常见的API访问控制方案有两种，分别是OAth 2.0, JWT(JSON Web Tokens)，它们完全不同，对应的应用场景也不一样，具体如下：
      - OAuth 2.0：本质上 是一个授权的行业标准协议，提供了一整套授权机制的指导标准，常用于使用第三方站点账号关联登录的情况。相对会重一些，常常还会授予第三方应用获取对应账号的个人基本信息等。
      - JWT：与OAuth 2.0 完全不同，常用于前后端分离的情况，能够非常便捷地给API接口提供安全鉴权。用于在各方之间以JSON对象安全地传输信息，且信息是经过数字签名的。JWT令牌的内容是非严格加密的。
        JWT 以紧凑的形式由三部分组成，以点分隔
        - Header 头部通常由两部分组成，分别是令牌的类型和所使用的签名算法，用于描述其元数据，最后用base64UrlEncode算法对此JSON对象进行转换。
        - Payload 有效载荷，用于存储在JWT中实际传输的数据，包括过期时间、签发者等，最后用base64UrlEncode算法对此JSON对象进行转换。对于敏感信息建议不要放到JWT中，如果一定要放，则应进行不可逆加密处理，确保信息的安全性。
        - Signature 签名是对Header,Payload进行约定算法和规则的签名，用于校验消息在整个过程中有没有被篡改。生成签名：HMACSHA256(base64UrlEncode(header)+"."+base64UrlEncode(payload),secret)
   - 常见应用中间件
   - 链路追踪
   - 应用配置问题
     1. 命令行参数
     2. 系统环境变量
        也可以将配置文件存放在系统自带的全局变量中，如$HOME/conf或/etc/conf中，好处是不需要重新自定义一个新的系统环境变更
        内置一些系统环境变量的读取，优先级低于命令行参数，但高于文件配置。
     3. 打包进二进制文件中
        或者将配置文件打包到二进制文件中，通过 go-bindata 库可以将数据文件转换为Go代码，就可以摆脱静态资源文件了。
        ```shell
         go get -u github.com/go-bindata/go-bindata/...
         go-bindata -o configs/config.go -pkg=configs configs/config.yaml
         # 执行如下代码，读取对应的文件内容
         b, _ := configs.Assert("configs/config.yaml")
        ```
     4. 其他方案
        或直接使用集中式的配置中心。
   - 编译程序应用
   - 优雅重启和停止
     通过信号量的方式来解决问题：优雅重启和停止。
     信号量是一种异步通知机制，用来提醒进程一个事件(硬件异常、程序执行异常、外部发生信息)已发生。如果进程定义的信号的处理函数，那么它将被执行，否则执行默认的处理函数。
   - 思考

   技术框架选型、需求分析
   数据库设计
   - 标签管理：文章所归属的分类，用于标识文章内容的要点和要素，以便读者识别和SEO收录等
   - 文章管理：对整个文章内容的管理，并把文章和标签进行关联
   项目工程设计、API 文件编写、业务接口编写、接口访问控制、链路追踪
   TODO: 
   - 实现标签和标签接口去重的判断逻辑：在新增标签和文章时，判断是否已经存在
   - 支持给一篇文章设置多个标签信息
   - 支持多张图片上传
   - 利用Redis实现Token Bucket，支持分布式的限流器

3. tag-service RPC 程序
   
   RPC应用(gRPC)实现：
   - gRPC 和 Protobuf 简介
   - Protobuf 的使用
   
     protoc 是 protobuf 的编译器，是用 C++ 编写的，其主要功能是编译 .proto 文件。
     参照 `protoc-installation` 安装 protoc。在找不到安装的动态链接库的特定情况下，需要手动执行 `ldconfig` 命令，让动态链接库为系统所共享。也就是说，ldconfig 是一个动态链接库管理命令。
     ```shell
     wget https://github.com/protocolbuffers/protobuf/releases/download/v3.15.7/protoc-3.15.7-osx-x86_64.zip
     unzip -d protoc-3.15.7-osx protoc-3.15.7-osx-x86_64.zip
     ln -fs protoc-3.15.7-osx current
     ```
     
     仅安装protoc编译器是不够的，针对不同的语音，还需要安装运行时的 protoc 插件，而对就 Go 的是 protoc-gen-go 插件、protoc-gen-go-grpc 插件。
     ```shell
     # module github.com/golang/protobuf is deprecated, use the "google.golang.org/protobuf" module instead.
     # go get -d -u -v github.com/golang/protobuf/protoc-gen-go
     # google.golang.org/protobuf=github.com/golang/protobuf
     go install google.golang.org/protobuf/cmd/protoc-gen-go@v2.15.2
     
     # module declares its path as: google.golang.org/grpc/cmd/protoc-gen-go-grpc, but was required as: github.com/grpc/grpc-go/cmd/protoc-gen-go-grpc
     # go get -d -u -v github.com/grpc/grpc-go/cmd/protoc-gen-go-grpc
     # google.golang.org/grpc=github.com/grpc/grpc-go
     go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
     #export PATH="$PATH:$(go env GOPATH)/bin"
     ```

   - gRPC 的使用

     在项目下，安装 Go gRPC 库：
     ```shell
     # github.com/grpc/grpc-go v1.59.0 
     # google.golang.org/grpc=github.com/grpc/grpc-go
     go get -u -v google.golang.org/grpc
     ```

   - 运行一个 gRPC 服务

     生成 proto 文件：`protoc --go_out=plugins=grpc:. ./proto/*.proto`
     gRPC 是基于 HTTP/2 协议的，不能直接通过 postman 或普通的 curl 进行调用，目前开源社区的方案：命令行工具 grpcurl，安装及使用命令如下：
     ```shell
     go get -u -v github.com/fullstorydev/grpcurl
     go install github.com/fullstorydev/grpcurl/cmd/grpcurl
     # 默认使用 TLS 认证(-cert,-key 设置公钥和密钥)，-plaintext 用来忽略TLS认证
     grpcurl -plaintext localhost:8001 list
     grpcurl -plaintext localhost:8001 list proto.TagService
     grpcurl -plaintext -d '{"name": "Go"}' localhost:8001 list proto.TagService.GetTagList
     ```

   - gRPC 服务间的内调
   - 提供 HTTP 接口

     开源社区的 grpc-gateway 是 protoc 的一个插件，能够读取 protobuf 的服务定义，生成一个反向代理服务器，将 Restful JSON API 转换为 gRPC。它主要是根据 protobuf 的服务定义中的 google.api.http 来生成的。
     简单来说，grpc-gateway 能够将 Restful 转换为 gRPC 请求，实现用同一个RPC方法提供对gRPC协议和HTTP/1.1 的双流量支持。
     ```shell
     # https://grpc-ecosystem.github.io/grpc-gateway/
     # github.com/grpc-ecosystem/grpc-gateway
     go install \
     github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
     github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
     google.golang.org/protobuf/cmd/protoc-gen-go \
     google.golang.org/grpc/cmd/protoc-gen-go-grpc
     
     # -I 参数的格式：-IPATH,--proto_path=PATH，用来指定 proto 文件中 import 搜索的目录，可指定多个。如果不指定，则默认是当前工作目录。
     protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. ./proto/*.proto
     ```

   - 接口文档
   - gRPC 拦截器
   - metadata 和 RPC 自定义认证
   - 链路追踪
   - gRPC 服务注册和发现
   - 实现自定义的 protoc 插件
   - 对 gRPC 接口进行版本管理
   - 常见问题讨论

4. chatroom IM聊天室
5. cache-example 进程内缓存