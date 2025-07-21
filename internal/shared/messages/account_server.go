package messages

import (
	"bytes"
	"encoding/binary"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages/protocol"
	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
)

type CharacterInfo struct {
	Name     [0x15]byte
	SlotUsed byte
	Class    byte
	Nation   byte
	Level    uint32
	Wear     [0xA]AclCharacterWear
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
	Wear  [0xA]AclCharacterWear
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

func NewMsgS2CAnsCreatePlayer(pcId uint32, class byte, name string, wear [0xA]AclCharacterWear) *MsgS2CAnsCreatePlayer {
	msg := MsgS2CAnsCreatePlayer{
		MsgHead: MsgHead{
			Protocol:          protocol.S2CAnsCreatePlayer,
			MsgHeadNoProtocol: MsgHeadNoProtocol{Ctrl: 0x03, Cmd: 0xFF, PcId: pcId},
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
				Ctrl: 0x03,
				Cmd:  0xFF,
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

type MsgS2MCharacterLogin struct {
	MsgHeadMs
	Account       [0x15]byte
	Password      [0x15]byte
	CharacterName [0x15]byte
	ClientIp      [0x10]byte
	Unknown       [0x4E]byte
}

func (msg *MsgS2MCharacterLogin) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgS2MCharacterLogin) SetSize() {
	msg.Size = uint16(msg.GetSize())
}

func (msg *MsgS2MCharacterLogin) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgS2MCharacterLogin(pcId uint32, account string, password string, characterName string, clientIp string, gateServerId byte) *MsgS2MCharacterLogin {
	msgS2MCharLogin := MsgS2MCharacterLogin{
		MsgHeadMs: MsgHeadMs{Protocol: protocol.S2MCharacterLogin, GateServerId: gateServerId, PcId: pcId},
	}
	copy(msgS2MCharLogin.Account[:], utils.MakeFixedLengthStringBytes(account, 0x15))
	copy(msgS2MCharLogin.Password[:], utils.MakeFixedLengthStringBytes(password, 0x15))
	copy(msgS2MCharLogin.CharacterName[:], utils.MakeFixedLengthStringBytes(characterName, 0x15))
	copy(msgS2MCharLogin.ClientIp[:], utils.MakeFixedLengthStringBytes(clientIp, 0x10))
	msgS2MCharLogin.SetSize()
	return &msgS2MCharLogin
}

func ReadMsgS2MCharacterLogin(packet []byte) (*MsgS2MCharacterLogin, error) {
	var msg MsgS2MCharacterLogin
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
