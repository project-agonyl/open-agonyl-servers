package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/constants"
)

type DBService interface {
	GetCharacter(id uint32, name string) (*Character, error)
	GetDB() *sqlx.DB
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

func (s *dbService) GetDB() *sqlx.DB {
	return s.db
}

func (s *dbService) Close() error {
	return s.db.Close()
}

func (s *dbService) GetCharacter(id uint32, name string) (*Character, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	qb := psql.Select(
		"characters.id",
		"characters.name",
		"characters.class",
		"characters.level",
		"characters.character_data",
		"accounts.username as account",
	).
		From("characters").
		LeftJoin("accounts ON accounts.id = characters.account_id").
		LeftJoin("characters_accounts ON characters_accounts.character_id = characters.id").
		Where(sq.And{
			sq.Eq{"account_id": id},
			sq.Eq{"status": constants.CharacterStatusActive},
			sq.Eq{"name": name},
		}).
		Limit(1)

	query, args, err := qb.ToSql()
	if err != nil {
		s.logger.Error("Failed to build get character query", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	character := &Character{}
	err = s.db.Get(character, query, args...)
	if err != nil {
		s.logger.Error("Failed to execute get character query", shared.Field{Key: "error", Value: err})
		return nil, err
	}

	return character, nil
}

type Character struct {
	ID      uint32        `db:"id"`
	Name    string        `db:"name"`
	Class   byte          `db:"class"`
	Level   uint16        `db:"level"`
	Account string        `db:"account"`
	Data    CharacterData `db:"character_data"`
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
	SkillID byte `json:"skill_id"`
	Level   byte `json:"level"`
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
