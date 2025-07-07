package messages

import (
	"bytes"
	"encoding/binary"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages/protocol"
	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
)

type MsgM2SError struct {
	MsgHeadMs
	Code uint16
	Msg  [0x40]byte
}

func (msg *MsgM2SError) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgM2SError) SetSize() {
	msg.Size = uint16(msg.GetSize())
}

func (msg *MsgM2SError) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgM2SError(pcId uint32, code uint16, msg string, gateServerId byte) *MsgM2SError {
	msgM2SError := MsgM2SError{
		MsgHeadMs: MsgHeadMs{Protocol: 0xA000, GateServerId: gateServerId, PcId: pcId},
	}
	copy(msgM2SError.Msg[:], utils.MakeFixedLengthStringBytes(msg, 0x40))
	msgM2SError.Code = code
	msgM2SError.SetSize()
	return &msgM2SError
}

func ReadMsgM2SError(packet []byte) (*MsgM2SError, error) {
	var msg MsgM2SError
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

type MsgM2SAnsCharacterLogin struct {
	MsgHeadMs
	ZoneId byte
	MapId  uint16
}

func (msg *MsgM2SAnsCharacterLogin) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgM2SAnsCharacterLogin) SetSize() {
	msg.Size = uint16(msg.GetSize())
}

func (msg *MsgM2SAnsCharacterLogin) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgM2SAnsCharacterLogin(pcId uint32, serverId byte, mapId uint16, gateServerId byte) *MsgM2SAnsCharacterLogin {
	msgM2SAnsCharacterLogin := MsgM2SAnsCharacterLogin{
		MsgHeadMs: MsgHeadMs{Protocol: protocol.S2MCharacterLogin, GateServerId: gateServerId, PcId: pcId},
	}
	msgM2SAnsCharacterLogin.ZoneId = serverId
	msgM2SAnsCharacterLogin.MapId = mapId
	msgM2SAnsCharacterLogin.SetSize()
	return &msgM2SAnsCharacterLogin
}

func ReadMsgM2SAnsCharacterLogin(packet []byte) (*MsgM2SAnsCharacterLogin, error) {
	var msg MsgM2SAnsCharacterLogin
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
