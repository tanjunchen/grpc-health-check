package main

import (
	"net"

	"github.com/sirupsen/logrus"
	"github.com/tanjunchen/grpc-health/controller"
	"github.com/tanjunchen/grpc-health/middleware"
	"github.com/tanjunchen/grpc-health/router"
	"github.com/tanjunchen/grpc-health/server/healthcheck"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
}

func main() {
	serverAdr := "0.0.0.0:8989"
	listenAddr, err := net.Listen("tcp", serverAdr)
	if err != nil {
		logrus.Fatalf("Error while starting the listening service %v", err.Error())
	}
	var opts []grpc.ServerOption

	opts = append(opts, grpc.UnaryInterceptor(middleware.Interceptor))

	grpcServer := grpc.NewServer(opts...)

	router.Init(grpcServer)

	reflection.Register(grpcServer)

	healthService := healthcheck.NewHealthChecker()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthService)

	stopCh := signals.SetupSignalHandler()
	go controller.StartWatchController(stopCh)

	logrus.Infof("Server starting to listen on %s", serverAdr)
	if err = grpcServer.Serve(listenAddr); err != nil {
		logrus.Fatalf("Error while starting the gRPC server on the %s listen address %v", listenAddr, err.Error())
	}
}
