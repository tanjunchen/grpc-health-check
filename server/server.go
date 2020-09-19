package main

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/tanjunchen/grpc-health/proto"
	"github.com/tanjunchen/grpc-health/server/healthcheck"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type server struct{}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
}

func (s *server) Hello(helloReq *proto.HelloRequest, srv proto.HelloService_HelloServer) error {
	logrus.Infof("Server received an rpc request with the following parameter %v", helloReq.Hello)
	for i := 0; i <= 10; i++ {
		resp := &proto.HelloResponse{
			Greet: fmt.Sprintf("Hello %s for %d time", helloReq.Hello, i),
		}
		srv.SendMsg(resp)
	}
	return nil
}

func main() {
	serverAdr := ":8989"
	listenAddr, err := net.Listen("tcp", serverAdr)
	if err != nil {
		logrus.Fatalf("Error while starting the listening service %v", err.Error())
	}
	grpcServer := grpc.NewServer()
	proto.RegisterHelloServiceServer(grpcServer, &server{})
	healthService := healthcheck.NewHealthChecker()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthService)
	logrus.Infof("Server starting to listen on %s", serverAdr)
	if err = grpcServer.Serve(listenAddr); err != nil {
		logrus.Fatalf("Error while starting the gRPC server on the %s listen address %v", listenAddr, err.Error())
	}
}
