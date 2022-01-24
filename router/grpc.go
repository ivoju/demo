package router

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/demo/api/v1.0/accounts"
	"github.com/demo/config"
	actpb "github.com/demo/proto/v1.0/accounts"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

// newGRPCServer creates the grpc server serve mux.
func NewGRPCServer(config *config.Config, logger *logrus.Logger) error {
	lis, err := net.Listen("tcp", ":"+string(config.Env.GetKey("grpc.port")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	// Logrus entry is used, allowing pre-definition of certain fields by the user.
	logrusEntry := logrus.NewEntry(logger)

	opts := []grpc_logrus.Option{grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
		return "grpc.time_ns", duration.Nanoseconds()
	})}

	alwaysLoggingDeciderServer := func(ctx context.Context, fullMethodName string, servingObject interface{}) bool { return true }

	// register grpc service server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			AuthInterceptor,
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
			grpc_logrus.PayloadUnaryServerInterceptor(logrusEntry, alwaysLoggingDeciderServer),
		)))

	actpb.RegisterAccountsServiceServer(grpcServer, accounts.New(config, logger))

	// add reflection service
	reflection.Register(grpcServer)

	// running gRPC server
	log.Println("[SERVER] GRPC server is ready")
	grpcServer.Serve(lis)

	return nil
}
