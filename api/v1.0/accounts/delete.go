package accounts

import (
	"context"

	pg "github.com/demo/pkg/v1.0/postgres"
	"github.com/demo/pkg/v1.0/utils/errors"
	actpb "github.com/demo/proto/v1.0/accounts"
	"google.golang.org/grpc/codes"
)

// ProcessController acts as the main entry point for this get service
func (s *Server) Delete(ctx context.Context, req *actpb.Request) (*actpb.Response, error) {
	err := isValidRequest(Delete, req)
	if err != nil {
		s.logger.Errorf("[ACCOUNT][DELETE] ERROR %v", err)
		return nil, err
	}

	_, err = s.config.PgConn.CustomAccountsUpdate(&pg.CustomAccounts{
		DelFlag: true,
	}, &pg.CustomAccounts{
		UserId: req.GetUserId(),
	})
	if err != nil {
		return nil, errors.FormatError(codes.InvalidArgument, "1004", err.Error())
	}

	s.logger.Infof("[ACCOUNT][DELETE] SUCCESS")

	return &actpb.Response{
		Success:  true,
		RespCode: "0000",
		RespDesc: "Success",
	}, nil

}
