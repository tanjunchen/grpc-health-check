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
	// serverAddr := "10.20.11.116:30380"
	serverAddr := "0.0.0.0:8989"
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

	WatchService(conn)

	logrus.Println("Doing a health check on the server")
}

func Hello(conn *grpc.ClientConn) {
	helloClient := proto.NewHelloServiceClient(conn)
	stream, _ := helloClient.Hello(context.Background(), &proto.HelloRequest{
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
}

func WatchService(conn *grpc.ClientConn) {
	client := proto.NewServiceServiceClient(conn)
	stream, _ := client.SyncServiceWatchListService(context.TODO())
	for {
		// 接收从 服务端返回的数据流
		response, err := stream.Recv()
		if err != nil {
			fmt.Println("接收数据出错:", err)
			break
		}

		if response != nil {
			// 没有错误的情况下，打印来自服务端的消息
			fmt.Printf("[客户端收到]: %s \n", response)
		} else {
			fmt.Printf("[客户端收到]: %s \n", response)
			break
		}
	}
}
