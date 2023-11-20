# grpc demo

```shell
#protoc --go_out=plugins=grpc:./protos/ ./protos/*proto
protoc --go_out=./protos/ --go_opt=Mprotos/helloworld.proto=github.com/dapings/examples/grpc-demo/helloworld \
--go-grpc_out=./protos/ --go-grpc_opt=Mprotos/helloworld.proto=github.com/dapings/examples/grpc-demo/helloworld \
protos/*.proto

#proto 文件中添加 option go_package
protoc --go_out=./protos/ --go-grpc_out=./protos/ protos/helloworld.proto
```

The Go import path format: --go_opt=M${PROTO_FILE}=${GO_IMPORT_PATH}, for example: --go_opt=Mprotos/bar.proto=example.com/protos/foo;package_name