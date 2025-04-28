package server

import (
	"context"
	pb "github.com/lrx0014/ScalableFlake/api/v1"
	allocator "github.com/lrx0014/ScalableFlake/pkg/machine"
)

type GRPCServer struct {
	pb.UnimplementedUIDGeneratorServer
	allocator allocator.Allocator
}

func NewGRPCServer(a allocator.Allocator) *GRPCServer {
	return &GRPCServer{allocator: a}
}

func (s *GRPCServer) AcquireUID(ctx context.Context, req *pb.AcquireUIDReq) (resp *pb.AcquireUIDResp, err error) {
	resp = &pb.AcquireUIDResp{}
	uid, err := s.allocator.Acquire(ctx, req.TenantId)
	if err != nil {
		return
	}

	resp.Uid = uid
	return
}
