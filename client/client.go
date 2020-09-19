package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tanjunchen/grpc-health/proto"
	"google.golang.org/grpc"
)

// customCredential 自定义认证
type customCredential struct{}

func (customCredential customCredential) RequireTransportSecurity() bool {
	return false
}

// GetRequestMetadata 实现自定义认证接口
func (customCredential customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid":  "admin",
		"appkey": "admin",
	}, nil
}

// interceptor 客户端拦截器
func interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	fmt.Printf("method=%s req=%v rep=%v duration=%s error=%v\n", method, req, reply, time.Since(start), err)
	return err
}

func main() {
	//serverAddr := "10.20.11.116:30380"
	serverAddr := ":8989"
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))

	// 指定客户端 interceptor
	opts = append(opts, grpc.WithUnaryInterceptor(interceptor))

	conn, err := grpc.Dial(serverAddr, opts...)

	if err != nil {
		logrus.Fatalf("Couldn't dial server at %s", serverAddr)
	}
	defer conn.Close()
	helloClient := proto.NewHelloServiceClient(conn)

	stream, err := helloClient.Hello(context.Background(), &proto.HelloRequest{
		Hello: "World",
	})

	for {
		streamData, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Fatalf("%v.Greet = _, %v", helloClient, err)
		}
		logrus.Println(streamData)
	}

	logrus.Println("Doing a health check on the server")
}
