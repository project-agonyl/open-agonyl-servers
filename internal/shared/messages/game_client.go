package messages

import (
	"bytes"
	"encoding/binary"

	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
)

type MsgC2SLogin struct {
	MsgHeadNoProtocol
	Username [0x15]byte
	Password [0x15]byte
}

func (msg *MsgC2SLogin) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgC2SLogin) SetSize() {
	msg.Size = msg.GetSize()
}

func (msg *MsgC2SLogin) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgC2SLogin(username, password string) *MsgC2SLogin {
	msg := MsgC2SLogin{
		MsgHeadNoProtocol: MsgHeadNoProtocol{Ctrl: 0x01, Cmd: 0x01},
	}
	copy(msg.Username[:], utils.MakeFixedLengthStringBytes(username, 0x15))
	copy(msg.Password[:], utils.MakeFixedLengthStringBytes(password, 0x15))
	msg.SetSize()
	return &msg
}

func ReadMsgC2SLogin(packet []byte) (*MsgC2SLogin, error) {
	var msg MsgC2SLogin
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
