# 一个最简单的自定义protoc插件

1. 写入固定的公共generator逻辑
2. tour 插件实现generator.Plugin接口
3. link_tour.go 文件，import 中初始化tour插件
4. 编译二进制文件：go build .
5. 根据 proto 文件，生成Go文件
    ```shell
    protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --tour_out=plugins=tour:./protos/ ./protos/*.proto
    ```
   --tour_out=plugins=tour，--tour_out会告诉protoc编译器去查找并使用名为protoc-gen-tour的插件，而plugins=tour则指定使用protoc-gen-tour插件中的tour子插件(允许在插件中自定义多个子插件)。