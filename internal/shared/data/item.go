package data

import (
	"encoding/binary"
	"os"

	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
)

type IT0Level struct {
	Level               byte
	AttributeRange      uint16
	Attribute           uint16
	Strength            uint16
	Intelligence        uint16
	Dexterity           uint16
	AdditionalAttribute uint16
	RedOption           uint16
	GreyOption          uint16
	BlueOption          uint16
}

type IT0Property struct {
	Levels []IT0Level `json:"levels"`
}

type IT1Property struct {
	RequiredLevel uint16
	Attribute     uint16
	RedOption     uint16
	GreyOption    uint16
	BlueOption    uint16
}

type IT2Property struct {
	RequiredLevel uint16
	SkillLevel    uint16
}

type Item struct {
	ItemCode    uint32
	SlotIndex   byte
	ItemName    string
	Itemtype    byte
	NPCPrice    uint32
	IT0Property *IT0Property
	IT1Property *IT1Property
	IT2Property *IT2Property
}

type IT0RawLevelProperties struct {
	AdditionalAttribute uint16
	Strength            uint16
	Dexterity           uint16
	Intelligence        uint16
	Attribute           uint16
	Range               uint16
	BlueOption          uint16
	RedOption           uint16
	GreyOption          uint16
}

type IT0Raw struct {
	ItemCodeBase uint16
	Row          uint16
	Slot         uint16
	Type         uint16
	Name         [32]byte
	NPCPrice     uint32
	Unknown2     [9]uint16
	Levels       [10]IT0RawLevelProperties
}

type IT0ExRaw struct {
	Row    uint16
	Levels [5]IT0RawLevelProperties
}

func LoadIT0Items(it0FilePath string, it0ExFilePath string) ([]Item, error) {
	it0File, err := os.Open(it0FilePath)
	if err != nil {
		return nil, err
	}

	defer it0File.Close()

	it0ExFile, err := os.Open(it0ExFilePath)
	if err != nil {
		return nil, err
	}

	defer it0ExFile.Close()

	it0FileStat, err := it0File.Stat()
	if err != nil {
		return nil, err
	}

	totalItems := it0FileStat.Size() / 242
	items := make([]Item, totalItems)
	for i := int64(0); i < totalItems; i++ {
		it0Raw := IT0Raw{}
		err = binary.Read(it0File, binary.LittleEndian, &it0Raw)
		if err != nil {
			return nil, err
		}

		it0property := &IT0Property{}
		it0property.Levels = make([]IT0Level, 10)
		for j, property := range it0Raw.Levels {
			it0property.Levels[j] = IT0Level{
				Level:               byte(j + 1),
				AttributeRange:      property.Range,
				Attribute:           property.Attribute,
				Strength:            property.Strength,
				Intelligence:        property.Intelligence,
				Dexterity:           property.Dexterity,
				AdditionalAttribute: property.AdditionalAttribute,
				RedOption:           property.RedOption,
				GreyOption:          property.GreyOption,
			}
		}

		items[it0Raw.Row] = Item{
			ItemCode:    uint32((it0Raw.ItemCodeBase << 10) + it0Raw.Row),
			SlotIndex:   byte(it0Raw.Slot),
			ItemName:    utils.ReadStringFromBytes(it0Raw.Name[:]),
			Itemtype:    byte(it0Raw.Type),
			NPCPrice:    it0Raw.NPCPrice,
			IT0Property: it0property,
		}
	}

	it0ExFileStat, err := it0ExFile.Stat()
	if err != nil {
		return nil, err
	}

	totalItemsEx := it0ExFileStat.Size() / 92
	for i := int64(0); i < totalItemsEx; i++ {
		it0ExRaw := IT0ExRaw{}
		err = binary.Read(it0ExFile, binary.LittleEndian, &it0ExRaw)
		if err != nil {
			return nil, err
		}

		for j, property := range it0ExRaw.Levels {
			items[it0ExRaw.Row].IT0Property.Levels = append(
				items[it0ExRaw.Row].IT0Property.Levels, IT0Level{
					Level:               byte(j + 1),
					AttributeRange:      property.Range,
					Attribute:           property.Attribute,
					Strength:            property.Strength,
					Intelligence:        property.Intelligence,
					Dexterity:           property.Dexterity,
					AdditionalAttribute: property.AdditionalAttribute,
					RedOption:           property.RedOption,
					GreyOption:          property.GreyOption,
				})
		}
	}

	return items, nil
}
