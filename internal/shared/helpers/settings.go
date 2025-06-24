package helpers

import (
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func GetSettingsByServerName(db *sqlx.DB, serverName string) (*ServerSettings, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qb := psql.Select("id", "server_name", "settings", "metadata").
		From("server_settings").
		Where(sq.Eq{"server_name": serverName})

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	settings := &ServerSettings{}
	err = db.Get(settings, sql, args...)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

type ServerSettings struct {
	ID         int             `db:"id"`
	ServerName string          `db:"server_name"`
	Settings   json.RawMessage `db:"settings"`
	Metadata   json.RawMessage `db:"metadata"`
}
