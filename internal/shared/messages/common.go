package messages

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
