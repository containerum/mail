package cherrygrpc

import (
	"context"
	"encoding/json"
	"net/http"

	"git.containerum.net/ch/kube-client/pkg/cherry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	JSONMarshal   = json.Marshal
	JSONUnmarshal = json.Unmarshal
)

var httpToGRPCCode = map[int]codes.Code{
	444:                                     codes.Canceled,
	http.StatusOK:                           codes.OK,
	http.StatusBadRequest:                   codes.InvalidArgument,
	http.StatusRequestTimeout:               codes.DeadlineExceeded,
	http.StatusNotFound:                     codes.NotFound,
	http.StatusConflict:                     codes.AlreadyExists,
	http.StatusForbidden:                    codes.PermissionDenied,
	http.StatusInsufficientStorage:          codes.ResourceExhausted,
	http.StatusPreconditionFailed:           codes.FailedPrecondition,
	http.StatusGatewayTimeout:               codes.Aborted,
	http.StatusRequestedRangeNotSatisfiable: codes.OutOfRange,
	http.StatusNotImplemented:               codes.Unimplemented,
	http.StatusInternalServerError:          codes.Internal,
	http.StatusServiceUnavailable:           codes.Unavailable,
	http.StatusUnauthorized:                 codes.Unauthenticated,
}

func toGRPC(errToPass *cherry.Err) error {
	data, err := JSONMarshal(errToPass)
	if err != nil {
		data = append(data, []byte("; with error "+err.Error())...)
	}
	code, mapped := httpToGRPCCode[errToPass.StatusHTTP]
	if !mapped {
		code = codes.Unknown
	}
	return status.Error(code, string(data))
}

// UnaryServerInterceptor -- middleware for grpc server to encode cherry error to grpc error.
// If used with "github.com/grpc-ecosystem/go-grpc-middleware".ChainUnaryServer it must be first in chain to work correctly.
func UnaryServerInterceptor(defaultErr func() *cherry.Err) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)

		if err == nil {
			return
		}

		switch err.(type) {
		case *cherry.Err:
			err = toGRPC(err.(*cherry.Err))
		default:
			err = toGRPC(defaultErr().AddDetailsErr(err))
		}

		return
	}
}

func fromGRPC(errToCheck error) (ret *cherry.Err, ok bool) {
	grpcErr, ok := status.FromError(errToCheck)
	if !ok {
		return
	}
	err := JSONUnmarshal([]byte(grpcErr.Message()), &ret)
	ok = err == nil
	return
}

// UnaryClientInterceptor -- grpc client middleware to decode cherry error from grpc error.
// If used with "github.com/grpc-ecosystem/go-grpc-middleware".ChainUnaryClient it must be first in chain to work correctly.
func UnaryClientInterceptor(defaultErr func() *cherry.Err) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)

		if err == nil {
			return nil
		}

		cherryErr, ok := fromGRPC(err)
		if !ok {
			return defaultErr().AddDetailsErr(err)
		}
		return cherryErr
	}
}
