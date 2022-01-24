package accounts

import (
	"context"

	pg "github.com/demo/pkg/v1.0/postgres"
	"github.com/demo/pkg/v1.0/utils/errors"
	actpb "github.com/demo/proto/v1.0/accounts"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
)

// ProcessController acts as the main entry point for this get service
func (s *Server) Get(ctx context.Context, req *empty.Empty) (*actpb.Response, error) {
	rows, err := s.config.PgConn.CustomAccountsSelect(&pg.CustomAccounts{})
	if err != nil {
		return nil, errors.FormatError(codes.InvalidArgument, "103", err.Error())
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

	s.logger.Infof("[ACCOUNT][GET] SUCCESS")

	return &actpb.Response{
		Success:  true,
		RespCode: "000",
		RespDesc: "Success",
		Data:     data,
	}, nil

}
