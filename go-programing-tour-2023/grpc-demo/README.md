# grpc demo

```shell
#protoc --go_out=plugins=grpc:./protos/ ./protos/*proto
protoc --go_out=./protos/ --go_opt=Mprotos/helloworld.proto=. \
protos/*.proto
```

The Go import path format: --go_opt=M${PROTO_FILE}=${GO_IMPORT_PATH}, for example: --go_opt=Mprotos/bar.proto=example.com/protos/foo;package_name