package model

type WorldID string

type World struct {
	WorldID WorldID
	Name    string
	Rooms   map[RoomID]*Room

	// Alone means that no other players will be seen here, will not support item uniqueness
	Alone bool
	// Instancable marks this world has one that is only ever run as an instance
	Instancable bool
	// Instance marks this world has an instance
	Instance bool
}

func NewWorld(id WorldID, instancable bool, instance, alone bool) *World {
	return &World{
		WorldID: id,
		Name:    "",
		Rooms:   make(map[RoomID]*Room),

		Alone:       alone,
		Instancable: instancable,
	}
}
