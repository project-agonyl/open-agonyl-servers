package messages

import (
	"bytes"
	"encoding/binary"

	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
)

type MsgGate2LsConnect struct {
	MsgHeadNoProtocol
	ServerId  byte
	AgentId   byte
	IpAddress [0x10]byte
	Port      uint32
	Name      [0x11]byte
}

func (msg *MsgGate2LsConnect) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgGate2LsConnect) SetSize() {
	msg.Size = msg.GetSize()
}

func (msg *MsgGate2LsConnect) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgGate2LsConnect(serverId byte, agentId byte, ipAddress string, port uint32, name string) *MsgGate2LsConnect {
	msg := MsgGate2LsConnect{
		MsgHeadNoProtocol: MsgHeadNoProtocol{Ctrl: 0x02, Cmd: 0xE0},
		ServerId:          serverId,
		AgentId:           agentId,
		Port:              port,
	}
	copy(msg.IpAddress[:], utils.MakeFixedLengthStringBytes(ipAddress, 0x10))
	copy(msg.Name[:], utils.MakeFixedLengthStringBytes(name, 0x11))
	msg.SetSize()
	return &msg
}

func ReadMsgGate2LsConnect(packet []byte) (*MsgGate2LsConnect, error) {
	var msg MsgGate2LsConnect
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

type MsgGate2LsAccLogout struct {
	MsgHeadNoProtocol
	Reason     byte
	Account    [0x15]byte
	LogoutDate [0x09]byte
	LogoutTime [0x07]byte
}

func (msg *MsgGate2LsAccLogout) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgGate2LsAccLogout) SetSize() {
	msg.Size = msg.GetSize()
}

func (msg *MsgGate2LsAccLogout) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgGate2LsAccLogout(reason byte, account string) *MsgGate2LsAccLogout {
	msg := MsgGate2LsAccLogout{
		MsgHeadNoProtocol: MsgHeadNoProtocol{Ctrl: 0x02, Cmd: 0xE2},
		Reason:            reason,
	}
	copy(msg.Account[:], utils.MakeFixedLengthStringBytes(account, 0x15))
	msg.SetSize()
	return &msg
}

func ReadMsgGate2LsAccLogout(packet []byte) (*MsgGate2LsAccLogout, error) {
	var msg MsgGate2LsAccLogout
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

type MsgGate2LsPreparedAccLogin struct {
	MsgHeadNoProtocol
	Account [0x15]byte
}

func (msg *MsgGate2LsPreparedAccLogin) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgGate2LsPreparedAccLogin) SetSize() {
	msg.Size = msg.GetSize()
}

func (msg *MsgGate2LsPreparedAccLogin) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgGate2LsPreparedAccLogin(account string) *MsgGate2LsPreparedAccLogin {
	msg := MsgGate2LsPreparedAccLogin{
		MsgHeadNoProtocol: MsgHeadNoProtocol{Ctrl: 0x02, Cmd: 0xE3},
	}
	copy(msg.Account[:], utils.MakeFixedLengthStringBytes(account, 0x15))
	msg.SetSize()
	return &msg
}

func ReadMsgGate2LsPreparedAccLogin(packet []byte) (*MsgGate2LsPreparedAccLogin, error) {
	var msg MsgGate2LsPreparedAccLogin
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
