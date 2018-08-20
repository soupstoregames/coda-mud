package database

import (
	"github.com/go-pg/pg"
	"github.com/soupstore/coda/simulation/model"
)

type Container struct {
	ID    model.ContainerID
	Type  string
	Items []model.ItemID
}

type RoomContainerLink struct {
	ID    model.ContainerID
	World model.WorldID
	Room  model.RoomID
}

func GetContainers(db *pg.DB) ([]*Container, []*RoomContainerLink, error) {
	var containers []*Container
	if err := db.Model(&containers).Select(); err != nil {
		return nil, nil, err
	}

	var roomContainerLinks []*RoomContainerLink
	if err := db.Model(&roomContainerLinks).Select(); err != nil {
		return nil, nil, err
	}

	return containers, roomContainerLinks, nil
}

func RemoveItemFromContainer(db *pg.DB, containerID model.ContainerID, itemID model.ItemID) error {
	container := &Container{ID: containerID}
	if err := db.Model(container).Select(); err != nil {
		return err
	}

	container.Items = removeFromItemList(container.Items, itemID)

	return db.Update(container)
}

func removeFromItemList(l []model.ItemID, i model.ItemID) []model.ItemID {
	result := []model.ItemID{}
	for _, id := range l {
		if id == i {
			continue
		}
		result = append(result, id)
	}
	return result
}
