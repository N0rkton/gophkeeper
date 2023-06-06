// Package storage provides implementations for data storage functions.
package storage

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gophkeeper/internal/datamodels"
	"gophkeeper/internal/sessionstorage"
	"gophkeeper/internal/utils"
	pb "gophkeeper/proto"
	"log"
	"time"
)

//todo obratnaya sync

// Client - grpc default client
var Client pb.GophkeeperClient

// Module errors
var (
	ErrNotFound      = errors.New("not found")
	ErrWrongPassword = errors.New("invalid password")
	ErrInternal      = errors.New("server error")
	ErrDuplicate     = errors.New("login already exists")
)

// Storage an interface that defines the following methods:
type Storage interface {
	//Auth - adds new user
	Auth(login string, password string) error
	// Login verifies the login credentials.
	Login(login string, password string) (uint32, error)
	// AddData adds data to the storage.
	AddData(data datamodels.Data) error
	// GetData retrieves data from the storage.
	GetData(dataID string, userID uint32) (datamodels.Data, error)
	// DelData deletes data from the storage.
	DelData(dataID string, userID uint32) error
	// Sync synchronizes data from server for a specific user.
	Sync(userId uint32) ([]datamodels.Data, error)
	//ClientSync - synchronize client data with server
	ClientSync(userID uint32, data []*pb.Data) error
}

// Users represents user sessions.
var Users sessionstorage.UserSession
var md metadata.MD

// Init initializes the storage package by establishing a gRPC connection.
func Init() {
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	Client = pb.NewGophkeeperClient(conn)
}

type uniqueData struct {
	dataID string
	userId uint32
}

// MemoryStorage a struct that implements the Storage interface and stores data in the computer's memory.
type MemoryStorage struct {
	localMem map[uniqueData]datamodels.Data
}

// NewMemoryStorage creates a new MemoryStorage instance.
func NewMemoryStorage() Storage {
	Users = sessionstorage.Init()
	return &MemoryStorage{localMem: make(map[uniqueData]datamodels.Data)}
}

// Auth adds a new user.
// If the user already exists, it returns an error.
func (ms *MemoryStorage) Auth(login string, password string) error {
	var header metadata.MD
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	_, err := Client.Auth(context.Background(), &pb.AuthLoginRequest{Login: login, Password: password}, grpc.Header(&header))
	md = header
	st := status.Convert(err)
	if st.Err() == nil {

		ctx = metadata.NewOutgoingContext(context.Background(), md)
		id, err := Client.Login(ctx, &pb.AuthLoginRequest{Login: login, Password: password}, grpc.Header(&header))
		md = header

		st = status.Convert(err)
		if st.Err() != nil {
			return st.Err()
		}
		passHash := utils.GetMD5Hash(password)
		err = Users.AddUser(login, passHash, id.Id)
		if err != nil {
			return errors.New("user already exists")
		}
		return nil
	}
	return st.Err()
}

// Login verifies the login credentials.
func (ms *MemoryStorage) Login(login string, password string) (uint32, error) {
	var header metadata.MD
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	id, err := Client.Login(ctx, &pb.AuthLoginRequest{Login: login, Password: password}, grpc.Header(&header))
	md = header
	if err == nil {
		return id.Id, nil
	}
	user, ok := Users.GetUser(login)
	if !ok {
		return 0, errors.New("user not found")
	}
	passHash := utils.GetMD5Hash(password)
	if user.Password != passHash {
		return 0, errors.New("wrong password")
	}
	return user.Id, nil
}

// AddData adds data to the storage.
func (ms *MemoryStorage) AddData(data datamodels.Data) error {
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	_, err := Client.AddData(ctx, &pb.AddDataRequest{Data: &pb.Data{DataId: data.DataID, Data: data.Data, MetaInfo: data.Metadata}})
	if err == nil {
		return nil
	}

	ms.localMem[uniqueData{dataID: data.DataID, userId: data.UserID}] = datamodels.Data{UserID: data.UserID, Data: data.Data, Metadata: data.Metadata, Deleted: false, ChangedAt: time.Now()}
	//TODO записывать в файл
	return nil
}

// DelData deletes data from the storage.
func (ms *MemoryStorage) DelData(dataID string, userID uint32) error {
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	Client.DelData(ctx, &pb.GetDataRequest{DataId: dataID})
	user, _ := ms.localMem[uniqueData{dataID: dataID, userId: userID}]
	if user.UserID == userID {
		user.Deleted = true
		ms.localMem[uniqueData{dataID: dataID, userId: userID}] = user
	}
	//TODO записывать в файл
	return nil
}

// GetData retrieves data from the storage.
func (ms *MemoryStorage) GetData(dataID string, userID uint32) (datamodels.Data, error) {
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resp, err := Client.GetData(ctx, &pb.GetDataRequest{DataId: dataID})
	if err == nil {
		response := datamodels.Data{DataID: resp.Data.DataId, Data: resp.Data.Data, UserID: userID, Metadata: resp.Data.MetaInfo}
		return response, nil
	}

	data, ok := ms.localMem[uniqueData{dataID: dataID, userId: userID}]
	if !ok {
		return datamodels.Data{}, errors.New("no data find")
	}
	if data.UserID == userID {
		return data, nil
	}
	return datamodels.Data{}, errors.New("no data find")
}

// Sync synchronizes data from server for a specific user.
func (ms *MemoryStorage) Sync(userId uint32) ([]datamodels.Data, error) {
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resp, err := Client.Sync(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	var response []datamodels.Data
	for _, v := range resp.Data {
		data, _ := ms.localMem[uniqueData{dataID: v.DataId, userId: userId}]
		if data.ChangedAt.Before(data.ChangedAt) {
			ms.localMem[uniqueData{dataID: v.DataId, userId: userId}] = datamodels.Data{DataID: v.DataId, Data: v.Data, UserID: userId, Metadata: v.MetaInfo, Deleted: v.Deleted, ChangedAt: v.ChangedAt.AsTime()}
			response = append(response, datamodels.Data{DataID: v.DataId, Data: v.Data, UserID: userId, Metadata: v.MetaInfo, Deleted: v.Deleted, ChangedAt: v.ChangedAt.AsTime()})
		}
	}
	return response, nil
}

// ClientSync - synchronize client data with server
func (ms *MemoryStorage) ClientSync(userID uint32, data []*pb.Data) error {
	var req []*pb.Data
	for k, v := range ms.localMem {
		if k.userId == userID {
			req = append(req, &pb.Data{Data: v.Data, DataId: v.DataID, MetaInfo: v.Metadata, Deleted: v.Deleted, ChangedAt: timestamppb.New(v.ChangedAt)})
		}
	}
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	Client.ClientSync(ctx, &pb.ClientSyncRequest{Data: req})
	return nil
}
