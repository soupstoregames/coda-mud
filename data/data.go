package data

import (
	"fmt"
	"path"

	"github.com/nicklanng/fsdiff"
	"github.com/soupstore/coda-world/log"
	"github.com/soupstore/coda-world/simulation"
	"github.com/soupstore/coda-world/simulation/model"
)

type Data struct {
	Worlds map[string]map[int]*Room
}

func WatchDataFolder(rootPath string, sim simulation.WorldController) error {
	initialLoad(rootPath, sim)
	_, err := fsdiff.BuildTree(rootPath)
	if err != nil {
		return err
	}

	return nil
}

func initialLoad(rootPath string, sim simulation.WorldController) error {
	worlds, err := loadRooms(path.Join(rootPath, "rooms"))
	if err != nil {
		return err
	}

	// load worlds
	for worldID, rooms := range worlds {
		wID := model.WorldID(worldID)
		sim.AddWorld(wID)

		// load rooms
		for roomID, room := range rooms {
			rID := model.RoomID(roomID)
			sim.MakeRoom(wID, rID, room.Name, room.Description)

			// load room exits
			for direction, exit := range room.Exits {
				d, err := model.StringToDirection(direction)
				if err != nil {
					return err
				}

				// if no worldID is provided, it defaults to the same as the room lives in
				if exit.WorldID == "" {
					exit.WorldID = worldID
				}

				sim.LinkRoom(wID, rID, d, model.WorldID(exit.WorldID), model.RoomID(exit.RoomID))
			}
		}

		log.Logger().Info(fmt.Sprintf("Loaded world '%s' with %d rooms", worldID, len(rooms)))
	}

	return nil
}
