# Book: [Go 语言编程之旅：一起用 Go 做项目](https://golang2.eddycjy.com/)

[tour-book source code](https://github.com/go-programming-tour-book)
[how-can-i-track-tool-dependencies-for-a-module](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module)

[protocol buffers documentation](https://protobuf.dev/)
[protocol buffers go tutorial](https://protobuf.dev/getting-started/gotutorial/)
[protocol buffers - protobuf](https://github.com/protocolbuffers/protobuf)
[go generated code guide](https://protobuf.dev/reference/go/go-generated/#package)

[grpc.io](https://grpc.io/)
[protoc-installation](https://grpc.io/docs/protoc-installation/)
[grpc-go-quick-start](https://grpc.io/docs/languages/go/quickstart/)
  - [regenerate grpc code](https://grpc.io/docs/languages/go/quickstart/#regenerate-grpc-code)
  - [examples](https://github.com/grpc/grpc-go/tree/master/examples)
[golang-protobuf](https://github.com/golang/protobuf/)
[grpc/grpc](https://github.com/grpc/grpc)
[grpc/grpc-go](https://github.com/grpc/grpc-go)

[grpc-ecosystem grpc-gateway](https://grpc-ecosystem.github.io/grpc-gateway/)
  - [customizing your gateway](https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/customizing_your_gateway/)
[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
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
      - JWT(利用Token进行身份认证)：与OAuth 2.0 完全不同，常用于前后端分离的情况，能够非常便捷地给API接口提供安全鉴权。用于在各方之间以JSON对象安全地传输信息，且信息是经过数字签名的。JWT令牌的内容是非严格加密的。
        JWT 以紧凑的形式由三部分组成，以点分隔
        - Header 头部通常由两部分组成，分别是令牌的类型和所使用的签名算法，用于描述其元数据，最后用base64UrlEncode算法对此JSON对象进行转换。
        - Payload 有效载荷，用于存储在JWT中实际传输的数据，包括过期时间、签发者等，可以自定义添加字段，最后用base64UrlEncode算法对此JSON对象进行转换。对于敏感信息建议不要放到JWT中，如果一定要放，则应进行不可逆加密处理，确保信息的安全性。
        - Signature 签名是对Header,Payload进行约定算法和规则的签名，用于校验消息在整个过程中有没有被篡改。生成签名：HMACSHA256(base64UrlEncode(header)+"."+base64UrlEncode(payload),secret)
        JWT可以用于下面几种场景：
        - 单点登录(SSO)：开销很小，能够轻松跨域使用。
        - 信息交换：使用公私钥对签名进行验证，确定发送者的真实身份；签名中包含了具体的内容，可以验证内容是否被篡改；
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
         #go get -u github.com/go-bindata/go-bindata/...
         go install github.com/go-bindata/go-bindata/...@v3.1
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
   
     Protobuf 是强规范的，其本身就包含字段名和字段类型等信息。
     gRPC基于HTTP/2标准设计，拥有双向流、流控、头部压缩、单TCP连接上的多路复用请求等特性，具有以下优势：
     - 使用HTTP/2来传输序列化后的二进制信息，字节数组比JSON更小，序列化与反序列化速度更快，传输速度更快。
     - 与语言、框架无关的序列化与反序列化能力。
     - 支持双向的流式传输

   - Protobuf 的使用
   
     protoc 是 protobuf 的编译器，是用 C++ 编写的，其主要功能是编译 .proto 文件。
     参照 `protoc-installation` 安装 protoc。在找不到安装的动态链接库的特定情况下，需要手动执行 `ldconfig` 命令，让动态链接库为系统所共享。也就是说，ldconfig 是一个动态链接库管理命令。
     ```shell
     # 3.15.7 to 25.1
     PROTOC_VERSION=25.1
     #wget https://github.com/protocolbuffers/protobuf/releases/download/v3.15.7/protoc-3.15.7-osx-x86_64.zip
     wget https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/protoc-$PROTOC_VERSION-osx-x86_64.zip
     unzip -d protoc-$PROTOC_VERSION-osx protoc-$PROTOC_VERSION-osx-x86_64.zip
     ln -fs protoc-$PROTOC_VERSION-osx current
     ```
     
     仅安装protoc编译器是不够的，针对不同的语音，还需要安装运行时的 protoc 插件，而对就 Go 的是 protoc-gen-go 插件、protoc-gen-go-grpc 插件。
     ```shell
     # module github.com/golang/protobuf is deprecated, use the "google.golang.org/protobuf" module instead.
     # go get -d -u -v github.com/golang/protobuf/protoc-gen-go
     # google.golang.org/protobuf=github.com/golang/protobuf
     #go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
     go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
     
     # module declares its path as: google.golang.org/grpc/cmd/protoc-gen-go-grpc, but was required as: github.com/grpc/grpc-go/cmd/protoc-gen-go-grpc
     # go get -d -u -v github.com/grpc/grpc-go/cmd/protoc-gen-go-grpc
     # google.golang.org/grpc=github.com/grpc/grpc-go
     #go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
     go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
     #export PATH="$PATH:$(go env GOPATH)/bin"
     ```

   - gRPC 的使用

     在项目下，安装 Go gRPC 库：
     ```shell
     # github.com/grpc/grpc-go v1.59.0, v1.29.1
     # google.golang.org/grpc=github.com/grpc/grpc-go
     # go get -u -v google.golang.org/grpc
     go get -u -v google.golang.org/grpc@v1.29
     ```
     在gRPC中，一共包含四种调用方式，分别是：
     1. Unary RPC 一元RPC，也就是单次RPC调用，简单来讲就是客户端发起一次普通的RPC请求、响应。
     2. Server-side streaming RPC 服务端流式RPC，是单向流，服务端流式响应，客户端为普通的一元RPC请求，简单来讲就是客户端发起一次普通的RPC请求，服务端通过流式响应多次发送数据集，客户端 Recv 接收数据集。
     3. Client-side streaming RPC 客户端流式RPC，是单向流，客户端通过流式发起多次RPC请求给服务端，服务端发起一次响应给客户端。
     4. Bidirectional streaming RPC 双向流式RPC，是双向流，由客户端以流式的方式发起请求，服务端同样以流式的方式响应请求。
        第一个请求一定是客户端发起，但具体交互方式(谁先谁后、一次发多少、响应多少、什么时候关闭)，根据程序编写的方式来确定(可结合协程)。
   
     在使用 Unary RPC时，会有如下的问题：
     1. 在一些业务场景下，数据包过大，可能会造成瞬时压力。
     2. 接收数据包时，需要所有数据包都接受成功且正确后，才能够回调响应，进行业务处理（无法客户端边发送，服务端边处理）。

     Streaming RPC 的场景：
     1. 持续且大数据包场景
     2. 实时交互场景
     
     gRPC在建立连接前，客户端/服务端都会发送连接前言(Magic+SETTINGS)，确立协议和配置项。
     gRPC在传输数据时，会涉及滑动窗口（WINDOW_UPDATE）等流控策略的。
     传播 gRPC 附加信息时，是基于 HEADERS 帧进行传播和设置；而具体的请求/响应数据是存储的 DATA 帧中的。
     gRPC 请求/响应结果会分为 HTTP 和 gRPC 状态响应（grpc-status、grpc-message）两种类型。
     客户端发起 PING，服务端就会回应 PONG，反之亦可。

   - 运行一个 gRPC 服务

     生成 proto 文件：
     ```shell
     #protoc --go_out=plugins=grpc:. ./protos/*.proto
     protoc --go_out=./protos/ --go-grpc_out=./protos/ ./protos/*.proto
     
     protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
     --go_out=./protos/ --go-grpc_out=./protos/ \
     --grpc-gateway_out=logtostderr=true:./protos ./protos/*.proto
     ```
     如果存在多层级目录的情况，可以利用 protoc 命令的 -I, M 指令来进行特定的处理。
     googleapis 是 google API 的公共接口定义，在 github 上搜索 googleapis 可找到对应的仓库。需要注意的是，由于 Go 具有不同的目录结构，很难在原始的googleapis库中存储和生成go grpc代码，所以go grpc实际上使用的是`go-genproto`库，此库有如下两个主要使用来源：
     - google/protobuf: protobuf, ptypes 子目录中的代码，均从存储库派生的，protobuf中的消息休用于描述 protobuf 本身，ptypes 下的消息体定义了常见的类型。
     - googleapis/googleapis: 专门用于与google API进行交互的类型。
     
     gRPC 是基于 HTTP/2 协议的，不能直接通过 postman 或普通的 curl 进行调用，目前开源社区的方案：命令行工具 grpcurl，安装及使用命令如下：
     ```shell
     go get -u -v github.com/fullstorydev/grpcurl
     go install github.com/fullstorydev/grpcurl/cmd/grpcurl
     
     # 默认使用 TLS 认证(-cert,-key 设置公钥和密钥)，-plaintext 用来忽略TLS认证
     grpcurl -plaintext localhost:8001 list
     grpcurl -plaintext localhost:8001 list proto.TagService
     grpcurl -plaintext -d '{"name": "Go"}' localhost:8001 list proto.TagService.GetTagList
     ```
     但使用此工具的前提是gRPC Server 已经注册了反射服务：`s := grpc.NewServer() reflection.Register(s)`
   
     > gRPC Server/Client 在启动和调用时，必须明确其是否加密，`DialOpton grpc.WithInsecure`指定为非加密模式。
     > grpc(HTTP/2)和HTTP/1.1通过Header中的Content-Type和ProtoMajor标识进行分流 

   - gRPC 服务间的内调
   - 提供 HTTP 接口

     grpc协议的本质是HTTP/2协议，如果服务需要在同端口适配两种协议流量，则需要进行特殊处理。
     - 不同的两个端口，监听不同协议的流量：使用两个协程分别监听 http endpoint, grpc endpoint 实际是一个阻塞行为。
     - 同端口上兼容多种协议：使用第三方开源库cmux来实现对多协议的支持。
     
        cmux根据有效负载(payload)对连接进行多路复用，即只匹配连接的前面的字节来区分当前连接的类型，可以在同一tcp listener上提供grpc,ssh,https,http,go rpc等几乎所有其他协议的服务，是一个相对通用的方案。
        需要注意，一个连接可以是grpc或http，但不能同时是两者。
        ```shell
        go get -u github.com/soheilhy/cmux@v0.1.5
        ```
        grpc(http/2),http/1.1在网络分层上都是基于tcp协议的，拆分为tcp,grpc,http逻辑，便于连接多路复用器。cmux基于content-type头信息标识进行分流，grpc特定标识：application/grpc。

     - 同端口同方法双流量支持：应用代理 grpc-gateway 能够将 Restful 转换为 gRPC 请求，实现用同一个RPC方法提供对gRPC协议和HTTP/1.1的双流量支持。

       开源社区的 grpc-gateway 是 protoc 的一个插件，能够读取 protobuf 的服务定义，生成一个反向代理服务器，将 Restful JSON API 转换为 gRPC。它主要是根据 protobuf 的服务定义中的 google.api.http 来生成的。
       grpc-gateway的功能就是生成一个HTTP的代理服务，将HTTP请求转换为gRPC协议，并转发到gRPC服务器中。
       简单来说，grpc-gateway 能够将 Restful 转换为 gRPC 请求，实现用同一个RPC方法提供对gRPC协议和HTTP/1.1的双流量支持。
       ```shell
       # https://grpc-ecosystem.github.io/grpc-gateway/
       # github.com/grpc-ecosystem/grpc-gateway
       go install \
       github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15 \
       github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15 \
       google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 \
       google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
     
       # -I 参数的格式：-IPATH,--proto_path=PATH，用来指定 proto 文件中 import 搜索的目录，可指定多个。如果不指定，则默认是当前工作目录。
       # Mfoo/bar.proto=quux/shme，则在生成、编译proto时，将指定的包名替换为要求的名字，例子中将把foo/bar.proto编译后的包名替换为quux/shme。
       protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:./protos/ ./protos/*.proto
       ```
       google/api/annotations.proto 文件，是 googleapis 的产物。使用grpc-gateway下的annotations.proto，可以保证兼容性和稳定性。
       http.proto HttpRule 对HTTP转换提供支持，可用于定义API服务的HTTP的相关配置，并可以指定每一个RPC方法都映射到一个或多个HTTP REST API方法上。
       如果没有引入annotations.proto文件，且填写了相关的HTTP Option，则执行生成命令后，虽然不会报错，但不会生成任何相关内容。
       
    - 其他方案
      - 外部网关组件
        - envoy gRPC-JSON transcoder
        - apache APISIX: etcd grpc proxy(A stateless etcd reverse proxy operating at the gRPC layer)

   - 接口文档
   
     proto 的插件 protoc-gen-swagger，作用是通过 proto 文件生成 swagger 定义(.swagger.json)：
     ```shell
     # github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
     go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest
     ```

     proto 的插件 protoc-gen-openapiv2，作用是通过 proto 文件生成 OpenAPI 定义(.swagger.json):
     ```shell
     go install \
     github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2
     
     # generating OpenAPI definitions for unannotated methods, use the `generate_unbound_methods` option
     # --openapiv2_out=generate_unbound_methods=true;./gen/openapiv2 
     # or --openapiv2_out ./gen/openapiv2 --openapiv2_opt generate_unbound_methods=true
     protoc -I. --openapiv2_out ./gen/openapiv2 --openapiv2_opt generate_unbound_methods=true \
     ./protos/tag.proto
     ```
     
     Swagger 提供了可视化的接口管理平台-[Swagger UI](https://swagger.io/tools/swagger-ui/)。从其管理平台下载源码压缩包，将其目录下的dist目录下的所有资源文件拷贝到项目的 third_part/swagger-ui 目录中。
     使用 go-bindata 库将资源应文件转换为 Go 代码，转换命令：`go-bindata --nocompress -pkg swagger -o pkg/swagger/data.go third_party/swagger-ui/...`，命令自动在目录 pkg/swagger 下创建 data.go 文件。
     为了让转换的静态资源代码能够被外部访问，需安装 go-bindata-assetfs 库，它能够结合 net/http、go-bindata 库生成 swagger-ui 的go代码供外部访问：
     ```shell
     go get -u github.com/elazarl/go-bindata-assetfs/...
     ```
     在 HTTP Server 中添加AssetFS相关的http.FileServer相关的处理逻辑：
     ```shell
     fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:     swagger.Asset,
		AssetDir:  swagger.AssetDir,
		Prefix: "third_party/swagger-ui",
     })
     serveMux.Handle(prefix, http.StripPrefix(prefix, fileServer))
     ```
     
     基于 proto 文件生成 swagger 定义文件 .swagger.json:
     ```shell
     # protos 目录下生成 .swagger.json 定义文件
     protoc -I. -I$GOPATH/src \
     -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
     --swagger_out=logtostderr=true:. ./protos/*.proto
     ```
     在实际环境中，让每个服务仅提供swagger定义，然后在统一的平台上提供swagger站点来读取swagger定义，这样就不需要每个服务都运行swagger站点了，同时由于入口是统一的，所以鉴权也能在这个基础上完成。

   - gRPC 拦截器
   
     在每个RPC方法的前面或后面做统一的特殊处理，如鉴权校验、上下文的超时控制、请求的日志记录等，使用拦截器(Interceptor)定制，不直接侵入业务代码。
     一种类型的拦截器只允许设置一个。gpc-go issues #935明确得知：官方仅提供了一个拦截器的钩子，以便在其中构建各种复杂的拦截器模式，而不会遇到多个拦截器的执行顺序问题，同时还能保持grpc-go自身的简介性，尽可能最小化公共API。
     将不同的功能设计为不同的拦截器：
     - 自己实现一套多拦截器的逻辑(拦截器中调用拦截器)
     - 直接使用grpc应用生态(grpc-ecosystem)中的go-grpc-middleware提供的grpc.UnaryInterceptor,grpc.StreamInterceptor， 在grpc.*Interceptor中嵌套grpc_middleware.ChainUnaryServer或grpc_middleware.ChainUnaryClient(拦截器数量大于1时，每个递归的拦截器都会不断地执行，最后才去真正执行代表RPC方法的handler)，以链式方式达到用多个拦截器的目的。
       ```
       // 服务端拦截器
       // grpc.ServerOption 设置Server的相关属性，如credentials,keepalive等参数。
       // 拦截器在此注册，需以指定的类型封装，如一元拦截器的类型必须为grpc.UnaryInterceptor。
       // 在grpc.UnaryInterceptor中嵌套grpc_middleware.ChainUnaryServer，以链式方式达到用多个拦截器的目的。
       opts := []grpc.ServerOption{
           grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
               middleware.AccessLog,
               middleware.ErrorLog,
               middleware.Recovery,
               middleware.ServerTracing,
           )),
       }
       s := grpc.NewServer(opts...)
       
       // 客户端拦截器的相关注册行为是在调用grpc.Dial或grpc.DialContext之前，通过DialOption配置选项进行注册的。
       // 超时控制(对上下文超时时间的设置和适当控制)，是在微服务架构中非常重要的一个保命项。
       // 当在服务调用中没有设置超时时间，或设置的超时时间过长时，就会导致多服务下的整个调用链雪崩(响应缓慢)。
       // grpc中建议始终设置截止日期(context.Deadline方法检查，若未设置截止时间，则返回false；context.WithTimeout方法设置默认超时时间)。
       var opts []grpc.DialOption
       opts = append(opts, grpc.WithUnaryInterceptor(
           grpc_middleware.ChainUnaryClient(
               middleware.UnaryCtxTimeout(),
               middleware.ClientTracing(),
           ),
       ))
       
       opts = append(opts, grpc.WithStreamInterceptor(
           grpc_middleware.ChainStreamClient(
               middleware.UnaryCtxTimeout(),
               middleware.ClientTracing(),
           ),
       ))

       clientConn, err := grpc.DialContext(ctx, target, opts...)
       ```   
       关于链式拦截器上，也就是多拦截器的使用，推荐go-grpc-middleware的方案，不过从`grpc v1.28.0`起，有社区朋友贡献并合并了链式拦截器的相关方法（参见 issues #935）。

     由于客户端和服务端有各自的一元拦截器和流拦截器，因此，在gRPC中，共有四种类型的拦截器：
     - 一元拦截器 Unary Interceptor：拦截和处理一元RPC调用
     - 流拦截器 Stream Interceptor: 拦截和处理流式RPC调用
     - 客户端
       - 一元拦截器：类型为UnaryClientInterceptor，实现通常分为三部分：预处理、调用RPC方法和后处理。
       
       ```
       type UnaryClientInterceptor func(
           ctx context.Context,
           method string,
           req,
           reply interface{},
           cc *ClientConn, // 客户端连接句柄
           invoker UnaryInvoker, // 所调用的RPC方法
           opts ...CallOption, // 调用的配置
       ) error
       ```
       
       - 流拦截器：类型为StreamClientInterceptor，实现包括预处理和流操作拦截两部分，不能在事后进行RPC方法调用和后处理，只能拦截用户对流的操作。
       
       ```
       type StreamClientInterceptor func(
           ctx context.Context,
           desc *StreamDesc,
           cc *ClientConn, // 客户端连接句柄
           method string, 
           streamer Streamer,
           opts ...CallOption, // 调用的配置
       ) (ClientStream, error)
       ```

     - 服务端
       - 一元拦截器：类型为UnaryServerInterceptor。
       
       ```
       type UnaryClientInterceptor func(
           ctx context.Context,
           req, 
           info *UnaryServerInfo, // RPC方法的所有信息
           handler UnaryHandler, // RPC方法本身
       ) (resp interface{}, err error)
       ```

       - 流拦截器
       
       ```
       type UnaryClientInterceptor func(
           srv interface{}, 
           ss ServerStream,
           info *StreamServerInfo, // RPC方法的所有信息
           handler StreamHandler, // RPC方法本身
       ) error
       ```

   - metadata 和 RPC 自定义认证
   
     metadata:
     在HTTP/1.1中，通常直接操纵Header来传递数据，而对于gRPC来讲，基于HTTP/2，本质上也可以通过Header来进行传递，但不会直接的去操纵它，而是通过gRPC中的metadata进行调用过程中的数据传递和操纵。
     需要注意，使用metadata的前提，需要使用的库进行支持。
     在gRPC中，metadata实际上是一个map结构，一个字符串与字符串切片的映射结构：`type MD map[string][]string`。
     `metadata.New`方法创建的metadata会直接被转换为对应的MD结构；`metadata.Pairs`方法创建的metadata会以奇数来配对，并且所有的Key默认都被转为小写，若有同名Key，将会追加到对应 Key 的切片（slice）上。
     在gRPC中，metadata可分为传入用的metadata和传出用的metadata两种，为了防止metadata从入站RPC直接转发到其出站RPC的情况(issues #1148)，并提供了两种方法分别进行设置：
     - NewIncomingContext：创建一个附加了传入 metadata 的新上下文，仅供自身的 gRPC 服务端内部使用(IncomingContext)。
     - NewOutgoingContext：创建一个附加了传出 metadata 的新上下文，可供外部的 gRPC 客户端、服务端使用(OutgoingContext)。
     在metadata获取上，也分为两个方法，分别是 FromIncomingContext 和 FromOutgoingContext。
     在内部进行了Key的区分，用指定的Key读取相应的metadata，以防止造成脏读。`推荐对Key的设置，使用一个结构体去定义`。
     注意：新增metadata信息时，务必使用Append类别的方法(e.g. AppendToOutgoingContext)，若直接创建一个全新的metadata，则会导致原有的metadata信息丢失。
     
     自定义认证：在实际场景中，对某些模块的RPC方法，做特殊认证或校验，可以使用gRPC Token 接口：`PerRPCCredentials`
     gRPC PerRPCCredentials，是用于自定义认证Token的默认接口，作用是将所需的安全认证信息添加到每个RPC方法的上下文中。客户端注册：在 DialOption 配置中调用 grpc.WithPerRPCCredentials 方法。服务端：调用 metadata.FromIncomingContext 从上下文中获取 metadata，再在不同的 RPC 方法中进行认证检查。
     实际上，metadata在应用传输上做了严格的进出隔离，即在在上下文中分隔传入和传出的metadata。在使用metadata时需要多思考一下，到底应该是出还是入，以此调用不同的处理方法。

   - 链路追踪
   
     链路追踪通常会涉及多个服务，链路信息会更多，因此精准的链路信息是非常有帮助的。
     做链路追踪的基本条件是要注入追踪信息，而最简单的方法就是使用服务端和客户端拦截器组成完整的链路信息：
     - 客户端拦截器：从metadata中提取链路信息，将其设置并追加到服务端的调用上下文中。若本次调用没有上一级链路信息，则生成对应的父级信息，自己成为父级；若本次调用存在上一级链路信息，则会根据上一级链路信息进行设置，成为其子级。
     - 服务端拦截器：从调用的上下文中提取链路信息，并将其作为metadata追加到RPC调用中。
     在opentracing中，可以指定SpanContexts的三种传输表示方式：Binary(不透明的二进制数据),TextMap(键值字符串对),HTTPHeaders(HTTP Header格式的字符串)

   - gRPC 服务注册和发现
   
     多个服务频繁发生网络位置变更和调整的情况，服务端并不会固定在某一个网络环境中，客户端需要动态获取实例地址，然后再次进行请求。
     服务注册和发现，本质上就是增加多种角色对服务信息进行协调处理，常见的角色如下：
     - 注册中心：承担对服务信息进行注册、协调、管理等工作。
     - 服务提供者：暴露特定端口，并提供一个到多个的服务允许外部访问。
     - 服务消费者：服务提供者的调用方。
     假设服务注册模式为自行上报，则服务提供者在启动服务时，会将自己的服务信息(如IP地址、端口号、版本号等)注册到注册中心。服务消费者在进行调用时，会以约定的命名标识(如服务名)到注册中心查询，发现当前那些具体的服务可以调用。注册中心再根据约定的负载均衡算法进行调度，最终请求到服务提供者。
     另外，当服务提供者出现问题时，或当定期的心跳检测发现服务提供者无正确响应时，那么这个出现问题的服务就会被下线，并标识为不可用。即在启动时上报注册中心进行注册，把被检测到出问题的服务下线，以此维护服务注册和发现的一个基本模型。
     以上是一个基本思路，除此之外，还有其他多种实现的思路和方案。
   
     通常，gRPC的负载均衡是基于每个调用去进行的负载，会在每次调用时根据具体的负载均衡策略，选择最优的服务进行访问，而不是基于每个连接来进行的。简单来说，即使所有的调用都来自同一个gRPC客户端，仍希望是在所有服务端中实现负载均衡，即把所有的调用平均地分配在各个服务端上，而不是固定在某一个服务端上。
   
     常见的负载均衡类型：
     - 客户端负载
       调用时，由客户端到注册中心对服务提供者进行查询，并获取所需的服务清单(包含各个服务的实际信息，如IP、端口号、集群命名空间等)，由客户端使用特定的负载均衡策略(如轮询)，在服务清单中选择一个或多个服务进行调用。
       优点：高性能、去中心化，并且不需要借助独立的外部负载均衡组件。缺点：实现成本高，并且可能需要针对当时获得的服务清单进行过期反补处理等。
     - 服务端负载
       在服务端搭设独立的负载均衡器，负载均衡器再根据给定的目标名称(如服务名)，找到适合调用的服务实例，具备负载均衡和反向代理的功能。
       优点：简单、透明，客户端不需要知道背后的逻辑，只需按给定的目标名称调用、访问就行，由服务端侧管理负载、均衡策略及代理。缺点：外部的负载均衡器可能成为性能瓶颈，可能会出现更高的网络延迟，同时必须保持高可用。
     
     gRPC官方基本设计思路(client: NameResolver(e.g. DNS)+grpclb policy, LB, Server)
     - 在进行gRPC调用时，gRPC客户端会向名称解析器(Name Resolver)发出服务端名称(服务名)的名称解析请求。名称解析器会将服务名解析成一个或多个IP地址，每个IP地址都会标识它是服务端地址，还是负载均衡器地址，及客户端使用的负载均衡策略(如round robin or grpclb)。
     - 客户端实例化负载均衡策略，如果gRPC客户端获取的地址是负载均衡器地址，那么客户端将使用grpclb策略，否则使用服务配置请求的负载均衡策略。如果服务配置未请求负载均衡策略，则客户端默认选择第一个可用的服务端地址。
     - 负载均衡策略会为每个服务端地址创建一个grpc服务端。对于每个请求，都由负载均衡策略决定将其发送到那个grpc服务端。
     设计思路的核心如下：
     - 客户端根据服务名发起请求
     - 名称解析器解析服务名并返回
     - 客户端根据服务端类型选择相应的策略
     - 最后根据不同的策略进行实际调用
     
     实际上，gRPC官方提供了两个接口(Resolver,Watcher)用于对外部组件进行扩展，其中Resolver接口的主要作用是解析目标创建的观察程序，跟踪其信息的更改，Watcher接口的主要作用是监视指定目标上的变更。
     目前可实现负载均衡的技术方案有很多，如在客户端负载上，可基于consul实现，而在服务端负载上，可基于nginx(1.13起支持grpc)、kubernetes、istio、linkerd等外部组件实现类似的负载均衡，或服务发现或注册的相关功能。

   - 实现自定义的 protoc 插件

     在开发grpc+protobuf的相关服务时，用到了许多protoc相关的插件，来实现各种功能：
     |插件名称|对应的命令|
     |---|---|
     |protoc-gen-go|--go_out|
     |protoc-gen-grpc-gateway|--grpc-gateway_out|
     |protoc-gen-swagger|--swagger_out|
   
     在查看grpc插件时，发现其实现了很多方法，但只需要重点关注`generator.Plugin`的相关接口。在实现自定义插件时，只需满足Plugin接口的定义，就可以无缝地接入protoc。grpc插件是根据Plugin接口的主体流程来流转的。

   - 对 gRPC 接口进行版本管理

     gRPC接口是以protobuf作为IDL的，如何改动才能兼容使用先前版本的客户端，且影响其他客户端？ 又如何更好地管理gRP接口版本？
     1. 可兼容性修改
        - 新增新的service
        - 在原service中新增新的rpc方法
        - 在原请求message中新增新的字段
        - 在原响应message中新增新的字段
     2. 破坏性修改
        - 修改原有的字段数据类型
        - 修改原有的字段标识位：字段标识是用来标识其在protobuf payload中的位置的，一旦发生改变，就会出现找不到或错位的情况
        - 修改原有的字段名
        - 修改message原有的命名
        - 删除service或RPC方法
        - 删除原有的字段
     3. 如何设计gRPC接口
        在变更gRPC接口时，要尽可能地不影响现有的客户端。
        在研判需求、开接口设计会议时，应把变更锁定在“可兼容性修改”上。如果遇到非常大的接口出入参变动，则往往会选择两种方式：
        1. 编写新的RPC方法
        2. 在原接口内对入参或出参进行转换，然后在内部将流量导向新的RPC方法
     4. 版本号管理
        1. 在package上指定版本号，格式：package proto.v1
        2. 若存在HTTP路由，则可以将大版本号定义在路由中，格式：/api/v1/tag
        3. 通过Header进行传播，但在gRPC下，并不够显性化，并且可能需要进行额外的特殊处理，需要根据实际场景判定。

   - 常见问题讨论

4. chatroom IM聊天室

   - 聊天室需求分析和设计
   - 核心流程
   - 广播器
   - @提醒功能
   - 敏感词处理
   - 离线消息处理
   - 关键性能分析和优化

5. cache-example 进程内缓存

## [go-micro](https://github.com/asim/go-micro) a Go microservices framework

see https://github.com/go-micro for tooling.
- [go-micro 样例](https://github.com/go-micro/examples)

go-micro 拥有更丰富的生态和功能，更方便的工具和API。如服务注册可以方便地切换到etcd,zk,gossip,nats等注册中心，方便实现服务注册功能。Server支持gRPC,HTTP等协议。
使用`protoc-gen-micro`插件来生成micro适用的协议文件，但插件版本需和go-micro版本相同。
```text
go install github.com/asim/go-micro/cmd/protoc-gen-micro/v4@latest

# 手动下载依赖文件并放入GOTPATH下
git clone git@github.com:googleapis/googleapis.git
mv googleapis/google $(go env GOPATH)/src/github.com/googleapis/google
# 或者使用 grpc-ecosystem/grpc-gateway v1.14.5
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.14.5

# go-micro,grpc-gateway插件，可能存在命名冲突，指定grpc-gateway的选项register_func_suffix，让生成的函数名包含后缀，解决命令冲突。
protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--grpc-gateway_out=logtostderr=true,register_func_suffix=Gw:./protos/ \
--micro_out=./protos/ --go_out=./protos/ --go-grpc_out=./protos/ \
./protos/*.proto

# gateway when http method is DELETE except allow_delete_body option is true
protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--grpc-gateway_out=logtostderr=true,allow_delete_body=true,register_func_suffix=Gw:./protos/ \
--go_out=./protos/ --go-grpc_out=./protos/ --micro_out=./protos/ \
./protos/crawler/crawler.proto
```

为了生成gRPC服务器，需要导入go-micro的gRP插件库(github.com/go-micro/plugins/v4/server/grpc)，生成一个gRPC Server注入micro.NewService 中。