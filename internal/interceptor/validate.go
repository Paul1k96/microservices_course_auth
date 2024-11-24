package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

type validator interface {
	Validate() error
}

// ValidateInterceptor validates request.
func ValidateInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handlerFunc grpc.UnaryHandler,
) (interface{}, error) {
	if v, ok := req.(validator); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	return handlerFunc(ctx, req)
}
