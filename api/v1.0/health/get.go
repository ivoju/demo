package health

import (
	"context"

	hlpb "github.com/demo/proto/v1.0/health"
	"github.com/golang/protobuf/ptypes/empty"
)

// ProcessController acts as the main entry point for this get service
func (s *Server) Get(ctx context.Context, req *empty.Empty) (*hlpb.Response, error) {
	s.logger.Infof("[HEALTH][GET] SUCCESS")

	return &hlpb.Response{
		Success:  true,
		RespCode: "0000",
		RespDesc: "Success",
	}, nil

}
