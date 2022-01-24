package router

import (
	"context"
	"log"
	"net/http"

	"github.com/demo/config"
	"github.com/demo/pkg/v1.0/utils/errors"
	actpb "github.com/demo/proto/v1.0/accounts"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// newHTTPServer creates the http server serve mux.
func NewHTTPServer(config *config.Config, logger *logrus.Logger) error {
	runtime.HTTPError = errors.CustomHTTPError

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Connect to the GRPC server
	addr := "0.0.0.0:" + config.Env.GetKey("grpc.port")
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		return err
	}
	defer conn.Close()

	// Create new grpc-gateway
	rmux := runtime.NewServeMux()

	// register gateway endpoints
	for _, f := range []func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error{
		// register grpc service handler
		actpb.RegisterAccountsServiceHandler,
	} {
		if err = f(ctx, rmux, conn); err != nil {
			log.Fatal(err)
			return err
		}
	}

	// create http server mux
	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	// run swagger server
	if config.Env.GetKey("runtime.environment") == "development" {
		CreateSwagger(mux)
	}

	// running rest http server
	log.Println("[SERVER] REST HTTP server is ready")
	err = http.ListenAndServe("0.0.0.0:"+config.Env.PortString(), mux)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// newHTTPEncodedAndSwaggerServer creates the swagger server serve mux.
func CreateSwagger(gwmux *http.ServeMux) {
	// register swagger service server
	gwmux.HandleFunc("/api/demo/docs.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger/docs.json")
	})

	// load swagger-ui file
	fs := http.FileServer(http.Dir("swagger/swagger-ui"))
	gwmux.Handle("/api/demo/docs/", http.StripPrefix("/api/demo/docs", fs))
}
