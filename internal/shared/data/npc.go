package data

import (
	"encoding/binary"
	"os"
)

type NPCAttack struct {
	Range            byte
	Area             byte
	Damage           uint16
	AdditionalDamage uint16
}

type NPCData struct {
	Name                [0x14]byte
	Id                  uint16
	RespawnRate         uint16
	AttackTypeInfo      byte
	TargetSelectionInfo byte
	Defense             byte
	AdditionalDefense   byte
	Attacks             [0x3]NPCAttack
	AttackSpeedLow      uint16
	AttackSpeedHigh     uint16
	AttackSpeed         uint16
	Level               byte
	PlayerExp           uint16
	Appearance          byte
	HP                  uint32
	BlueAttackDefense   uint16
	RedAttackDefense    uint16
	GreyAttackDefense   uint16
	MercenaryExp        uint16
}

func LoadNPCData(npcFilePath string) (*NPCData, error) {
	npcFile, err := os.Open(npcFilePath)
	if err != nil {
		return nil, err
	}

	defer npcFile.Close()
	npcData := NPCData{}
	err = binary.Read(npcFile, binary.LittleEndian, &npcData)
	if err != nil {
		return nil, err
	}

	return &npcData, nil
}
