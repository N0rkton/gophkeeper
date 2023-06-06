package grpcfuncs

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gophkeeper/internal/datamodels"
	"gophkeeper/internal/sessionstorage"
	"gophkeeper/internal/storage"
	"gophkeeper/internal/utils"
	pb "gophkeeper/proto"
	"log"
	"time"
)

// GetUserId - search UserId key in metadata
func GetUserId(ctx context.Context) string {
	var userId string
	var value []string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		value = md.Get("UserId")
		if len(value) > 0 {
			// ключ содержит слайс строк, получаем первую строку
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

type GophKeeperServer struct {
	pb.UnimplementedGophkeeperServer
}

var db storage.Storage
var users sessionstorage.SessionStorage

func Init() {
	var err error
	db, err = storage.NewDBStorage("postgresql://localhost:5432/shvm")
	if err != nil {
		log.Fatalf("err pinging db")
	}
}
func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing token")
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	return handler(newCtx, req)

}
func (g *GophKeeperServer) Auth(ctx context.Context, in *pb.AuthLoginRequest) (*pb.AuthLoginResponse, error) {
	var resp pb.AuthLoginResponse
	passHash := utils.GetMD5Hash(in.Password)
	err := db.Auth(in.Login, passHash)
	if err != nil {
		return nil, mapErr(err)
	}
	id, err := db.Login(in.Login, in.Password)
	if err != nil {
		return nil, mapErr(err)
	}
	token := utils.GenerateRandomString(5)
	if err = users.AddUser(token, id); err != nil {
		return nil, err
	}
	resp.Id = id
	md2 := metadata.New(map[string]string{"UserId": token})
	metadata.NewIncomingContext(context.Background(), md2)
	err = grpc.SetHeader(ctx, md2)
	if err != nil {
		return nil, status.Error(codes.Internal, "SetHeader err")
	}
	return &resp, nil
}
func (g *GophKeeperServer) Login(ctx context.Context, in *pb.AuthLoginRequest) (*pb.AuthLoginResponse, error) {
	var resp pb.AuthLoginResponse
	id, err := db.Login(in.Login, in.Password)
	if err != nil {
		return nil, mapErr(err)
	}
	token := utils.GenerateRandomString(5)
	if err = users.AddUser(token, id); err != nil {
		return nil, err
	}
	resp.Id = id
	md2 := metadata.New(map[string]string{"UserId": token})
	metadata.NewIncomingContext(context.Background(), md2)
	err = grpc.SetHeader(ctx, md2)
	if err != nil {
		return nil, status.Error(codes.Internal, "SetHeader err")
	}
	return &resp, nil
}
func (g *GophKeeperServer) AddData(ctx context.Context, in *pb.AddDataRequest) (*pb.AddDelDataResponse, error) {
	//TODO хранить зашифровано
	//TODO vozvrashat' id zapisi?
	token := GetUserId(ctx)
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "token is empty")
	}
	id, err := users.GetUser(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user unauthenticated")
	}
	err = db.AddData(datamodels.Data{UserID: id, DataID: in.Data.DataId, Data: in.Data.Data, Metadata: in.Data.MetaInfo, ChangedAt: time.Now()})
	if err != nil {
		return nil, mapErr(err)
	}
	return nil, nil
}
func (g *GophKeeperServer) GetData(ctx context.Context, in *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	var resp pb.GetDataResponse
	token := GetUserId(ctx)
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "token is empty")
	}
	id, err := users.GetUser(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user unauthenticated")
	}
	data, err := db.GetData(in.DataId, id)
	if err != nil {
		return nil, mapErr(err)
	}
	resp.Data.DataId = in.DataId
	resp.Data.Data = data.Data
	resp.Data.MetaInfo = data.Metadata
	return &resp, nil
}
func (g *GophKeeperServer) DelData(ctx context.Context, in *pb.GetDataRequest) (*pb.AddDelDataResponse, error) {
	token := GetUserId(ctx)
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "token is empty")
	}
	id, err := users.GetUser(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user unauthenticated")
	}
	err = db.DelData(in.DataId, id)
	if err != nil {
		return nil, mapErr(err)
	}
	return nil, nil
}
func (g *GophKeeperServer) Sync(ctx context.Context, in *emptypb.Empty) (*pb.SynchronizationResponse, error) {
	var resp pb.SynchronizationResponse
	token := GetUserId(ctx)
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "token is empty")
	}
	id, err := users.GetUser(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user unauthenticated")
	}
	data, err := db.Sync(id)
	if err != nil {
		return nil, mapErr(err)
	}

	for _, v := range data {

		resp.Data = append(resp.Data, &pb.Data{DataId: v.DataID, Data: v.Data, MetaInfo: v.Metadata})
	}
	return &resp, nil
}
