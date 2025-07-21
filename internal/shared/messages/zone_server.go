package messages

import (
	"bytes"
	"encoding/binary"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages/protocol"
	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
)

type MsgS2MWorldLogin struct {
	MsgHeadMs
	CharacterName [0x15]byte
}

func (msg *MsgS2MWorldLogin) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgS2MWorldLogin) SetSize() {
	msg.Size = uint16(msg.GetSize())
}

func (msg *MsgS2MWorldLogin) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgS2MWorldLogin(pcId uint32, characterName string) *MsgS2MWorldLogin {
	msg := MsgS2MWorldLogin{
		MsgHeadMs: MsgHeadMs{
			PcId:     pcId,
			Protocol: protocol.S2MWorldLogin,
		},
	}
	copy(msg.CharacterName[:], utils.MakeFixedLengthStringBytes(characterName, 0x15))
	msg.SetSize()
	return &msg
}

func ReadMsgS2MWorldLogin(packet []byte) (*MsgS2MWorldLogin, error) {
	var msg MsgS2MWorldLogin
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

type MsgS2CWorldLogin struct {
	MsgHead
	CharacterName         [0x15]byte
	Class                 byte
	Level                 uint16
	Exp                   uint32
	MapIndex              uint32
	MapCell               uint32
	Skill                 SkillInfo
	PKCount               uint32
	RTime                 uint32
	SocialInfo            SocialInfo
	Woonz                 uint32
	HPStore               uint32
	MPStore               uint32
	Lore                  uint32
	RemainingPoints       uint16
	Strength              uint16
	Intelligence          uint16
	Dexterity             uint16
	Vitality              uint16
	Mana                  uint16
	HPCapacity            uint32
	MPCapacity            uint32
	HP                    uint16
	MP                    uint16
	HitAttack             uint16
	MagicAttack           uint16
	Defense               uint16
	FireAttack            uint16
	FireDefence           uint16
	IceAttack             uint16
	IceDefense            uint16
	LightAttack           uint16
	LightDefense          uint16
	MaxHp                 uint16
	MaxMp                 uint16
	AdditionalHitAttack   uint16
	AdditionalMagicAttack uint16
	Unknown               uint16
	WearList              [0xA]CharacterWear
	CharacterInventory    [0x1E]CharacterInventory
	ActivePet             Pet
	PetInventory          [0x5]Pet
}

func (msg *MsgS2CWorldLogin) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgS2CWorldLogin) SetSize() {
	msg.Size = msg.GetSize()
}

func (msg *MsgS2CWorldLogin) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgS2CWorldLogin(pcId uint32, characterName string) *MsgS2CWorldLogin {
	msg := MsgS2CWorldLogin{
		MsgHead: MsgHead{
			MsgHeadNoProtocol: MsgHeadNoProtocol{
				PcId: pcId,
				Ctrl: 0x03,
				Cmd:  0xFF,
			},
			Protocol: protocol.S2CWorldLogin,
		},
	}

	copy(msg.CharacterName[:], utils.MakeFixedLengthStringBytes(characterName, 0x15))
	msg.SetSize()
	return &msg
}

func ReadMsgS2CWorldLogin(packet []byte) (*MsgS2CWorldLogin, error) {
	var msg MsgS2CWorldLogin
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
