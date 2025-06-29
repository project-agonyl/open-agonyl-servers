package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
)

type DBService interface {
	GetAccount(id uint32) (*Account, error)
	Close() error
}

type dbService struct {
	db     *sqlx.DB
	logger shared.Logger
}

func NewDbService(dbUrl string, logger shared.Logger) (DBService, error) {
	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	return &dbService{
		db:     db,
		logger: logger,
	}, nil
}

func (s *dbService) Close() error {
	return s.db.Close()
}

func (s *dbService) GetAccount(id uint32) (*Account, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qb := psql.Select("id", "username", "password_hash", "status", "is_online").
		From("accounts").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"status": "active"}})

	query, args, err := qb.ToSql()
	if err != nil {
		s.logger.Error("Failed to build get account query", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	account := &Account{}
	err = s.db.Get(account, query, args...)
	if err != nil {
		s.logger.Error("Failed to execute get account query", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	return account, nil
}

type Account struct {
	ID           uint32 `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	Status       string `db:"status"`
	IsOnline     bool   `db:"is_online"`
}
