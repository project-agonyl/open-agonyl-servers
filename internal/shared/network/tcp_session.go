package network

type TCPServerSession interface {
	ID() uint32
	Handle()
	Close() error
	Send(data []byte) error
}
