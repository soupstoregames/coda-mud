package simulation

import (
	"github.com/soupstore/coda-world/simulation/model"
)

// WorldController is an interface over Simulation for modifying the world itself
type WorldController interface {
	AddWorld(worldID model.WorldID) error
	RemoveWorld(worldID model.WorldID)
	MakeRoom(worldID model.WorldID, roomID model.RoomID, name, description string) error
	GetRoom(worldID model.WorldID, roomID model.RoomID) (*model.Room, error)
	RemoveRoom(worldID model.WorldID, roomID model.RoomID) error
	SetSpawnRoom(worldID model.WorldID, roomID model.RoomID) error
	LinkRoom(originWorldName model.WorldID, origin model.RoomID, direction model.Direction, destinationWorldID model.WorldID, destination model.RoomID) error
	SpawnItem(*model.Item, model.ContainerID) error
}

func (s *Simulation) AddWorld(worldID model.WorldID) error {
	s.worlds[worldID] = model.NewWorld(worldID)
	return nil
}

func (s *Simulation) RemoveWorld(worldID model.WorldID) {
	delete(s.worlds, worldID)
}

// MakeRoom creates a new room at the next available ID
func (s *Simulation) MakeRoom(worldID model.WorldID, roomID model.RoomID, name, description string) error {
	containerID := s.getNextContainerID()

	world, ok := s.worlds[worldID]
	if !ok {
		return ErrWorldNotFound
	}

	room := model.NewRoom(roomID, worldID, containerID, name, description)
	world.Rooms[roomID] = room

	container := room.Container
	s.containers[container.ID] = container

	return nil
}

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

func (s *Simulation) RemoveRoom(worldID model.WorldID, roomID model.RoomID) error {
	world, ok := s.worlds[worldID]
	if !ok {
		return ErrWorldNotFound
	}

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

// LinkRoom creates an exit from the origin room to the destination room in the direction specified.
func (s *Simulation) LinkRoom(originWorldName model.WorldID, origin model.RoomID, direction model.Direction, destinationWorldID model.WorldID, destination model.RoomID) error {
	originWorld, ok := s.worlds[originWorldName]
	if !ok {
		return ErrWorldNotFound
	}

	originRoom, ok := originWorld.Rooms[origin]
	if !ok {
		return ErrRoomNotFound
	}

	originRoom.Exits[direction] = &model.Exit{
		WorldID: destinationWorldID,
		RoomID:  destination,
	}

	return nil
}

func (s *Simulation) SpawnItem(item *model.Item, containerID model.ContainerID) error {
	container, ok := s.containers[containerID]
	if !ok {
		return ErrContainerNotFound
	}

	container.PutItem(item)
	return nil
}
