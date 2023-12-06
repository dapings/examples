# 自定义protoc插件

## 一个最简单的自定义protoc插件

1. 写入固定的公共generator逻辑
2. tour 插件实现generator.Plugin接口
3. link_tour.go 文件，import 中初始化tour插件
4. 编译二进制文件：go build .
5. 根据 proto 文件，生成Go文件
    ```shell
    protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --tour_out=plugins=tour:./protos/ ./protos/*.proto
    ```
   --tour_out=plugins=tour，--tour_out会告诉protoc编译器去查找并使用名为protoc-gen-tour的插件，而plugins=tour则指定使用protoc-gen-tour插件中的tour子插件(允许在插件中自定义多个子插件)。

## 实现定制化的gRPC自定义插件

需要自定义插件的情况：
1. 基于gRPC插件实现一些定制化功能
2. 基于proto文件生成一些扩展性的代码和功能(如在社区中很常见的protoc-gen-grpc-gateway,protoc-gen-swagger)

基于官方protoc gRPC插件，实现一个最简单的定制化需求：
1. 确认需求

   在业务开发中，有时会有多租户的概念，即根据租户标识的不同，获取其对应的租户实例信息和数据，如根据租户标识，判定当前部署的环境，以便进行更多的精确调度。
   因此须在调用中传播租户标识。

2. 解决方法

   方法之一，开发自定义插件，在调用gRPC client时，必须传入租户标识，并可以在内部进行入参校验和节点发现。此方案能够较好地解决这个问题，并对开发人员的直接侵入性较小。

3. 实现插件

   - 拷贝插件模板：golang/protobuf，进入 protoc-gen-go 目录，打开grpc目录中的grpc.go文件，调整包名、原本的gRPC结构体及其方法，并修改插件的名称。
   - 进行二次开发：对gRPC client方法进行租户标识(orgcode)的获取和判断。若不存在，则直接返回相应的错误信息。
     因此，需要编写获取和设置租户标识值的方法，编写判定租户标识值是否正确的方法。

小结：
上下文设置的方法，实际上，侵入了client interface。而当拥有顶级公共库时，可以将这类方法抽到公共库中，然后直接调用。这样既不会影响client interface的默认定义，只保留核心的调度逻辑在自定义插件的生成逻辑中，又能实现相对的解耦。
另外，还可以通过初始化的client直接调用该方法，不需要再引用另外一个包。