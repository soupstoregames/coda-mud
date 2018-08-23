package static

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"time"

	"github.com/nicklanng/fsdiff"
	"github.com/soupstore/coda/simulation"
	"github.com/soupstore/coda/simulation/model"
	"github.com/soupstore/go-core/logging"
)

const roomExtension = ".toml"

type DataWatcher struct {
	Errors        chan error
	dataFolder    string
	lastDataState *fsdiff.Node
	sim           simulation.WorldController
}

func NewDataWatcher(rootPath string, sim simulation.WorldController) *DataWatcher {
	dw := &DataWatcher{
		Errors:     make(chan error, 1),
		dataFolder: rootPath,
		sim:        sim,
	}

	return dw
}

func (dw *DataWatcher) InitialLoad() error {
	dataState, err := dw.initialLoad()
	if err != nil {
		return err
	}
	dw.lastDataState = dataState
	return nil
}

func (dw *DataWatcher) Watch() {
	// regularly load data folder and check for differences
	t := time.NewTicker(time.Minute)
	go func() {
		for range t.C {
			// get current state of data folder
			newState, err := fsdiff.BuildTree(dw.dataFolder)
			if err != nil {
				dw.Errors <- err
				continue
			}

			// find the diffed files
			diff := fsdiff.Compare(dw.lastDataState, newState)

			// apply diffs to simulation
			dw.applyDiff(diff)

			// save state for next time
			dw.lastDataState = newState
		}
	}()
}

func (dw *DataWatcher) initialLoad() (*fsdiff.Node, error) {
	items, err := loadAllItems(path.Join(dw.dataFolder, "items"))
	if err != nil {
		return nil, err
	}

	worlds, err := loadAllWorlds(path.Join(dw.dataFolder, "rooms"))
	if err != nil {
		return nil, err
	}

	state, err := fsdiff.BuildTree(dw.dataFolder)
	if err != nil {
		return nil, err
	}

	// load items
	for fileID, item := range items {
		itemDefinitionID := model.ItemDefinitionID(fileID)
		dw.addItemToSim(itemDefinitionID, item)
	}

	// load worlds
	for worldID, rooms := range worlds {
		wID := model.WorldID(worldID)
		dw.addWorldToSim(wID, rooms)
	}

	return state, nil
}

func (dw *DataWatcher) applyDiff(diff *fsdiff.Diff) {
	if diff.DiffType == fsdiff.DiffTypeNone {
		return
	}

	// get room folder
	rooms, ok := searchChildrenForName(diff, "rooms")
	if !ok {
		dw.Errors <- errors.New("no rooms folder")
		return
	}

	dw.applyWorldDiffs(rooms)

	return
}

func (dw *DataWatcher) applyWorldDiffs(diff *fsdiff.Diff) {
	// has room folder changed?
	if diff.DiffType != fsdiff.DiffTypeChanged {
		return
	}

	for _, world := range diff.Children {
		// if this world hasnt changed - move on
		if world.DiffType == fsdiff.DiffTypeNone {
			continue
		}

		worldID := model.WorldID(filepath.Base(world.Path))

		// the world has been added - load all rooms into the sim
		if world.DiffType == fsdiff.DiffTypeAdded {
			rooms, err := loadWorldFolder(world.Path)
			if err != nil {
				dw.Errors <- errors.New("failed to load world")
			}

			for roomID, room := range rooms {
				rID := model.RoomID(roomID)
				dw.addRoomToSim(worldID, rID, room)
			}

			logging.Info(fmt.Sprintf("Loaded world '%s' with %d rooms", worldID, len(rooms)))
		}

		// the world has been added - load all rooms into the sim
		if world.DiffType == fsdiff.DiffTypeRemoved {
			dw.sim.DestroyWorld(worldID)

			logging.Info(fmt.Sprintf("Removed world '%s'", worldID))
		}

		// the world have been changed, move down to the room level
		if world.DiffType == fsdiff.DiffTypeChanged {
			for _, room := range world.Children {
				if filepath.Ext(room.Path) != roomExtension {
					continue
				}

				switch room.DiffType {
				case fsdiff.DiffTypeAdded:
					roomID, err := getRoomID(filepath.Base(room.Path))
					if err != nil {
						dw.Errors <- err
						continue
					}

					room, err := loadRoom(room.Path)
					if err != nil {
						dw.Errors <- err
						continue
					}

					dw.addRoomToSim(worldID, model.RoomID(roomID), room)
					logging.Info(fmt.Sprintf("Added room %d to world '%s'", roomID, worldID))

				case fsdiff.DiffTypeRemoved:
					roomID, err := getRoomID(filepath.Base(room.Path))
					if err != nil {
						dw.Errors <- err
						continue
					}

					dw.sim.DestroyRoom(worldID, model.RoomID(roomID))
					logging.Info(fmt.Sprintf("Removed room %d in world '%s'", roomID, worldID))

				case fsdiff.DiffTypeChanged:
					roomID, err := getRoomID(filepath.Base(room.Path))
					if err != nil {
						dw.Errors <- err
						continue
					}

					room, err := loadRoom(room.Path)
					if err != nil {
						dw.Errors <- err
						continue
					}

					dw.updateRoomInSim(worldID, model.RoomID(roomID), room)
					logging.Info(fmt.Sprintf("Updated room %d in world '%s'", roomID, worldID))
				}
			}
		}
	}
}

func (dw *DataWatcher) addWorldToSim(worldID model.WorldID, rooms map[int]*Room) {
	dw.sim.CreateWorld(worldID)

	// load rooms
	for roomID, room := range rooms {
		rID := model.RoomID(roomID)

		dw.addRoomToSim(worldID, rID, room)
	}

	logging.Info(fmt.Sprintf("Loaded world '%s' with %d rooms", worldID, len(rooms)))
}

func (dw *DataWatcher) addRoomToSim(worldID model.WorldID, roomID model.RoomID, room *Room) error {
	r, err := dw.sim.CreateRoom(worldID, roomID, room.Name, room.Region, room.Description, room.Script)
	if err != nil {
		return err
	}

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

		r.Exits[d] = &model.Exit{
			WorldID: model.WorldID(exit.WorldID),
			RoomID:  model.RoomID(exit.RoomID),
		}
	}

	return nil
}

func (dw *DataWatcher) updateRoomInSim(worldID model.WorldID, roomID model.RoomID, room *Room) error {
	r, err := dw.sim.GetRoom(worldID, roomID)
	if err != nil {
		return err
	}

	r.Name = room.Name
	r.Description = room.Description
	r.Exits = make(map[model.Direction]*model.Exit)

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

		r.Exits[d] = &model.Exit{
			WorldID: model.WorldID(exit.WorldID),
			RoomID:  model.RoomID(exit.RoomID),
		}
	}

	return nil
}

func (dw *DataWatcher) addItemToSim(itemDefinitionID model.ItemDefinitionID, item *Item) error {
	var rigSlot model.RigSlot
	switch item.RigSlot {
	case "backpack":
		rigSlot = model.RigSlotBackpack
	}

	var container *model.ContainerDefinition
	if item.Container != nil {
		container = &model.ContainerDefinition{}
	}

	_, err := dw.sim.CreateItemDefinition(itemDefinitionID, item.Name, item.Aliases, rigSlot, container)
	return err
}

func searchChildrenForName(parent *fsdiff.Diff, name string) (*fsdiff.Diff, bool) {
	for _, ch := range parent.Children {
		if filepath.Base(ch.Path) == name {
			return ch, true
		}
	}
	return nil, false
}
