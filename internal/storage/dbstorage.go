package storage

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gophkeeper/internal/datamodels"
	pb "gophkeeper/proto"
	"time"
)

type DBStorage struct {
	db *sql.DB
}

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
func (dbs *DBStorage) Auth(login string, password string) error {
	_, err := dbs.db.Exec("insert into users (login, password) values ($1, $2);", login, password)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return ErrDuplicate
	}
	return err
}
func (dbs *DBStorage) Login(login string, password string) (uint32, error) {
	rows := dbs.db.QueryRow("select id,password from users where login=$1 limit 1;", login)
	var v datamodels.Auth
	err := rows.Scan(&v.ID, &v.Password)
	if err != nil {
		return 0, ErrNotFound
	}
	if v.Password != password {
		return 0, ErrWrongPassword
	}
	return v.ID, nil
}
func (dbs *DBStorage) AddData(data datamodels.Data) error {
	query := `insert into keeper (data_id,user_id, data_info,meta_info, changed_at) values ($1, $2,$3,$4,$5) ON CONFLICT (user_id, data_id) DO UPDATE SET data_info=EXCLUDED.data_info, meta_info=EXCLUDED.meta_info, changed_at=EXCLUDED.changed_at where keeper.changed_at < $5;`
	_, err := dbs.db.Exec(query, data.DataID, data.UserID, data.Data, data.Metadata, data.ChangedAt.Format(time.RFC3339))
	if err != nil {
		return ErrInternal
	}
	return nil
}
func (dbs *DBStorage) GetData(dataID string, userID uint32) (datamodels.Data, error) {
	rows := dbs.db.QueryRow("select data_info,meta_info from keeper where data_id=$1 and user_id=$2 limit 1;", dataID, userID)
	var v datamodels.Data
	err := rows.Scan(&v.Data, &v.Metadata)
	if err != nil {
		return datamodels.Data{}, ErrNotFound
	}
	return v, nil
}
func (dbs *DBStorage) DelData(dataID string, userID uint32) error {
	_, err := dbs.db.Exec("UPDATE  keeper set deleted=true where data_id=$1 and user_id=$2;", dataID, userID)
	if err != nil {
		return ErrInternal
	}
	return nil
}
func (dbs *DBStorage) Sync(userID uint32) ([]datamodels.Data, error) {
	rows, err := dbs.db.Query("SELECT(data_id,data_info,meta_info,deleted,changed_at) from keeper where  user_id=$1;", userID)
	if err != nil {
		return nil, ErrInternal
	}
	var resp []datamodels.Data
	var tmp datamodels.Data

	for rows.Next() {
		err = rows.Scan(&tmp.DataID, &tmp.Data, &tmp.Metadata, &tmp.Deleted, &tmp.ChangedAt)
		if err == nil {
			resp = append(resp, tmp)
		}
	}
	return resp, nil
}

// ClientSync - synchronize client data with server
func (dbs *DBStorage) ClientSync(userID uint32, data []*pb.Data) error {
	query := `insert into keeper (data_id,user_id, data_info,meta_info, changed_at) values ($1, $2,$3,$4,$5) ON CONFLICT (user_id, data_id) DO UPDATE SET data_info=EXCLUDED.data_info, meta_info=EXCLUDED.meta_info, changed_at=EXCLUDED.changed_at where keeper.changed_at < $5;`
	for i := range data {
		_, err := dbs.db.Exec(query, data[i].DataId, userID, data[i].Data, data[i].MetaInfo, data[i].ChangedAt.AsTime().Format(time.RFC3339))
		if err != nil {
			return ErrInternal
		}
	}
	return nil
}
