package rpc

import (
	"context"

	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
)

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, GetGRPCDialOptionWithInsecure())
	// DialContext 方法是异步建立连接的，并不会马上就成为可用连接，仅处于 connecting 状态，只能正式达到ready状态，才算是可用的。
	// 创建给定目标的客户端连接，但不是马上建立了可用连接
	// 如果需要立刻打通与服务端的连接，需要设置WithBlock DialOption，这样当发起连接时会阻塞等待连接完成，使最终连接到达Ready状态
	return grpc.DialContext(ctx, target, opts...)
}

func GetGRPCDialOptionWithInsecure() grpc.DialOption {
	return grpc.WithInsecure()
	// return grpc.WithTransportCredentials(insecure.NewCredentials())
}
