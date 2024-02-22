package simulation

import (
	"github.com/soupstoregames/coda-mud/simulation/model"
)

// WorldController is an interface over Simulation for modifying the world itself
type WorldController interface {
	CreateWorld(worldID model.WorldID, instance bool, alone bool) error
	DestroyWorld(worldID model.WorldID)
	CreateRoom(worldID model.WorldID, roomID model.RoomID, name, region, description, script string) (*model.Room, error)
	GetRoom(worldID model.WorldID, roomID model.RoomID) (*model.Room, error)
	DestroyRoom(worldID model.WorldID, roomID model.RoomID) error
	SetSpawnRoom(worldID model.WorldID, roomID model.RoomID) error
	CreateItemDefinition(itemID model.ItemDefinitionID, name string, aliases []string, weight int64, rigSlot model.RigSlot, container *model.ContainerDefinition) (*model.ItemDefinition, error)
	SpawnItem(itemDefinitionID model.ItemDefinitionID, containerID model.ContainerID) error
}

// CreateWorld creates a new world in the simulation.
// Every world must have a unique WorldID, which is a type aliased sting.
func (s *Simulation) CreateWorld(worldID model.WorldID, instancable bool, alone bool) error {
	// TODO: check for uniqueness
	s.worlds[worldID] = model.NewWorld(worldID, instancable, false, alone)
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

	room := model.NewRoom(roomID, worldID, name, region, description, script, world.Alone)
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
func (s *Simulation) CreateItemDefinition(itemID model.ItemDefinitionID, name string, aliases []string, weight int64, rigSlot model.RigSlot, container *model.ContainerDefinition) (*model.ItemDefinition, error) {
	item := model.NewItemDefinition(itemID, name, aliases, weight, rigSlot, container)
	s.itemDefinitions[itemID] = item
	return item, nil
}

// SpawnItem creates a new instance of the item definition in the desired container.
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
