package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/constants"
)

type DBService interface {
	GetCharacterMapInfo(accountID uint32, characterName string) (uint16, error)
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

func (s *dbService) GetCharacterMapInfo(accountID uint32, characterName string) (uint16, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qb := psql.Select("id", "name", "character_data").
		From("characters").
		Where(sq.And{sq.Eq{"account_id": accountID}, sq.Eq{"status": constants.CharacterStatusActive}, sq.Eq{"name": characterName}})

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	var character CharacterForMap
	err = s.db.Get(&character, query, args...)
	if err != nil {
		return 0, err
	}

	return character.Data.Location.MapCode, nil
}

type CharacterForMap struct {
	ID   uint32        `db:"id"`
	Name string        `db:"name"`
	Data CharacterData `db:"character_data"`
}

type CharacterData struct {
	Location Location `json:"location"`
}

func (c *CharacterData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("CharacterData: type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, c)
}

func (c *CharacterData) Value() (driver.Value, error) {
	return json.Marshal(c)
}

type Location struct {
	MapCode uint16 `json:"map_code"`
}
