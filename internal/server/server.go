package server

import (
	pb "github.com/lrx0014/ScalableFlake/api/v1"
	allocator "github.com/lrx0014/ScalableFlake/pkg/machine"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func RunServer(allocator allocator.Allocator) {
	go func() {
		log.Println("Starting grpc server")
		lis, _ := net.Listen("tcp", ":9000")
		grpcServer := grpc.NewServer()
		pb.RegisterUIDGeneratorServer(grpcServer, NewGRPCServer(allocator))
		_ = grpcServer.Serve(lis)
	}()

	log.Println("Starting http server")
	httpServer := NewHTTPServer(allocator)
	_ = http.ListenAndServe(":8000", httpServer)
}
