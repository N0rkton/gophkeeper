// Package storage provides implementations for data storage functions.
package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gophkeeper/internal/datamodels"
	"gophkeeper/internal/utils"
	pb "gophkeeper/proto"
	"time"
)

// dbSecret - secret key for cipher
var dbSecret = []byte("alskdjfhgnbvcmrt")

// DBStorage is a struct that represents a storage implementation using a PostgreSQL database.
type DBStorage struct {
	db *sql.DB
}

// NewDBStorage creates a new DBStorage instance with the provided database path.
func NewDBStorage(path string) (Storage, error) {
	if path == "" {
		return nil, errors.New("invalid db address")
	}
	db, err := sql.Open("pgx", path)
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./database/migration",
		"postgres", driver)
	if err != nil {
		return nil, err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}
	return &DBStorage{db: db}, nil
}

// Auth adds a new user with the provided login and password to the storage.
func (dbs *DBStorage) Auth(login string, password string) error {
	_, err := dbs.db.Exec("insert into users (login, password) values ($1, $2);", login, password)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return ErrDuplicate
	}
	return err
}

// Login verifies the login credentials of a user and returns the user ID if successful.
func (dbs *DBStorage) Login(login string, password string) (uint32, error) {
	rows := dbs.db.QueryRow("select id,password from users where login=$1 limit 1;", login)
	var v datamodels.Login
	err := rows.Scan(&v.ID, &v.Password)
	if err != nil {
		return 0, ErrNotFound
	}
	if v.Password != password {
		return 0, ErrWrongPassword
	}
	return v.ID, nil
}

// AddData adds new data to the storage.
func (dbs *DBStorage) AddData(data datamodels.Data) error {
	query := `insert into keeper (data_id,user_id, data_info,meta_info, changed_at) values ($1, $2,$3,$4,$5) ON CONFLICT (user_id, data_id) DO UPDATE SET data_info=EXCLUDED.data_info, meta_info=EXCLUDED.meta_info, changed_at=EXCLUDED.changed_at where keeper.changed_at < $5;`
	data.Data = utils.Encrypt(data.Data, dbSecret)
	data.Metadata = utils.Encrypt(data.Metadata, dbSecret)
	_, err := dbs.db.Exec(query, data.DataID, data.UserID, data.Data, data.Metadata, data.ChangedAt.Format(time.RFC3339))
	if err != nil {
		return ErrInternal
	}
	return nil
}

// GetData retrieves data from the storage based on the data ID and user ID.
func (dbs *DBStorage) GetData(dataID string, userID uint32) (datamodels.Data, error) {
	rows := dbs.db.QueryRow("select data_info,meta_info, changed_at from keeper where data_id=$1 and user_id=$2 and deleted=false limit 1;", dataID, userID)
	var v datamodels.Data
	err := rows.Scan(&v.Data, &v.Metadata, &v.ChangedAt)
	v.Data = utils.Decrypt(v.Data, dbSecret)
	v.Metadata = utils.Decrypt(v.Metadata, dbSecret)
	if err != nil {
		return datamodels.Data{}, ErrNotFound
	}
	return v, nil
}

// DelData marks data as deleted in the storage based on the data ID and user ID.
func (dbs *DBStorage) DelData(dataID string, userID uint32) error {
	_, err := dbs.db.Exec("UPDATE  keeper set deleted=true where data_id=$1 and user_id=$2;", dataID, userID)
	if err != nil {
		return ErrInternal
	}
	return nil
}

// Sync retrieves all data associated with a user from the storage.
func (dbs *DBStorage) Sync(userID uint32) ([]datamodels.Data, error) {
	rows, err := dbs.db.Query("SELECT(data_id,data_info,meta_info,deleted,changed_at) from keeper where  user_id=$1;", userID)
	if err != nil {
		return nil, ErrInternal
	}
	var resp []datamodels.Data
	var tmp datamodels.Data

	for rows.Next() {
		err = rows.Scan(&tmp.DataID, &tmp.Data, &tmp.Metadata, &tmp.Deleted, &tmp.ChangedAt)
		tmp.Data = utils.Decrypt(tmp.Data, dbSecret)
		tmp.Metadata = utils.Decrypt(tmp.Metadata, dbSecret)
		if err == nil {
			resp = append(resp, tmp)
		}
	}
	if resp != nil {
		return resp, nil
	}
	return nil, nil
}

// ClientSync synchronizes client data with the server in the storage.
func (dbs *DBStorage) ClientSync(userID uint32, data []*pb.Data) error {
	query := `insert into keeper (data_id,user_id, data_info,meta_info, changed_at,deleted) values ($1, $2,$3,$4,$5,$6) ON CONFLICT (user_id, data_id) DO UPDATE SET data_info=EXCLUDED.data_info, meta_info=EXCLUDED.meta_info, changed_at=EXCLUDED.changed_at where keeper.changed_at < $5;`
	for i := range data {

		fmt.Println(i, data[i].Data)
		data[i].Data = utils.Encrypt(data[i].Data, dbSecret)
		fmt.Println(i, data[i].Data)
		fmt.Println(i, data[i].MetaInfo)
		data[i].MetaInfo = utils.Encrypt(data[i].MetaInfo, dbSecret)
		fmt.Println(i, data[i].MetaInfo)
		_, err := dbs.db.Exec(query, data[i].DataId, userID, data[i].Data, data[i].MetaInfo, data[i].ChangedAt.AsTime().Format(time.RFC3339), data[i].Deleted)
		if err != nil {
			return ErrInternal
		}
	}
	return nil
}
