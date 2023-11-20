# grpc demo

## 根据 .proto 文件生成 Go Code

```shell
#protoc --go_out=plugins=grpc:./protos/ ./protos/*proto
protoc --go_out=./protos/ --go_opt=Mprotos/helloworld.proto=github.com/dapings/examples/grpc-demo/helloworld \
--go-grpc_out=./protos/ --go-grpc_opt=Mprotos/helloworld.proto=github.com/dapings/examples/grpc-demo/helloworld \
protos/*.proto

#proto 文件中添加 option go_package
protoc --go_out=./protos/ --go-grpc_out=./protos/ protos/helloworld.proto
```

The Go import path format: --go_opt=M${PROTO_FILE}=${GO_IMPORT_PATH}, for example: --go_opt=Mprotos/bar.proto=example.com/protos/foo;package_name

## 运行 Server

```shell
go run server/server.go
```

## 运行 Client

```shell
go run client/client.go
greeting: hello dapings world
```