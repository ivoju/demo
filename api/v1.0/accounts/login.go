package accounts

import (
	"context"

	pg "github.com/demo/pkg/v1.0/postgres"
	"github.com/demo/pkg/v1.0/utils/errors"
	"github.com/demo/pkg/v1.0/utils/jwt"
	actpb "github.com/demo/proto/v1.0/accounts"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
)

// ProcessController acts as the main entry point for this get service
func (s *Server) Login(ctx context.Context, req *actpb.Request) (*actpb.Response, error) {
	err := isValidRequest(Login, req)
	if err != nil {
		s.logger.Errorf("[ACCOUNT][LOGIN] ERROR %v", err)
		return nil, err
	}

	rows, err := s.config.PgConn.CustomAccountsSelect(&pg.CustomAccounts{
		UserId: req.GetUserId(),
	})
	if err != nil {
		return nil, errors.FormatError(codes.InvalidArgument, "1004", err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(rows[0].Pass), []byte(req.GetPass()))
	if err != nil {
		return nil, errors.FormatError(codes.InvalidArgument, "1005", err.Error())
	}

	jwtd := jwt.New(s.config.Env)

	token, err := jwtd.GenerateToken(req.GetUserId())
	if err != nil {
		return nil, errors.FormatError(codes.InvalidArgument, "1006", err.Error())
	}

	s.logger.Infof("[ACCOUNT][LOGIN] SUCCESS")

	return &actpb.Response{
		Success:  true,
		RespCode: "0000",
		RespDesc: "Success",
		Token:    token,
	}, nil

}
