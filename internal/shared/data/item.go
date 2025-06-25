package data

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
	Unknown1 uint16
	Row      uint16
	Slot     uint16
	Type     uint16
	Name     [32]byte
	NPCPrice uint32
	Unknown2 [9]uint16
	Levels   [10]IT0RawLevelProperties
}

type IT0ExRaw struct {
	Row    uint16
	Levels [5]IT0RawLevelProperties
}
