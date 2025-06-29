package messages

import (
	"bytes"
	"encoding/binary"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages/protocol"
	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
)

type CharacterWear struct {
	ItemPtr    uint32
	ItemCode   uint32
	ItemOption uint32
	WearIndex  uint32
}

type CharacterInfo struct {
	Name     [0x15]byte
	LastUsed byte
	Class    byte
	Town     byte
	Level    uint32
	Wear     [0xA]CharacterWear
}

type MsgS2CCharacterList struct {
	MsgHead
	CharacterList [0x5]CharacterInfo
}

func (msg *MsgS2CCharacterList) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgS2CCharacterList) SetSize() {
	msg.Size = msg.GetSize()
}

func (msg *MsgS2CCharacterList) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgS2CCharacterList(pcId uint32, characterList []CharacterInfo) *MsgS2CCharacterList {
	msgS2CCharacterList := &MsgS2CCharacterList{
		MsgHead: MsgHead{
			Protocol: protocol.S2CCharacterList,
			MsgHeadNoProtocol: MsgHeadNoProtocol{
				Ctrl: 0x03,
				Cmd:  0xFF,
				PcId: pcId,
			},
		},
		CharacterList: [5]CharacterInfo{},
	}

	for i := 0; i < 5; i++ {
		if i < len(characterList) {
			msgS2CCharacterList.CharacterList[i] = characterList[i]
		} else {
			msgS2CCharacterList.CharacterList[i].Class = 255
		}
	}

	msgS2CCharacterList.SetSize()
	return msgS2CCharacterList
}

func NewMsgS2CCharacterListEmpty(pcId uint32) *MsgS2CCharacterList {
	msgS2CCharacterList := &MsgS2CCharacterList{
		MsgHead: MsgHead{
			Protocol: protocol.S2CCharacterList,
			MsgHeadNoProtocol: MsgHeadNoProtocol{
				Ctrl: 0x03,
				Cmd:  0xFF,
				PcId: pcId,
			},
		},
		CharacterList: [5]CharacterInfo{},
	}

	for i := range msgS2CCharacterList.CharacterList {
		msgS2CCharacterList.CharacterList[i].Class = 255
	}

	msgS2CCharacterList.SetSize()
	return msgS2CCharacterList
}

func ReadMsgS2CCharacterList(packet []byte) (*MsgS2CCharacterList, error) {
	var msg MsgS2CCharacterList
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

type MsgS2CAnsCreatePlayer struct {
	MsgHead
	Class byte
	Name  [0x15]byte
	Wear  [0xA]CharacterWear
}

func (msg *MsgS2CAnsCreatePlayer) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgS2CAnsCreatePlayer) SetSize() {
	msg.Size = msg.GetSize()
}

func (msg *MsgS2CAnsCreatePlayer) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgS2CAnsCreatePlayer(pcId uint32, class byte, name string, wear [0xA]CharacterWear) *MsgS2CAnsCreatePlayer {
	msg := MsgS2CAnsCreatePlayer{
		MsgHead: MsgHead{
			Protocol:          protocol.S2CAnsCreatePlayer,
			MsgHeadNoProtocol: MsgHeadNoProtocol{Ctrl: 0x03, Cmd: 0x01, PcId: pcId},
		},
		Class: class,
		Wear:  wear,
	}

	copy(msg.Name[:], utils.MakeFixedLengthStringBytes(name, 0x15))
	msg.SetSize()
	return &msg
}

type MsgS2CAnsDeletePlayer struct {
	MsgHead
	Name [0x15]byte
}

func (msg *MsgS2CAnsDeletePlayer) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgS2CAnsDeletePlayer) SetSize() {
	msg.Size = msg.GetSize()
}

func (msg *MsgS2CAnsDeletePlayer) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgS2CAnsDeletePlayer(pcId uint32, name string) *MsgS2CAnsDeletePlayer {
	msg := MsgS2CAnsDeletePlayer{
		MsgHead: MsgHead{
			Protocol: protocol.S2CAnsDeletePlayer,
			MsgHeadNoProtocol: MsgHeadNoProtocol{
				Ctrl: 0x01,
				Cmd:  0x01,
				PcId: pcId,
			},
		},
	}
	copy(msg.Name[:], utils.MakeFixedLengthStringBytes(name, 0x15))
	msg.SetSize()
	return &msg
}

func ReadMsgS2CAnsDeletePlayer(packet []byte) (*MsgS2CAnsDeletePlayer, error) {
	var msg MsgS2CAnsDeletePlayer
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
