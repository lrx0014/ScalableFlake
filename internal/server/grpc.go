package server

import (
	"context"
	pb "github.com/lrx0014/ScalableFlake/api/v1"
	allocator "github.com/lrx0014/ScalableFlake/pkg/machine"
	"github.com/lrx0014/ScalableFlake/pkg/snowflake"
	log "github.com/sirupsen/logrus"
)

type GRPCServer struct {
	pb.UnimplementedUIDGeneratorServer
	allocator allocator.Allocator
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}

func (s *GRPCServer) AcquireUID(ctx context.Context, req *pb.AcquireUIDReq) (resp *pb.AcquireUIDResp, err error) {
	resp = &pb.AcquireUIDResp{}
	uid, err := snowflake.GenerateUID()
	if err != nil {
		log.Errorf("failed to generate uid: %v", err)
		return
	}

	resp.Uid = uid
	return
}
