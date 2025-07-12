package session

import (
	"time"

	"github.com/project-agonyl/open-agonyl-servers/internal/webserver/db"
)

type SessionStorage interface {
	GetSessionBySessionID(sessionId string, isActiveOnly bool) (*Session, error)
	GetSessionsByAccountID(accountID uint32, isActiveOnly bool) ([]Session, error)
	GetSessionsByUsername(username string, isActiveOnly bool) ([]Session, error)
	CreateSession(accountID uint32, userAgent string, ipAddress string, expiresAt time.Time) (string, error)
	RevokeSession(sessionID string) error
}

type dbSessionStorage struct {
	dbService db.DBService
}

func NewDBSessionStorage(dbService db.DBService) SessionStorage {
	return &dbSessionStorage{
		dbService: dbService,
	}
}

func (s *dbSessionStorage) GetSessionBySessionID(sessionId string, isActiveOnly bool) (*Session, error) {
	dbSession, err := s.dbService.GetSessionBySessionID(sessionId, isActiveOnly)
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:        dbSession.ID,
		AccountID: dbSession.AccountID,
		SessionID: dbSession.SessionID,
		UserAgent: dbSession.UserAgent,
		IPAddress: dbSession.IPAddress,
		IssuedAt:  dbSession.IssuedAt,
		ExpiresAt: dbSession.ExpiresAt,
		RevokedAt: dbSession.RevokedAt,
	}, nil
}

func (s *dbSessionStorage) GetSessionsByAccountID(accountID uint32, isActiveOnly bool) ([]Session, error) {
	dbSessions, err := s.dbService.GetSessionsByAccountID(accountID, isActiveOnly)
	if err != nil {
		return nil, err
	}

	sessions := make([]Session, len(dbSessions))
	for i, dbSession := range dbSessions {
		sessions[i] = Session{
			ID:        dbSession.ID,
			AccountID: dbSession.AccountID,
			SessionID: dbSession.SessionID,
			UserAgent: dbSession.UserAgent,
			IPAddress: dbSession.IPAddress,
			IssuedAt:  dbSession.IssuedAt,
			ExpiresAt: dbSession.ExpiresAt,
			RevokedAt: dbSession.RevokedAt,
		}
	}

	return sessions, nil
}

func (s *dbSessionStorage) GetSessionsByUsername(username string, isActiveOnly bool) ([]Session, error) {
	dbSessions, err := s.dbService.GetSessionsByUsername(username, isActiveOnly)
	if err != nil {
		return nil, err
	}

	sessions := make([]Session, len(dbSessions))
	for i, dbSession := range dbSessions {
		sessions[i] = Session{
			ID:        dbSession.ID,
			AccountID: dbSession.AccountID,
			SessionID: dbSession.SessionID,
			UserAgent: dbSession.UserAgent,
			IPAddress: dbSession.IPAddress,
			IssuedAt:  dbSession.IssuedAt,
			ExpiresAt: dbSession.ExpiresAt,
			RevokedAt: dbSession.RevokedAt,
		}
	}

	return sessions, nil
}

func (s *dbSessionStorage) CreateSession(accountID uint32, userAgent string, ipAddress string, expiresAt time.Time) (string, error) {
	return s.dbService.CreateSession(accountID, userAgent, ipAddress, expiresAt)
}

func (s *dbSessionStorage) RevokeSession(sessionID string) error {
	return s.dbService.RevokeSession(sessionID)
}

type Session struct {
	ID        uint32     `json:"id" db:"id"`
	AccountID uint32     `json:"account_id" db:"account_id"`
	SessionID string     `json:"session_id" db:"session_id"`
	UserAgent string     `json:"user_agent" db:"user_agent"`
	IPAddress string     `json:"ip_address" db:"ip_address"`
	IssuedAt  time.Time  `json:"issued_at" db:"issued_at"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at" db:"revoked_at"`
}

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func (s *Session) IsRevoked() bool {
	return s.RevokedAt != nil
}

func (s *Session) IsActive() bool {
	return !s.IsExpired() && !s.IsRevoked()
}
