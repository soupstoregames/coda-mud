package data

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"time"

	"github.com/nicklanng/fsdiff"
	"github.com/soupstore/coda-world/log"
	"github.com/soupstore/coda-world/simulation"
	"github.com/soupstore/coda-world/simulation/model"
)

type DataWatcher struct {
	Errors        chan error
	dataFolder    string
	lastDataState *fsdiff.Node
	sim           simulation.WorldController
}

func NewDataWatcher(rootPath string, sim simulation.WorldController) *DataWatcher {
	return &DataWatcher{
		Errors:     make(chan error, 1),
		dataFolder: rootPath,
		sim:        sim,
	}
}

func (dw *DataWatcher) Watch() {
	// load initial state
	dataState, err := dw.initialLoad()
	if err != nil {
		dw.Errors <- err
	}
	dw.lastDataState = dataState

	// regularly load data folder and check for differences
	t := time.NewTicker(time.Minute)
	go func() {
		for _ = range t.C {
			// get current state of data folder
			newState, err := fsdiff.BuildTree(dw.dataFolder)
			if err != nil {
				dw.Errors <- err
				continue
			}

			// find the diffed files
			diff := fsdiff.Compare(dw.lastDataState, newState)

			// find changes to data folder
			dw.walkDiff(diff)

			// save state for next time
			dw.lastDataState = newState
		}
	}()
}

func (dw *DataWatcher) initialLoad() (*fsdiff.Node, error) {
	worlds, err := loadAllWorlds(path.Join(dw.dataFolder, "rooms"))
	if err != nil {
		return nil, err
	}

	state, err := fsdiff.BuildTree(dw.dataFolder)
	if err != nil {
		return nil, err
	}

	// load worlds
	for worldID, rooms := range worlds {
		wID := model.WorldID(worldID)
		dw.addWorldToSim(wID, rooms)
	}

	return state, nil
}

func (dw *DataWatcher) walkDiff(diff *fsdiff.Diff) {
	if diff.DiffType == fsdiff.DiffTypeNone {
		return
	}

	// get room folder
	rooms, ok := searchChildrenForName(diff, "rooms")
	if !ok {
		dw.Errors <- errors.New("no rooms folder")
		return
	}

	dw.walkWorlds(rooms)

	return
}

func (dw *DataWatcher) walkWorlds(diff *fsdiff.Diff) {
	// has room folder changed?
	if diff.DiffType != fsdiff.DiffTypeChanged {
		return
	}

	for _, world := range diff.Children {
		// if this world hasnt changed - move on
		if world.DiffType == fsdiff.DiffTypeNone {
			continue
		}

		// the world has been added - load all rooms into the sim
		if world.DiffType == fsdiff.DiffTypeAdded {
			worldID := model.WorldID(filepath.Base(world.Path))
			rooms, err := loadWorldFolder(world.Path)
			if err != nil {
				dw.Errors <- errors.New("failed to load world")
			}

			for roomID, room := range rooms {
				rID := model.RoomID(roomID)
				dw.addRoomToSim(worldID, rID, room)
			}

			log.Logger().Info(fmt.Sprintf("Loaded world '%s' with %d rooms", worldID, len(rooms)))
		}

		// the world has been added - load all rooms into the sim
		if world.DiffType == fsdiff.DiffTypeRemoved {
			worldID := model.WorldID(filepath.Base(world.Path))
			dw.sim.RemoveWorld(worldID)

			log.Logger().Info(fmt.Sprintf("Removed world '%s'", worldID))
		}

		// the world have been changed, move down to the room level
		if world.DiffType == fsdiff.DiffTypeChanged {
			// TODO: test rooms
		}
	}
}

func (dw *DataWatcher) addWorldToSim(worldID model.WorldID, rooms map[int]*Room) {
	dw.sim.AddWorld(worldID)

	// load rooms
	for roomID, room := range rooms {
		rID := model.RoomID(roomID)

		dw.addRoomToSim(worldID, rID, room)
	}

	log.Logger().Info(fmt.Sprintf("Loaded world '%s' with %d rooms", worldID, len(rooms)))
}

func (dw *DataWatcher) addRoomToSim(worldID model.WorldID, roomID model.RoomID, room *Room) error {
	dw.sim.MakeRoom(worldID, roomID, room.Name, room.Description)

	// load room exits
	for direction, exit := range room.Exits {
		d, err := model.StringToDirection(direction)
		if err != nil {
			return err
		}

		// if no worldID is provided, it defaults to the same as the room lives in
		if exit.WorldID == "" {
			exit.WorldID = string(worldID)
		}

		dw.sim.LinkRoom(worldID, roomID, d, model.WorldID(exit.WorldID), model.RoomID(exit.RoomID))
	}

	return nil
}

func searchChildrenForName(parent *fsdiff.Diff, name string) (*fsdiff.Diff, bool) {
	for _, ch := range parent.Children {
		if filepath.Base(ch.Path) == name {
			return ch, true
		}
	}
	return nil, false
}
