package service

import (
	"context"
	pb "im-server/api/v1"

	"im-server/internal/biz"
	"im-server/internal/biz/do"
)

// GreeterService is a greeter service.
type GreeterService struct {
	pb.UnimplementedGreeterSvcServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

func (g *GreeterService) ListGreeter(ctx context.Context, req *pb.ListGreeterReq) (*pb.ListGreeterRsp, error) {
	greeters, err := g.uc.ListAllGreeter(ctx)
	if err != nil {
		return nil, err
	}
	rsp := &pb.ListGreeterRsp{
		Body: &pb.ListGreeterRsp_Body{
			Greeters: convertGreeterDoToPbs(greeters),
		},
	}
	return rsp, nil
}

func convertGreeterDoToPbs(gs []*do.Greeter) []*pb.Greeter {
	ret := make([]*pb.Greeter, 0)
	for _, g := range gs {
		ret = append(ret, convertGreeterDoToPb(g))
	}
	return ret
}

func convertGreeterDoToPb(g *do.Greeter) *pb.Greeter {
	return &pb.Greeter{
		Name: g.Name,
		Age:  g.Age,
	}
}
