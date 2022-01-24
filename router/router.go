package router

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/demo/pkg/v1.0/utils/errors"
	"github.com/demo/pkg/v1.0/utils/jwt"
	"github.com/fatih/structs"
	"github.com/kenshaw/envcfg"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// create config and logger
	env, err := envcfg.New()
	if err != nil {
		return nil, err
	}

	// get jwt is enable config
	isJwtEnable, err := strconv.ParseBool(env.GetKey("jwt.isenable"))
	if err != nil {
		return nil, errors.FormatError(codes.InvalidArgument, "100", err.Error())
	}

	if isJwtEnable {
		// read header from incoming request
		headers, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.FormatError(codes.InvalidArgument, "101", "failed to read header from incoming request")
		}

		if len(headers.Get("Authorization")) < 1 || headers.Get("Authorization")[0] == "" {
			return nil, errors.FormatError(codes.InvalidArgument, "102", "invalid header: missing Authorization")
		}

		// claim jwt token
		jwtd := jwt.New(env)

		claims, err := jwtd.ClaimToken(headers.Get("Authorization")[0], info)
		if err != nil {
			return nil, errors.FormatError(codes.Unauthenticated, "401", err.Error())
		}

		r := structs.Map(req)

		if r["UserId"] != nil {
			if !(strings.Contains(info.FullMethod, "/Login") || strings.Contains(info.FullMethod, "/Register")) {
				if r["UserId"] != claims["userId"] {
					return nil, errors.FormatError(codes.Unauthenticated, "401", "The user is not authorized")
				}
			}
		}

		// ctx = context.WithValue(ctx, "userId", claims["userId"].(string))
	}

	return handler(ctx, req)
}

// ignoreErr returns true when err can be safely ignored.
func IgnoreErr(err error) bool {
	switch {
	case err == nil || err == http.ErrServerClosed || err == cmux.ErrListenerClosed:
		return true
	}
	if opErr, ok := err.(*net.OpError); ok {
		return opErr.Err.Error() == "use of closed network connection"
	}
	return false
}
