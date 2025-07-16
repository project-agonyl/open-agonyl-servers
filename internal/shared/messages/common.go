package messages

import (
	"bytes"
	"encoding/binary"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages/protocol"
	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
)

type Msg interface {
	GetSize() uint32
	SetSize()
	GetBytes() []byte
}

type MsgHeadNoProtocol struct {
	Size uint32
	PcId uint32
	Ctrl byte
	Cmd  byte
}

type MsgHead struct {
	MsgHeadNoProtocol
	Protocol uint16
}

type MsgHeadMs struct {
	Protocol     uint16
	Size         uint16
	PcId         uint32
	GateServerId byte
}

type MsgS2CError struct {
	MsgHead
	Code uint16
	Msg  [64]byte
}

func (msg *MsgS2CError) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgS2CError) SetSize() {
	msg.Size = msg.GetSize()
}

func (msg *MsgS2CError) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgS2CError(pcId uint32, code uint16, msg string) *MsgS2CError {
	msgS2CError := MsgS2CError{
		MsgHead: MsgHead{Protocol: protocol.S2CError, MsgHeadNoProtocol: MsgHeadNoProtocol{Ctrl: 0x03, Cmd: 0xFF}},
		Code:    code,
	}

	copy(msgS2CError.Msg[:], msg)
	msgS2CError.PcId = pcId
	msgS2CError.SetSize()
	return &msgS2CError
}

func ReadMsgS2CError(packet []byte) (*MsgS2CError, error) {
	var msg MsgS2CError
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

type MsgS2GZoneChange struct {
	MsgHeadNoProtocol
	ZoneId byte
}

func (msg *MsgS2GZoneChange) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgS2GZoneChange) SetSize() {
	msg.Size = msg.GetSize()
}

func (msg *MsgS2GZoneChange) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgS2GZoneChange(pcId uint32, zoneId byte) MsgS2GZoneChange {
	msg := MsgS2GZoneChange{
		MsgHeadNoProtocol: MsgHeadNoProtocol{Ctrl: 0x01, Cmd: 0xE1, PcId: pcId},
		ZoneId:            zoneId,
	}

	msg.SetSize()
	return msg
}

func ReadMsgS2GZoneChange(packet []byte) (*MsgS2GZoneChange, error) {
	var msg MsgS2GZoneChange
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

type MsgS2MCharacterLogout struct {
	MsgHeadMs
	CharacterName [0x15]byte
}

func (msg *MsgS2MCharacterLogout) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgS2MCharacterLogout) SetSize() {
	msg.Size = uint16(msg.GetSize())
}

func (msg *MsgS2MCharacterLogout) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgS2MCharacterLogout(pcId uint32, characterName string) *MsgS2MCharacterLogout {
	msg := MsgS2MCharacterLogout{
		MsgHeadMs: MsgHeadMs{
			Protocol: protocol.S2MCharacterLogout,
			PcId:     pcId,
		},
	}

	copy(msg.CharacterName[:], utils.MakeFixedLengthStringBytes(characterName, 0x15))
	msg.SetSize()
	return &msg
}

func ReadMsgS2MCharacterLogout(packet []byte) (*MsgS2MCharacterLogout, error) {
	var msg MsgS2MCharacterLogout
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
