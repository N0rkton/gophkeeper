// Package grpcfuncs provides server-side gRPC functions for authentication, data management, and synchronization.
package grpcfuncs

import (
	"context"
	"gophkeeper/internal/datamodels"
	"gophkeeper/internal/sessionstorage"
	"gophkeeper/internal/storage"
	"gophkeeper/internal/utils"
	pb "gophkeeper/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetUserId - search UserID key in metadata
func GetUserId(ctx context.Context) string {
	var userId string
	var value []string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		value = md.Get("userid")
		if len(value) > 0 {
			userId = value[0]
			return userId
		}
	}
	return ""
}

// mapErr - maps err from storage to grpc error codes
func mapErr(err error) error {
	if err == storage.ErrDuplicate {
		return status.Errorf(codes.AlreadyExists, "login already exists")
	}
	if err == storage.ErrWrongPassword {
		return status.Errorf(codes.InvalidArgument, "wrong password")
	}
	if err == storage.ErrNotFound {
		return status.Errorf(codes.NotFound, "not found")
	}
	return status.Errorf(codes.Internal, "internal error")
}

// GophKeeperServer is the gRPC server implementation for GophKeeper.
type GophKeeperServer struct {
	pb.UnimplementedGophkeeperServer
	db    storage.Storage
	users sessionstorage.SessionStorage
}

// Init initializes the gRPC server.
func NewGophKeeperServer() GophKeeperServer {
	var err error
	var g GophKeeperServer
	g.db, err = storage.NewDBStorage("postgresql://localhost:5432/shvm")
	g.users = sessionstorage.NewAuthUsersStorage()
	if err != nil {
		log.Fatalf("err pinging db")
	}
	return g
}

// Auth handles the authentication request.
func (g *GophKeeperServer) Auth(ctx context.Context, in *pb.AuthLoginRequest) (*pb.AuthLoginResponse, error) {
	var resp pb.AuthLoginResponse
	passHash := utils.GetMD5Hash(in.Password)
	err := g.db.Auth(in.Login, passHash)
	if err != nil {
		return nil, mapErr(err)
	}
	id, err := g.db.Login(in.Login, passHash)
	if err != nil {
		return nil, mapErr(err)
	}
	token := utils.GenerateRandomString(5)
	if err = g.users.AddUser(token, id); err != nil {
		return nil, err
	}
	resp.Id = id
	md2 := metadata.New(map[string]string{"userid": token})
	outgoingCtx := metadata.NewOutgoingContext(ctx, md2)
	err = grpc.SetHeader(outgoingCtx, md2)
	if err != nil {
		return nil, status.Error(codes.Internal, "SetHeader err")
	}
	return &resp, nil
}

// Login handles the login request.
func (g *GophKeeperServer) Login(ctx context.Context, in *pb.AuthLoginRequest) (*pb.AuthLoginResponse, error) {
	var resp pb.AuthLoginResponse
	passHash := utils.GetMD5Hash(in.Password)
	id, err := g.db.Login(in.Login, passHash)
	if err != nil {
		return nil, mapErr(err)
	}
	token := utils.GenerateRandomString(5)
	err = g.users.AddUser(token, id)
	if err != nil {
		return nil, err
	}
	resp.Id = id
	header := metadata.Pairs("userid", token)
	grpc.SetHeader(ctx, header)

	return &resp, nil
}

// AddData handles the request to add data.
func (g *GophKeeperServer) AddData(ctx context.Context, in *pb.AddDataRequest) (*emptypb.Empty, error) {
	//TODO хранить зашифровано
	token := GetUserId(ctx)
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "token is empty")
	}
	id, err := g.users.GetUser(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user unauthenticated")
	}
	err = g.db.AddData(datamodels.Data{UserID: id, DataID: in.Data.DataId, Data: in.Data.Data, Metadata: in.Data.MetaInfo, ChangedAt: time.Now()})
	if err != nil {
		return nil, mapErr(err)
	}
	return new(emptypb.Empty), nil
}

// GetData handles the request to get data.
func (g *GophKeeperServer) GetData(ctx context.Context, in *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	var resp pb.GetDataResponse
	token := GetUserId(ctx)
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "token is empty")
	}
	id, err := g.users.GetUser(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user unauthenticated")
	}
	data, err := g.db.GetData(in.DataId, id)
	if err != nil {
		return nil, mapErr(err)
	}
	resp.Data = &pb.Data{DataId: in.DataId, Data: data.Data, MetaInfo: data.Metadata, ChangedAt: timestamppb.New(data.ChangedAt), Deleted: false}
	return &resp, nil
}

// DelData handles the request to delete data.
func (g *GophKeeperServer) DelData(ctx context.Context, in *pb.GetDataRequest) (*emptypb.Empty, error) {
	token := GetUserId(ctx)
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "token is empty")
	}
	id, err := g.users.GetUser(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user unauthenticated")
	}
	err = g.db.DelData(in.DataId, id)
	if err != nil {
		return nil, mapErr(err)
	}
	return new(emptypb.Empty), nil
}

// Sync handles the synchronization request.
func (g *GophKeeperServer) Sync(ctx context.Context, in *emptypb.Empty) (*pb.SynchronizationResponse, error) {
	var resp pb.SynchronizationResponse
	token := GetUserId(ctx)
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "token is empty")
	}
	id, err := g.users.GetUser(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user unauthenticated")
	}
	data, err := g.db.Sync(id)
	if err != nil {
		return nil, mapErr(err)
	}
	if data != nil {
		for _, v := range data {
			resp.Data = append(resp.Data, &pb.Data{DataId: v.DataID, Data: v.Data, MetaInfo: v.Metadata, Deleted: v.Deleted, ChangedAt: timestamppb.New(v.ChangedAt)})
		}
	}
	return &resp, nil

}

// ClientSync handles the client synchronization request.
func (g *GophKeeperServer) ClientSync(ctx context.Context, in *pb.ClientSyncRequest) (*emptypb.Empty, error) {
	token := GetUserId(ctx)
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "token is empty")
	}
	id, err := g.users.GetUser(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user unauthenticated")
	}
	err = g.db.ClientSync(id, in.Data)
	if err != nil {
		return nil, mapErr(err)
	}
	return new(emptypb.Empty), nil
}
