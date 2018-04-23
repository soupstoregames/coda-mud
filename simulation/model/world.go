package model

type World struct {
	Name  string
	Rooms map[RoomID]*Room
}

func NewWorld(name string) *World {
	return &World{
		Name:  name,
		Rooms: make(map[RoomID]*Room),
	}
}
