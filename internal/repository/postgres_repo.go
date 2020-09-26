package repository

import (
	_ "github.com/lib/pq"
	"gopkg.in/jackc/pgx.v2"
	"time"
)

const (
	OK           = 0
	CREATED      = 1
	DB_ERROR     = 2
	EMPTY_RESULT = 3
	FORBIDDEN    = 4
	CONFLICT     = 5
	WRONG_INPUT  = 6
)

type PostgresRepository struct {
	db           *pgx.ConnPool
	user         string
	password     string
	databaseName string
	host         string
	port         uint16
}

func NewRepository(user string, password string, dataBaseName string,
	host string, port uint16) *PostgresRepository {
	repo := new(PostgresRepository)
	repo.user = user
	repo.databaseName = dataBaseName
	repo.password = password
	repo.host = host
	repo.port = port

	return repo
}

func (repo *PostgresRepository) Start(maxConn, acquireTimeout int) error {
	conf := pgx.ConnConfig{
		Host:     repo.host,
		Port:     repo.port,
		User:     repo.user,
		Password: repo.password,
		Database: repo.databaseName,
	}

	duration := time.Duration(acquireTimeout) * time.Second
	poolConf := pgx.ConnPoolConfig{
		ConnConfig:     conf,
		MaxConnections: maxConn,
		AcquireTimeout: duration,
	}
	dataBase, err := pgx.NewConnPool(poolConf)
	if err != nil {
		return err
	}

	repo.db = dataBase
	return nil
}

func (repo *PostgresRepository) Close() {
	repo.db.Close()
}
