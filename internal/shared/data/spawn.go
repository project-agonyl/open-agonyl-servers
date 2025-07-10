package data

import (
	"encoding/binary"
	"os"
)

type MonsterSpawnData struct {
	Id          uint16
	X           byte
	Y           byte
	Unknown1    uint16
	Orientation byte
	SpwanStep   byte
}

func LoadMonsterSpawnData(spawnFilePath string) ([]MonsterSpawnData, error) {
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
	spawnData := make([]MonsterSpawnData, totalSpawns)
	for i := int64(0); i < totalSpawns; i++ {
		spawnData[i] = MonsterSpawnData{}
		err = binary.Read(spawnFile, binary.LittleEndian, &spawnData[i])
		if err != nil {
			return nil, err
		}
	}

	return spawnData, nil
}

func (m *MonsterSpawnData) IsMonster() bool {
	return m.Id < 1000
}
