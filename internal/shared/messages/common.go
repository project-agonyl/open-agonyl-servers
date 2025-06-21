package messages

import (
	"bytes"
	"encoding/binary"
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
		MsgHead: MsgHead{Protocol: 0x0FFF, MsgHeadNoProtocol: MsgHeadNoProtocol{Ctrl: 0x03, Cmd: 0xFF}},
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
