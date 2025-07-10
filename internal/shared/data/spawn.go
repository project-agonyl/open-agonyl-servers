package data

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

type NPCSpawnData struct {
	Id          uint16
	X           byte
	Y           byte
	Unknown1    uint16
	Orientation byte
	SpwanStep   byte
}

func LoadNPCSpawnData(spawnFilePath string) ([]NPCSpawnData, error) {
	if !strings.HasSuffix(spawnFilePath, ".n_ndt") {
		return nil, fmt.Errorf("invalid spawn file path: %s", spawnFilePath)
	}

	spawnFile, err := os.Open(spawnFilePath)
	if err != nil {
		return nil, err
	}

	defer spawnFile.Close()

	spawnFileStat, err := spawnFile.Stat()
	if err != nil {
		return nil, err
	}

	totalSpawns := spawnFileStat.Size() / 8
	spawnData := make([]NPCSpawnData, totalSpawns)
	for i := int64(0); i < totalSpawns; i++ {
		spawnData[i] = NPCSpawnData{}
		err = binary.Read(spawnFile, binary.LittleEndian, &spawnData[i])
		if err != nil {
			return nil, err
		}
	}

	return spawnData, nil
}

func (m *NPCSpawnData) IsMonster() bool {
	return m.Id < 1000
}
