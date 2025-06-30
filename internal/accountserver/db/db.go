package db

import (
	"database/sql"
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
	GetAccount(id uint32) (*Account, error)
	GetCharactersForListing(accountID uint32) ([]CharacterForListing, error)
	DoesCharacterExist(name string) (bool, error)
	GetCharacterCount(accountID uint32) (int, error)
	CreateCharacter(accountID uint32, name string, class byte, characterData []byte) (uint32, error)
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

func (s *dbService) GetCharactersForListing(accountID uint32) ([]CharacterForListing, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qb := psql.Select("id", "name", "class", "level", "character_data").
		From("characters").
		Where(sq.And{sq.Eq{"account_id": accountID}, sq.Eq{"status": constants.CharacterStatusActive}}).
		OrderBy("last_login DESC")

	query, args, err := qb.ToSql()
	if err != nil {
		s.logger.Error("Failed to build get characters for listing query", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	characters := []CharacterForListing{}
	err = s.db.Select(&characters, query, args...)
	if err != nil {
		s.logger.Error("Failed to execute get characters for listing query", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	return characters, nil
}

func (s *dbService) DoesCharacterExist(name string) (bool, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qb := psql.Select("id").
		From("characters").
		Where(sq.Eq{"name": name}).
		Limit(1)

	query, args, err := qb.ToSql()
	if err != nil {
		s.logger.Error("Failed to build does character exist query", shared.Field{Key: "error", Value: err})
		return false, err
	}

	count := 0
	err = s.db.Get(&count, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		s.logger.Error("Failed to execute does character exist query", shared.Field{Key: "error", Value: err})
		return false, err
	}

	return count > 0, nil
}

func (s *dbService) GetCharacterCount(accountID uint32) (int, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qb := psql.Select("COUNT(*)").
		From("characters").
		Where(sq.And{sq.Eq{"account_id": accountID}, sq.Eq{"status": constants.CharacterStatusActive}})

	query, args, err := qb.ToSql()
	if err != nil {
		s.logger.Error("Failed to build get character count query", shared.Field{Key: "error", Value: err})
		return 0, err
	}

	count := 0
	err = s.db.Get(&count, query, args...)
	if err != nil {
		s.logger.Error("Failed to execute get character count query", shared.Field{Key: "error", Value: err})
		return 0, err
	}

	return count, nil
}

func (s *dbService) CreateCharacter(accountID uint32, name string, class byte, characterData []byte) (uint32, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qb := psql.Insert("characters").
		Columns("account_id", "name", "class", "character_data").
		Values(accountID, name, class, characterData).
		Suffix("RETURNING id")

	query, args, err := qb.ToSql()
	if err != nil {
		s.logger.Error("Failed to build create character query", shared.Field{Key: "error", Value: err})
		return 0, err
	}

	id := uint32(0)
	err = s.db.Get(&id, query, args...)
	if err != nil {
		s.logger.Error("Failed to execute create character query", shared.Field{Key: "error", Value: err})
		return 0, err
	}

	return id, nil
}

type Account struct {
	ID           uint32 `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	Status       string `db:"status"`
	IsOnline     bool   `db:"is_online"`
}

type CharacterForListing struct {
	ID    uint32        `db:"id"`
	Name  string        `db:"name"`
	Class byte          `db:"class"`
	Level uint32        `db:"level"`
	Data  CharacterData `db:"character_data"`
}

type CharacterData struct {
	Parole       uint32          `json:"parole"`
	SocialInfo   SocialInfo      `json:"social_info"`
	Wear         []WearItem      `json:"wear"`
	Inventory    []InventoryItem `json:"inventory"`
	Lore         uint32          `json:"lore"`
	Location     Location        `json:"location"`
	CurrentQuest QuestInfo       `json:"current_quest"`
	Skills       []SkillInfo     `json:"skills"`
	Stats        Stats           `json:"stats"`
	NPCFavors    []NPCFavor      `json:"npc_favors"`
	ActivePet    Pet             `json:"active_pet"`
	PetInventory []PetInventory  `json:"pet_inventory"`
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

type SocialInfo struct {
	Nation  byte `json:"nation"`
	KHIndex byte `json:"kh_index"`
}

type WearItem struct {
	ItemCode       uint32 `json:"item_code"`
	ItemOption     uint32 `json:"item_option"`
	ItemUniqueCode uint32 `json:"item_unique_code"`
}

type InventoryItem struct {
	ItemCode       uint32 `json:"item_code"`
	ItemOption     uint32 `json:"item_option"`
	ItemUniqueCode uint32 `json:"item_unique_code"`
	Slot           byte   `json:"slot"`
}

type Location struct {
	MapCode  uint16   `json:"map_code"`
	Position Position `json:"position"`
}

type Position struct {
	X byte `json:"x"`
	Y byte `json:"y"`
}

type QuestInfo struct {
	QuestID   uint32 `json:"quest_id"`
	Progress1 uint32 `json:"progress1"`
	Progress2 uint32 `json:"progress2"`
	Progress3 uint32 `json:"progress3"`
}

type SkillInfo struct {
	SkillID uint32 `json:"skill_id"`
	Level   uint32 `json:"level"`
}

type Stats struct {
	Strength        uint16 `json:"strength"`
	Intelligence    uint16 `json:"intelligence"`
	Dexterity       uint16 `json:"dexterity"`
	Vitality        uint16 `json:"vitality"`
	Mana            uint16 `json:"mana"`
	RemainingPoints uint16 `json:"remaining_points"`
	HP              uint16 `json:"hp"`
	MP              uint16 `json:"mp"`
	HPCapacity      uint16 `json:"hp_capacity"`
	MPCapacity      uint16 `json:"mp_capacity"`
}

type NPCFavor struct {
	NPCID uint32 `json:"npc_id"`
	Favor uint16 `json:"favor"`
}

type Pet struct {
	PetCode       uint32 `json:"pet_code"`
	PetHP         uint32 `json:"pet_hp"`
	PetOption     uint32 `json:"pet_option"`
	PetUniqueCode uint32 `json:"pet_unique_code"`
}

type PetInventory struct {
	PetCode       uint32 `json:"pet_code"`
	PetHP         uint32 `json:"pet_hp"`
	PetOption     uint32 `json:"pet_option"`
	PetUniqueCode uint32 `json:"pet_unique_code"`
	Slot          byte   `json:"slot"`
}
