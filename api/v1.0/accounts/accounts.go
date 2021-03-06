package accounts

import (
	"github.com/demo/config"
	"github.com/demo/pkg/v1.0/utils/errors"
	actpb "github.com/demo/proto/v1.0/accounts"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
)

// Method is the method type.
type Method int

const (
	// List of different Methods
	Inquiry Method = iota
	Register
	Login
	Delete
	GetList
)

// Server is the server object for this api service.
type Server struct {
	config *config.Config
	logger *logrus.Logger
}

// New creates a new server.
func New(config *config.Config, logger *logrus.Logger) *Server {
	return &Server{
		config: config,
		logger: logger,
	}
}

// isValidRequest validates the status request
func isValidRequest(m Method, req *actpb.Request) error {
	switch m {
	case Register:
		if req.GetPass() == "" {
			return errors.FormatError(codes.InvalidArgument, "1003", "invalid request: missing pass")
		}
	case Login:
		if req.GetPass() == "" {
			return errors.FormatError(codes.InvalidArgument, "1003", "invalid request: missing pass")
		}
	}

	return nil
}
