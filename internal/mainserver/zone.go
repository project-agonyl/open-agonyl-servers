package mainserver

import "slices"

type Zone struct {
	serverId byte
	session  *mainServerSession
	maps     []uint16
}

func NewZone(serverId byte, session *mainServerSession) *Zone {
	return &Zone{
		serverId: serverId,
		session:  session,
		maps:     make([]uint16, 0),
	}
}

func (z *Zone) HasMap(mapId uint16) bool {
	return slices.Contains(z.maps, mapId)
}

func (z *Zone) AddMap(mapId uint16) {
	if z.HasMap(mapId) {
		return
	}

	z.maps = append(z.maps, mapId)
}

func (z *Zone) RemoveMap(mapId uint16) {
	z.maps = slices.DeleteFunc(z.maps, func(m uint16) bool {
		return m == mapId
	})
}

func (z *Zone) GetMaps() []uint16 {
	return z.maps
}

func (z *Zone) Send(packet []byte) error {
	return z.session.Send(packet)
}

func (z *Zone) GetServerId() byte {
	return z.serverId
}

func (z *Zone) SetMaps(maps []uint16) {
	z.maps = maps
}
