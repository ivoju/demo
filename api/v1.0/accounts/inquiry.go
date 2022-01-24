package accounts

import (
	"context"

	pg "github.com/demo/pkg/v1.0/postgres"
	"github.com/demo/pkg/v1.0/utils/errors"
	actpb "github.com/demo/proto/v1.0/accounts"
	"google.golang.org/grpc/codes"
)

// ProcessController acts as the main entry point for this get service
func (s *Server) Inquiry(ctx context.Context, req *actpb.Request) (*actpb.Response, error) {
	err := isValidRequest(Inquiry, req)
	if err != nil {
		s.logger.Errorf("[ACCOUNT][INQUIRY] ERROR %v", err)
		return nil, err
	}

	rows, err := s.config.PgConn.CustomAccountsSelect(&pg.CustomAccounts{
		UserId: req.GetUserId(),
	})
	if err != nil {
		return nil, errors.FormatError(codes.InvalidArgument, "1004", err.Error())
	}

	var data []*actpb.Data

	for _, row := range rows {
		data = append(data, &actpb.Data{
			UserId:  row.UserId,
			DelFlag: row.DelFlag,
			Desc:    row.Desc.String,
			CreId:   row.CreId,
			CreTime: row.CreTime.String(),
			ModId:   row.ModId,
			ModTime: row.ModTime.String(),
		})
	}

	s.logger.Infof("[ACCOUNT][INQUIRY] SUCCESS")

	return &actpb.Response{
		Success:  true,
		RespCode: "0000",
		RespDesc: "Success",
		Data:     data,
	}, nil

}
