package storage

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"gophkeeper/internal/datamodels"
	"time"
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(path string) (Storage, error) {
	if path == "" {
		return nil, errors.New("invalid db address")
	}
	db, err := sql.Open("pgx",
		path)
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
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
	query := `insert into keeper (id,user_id, data_info,meta_info, changed_at)where changed_at<$5  values ($1, $2,$3,$4,$5)
		ON CONFLICT ON CONSTRAINT data_id DO UPDATE SET data_info=EXCLUDED.data_info, meta_info=EXCLUDED.meta_info, changed_at=EXCLUDED.changed_at;`
	_, err := dbs.db.Exec(query, data.DataID, data.UserID, data.Data, data.Metadata, data.ChangedAt.Format(time.RFC3339))
	if err != nil {
		return ErrInternal
	}
	return nil
}
func (dbs *DBStorage) GetData(dataID uint32, userID uint32) (datamodels.Data, error) {
	rows := dbs.db.QueryRow("select data_info,meta_info from keeper where id=$1 and user_id=$2limit 1;", dataID, userID)
	var v datamodels.Data
	err := rows.Scan(&v.Data, &v.Metadata)
	if err != nil {
		return datamodels.Data{}, ErrNotFound
	}
	return v, nil
}
func (dbs *DBStorage) DelData(dataID uint32, userID uint32) error {
	_, err := dbs.db.Exec("DELETE from keeper where id=$1 and user_id=$2;", dataID, userID)
	if err != nil {
		return ErrInternal
	}
	return nil
}
