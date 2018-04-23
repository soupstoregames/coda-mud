package data

import (
	"path"
)

type Data struct {
	Worlds map[string]map[int]*Room
}

func Load(rootPath string) (*Data, error) {
	worlds, err := loadRooms(path.Join(rootPath, "rooms"))
	if err != nil {
		return nil, err
	}

	return &Data{
		Worlds: worlds,
	}, nil
}
