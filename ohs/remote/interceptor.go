package remote

import (
	"context"
	"log"
	"order-service/ohs/local/pl/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryInterceptor unary interceptor
func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	var siteCode string
	if !ok {
		return nil, errors.MissingMetadataError("missing metadata")
	}

	if val, ok := md["site-code"]; ok {
		siteCode = val[0]
		// md.Append("site-code", siteCode)
	}

	if !validSiteCode(siteCode) {
		return nil, errors.PreconditionFailed("precondition failed")
	}
	// newCtx := metadata.NewIncomingContext(ctx, md)
	m, err := handler(ctx, req)
	if err != nil {
		log.Fatalf("RPC failed with error: %v", err)
	}

	return m, err
}

func validSiteCode(siteCode string) bool {
	if len(siteCode) == 0 {
		return false
	}
	return true
}
