package server

import (
	pb "github.com/lrx0014/ScalableFlake/api/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func RunServer() {
	log.Infof("starting server [grpc, http]")
	go func() {
		lis, _ := net.Listen("tcp", ":9000")
		grpcServer := grpc.NewServer()
		pb.RegisterUIDGeneratorServer(grpcServer, NewGRPCServer())
		_ = grpcServer.Serve(lis)
	}()

	_ = NewHTTPServer().Run(":8000")
}
