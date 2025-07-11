package db

import (
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	sharedConstants "github.com/project-agonyl/open-agonyl-servers/internal/shared/constants"
	"golang.org/x/crypto/bcrypt"
)

type DBService interface {
	GetAccountByUsername(username string) (*Account, error)
	GetOrCreateAccount(username string, password string) (*Account, error)
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

func (s *dbService) GetAccountByUsername(username string) (*Account, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qb := psql.Select("id", "username", "password_hash", "status", "is_online").
		From("accounts").
		Where(sq.Eq{"username": username})

	query, args, err := qb.ToSql()
	if err != nil {
		s.logger.Error("Failed to build get account by username query", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	account := &Account{}
	err = s.db.Get(account, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		s.logger.Error("Failed to execute get account by username query", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	return account, nil
}

func (s *dbService) GetOrCreateAccount(username string, password string) (*Account, error) {
	account, err := s.GetAccountByUsername(username)
	if err == nil {
		return account, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qb := psql.Insert("accounts").
		Columns("username", "password_hash", "status", "email").
		Values(username, string(passwordHash), sharedConstants.AccountStatusActive, username+"@a3agonyl.com").
		Suffix("RETURNING id, username, password_hash, status, is_online")

	query, args, err := qb.ToSql()
	if err != nil {
		s.logger.Error("Failed to build create account query", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	account = &Account{}
	err = s.db.Get(account, query, args...)
	if err != nil {
		s.logger.Error("Failed to execute create account query", shared.Field{Key: "error", Value: err})
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
