package simulation

import (
	"github.com/soupstore/coda/simulation/model"
)

// WorldController is an interface over Simulation for modifying the world itself
type WorldController interface {
	CreateWorld(worldID model.WorldID) error
	DestroyWorld(worldID model.WorldID)
	CreateRoom(worldID model.WorldID, roomID model.RoomID, name, region, description, script string) (*model.Room, error)
	GetRoom(worldID model.WorldID, roomID model.RoomID) (*model.Room, error)
	DestroyRoom(worldID model.WorldID, roomID model.RoomID) error
	SetSpawnRoom(worldID model.WorldID, roomID model.RoomID) error
	CreateItemDefinition(itemID model.ItemDefinitionID, name string, aliases []string, rigSlot model.RigSlot, container *model.ContainerDefinition) (*model.ItemDefinition, error)
	SpawnItem(itemDefinitionID model.ItemDefinitionID, containerID model.ContainerID) error
}

// CreateWorld creates a new world in the simulation.
// Every world must have a unique WorldID, which is a type aliased sting.
func (s *Simulation) CreateWorld(worldID model.WorldID) error {
	// TODO: check for uniqueness
	s.worlds[worldID] = model.NewWorld(worldID)
	return nil
}

// DestroyWorld unloads a world and all of its rooms from the simulation.
func (s *Simulation) DestroyWorld(worldID model.WorldID) {
	// TODO: Move all characters in this world to a safe location
	delete(s.worlds, worldID)
}

// CreateRoom creates a new room in the specified world with the specified room ID
func (s *Simulation) CreateRoom(worldID model.WorldID, roomID model.RoomID, name, region, description, script string) (*model.Room, error) {
	world, ok := s.worlds[worldID]
	if !ok {
		return nil, ErrWorldNotFound
	}

	// TODO: Check that room with ID does not already exist

	room := model.NewRoom(roomID, worldID, name, region, description, script)
	world.Rooms[roomID] = room

	container := room.Container
	s.containers[container.ID()] = container

	return room, nil
}

// GetRoom returns the room object as the specified world ID and room ID
func (s *Simulation) GetRoom(worldID model.WorldID, roomID model.RoomID) (*model.Room, error) {
	world, ok := s.worlds[worldID]
	if !ok {
		return nil, ErrWorldNotFound
	}

	room, ok := world.Rooms[roomID]
	if !ok {
		return nil, ErrRoomNotFound
	}

	return room, nil
}

// DestroyRoom removes a room from a world.
func (s *Simulation) DestroyRoom(worldID model.WorldID, roomID model.RoomID) error {
	world, ok := s.worlds[worldID]
	if !ok {
		return ErrWorldNotFound
	}

	// TODO: Clean up broken exits in the rest of the sim.

	delete(world.Rooms, roomID)

	return nil
}

// SetSpawnRoom sets the room that all new characters will start in
func (s *Simulation) SetSpawnRoom(worldID model.WorldID, roomID model.RoomID) error {
	room, err := s.GetRoom(worldID, roomID)
	if err != nil {
		return ErrRoomNotFound
	}

	s.spawnRoom = room

	return nil
}

// CreateItemDefinition creates a new item definition
func (s *Simulation) CreateItemDefinition(itemID model.ItemDefinitionID, name string, aliases []string, rigSlot model.RigSlot, container *model.ContainerDefinition) (*model.ItemDefinition, error) {
	item := model.NewItemDefinition(itemID, name, aliases, rigSlot, container)
	s.itemDefinitions[itemID] = item
	return item, nil
}

func (s *Simulation) SpawnItem(itemDefinitionID model.ItemDefinitionID, containerID model.ContainerID) error {
	container, ok := s.containers[containerID]
	if !ok {
		return ErrContainerNotFound
	}

	definition := s.itemDefinitions[itemDefinitionID]

	instance := definition.Spawn()
	s.items[instance.ID] = instance

	container.PutItem(instance)

	return nil
}

// func (s *Simulation) LoadCharacters(characters []*state.Character) {
// 	loaded := 0
// 	for _, ch := range characters {
// 		room, err := s.GetRoom(ch.World, ch.Room)
// 		if err != nil {
// 			logging.Logger().Error("failed to load character")
// 		}
// 		character := model.NewCharacter(ch.ID, ch.Name, room)
// 		s.characters[ch.ID] = character

// 		// add character to room
// 		room.AddCharacter(character)

// 		loaded++
// 	}
// 	logging.Logger().Info(fmt.Sprintf("Loaded %d characters", loaded))
// }

// func (s *Simulation) LoadContainers(containers []*database.Container, roomContainerLinks []*database.RoomContainerLink) {
// 	loaded := 0
// 	for _, c := range containers {
// 		switch c.Type {
// 		case "room":
// 			s.containers[c.ID] = model.NewRoomContainer(c.ID)
// 			for _, i := range c.Items {
// 				s.containers[c.ID].PutItem(s.items[i])
// 			}
// 		}

// 		loaded++
// 	}

// 	// link containers to rooms
// 	for _, l := range roomContainerLinks {
// 		room, err := s.GetRoom(l.World, l.Room)
// 		if err != nil {
// 			fmt.Println("did not find room")
// 		}
// 		container, ok := s.containers[l.ID]
// 		if !ok {
// 			fmt.Println("did not find container")
// 		}

// 		room.Container = container
// 	}

// 	logging.Logger().Info(fmt.Sprintf("Loaded %d containers", loaded))
// }

// func (s *Simulation) LoadItems(items []*database.Item) {
// 	loaded := 0
// 	for _, i := range items {
// 		definition := s.itemDefinitions[i.ItemDefinitionID]
// 		s.items[i.ID] = definition.Spawn(i.ID)

// 		loaded++
// 	}
// 	logging.Logger().Info(fmt.Sprintf("Loaded %d items", loaded))
// }
