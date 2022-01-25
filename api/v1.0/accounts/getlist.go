package accounts

import (
	"context"
	"time"

	pg "github.com/demo/pkg/v1.0/postgres"
	"github.com/demo/pkg/v1.0/utils/errors"
	"github.com/demo/pkg/v1.0/utils/requestapi"
	actpb "github.com/demo/proto/v1.0/accounts"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
)

// ProcessController acts as the main entry point for this get service
func (s *Server) GetList(ctx context.Context, req *empty.Empty) (*actpb.Response, error) {
	rows, err := s.config.PgConn.CustomAccountsSelect(&pg.CustomAccounts{})
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

	header := make(map[string]interface{})
	header["ContentType"] = "Application/Json"

	_, _ = requestapi.POST(requestapi.ReqInfo{
		URL:        "http://google.com/login",
		HeaderInfo: header,
		Body:       []byte("Halo"),
	}, 30*time.Second)

	s.logger.Infof("[ACCOUNT][GETLIST] SUCCESS")

	return &actpb.Response{
		Success:  true,
		RespCode: "0000",
		RespDesc: "Success",
		Data:     data,
	}, nil

}
