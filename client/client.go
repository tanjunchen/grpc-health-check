package main

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/tanjunchen/grpc-health/proto"
	"google.golang.org/grpc"
)

func main() {
	serverAddr := ":8989"
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
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
