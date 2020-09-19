package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/xid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// auth 验证Token
func auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "No Token")
	}

	var (
		appid  string
		appkey string
	)

	if val, ok := md["appid"]; ok {
		appid = val[0]
	}

	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}

	if appid != "admin" || appkey != "admin" {
		return status.Errorf(codes.Unauthenticated, "Token invalid: appid=%s, appkey=%s", appid, appkey)
	}

	return nil
}

// interceptor 拦截器
func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	p, _ := peer.FromContext(ctx)

	// 是否检查 secret 密钥
	err := auth(ctx)
	if err != nil {
		return nil, err
	}

	m, err := handler(ctx, req)
	if err != nil {
		fmt.Printf("failed to handler Unary RPC: %v", err)
	}

	fmt.Sprintf(":: uuid=[%v], source=[%v], method=[%s], request=[%v], duration=[%s], error=[%v];", xid.New().String(), p.Addr, info.FullMethod, req, time.Since(start), err)
	return m, err
}
