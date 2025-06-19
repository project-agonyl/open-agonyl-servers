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

type GateServerInfo struct {
	ServerID     byte
	ServerName   [0x11]byte
	ServerStatus [0x51]byte
}

type MsgLs2GateLogin struct {
	MsgHeadNoProtocol
	Account [0x15]byte
	Unknown [0x09]byte
}

func (msg *MsgLs2GateLogin) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgLs2GateLogin) SetSize() {
	msg.Size = msg.GetSize()
}

func (msg *MsgLs2GateLogin) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgLs2GateLogin(account string, pcId uint32) *MsgLs2GateLogin {
	msg := MsgLs2GateLogin{
		MsgHeadNoProtocol: MsgHeadNoProtocol{Ctrl: 0x01, Cmd: 0xE1, PcId: pcId},
	}
	copy(msg.Account[:], utils.MakeFixedLengthStringBytes(account, 0x15))
	msg.SetSize()
	return &msg
}

func ReadMsgLs2GateLogin(packet []byte) (*MsgLs2GateLogin, error) {
	var msg MsgLs2GateLogin
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

type MsgS2CGateInfo struct {
	MsgHeadNoProtocol
	PcId   uint32
	ZaIP   [0x10]byte
	ZaPort uint32
}

func (msg *MsgS2CGateInfo) GetSize() uint32 {
	return uint32(binary.Size(msg))
}

func (msg *MsgS2CGateInfo) SetSize() {
	msg.Size = msg.GetSize()
}

func (msg *MsgS2CGateInfo) GetBytes() []byte {
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.LittleEndian, msg)
	return buffer.Bytes()
}

func NewMsgS2CGateInfo(pcId uint32, zaIP string, zaPort uint32) *MsgS2CGateInfo {
	msg := MsgS2CGateInfo{
		MsgHeadNoProtocol: MsgHeadNoProtocol{Ctrl: 0x01, Cmd: 0xE2, PcId: pcId},
		PcId:              pcId,
		ZaPort:            zaPort,
	}
	copy(msg.ZaIP[:], utils.MakeFixedLengthStringBytes(zaIP, 0x10))
	msg.SetSize()
	return &msg
}

func ReadMsgS2CGateInfo(packet []byte) (*MsgS2CGateInfo, error) {
	var msg MsgS2CGateInfo
	if err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
