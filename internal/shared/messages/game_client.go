package messages

import (
	"bytes"
	"encoding/binary"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages/protocol"
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

type MsgC2SGateLogin struct {
	MsgHeadNoProtocol
	PcId     uint32
	Account  [0x15]byte
	Password [0x15]byte
}

func (msg *MsgC2SGateLogin) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgC2SGateLogin) SetSize() {
	msg.Size = msg.GetSize()
}

func (msg *MsgC2SGateLogin) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgC2SGateLogin(pcId uint32, account string, password string) *MsgC2SGateLogin {
	msg := MsgC2SGateLogin{
		MsgHeadNoProtocol: MsgHeadNoProtocol{Ctrl: 0x01, Cmd: 0xE2, PcId: pcId},
		PcId:              pcId,
	}

	copy(msg.Account[:], utils.MakeFixedLengthStringBytes(account, 0x15))
	copy(msg.Password[:], utils.MakeFixedLengthStringBytes(password, 0x15))
	msg.SetSize()
	return &msg
}

func ReadMsgC2SGateLogin(packet []byte) (*MsgC2SGateLogin, error) {
	var msg MsgC2SGateLogin
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

type MsgC2SAskCreatePlayer struct {
	MsgHead
	Class byte
	Town  byte
	Name  [0x15]byte
}

func (msg *MsgC2SAskCreatePlayer) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgC2SAskCreatePlayer) SetSize() {
	msg.Size = msg.GetSize()
}

func (msg *MsgC2SAskCreatePlayer) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgC2SAskCreatePlayer(pcId uint32, class byte, town byte, name string) *MsgC2SAskCreatePlayer {
	msg := MsgC2SAskCreatePlayer{
		MsgHead: MsgHead{
			Protocol: protocol.C2SAskCreatePlayer,
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

func ReadMsgC2SAskCreatePlayer(packet []byte) (*MsgC2SAskCreatePlayer, error) {
	var msg MsgC2SAskCreatePlayer
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

type MsgC2SAskDeletePlayer struct {
	MsgHead
	Name [0x15]byte
}

func (msg *MsgC2SAskDeletePlayer) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgC2SAskDeletePlayer) SetSize() {
	msg.Size = msg.GetSize()
}

func (msg *MsgC2SAskDeletePlayer) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgC2SAskDeletePlayer(pcId uint32, name string) *MsgC2SAskDeletePlayer {
	msg := MsgC2SAskDeletePlayer{
		MsgHead: MsgHead{
			Protocol: protocol.C2SAskDeletePlayer,
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

func ReadMsgC2SAskDeletePlayer(packet []byte) (*MsgC2SAskDeletePlayer, error) {
	var msg MsgC2SAskDeletePlayer
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
