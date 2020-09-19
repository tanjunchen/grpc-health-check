package router

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tanjunchen/grpc-health/proto"
	"google.golang.org/grpc"
)

func Init(grpcServer *grpc.Server) {
	proto.RegisterHelloServiceServer(grpcServer, &server{})
}
func (s *server) Hello(helloReq *proto.HelloRequest, srv proto.HelloService_HelloServer) error {
	logrus.Infof("Server received an rpc request with the following parameter %v", helloReq.Hello)
	for i := 0; i <= 10; i++ {
		resp := &proto.HelloResponse{
			Greet: fmt.Sprintf("Hello %s for %d time", helloReq.Hello, i),
		}
		fmt.Println(resp.Greet)
		srv.SendMsg(resp)
	}
	return nil
}

type server struct{}
