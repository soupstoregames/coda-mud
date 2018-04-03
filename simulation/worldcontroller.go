package simulation

import "github.com/soupstore/coda-world/simulation/model"

// WorldController is an interface over Simulation for modifying the world itself
type WorldController interface {
	MakeRoom(name, description string) model.RoomID
	SetSpawnRoom(id model.RoomID)
}

// MakeRoom creates a new room at the next available ID
func (s *Simulation) MakeRoom(name, description string) model.RoomID {
	roomID := s.getNextRoomID()
	room := model.NewRoom(roomID, name, description)
	s.rooms[roomID] = room
	return roomID
}

// SetSpawnRoom sets the room that all new characters will start in
func (s *Simulation) SetSpawnRoom(id model.RoomID) {
	spawnRoom, ok := s.rooms[id]
	if !ok {
		// error
	}
	s.spawnRoom = spawnRoom
}
