package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type metaData struct {
	userAgent string
	clientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *metaData {
	meta := &metaData{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md[grpcGatewayUserAgentHeader]; ok {
			meta.userAgent = val[0]
		}
		if val, ok := md[userAgentHeader]; ok {
			meta.userAgent = val[0]
		}
		if val, ok := md[xForwardedForHeader]; ok {
			meta.clientIP = val[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		meta.clientIP = p.Addr.String()
	}
	return meta
}
