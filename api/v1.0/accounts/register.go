package accounts

import (
	"context"

	pg "github.com/demo/pkg/v1.0/postgres"
	"github.com/demo/pkg/v1.0/utils/errors"
	actpb "github.com/demo/proto/v1.0/accounts"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
)

// ProcessController acts as the main entry point for this get service
func (s *Server) Register(ctx context.Context, req *actpb.Request) (*actpb.Response, error) {
	err := isValidRequest(Register, req)
	if err != nil {
		s.logger.Errorf("[ACCOUNT][REGISTER] ERROR %v", err)
		return nil, err
	}

	passHashed, err := bcrypt.GenerateFromPassword([]byte(req.GetPass()), bcrypt.MinCost)
	if err != nil {
		return nil, errors.FormatError(codes.InvalidArgument, "103", err.Error())
	}

	_, err = s.config.PgConn.CustomAccountsInsert(&pg.CustomAccounts{
		UserId: req.GetUserId(),
		Pass:   string(passHashed),
		CreId:  req.GetUserId(),
		ModId:  req.GetUserId(),
	})
	if err != nil {
		return nil, errors.FormatError(codes.InvalidArgument, "104", err.Error())
	}

	s.logger.Infof("[ACCOUNT][REGISTER] SUCCESS")

	return &actpb.Response{
		Success:  true,
		RespCode: "000",
		RespDesc: "Success",
	}, nil

}
