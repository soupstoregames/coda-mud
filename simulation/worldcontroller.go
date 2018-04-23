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

func (s *Simulation) GetRoom(worldName string, roomID model.RoomID) (*model.Room, error) {
	world, ok := s.worlds[worldName]
	if !ok {
		return nil, ErrWorldNotFound
	}

	spawnRoom, ok := world.Rooms[roomID]
	if !ok {
		return nil, ErrRoomNotFound
	}
	return spawnRoom, nil
}

func (s *Simulation) AddWorld(worldName string) error {
	s.worlds[worldName] = model.NewWorld(worldName)
	return nil
}

// MakeRoom creates a new room at the next available ID
func (s *Simulation) loadRoom(worldName string, roomID model.RoomID, r *data.Room) error {
	containerID := s.getNextContainerID()

	world, ok := s.worlds[worldName]
	if !ok {
		return ErrWorldNotFound
	}

	room := model.NewRoom(roomID, containerID, r.Name, r.Description)
	world.Rooms[roomID] = room

	container := room.Container
	s.containers[container.ID] = container

	return nil
}

// SetSpawnRoom sets the room that all new characters will start in
func (s *Simulation) SetSpawnRoom(worldName string, id model.RoomID) error {
	world, ok := s.worlds[worldName]
	if !ok {
		return ErrWorldNotFound
	}

	spawnRoom, ok := world.Rooms[id]
	if !ok {
		return ErrRoomNotFound
	}

	s.spawnRoom = spawnRoom

	return nil
}

// LinkRoom creates an exit from the origin room to the destination room in the direction specified.
// If the link is bidirectional then an exit is created from the destination room to the origin room in the opposite direction.
func (s *Simulation) LinkRoom(originWorldName string, origin model.RoomID, direction model.Direction, destinationWorldName string, destination model.RoomID, bidirectional bool) error {
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

	originRoom.Exits[direction] = destinationRoom

	if bidirectional {
		destinationRoom.Exits[direction.Opposite()] = originRoom
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
