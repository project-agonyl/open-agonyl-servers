package shared

import "sync/atomic"

type UidGenerator struct {
	start uint32
	uid   atomic.Uint32
}

func NewUidGenerator(startValue uint32) *UidGenerator {
	gen := &UidGenerator{
		start: startValue,
	}
	gen.uid.Store(startValue)
	return gen
}

func (l *UidGenerator) Uid() uint32 {
	return l.uid.Add(1)
}
