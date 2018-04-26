package model

type WorldID string

type World struct {
	WorldID WorldID
	Name    string
	Rooms   map[RoomID]*Room
}

func NewWorld(id WorldID) *World {
	return &World{
		WorldID: id,
		Name:    "",
		Rooms:   make(map[RoomID]*Room),
	}
}
