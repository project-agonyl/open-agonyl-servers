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
