package data

import (
	"encoding/binary"
	"os"

	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
)

type WarpDataRaw struct {
	MapId    uint16
	X        byte
	Y        byte
	Unknown1 uint16
}

type MapDataRaw struct {
	Id        uint16
	Name      [0x14]byte
	WarpCount byte
}

type NavigationDataRaw struct {
	IsMovable byte
	Unknown1  [0x3]byte
}

type WarpData struct {
	MapId uint16
	X     byte
	Y     byte
}

type NavigationData struct {
	IsMovable bool
}

type MapData struct {
	Id             uint16
	Name           string
	WarpCount      byte
	WarpData       []WarpData
	NavigationMesh [0xFF][0xFF]NavigationData
}

func LoadMapData(mapFilePath string) (*MapData, error) {
	mapFile, err := os.Open(mapFilePath)
	if err != nil {
		return nil, err
	}
	defer mapFile.Close()

	var mapData MapData
	var mapDataRaw MapDataRaw
	err = binary.Read(mapFile, binary.LittleEndian, &mapDataRaw)
	if err != nil {
		return nil, err
	}

	mapData = MapData{
		Id:        mapDataRaw.Id,
		Name:      utils.ReadStringFromBytes(mapDataRaw.Name[:]),
		WarpCount: mapDataRaw.WarpCount,
		WarpData:  make([]WarpData, mapDataRaw.WarpCount),
	}

	for i := 0; i < int(mapDataRaw.WarpCount); i++ {
		var warpDataRow WarpDataRaw
		err = binary.Read(mapFile, binary.LittleEndian, &warpDataRow)
		if err != nil {
			return nil, err
		}

		mapData.WarpData[i] = WarpData{
			MapId: warpDataRow.MapId,
			X:     warpDataRow.X,
			Y:     warpDataRow.Y,
		}
	}

	for x := 0; x < 0xFF; x++ {
		for y := 0; y < 0xFF; y++ {
			var navDataRaw NavigationDataRaw
			err = binary.Read(mapFile, binary.LittleEndian, &navDataRaw)
			if err != nil {
				return nil, err
			}

			mapData.NavigationMesh[x][y] = NavigationData{
				IsMovable: navDataRaw.IsMovable != 0,
			}
		}
	}

	return &mapData, nil
}
