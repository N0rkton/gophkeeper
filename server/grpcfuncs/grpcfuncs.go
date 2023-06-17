package grpcfuncs

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	pb "gophkeeper/proto"
	"sync"
)

type GophKeeperServer struct {
	pb.UnimplementedGophkeeperServer
	users sync.Map
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// допишите код
	// ...
	var token string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing token")
	}

	values := md.Get("token")
	if len(values) > 0 {
		// ключ содержит слайс строк, получаем первую строку
		token = values[0]
	}
	if token != "test" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return handler, nil
}
