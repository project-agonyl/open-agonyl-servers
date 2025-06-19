package messages

import (
	"bytes"
	"encoding/binary"

	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
)

type MsgLs2ClSay struct {
	MsgHeadNoProtocol
	Type  byte
	Words [0x51]byte
}

func (msg *MsgLs2ClSay) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgLs2ClSay) SetSize() {
	msg.Size = msg.GetSize()
}

func NewMsgLs2ClSay(words string) *MsgLs2ClSay {
	msg := MsgLs2ClSay{
		MsgHeadNoProtocol: MsgHeadNoProtocol{Ctrl: 0x01, Cmd: 0xE0},
		Type:              0x00,
	}
	copy(msg.Words[:], utils.MakeFixedLengthStringBytes(words, 0x51))
	msg.SetSize()
	return &msg
}

func ReadMsgLs2ClSay(packet []byte) (*MsgLs2ClSay, error) {
	var msg MsgLs2ClSay
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
