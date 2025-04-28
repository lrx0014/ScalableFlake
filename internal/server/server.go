package server

import (
	pb "github.com/lrx0014/ScalableFlake/api/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

func RunServer() {
	log.Infof("starting server [grpc, http]")
	go func() {
		lis, _ := net.Listen("tcp", ":9000")
		grpcServer := grpc.NewServer()
		pb.RegisterUIDGeneratorServer(grpcServer, NewGRPCServer())
		_ = grpcServer.Serve(lis)
	}()

	httpServer := NewHTTPServer()
	_ = http.ListenAndServe(":8000", httpServer)
}
