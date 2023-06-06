package storage

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"gophkeeper/internal/datamodels"
	"gophkeeper/internal/sessionstorage"
	pb "gophkeeper/proto"
	"log"
)

// TODO время изменения тоже передавать в воркер
var Client pb.GophkeeperClient
var (
	ErrNotFound      = errors.New("not found")
	ErrWrongPassword = errors.New("invalid password")
	ErrInternal      = errors.New("server error")
	ErrDuplicate     = errors.New("login already exists")
)

type Storage interface {
	Auth(login string, password string) error
	Login(login string, password string) (uint32, error)
	AddData(data datamodels.Data) error
	GetData(dataID string, userID uint32) (datamodels.Data, error)
	DelData(dataID string, userID uint32) error
	Sync(userId uint32) ([]datamodels.Data, error)
}
type storeInfo struct {
	userID   uint32
	dataID   string
	data     string
	metaInfo string
}

var Users sessionstorage.UserSession

func Init() {
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	Client = pb.NewGophkeeperClient(conn)
}

// MemoryStorage a struct that implements the Storage interface and stores data in the computer's memory.
type MemoryStorage struct {
	localMem map[string]storeInfo
}

// NewMemoryStorage creates a new MemoryStorage instance.
func NewMemoryStorage() Storage {
	//todo init
	//if err != nil {
	//	log.Fatalf("err reading secretfile")
	//}
	Users = sessionstorage.Init()
	return &MemoryStorage{localMem: make(map[string]storeInfo)}
}

func (ms *MemoryStorage) Auth(login string, password string) error {
	//TODO асинхронный воркер в который добавляются все запросы
	//все что снизу по должно уехать в воркер
	_, err := Client.Auth(context.Background(), &pb.AuthLoginRequest{Login: login, Password: password})
	st := status.Convert(err)
	if st.Err() == nil {
		return nil
	}
	if st.Code() == codes.Unauthenticated {
		//TODO узнать код если сервер не отвечает
	}
	ok := Users.AddUser(login, password)
	if ok != nil {
		return errors.New("user already exists")
	}
	return nil
}
func (ms *MemoryStorage) Login(login string, password string) (uint32, error) {
	realPass, ok := Users.GetUser(login)
	if !ok {
		return 0, errors.New("user not found")
	}
	if realPass != password {
		return 0, errors.New("wrong password")
	}
	return 0, nil
}

func (ms *MemoryStorage) AddData(data datamodels.Data) error {
	ms.localMem[data.DataID] = storeInfo{userID: data.UserID, data: data.Data, metaInfo: data.Metadata}
	//TODO записывать в файл
	return nil
}

func (ms *MemoryStorage) DelData(dataID string, userID uint32) error {
	user, ok := ms.localMem[dataID]
	if !ok {
		return errors.New("no data to delete")
	}
	if user.userID == userID {
		delete(ms.localMem, dataID)
	}
	//TODO записывать в файл
	return nil
}
func (ms *MemoryStorage) GetData(dataID string, userID uint32) (datamodels.Data, error) {
	data, ok := ms.localMem[dataID]
	if !ok {
		return datamodels.Data{}, errors.New("no data find")
	}
	if data.userID == userID {
		return datamodels.Data{DataID: dataID, Data: data.data, Metadata: data.metaInfo}, nil
	}
	return datamodels.Data{}, errors.New("no data find")
}
func (ms *MemoryStorage) Sync(userId uint32) ([]datamodels.Data, error) {
	return nil, nil
}
