package simulation

import (
	"github.com/soupstore/coda-world/data"
	"github.com/soupstore/coda-world/simulation/model"
)

// WorldController is an interface over Simulation for modifying the world itself
type WorldController interface {
	MakeRoom(name, description string) model.RoomID
	SetSpawnRoom(id model.RoomID) error
	LinkRoom(model.RoomID, model.Direction, model.RoomID, bool) error
	SpawnItem(*model.Item, model.ContainerID) error
}

func (s *Simulation) AddWorld(worldID model.WorldID) error {
	s.worlds[worldID] = model.NewWorld(worldID)
	return nil
}

// MakeRoom creates a new room at the next available ID
func (s *Simulation) loadRoom(worldID model.WorldID, roomID model.RoomID, r *data.Room) error {
	containerID := s.getNextContainerID()

	world, ok := s.worlds[worldID]
	if !ok {
		return ErrWorldNotFound
	}

	room := model.NewRoom(roomID, worldID, containerID, r.Name, r.Description)
	world.Rooms[roomID] = room

	container := room.Container
	s.containers[container.ID] = container

	return nil
}

// SetSpawnRoom sets the room that all new characters will start in
func (s *Simulation) SetSpawnRoom(worldID model.WorldID, roomID model.RoomID) error {
	room, err := s.getRoom(worldID, roomID)
	if err != nil {
		return ErrRoomNotFound
	}

	s.spawnRoom = room

	return nil
}

// LinkRoom creates an exit from the origin room to the destination room in the direction specified.
func (s *Simulation) LinkRoom(originWorldName model.WorldID, origin model.RoomID, direction model.Direction, destinationWorldName model.WorldID, destination model.RoomID) error {
	originWorld, ok := s.worlds[originWorldName]
	if !ok {
		return ErrWorldNotFound
	}

	originRoom, ok := originWorld.Rooms[origin]
	if !ok {
		return ErrRoomNotFound
	}

	destinationWorld, ok := s.worlds[destinationWorldName]
	if !ok {
		return ErrWorldNotFound
	}

	destinationRoom, ok := destinationWorld.Rooms[destination]
	if !ok {
		return ErrRoomNotFound
	}

	originRoom.Exits[direction] = &model.Exit{
		WorldID: originRoom.WorldID,
		RoomID:  destinationRoom.ID,
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
