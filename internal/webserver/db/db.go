package db

import (
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
)

type DBService interface {
	GetAccountByUsername(username string) (*Account, error)
	GetSessionBySessionID(sessionID string, isActiveOnly bool) (*Session, error)
	GetSessionsByAccountID(accountID uint32, isActiveOnly bool) ([]Session, error)
	GetSessionsByUsername(username string, isActiveOnly bool) ([]Session, error)
	CreateSession(accountID uint32, userAgent string, ipAddress string, expiresAt time.Time) (string, error)
	RevokeSession(sessionID string) error
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
	qb := psql.Select("id", "account_id", "username", "password_hash", "status", "is_online").
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

type Account struct {
	ID           uint32 `db:"id"`
	AccountID    string `db:"account_id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	Status       string `db:"status"`
	IsOnline     bool   `db:"is_online"`
}

type Session struct {
	ID        uint32     `db:"id"`
	SessionID string     `db:"session_id"`
	AccountID uint32     `db:"account_id"`
	UserAgent string     `db:"user_agent"`
	IPAddress string     `db:"ip_address"`
	IssuedAt  time.Time  `db:"issued_at"`
	ExpiresAt time.Time  `db:"expires_at"`
	RevokedAt *time.Time `db:"revoked_at"`
	Metadata  *string    `db:"metadata"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}

func (s *dbService) GetSessionBySessionID(sessionID string, isActiveOnly bool) (*Session, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qb := psql.Select("id", "session_id", "account_id", "user_agent", "ip_address", "issued_at", "expires_at", "revoked_at", "metadata", "created_at", "updated_at").
		From("web_sessions").
		Where(sq.Eq{"session_id": sessionID})

	if isActiveOnly {
		qb = qb.Where(sq.And{
			sq.Eq{"revoked_at": nil},
			sq.Gt{"expires_at": time.Now()},
		})
	}

	query, args, err := qb.ToSql()
	if err != nil {
		s.logger.Error("Failed to build get session by session ID query", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	session := &Session{}
	err = s.db.Get(session, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		s.logger.Error("Failed to execute get session by session ID query", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	return session, nil
}

func (s *dbService) GetSessionsByAccountID(accountID uint32, isActiveOnly bool) ([]Session, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qb := psql.Select("id", "session_id", "account_id", "user_agent", "ip_address", "issued_at", "expires_at", "revoked_at", "metadata", "created_at", "updated_at").
		From("web_sessions").
		Where(sq.Eq{"account_id": accountID}).
		OrderBy("issued_at DESC")

	if isActiveOnly {
		qb = qb.Where(sq.And{
			sq.Eq{"revoked_at": nil},
			sq.Gt{"expires_at": time.Now()},
		})
	}

	query, args, err := qb.ToSql()
	if err != nil {
		s.logger.Error("Failed to build get sessions by account ID query", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	var sessions []Session
	err = s.db.Select(&sessions, query, args...)
	if err != nil {
		s.logger.Error("Failed to execute get sessions by account ID query", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	return sessions, nil
}

func (s *dbService) GetSessionsByUsername(username string, isActiveOnly bool) ([]Session, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qb := psql.Select("ws.id", "ws.session_id", "ws.account_id", "ws.user_agent", "ws.ip_address", "ws.issued_at", "ws.expires_at", "ws.revoked_at", "ws.metadata", "ws.created_at", "ws.updated_at").
		From("web_sessions ws").
		Join("accounts a ON ws.account_id = a.id").
		Where(sq.Eq{"a.username": username}).
		OrderBy("ws.issued_at DESC")

	if isActiveOnly {
		qb = qb.Where(sq.And{
			sq.Eq{"ws.revoked_at": nil},
			sq.Gt{"ws.expires_at": time.Now()},
		})
	}

	query, args, err := qb.ToSql()
	if err != nil {
		s.logger.Error("Failed to build get sessions by username query", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	var sessions []Session
	err = s.db.Select(&sessions, query, args...)
	if err != nil {
		s.logger.Error("Failed to execute get sessions by username query", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	return sessions, nil
}

func (s *dbService) CreateSession(accountID uint32, userAgent string, ipAddress string, expiresAt time.Time) (string, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qb := psql.Insert("web_sessions").
		Columns("account_id", "user_agent", "ip_address", "expires_at").
		Values(accountID, userAgent, ipAddress, expiresAt).
		Suffix("RETURNING session_id")

	query, args, err := qb.ToSql()
	if err != nil {
		s.logger.Error("Failed to build create session query", shared.Field{Key: "error", Value: err})
		return "", err
	}

	var sessionID string
	err = s.db.QueryRow(query, args...).Scan(&sessionID)
	if err != nil {
		s.logger.Error("Failed to execute create session query", shared.Field{Key: "error", Value: err})
		return "", err
	}

	return sessionID, nil
}

func (s *dbService) RevokeSession(sessionID string) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qb := psql.Update("web_sessions").
		Set("revoked_at", time.Now()).
		Where(sq.Eq{"session_id": sessionID})

	query, args, err := qb.ToSql()
	if err != nil {
		s.logger.Error("Failed to build revoke session query", shared.Field{Key: "error", Value: err})
		return err
	}

	result, err := s.db.Exec(query, args...)
	if err != nil {
		s.logger.Error("Failed to execute revoke session query", shared.Field{Key: "error", Value: err})
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("Failed to get rows affected for revoke session", shared.Field{Key: "error", Value: err})
		return err
	}

	if rowsAffected == 0 {
		return errors.New("session not found")
	}

	return nil
}
