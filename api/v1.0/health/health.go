package health

import (
	"github.com/demo/config"
	"github.com/sirupsen/logrus"
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
