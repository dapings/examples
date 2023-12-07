module github.com/dapings/examples/go-programing-tour-2023/tag-service

go 1.20

require (
	github.com/elazarl/go-bindata-assetfs v1.0.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.6 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.0.0-beta.3
	github.com/opentracing/opentracing-go v1.2.0
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	golang.org/x/exp v0.0.0-20231108232855-2478ac86f678
	golang.org/x/net v0.19.0
	// v1.29.1
	google.golang.org/grpc v1.29.1
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.3.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/golang/glog v1.1.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/soheilhy/cmux v0.1.5
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	//go.etcd.io/etcd/client/v3 v3.5.10 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto v0.0.0-20200521103424-e9a78aa275b7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require github.com/coreos/etcd v3.3.27+incompatible

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coreos/bbolt v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/pkg v0.0.0-20230601102743-20bbbf26f4d8 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/uuid v1.4.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/prometheus/client_golang v1.17.0 // indirect
	github.com/prometheus/client_model v0.4.1-0.20230718164431-9a2bf3000d16 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20220101234140-673ab2c3ae75 // indirect
	github.com/xiang90/probing v0.0.0-20221125231312-a49e3df8f510 // indirect
	go.etcd.io/bbolt v0.0.0-00010101000000-000000000000 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.18.1 // indirect
	golang.org/x/crypto v0.16.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

replace (
	//	two different module paths (github.com/coreos/bbolt and go.etcd.io/bbolt)
	// replace 两次，具体原因待排查。
	// 410 Gone，通常是Go版本不同导致的，1.13以上版本，通过GOSUMDB的环境变量调整：export GOSUMDB=off; go mod download
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.8
	go.etcd.io/bbolt => github.com/coreos/bbolt v1.3.8
)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.3.2
